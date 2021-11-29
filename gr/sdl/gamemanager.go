package sdl

import "github.com/dragonfax/gunroar/gr/sdl/screen"

/**
 * Manage the lifecycle of the game.
 */

type GameManager interface {
	Init()
	Start()
	Move()
	Draw()

	SetMainLoop(*MainLoop)
	SetUIs(screen.Screen, Input)
	SetPrefManager(PrefManager)
}

type GameManagerBase struct {
	mainLoop        *MainLoop
	abstScreen      screen.Screen
	input           Input
	abstPrefManager PrefManager
}

func NewGameManagerBaseInternal() GameManagerBase {
	return GameManagerBase{}
}

func (this *GameManagerBase) SetMainLoop(mainLoop *MainLoop) {
	this.mainLoop = mainLoop
}

func (this *GameManagerBase) SetUIs(screen screen.Screen, input Input) {
	this.abstScreen = screen
	this.input = input
}

func (this *GameManagerBase) SetPrefManager(prefManager PrefManager) {
	this.abstPrefManager = prefManager
}
