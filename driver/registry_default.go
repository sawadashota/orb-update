package driver

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sawadashota/orb-update/driver/configuration"
	"github.com/sawadashota/orb-update/handler"
	"github.com/sawadashota/orb-update/internal/extraction"
	"github.com/sawadashota/orb-update/internal/filesystem"
	"github.com/sawadashota/orb-update/internal/git"
	"github.com/sawadashota/orb-update/internal/pullrequest"
	"github.com/sawadashota/orb-update/internal/vcsuser"
	"github.com/spf13/viper"
)

// Logger  of default
var Logger io.Writer = os.Stdout

// DefaultRegistry .
type DefaultRegistry struct {
	l    io.Writer
	c    configuration.Provider
	g    git.Git
	repo *git.Repository
	fs   filesystem.Filesystem
	pr   pullrequest.Creator

	h *handler.Handler
}

// Check satisfy Registry interface
var _ Registry = new(DefaultRegistry)

// NewDefaultRegistry .
func NewDefaultRegistry(c configuration.Provider) (*DefaultRegistry, error) {
	dr := &DefaultRegistry{
		l: Logger,
		c: c,
	}

	if c.GitAuthorEmail() == "" || c.GitAuthorName() == "" {
		if err := dr.setGitAuthorFromVCS(context.Background()); err != nil {
			return nil, err
		}
	}

	if err := dr.setupRepository(); err != nil {
		return nil, err
	}

	dr.pr = pullrequest.NewGitHubPullRequest(context.Background(), dr, c)

	return dr, nil
}

// Logger .
func (d *DefaultRegistry) Logger() io.Writer {
	return d.l
}

// setGitAuthorFromVCS when git author is not configured
func (d *DefaultRegistry) setGitAuthorFromVCS(ctx context.Context) error {
	cl := vcsuser.NewGithubClient(ctx, d.c)
	user, err := cl.Fetch(ctx)
	if err != nil {
		return err
	}

	viper.Set(configuration.ViperGitAuthorName, user.Name())
	viper.Set(configuration.ViperGitAuthorEmail, user.Email())

	_, _ = fmt.Fprintf(d.l, "git author has set GitHub name(%s) and email(%s)\n", user.Name(), user.Email())

	return nil
}

func (d *DefaultRegistry) setupRepository() error {
	if d.c.FilesystemStrategy() == configuration.OsFileSystemStrategy {
		return d.setupOsFileSystem()
	}

	return d.setupInMemoryFilesystem()
}

func (d *DefaultRegistry) setupOsFileSystem() error {
	g, fs, err := git.OpenCurrentDirectoryRepository(d.c)
	if err != nil {
		return err
	}

	d.g = g
	d.fs = fs
	return nil
}

func (d *DefaultRegistry) setupInMemoryFilesystem() error {
	repo, err := git.ParseRepository(d.c.RepositoryName())
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(d.l, "cloning %s ...\n", repo)

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

// Handler .
func (d *DefaultRegistry) Handler() *handler.Handler {
	if d.h == nil {
		opts := make([]handler.Option, 0)
		if d.c.GitHubPullRequest() {
			opts = append(opts, handler.WithPullRequestCreation())
		}

		d.h = handler.New(d, d.c, opts...)
	}

	return d.h
}

// Extraction for orb instance
func (d *DefaultRegistry) Extraction() (*extraction.Extraction, error) {
	reader, err := d.targetFilesReader()
	if err != nil {
		return nil, err
	}
	e, err := extraction.NewExtraction(reader)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// targetFilesReader is bundled io.Reader
func (d *DefaultRegistry) targetFilesReader() (io.Reader, error) {
	files := make([]io.Reader, len(d.c.TargetFiles()))

	var err error
	for i, target := range d.c.TargetFiles() {
		files[i], err = d.Filesystem().Reader(target)
		if err != nil {
			return nil, err
		}
	}

	return io.MultiReader(files...), nil
}
