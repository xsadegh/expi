package event

type Receiver func(*Event)

type Event struct {
	Topic    string
	Event    any
	Response any
}
