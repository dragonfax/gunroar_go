package main

import "github.com/veandco/go-sdl2/sdl"

type EventC chan sdl.Event

var eventSendC = make(EventC)

var eventReceivers = make([]EventC, 0)

// to run in a goroutine to make sure all event listeners get all events.
func MuxEvents() {
	for {
		select {
		case event := <-eventSendC:
			for _, ec := range eventReceivers {
				ec <- event
			}
		}
	}
}

// for anyone to get an event queue of their own
func GetEventReceiver() EventC {
	eventReceiver := make(EventC)
	eventReceivers = append(eventReceivers, eventReceiver)
	return eventReceiver
}
