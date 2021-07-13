package eventbus

import (
	"sync"
)

// Bus event bus instance
var Bus = newBus()

// Event event struct
type Event struct {
	Topic string
	Data  interface{}
}

// EventCh represents event channel
type EventCh chan Event

// EventChSlice represents event channel slice
type EventChSlice []EventCh

// event bus struct
type bus struct {
	subscribers map[string]EventChSlice
	closed      bool
	rw          sync.RWMutex
}

// returns a new event bus instance
func newBus() *bus {
	return &bus{
		subscribers: make(map[string]EventChSlice),
		rw:          sync.RWMutex{},
	}
}

// Publish topic with data
func (eb *bus) Publish(topic string, data interface{}) {
	eb.rw.RLock()
	defer eb.rw.RUnlock()

	if eb.closed {
		panic("event bus is closed")
	}

	if chs, exist := eb.subscribers[topic]; exist {
		ev := Event{
			Topic: topic,
			Data:  data,
		}
		newChs := append(EventChSlice{}, chs...)
		go func(ev Event, chs EventChSlice) {
			for _, ch := range chs {
				ch <- ev
			}
		}(ev, newChs)
	}
}

// Subscribe subscribe topic with a specific event channel
func (eb *bus) Subscribe(topic string, ch EventCh) {
	eb.rw.Lock()
	defer eb.rw.Unlock()

	if eb.closed {
		panic("event bus is closed")
	}

	if chs, exist := eb.subscribers[topic]; exist {
		eb.subscribers[topic] = append(chs, ch)
	} else {
		eb.subscribers[topic] = append(EventChSlice{}, ch)
	}
}

// Unsubscribe unsubscribe the specific topic
func (eb *bus) Unsubscribe(topic string, ch EventCh) {
	eb.rw.Lock()
	defer eb.rw.Unlock()

	if chs, exist := eb.subscribers[topic]; exist {
		for i := range chs {
			if chs[i] == ch {
				if len(chs) == 1 {
					delete(eb.subscribers, topic)
				} else {
					chs = append(chs[:i], chs[i+1:]...)
				}
			}
		}
	}
}

// Close close the eventbus and release all subscribers
func (eb *bus) Close() {
	eb.rw.Lock()
	defer eb.rw.Unlock()

	eb.closed = true

	for k := range eb.subscribers {
		eb.subscribers[k] = nil
	}

	eb.subscribers = nil
}
