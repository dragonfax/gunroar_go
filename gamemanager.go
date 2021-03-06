/*
 * $Id: gamemanager.d,v 1.5 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var state GameState
var showFps bool

type GameState interface {
	start()
	move()
	draw()
	drawFront()
	drawOrtho()
	drawLuminous()
}

/**
 * Manage the game state
 */

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
	j, err := pad.openJoystick(nil)
	if err != nil {
		panic(err)
	}
	twinStick.openJoystick(j)
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

func (this *GameManager) startTitle() {
	state = titleState
	state.start()
}

func (this *GameManager) startInGame(gameMode GameMode) {
	state = inGameState
	inGameState.gameMode = gameMode
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
	if showFps {
		limiter.draw()
	}
	viewPerspective()
}
