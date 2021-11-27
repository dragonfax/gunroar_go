package main

import (
	"math"
	"os"
	"strconv"
)

// "Usage: " ~ progName ~ " [-window] [-res x y] [-brightness [0-100]] [-luminosity [0-100]] [-nosound] [-exchange] [-turnspeed [0-500]] [-firerear] [-rotatestick2 deg] [-reversestick2] [-enableaxis5] [-nowait]");

/**
 * Boot the game.
 */
var screen Screen
var input sdl.MultipleInputDevice
var pad sdl.RecordablePad
var twinStick sdl.RecordableTwinStick

// RecordableMouse mouse;
var gameManager GameManager
var prefManager PrefManager
var mainLoop MainLoop

func main() {
	return boot(os.Args)
}

func boot(args []string) {
	screen = NewScreen
	input = NewMultipleInputDevice
	// pad = new RecordablePad;
	twinStick = sdl.NewRecordableTwinStick()
	// mouse = new RecordableMouse(screen);
	// input.inputs ~= pad;
	input.inputs = append(input.inputs, twinStick)
	// input.inputs ~= mouse;
	gameManager = NewGameManager()
	prefManager = NewPrefManager()
	mainLoop = NewMainLoop(screen, input, gameManager, prefManager)
	err := parseArgs(args)
	if err != nil {
		panic(err)
	}

	err = mainLoop.loop()
	if err != nil {
		panic(err)
	}
	return
}

func parseArgs(commandArgs []string) {
	args := make([]string, len(commandArgs)-1, len(commandArgs)-1)
	for i := 1; i < commandArgs.length; i++ {
		args[i] = commandArgs[i]
	}
	progName := commandArgs[0]
	for i := 0; i < args.length; i++ {
		switch args[i] {
		case "-brightness":
			if i >= len(args)-1 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			b := float64(strconv.Atoi(args[i])) / 100
			if b < 0 || b > 1 {
				usage(args[0])
				panic("Invalid options")
			}
			Screen.brightness = b
		case "-luminosity", "-luminous":
			if i >= len(args)-1 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			l = float64(strconv.Atoi(args[i])) / 100
			if l < 0 || l > 1 {
				usage(progName)
				panic("Invalid options")
			}
			screen.luminosity = l
		case "-window":
			screen.windowMode = true
		case "-res":
			if i >= len(args)-2 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			w := strconv.Atoi(args[i])
			i++
			h := strconv.Atoi(args[i])
			screen.width = w
			screen.height = h
		case "-nosound":
			SoundManager.noSound = true
		//case "-exchange":
		//  pad.buttonReversed = true;
		case "-nowait":
			mainLoop.nowait = true
		case "-accframe":
			mainLoop.accframe = 1
		case "-turnspeed":
			if i >= len(args)-1 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			s := float64(strconv.Atoi(args[i])) / 100
			if s < 0 || s > 5 {
				usage(progName)
				panic("Invalid options")
			}
			GameManager.shipTurnSpeed = s
		case "-firerear":
			GameManager.shipReverseFire = true
		case "-rotatestick2", "-rotaterightstick":
			if i >= len(args)-1 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			twinStick.rotate = float64(strconv.Atoi(args[i])) * math.Pi / 180.0
		case "-reversestick2", "-reverserightstick":
			twinStick.reverse = -1
		case "-enableaxis5":
			twinStick.enableAxis5 = true
		default:
			usage(progName)
			panic("Invalid options")
		}
	}
}

func usage(progName string) {
	Logger.error
}
