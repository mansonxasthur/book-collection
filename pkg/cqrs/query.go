package cqrs

import "context"

type Query interface {
	QueryName() string
}

type QueryHandler[T Query, R any] interface {
	Handle(context.Context, T) (R, error)
}

type QueryBus interface {
	Dispatch(ctx context.Context, query Query) (any, error)
}
