package events

import "context"

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
	Callback
)

type Event struct {
	Type Type
	Text string
	Meta any
}

type RoutineData struct {
	channel *chan Event
	context *context.Context
}
