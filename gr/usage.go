package main

/* Command line options and usage */

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var resRegex = regexp.MustCompile(`^(\d+)x(\d+)$`)

func parseArgs() {

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
		flag.Usage()
		os.Exit(0)
	}

	if *brightnessP < 0 {
		fmt.Println("brightness set too low")
		flag.Usage()
		os.Exit(1)
	}
	if *brightnessP > 100 {
		fmt.Println("brightness set too high")
		flag.Usage()
		os.Exit(1)
	}
	brightness = float32(*brightnessP) / 100

	if *luminosityP < 0 {
		fmt.Println("luminosity set too low")
		flag.Usage()
		os.Exit(1)
	}
	if *luminosityP > 100 {
		fmt.Println("luminosity set too high")
		flag.Usage()
		os.Exit(1)
	}
	screen.luminosity = float32(*luminosityP) / 100

	screen.windowMode = *windowModeP

	if *resP != "" {
		matches := resRegex.FindStringSubmatch(*resP)
		if len(matches) == 0 {
			fmt.Println("resolution provided does not match the required format of ###x###")
			flag.Usage()
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
			flag.Usage()
			os.Exit(1)
		}
		screen.width = w
		screen.height = h
	}

	noSound = *noSoundP

	pad.buttonReversed = *buttonReverseP

	limiter.nowait = *nowaitP

	limiter.accframe = *accframeP

	if *turnspeedP < 0 {
		fmt.Println("ship turning speed is too low")
		flag.Usage()
		os.Exit(1)
	}
	if *turnspeedP > 500 {
		fmt.Println("ship turning speed is too high")
		flag.Usage()
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
