package orb

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var semanticVersionRegex *regexp.Regexp

func init() {
	semanticVersionRegex = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+`)
}

// Orb .
type Orb struct {
	namespace string
	name      string
	version   Version
}

// Version type
type Version string

// String .
func (v Version) String() string {
	return string(v)
}

// IsSemantic version or not
func (v Version) IsSemantic() bool {
	return semanticVersionRegex.MatchString(v.String())
}

// New orb instance.
func New(namespace, name, version string) *Orb {
	return &Orb{
		namespace: namespace,
		name:      name,
		version:   Version(version),
	}
}

// Parse format namespace/name@version
func Parse(orb string) (*Orb, error) {
	splitVersion := strings.Split(orb, "@")
	if len(splitVersion) != 2 {
		return nil, errors.Errorf("Incorrect orb format: %s", orb)
	}

	splitName := strings.Split(splitVersion[0], "/")
	if len(splitName) != 2 {
		return nil, errors.Errorf("Incorrect orb format: %s", orb)
	}

	return New(splitName[0], splitName[1], splitVersion[1]), nil
}

// Namespace .
func (o *Orb) Namespace() string {
	return o.namespace
}

// Name .
func (o *Orb) Name() string {
	return o.name
}

// Version .
func (o *Orb) Version() Version {
	return o.version
}

// String .
func (o *Orb) String() string {
	return fmt.Sprintf("%s/%s@%s", o.namespace, o.name, o.version)
}

// IsSameOrb or not
func (o *Orb) IsSameOrb(compare *Orb) bool {
	return o.Namespace() == compare.Namespace() && o.Name() == compare.Name()
}
