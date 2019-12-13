package pullrequest

import "context"

type Creator interface {
	AlreadyCreated(ctx context.Context, branch string) (bool, error)
	Create(ctx context.Context, message, baseBranch string) error
}
