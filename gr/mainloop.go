/*
 * $Id: mainloop.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

const INTERVAL_BASE = 16

type MainLoop struct {
	nowait       bool
	accframe     bool
	maxSkipFrame int
	event        *sdl.Event

	screen             *Screen
	input              *Input
	gameManager        *GameManager
	slowdownRatio      float32
	interval           float32
	slowdownStartRatio float32
	slowdownMaxRatio   float32

	done bool
}

func NewMainLoop(screen *Screen, input *Input, gm *GameManager) {
	this = MainLoop{}
	this.maxSkipFrame = 5
	this.slowdownStartRatio = 1
	this.slowdownMaxRatio = 1.75
	this.screen = screen
	this.input = input
	this.interval = INTERVAL_BASE
	gameManager.setMainLoop(*this)
	gameManager.setUIs(screen, input)
	gameManager.setPrefManager(prefManager)
	this.gameManager = gameManager
	return this
}

func (m *MainLoop) initFirst() {
	SoundManager.init()
	gameManager.init()
	initInterval()
}

func (m *MainLoop) quitLast() {
	gameManager.close()
	SoundManager.close()
	screen.closeSDL()
	sdl.Quit()
}

func (m *MainLoop) breakLoop() {
	m.done = true
}

func (m *MainLoop) loop() {
	m.done = false
	var prvTickCount long = 0
	var i int
	var nowTick long
	var frame int
	m.screen.initSDL()
	m.initFirst()
	m.gameManager.start()
	for !done {
		event := sdl.PollEvent()
		if event != nil {
			event.Type = sdl.USEREVENT
		}
		m.input.handleEvent(event)
		if event.Type == sdl.QUIT {
			breakLoop()
		}
		nowTick := sdl.GetTicks()
		var itv int = int(interval)
		var frame int = int((nowTick - prvTickCount) / itv)
		if frame <= 0 {
			frame = 1
			sdl.Delay(prvTickCount + itv - nowTick)
			if accframe {
				prvTickCount = sdl.GetTicks()
			} else {
				prvTickCount += interval
			}
		} else if frame > maxSkipFrame {
			frame = maxSkipFrame
			prvTickCount = nowTick
		} else {
			prvTickCount = nowTick
		}
		m.slowdownRatio = 0
		for i := 0; i < frame; i++ {
			m.gameManager.move()
		}
		m.slowdownRatio /= frame
		m.screen.clear()
		m.gameManager.draw()
		m.screen.flip()
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
