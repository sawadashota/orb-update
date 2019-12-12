package pullrequest

import "context"

type Creator interface {
	AlreadyCreated(ctx context.Context) (bool, error)
	Create(ctx context.Context, message string) error
}
