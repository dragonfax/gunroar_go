/*
 * $Id: gamemanager.d,v 1.5 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Manage the game state
 */

var shipTurnSpeed float32 = 1
var shipReverseFire bool = false

type GameManager struct {
	pad Pad
	// twinStick TwinStick
	mouse        Mouse
	mouseAndPad  MouseAndPad
	screen       Screen
	field        Field
	ship         Ship
	stageManager StageManager
	titleManager TitleManager
	scoreReel    ScoreReel
	state        GameState
	titleState   TitleState
	inGameState  InGameState
	escPressed   bool
}

func (this *GameManager) init() {
	InitLetter.init()
	InitShot.init()
	InitBulletShape.init()
	InitEnemyShape.init()
	InitTurret.init()
	InitTurretShape.init()
	InitFragment.init()
	InitSparkFragment.init()
	InitCrystal.init()
	this.pad = input.inputs[0]
	// twinStick = cast(TwinStick) (cast(MultipleInputDevice) input).inputs[1]
	// twinStick.openJoystick(pad.openJoystick())
	this.mouse = input.inputs[2]
	this.mouse.init(screen)
	this.mouseAndPad = NewMouseAndPad(mouse, pad)
	this.field = NewField()
	this.ship = NewShip(pad, twinStick, mouse, mouseAndPad, field, screen)
	this.scoreReel = NewScoreReel()
	this.stageManager = NewStageManager(field, enemies, ship)
	this.ship.setStageManager(stageManager)
	this.field.setStageManager(stageManager)
	this.field.setShip(ship)
	loadSounds()
	this.titleManager = NewTitleManager(pad, mouse, field, this)
	this.inGameState = NewInGameState(this, screen, pad /*twinStick, */, mouse, mouseAndPad,
		field, ship, stageManager, scoreReel)
	this.titleState = NewTitleState(this, screen, pad /*twinStick, */, mouse, mouseAndPad,
		field, ship, stageManager, scoreReel,
		titleManager, inGameState)
	this.ship.setGameState(this.inGameState)
}

func (this *GameManager) close() {
	this.ship.close()
	CloseBulletShape()
	CloseEnemyShape()
	CloseTurretShape()
	CloseFragment()
	CloseSparkFragment()
	CloseCrystal()
	this.titleState()
	CloseLetter()
}

func (this *GameManager) start() {
	this.startTitle()
}

func (this *GameManager) startTitle(fromGameover bool /*= false*/) {
	this.state = this.titleState
	this.startState()
}

func (this *GameManager) startInGame(gameMode GameMode) {
	this.state = this.inGameState
	this.inGameState.gameMode = gameMode
	this.startState()
}

func (this *GameManager) startState() {
	this.state.start()
}

func (this *GameManager) initInterval() {
	mainLoop.initInterval()
}

func (this *GameManager) addSlowdownRatio(sr float32) {
	mainLoop.addSlowdownRatio(sr)
}

func (this *GameManager) move() {
	if this.pad.keys[SDL.K_ESCAPE] == sdl.PRESSED {
		if !escPressed {
			this.escPressed = true
			if this.state == this.inGameState {
				this.startTitle()
			} else {
				mainLoop.breakLoop()
			}
			return
		}
	} else {
		this.escPressed = false
	}
	this.state.move()
}

func (this *GameManager) draw() {
	/*
			e := mainLoop.event
			switch (e.type) {
				case sdl.WindowEvent:
					re := e.resize
					if (re.w > 150 && re.h > 100) {
						this.screen.resized(re.w, re.h)
					}
			}
		 }
	*/
	if this.screen.startRenderToLuminousScreen() {
		glPushMatrix()
		this.screen.setEyepos()
		this.state.drawLuminous()
		glPopMatrix()
		this.screen.endRenderToLuminousScreen()
	}
	this.screen.clear()
	glPushMatrix()
	this.screen.setEyepos()
	this.state.draw()
	glPopMatrix()
	this.screen.drawLuminous()
	glPushMatrix()
	this.screen.setEyepos()
	this.field.drawSideWalls()
	this.state.drawFront()
	glPopMatrix()
	this.screen.viewOrthoFixed()
	this.state.drawOrtho()
	this.screen.viewPerspective()
}

/**
 * Manage the game state.
 * (e.g. title, in game, gameover, pause, ...)
 */
type GameState struct {
	gameManager GameManager
	screen      Screen
	pad         Pad
	// TwinStick twinStick
	mouse        Mouse
	mouseAndPad  MouseAndPad
	field        Field
	ship         Ship
	stageManager StageManager
	scoreReel    ScoreReel
}

func NewGameState(gameManager GameManager, screen Screen,
	pad Pad /*twinStick twinStick,*/, mouse Mouse, mouseAndPad MouseAndPad,
	field Field, ship Ship, stageManager StageManager, scoreReel ScoreReel) *GameState {
	this := new(GameState)
	this.gameManager = gameManager
	this.screen = screen
	this.pad = pad
	this.twinStick = twinStick
	this.mouse = mouse
	this.mouseAndPad = mouseAndPad
	this.field = field
	this.ship = ship
	this.stageManager = stageManager
	this.scoreReel = scoreReel
	return this
}

type GameMode int

const (
	GameModeNORMAL GameMode = iota
	GameModeTWIN_STICK
	GameModeDOUBLE_PLAY
	GameModeMOUSE
)

const GAME_MODE_NUM = 2

var gameModeText []string = []string{"NORMAL" /*"TWIN STICK",*/ /*"DOUBLE PLAY",*/, "MOUSE"}
var isGameOver bool

const SCORE_REEL_SIZE_DEFAULT = 0.5
const SCORE_REEL_SIZE_SMALL = 0.01

type InGameState struct {
	*GameState

	left, time, gameOverCnt int
	btnPressed              bool
	pauseCnt                int
	pausePressed            bool
	scoreReelSize           float32
	gameMode                GameMode
}

func NewInGameState(gameManager GameManager, screen Screen,
	pad Pad /*TwinStick twinStick,*/, mouse Mouse, mouseAndPad MouseAndPad,
	field Field, ship Ship, stageManager StageManager, scoreReel ScoreReel) *InGameState {

	this := InGameState{NewGameState(gameManager, screen, pad /*twinStick, */, mouse, mouseAndPad, field, ship, stageManager, scoreReel)}
	this.scoreReelSize = SCORE_REEL_SIZE_DEFAULT
	return this
}

func (this *InGameState) start() {
	enableBgm()
	enableSe()
	this.startInGame()
}

func (this *InGameState) startInGame() {
	this.clearAll()
	this.stageManager.start(1)
	this.field.start()
	this.ship.start(this.gameMode)
	this.initGameState()
	this.screen.setScreenShake(0, 0)
	this.gameOverCnt = 0
	this.pauseCnt = 0
	this.scoreReelSize = SCORE_REEL_SIZE_DEFAULT
	this.isGameOver = false
	playBgm()
}

func (this *InGameState) initGameState() {
	this.time = 0
	this.left = 2
	this.scoreReel.clear(9)
	initTargetY()
}

func (this *InGameState) move() {
	if this.pad.keys[SDL.K_p] == sdl.PRESSED {
		if !this.pausePressed {
			if this.pauseCnt <= 0 && !this.isGameOver {
				this.pauseCnt = 1
			} else {
				this.pauseCnt = 0
			}
		}
		this.pausePressed = true
	} else {
		this.pausePressed = false
	}
	if this.pauseCnt > 0 {
		this.pauseCnt++
		return
	}
	this.moveInGame()
	if isGameOver {
		this.gameOverCnt++
		input := pad.getState(false)
		mouseInput := mouse.getState(false)
		if (input.button & PadState.Button.A) || (gameMode == InGameState.GameMode.MOUSE && (mouseInput.button & MouseState.Button.LEFT)) {
			if this.gameOverCnt > 60 && !this.btnPressed {
				this.gameManager.startTitle(true)
			}
			this.btnPressed = true
		} else {
			this.btnPressed = false
		}
		if this.gameOverCnt == 120 {
			fadeBgm()
			disableBgm()
		}
		if this.gameOverCnt > 1200 {
			this.gameManager.startTitle(true)
		}
	}
}

func (this *InGameState) moveInGame() {
	this.field.move()
	this.ship.move()
	this.stageManager.move()
	enemiesMove()
	shotsMove()
	bulletsMove()
	crystalsMove()
	numIndicatorsMove()
	sparksMove()
	smokesMove()
	fragmentsMove()
	sparkFragmentsMove()
	wakesMove()
	this.screen.move()
	this.scoreReelSize += (SCORE_REEL_SIZE_DEFAULT - this.scoreReelSize) * 0.05
	this.scoreReel.move()
	if !this.isGameOver {
		this.time += 17
	}
	playMarkedSe()
}

func (this *InGameState) draw() {
	this.field.draw()
	glBegin(GL_TRIANGLES)
	wakesDraw()
	sparksDraw()
	glEnd()
	glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)
	glBegin(GL_QUADS)
	smokesDraw()
	glEnd()
	fragmentsDraw()
	sparkFragmentsDraw()
	crystalsDraw()
	glBlendFunc(GL_SRC_ALPHA, GL_ONE)
	enemiesDraw()
	shotsDraw()
	shipDraw()
	bulletsDraw()
}

func (this *InGameState) drawFront() {
	this.ship.drawFront()
	this.scoreReel.draw(11.5+(SCORE_REEL_SIZE_DEFAULT-this.scoreReelSize)*3,
		-8.2-(SCORE_REEL_SIZE_DEFAULT-this.scoreReelSize)*3,
		this.scoreReelSize)
	var x float32 = -12
	for i = 0; i < this.left; i++ {
		glPushMatrix()
		glTranslatef(x, -9, 0)
		glScalef(0.7, 0.7, 0.7)
		this.ship.drawShape()
		glPopMatrix()
		x += 0.7
	}
	numIndicatorsDraw()
}

func (this *InGameState) drawGameParams() {
	this.stageManager.draw()
}

func (this *InGameState) drawOrtho() {
	this.drawGameParams()
	if this.isGameOver {
		drawString("GAME OVER", 190, 180, 15)
	}
	if this.pauseCnt > 0 && (this.pauseCnt%64) < 32 {
		drawString("PAUSE", 265, 210, 12)
	}
}

func (this *InGameState) drawLuminous() {
	glBegin(GL_TRIANGLES)
	sparksDrawLuminous()
	glEnd()
	sparkFragmentsDrawLuminous()
	glBegin(GL_QUADS)
	smokesDrawLuminous()
	glEnd()
}

func (this *InGameState) shipDestroyed() {
	clearBullets()
	this.stageManager.shipDestroyed()
	this.gameManager.initInterval()
	this.left--
	if this.left < 0 {
		this.isGameOver = true
		this.btnPressed = true
		fadeBgm()
		this.scoreReel.accelerate()
	}
}

func (this *InGameState) shrinkScoreReel() {
	this.scoreReelSize += (SCORE_REEL_SIZE_SMALL - this.scoreReelSize) * 0.08
}

type TitleState struct {
	*GameState

	titleManager TitleManager
	inGameState  InGameState
	gameOverCnt  int
}

func NewTitleState(gameManager GameManager, screen Screen,
	pad Pad /*TwinStick twinStick, */, mouse Mouse, mouseAndPad MouseAndPad,
	field Field, ship Ship, stageManager StageManager, scoreReel ScoreReel,
	titleManager TitleManager, inGameState InGameState) *TitleState {

	this := TitleState{NewGameState(gameManager, screen, pad /*twinStick, */, mouse, mouseAndPad, field, ship, stageManager, scoreReel)}

	this.titleManager = titleManager
	this.inGameState = inGameState

	return this
}

func (this *TitleState) close() {
	this.titleManager.close()
}

func (this *TitleState) start() {
	haltBgm()
	disableBgm()
	disableSe()
	this.titleManager.start()
}

func (this *TitleState) move() {
	this.titleManager.move()
}

func (this *TitleState) draw() {
	this.field.draw()
}

func (this *TitleState) drawFront() {
}

func (this *TitleState) drawOrtho() {
	this.titleManager.draw()
}

func (this *TitleState) drawLuminous() {
	this.inGameState.drawLuminous()
}