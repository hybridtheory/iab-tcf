package cmp

import (
	"encoding/json"
	"io"
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
	URL  string
	JSON string
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

// WithJSON allows to pass a JSON string loaded externally so our loader
// will parse the cmp-list JSON format for you.
func WithJSON(json string) Option {
	return func(cmp *Loader) {
		cmp.JSON = json
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
func (loader *Loader) Unmarshal(data []byte) ([]CMP, error) {
	type Response struct {
		CMPS map[string]CMP `json:"cmps"`
	}
	response := Response{}
	if err := json.Unmarshal(data, &response); err != nil {
		return []CMP{}, err
	}
	return maps.Values(response.CMPS), nil
}

// LoadHTTP is used to load the vendor list from a HTTP url.
func (loader *Loader) LoadHTTP() ([]CMP, error) {
	response, err := http.Get(loader.URL)
	if err == nil {
		if data, err := io.ReadAll(response.Body); err == nil {
			return loader.Unmarshal(data)
		}
	}
	return []CMP{}, err
}

// LoadJSON is used to load the vendor list from a received JSON string.
func (loader *Loader) LoadJSON() ([]CMP, error) {
	return loader.Unmarshal([]byte(loader.JSON))
}

// Load decides which CMP list we are going to load.
func (loader *Loader) Load() ([]CMP, error) {
	if loader.JSON != "" {
		return loader.LoadJSON()
	}
	return loader.LoadHTTP()
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
