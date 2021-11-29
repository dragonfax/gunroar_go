package main

import (
	"github.com/dragonfax/gunroar/gr/letter"
	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Title screen.
 */

const TITLE_SCROLL_SPEED_BASE = 0.025

type TitleManager struct {
	prefManager *PrefManager
	pad         *sdl.RecordablePad
	// RecordableMouse mouse;
	field         *Field
	gameManager   *GameManager
	displayList   *sdl.DisplayList
	logo          *sdl.Texture
	cnt           int
	_replayData   *ReplayData
	btnPressedCnt int
	gameMode      GameMode
}

func NewTitleManager(prefManager *PrefManager, pad *sdl.RecordablePad /* Mouse mouse, */, field *Field, gameManager *GameManager) *TitleManager {
	this := &TitleManager{}
	this.prefManager = prefManager
	// this.pad = cast(RecordablePad) pad;
	// this.mouse = cast(RecordableMouse) mouse;
	this.field = field
	this.pad = pad
	this.gameManager = gameManager
	this.init()
	return this
}

func (this *TitleManager) init() {
	this.logo = sdl.NewTexture("title.bmp")
	this.displayList = sdl.NewDisplayList(1)
	this.displayList.BeginNewList()
	gl.Enable(gl.TEXTURE_2D)
	this.logo.Bind(0)
	sdl.SetColor(1, 1, 1, 1)
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
	LineWidth(3)
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
	sdl.SetColor(1, 1, 1, 1)
	gl.Vertex2f(-19, -6)
	sdl.SetColor(0, 0, 0, 1)
	gl.Vertex2f(-79, -6)
	gl.Vertex2f(11, -69)
	gl.End()
	gl.Begin(gl.TRIANGLE_FAN)
	sdl.SetColor(1, 1, 1, 1)
	gl.Vertex2f(-16, -3)
	sdl.SetColor(0, 0, 0, 1)
	gl.Vertex2f(44, -3)
	gl.Vertex2f(-46, 60)
	gl.End()
	LineWidth(1)
	this.displayList.EndNewList()
	this.gameMode = prefManager.prefData().gameMode()
}

func (this *TitleManager) start() {
	this.cnt = 0
	this.field.start()
	this.btnPressedCnt = 1
}

func (this *TitleManager) move() {
	if this._replayData != nil {
		this.field.move()
		this.field.scroll(TITLE_SCROLL_SPEED_BASE, true)
	}
	input := this.pad.GetState(false)
	// MouseState mouseInput = mouse.getState(false);
	if this.btnPressedCnt <= 0 {
		if int(input.Button&sdl.ButtonA) > 0 && int(this.gameMode) >= 0 {
			/* (this.gameMode == MOUSE &&
			(this.mouseInput.button & MouseState.Button.LEFT)) ) && */
			this.gameManager.startInGame(this.gameMode)
		}
		gmc := 0
		if int(input.Button&sdl.ButtonB) > 0 || int(input.Dir&sdl.DOWN) > 0 {
			gmc = 1
		} else if int(input.Dir&sdl.UP) > 0 {
			gmc = -1
		}
		if gmc != 0 {
			this.gameMode = GameMode(int(this.gameMode) + gmc)
			if this.gameMode >= GAME_MODE_NUM {
				this.gameMode = -1
			} else if this.gameMode < -1 {
				this.gameMode = GAME_MODE_NUM - 1
			}
			if int(this.gameMode) == -1 && this._replayData != nil {
				enableBgm()
				enableSe()
				playCurrentBgm()
			} else {
				fadeBgm()
				disableBgm()
				disableSe()
			}
		}
	}
	if int(input.Button&(sdl.ButtonA|sdl.ButtonB)) > 0 ||
		int(input.Dir&(sdl.UP|sdl.DOWN)) > 0 { /* ||
		(mouseInput.button & MouseState.Button.LEFT) */
		this.btnPressedCnt = 6
	} else {
		this.btnPressedCnt--
	}
	this.cnt++
}

func (this *TitleManager) draw() {
	if this.gameMode < 0 {
		letter.DrawString("REPLAY", 3, 400, 5, letter.TO_RIGHT, 0, false, 0)
		return
	}
	ts := 1.0
	if this.cnt > 120 {
		ts -= float64(this.cnt-120) * 0.015
		if ts < 0.5 {
			ts = 0.5
		}
	}
	gl.PushMatrix()
	gl.Translated(80*ts, 240, 0)
	gl.Scaled(ts, ts, 0)
	this.displayList.Call(0)
	gl.PopMatrix()
	if this.cnt > 150 {
		letter.DrawString("HIGH", 3, 305, 4, letter.TO_RIGHT, 1, false, 0)
		letter.DrawNum(prefManager.prefData().highScore(this.gameMode), 80, 320, 4, 0, 9, -1, -1)
	}
	if this.cnt > 200 {
		letter.DrawString("LAST", 3, 345, 4, letter.TO_RIGHT, 1, false, 0)
		ls := 0
		if this._replayData != nil {
			ls = this._replayData.score
		}
		letter.DrawNum(ls, 80, 360, 4, 0, 9, -1, -1)
	}
	letter.DrawString(gameModeText[this.gameMode], 3, 400, 5, letter.TO_RIGHT, 0, false, 0)
}

func (this *TitleManager) replayData(v *ReplayData) *ReplayData {
	this._replayData = v
	return v
}
