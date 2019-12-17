package driver

import "github.com/sawadashota/orb-update/driver/configuration"

// Driver .
type Driver interface {
	Configuration() configuration.Provider
}

// DefaultDriver .
type DefaultDriver struct {
	c configuration.Provider
}

// NewDefaultDriver .
func NewDefaultDriver() Driver {
	c := configuration.NewViperProvider()

	return &DefaultDriver{
		c: c,
	}
}

// Configuration .
func (d *DefaultDriver) Configuration() configuration.Provider {
	return d.c
}
