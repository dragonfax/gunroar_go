/*
 * $Id: mainloop.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/veandco/go-sdl2/sdl"

const INTERVAL_BASE = 16

var mainLoop *MainLoop

func main() {
	mainLoop = NewMainLoop()
	mainLoop.run()
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
	prvTickCount       uint32

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
	parseArgs()
	m.done = false
	m.prvTickCount = 0
	screen.initSDL()
	InitSoundManager()
	gameManager.init()
	m.initInterval()
	gameManager.start()
	displayListsFinalized = true
}

func (m *MainLoop) loop() {
	for !m.done {
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
						screen.resized(int(w), int(h))
					}
				}
			}
		}
		mouseAndPad.update()
		twinStick.update()
		nowTick := sdl.GetTicks()
		var itv uint32 = m.interval
		var frame = (nowTick - m.prvTickCount) / itv
		if frame <= 0 {
			frame = 1
			sdl.Delay(m.prvTickCount + itv - nowTick)
			if m.accframe {
				m.prvTickCount = sdl.GetTicks()
			} else {
				m.prvTickCount += m.interval
			}
		} else if frame > m.maxSkipFrame {
			frame = m.maxSkipFrame
			m.prvTickCount = nowTick
		} else {
			m.prvTickCount = nowTick
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
}

func (m *MainLoop) tearDown() {
	gameManager.close()
	CloseSoundManager()
	screen.closeSDL()
	sdl.Quit()
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
