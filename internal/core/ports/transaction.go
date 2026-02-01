package ports

import "context"

type (
	Transaction interface {
		Begin(ctx context.Context) context.Context
		End(ctx context.Context, err *error) error
	}
)
