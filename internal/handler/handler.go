package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/wpalmer/gozone"
)

type Config struct {
	StaticDir string
	Zonefile  string
	ApplyCmd  []string
}

type Handler struct {
	zonefile string
	applyCmd []string
	mux      *http.ServeMux
}

type dnsRecord struct {
	Type   string `json:"type"`
	Domain string `json:"domain"`
	Data   string `json:"data"`
}

func (h *Handler) ApplyRecords(w http.ResponseWriter, r *http.Request) {
	var changedRecords []dnsRecord
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&changedRecords); err != nil {
		http.Error(w, fmt.Sprintf("failed to read records: %v", err), http.StatusBadRequest)
		return
	}
	zoneRecords, err := h.readZonefile()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to apply records: %v", err), http.StatusInternalServerError)
		return
	}
	// remove all A records
	filteredZoneRecords := make([]gozone.Record, 0, len(zoneRecords))
	for i := range zoneRecords {
		if zoneRecords[i].Type != gozone.RecordType_A || zoneRecords[i].DomainName == "IN" {
			filteredZoneRecords = append(filteredZoneRecords, zoneRecords[i])
		}
	}
	// insert all A records
	for i := range changedRecords {
		if changedRecords[i].Type != "A" {
			continue
		}
		filteredZoneRecords = append(filteredZoneRecords, gozone.Record{
			DomainName: changedRecords[i].Domain,
			Data:       strings.Split(changedRecords[i].Data, " "),
			TimeToLive: -1,
			Class:      gozone.RecordClass_IN,
			Type:       gozone.RecordType_A,
		})
	}
	// write changes to file
	if err := h.writeZonefile(filteredZoneRecords); err != nil {
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

func (h *Handler) writeZonefile(records []gozone.Record) error {
	file, err := os.OpenFile(h.zonefile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return errors.Wrap(err, "failed to create zonefile")
	}
	defer file.Close()
	for _, rec := range records {
		if rec.DomainName == "@" {
			rev, _ := strconv.Atoi(rec.Data[3])
			rec.Data[3] = strconv.Itoa(rev + 1)
		} else if rec.DomainName == "IN" {
			rec.DomainName = "\t\tIN"
		}
		_, err := fmt.Fprintln(file, rec)
		if err != nil {
			return errors.Wrap(err, "failed to write record")
		}
	}
	return nil
}

func (h *Handler) readZonefile() ([]gozone.Record, error) {
	file, err := os.Open(h.zonefile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open zonefile")
	}
	defer file.Close()
	var (
		records []gozone.Record
		scanner = gozone.NewScanner(file)
	)
	for {
		var r gozone.Record
		if err := scanner.Next(&r); err != nil {
			break
		}
		records = append(records, r)
	}
	return records, nil
}

func (h *Handler) ListRecords(w http.ResponseWriter, r *http.Request) {
	zoneRecords, err := h.readZonefile()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list zones: %v", err), http.StatusInternalServerError)
		return
	}
	records := make([]dnsRecord, 0)
	for _, rec := range zoneRecords {
		if rec.Type != gozone.RecordType_A {
			continue
		} else if rec.DomainName == "IN" {
			continue
		}
		records = append(records, dnsRecord{
			Type:   "A",
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
	log.Printf("%s %s", r.Method, r.URL)
	h.mux.ServeHTTP(w, r)
}

func New(cfg Config) *Handler {
	handler := &Handler{
		mux:      http.NewServeMux(),
		zonefile: cfg.Zonefile,
		applyCmd: cfg.ApplyCmd,
	}
	handler.mux.Handle("/", http.FileServer(http.Dir(cfg.StaticDir)))
	handler.mux.HandleFunc("/api/list", handler.ListRecords)
	handler.mux.HandleFunc("/api/apply", handler.ApplyRecords)
	return handler
}
