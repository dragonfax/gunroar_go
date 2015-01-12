package main

import (
	"github.com/dragonfax/gunroar_go/gr"
	"time"
)

func main() {

	var screen gr.Screen
	screen.InitSDL()
	gr.InitLetter()
	gr.InitNumIndicator()
	gr.InitTargetY()
	defer gr.CloseLetter()

	var reel gr.ScoreReel
	reel.Init()

	reel.AddActualScore(20)
	reel.Draw(0, 0, 0.5)

	screen.Flip()

	time.Sleep(10 * time.Second)

	screen.Close()

}
