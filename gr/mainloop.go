/*
 * $Id: mainloop.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const INTERVAL_BASE = 16

var input *MouseAndPad

// twinStick TwinStick
var mainLoop *MainLoop

func main() {
	screen = NewScreen()
	pad = NewPad()
	// twinStick = new RecordableTwinStick
	mouse = NewMouse()
	input = NewMouseAndPad(mouse, pad)
	gameManager = NewGameManager()
	mainLoop = NewMainLoop()
	mainLoop.loop()
}

type MainLoop struct {
	nowait       bool
	accframe     bool
	maxSkipFrame int
	event        *sdl.Event

	slowdownRatio      float32
	interval           float32
	slowdownStartRatio float32
	slowdownMaxRatio   float32

	done bool
}

func NewMainLoop() *MainLoop {
	this := new(MainLoop)
	this.maxSkipFrame = 5
	this.slowdownStartRatio = 1
	this.slowdownMaxRatio = 1.75
	this.interval = INTERVAL_BASE
	return this
}

func (this *MainLoop) initFirst() {
	InitSoundManager()
	gameManager.init()
	this.initInterval()
}

func (m *MainLoop) quitLast() {
	gameManager.close()
	CloseSoundManager()
	screen.closeSDL()
	sdl.Quit()
}

func (m *MainLoop) breakLoop() {
	m.done = true
}

func (m *MainLoop) loop() {
	m.done = false
	var prvTickCount int32 = 0
	var i int
	var nowTick int32
	var frame int
	screen.initSDL()
	m.initFirst()
	gameManager.start()
	for !m.done {
		event := sdl.PollEvent()
		/*if event != nil {
			event.Type = sdl.USEREVENT
		}
		*/
		input.handleEvent(event)
		if event == sdl.QUIT {
			m.breakLoop()
		}
		nowTick := sdl.GetTicks()
		var itv int = int(m.interval)
		var frame int = int((nowTick - uint32(prvTickCount)) / itv)
		if frame <= 0 {
			frame = 1
			sdl.Delay(prvTickCount + itv - nowTick)
			if m.accframe {
				prvTickCount = sdl.GetTicks()
			} else {
				prvTickCount += m.interval
			}
		} else if frame > m.maxSkipFrame {
			frame = m.maxSkipFrame
			prvTickCount = nowTick
		} else {
			prvTickCount = nowTick
		}
		m.slowdownRatio = 0
		for i := 0; i < frame; i++ {
			gameManager.move()
		}
		m.slowdownRatio /= frame
		screen.clear()
		gameManager.draw()
		screen.flip()
		if !m.nowait {
			m.calcInterval()
		}
	}
	m.quitLast()
}

// Intentional slowdown.

func (m *MainLoop) initInterval() {
	m.interval = INTERVAL_BASE
}

func (m *MainLoop) addSlowdownRatio(sr float32) {
	m.slowdownRatio += sr
}

func (m *MainLoop) calcInterval() {
	if m.slowdownRatio > m.slowdownStartRatio {
		sr := m.slowdownRatio / m.slowdownStartRatio
		if sr > m.slowdownMaxRatio {
			sr = m.slowdownMaxRatio
		}
		m.interval += (sr*INTERVAL_BASE - m.interval) * 0.1
	} else {
		m.interval += (INTERVAL_BASE - m.interval) * 0.08
	}
}
