package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/lnsp/koala/internal/model"
	"github.com/wpalmer/gozone"
)

var (
	RecordA            = "A"
	RecordCNAME        = "CNAME"
	AllowedRecordTypes = map[string]gozone.RecordType{
		RecordA:     gozone.RecordType_A,
		RecordCNAME: gozone.RecordType_CNAME,
	}
)

type Config struct {
	StaticDir string
	Zonefile  string
	ApplyCmd  []string
}

type Handler struct {
	model    *model.Model
	applyCmd []string
	mux      *http.ServeMux
}

type dnsRecord struct {
	Type   string `json:"type"`
	Domain string `json:"domain"`
	Data   string `json:"data"`
}

func IsAllowedType(t string) bool {
	for at, _ := range AllowedRecordTypes {
		if at == t {
			return true
		}
	}
	return false
}

func IsAllowedTypeID(t gozone.RecordType) bool {
	for _, at := range AllowedRecordTypes {
		if at == t {
			return true
		}
	}
	return false
}

func GetRecordTypeDesc(t gozone.RecordType) string {
	for s, at := range AllowedRecordTypes {
		if at == t {
			return s
		}
	}
	return ""
}

// ApplyRecords reads the A-records from the request body,
// reads all records from the zonefile, removes all A-records from the zonefile
// and inserts the A-records from the request body.
func (h *Handler) ApplyRecords(w http.ResponseWriter, r *http.Request) {
	var changedRecords []dnsRecord
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&changedRecords); err != nil {
		http.Error(w, fmt.Sprintf("failed to read records: %v", err), http.StatusBadRequest)
		return
	}
	zoneRecords, err := h.model.ReadAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to apply records: %v", err), http.StatusInternalServerError)
		return
	}
	// remove all A records
	filteredZoneRecords := make([]gozone.Record, 0, len(zoneRecords))
	for i := range zoneRecords {
		if !IsAllowedTypeID(zoneRecords[i].Type) || zoneRecords[i].DomainName == "IN" {
			filteredZoneRecords = append(filteredZoneRecords, zoneRecords[i])
		}
	}
	// insert all A records
	for _, rec := range changedRecords {
		if !IsAllowedType(rec.Type) {
			continue
		}
		filteredZoneRecords = append(filteredZoneRecords, gozone.Record{
			DomainName: rec.Domain,
			Data:       strings.Split(rec.Data, " "),
			TimeToLive: -1,
			Class:      gozone.RecordClass_IN,
			Type:       AllowedRecordTypes[rec.Type],
		})
		log.Printf("Inserting record '%s' (%s) -> %s", rec.Domain, rec.Type, rec.Data)
	}
	// write changes to file
	if err := h.model.Update(filteredZoneRecords); err != nil {
		http.Error(w, fmt.Sprintf("failed to apply records: %v", err), http.StatusInternalServerError)
		return
	}
	// run post-update command
	if len(h.applyCmd) != 0 {
		cmd := exec.Command(h.applyCmd[0], h.applyCmd[1:]...)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to start post-apply cmd: %v", err), http.StatusInternalServerError)
			return
		}
		if err := cmd.Run(); err != nil {
			http.Error(w, fmt.Sprintf("failed to execute post-apply cmd: %v", err), http.StatusInternalServerError)
		}
	}
	w.Write([]byte("OK"))
}

func (h *Handler) ListRecords(w http.ResponseWriter, r *http.Request) {
	zoneRecords, err := h.model.ReadAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list zones: %v", err), http.StatusInternalServerError)
		return
	}
	records := make([]dnsRecord, 0)
	for _, rec := range zoneRecords {
		if rec.DomainName == "IN" {
			continue
		} else if GetRecordTypeDesc(rec.Type) == "" {
			continue
		}
		records = append(records, dnsRecord{
			Type:   GetRecordTypeDesc(rec.Type),
			Domain: rec.DomainName,
			Data:   strings.Join(rec.Data, " "),
		})
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&records); err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	h.mux.ServeHTTP(w, r)
	elapsed := time.Since(start)
	log.Printf("%s %s %.3fs", r.Method, r.URL, elapsed.Seconds())
}

func New(cfg Config) *Handler {
	model, err := model.FromZonefile(cfg.Zonefile)
	if err != nil {
		log.Fatalf("failed to create model: %v", err)
	}
	handler := &Handler{
		mux:      http.NewServeMux(),
		applyCmd: cfg.ApplyCmd,
		model:    model,
	}
	handler.mux.Handle("/", http.FileServer(http.Dir(cfg.StaticDir)))
	handler.mux.HandleFunc("/api/list", handler.ListRecords)
	handler.mux.HandleFunc("/api/apply", handler.ApplyRecords)
	return handler
}
