package sdl

import "github.com/dragonfax/gunroar/gr/sdl/screen"

/**
 * Manage the lifecycle of the game.
 */

type GameManager interface {
	init()
	start()
	move()
	draw()

	setMainLoop(*MainLoop)
	setUIs(screen.Screen, Input)
	setPrefManager(PrefManager)
}

type GameManagerBase struct {
	mainLoop        MainLoop
	abstScreen      screen.Screen
	input           Input
	abstPrefManager PrefManager
}

func NewGameManagerBaseInternal() GameManagerBase {
	return GameManagerBase{}
}

func (this *GameManagerBase) setMainLoop(mainLoop MainLoop) {
	this.mainLoop = mainLoop
}

func (this *GameManagerBase) setUIs(screen screen.Screen, input Input) {
	this.abstScreen = screen
	this.input = input
}

func (this *GameManagerBase) setPrefManager(prefManager PrefManager) {
	this.abstPrefManager = prefManager
}
