package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/lnsp/koala/api/pkg/security"
	"github.com/sirupsen/logrus"

	"github.com/lnsp/koala/api/pkg/model"
)

type Config struct {
	Zonefile  string
	Origin    string
	TTL       int64
	ApplyCmd  []string
	JWTSecret string
	Htpasswd  string
	CORS      bool
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.corsEnabled {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, X-Requested-With, remember-me, Authorization")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
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
		applyCmd:    cfg.ApplyCmd,
		model:       dataModel,
		corsEnabled: cfg.CORS,
	}

	// Setup routing
	mux := http.NewServeMux()
	mux.HandleFunc("/list", handler.ListRecords)
	mux.HandleFunc("/apply", handler.ApplyRecords)
	mux.HandleFunc("/", handler.OK)
	handler.mux = mux

	// Inject security middleware
	if cfg.JWTSecret != "" {
		logrus.Info("Enabled JWT authentication")
		handler.mux = security.NewJWTGuard([]byte(cfg.JWTSecret), mux)
	}

	if cfg.CORS {
		logrus.Info("Enabled support for CORS")
	}

	return handler
}
