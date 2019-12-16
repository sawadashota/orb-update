package orb

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

var orbFormatRegex *regexp.Regexp

func init() {
	orbFormatRegex = regexp.MustCompile(`([^\s]+?/[^\s]+?)@[^\s].+`)
}

type Config struct {
	Orbs []*Orb
}

type ConfigFile struct {
	bytes []byte
}

func NewConfigFile(r io.Reader) (*ConfigFile, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &ConfigFile{
		bytes: b,
	}, nil
}

func (cf *ConfigFile) reader() io.Reader {
	return bytes.NewReader(cf.bytes)
}

// Parse configuration file
func (cf *ConfigFile) Parse() (*Config, error) {
	var mapConfig struct {
		Orbs map[string]string
	}
	if err := yaml.NewDecoder(cf.reader()).Decode(&mapConfig); err != nil {
		return nil, err
	}

	var config Config
	for _, orb := range mapConfig.Orbs {
		o, err := ParseOrb(orb)
		if err != nil {
			return nil, err
		}

		config.Orbs = append(config.Orbs, o)
	}

	return &config, nil
}

func (cf *ConfigFile) Update(w io.Writer, diff *Difference) error {
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
