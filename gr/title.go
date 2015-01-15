/*
 * $Id: title.d,v 1.4 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Title screen.
 */

const TITLE_SCROLL_SPEED_BASE = 0.025

type TitleManager struct {
	pad           Pad
	mouse         Mouse
	field         Field
	gameManager   GameManager
	displayList   DisplayList
	logo          Texture
	cnt           int
	btnPressedCnt int
	gameMode      int
}

func NewTitleManager(pad Pad, mouse Mouse, field Field, gameManager GameManager) *TitleManager {
	tm := new(TitleManager)

	tm.pad = pad
	tm.mouse = mouse
	tm.field = field
	tm.gameManager = gameManager

	tm.logo = Texture.load("title.bmp")
	tm.displayList = NewDisplayList(1)
	tm.displayList.beginNewList()
	gl.Enable(GL_TEXTURE_2D)
	tm.logo.bind()
	setScreenColor(1, 1, 1)
	gl.Begin(GL_TRIANGLE_FAN)
	gl.TexCoord2(0, 0)
	gl.Vertex2(0, -63)
	gl.TexCoord2(1, 0)
	gl.Vertex2(255, -63)
	gl.TexCoord2(1, 1)
	gl.Vertex2(255, 0)
	gl.TexCoord2(0, 1)
	gl.Vertex2(0, 0)
	gl.End()
	lineWidth(3)
	gl.Disable(GL_TEXTURE_2D)
	gl.Begin(GL_LINE_STRIP)
	gl.Vertex2(-80, -7)
	gl.Vertex2(-20, -7)
	gl.Vertex2(10, -70)
	gl.End()
	gl.Begin(GL_LINE_STRIP)
	gl.Vertex2(45, -2)
	gl.Vertex2(-15, -2)
	gl.Vertex2(-45, 61)
	gl.End()
	gl.Begin(GL_TRIANGLE_FAN)
	setScreenColor(1, 1, 1)
	gl.Vertex2(-19, -6)
	setScreenColor(0, 0, 0)
	gl.Vertex2(-79, -6)
	gl.Vertex2(11, -69)
	gl.End()
	gl.Begin(GL_TRIANGLE_FAN)
	setScreenColor(1, 1, 1)
	gl.Vertex2(-16, -3)
	setScreenColor(0, 0, 0)
	gl.Vertex2(44, -3)
	gl.Vertex2(-46, 60)
	gl.End()
	lineWidth(1)
	tm.displayList.endNewList()
}

func (this *TitleManager) close() {
	this.displayList.close()
	this.logo.close()
}

func (this *TitleManager) start() {
	this.cnt = 0
	this.field.start()
	this.btnPressedCnt = 1
}

func (this *TitleManager) move() {
	this.field.move()
	this.field.scroll(TITLE_SCROLL_SPEED_BASE, true)
	input := this.pad.getState(false)
	mouseInput := this.mouse.getState(false)
	if this.btnPressedCnt <= 0 {
		if ((input.button & PadStateButto.A) || (gameMode == GameModeMOUSE && (mouseInput.button & MouseStateButtonLEFT))) && this.gameMode >= 0 {
			this.gameManager.startInGame(this.gameMode)
		}
		gmc := 0
		if (input.button & PadStateButtonB) || (input.dir & PadStateDirDOWN) {
			gmc = 1
		} else if input.dir & PadStateDirUP {
			gmc = -1
		}
		if gmc != 0 {
			this.gameMode += gmc
			if this.gameMode >= InGameState.GAME_MODE_NUM {
				this.gameMode = -1
			} else if this.gameMode < -1 {
				this.gameMode = InGameState.GAME_MODE_NUM - 1
			}
			fadeBgm()
			disableBgm()
			disableSe()
		}
	}
	if (input.button & (PadState.Button.A | PadState.Button.B)) || (input.dir & (PadState.Dir.UP | PadState.Dir.DOWN)) || (mouseInput.button & MouseState.Button.LEFT) {
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
		ts -= (this.cnt - 120) * 0.015
		if ts < 0.5 {
			ts = 0.5
		}
	}
	gl.PushMatrix()
	gl.Translatef(80*ts, 240, 0)
	gl.Scalef(ts, ts, 0)
	this.displayList.call()
	gl.PopMatrix()
	if this.cnt > 200 {
		drawString("LAST", 3, 345, 4, Letter.Direction.TO_RIGHT, 1)
		ls := 0
		drawNum(ls, 80, 360, 4, 0, 9)
	}
	drawString(InGameState.gameModeText[gameMode], 3, 400, 5)
}
