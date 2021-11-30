package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/dragonfax/gunroar/gr/sdl"
)

// "Usage: " ~ progName ~ " [-window] [-res x y] [-brightness [0-100]] [-luminosity [0-100]] [-nosound] [-exchange] [-turnspeed [0-500]] [-firerear] [-rotatestick2 deg] [-reversestick2] [-enableaxis5] [-NoWait]");

/**
 * Boot the game.
 */
var screen *Screen
var input *sdl.MultipleInputDevice
var pad *sdl.RecordablePad
var twinStick *sdl.RecordableTwinStick

// RecordableMouse mouse;
var gameManager *GameManager
var prefManager *PrefManager
var mainLoop *sdl.MainLoop

func main() {
	boot(os.Args)
}

func boot(args []string) {
	screen = NewScreen()
	input = sdl.NewMultipleInputDevice()
	pad = sdl.NewRecordablePad()
	twinStick = sdl.NewRecordableTwinStick()
	input.Inputs = append(input.Inputs, pad)
	input.Inputs = append(input.Inputs, twinStick)
	gameManager = NewGameManager()
	prefManager = NewPrefManager()
	mainLoop = sdl.NewMainLoop(screen, input, gameManager, prefManager)
	parseArgs(args)

	mainLoop.Loop()
}

func parseArgs(commandArgs []string) {
	args := make([]string, len(commandArgs)-1)
	for i := 1; i < len(commandArgs); i++ {
		args[i] = commandArgs[i]
	}
	progName := commandArgs[0]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-brightness":
			if i >= len(args)-1 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			is, err := strconv.Atoi(args[i])
			if err != nil {
				panic(err)
			}
			b := float64(is) / 100
			if b < 0 || b > 1 {
				usage(args[0])
				panic("Invalid options")
			}
			sdl.Brightness(b)
		case "-luminosity", "-luminous":
			if i >= len(args)-1 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			is, err := strconv.Atoi(args[i])
			if err != nil {
				panic(err)
			}
			l := float64(is) / 100
			if l < 0 || l > 1 {
				usage(progName)
				panic("Invalid options")
			}
			screen.luminosity(l)
		case "-window":
			screen.SetWindowMode(true)
		case "-res":
			if i >= len(args)-2 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			w, err := strconv.Atoi(args[i])
			if err != nil {
				panic(err)
			}
			i++
			h, err := strconv.Atoi(args[i])
			if err != nil {
				panic(err)
			}
			screen.SetWidth(w)
			screen.SetHeight(h)
		case "-nosound":
			sdl.NoSound = true
		//case "-exchange":
		//  pad.buttonReversed = true;
		case "-NoWait":
			mainLoop.NoWait = true
		case "-AccFrame":
			mainLoop.AccFrame = true
		case "-turnspeed":
			if i >= len(args)-1 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			is, err := strconv.Atoi(args[i])
			if err != nil {
				panic(err)
			}
			s := float64(is) / 100
			if s < 0 || s > 5 {
				usage(progName)
				panic("Invalid options")
			}
			shipTurnSpeed = s
		case "-firerear":
			shipReverseFire = true
		case "-rotatestick2", "-rotaterightstick":
			if i >= len(args)-1 {
				usage(progName)
				panic("Invalid options")
			}
			i++
			is, err := strconv.Atoi(args[i])
			if err != nil {
				panic(err)
			}
			twinStick.Rotate = float64(is) * math.Pi / 180.0
		case "-reversestick2", "-reverserightstick":
			twinStick.Reverse = -1
		case "-enableaxis5":
			twinStick.EnableAxis5 = true
		default:
			usage(progName)
			panic("Invalid options")
		}
	}
}

func usage(progName string) {
	fmt.Println("usage: ")
	os.Exit(1)
}
