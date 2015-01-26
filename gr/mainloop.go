/*
 * $Id: mainloop.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/veandco/go-sdl2/sdl"

var mainLoop *MainLoop
var limiter FrameLimiter

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
	gameManager.start()
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
					screen.resized(int(w), int(h))
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

/* FrameLimit handles skipping draw frames, and slowing down for performance
 */

const INTERVAL_BASE = 16 // how many milliseconds in a "frame"
const SlowdownStartRatio = 1
const SlowdownMaxRatio = 1.75

type FrameLimiter struct {
	moveFrame     func()
	drawFrame     func()
	nowait        bool
	accframe      bool
	maxSkipFrame  uint32
	slowdownRatio float32
	interval      uint32
	prvTickCount  uint32
}

func NewFrameLimiter(moveFrame func(), drawFrame func()) FrameLimiter {
	this := FrameLimiter{}
	this.moveFrame = moveFrame
	this.drawFrame = drawFrame
	this.maxSkipFrame = 5
	this.interval = INTERVAL_BASE
	this.initInterval()
	return this
}

func (this *FrameLimiter) cycle() {
	nowTick := sdl.GetTicks()
	var itv uint32 = this.interval
	var frame = (nowTick - this.prvTickCount) / itv
	if frame <= 0 {
		frame = 1
		sdl.Delay(this.prvTickCount + itv - nowTick)
		if this.accframe {
			this.prvTickCount = sdl.GetTicks()
		} else {
			this.prvTickCount += this.interval
		}
	} else if frame > this.maxSkipFrame {
		frame = this.maxSkipFrame
		this.prvTickCount = nowTick
	} else {
		this.prvTickCount = nowTick
	}
	this.slowdownRatio = 0
	for i := uint32(0); i < frame; i++ {
		this.moveFrame()
	}
	this.slowdownRatio = this.slowdownRatio / float32(frame)

	this.drawFrame()

	if !this.nowait {
		this.calcInterval()
	}
}

// Intentional slowdown.

func (this *FrameLimiter) initInterval() {
	this.interval = INTERVAL_BASE
}

func (this *FrameLimiter) addSlowdownRatio(sr float32) {
	this.slowdownRatio += sr
}

func (this *FrameLimiter) calcInterval() {
	if this.slowdownRatio > SlowdownStartRatio {
		sr := this.slowdownRatio / SlowdownStartRatio
		if sr > SlowdownMaxRatio {
			sr = SlowdownMaxRatio
		}
		this.interval += uint32((sr*INTERVAL_BASE - float32(this.interval)) * 0.1)
	} else {
		this.interval += uint32((INTERVAL_BASE - float32(this.interval)) * 0.08)
	}
}
