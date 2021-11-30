package sdl

import (
	"github.com/dragonfax/gunroar/gr/sdl/screen"
	"github.com/veandco/go-sdl2/sdl"
)

/**
 * SDL main loop.
 */

const INTERVAL_BASE = 16

type PrefManager interface {
	Save()
	Load()
}

type MainLoop struct {
	NoWait              bool
	AccFrame            bool
	maxSkipFrame        uint32
	Event               sdl.Event
	screen              screen.Screen
	input               Input
	gameManager         GameManager
	prefManager         PrefManager
	slowdownRatio       float64
	interval            float64
	_slowdownStartRatio float64
	_slowdownMaxRatio   float64
	done                bool
}

func NewMainLoop(screen screen.Screen, input Input, gameManager GameManager, prefManager PrefManager) *MainLoop {
	this := &MainLoop{
		maxSkipFrame:        5,
		interval:            INTERVAL_BASE,
		_slowdownStartRatio: 1,
		_slowdownMaxRatio:   1.75,
		screen:              screen,
		input:               input,
		gameManager:         gameManager,
		prefManager:         prefManager,
	}
	gameManager.SetMainLoop(this)
	gameManager.SetUIs(screen, input)
	gameManager.SetPrefManager(prefManager)
	return this
}

// Initialize and load preference.
func (this *MainLoop) initFirst() {
	this.prefManager.Load()
	SoundManagerInit()
	this.gameManager.Init()
	this.InitInterval()
}

func (this *MainLoop) quitLast() {
	sdl.Quit()
}

func (this *MainLoop) BreakLoop() {
	this.done = true
}

func (this *MainLoop) Loop() {
	this.done = false
	var prvTickCount uint32
	var nowTick uint32
	var frame uint32
	this.screen.InitSDL()
	this.initFirst()
	this.gameManager.Start()
	for !this.done {
		this.Event = sdl.PollEvent()
		if this.Event != nil {
			// TODO this.event.type = sdl.USEREVENT;
			this.input.HandleEvent(this.Event)
			if this.Event.GetType() == sdl.QUIT {
				this.BreakLoop()
			}
		}
		nowTick = sdl.GetTicks()
		itv := uint32(this.interval)
		frame = (nowTick - prvTickCount) / itv
		if frame <= 0 {
			frame = 1
			sdl.Delay(prvTickCount + itv - nowTick)
			if this.AccFrame {
				prvTickCount = sdl.GetTicks()
			} else {
				prvTickCount += uint32(this.interval)
			}
		} else if frame > this.maxSkipFrame {
			frame = this.maxSkipFrame
			prvTickCount = nowTick
		} else {
			prvTickCount = nowTick
		}
		this.slowdownRatio = 0
		for i := uint32(0); i < frame; i++ {
			this.gameManager.Move()
		}
		this.slowdownRatio /= float64(frame)
		this.screen.Clear()
		this.gameManager.Draw()
		this.screen.Flip()
		if !this.NoWait {
			this.calcInterval()
		}
	}
	this.quitLast()
}

// Intentional slowdown.

func (this *MainLoop) InitInterval() {
	this.interval = INTERVAL_BASE
}

func (this *MainLoop) AddSlowdownRatio(sr float64) {
	this.slowdownRatio += sr
}

func (this *MainLoop) calcInterval() {
	if this.slowdownRatio > this._slowdownStartRatio {
		sr := this.slowdownRatio / this._slowdownStartRatio
		if sr > this._slowdownMaxRatio {
			sr = this._slowdownMaxRatio
		}
		this.interval += (sr*INTERVAL_BASE - this.interval) * 0.1
	} else {
		this.interval += (INTERVAL_BASE - this.interval) * 0.08
	}
}

func (this *MainLoop) slowdownStartRatio(v float64) float64 {
	this._slowdownStartRatio = v
	return v
}

func (this *MainLoop) slowdownMaxRatio(v float64) float64 {
	this._slowdownMaxRatio = v
	return v
}
