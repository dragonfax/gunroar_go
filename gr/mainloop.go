/*
 * $Id: mainloop.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/veandco/go-sdl2/sdl"

const INTERVAL_BASE = 16

// twinStick TwinStick
var mainLoop *MainLoop

func main() {
	screen = NewScreen()
	pad = NewPad()
	// twinStick = new RecordableTwinStick
	mouse = NewMouse()
	mouseAndPad = NewMouseAndPad()
	gameManager = NewGameManager()
	mainLoop = NewMainLoop()
	mainLoop.loop()
}

type MainLoop struct {
	nowait       bool
	accframe     bool
	maxSkipFrame uint32
	event        *sdl.Event

	slowdownRatio      float32
	interval           uint32
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
	var prvTickCount uint32 = 0
	screen.initSDL()
	m.initFirst()
	gameManager.start()
	for !m.done {
		event := sdl.PollEvent()
		mouseAndPad.handleEvent(event)
		switch event.(type) {
		case *sdl.QuitEvent:
			m.breakLoop()
		}
		nowTick := sdl.GetTicks()
		var itv uint32 = m.interval
		var frame = (nowTick - prvTickCount) / itv
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
		for i := uint32(0); i < frame; i++ {
			gameManager.move()
		}
		m.slowdownRatio = m.slowdownRatio / float32(frame)
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
		m.interval += uint32((sr*INTERVAL_BASE - float32(m.interval)) * 0.1)
	} else {
		m.interval += uint32((INTERVAL_BASE - float32(m.interval)) * 0.08)
	}
}
