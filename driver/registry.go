package driver

import (
	"github.com/sawadashota/orb-update/filesystem"
	"github.com/sawadashota/orb-update/git"
	"github.com/sawadashota/orb-update/handler"
	"github.com/sawadashota/orb-update/orb"
	"github.com/sawadashota/orb-update/pullrequest"
)

// Registry .
type Registry interface {
	Git() git.Git
	Filesystem() filesystem.Filesystem
	PullRequest() pullrequest.Creator
	CircleCIClient() orb.Client
	VCSRepository() pullrequest.GitRepository

	Handler() *handler.Handler
}
