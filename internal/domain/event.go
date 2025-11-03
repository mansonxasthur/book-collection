package domain

import "time"

type Event interface {
	Name() string
	OccurredAt() time.Time
	Version() string
	AggregateID() string
}

type BaseEvent struct {
	NameVal       string
	OccurredAtVal time.Time
	VersionVal    string
	AggregateVal  string
}

func (e *BaseEvent) Name() string          { return e.NameVal }
func (e *BaseEvent) OccurredAt() time.Time { return e.OccurredAtVal }
func (e *BaseEvent) Version() string       { return e.VersionVal }
func (e *BaseEvent) AggregateID() string   { return e.AggregateVal }
