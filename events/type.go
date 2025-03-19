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
	//Meta any
	Meta Meta
}

type Meta struct {
	ChatID   int
	UserName string
}

type RoutineData struct {
	Channel chan Event
	Context context.Context
}
