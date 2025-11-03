package cqrs

import (
	"context"
	"errors"
	"reflect"
)

var (
	ErrCommandHasNoHandler = errors.New("command has no handler")
	ErrQueryHasNoHandler   = errors.New("query has no handler")
)

type InMemoryCommandBus struct {
	handlers map[string]interface{}
}

func NewInMemoryBus() *InMemoryCommandBus {
	return &InMemoryCommandBus{
		handlers: make(map[string]interface{}),
	}
}

func (b *InMemoryCommandBus) Register(commandName string, handler interface{}) {
	b.handlers[commandName] = handler
}

func (b *InMemoryCommandBus) Dispatch(ctx context.Context, command Command) error {
	handler, exists := b.handlers[command.CommandName()]
	if !exists {
		return ErrCommandHasNoHandler
	}

	handlerValue := reflect.ValueOf(handler)
	method := handlerValue.MethodByName("Handle")
	results := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(command)})
	if len(results) > 0 && !results[0].IsNil() {
		return results[0].Interface().(error)
	}

	return nil
}

type InMemoryQueryBus struct {
	handlers map[string]interface{}
}

func NewInMemoryQueryBus() *InMemoryQueryBus {
	return &InMemoryQueryBus{
		handlers: make(map[string]interface{}),
	}
}

func (b *InMemoryQueryBus) Register(queryName string, handler interface{}) {
	b.handlers[queryName] = handler
}

func (b *InMemoryQueryBus) Dispatch(ctx context.Context, query Query) (any, error) {
	handler, exists := b.handlers[query.QueryName()]
	if !exists {
		return nil, ErrQueryHasNoHandler
	}

	handlerValue := reflect.ValueOf(handler)
	method := handlerValue.MethodByName("Handle")
	results := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(query)})
	if len(results) > 1 && !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}

	return results[0].Interface(), nil
}
