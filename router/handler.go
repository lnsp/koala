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
	Zones    []Zone
	ApplyCmd []string
	CORS     bool
	Security security.Guard
	UI       fs.FS
	APIRoot  string
}

type Zone struct {
	Name   string
	Path   string
	Origin string
	TTL    int64
}

type Handler struct {
	zones       []string
	models      map[string]*model.Model
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
func (h *Handler) ApplyRecords(w http.ResponseWriter, r *http.Request, mod *model.Model) {
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
	mod.Update(converted)
	if err := mod.Flush(); err != nil {
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

type ModelHandler func(w http.ResponseWriter, r *http.Request, mod *model.Model)

func (h *Handler) MatchModel(mh ModelHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		zone := params.Get("zone")

		model, ok := h.models[zone]
		if !ok {
			http.Error(w, "zone not found", http.StatusNotFound)
			return
		}
		mh(w, r, model)
	})
}

func (h *Handler) ListRecords(w http.ResponseWriter, r *http.Request, mod *model.Model) {
	mod.Refresh()
	converted := make([]dnsRecord, len(mod.Records))
	for i, rr := range mod.Records {
		converted[i] = dnsRecord(rr)
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&converted); err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListZones(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.zones)
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

func New(cfg Config) (*Handler, error) {
	zones := make([]string, len(cfg.Zones))
	models := make(map[string]*model.Model)
	for i, zone := range cfg.Zones {
		dataModel, err := model.FromZonefile(zone.Path, zone.Origin, zone.TTL)
		if err != nil {
			return nil, err
		}
		models[zone.Name] = dataModel
		zones[i] = zone.Name
	}
	handler := &Handler{
		zones:       zones,
		models:      models,
		applyCmd:    cfg.ApplyCmd,
		corsEnabled: cfg.CORS,
	}

	// Setup routing

	rtr := mux.NewRouter()

	apiMux := rtr.PathPrefix(cfg.APIRoot).Subrouter()
	apiMux.Use(mux.MiddlewareFunc(cfg.Security))
	apiMux.HandleFunc("/", handler.ListZones)
	apiMux.Handle("/list", handler.MatchModel(handler.ListRecords))
	apiMux.Handle("/apply", handler.MatchModel(handler.ApplyRecords))

	rtr.PathPrefix("/").Handler(http.FileServer(http.FS(cfg.UI)))

	handler.mux = rtr

	if cfg.CORS {
		logrus.Info("enabled support for CORS")
	}

	return handler, nil
}
