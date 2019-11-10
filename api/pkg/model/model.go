package model

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wpalmer/gozone"
)

type Model struct {
	Zonefile string
}

// FromZonefile initializes a new data model with a zonefile backend.
func FromZonefile(zonefile string) (*Model, error) {
	return &Model{
		Zonefile: zonefile,
	}, nil
}

// ReadAll reads all records from the zonefile model.
func (model *Model) ReadAll() ([]gozone.Record, error) {
	file, err := os.Open(model.Zonefile)
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
			logrus.WithError(err).Warn("Could not parse zonefile entry")
			break
		}
		records = append(records, r)
	}
	return records, nil

}

// Update updates all records in the zonefile model.
func (model *Model) Update(records []gozone.Record) error {
	file, err := os.OpenFile(model.Zonefile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
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
