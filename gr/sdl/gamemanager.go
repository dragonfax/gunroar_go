package sdl

/**
 * Manage the lifecycle of the game.
 */

type GameManager interface {
	init()
	start()
	close()
	move()
	draw()

	setMainLoop(*MainLoop)
	setUIs(Screen, Input)
	setPrefManager(PrefManager)
}

type GameManagerBase struct {
	MainLoop    mainLoop
	Screen      abstScreen
	Input       input
	PrefManager abstPrefManager
}

func NewGameManagerBaseInternal() GameManagerBase {
	return GameManagerBase{}
}

func (this *GameManagerBase) setMainLoop(mainLoop MainLoop) {
	this.mainLoop = mainLoop
}

func (this *GameManagerBase) setUIs(screen Screen, input Input) {
	this.abstScreen = screen
	this.input = input
}

func (this *GameManagerBase) setPrefManager(prefManager PrefManager) {
	this.abstPrefManager = prefManager
}
