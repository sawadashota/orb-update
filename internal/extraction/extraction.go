package extraction

import (
	"bytes"
	"io"
	"regexp"

	"github.com/pkg/errors"

	"github.com/sawadashota/orb-update/internal/orb"
	"gopkg.in/yaml.v2"
)

// OrbFormatRegex .
var OrbFormatRegex *regexp.Regexp

func init() {
	OrbFormatRegex = regexp.MustCompile(`([^\s]+?/[^\s]+?)@[^\s].+`)
}

// Extraction orb instance
type Extraction struct {
	buf bytes.Buffer
}

// New Extraction .
func New(r io.Reader) (*Extraction, error) {
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
func (e *Extraction) Orbs() ([]*orb.Orb, error) {
	var mapConfig struct {
		Orbs map[string]string
	}
	if err := yaml.NewDecoder(e.Reader()).Decode(&mapConfig); err != nil {
		return nil, errors.Errorf(`failed to decode config file of orb-update because "%s"`, err)
	}

	var orbs []*orb.Orb
	for _, orbStr := range mapConfig.Orbs {
		o, err := orb.Parse(orbStr)
		if err != nil {
			return nil, err
		}

		orbs = append(orbs, o)
	}

	return orbs, nil
}
