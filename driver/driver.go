package driver

import "github.com/sawadashota/orb-update/driver/configuration"

// Driver .
type Driver interface {
	Registry() Registry
	Configuration() Configuration
}

type Configuration interface {
	configuration.Provider
}

// DefaultDriver .
type DefaultDriver struct {
	r Registry
	c configuration.Provider
}

// NewDefaultDriver .
func NewDefaultDriver() (Driver, error) {
	c := configuration.NewViperProvider()

	r, err := NewDefaultRegistry(c)
	if err != nil {
		return nil, err
	}

	return &DefaultDriver{
		r: r,
		c: c,
	}, nil
}

// Registry .
func (d *DefaultDriver) Registry() Registry {
	return d.r
}

// Configuration .
func (d *DefaultDriver) Configuration() Configuration {
	return d.c
}
