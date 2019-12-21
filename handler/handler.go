package handler

// Handler .
type Handler struct {
	r                     Registry
	c                     Configuration
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

// New Handler instance
func New(r Registry, c Configuration, opts ...Option) *Handler {
	h := &Handler{
		r:                     r,
		c:                     c,
		doesCreatePullRequest: false,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}
