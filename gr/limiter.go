package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

/* FrameLimit handles skipping draw frames, and slowing down for performance
 */

const INTERVAL_BASE = 2 // how many milliseconds in a "frame"

const NUM_FPS_TO_AVERAGE = 10

var limiter FrameLimiter

type FrameLimiter struct {
	moveFrame   func()
	drawFrame   func()
	thenTick    uint32
	previousFps []float32
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

	this.addFps(sdl.GetTicks() - this.thenTick)

	this.thenTick = nowTick
}

func (this *FrameLimiter) addFps(frameTime uint32) {
	fps := 1000 / float32(frameTime)
	this.previousFps = append(this.previousFps, fps)
	if len(this.previousFps) > NUM_FPS_TO_AVERAGE {
		this.previousFps = this.previousFps[1:len(this.previousFps)]
	}
}

func (this *FrameLimiter) draw() {
	var totalFps float32
	for _, fps := range this.previousFps {
		totalFps += fps
	}
	avgFps := totalFps / float32(len(this.previousFps))
	drawString(fmt.Sprintf("%3d", int(avgFps)), 10, 10, 3.0)
}
