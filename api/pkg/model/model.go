package model

import (
	"fmt"
	"os"
	"sync"

	"github.com/miekg/dns"
)

type Model struct {
	sync.Mutex
	Origin   string
	TTL      int64
	Zonefile string
	Header   []dns.RR
	Records  []Record
}

type Record struct {
	Type string
	Name string
	Data string
}

func (record Record) String() string {
	return fmt.Sprintf("%s IN %s %s", record.Name, record.Type, record.Data)
}

// FromZonefile initializes a new data model with a zonefile backend.
func FromZonefile(zonefile, origin string, ttl int64) (*Model, error) {
	model := &Model{
		Zonefile: zonefile,
		Origin:   origin,
		TTL:      ttl,
	}
	model.Refresh()
	return model, nil
}

func (model *Model) Update(records []Record) {
	model.Lock()
	model.Records = records
	model.Unlock()
}

func (model *Model) Flush() error {
	model.Lock()
	defer model.Unlock()
	// Open zonefile
	file, err := os.Create(model.Zonefile)
	if err != nil {
		return fmt.Errorf("flush zonefile: %w", err)
	}
	defer file.Close()
	// Write header
	for _, rr := range model.Header {
		fmt.Fprintln(file, rr.String())
	}
	// Write records
	directives := fmt.Sprintf("$ORIGIN %s\n$TTL %d\n", model.Origin, model.TTL)
	for _, record := range model.Records {
		rr, err := dns.NewRR(directives + record.String())
		if err != nil {
			return fmt.Errorf("convert record: %w", err)
		}
		fmt.Fprintln(file, rr.String())
	}
	return nil
}

// Refresh fetches the current state from disk.
func (model *Model) Refresh() error {
	model.Lock()
	defer model.Unlock()
	file, err := os.Open(model.Zonefile)
	if err != nil {
		return fmt.Errorf("open zonefile: %w", err)
	}
	defer file.Close()
	// Split zonefile into header and user-editable records
	model.Header = nil
	model.Records = nil
	parser := dns.NewZoneParser(file, "", "")
	for rr, ok := parser.Next(); ok; rr, ok = parser.Next() {
		// If type is A, AAAA or CNAME, add to records
		switch rr := rr.(type) {
		case *dns.A:
			model.Records = append(model.Records, Record{
				Type: "A",
				Name: rr.Hdr.Name,
				Data: rr.A.String(),
			})
		case *dns.AAAA:
			model.Records = append(model.Records, Record{
				Type: "AAAA",
				Name: rr.Hdr.Name,
				Data: rr.AAAA.String(),
			})
		case *dns.CNAME:
			model.Records = append(model.Records, Record{
				Type: "CNAME",
				Name: rr.Hdr.Name,
				Data: rr.Target,
			})
		default:
			model.Header = append(model.Header, rr)
		}
	}
	if err := parser.Err(); err != nil {
		return err
	}
	return nil
}
