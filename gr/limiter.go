package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

/* FrameLimit handles skipping draw frames, and slowing down for performance
 */

// how many milliseconds in a "frame"
var INTERVAL_BASE uint32 = 16

const NUM_FPS_TO_AVERAGE = 10

var limiter FrameLimiter
var drawLimiter FrameLimiter

type FrameLimiter struct {
	thenTick    uint32
	previousFps []float32
}

func NewFrameLimiter() FrameLimiter {
	this := FrameLimiter{}
	return this
}

func (this *FrameLimiter) cycle() {

	// sleep until next frame
	nowTick := sdl.GetTicks()
	timeTaken := nowTick - this.thenTick

	if timeTaken < INTERVAL_BASE {
		sleepTime := INTERVAL_BASE - timeTaken
		time.Sleep(time.Millisecond * time.Duration(sleepTime))
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

func (this *FrameLimiter) draw(px, py float32) {
	var totalFps float32
	for _, fps := range this.previousFps {
		totalFps += fps
	}
	avgFps := totalFps / float32(len(this.previousFps))
	drawString(fmt.Sprintf("%3d", int(avgFps)), 10+px, 10+py, 3.0)
}
