package gobus

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

var (
	rw sync.RWMutex
	b  *Bus
)

func init() {
	b = New()
}

// Our EventBus which contains all of the event funcs
type Bus struct {
	eventFuncs EventFuncsMap
}

// The options for this emit.
type emitOptions struct {
	async bool
	wait  bool
}

// NewEmitOptions returns a new instance of emitOptions.
//
// This function does not take any parameters.
// It returns a pointer to an emitOptions object.
func NewEmitOptions() *emitOptions {
	return &emitOptions{
		async: true,
		wait:  false,
	}
}

func Emit(event Event, args any, opts *emitOptions) error {

	rw.Lock()
	fns, ok := b.eventFuncs[event]
	rw.Unlock()

	if !ok {
		return fmt.Errorf("unable to find event funcs for event: %s", event)
	}

	if opts == nil {
		opts = NewEmitOptions()
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(fns))

	for _, fn := range fns {
		if opts.wait {
			wg.Add(1)
		}

		if opts.async {
			go func(fn EventFunc, args any, errCh chan<- error) {
				if opts.wait {
					defer wg.Done()
				}

				if err := fn(args); err != nil {
					fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
					errCh <- fmt.Errorf("[%s]: %w", fnName, err)
				}
			}(fn, args, errCh)
		} else {
			if err := fn(args); err != nil {
				fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
				errCh <- fmt.Errorf("[%s]: %w", fnName, err)
			}

			if opts.wait {
				wg.Done()
			}
		}
	}

	if opts.wait {
		wg.Wait()
	}

	close(errCh)

	errs := []error{}
	for err := range errCh {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// A method that eliminates manually calling `Emit`.
//
// Notes:
//
//	This works really well when you already know what to emit beforehand.
//	This may not always be the case.
func EmitAfter(event Event, opts *emitOptions, fn func()) error {
	fn()

	if opts == nil {
		opts = NewEmitOptions()
	}

	return Emit(event, nil, opts)
}

/*
Sets the event funcs on the Bus instance.

Notes:

	A declarative way to define the events and funcs for the already initialized Bus.

Example:

	var (
		EventUserSignedUp = Event("user-signed-up")
		EventUserLoggedIn = Event("user-logged-in")

		EventUserSignedUpSendWelcomeEmail EventFunc = func(args *any) error {
			return nil
		}
		EventUserSignedUpSendVerificationEmail EventFunc = func(args *any) error {
			return nil
		}

		EventUserLoggedInLog EventFunc = func(args *any) error {
			return nil
		}
		EventUserLoggedInUpdateCount EventFunc = func(args *any) error {
			return nil
		}
	)

	gobus.SetEventFuncs(EventFuncsMap{
		EventUserSignedUp: []EventFunc{
			EventUserSignedUpSendWelcomeEmail,
			EventUserSignedUpSendVerificationEmail,
		},
		EventUserLoggedIn: []EventFunc {
			EventUserLoggedInLog,
			EventUserLoggedInUpdateCount,
		},
	})
*/
func SetEventFuncs(eventFuncs EventFuncsMap) {
	rw.Lock()
	b.eventFuncs = eventFuncs
	rw.Unlock()
}

// Returns a new initialized Bus instance
func New() *Bus {
	b := new(Bus)
	b.eventFuncs = make(EventFuncsMap)
	return b
}
