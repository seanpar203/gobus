package gobus

// Custom type for string to make declaring, using and refactoring easier.
type Event string

// Our custom BusFunc
type EventFunc = func(args any) error

// Custom type for our map of event funcs
type EventFuncsMap = map[Event][]EventFunc
