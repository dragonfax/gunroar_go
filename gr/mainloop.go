/*
 * $Id: mainloop.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/veandco/go-sdl2/sdl"

var mainLoop *MainLoop

func main() {
	mainLoop = NewMainLoop()
	mainLoop.run()
}

type MainLoop struct {
	event *sdl.Event

	done bool
}

func NewMainLoop() *MainLoop {
	this := new(MainLoop)
	return this
}

func (this *MainLoop) run() {
	mainLoop.setup()
	MainQueueLoop()
	mainLoop.tearDown()
}

func (m *MainLoop) setup() {
	screen = NewScreen()
	pad = NewPad()
	twinStick = NewTwinStick()
	twinStick.update()
	mouse = NewMouse()
	mouseAndPad = NewMouseAndPad()
	mouseAndPad.update()
	gameManager = NewGameManager()
	limiter = NewFrameLimiter()
	drawLimiter = NewFrameLimiter()
	parseArgs()
	go MuxEvents()
	screen.initSDL()
	InitSoundManager()
	gameManager.init()
	gameManager.startTitle()
	displayListsFinalized = true
	go m.inputLoop()
	go m.drawLoop()
	go m.moveLoop()
}

func (m *MainLoop) moveLoop() {
	for {
		gameManager.move()
		limiter.cycle()
	}
}

func (m *MainLoop) inputLoop() {
	eventReceiver := GetEventReceiver()
	for {
		m.handleInput(eventReceiver)
	}
}

func (m *MainLoop) handleInput(eventReceiver EventC) {
	select {
	case event := <-eventReceiver:
		switch e := event.(type) {
		case *sdl.QuitEvent:
			m.done = true
		case *sdl.WindowEvent:
			switch e.Event {
			case sdl.WINDOWEVENT_RESIZED:
				w := e.Data1
				h := e.Data2
				if w > 150 && h > 100 {
					QueueMain(func() {
						screen.resized(int(w), int(h))
					}, nil)
				}
			}
		}
	}
	mouseAndPad.update()
	twinStick.update()
}

// goroutine for drawing the screen.
func (m *MainLoop) drawLoop() {
	w := make(WaitChannel)
	for {
		QueueMain(m.draw, w)
		<-w // wait until the drawing is complete
		drawLimiter.cycle()
	}
}

func (m *MainLoop) draw() {
	screen.clear()
	gameManager.draw()
	screen.flip()
}

func (m *MainLoop) tearDown() {
	gameManager.close()
	CloseSoundManager()
	screen.closeSDL()
	sdl.Quit()
}
