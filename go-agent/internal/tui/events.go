package tui

import "fmt"

type Event struct {
	Generation int64
	Name       string
	Content    string
	Data       any
}

type EventWriter struct {
	ch         chan<- Event
	generation int64
}

func NewEventWriter(ch chan<- Event, generation ...int64) EventWriter {
	var gen int64
	if len(generation) > 0 {
		gen = generation[0]
	}
	return EventWriter{ch: ch, generation: gen}
}

func (w EventWriter) Event(name string, data any) error {
	if w.ch == nil {
		return nil
	}
	w.ch <- Event{Generation: w.generation, Name: name, Content: eventContent(data), Data: data}
	return nil
}

func eventContent(data any) string {
	if m, ok := data.(map[string]any); ok {
		if value, ok := m["content"].(string); ok {
			return value
		}
		if value, ok := m["message"].(string); ok {
			return value
		}
	}
	return fmt.Sprint(data)
}
