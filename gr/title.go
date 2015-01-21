/*
 * $Id: title.d,v 1.4 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/go-gl/gl"
)

/**
 * Title screen.
 */

const TITLE_SCROLL_SPEED_BASE = 0.025

type TitleManager struct {
	displayList   *DisplayList
	logo          *Texture
	cnt           int
	btnPressedCnt int
	gameMode      GameMode
}

func NewTitleManager() *TitleManager {
	tm := new(TitleManager)

	tm.logo = NewTextureFromBMP("title.bmp")
	tm.displayList = NewDisplayList(1)
	tm.displayList.beginSingleList()
	gl.Enable(gl.TEXTURE_2D)
	tm.logo.Bind(0)
	setScreenColor(1, 1, 1, 1)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(0, -63)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(255, -63)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(255, 0)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(0, 0)
	gl.End()
	lineWidth(3)
	gl.Disable(gl.TEXTURE_2D)
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(-80, -7)
	gl.Vertex2f(-20, -7)
	gl.Vertex2f(10, -70)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(45, -2)
	gl.Vertex2f(-15, -2)
	gl.Vertex2f(-45, 61)
	gl.End()
	gl.Begin(gl.TRIANGLE_FAN)
	setScreenColor(1, 1, 1, 1)
	gl.Vertex2f(-19, -6)
	setScreenColor(0, 0, 0, 1)
	gl.Vertex2f(-79, -6)
	gl.Vertex2f(11, -69)
	gl.End()
	gl.Begin(gl.TRIANGLE_FAN)
	setScreenColor(1, 1, 1, 1)
	gl.Vertex2f(-16, -3)
	setScreenColor(0, 0, 0, 1)
	gl.Vertex2f(44, -3)
	gl.Vertex2f(-46, 60)
	gl.End()
	lineWidth(1)
	tm.displayList.endSingleList()
	return tm
}

func (this *TitleManager) close() {
	this.displayList.close()
	this.logo.Close()
}

func (this *TitleManager) start() {
	this.cnt = 0
	field.start()
	this.btnPressedCnt = 1
}

func (this *TitleManager) move() {
	field.move()
	field.scroll(TITLE_SCROLL_SPEED_BASE, true)
	input := pad.getState()
	mouseInput := mouse.getState()
	if this.btnPressedCnt <= 0 {
		if ((input.button&PadButtonA != 0) || (this.gameMode == GameModeMOUSE && (mouseInput.button&MouseButtonLEFT != 0))) && this.gameMode >= 0 {
			gameManager.startInGame(this.gameMode)
		}
		gmc := GameMode(0)
		if (input.button&PadButtonB != 0) || (input.dir&PadDirDOWN != 0) {
			gmc = 1
		} else if input.dir&PadDirUP != 0 {
			gmc = -1
		}
		if gmc != 0 {
			this.gameMode += gmc
			if this.gameMode >= GAME_MODE_NUM {
				this.gameMode = -1
			} else if this.gameMode < -1 {
				this.gameMode = GAME_MODE_NUM - 1
			}
			fadeBgm()
			disableBgm()
			disableSe()
		}
	}
	if (input.button&(PadButtonA|PadButtonB) != 0) || (input.dir&(PadDirUP|PadDirDOWN) != 0) || (mouseInput.button&MouseButtonLEFT != 0) {
		this.btnPressedCnt = 6
	} else {
		this.btnPressedCnt--
	}
	this.cnt++
}

func (this *TitleManager) draw() {
	if this.gameMode < 0 {
		drawString("REPLAY", 3, 400, 5)
		return
	}
	var ts float32 = 1
	if this.cnt > 120 {
		ts -= (float32(this.cnt) - 120) * 0.015
		if ts < 0.5 {
			ts = 0.5
		}
	}
	gl.PushMatrix()
	gl.Translatef(80*ts, 240, 0)
	gl.Scalef(ts, ts, 0)
	this.displayList.call(0)
	gl.PopMatrix()
	if this.cnt > 200 {
		drawStringOption("LAST", 3, 345, 4, TO_RIGHT, 1, false, 0)
		ls := 0
		drawNumOption(ls, 80, 360, 4, 0, 9, -1, -1)
	}
	drawString(gameModeText[this.gameMode], 3, 400, 5)
}
