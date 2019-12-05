package orb

import (
	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/client"
)

type Client interface {
	LatestVersion(orb *Orb) (*Orb, error)
}

type DefaultClient struct {
	*client.Client
}

func NewClient(host, endpoint string) *DefaultClient {
	return &DefaultClient{
		Client: client.NewClient(host, endpoint, "", false),
	}
}

func NewDefaultClient() *DefaultClient {
	return NewClient("https://circleci.com", "graphql-unstable")
}

func (c *DefaultClient) LatestVersion(orb *Orb) (*Orb, error) {
	version, err := api.OrbLatestVersion(c.Client, orb.namespace, orb.name)
	if err != nil {
		return nil, err
	}

	return NewOrb(orb.namespace, orb.name, version), nil
}
