package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/lnsp/koala/internal/handler"
)

var version = "dev-build"

type Specification struct {
	Addr        string `default:":8080" description:"Address the server will be listening on"`
	Zonefile    string `required:"true" description:"Zonefile to be edited"`
	StaticDir   string `default:"web/dist" description:"Directory of static assets"`
	ApplyCmd    string `default:"sleep 1" description:"Command executed after applying zonefile changes"`
	Certificate string `default:"" description:"Certificate file for HTTPS"`
	PrivateKey  string `default:"" description:"Private key file for HTTPS"`
}

func main() {
	var s Specification
	if err := envconfig.Process("koala", &s); err != nil {
		envconfig.Usage("koala", &s)
		return
	}
	srv := &http.Server{
		Addr: s.Addr,
		Handler: handler.New(handler.Config{
			StaticDir: s.StaticDir,
			Zonefile:  s.Zonefile,
			ApplyCmd:  strings.Split(s.ApplyCmd, " "),
		}),
	}
	log.Printf("koala %s ready to serve!", version)
	if s.Certificate != "" {
		log.Printf("enabled HTTPS using cert '%s' and key '%s'", s.Certificate, s.PrivateKey)
		if err := srv.ListenAndServeTLS(s.Certificate, s.PrivateKey); err != nil {
			log.Println("failed to serve:", err)
		}
	} else {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("failed to serve:", err)
		}
	}
}
