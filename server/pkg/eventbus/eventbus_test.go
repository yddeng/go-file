package eventbus

import (
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {
	eb := newBus()

	taskEventCh1 := make(EventCh)
	eb.Subscribe("topic 1", taskEventCh1)
	defer eb.Unsubscribe("topic 1", taskEventCh1)

	taskEventCh2 := make(EventCh)
	eb.Subscribe("topic 2", taskEventCh2)
	defer eb.Unsubscribe("topic 2", taskEventCh2)

	go func() {
		eb.Publish("topic 1", "mock data")
	}()

	for {
		select {
		case ev := <-taskEventCh1:
			t.Logf("recv 1: %v", ev)

		case <-taskEventCh2:
			t.Fatal("should not get topic 2")

		case <-time.After(1 * time.Second):
			t.Fatal("time out")
			return
		}
	}
}
