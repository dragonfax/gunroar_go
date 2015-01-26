package main

import (
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type MainFunc func()

var mainFuncC = make(chan MainFuncEvent)

type WaitChannel chan bool // the value means nothing

type MainFuncEvent struct {
	F MainFunc
	W WaitChannel
}

func QueueMain(f MainFunc, w WaitChannel) {
	mainFuncC <- MainFuncEvent{f, w}
}

// to run as the last step in the main thread. (not in a goroutine)
func MainQueueLoop() {
	runtime.LockOSThread()

	for !mainLoop.done {
		var didWork = false
		select {
		case receivedFuncEvent := <-mainFuncC:
			receivedFuncEvent.F()
			if receivedFuncEvent.W != nil {
				receivedFuncEvent.W <- true
			}
			didWork = true
		default:
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				if event != nil {
					eventSendC <- event
					didWork = true
				}
			}
		}

		if !didWork {
			time.Sleep(time.Millisecond * 10)
		}
	}

}
