package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/lnsp/koala/api/pkg/router"
)

var version = "dev-build"

type Specification struct {
	Addr      string `default:":8080" description:"Address the server will be listening on"`
	Zonefile  string `required:"true" description:"Zonefile to be edited"`
	Origin    string `default:"."`
	TTL       int64  `default:"3600"`
	ApplyCmd  string `default:"sleep 1" description:"Command executed after applying zonefile changes"`
	JWTSecret string `default:"" description:"Auth secret for JWT tokens"`
	Debug     bool   `default:"false" description:"Enable debug logging" envconfig:"debug"`
	CORS      bool   `default:"false" description:"Enable support for CORS"`
}

func main() {
	var s Specification
	if err := envconfig.Process("koala", &s); err != nil {
		envconfig.Usage("koala", &s)
		return
	}
	if s.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	srv := &http.Server{
		Addr: s.Addr,
		Handler: router.New(router.Config{
			Zonefile:  s.Zonefile,
			Origin:    s.Origin,
			TTL:       s.TTL,
			ApplyCmd:  strings.Split(s.ApplyCmd, " "),
			JWTSecret: s.JWTSecret,
			CORS:      s.CORS,
		}),
	}
	logrus.WithFields(logrus.Fields{
		"version": version,
	}).Info("Up and running")
	if err := srv.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("Could not serve")
	}
}
