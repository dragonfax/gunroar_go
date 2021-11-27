package sdl

import "github.com/veandco/go-sdl2/sdl"

/**
 * SDL main loop.
 */

const INTERVAL_BASE = 16

type MainLoop struct {
	nowait              bool
	accframe            bool
	maxSkipFrame        int
	event               *SDL_Event
	screen              *Screen
	input               Input
	gameManager         GameManager
	prefManager         PrefManager
	slowdownRatio       float64
	interval            float64
	_slowdownStartRatio float64
	_slowdownMaxRatio   float64
	done                bool
}

func NewMainLoop(screen Screen, input Input, gameManager GameManager, prefManager PrefManager) *MainLoop {
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
	gameManager.setMainLoop(this)
	gameManager.setUIs(screen, input)
	gameManager.setPrefManager(prefManager)
	return this
}

// Initialize and load preference.
func (this *MainLoop) initFirst() {
	this.prefManager.load()
	err := SoundManagerInit()
	if err != nil {
		Logger.error(err)
	}
	this.gameManager.init()
	this.initInterval()
}

// Quit and save preference.
func (this *MainLoop) quitLast() {
	this.gameManager.close()
	this.SoundManagerClose()
	this.prefManager.save()
	this.screen.closeSDL()
	this.SDL_Quit()
}

func (this *MainLoop) breakLoop() {
	this.done = true
}

func (this *MainLoop) loop() {
	this.done = false
	var prvTickCount uint32
	var i int
	var nowTick uint32
	var frame uint32
	this.screen.initSDL()
	this.initFirst()
	this.gameManager.start()
	for !this.done {
		this.event = sdl.PollEvent()
		if this.event != nil {
			// TODO this.event.type = sdl.USEREVENT;
			this.input.HandleEvent(this.event)
			if this.event.GetType() == sdl.QUIT {
				this.breakLoop()
			}
		}
		nowTick = sdl.GetTicks()
		itv := uint32(this.interval)
		frame := (nowTick - prvTickCount) / itv
		if frame <= 0 {
			frame = 1
			sdl.Delay(prvTickCount + itv - nowTick)
			if this.accframe {
				prvTickCount = sdl.GetTicks()
			} else {
				prvTickCount += this.interval
			}
		} else if frame > this.maxSkipFrame {
			frame = this.maxSkipFrame
			prvTickCount = nowTick
		} else {
			prvTickCount = nowTick
		}
		this.slowdownRatio = 0
		for i := 0; i < frame; i++ {
			this.gameManager.move()
		}
		this.slowdownRatio /= frame
		screen.clear()
		this.gameManager.draw()
		screen.flip()
		if !this.nowait {
			this.calcInterval()
		}
	}
	this.quitLast()
}

// Intentional slowdown.

func (this *MainLoop) initInterval() {
	this.interval = INTERVAL_BASE
}

func (this *MainLoop) addSlowdownRatio(sr float64) {
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
