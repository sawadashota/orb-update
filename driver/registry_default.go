package driver

import (
	"context"

	"github.com/sawadashota/orb-update/driver/configuration"
	"github.com/sawadashota/orb-update/handler"
	"github.com/sawadashota/orb-update/internal/filesystem"
	"github.com/sawadashota/orb-update/internal/git"
	"github.com/sawadashota/orb-update/internal/orb"
	"github.com/sawadashota/orb-update/internal/pullrequest"
)

// DefaultRegistry .
type DefaultRegistry struct {
	c    configuration.Provider
	g    git.Git
	repo *git.Repository
	fs   filesystem.Filesystem
	pr   pullrequest.Creator
	cl   orb.Client

	h *handler.Handler
}

// Check satisfy Registry interface
var _ Registry = new(DefaultRegistry)

// NewDefaultRegistry .
func NewDefaultRegistry(c configuration.Provider) (*DefaultRegistry, error) {
	dr := &DefaultRegistry{
		c:  c,
		cl: orb.NewDefaultClient(),
	}

	if err := dr.setupRepository(); err != nil {
		return nil, err
	}

	pr, err := pullrequest.NewGitHubPullRequest(context.Background(), dr, c)
	if err != nil {
		return nil, err
	}
	dr.pr = pr

	return dr, nil
}

func (d *DefaultRegistry) setupRepository() error {
	if d.c.FilesystemStrategy() == configuration.OsFileSystemStrategy {
		g, fs, err := git.OpenCurrentDirectoryRepository(d.c)
		if err != nil {
			return err
		}

		d.g = g
		d.fs = fs
		return nil
	}

	repo, err := git.ParseRepository(d.c.RepositoryName())
	if err != nil {
		return err
	}

	g, fs, err := git.Clone(d.c, repo.Owner(), repo.Name())
	if err != nil {
		return err
	}

	d.repo = repo
	d.g = g
	d.fs = fs

	return nil
}

// VCSRepository .
func (d *DefaultRegistry) VCSRepository() pullrequest.GitRepository {
	return d.repo
}

// Git .
func (d *DefaultRegistry) Git() git.Git {
	return d.g
}

// Filesystem .
func (d *DefaultRegistry) Filesystem() filesystem.Filesystem {
	return d.fs
}

// PullRequest .
func (d *DefaultRegistry) PullRequest() pullrequest.Creator {
	return d.pr
}

// CircleCIClient .
func (d *DefaultRegistry) CircleCIClient() orb.Client {
	return d.cl
}

// Handler .
func (d *DefaultRegistry) Handler() *handler.Handler {
	if d.h == nil {
		d.h = handler.New(d, d.c)
	}

	return d.h
}
