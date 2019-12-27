package orb

import (
	"bytes"
	"io"
	"regexp"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

// ExtractionRegex .
var ExtractionRegex *regexp.Regexp

func init() {
	ExtractionRegex = regexp.MustCompile(`([^\s]+?/[^\s]+?)@[^\s].+`)
}

// Extraction orb instance
type Extraction struct {
	buf bytes.Buffer
}

// NewExtraction .
func NewExtraction(r io.Reader) (*Extraction, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, errors.Errorf(`failed to read config file because "%s"`, err)
	}

	return &Extraction{
		buf: buf,
	}, nil
}

// Reader returns io.Reader
func (e *Extraction) Reader() io.Reader {
	return bytes.NewReader(e.buf.Bytes())
}

// Orbs extract orbs from configuration file
func (e *Extraction) Orbs() ([]*Orb, error) {
	var mapConfig struct {
		Orbs map[string]string
	}
	if err := yaml.NewDecoder(e.Reader()).Decode(&mapConfig); err != nil {
		return nil, errors.Errorf(`failed to decode config file of orb-update because "%s"`, err)
	}

	var orbs []*Orb
	for _, orbStr := range mapConfig.Orbs {
		o, err := Parse(orbStr)
		if err != nil {
			return nil, err
		}

		orbs = append(orbs, o)
	}

	return orbs, nil
}
