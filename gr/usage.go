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
	showFpsP := flag.Bool("showfps", false, "show fps counter")
	brightnessP := flag.Int("brightness", 100, "0-100")
	luminosityP := flag.Int("luminosity", 100, "lumonisity, 0-100")
	windowModeP := flag.Bool("window", false, "play in a window (instead of full screen)")
	resP := flag.String("res", "", "resolution to play at, ex 640x480")
	noSoundP := flag.Bool("nosound", false, "disable sound")
	buttonReverseP := flag.Bool("exchange", false, "swap buttons")
	rightStickRotationP := flag.Int("rotaterightstick", 0, "degree to rotate right stick control")
	reverseRightStickP := flag.Bool("reverserightstick", false, "reverse right stick controls")
	enableAxis5P := flag.Bool("enableaxis5", false, "enable the 5th axis for some controllers")

	flag.Parse()

	if *helpP {
		flag.Usage()
		os.Exit(0)
	}

	showFps = *showFpsP

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

	twinStick.rotate = float32(*rightStickRotationP) * Pi32 / 180

	if *reverseRightStickP {
		twinStick.reverse = -1
	}

	twinStick.enableAxis5 = *enableAxis5P

}
