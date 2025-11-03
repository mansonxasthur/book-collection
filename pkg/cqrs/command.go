package cqrs

import "context"

type Command interface {
	CommandName() string
}

type CommandHandler[T Command] interface {
	Handle(context.Context, T) error
}

type CommandBus interface {
	Dispatch(ctx context.Context, command Command) error
}
