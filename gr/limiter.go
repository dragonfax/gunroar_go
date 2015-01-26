package main

import "github.com/banthar/Go-SDL/sdl"

/* FrameLimit handles skipping draw frames, and slowing down for performance
 */

const INTERVAL_BASE = 16 // how many milliseconds in a "frame"

type FrameLimiter struct {
	moveFrame func()
	drawFrame func()
	thenTick  uint32
}

func NewFrameLimiter(moveFrame func(), drawFrame func()) FrameLimiter {
	this := FrameLimiter{}
	this.moveFrame = moveFrame
	this.drawFrame = drawFrame
	return this
}

func (this *FrameLimiter) cycle() {

	this.moveFrame()
	this.drawFrame()

	// sleep until next frame
	nowTick := sdl.GetTicks()
	timeTaken := nowTick - this.thenTick
	if timeTaken < INTERVAL_BASE {
		sdl.Delay(INTERVAL_BASE - timeTaken)
	}
	this.thenTick = nowTick
}
