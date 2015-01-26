package main

import "github.com/banthar/Go-SDL/sdl"

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
