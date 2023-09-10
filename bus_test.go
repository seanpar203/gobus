package gobus

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type errFailedToProcessEvent struct{}

func (err errFailedToProcessEvent) Error() string {
	return "Failed to process event"
}

func printWorld(args *any) error {
	fmt.Println("World!")
	return nil
}

func sleepAndReturnError(args *any) error {
	time.Sleep(time.Second * 2)
	return errFailedToProcessEvent{}
}

func TestNew(t *testing.T) {
	t.Parallel()

	bus := New()

	assert.IsType(t, &Bus{}, bus)
	assert.Empty(t, bus.eventFuncs, "Expecting event map to be empty")
}

func TestSetEventFuncs(t *testing.T) {
	t.Parallel()

	eventName := Event("set-event-funcs")
	eventFunc := func(args *any) error { return nil }

	SetEventFuncs(EventFuncsMap{
		eventName: []EventFunc{eventFunc},
	})

	assert.NotEmpty(t, b.eventFuncs, "The event funcs map shouldn't be empty")
	assert.Contains(t, b.eventFuncs, eventName, "The event funcs should contain %s", eventName)

	SetEventFuncs(EventFuncsMap{})
}

func TestEmit(t *testing.T) {
	t.Parallel()

	var hello = Event("testing-emit")

	SetEventFuncs(EventFuncsMap{
		hello: []EventFunc{
			printWorld,
			sleepAndReturnError,
		},
	})

	fmt.Println("Hello ")

	opts := NewEmitOptions()
	opts.async = false

	err := Emit(hello, nil, opts)

	assert.True(t, errors.As(err, &errFailedToProcessEvent{}), "should have custom error: errFailedToProcessEvent")
}

func TestEmitAfter(t *testing.T) {
	t.Parallel()

	var event = Event("emit-after")

	m := make(map[Event]bool)

	setEventToTrue := func(args *any) error {
		m[event] = true
		return nil
	}

	SetEventFuncs(EventFuncsMap{
		event: []EventFunc{
			setEventToTrue,
		},
	})

	EmitAfter(event, nil, func() {
		fmt.Println("Real work would happen here and all subsequent events would be emitted")
	})

	time.Sleep(1 * time.Second)

	val, ok := m[event]

	if !ok {
		t.Error("Event not emitted")
	}

	assert.True(t, val, fmt.Sprintf("map key: %s should be true", event))
}
