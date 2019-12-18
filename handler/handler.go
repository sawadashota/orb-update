package handler

import "io"

// Handler .
type Handler struct {
	r      Registry
	c      Configuration
	logger io.Writer

	doesCreatePullRequest bool
}

// New Handler instance
func New(r Registry, c Configuration) *Handler {
	return &Handler{
		r: r,
		c: c,
	}
}

func (h *Handler) UpdateAll() error {
	return nil
}

func (h *Handler) Update() error {
	return nil
}
