package handler

import (
	"io"

	"github.com/sawadashota/orb-update/driver/configuration"
	"github.com/sawadashota/orb-update/internal/filesystem"
	"github.com/sawadashota/orb-update/internal/git"
	"github.com/sawadashota/orb-update/internal/orb"
	"github.com/sawadashota/orb-update/internal/pullrequest"
)

// Registry for handler
type Registry interface {
	Logger() io.Writer
	Git() git.Git
	Filesystem() filesystem.Filesystem
	PullRequest() pullrequest.Creator
	CircleCIClient() orb.Client
}

// Configuration .
type Configuration interface {
	configuration.Provider
}
