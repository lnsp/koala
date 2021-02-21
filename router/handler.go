package router

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/mux"

	"github.com/lnsp/koala/model"
	"github.com/lnsp/koala/security"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Zonefile string
	Origin   string
	TTL      int64
	ApplyCmd []string
	CORS     bool
	Security security.Guard
	UI       fs.FS
	APIRoot  string
}

type Handler struct {
	model       *model.Model
	applyCmd    []string
	mux         http.Handler
	corsEnabled bool
}

type dnsRecord struct {
	Type string `json:"type"`
	Name string `json:"domain"`
	Data string `json:"data"`
}

// ApplyRecords reads the records from the request body,
// reads all records from the zonefile, removes all records from the zonefile
// and inserts the A-records from the request body.
func (h *Handler) ApplyRecords(w http.ResponseWriter, r *http.Request) {
	var records []dnsRecord
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&records); err != nil {
		http.Error(w, fmt.Sprintf("read records: %v", err), http.StatusBadRequest)
		return
	}
	converted := make([]model.Record, len(records))
	for i, rr := range records {
		converted[i] = model.Record(rr)
	}
	h.model.Update(converted)
	if err := h.model.Flush(); err != nil {
		http.Error(w, fmt.Sprintf("apply records: %v", err), http.StatusInternalServerError)
		return
	}
	// run post-update command
	if len(h.applyCmd) != 0 {
		cmd := exec.Command(h.applyCmd[0], h.applyCmd[1:]...)
		if err := cmd.Run(); err != nil {
			http.Error(w, fmt.Sprintf("execute post-apply cmd: %v", err), http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("OK"))
}

func (h *Handler) ListRecords(w http.ResponseWriter, r *http.Request) {
	h.model.Refresh()
	converted := make([]dnsRecord, len(h.model.Records))
	for i, rr := range h.model.Records {
		converted[i] = dnsRecord(rr)
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&converted); err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) OK(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.corsEnabled {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, X-Requested-With, remember-me, Authorization")
		w.Header().Add("Access-Control-Allow-Methods", "GET, HEAD, POST, OPTIONS")
	}
	if r.Method == "OPTIONS" {
		w.Header().Set("Allow", "GET, HEAD, POST, OPTIONS")
		return
	}

	start := time.Now()
	h.mux.ServeHTTP(w, r)
	elapsed := time.Since(start)
	logrus.WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.URL,
		"time":   elapsed.Seconds(),
	}).Debug("HTTP Request")
}

func New(cfg Config) *Handler {
	dataModel, err := model.FromZonefile(cfg.Zonefile, cfg.Origin, cfg.TTL)
	if err != nil {
		logrus.WithError(err).Fatal("Could not create model")
	}
	handler := &Handler{
		model:       dataModel,
		applyCmd:    cfg.ApplyCmd,
		corsEnabled: cfg.CORS,
	}

	// Setup routing

	rtr := mux.NewRouter()

	apiMux := rtr.PathPrefix(cfg.APIRoot).Subrouter()
	apiMux.Use(mux.MiddlewareFunc(cfg.Security))
	apiMux.HandleFunc("/", handler.OK)
	apiMux.HandleFunc("/list", handler.ListRecords)
	apiMux.HandleFunc("/apply", handler.ApplyRecords)

	rtr.PathPrefix("/").Handler(http.FileServer(http.FS(cfg.UI)))

	handler.mux = rtr

	if cfg.CORS {
		logrus.Info("Enabled support for CORS")
	}

	return handler
}
