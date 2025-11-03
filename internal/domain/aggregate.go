package domain

type Aggregate interface {
	ApplyEvent(event Event)
	Events() []Event
	ClearEvents()
}

type BaseAggregate struct {
	events []Event
}

func (a *BaseAggregate) ApplyEvent(event Event) {
	a.events = append(a.events, event)
}

func (a *BaseAggregate) Events() []Event {
	return a.events
}

func (a *BaseAggregate) ClearEvents() {
	a.events = []Event{}
}
