package tx_log

import (
	"fmt"
	"time"
)

type EventType byte

const (
	EventGet EventType = iota
	EventPut
	EventDelete
)

func (et EventType) String() string {
	switch et {
	case EventGet:
		return "GET"
	case EventPut:
		return "PUT"
	case EventDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

type Event struct {
	timestamp time.Time
	EventType EventType
	Key       string
	Value     string
}

type ConsoleTxStoreLogger struct {
	events chan Event
}

func (l *ConsoleTxStoreLogger) Run() {
	l.events = make(chan Event, 16)
	go func() {
		for ev := range l.events {
			fmt.Printf("%s %s %s %s\n", ev.timestamp.Format(time.RFC3339), ev.EventType, ev.Key, ev.Value)
		}
	}()
}

func (l *ConsoleTxStoreLogger) LogPut(key, value string) {
	l.events <- Event{
		timestamp: time.Now(),
		EventType: EventPut,
		Key:       key,
		Value:     value,
	}
}

func (l *ConsoleTxStoreLogger) LogGet(key string) {
	l.events <- Event{
		timestamp: time.Now(),
		EventType: EventGet,
		Key:       key,
		Value:     "",
	}
}

func (l *ConsoleTxStoreLogger) LogDelete(key string) {
	l.events <- Event{
		timestamp: time.Now(),
		EventType: EventDelete,
		Key:       key,
		Value:     "",
	}
}

func (l *ConsoleTxStoreLogger) Stop() {
	close(l.events)
}

