package orb

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Orb struct {
	namespace string
	name      string
	version   string
}

func NewOrb(namespace, name, version string) *Orb {
	return &Orb{
		namespace: namespace,
		name:      name,
		version:   version,
	}
}

// ParseOrb format namespace/name@version
func ParseOrb(orb string) (*Orb, error) {
	splitVersion := strings.Split(orb, "@")
	if len(splitVersion) != 2 {
		return nil, errors.Errorf("Incorrect orb format: %s", orb)
	}

	splitName := strings.Split(splitVersion[0], "/")
	if len(splitName) != 2 {
		return nil, errors.Errorf("Incorrect orb format: %s", orb)
	}

	return NewOrb(splitName[0], splitName[1], splitVersion[1]), nil
}

func (o *Orb) Namespace() string {
	return o.namespace
}

func (o *Orb) Name() string {
	return o.name
}

func (o *Orb) Version() string {
	return o.version
}

func (o *Orb) String() string {
	return fmt.Sprintf("%s/%s@%s", o.namespace, o.name, o.version)
}
