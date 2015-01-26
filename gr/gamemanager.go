/*
 * $Id: gamemanager.d,v 1.5 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/go-gl/gl"
	"github.com/veandco/go-sdl2/sdl"
)

/**
 * Manage the game state
 */

var shipTurnSpeed float32 = 1
var shipReverseFire bool

var field *Field
var pad *Pad

var twinStick *TwinStick
var mouse *Mouse
var mouseAndPad *MouseAndPad
var screen *Screen
var ship *Ship
var stageManager *StageManager
var titleManager *TitleManager
var scoreReel *ScoreReel
var state GameState
var titleState *TitleState
var inGameState *InGameState

var gameManager *GameManager

type GameManager struct {
	escPressed bool
}

func NewGameManager() *GameManager {
	return new(GameManager)
}

func (this *GameManager) init() {
	InitLetter()
	initShots()
	InitBulletShapes()
	InitEnemyShapes()
	// InitTurret.init()
	InitTurretShapes()
	InitBoats()
	InitShip()
	InitFragments()
	InitSparkFragments()
	InitCrystalShape()
	// InitCrystal.init()
	// twinStick = cast(TwinStick) (cast(MultipleInputDevice) input).inputs[1]
	twinStick.openJoystick(pad.openJoystick(nil))
	field = NewField()
	ship = NewShip()
	scoreReel = NewScoreReel()
	stageManager = NewStageManager()
	loadSounds()
	titleManager = NewTitleManager()
	inGameState = NewInGameState()
	titleState = NewTitleState()
}

func (this *GameManager) close() {
	ship.close()
	closeBoats()
	closeBulletShapes()
	closeEnemyShapes()
	closeTurretShapes()
	CloseFragments()
	CloseSparkFragments()
	CloseCrystalShape()
	// CloseCrystal()
	titleState.close()
	CloseLetter()
}

func (this *GameManager) start() {
	this.startTitle()
}

func (this *GameManager) startTitle() {
	state = titleState
	this.startState()
}

func (this *GameManager) startInGame(gameMode GameMode) {
	state = inGameState
	inGameState.gameMode = gameMode
	this.startState()
}

func (this *GameManager) startState() {
	state.start()
}

func (this *GameManager) move() {
	if pad.keys[sdl.SCANCODE_ESCAPE] == sdl.PRESSED {
		if !this.escPressed {
			this.escPressed = true
			if state == inGameState {
				this.startTitle()
			} else {
				mainLoop.done = true
			}
			return
		}
	} else {
		this.escPressed = false
	}
	state.move()
}

func (this *GameManager) draw() {
	if screen.startRenderToLuminousScreen() {
		gl.PushMatrix()
		screen.setEyepos()
		state.drawLuminous()
		gl.PopMatrix()
		screen.endRenderToLuminousScreen()
	}
	screen.clear()
	gl.PushMatrix()
	screen.setEyepos()
	state.draw()
	gl.PopMatrix()
	screen.drawLuminous()
	gl.PushMatrix()
	screen.setEyepos()
	field.drawSideWalls()
	state.drawFront()
	gl.PopMatrix()
	viewOrthoFixed()
	state.drawOrtho()
	viewPerspective()
}

/**
 * Manage the game state.
 * (e.g. title, in game, gameover, pause, ...)
 */
type GameState interface {
	start()
	move()
	draw()
	drawFront()
	drawOrtho()
	drawLuminous()
}

type GameMode int

const (
	GameModeNORMAL GameMode = iota
	GameModeTWIN_STICK
	GameModeDOUBLE_PLAY
	GameModeMOUSE
)

const GAME_MODE_NUM = 4

var gameModeText []string = []string{"NORMAL", "TWIN STICK", "DOUBLE PLAY", "MOUSE"}
var isGameOver bool

const SCORE_REEL_SIZE_DEFAULT = 0.5
const SCORE_REEL_SIZE_SMALL = 0.01

type InGameState struct {
	time, gameOverCnt int
	btnPressed        bool
	pauseCnt          int
	pausePressed      bool
	scoreReelSize     float32
	gameMode          GameMode
}

func NewInGameState() *InGameState {

	this := new(InGameState)
	this.scoreReelSize = SCORE_REEL_SIZE_DEFAULT
	return this
}

func (this *InGameState) start() {
	enableBgm()
	enableSe()
	this.startInGame()
}

func (this *InGameState) startInGame() {
	clearActors()
	field = NewField()
	ship = NewShip()
	scoreReel = NewScoreReel()
	stageManager = NewStageManager()
	stageManager.start(1)
	field.start()
	ship.start(this.gameMode)
	this.initGameState()
	screen.setScreenShake(0, 0)
	this.gameOverCnt = 0
	this.pauseCnt = 0
	this.scoreReelSize = SCORE_REEL_SIZE_DEFAULT
	isGameOver = false
	playBgm()
}

func (this *InGameState) initGameState() {
	this.time = 0
	scoreReel.clear(9)
	InitTargetY()
}

func (this *InGameState) move() {
	if pad.keys[sdl.SCANCODE_P] == sdl.PRESSED {
		if !this.pausePressed {
			if this.pauseCnt <= 0 && !isGameOver {
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
		input := pad.getState()
		mouseInput := mouse.getState()
		if (input.button&PadButtonA) != 0 || (this.gameMode == GameModeMOUSE && (mouseInput.button&MouseButtonLEFT) != 0) {
			if this.gameOverCnt > 60 && !this.btnPressed {
				gameManager.startTitle()
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
			gameManager.startTitle()
		}
	}
}

func (this *InGameState) moveInGame() {
	field.move()
	ship.move()
	stageManager.move()
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
	screen.move()
	this.scoreReelSize += (SCORE_REEL_SIZE_DEFAULT - this.scoreReelSize) * 0.05
	scoreReel.move()
	if !isGameOver {
		this.time += 17
	}
	playMarkedSe()
}

func enemiesMove() {
	for a, _ := range actors {
		e, ok := a.(*Enemy)
		if ok {
			e.move()
		}
	}
}

func shotsMove() {
	for a, _ := range actors {
		s, ok := a.(*Shot)
		if ok {
			s.move()
		}
	}
}

func bulletsMove() {
	for a, _ := range actors {
		b, ok := a.(*Bullet)
		if ok {
			b.move()
		}
	}
}

func crystalsMove() {
	for a, _ := range actors {
		b, ok := a.(*Crystal)
		if ok {
			b.move()
		}
	}
}

func numIndicatorsMove() {
	for a, _ := range actors {
		b, ok := a.(*NumIndicator)
		if ok {
			b.move()
		}
	}
}

func sparksMove() {
	for a, _ := range actors {
		b, ok := a.(*Spark)
		if ok {
			b.move()
		}
	}
}

func smokesMove() {
	for a, _ := range actors {
		b, ok := a.(*Smoke)
		if ok {
			b.move()
		}
	}
}

func fragmentsMove() {
	for a, _ := range actors {
		b, ok := a.(*Fragment)
		if ok {
			b.move()
		}
	}
}

func sparkFragmentsMove() {
	for a, _ := range actors {
		b, ok := a.(*SparkFragment)
		if ok {
			b.move()
		}
	}
}

func wakesMove() {
	for a, _ := range actors {
		b, ok := a.(*Wake)
		if ok {
			b.move()
		}
	}
}

func (this *InGameState) draw() {
	field.draw()
	gl.Begin(gl.TRIANGLES)
	wakesDraw()
	sparksDraw()
	gl.End()
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Begin(gl.QUADS)
	smokesDraw()
	gl.End()
	fragmentsDraw()
	sparkFragmentsDraw()
	crystalsDraw()
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	enemiesDraw()
	shotsDraw()
	ship.draw()
	bulletsDraw()
}

func wakesDraw() {
	for a, _ := range actors {
		b, ok := a.(*Wake)
		if ok {
			b.draw()
		}
	}
}

func sparksDraw() {
	for a, _ := range actors {
		b, ok := a.(*Spark)
		if ok {
			b.draw()
		}
	}
}

func smokesDraw() {
	for a, _ := range actors {
		b, ok := a.(*Smoke)
		if ok {
			b.draw()
		}
	}
}

func sparkFragmentsDraw() {
	for a, _ := range actors {
		b, ok := a.(*SparkFragment)
		if ok {
			b.draw()
		}
	}
}

func fragmentsDraw() {
	for a, _ := range actors {
		b, ok := a.(*Fragment)
		if ok {
			b.draw()
		}
	}
}

func crystalsDraw() {
	for a, _ := range actors {
		b, ok := a.(*Crystal)
		if ok {
			b.draw()
		}
	}
}

func enemiesDraw() {
	for a, _ := range actors {
		b, ok := a.(*Enemy)
		if ok {
			b.draw()
		}
	}
}

func shotsDraw() {
	for a, _ := range actors {
		b, ok := a.(*Shot)
		if ok {
			b.draw()
		}
	}
}

func bulletsDraw() {
	for a, _ := range actors {
		b, ok := a.(*Bullet)
		if ok {
			b.draw()
		}
	}
}

func numIndicatorsDraw() {
	for a, _ := range actors {
		b, ok := a.(*NumIndicator)
		if ok {
			b.draw()
		}
	}
}

func (this *InGameState) drawFront() {
	ship.drawFront()
	scoreReel.drawAtPos(11.5+(SCORE_REEL_SIZE_DEFAULT-this.scoreReelSize)*3,
		-8.2-(SCORE_REEL_SIZE_DEFAULT-this.scoreReelSize)*3,
		this.scoreReelSize)
	var x float32 = -12
	for i := 0; i < ship.livesLeft; i++ {
		gl.PushMatrix()
		gl.Translatef(x, -9, 0)
		gl.Scalef(0.7, 0.7, 0.7)
		ship.drawShape()
		gl.PopMatrix()
		x += 0.7
	}
	numIndicatorsDraw()
}

func (this *InGameState) drawGameParams() {
	stageManager.draw()
}

func (this *InGameState) drawOrtho() {
	this.drawGameParams()
	if isGameOver {
		drawString("GAME OVER", 190, 180, 15)
	}
	if this.pauseCnt > 0 && (this.pauseCnt%64) < 32 {
		drawString("PAUSE", 265, 210, 12)
	}
}

func (this *InGameState) drawLuminous() {
	gl.Begin(gl.TRIANGLES)
	sparksDrawLuminous()
	gl.End()
	sparkFragmentsDrawLuminous()
	gl.Begin(gl.QUADS)
	smokesDrawLuminous()
	gl.End()
}

func sparksDrawLuminous() {
	for a, _ := range actors {
		b, ok := a.(*Spark)
		if ok {
			b.drawLuminous()
		}
	}
}

func smokesDrawLuminous() {
	for a, _ := range actors {
		b, ok := a.(*Smoke)
		if ok {
			b.drawLuminous()
		}
	}
}

func sparkFragmentsDrawLuminous() {
	for a, _ := range actors {
		b, ok := a.(*SparkFragment)
		if ok {
			b.drawLuminous()
		}
	}
}

func (this *InGameState) shipDestroyed() {
	clearBullets()
	stageManager.shipDestroyed()
	limiter.initInterval()
	ship.livesLeft--
	if ship.livesLeft < 0 {
		isGameOver = true
		this.btnPressed = true
		fadeBgm()
		scoreReel.accelerate()
	}
}

func clearBullets() {
	for a, _ := range actors {
		b, ok := a.(*SparkFragment)
		if ok {
			b.close()
		}
	}
}

func (this *InGameState) shrinkScoreReel() {
	this.scoreReelSize += (SCORE_REEL_SIZE_SMALL - this.scoreReelSize) * 0.08
}

type TitleState struct {
	gameOverCnt int
}

func NewTitleState() *TitleState {
	this := new(TitleState)
	return this
}

func (this *TitleState) close() {
	titleManager.close()
}

func (this *TitleState) start() {
	haltBgm()
	disableBgm()
	disableSe()
	titleManager.start()
}

func (this *TitleState) move() {
	titleManager.move()
}

func (this *TitleState) draw() {
	field.draw()
}

func (this *TitleState) drawFront() {
}

func (this *TitleState) drawOrtho() {
	titleManager.draw()
}

func (this *TitleState) drawLuminous() {
	inGameState.drawLuminous()
}
