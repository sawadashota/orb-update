package driver

import (
	"io"

	"github.com/sawadashota/orb-update/internal/orb"

	"github.com/sawadashota/orb-update/handler"
	"github.com/sawadashota/orb-update/internal/filesystem"
	"github.com/sawadashota/orb-update/internal/git"
	"github.com/sawadashota/orb-update/internal/pullrequest"
)

// Registry .
type Registry interface {
	Logger() io.Writer

	Git() git.Git
	Filesystem() filesystem.Filesystem
	PullRequest() pullrequest.Creator
	VCSRepository() pullrequest.GitRepository

	Handler() *handler.Handler
	Extraction() (*orb.Extraction, error)
}
