package main

import (
	"embed"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kelseyhightower/envconfig"
	"github.com/lnsp/koala/pkg/router"
	"github.com/lnsp/koala/pkg/security"
)

var version = "dev"

// go:embed webui/dist
var staticFiles embed.FS

type spec struct {
	Addr     string `default:":8080" desc:"Address the server will be listening on"`
	Zonefile string `required:"true" desc:"Zonefile to be edited"`
	Origin   string `default:"." desc:"Zone to be edited"`
	TTL      int64  `default:"3600" desc:"Default TTL for records"`
	ApplyCmd string `default:"sleep 1" desc:"Command executed after applying zonefile changes"`

	Debug bool `default:"false" desc:"Enable debug logging" envconfig:"debug"`
	CORS  bool `default:"false" desc:"Enable support for CORS"`

	Security           string `default:"none" desc:"Security guard to use [none|oidc|jwt]"`
	OIDCClientID       string `default:"" desc:"OpenID Connect Client ID"`
	OIDCIdentityServer string `default:"" desc:"URL of identity provider"`
	JWTSecret          string `default:"" desc:"Auth secret for JWT tokens"`
}

func main() {
	var s spec
	if err := envconfig.Process("koala", &s); err != nil {
		envconfig.Usage("koala", &s)
		return
	}
	if s.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	// Set up security guard
	var guard security.Guard
	switch s.Security {
	case "none":
		guard = security.None()
	case "jwt":
		guard = security.JWT(s.JWTSecret)
	case "oidc":
		guard = security.OIDC(s.OIDCClientID, s.OIDCIdentityServer)
	default:
		logrus.Fatal("unknown security guard:", s.Security)
	}
	srv := &http.Server{
		Addr: s.Addr,
		Handler: router.New(router.Config{
			Zonefile: s.Zonefile,
			Origin:   s.Origin,
			TTL:      s.TTL,
			ApplyCmd: strings.Split(s.ApplyCmd, " "),
			CORS:     s.CORS,
			Security: guard,
		}),
	}
	logrus.WithFields(logrus.Fields{
		"version": version,
	}).Info("Up and running")
	if err := srv.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("Could not serve")
	}
}
