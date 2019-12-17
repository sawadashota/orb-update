package pullrequest

import "context"

// Creator of Pull Request
type Creator interface {
	AlreadyCreated(ctx context.Context, branch string) (bool, error)
	Create(ctx context.Context, message, baseBranch string) error
}
