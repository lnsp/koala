package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/lnsp/koala/internal/handler"
)

type Specification struct {
	Addr      string `default:":8080" description:"Address the server will be listening on"`
	Zonefile  string `required:"true" description:"Zonefile to be edited"`
	StaticDir string `default:"web/dist" description:"Directory of static assets"`
	ApplyCmd  string `default:"sleep 1" description:"Command executed after applying zonefile changes"`
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
	if err := srv.ListenAndServe(); err != nil {
		log.Println("failed to serve:", err)
	}
}
