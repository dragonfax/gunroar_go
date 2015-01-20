/*
 * $Id: mainloop.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

const INTERVAL_BASE = 16

var mainLoop *MainLoop

func main() {
	screen = NewScreen()
	pad = NewPad()
	twinStick = NewTwinStick()
	mouse = NewMouse()
	mouseAndPad = NewMouseAndPad()
	gameManager = NewGameManager()
	mainLoop = NewMainLoop()
	parseArgs()
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
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				m.breakLoop()
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

var resRegex = regexp.MustCompile(`^(\d+)x(\d+)$`)

func parseArgs() {

	progName := os.Args[0]

	helpP := flag.Bool("help", false, "show usage")
	brightnessP := flag.Int("brightness", 100, "0-100")
	luminosityP := flag.Int("luminosity", 100, "lumonisity, 0-100")
	windowModeP := flag.Bool("window", false, "play in a window (instead of full screen)")
	resP := flag.String("res", "", "resolution to play at, ex 640x480")
	noSoundP := flag.Bool("nosound", false, "disable sound")
	buttonReverseP := flag.Bool("exchange", false, "swap buttons")
	nowaitP := flag.Bool("nowait", false, "disable wait in loop")
	accframeP := flag.Bool("accframe", false, "")
	turnspeedP := flag.Int("turnspeed", 100, "ship turning speed, 0-500")
	firerearP := flag.Bool("firerear", false, "fire from back of ship")
	rightStickRotationP := flag.Int("rotaterightstick", 0, "degree to rotate right stick control")
	reverseRightStickP := flag.Bool("reverserightstick", false, "reverse right stick controls")
	enableAxis5P := flag.Bool("enableaxis5", false, "enable the 5th axis for some controllers")

	flag.Parse()

	if *helpP {
		usage(progName)
		os.Exit(0)
	}

	if *brightnessP < 0 {
		fmt.Println("brightness set too low")
		usage(progName)
		os.Exit(1)
	}
	if *brightnessP > 100 {
		fmt.Println("brightness set too high")
		usage(progName)
		os.Exit(1)
	}
	brightness = float32(*brightnessP) / 100

	if *luminosityP < 0 {
		fmt.Println("luminosity set too low")
		usage(progName)
		os.Exit(1)
	}
	if *luminosityP > 100 {
		fmt.Println("luminosity set too high")
		usage(progName)
		os.Exit(1)
	}
	screen.luminosity = float32(*luminosityP) / 100

	screen.windowMode = *windowModeP

	if *resP != "" {
		matches := resRegex.FindStringSubmatch(*resP)
		if len(matches) == 0 {
			fmt.Println("resolution provided does not match the required format of ###x###")
			usage(progName)
			os.Exit(1)
		}
		var w, h int
		var err error
		w, err = strconv.Atoi(matches[1])
		if err == nil {
			h, err = strconv.Atoi(matches[2])
		}
		if err != nil {
			fmt.Println("Error parsing the width and height values")
			usage(progName)
		}
		screen.width = w
		screen.height = h
	}

	noSound = *noSoundP

	pad.buttonReversed = *buttonReverseP

	mainLoop.nowait = *nowaitP

	mainLoop.accframe = *accframeP

	if *turnspeedP < 0 {
		fmt.Println("ship turning speed is too low")
		usage(progName)
		os.Exit(1)
	}
	if *turnspeedP > 500 {
		fmt.Println("ship turning speed is too high")
		usage(progName)
		os.Exit(1)
	}
	shipTurnSpeed = float32(*turnspeedP) / 100

	shipReverseFire = *firerearP

	twinStick.rotate = float32(*rightStickRotationP) * Pi32 / 180

	if *reverseRightStickP {
		twinStick.reverse = -1
	}

	twinStick.enableAxis5 = *enableAxis5P

}

func usage(progName string) {
	fmt.Printf("Usage: " + progName + " [-window] [-res x y] [-brightness [0-100]] [-luminosity [0-100]] [-nosound] [-exchange] [-turnspeed [0-500]] [-firerear] [-rotatestick2 deg] [-reversestick2] [-enableaxis5] [-nowait]\n")
}
