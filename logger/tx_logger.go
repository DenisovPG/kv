package logger

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

type TxStoreLogger struct {
	events chan Event
}

func (l *TxStoreLogger) Run() {
	l.events = make(chan Event, 16)
	go func() {
		for ev := range l.events {
			fmt.Printf("%s %s %s %s\n", ev.timestamp.Format(time.RFC3339), ev.EventType, ev.Key, ev.Value)
		}
	}()
}

func (l *TxStoreLogger) LogPut(key, value string) {
	l.events <- Event{
		timestamp: time.Now(),
		EventType: EventPut,
		Key:       key,
		Value:     value,
	}
}

func (l *TxStoreLogger) LogGet(key string) {
	l.events <- Event{
		timestamp: time.Now(),
		EventType: EventGet,
		Key:       key,
		Value:     "",
	}
}

func (l *TxStoreLogger) LogDelete(key string) {
	l.events <- Event{
		timestamp: time.Now(),
		EventType: EventDelete,
		Key:       key,
		Value:     "",
	}
}

func (l *TxStoreLogger) Stop() {
	close(l.events)
}
