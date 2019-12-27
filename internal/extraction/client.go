package extraction

import (
	"github.com/CircleCI-Public/circleci-cli/api"
	"github.com/CircleCI-Public/circleci-cli/client"
	"github.com/pkg/errors"
	"github.com/sawadashota/orb-update/internal/orb"
)

// Client of CircleCI API
type Client interface {
	LatestVersion(o *orb.Orb) (*orb.Orb, error)
}

// DefaultClient of CircleCI API
var DefaultClient Client = NewDefaultClient()

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
func (c *defaultClient) LatestVersion(o *orb.Orb) (*orb.Orb, error) {
	version, err := api.OrbLatestVersion(c.Client, o.Namespace(), o.Name())
	if err != nil {
		return nil, errors.Errorf(`failed to fetch o's latest version  because "%s"`, err)
	}

	return orb.New(o.Namespace(), o.Name(), version), nil
}
