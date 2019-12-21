package orb

import (
	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/client"
	"github.com/pkg/errors"
)

// Client of CircleCI API
type Client interface {
	LatestVersion(orb *Orb) (*Orb, error)
}

// DefaultClient of CircleCI API
type DefaultClient struct {
	*client.Client
}

// NewClient returns DefaultClient instance
func NewClient(host, endpoint string) *DefaultClient {
	return &DefaultClient{
		Client: client.NewClient(host, endpoint, "", false),
	}
}

// NewDefaultClient returns DefaultClient instance with default params
func NewDefaultClient() *DefaultClient {
	return NewClient("https://circleci.com", "graphql-unstable")
}

// LatestVersion of orb
func (c *DefaultClient) LatestVersion(orb *Orb) (*Orb, error) {
	version, err := api.OrbLatestVersion(c.Client, orb.namespace, orb.name)
	if err != nil {
		return nil, errors.Errorf(`failed to fetch orb's latest version  because "%s"`, err)
	}

	return New(orb.namespace, orb.name, version), nil
}
