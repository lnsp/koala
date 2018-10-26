package main

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/lnsp/koala/api/pkg/router"
)

var version = "dev-build"

type Specification struct {
	Addr        string `default:":8080" description:"Address the server will be listening on"`
	Zonefile    string `required:"true" description:"Zonefile to be edited"`
	ApplyCmd    string `default:"sleep 1" description:"Command executed after applying zonefile changes"`
	Certificate string `default:"" description:"Certificate file for HTTPS"`
	PrivateKey  string `default:"" description:"Private key file for HTTPS"`
	JWTSecret string `default:"" description:"Auth secret for JWT tokens"`
	Debug bool `default:"false" description:"Enable debug logging" envconfig:"debug"`
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
			Zonefile: s.Zonefile,
			ApplyCmd: strings.Split(s.ApplyCmd, " "),
			JWTSecret: s.JWTSecret,
		}),
	}
	logrus.WithFields(logrus.Fields{
		"version": version,
	}).Info("Up and running")
	if s.Certificate != "" {
		logrus.WithFields(logrus.Fields{
			"cert": s.Certificate,
			"key": s.PrivateKey,
		}).Info("Enabled HTTPS")
		if err := srv.ListenAndServeTLS(s.Certificate, s.PrivateKey); err != nil {
			logrus.WithError(err).Fatal("Could not serve")
		}
	} else {
		if err := srv.ListenAndServe(); err != nil {
			logrus.WithError(err).Fatal("Could not serve")
		}
	}
}
