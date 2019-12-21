package configfile

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/sawadashota/orb-update/internal/orb"
	"gopkg.in/yaml.v2"
)

var orbFormatRegex *regexp.Regexp

func init() {
	orbFormatRegex = regexp.MustCompile(`([^\s]+?/[^\s]+?)@[^\s].+`)
}

// ConfigFile of CircleCI
type ConfigFile struct {
	bytes []byte
}

// New .
func New(r io.Reader) (*ConfigFile, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Errorf(`failed to read config file because "%s"`, err)
	}

	return &ConfigFile{
		bytes: b,
	}, nil
}

func (cf *ConfigFile) reader() io.Reader {
	return bytes.NewReader(cf.bytes)
}

// Parse configuration file
func (cf *ConfigFile) Parse() ([]*orb.Orb, error) {
	var mapConfig struct {
		Orbs map[string]string
	}
	if err := yaml.NewDecoder(cf.reader()).Decode(&mapConfig); err != nil {
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

// Update writes updated orb version
func (cf *ConfigFile) Update(w io.Writer, diff *orb.Difference) error {
	var b bytes.Buffer

	scan := bufio.NewScanner(cf.reader())
	for scan.Scan() {
		func() {
			if strings.Contains(scan.Text(), diff.Old.String()) {
				b.WriteString(orbFormatRegex.ReplaceAllString(scan.Text(), "$1@"+diff.New.Version().String()))
				b.WriteString("\n")
				return
			}
			b.Write(scan.Bytes())
			b.WriteString("\n")
		}()
	}

	_, err := io.Copy(w, &b)
	return err
}
