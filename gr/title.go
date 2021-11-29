package main

import (
	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Title screen.
 */

const TITLE_SCROLL_SPEED_BASE = 0.025

type TitleManager struct {
	prefManager *PrefManager
	// RecordablePad pad;
	// RecordableMouse mouse;
	field         *Field
	gameManager   *GameManager
	displayList   *sdl.DisplayList
	logo          *sdl.Texture
	cnt           int
	_replayData   *ReplayData
	btnPressedCnt int
	gameMode      int
}

func NewTitleManager(prefManager *PrefManager /* Pad pad, Mouse mouse, */, field *Field, gameManager *GameManager) *TitleManager {
	this := &TitleManager{}
	this.prefManager = prefManager
	// this.pad = cast(RecordablePad) pad;
	// this.mouse = cast(RecordableMouse) mouse;
	this.field = field
	this.gameManager = gameManager
	this.init()
	return this
}

func (this *TitleManager) init() {
	this.logo = sdl.NewTexture("title.bmp")
	this.displayList = sdl.NewDisplayList(1)
	this.displayList.BeginNewList()
	gl.Enable(gl.TEXTURE_2D)
	this.logo.bind()
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
	sdl.lineWidth(3)
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
	sdl.SetColor(1, 1, 1)
	gl.Vertex2f(-19, -6)
	sdl.SetColor(0, 0, 0)
	gl.Vertex2f(-79, -6)
	gl.Vertex2f(11, -69)
	gl.End()
	gl.Begin(gl.TRIANGLE_FAN)
	sdl.SetColor(1, 1, 1)
	gl.Vertex2f(-16, -3)
	sdl.SetColor(0, 0, 0)
	gl.Vertex2f(44, -3)
	gl.Vertex2f(-46, 60)
	gl.End()
	sdl.lineWidth(1)
	this.displayList.endNewList()
	this.gameMode = prefManager.prefData.gameMode
}

func (this *TitleManager) start() {
	this.cnt = 0
	this.field.start()
	this.btnPressedCnt = 1
}

func (this *TitleManager) move() {
	if !this._replayData {
		this.field.move()
		this.field.scroll(SCROLL_SPEED_BASE, true)
	}
	// PadState input = pad.getState(false);
	// MouseState mouseInput = mouse.getState(false);
	if this.btnPressedCnt <= 0 {
		if ((this.input.button & PadState.Button.A) ||
			(this.gameMode == InGameState.GameMode.MOUSE &&
				(this.mouseInput.button & MouseState.Button.LEFT))) &&
			this.gameMode >= 0 {
			this.gameManager.startInGame(this.gameMode)
		}
		gmc := 0
		if (this.input.button & PadState.Button.B) || (this.input.dir & PadState.Dir.DOWN) {
			gmc = 1
		} else if this.input.dir & PadState.Dir.UP {
			gmc = -1
		}
		if gmc != 0 {
			this.gameMode += gmc
			if this.gameMode >= InGameState.GAME_MODE_NUM {
				this.gameMode = -1
			} else if this.gameMode < -1 {
				this.gameMode = InGameState.GAME_MODE_NUM - 1
			}
			if this.gameMode == -1 && this._replayData {
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
	if (this.input.button & (PadState.Button.A | PadState.Button.B)) ||
		(this.input.dir & (PadState.Dir.UP | PadState.Dir.DOWN)) ||
		(mouseInput.button & MouseState.Button.LEFT) {
		this.btnPressedCnt = 6
	} else {
		this.btnPressedCnt--
	}
	this.cnt++
}

func (this *TitleManager) draw() {
	if this.gameMode < 0 {
		letter.drawString("REPLAY", 3, 400, 5)
		return
	}
	ts := 1.0
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
	if this.cnt > 150 {
		letter.drawString("HIGH", 3, 305, 4, TO_RIGHT, 1)
		letter.drawNum(prefManager.prefData.highScore(gameMode), 80, 320, 4, 0, 9)
	}
	if this.cnt > 200 {
		letter.drawString("LAST", 3, 345, 4, TO_RIGHT, 1)
		ls := 0
		if this._replayData {
			ls = this._replayData.score
		}
		letter.drawNum(ls, 80, 360, 4, 0, 9)
	}
	letter.drawString(InGameState.gameModeText[this.gameMode], 3, 400, 5)
}

func (this *TitleManager) replayData(v ReplayData) ReplayData {
	this._replayData = v
	return v
}
