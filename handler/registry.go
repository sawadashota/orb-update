package handler

import (
	"github.com/sawadashota/orb-update/driver/configuration"
	"github.com/sawadashota/orb-update/filesystem"
	"github.com/sawadashota/orb-update/git"
	"github.com/sawadashota/orb-update/orb"
	"github.com/sawadashota/orb-update/pullrequest"
)

// Registry for handler
type Registry interface {
	Git() git.Git
	Filesystem() filesystem.Filesystem
	PullRequest() pullrequest.Creator
	CircleCIClient() orb.Client
}

// Configuration .
type Configuration interface {
	configuration.Provider
}
