package event

type Receiver func(*Event)

type Event struct {
	Topic    string
	Error    error
	Event    any
	Response any
}
