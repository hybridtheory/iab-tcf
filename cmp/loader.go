package cmp

import (
	"encoding/json"
	"net/http"

	"golang.org/x/exp/maps"
)

const (
	DefaultCMPVendorList = "https://cmplist.consensu.org/v2/cmp-list.json"
)

var (
	ValidCMPs []int
)

// Option is the type that allows us to configure the Loader dynamically.
type Option func(loader *Loader)

// Loader is the type that contains the logic to load and parse a CMP JSON list.
type Loader struct {
	URL string
}

// CMP contains the structure of the CMP info that comes inside the JSON
type CMP struct {
	ID           int
	Name         string
	IsCommercial bool
	Environments []string
}

// WithURL allows to configure a different URL for the CMP JSON list.
func WithURL(url string) Option {
	return func(cmp *Loader) {
		cmp.URL = url
	}
}

// NewLoader returns a CMP vendor list loader instance.
func NewLoader(options ...Option) *Loader {
	loader := &Loader{
		URL: DefaultCMPVendorList,
	}
	for _, option := range options {
		option(loader)
	}
	return loader
}

// Unmarshal parses the JSON vendor list into a struct so we can use them.
func (loader *Loader) Unmarshal(response *http.Response) ([]CMP, error) {
	type Response struct {
		CMPS map[string]CMP `json:"cmps"`
	}
	data := Response{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return []CMP{}, err
	}
	return maps.Values(data.CMPS), nil
}

// Load is used to load the vendor list into a list of CMP information.
func (loader *Loader) Load() ([]CMP, error) {
	response, err := http.Get(loader.URL)
	if err == nil {
		return loader.Unmarshal(response)
	}
	return []CMP{}, err
}

// LoadIDs loads the list of vendor CMP ids globally so we can reuse it
// with subsequent calls.
func (loader *Loader) LoadIDs() error {
	cmps, err := loader.Load()
	if err == nil {
		ValidCMPs = []int{}
		for _, cmp := range cmps {
			ValidCMPs = append(ValidCMPs, cmp.ID)
		}
	}
	return err
}
