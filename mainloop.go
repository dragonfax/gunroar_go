/*
 * $Id: mainloop.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

var mainLoop *MainLoop

func main() {
	runtime.LockOSThread()
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
	mainLoop.loop()
	mainLoop.tearDown()
}

func (m *MainLoop) setup() {
	screen = NewScreen()
	pad = NewPad()
	twinStick = NewTwinStick()
	mouse = NewMouse()
	mouseAndPad = NewMouseAndPad()
	gameManager = NewGameManager()
	limiter = NewFrameLimiter(gameManager.move, m.draw)
	parseArgs()
	screen.initSDL()
	InitSoundManager()
	gameManager.init()
	gameManager.startTitle()
	displayListsFinalized = true
}

func (m *MainLoop) loop() {
	for !m.done {
		m.handleInput()
		limiter.cycle()
	}
}

func (m *MainLoop) handleInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			m.done = true
		case *sdl.WindowEvent:
			switch e.Event {
			case sdl.WINDOWEVENT_RESIZED:
				w := e.Data1
				h := e.Data2
				if w > 150 && h > 100 {
					screen.resized(uint32(w), uint32(h))
				}
			}
		}
	}
	mouseAndPad.update()
	twinStick.update()
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
