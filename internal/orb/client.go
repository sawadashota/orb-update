package orb

import (
	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/client"
	"github.com/pkg/errors"
)

// Client of CircleCI API
type Client interface {
	LatestVersion(o *Orb) (*Orb, error)
}

// CircleCIClient of CircleCI API
var CircleCIClient Client = NewDefaultClient()

// defaultClient of CircleCI API
type defaultClient struct {
	*client.Client
}

// NewDefaultClient returns defaultClient instance with default params
func NewDefaultClient() Client {
	return &defaultClient{
		Client: client.NewClient("https://circleci.com", "graphql-unstable", "", false),
	}

}

// LatestVersion of orb
func (c *defaultClient) LatestVersion(o *Orb) (*Orb, error) {
	version, err := api.OrbLatestVersion(c.Client, o.Namespace(), o.Name())
	if err != nil {
		return nil, errors.Errorf(`failed to fetch o's latest version  because "%s"`, err)
	}

	return New(o.Namespace(), o.Name(), version), nil
}
