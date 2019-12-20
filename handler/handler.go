package handler

import (
	"io"
	"os"
)

// Handler .
type Handler struct {
	r                     Registry
	c                     Configuration
	logger                io.Writer
	doesCreatePullRequest bool
}

// Option for Handler
type Option func(handler *Handler)

// WithPullRequestCreation .
func WithPullRequestCreation() Option {
	return func(handler *Handler) {
		handler.doesCreatePullRequest = true
	}
}

// WithLogger replace default logger
func WithLogger(w io.Writer) Option {
	return func(handler *Handler) {
		handler.logger = w
	}
}

// New Handler instance
func New(r Registry, c Configuration, opts ...Option) *Handler {
	h := &Handler{
		r:                     r,
		c:                     c,
		logger:                os.Stdout,
		doesCreatePullRequest: false,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}
