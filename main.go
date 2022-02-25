package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/sirupsen/logrus"

	"github.com/lnsp/koala/router"
	"github.com/lnsp/koala/security"
	"github.com/lnsp/koala/webui"
)

var version = "dev"

type Configuration struct {
	Addr     string
	Debug    bool
	CORS     bool
	ApplyCmd string `toml:"apply_cmd"`
	APIRoot  string `toml:"api_root"`

	Zones []struct {
		Name   string
		Path   string
		Origin string
		TTL    int64
	}

	Security struct {
		Mode string
		OIDC struct {
			ClientID       string `toml:"client_id"`
			IdentityServer string `toml:"identity_server"`
		}
		JWT struct {
			Secret string
		}
	}
}

var (
	configPath = flag.String("config", "config.toml", "path to configuration")
)

func main() {
	flag.Parse()

	// Read config data
	configData, err := os.ReadFile(*configPath)
	if err != nil {
		logrus.Fatal("read config file:", err)
	}
	// ... and unmarshal
	var s Configuration
	if err := toml.Unmarshal(configData, &s); err != nil {
		logrus.Fatal("unmarshal config:", err)
	}
	// Configure debug level
	if s.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugf("with configuration %+v", s)
	}
	// Set up security guard
	var guard security.Guard
	switch s.Security.Mode {
	case "none":
		guard = security.None()
	case "jwt":
		guard = security.JWT(s.Security.JWT.Secret)
	case "oidc":
		guard = security.OIDC(s.Security.OIDC.ClientID, s.Security.OIDC.IdentityServer)
	default:
		logrus.Fatal("unknown security guard:", s.Security)
	}
	// Generate zone list
	routerZones := make([]router.Zone, len(s.Zones))
	for i, zone := range s.Zones {
		routerZones[i] = router.Zone(zone)
	}
	rtr, err := router.New(router.Config{
		Zones:    routerZones,
		ApplyCmd: strings.Split(s.ApplyCmd, " "),
		CORS:     s.CORS,
		Security: guard,
		UI:       webui.FS,
		APIRoot:  s.APIRoot,
	})
	if err != nil {
		logrus.Fatal("setup router:", err)
	}
	srv := &http.Server{
		Addr:         s.Addr,
		Handler:      rtr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logrus.WithFields(logrus.Fields{
		"version": version,
	}).Info("up and running")
	if err := srv.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("could not serve")
	}
}
