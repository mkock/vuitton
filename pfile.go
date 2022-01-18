package vuitton

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// psMax is the maximum number of products we are willing to track.
const psMax = 10

var errMaxExceeded = errors.New("too many products being tracked")

// PFileModifiedSince returns true if the PFile has been modified since the given time.
func (m *MainLoop) PFileModifiedSince(when time.Time) bool {
	info, err := os.Stat(m.PFileName)
	if err != nil {
		return true // Proceeds with reading the file, which is okay. We'll handle file read issues at that point.
	}
	return info.ModTime().After(when)
}

// ReadPFile reads the "P" file (products) file from disk and converts the URL's into a slice of products.
// ReadPFile returns an error if it was unable to read the file. If the file was empty, ReadPFile just returns an empty
// slice and a nil error.
func (m *MainLoop) ReadPFile() ([]product, error) {
	if m.PFileName == "" {
		return []product{}, nil
	}
	f, err := os.Open(m.PFileName)
	if err != nil {
		return []product{}, err
	}
	defer func() { _ = f.Close() }()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return []product{}, err
	}
	lines := strings.Split(string(bytes), "\n")

	ps := make([]product, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		ps = append(ps, product{URL: l})
	}

	if len(ps) > psMax {
		return []product{}, errMaxExceeded
	}

	return ps, nil
}
