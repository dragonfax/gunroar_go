/*
 * $Id: field.d,v 1.3 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/go-gl/gl"
)

var field *Field

type PlatformPos struct {
	pos  Vector
	deg  float32
	used bool
}

type Panel struct {
	x, y, z    float32
	ci         int
	or, og, ob float32
}

const BLOCK_SIZE_X = 20
const BLOCK_SIZE_Y = 64
const ON_BLOCK_THRESHOLD = 1
const NEXT_BLOCK_AREA_SIZE = 16
const SIDEWALL_X1 = 18
const SIDEWALL_X2 = 9.3
const SIDEWALL_Y = 15
const TIME_COLOR_INDEX = 5
const TIME_CHANGE_RATIO = 0.00033
const SCREEN_BLOCK_SIZE_X = 20
const SCREEN_BLOCK_SIZE_Y = 24
const BLOCK_WIDTH = 1
const PANEL_WIDTH = 1.8
const PANEL_HEIGHT_BASE = 0.66

var baseColorTime = [5][6][3]float32{
	[6][3]float32{[3]float32{0.15, 0.15, 0.3}, [3]float32{0.25, 0.25, 0.5}, [3]float32{0.35, 0.35, 0.45}, [3]float32{0.6, 0.7, 0.35}, [3]float32{0.45, 0.8, 0.3}, [3]float32{0.2, 0.6, 0.1}},
	[6][3]float32{[3]float32{0.1, 0.1, 0.3}, [3]float32{0.2, 0.2, 0.5}, [3]float32{0.3, 0.3, 0.4}, [3]float32{0.5, 0.65, 0.35}, [3]float32{0.4, 0.7, 0.3}, [3]float32{0.1, 0.5, 0.1}},
	[6][3]float32{[3]float32{0.1, 0.1, 0.3}, [3]float32{0.2, 0.2, 0.5}, [3]float32{0.3, 0.3, 0.4}, [3]float32{0.5, 0.65, 0.35}, [3]float32{0.4, 0.7, 0.3}, [3]float32{0.1, 0.5, 0.1}},
	[6][3]float32{[3]float32{0.2, 0.15, 0.25}, [3]float32{0.35, 0.2, 0.4}, [3]float32{0.5, 0.35, 0.45}, [3]float32{0.7, 0.6, 0.3}, [3]float32{0.6, 0.65, 0.25}, [3]float32{0.2, 0.45, 0.1}},
	[6][3]float32{[3]float32{0.0, 0.0, 0.1}, [3]float32{0.1, 0.1, 0.3}, [3]float32{0.2, 0.2, 0.3}, [3]float32{0.2, 0.3, 0.15}, [3]float32{0.2, 0.2, 0.1}, [3]float32{0.0, 0.15, 0.0}},
}

/**
 * Game field.
 */
type Field struct {
	size, outerSize         Vector
	block                   [][]int   /* BLOCK_SIZE_Y x BLOCK_SIZE_X */
	panel                   [][]Panel /* BLOCK_SIZE_Y x BLOCK_SIZE_X */
	nextBlockY              int
	screenY, blockCreateCnt float32
	lastScrollY             float32
	screenPos               Vector
	platformPos             []PlatformPos /* SCREEN_BLOCK_SIZE_X * NEXT_BLOCK_AREA_SIZE */
	platformPosNum          int
	baseColor               [6][3]float32
	time                    float32
}

func NewField() *Field {
	this := new(Field)
	this.size = Vector{SCREEN_BLOCK_SIZE_X / 2 * 0.9, SCREEN_BLOCK_SIZE_Y / 2 * 0.8}
	this.outerSize = Vector{SCREEN_BLOCK_SIZE_X / 2, SCREEN_BLOCK_SIZE_Y / 2}
	this.block = make([][]int, BLOCK_SIZE_X, BLOCK_SIZE_X)
	this.panel = make([][]Panel, BLOCK_SIZE_X, BLOCK_SIZE_X)
	for i, _ := range this.block {
		this.block[i] = make([]int, BLOCK_SIZE_Y, BLOCK_SIZE_Y)
		this.panel[i] = make([]Panel, BLOCK_SIZE_Y, BLOCK_SIZE_Y)
	}
	this.platformPos = make([]PlatformPos, SCREEN_BLOCK_SIZE_X*NEXT_BLOCK_AREA_SIZE, SCREEN_BLOCK_SIZE_X*NEXT_BLOCK_AREA_SIZE)
	/* for i, _ := range this.platformPos {
		this.platformPos[i].pos = Vector{}
	} */
	return this
}

func (this *Field) start() {
	this.lastScrollY = 0
	this.nextBlockY = 0
	this.screenY = NEXT_BLOCK_AREA_SIZE
	this.blockCreateCnt = 0
	for y := 0; y < BLOCK_SIZE_Y; y++ {
		for x := 0; x < BLOCK_SIZE_X; x++ {
			this.block[x][y] = -3
			this.createPanel(x, y)
		}
	}
	this.time = nextFloat(TIME_COLOR_INDEX)
}

func (this *Field) createPanel(x int, y int) {
	p := &(this.panel[x][y])
	p.x = nextFloat(1) - 0.75
	p.y = nextFloat(1) - 0.75
	p.z = float32(this.block[x][y])*PANEL_HEIGHT_BASE + nextFloat(PANEL_HEIGHT_BASE)
	p.ci = this.block[x][y] + 3
	p.or = 1 + nextSignedFloat(0.1)
	p.og = 1 + nextSignedFloat(0.1)
	p.ob = 1 + nextSignedFloat(0.1)
	p.or *= 0.33
	p.og *= 0.33
	p.ob *= 0.33
}

func (this *Field) scroll(my float32, isDemo bool /*= false*/) {
	this.lastScrollY = my
	this.screenY -= my
	if this.screenY < 0 {
		this.screenY += BLOCK_SIZE_Y
	}
	this.blockCreateCnt -= my
	if this.blockCreateCnt < 0 {
		stageManager.gotoNextBlockArea()
		var bd int
		if stageManager.bossMode {
			bd = 0
		} else {
			bd = stageManager.blockDensity
		}
		this.createBlocks(bd)
		if !isDemo {
			stageManager.addBatteries(this.platformPos, this.platformPosNum)
		}
		this.gotoNextBlockArea()
	}
}

func (this *Field) createBlocks(groundDensity int) {
	for y := this.nextBlockY; y < this.nextBlockY+NEXT_BLOCK_AREA_SIZE; y++ {
		by := y % BLOCK_SIZE_Y
		for bx := 0; bx < BLOCK_SIZE_X; bx++ {
			this.block[bx][by] = -3
		}
	}
	this.platformPosNum = 0
	groundType := nextInt(3)
	for i := 0; i < groundDensity; i++ {
		this.addGround(groundType)
	}
	for y := this.nextBlockY; y < this.nextBlockY+NEXT_BLOCK_AREA_SIZE; y++ {
		by := y % BLOCK_SIZE_Y
		for bx := 0; bx < BLOCK_SIZE_X; bx++ {
			if y == this.nextBlockY || y == this.nextBlockY+NEXT_BLOCK_AREA_SIZE-1 {
				this.block[bx][by] = -3
			}
		}
	}
	for y := this.nextBlockY; y < this.nextBlockY+NEXT_BLOCK_AREA_SIZE; y++ {
		by := y % BLOCK_SIZE_Y
		for bx := 0; bx < BLOCK_SIZE_X-1; bx++ {
			if this.block[bx][by] == 0 {
				if this.countAroundBlock(bx, by, 0) <= 1 {
					this.block[bx][by] = -2
				}
			}
		}
		for bx := BLOCK_SIZE_X - 1; bx >= 0; bx-- {
			if this.block[bx][by] == 0 {
				if this.countAroundBlock(bx, by, 0) <= 1 {
					this.block[bx][by] = -2
				}
			}
		}
		for bx := 0; bx < BLOCK_SIZE_X; bx++ {
			var b int
			c := this.countAroundBlock(bx, by, 0)
			if this.block[bx][by] >= 0 {
				switch c {
				case 0:
					b = -2
				case 1, 2, 3:
					b = 0
				case 4:
					b = 2
				}
			} else {
				switch c {
				case 0:
					b = -3
				case 1, 2, 3, 4:
					b = -1
				}
			}
			this.block[bx][by] = b
			if b == -1 && bx >= 2 && bx < BLOCK_SIZE_X-2 {
				pd := this.calcPlatformDeg(bx, by)
				if pd >= -Pi32*2 {
					this.platformPos[this.platformPosNum].pos.x = float32(bx)
					this.platformPos[this.platformPosNum].pos.y = float32(by)
					this.platformPos[this.platformPosNum].deg = pd
					this.platformPos[this.platformPosNum].used = false
					this.platformPosNum++
				}
			}
		}
	}
	for y := this.nextBlockY; y < this.nextBlockY+NEXT_BLOCK_AREA_SIZE; y++ {
		by := y % BLOCK_SIZE_Y
		for bx := 0; bx < BLOCK_SIZE_X; bx++ {
			if this.block[bx][by] == -3 {
				if this.countAroundBlock(bx, by, -1) > 0 {
					this.block[bx][by] = -2
				}
			} else if this.block[bx][by] == 2 {
				if this.countAroundBlock(bx, by, 1) < 4 {
					this.block[bx][by] = 1
				}
			}
			this.createPanel(bx, by)
		}
	}
}

func (this *Field) addGround(groundType int) {
	var cx int
	switch groundType {
	case 0:
		cx = nextInt(int(BLOCK_SIZE_X*0.4)) + int(BLOCK_SIZE_X*0.1)
	case 1:
		cx = nextInt(int(BLOCK_SIZE_X*0.4)) + int(BLOCK_SIZE_X*0.5)
	case 2:
		if nextInt(2) == 0 {
			cx = nextInt(int(BLOCK_SIZE_X*0.4)) - int(BLOCK_SIZE_X*0.2)
		} else {
			cx = nextInt(int(BLOCK_SIZE_X*0.4)) + int(BLOCK_SIZE_X*0.8)
		}
	}
	/* this crazy bit is required to get dlang to do what I want
	 * otherwise type and constant conversion and truncation errors, ahoy.
	 */
	cx1 := float32(NEXT_BLOCK_AREA_SIZE) * 0.6
	cx1r := nextInt(int(cx1))
	cx2 := float32(NEXT_BLOCK_AREA_SIZE) * 0.2
	cy := cx1r + int(cx2)
	cy += this.nextBlockY
	w1 := float32(BLOCK_SIZE_X) * 0.33
	w1r := nextInt(int(w1))
	w := w1r + int(w1)
	h1 := float32(NEXT_BLOCK_AREA_SIZE) * 0.24
	h1r := nextInt(int(h1))
	h2 := float32(NEXT_BLOCK_AREA_SIZE) * 0.33
	h := h1r + int(h2)
	cx -= w / 2
	cy -= h / 2
	var wr, hr float32
	for y := this.nextBlockY; y < this.nextBlockY+NEXT_BLOCK_AREA_SIZE; y++ {
		by := y % BLOCK_SIZE_Y
		for bx := 0; bx < BLOCK_SIZE_X; bx++ {
			if bx >= cx && bx < cx+w && y >= cy && y < cy+h {
				var o, to float32
				wr = nextFloat(0.2) + 0.2
				hr = nextFloat(0.3) + 0.4
				o = float32(bx-cx)*wr + float32(y-cy)*hr
				wr = nextFloat(0.2) + 0.2
				hr = nextFloat(0.3) + 0.4
				to = float32(cx+w-1-bx)*wr + float32(y-cy)*hr
				if to < o {
					o = to
				}
				wr = nextFloat(0.2) + 0.2
				hr = nextFloat(0.3) + 0.4
				to = float32(bx-cx)*wr + float32(cy+h-1-y)*hr
				if to < o {
					o = to
				}
				wr = nextFloat(0.2) + 0.2
				hr = nextFloat(0.3) + 0.4
				to = float32(cx+w-1-bx)*wr + float32(cy+h-1-y)*hr
				if to < o {
					o = to
				}
				if o > 1 {
					this.block[bx][by] = 0
				}
			}
		}
	}
}

func (this *Field) gotoNextBlockArea() {
	this.blockCreateCnt += NEXT_BLOCK_AREA_SIZE
	this.nextBlockY -= NEXT_BLOCK_AREA_SIZE
	if this.nextBlockY < 0 {
		this.nextBlockY += BLOCK_SIZE_Y
	}
}

func (this *Field) getBlockVector(p Vector) int {
	return this.getBlock(p.x, p.y)
}

func (this *Field) getBlock(x float32, y float32) int {
	y -= this.screenY - floor32(this.screenY)
	bx := int((x + BLOCK_WIDTH*SCREEN_BLOCK_SIZE_X/2) / BLOCK_WIDTH)
	by := int(this.screenY) + int((-y+BLOCK_WIDTH*SCREEN_BLOCK_SIZE_Y/2)/BLOCK_WIDTH)
	if bx < 0 || bx >= BLOCK_SIZE_X {
		return -1
	}
	if by < 0 {
		by += BLOCK_SIZE_Y
	} else if by >= BLOCK_SIZE_Y {
		by -= BLOCK_SIZE_Y
	}
	return this.block[bx][by]
}

func (this *Field) convertToScreenPos(bx int, y int) Vector {
	oy := this.screenY - this.screenY
	by := y - int(this.screenY)
	if by <= -BLOCK_SIZE_Y {
		by += BLOCK_SIZE_Y
	}
	if by > 0 {
		by -= BLOCK_SIZE_Y
	}
	this.screenPos.x = float32(bx*BLOCK_WIDTH - BLOCK_WIDTH*SCREEN_BLOCK_SIZE_X/2 + BLOCK_WIDTH/2)
	this.screenPos.y = float32(by*-BLOCK_WIDTH + BLOCK_WIDTH*SCREEN_BLOCK_SIZE_Y/2 + int(oy) - BLOCK_WIDTH/2)
	return this.screenPos
}

func (this *Field) move() {
	this.time += TIME_CHANGE_RATIO
	if this.time >= TIME_COLOR_INDEX {
		this.time -= TIME_COLOR_INDEX
	}
}

func (this *Field) draw() {
	this.drawPanel()
}

func (this *Field) drawSideWalls() {
	gl.Disable(gl.BLEND)
	setScreenColor(0, 0, 0, 1)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex3f(SIDEWALL_X1, SIDEWALL_Y, 0)
	gl.Vertex3f(SIDEWALL_X2, SIDEWALL_Y, 0)
	gl.Vertex3f(SIDEWALL_X2, -SIDEWALL_Y, 0)
	gl.Vertex3f(SIDEWALL_X1, -SIDEWALL_Y, 0)
	gl.End()
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex3f(-SIDEWALL_X1, SIDEWALL_Y, 0)
	gl.Vertex3f(-SIDEWALL_X2, SIDEWALL_Y, 0)
	gl.Vertex3f(-SIDEWALL_X2, -SIDEWALL_Y, 0)
	gl.Vertex3f(-SIDEWALL_X1, -SIDEWALL_Y, 0)
	gl.End()
	gl.Enable(gl.BLEND)
}

func (this *Field) drawPanel() {
	ci := this.time
	nci := ci + 1
	if nci >= TIME_COLOR_INDEX {
		nci = 0
	}
	var co float32 = this.time - ci
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			this.baseColor[i][j] = baseColorTime[int(ci)][i][j]*(1-co) + baseColorTime[int(nci)][i][j]*co
		}
	}
	var by int = int(this.screenY)
	var oy float32 = this.screenY - float32(by)
	var sx float32
	var sy float32 = BLOCK_WIDTH*SCREEN_BLOCK_SIZE_Y/2 + oy
	by--
	if by < 0 {
		by += BLOCK_SIZE_Y
	}
	sy += BLOCK_WIDTH
	gl.Begin(gl.QUADS)
	for y := -1; y < SCREEN_BLOCK_SIZE_Y+NEXT_BLOCK_AREA_SIZE; y++ {
		if by >= BLOCK_SIZE_Y {
			by -= BLOCK_SIZE_Y
		}
		sx = -BLOCK_WIDTH * SCREEN_BLOCK_SIZE_X / 2
		for bx := 0; bx < SCREEN_BLOCK_SIZE_X; bx++ {
			p := &(this.panel[bx][int(by)])
			setScreenColor(this.baseColor[p.ci][0]*p.or*0.66,
				this.baseColor[p.ci][1]*p.og*0.66,
				this.baseColor[p.ci][2]*p.ob*0.66, 1)
			gl.Vertex3f(sx+p.x, sy-p.y, p.z)
			gl.Vertex3f(sx+p.x+PANEL_WIDTH, sy-p.y, p.z)
			gl.Vertex3f(sx+p.x+PANEL_WIDTH, sy-p.y-PANEL_WIDTH, p.z)
			gl.Vertex3f(sx+p.x, sy-p.y-PANEL_WIDTH, p.z)
			setScreenColor(this.baseColor[p.ci][0]*0.33,
				this.baseColor[p.ci][1]*0.33,
				this.baseColor[p.ci][2]*0.33, 1)
			gl.Vertex2f(sx, sy)
			gl.Vertex2f(sx+BLOCK_WIDTH, sy)
			gl.Vertex2f(sx+BLOCK_WIDTH, sy-BLOCK_WIDTH)
			gl.Vertex2f(sx, sy-BLOCK_WIDTH)
			sx += BLOCK_WIDTH
		}
		sy -= BLOCK_WIDTH
		by++
	}
	gl.End()
}

var degBlockOfs = [4][2]int{[2]int{0, -1}, [2]int{1, 0}, [2]int{0, 1}, [2]int{-1, 0}}

func (this *Field) calcPlatformDeg(x int, y int) float32 {
	d := nextInt(4)
	for i := 0; i < 4; i++ {
		if !this.checkBlock(x+degBlockOfs[d][0], y+degBlockOfs[d][1], -1, true) {
			pd := float32(d) * Pi32 / 2
			ox := x + degBlockOfs[d][0]
			oy := y + degBlockOfs[d][1]
			td := d
			td--
			if td < 0 {
				td = 3
			}
			b1 := this.checkBlock(ox+degBlockOfs[td][0], oy+degBlockOfs[td][1], -1, true)
			td = d
			td++
			if td >= 4 {
				td = 0
			}
			b2 := this.checkBlock(ox+degBlockOfs[td][0], oy+degBlockOfs[td][1], -1, true)
			if !b1 && b2 {
				pd -= Pi32 / 4
			}
			if b1 && !b2 {
				pd += Pi32 / 4
			}
			pd = normalizeDeg(pd)
			return pd
		}
		d++
		if d >= 4 {
			d = 0
		}
	}
	return -99999
}

func (this *Field) countAroundBlock(x int, y int, th int /*= 0*/) int {
	c := 0
	if this.checkBlock(x, y-1, th, false) {
		c++
	}
	if this.checkBlock(x+1, y, th, false) {
		c++
	}
	if this.checkBlock(x, y+1, th, false) {
		c++
	}
	if this.checkBlock(x-1, y, th, false) {
		c++
	}
	return c
}

func (this *Field) checkBlock(x int, y int, th int /*= 0*/, outScreen bool /*= false*/) bool {
	if x < 0 || x >= BLOCK_SIZE_X {
		return outScreen
	}
	by := y
	if by < 0 {
		by += BLOCK_SIZE_Y
	}
	if by >= BLOCK_SIZE_Y {
		by -= BLOCK_SIZE_Y
	}
	return (this.block[x][by] >= th)
}

func (this *Field) checkInFieldVector(p Vector) bool {
	return this.size.containsVector(p, 1)
}

func (this *Field) checkInField(x float32, y float32) bool {
	return this.size.contains(x, y, 1)
}

func (this *Field) checkInOuterFieldVector(p Vector) bool {
	return this.outerSize.containsVector(p, 1)
}

func (this *Field) checkInOuterField(x float32, y float32) bool {
	return this.outerSize.contains(x, y, 1)
}

func (this *Field) checkInOuterHeightField(p Vector) bool {
	return p.x >= -this.size.x && p.x <= this.size.x && p.y >= -this.outerSize.y && p.y <= this.outerSize.y
}

func (this *Field) checkInFieldExceptTop(p Vector) bool {
	return p.x >= -this.size.x && p.x <= this.size.x && p.y >= -this.size.y
}

func (this *Field) checkInOuterFieldExceptTop(p Vector) bool {
	return p.x >= -this.outerSize.x && p.x <= this.outerSize.x && p.y >= -this.outerSize.y && p.y <= this.outerSize.y*2
}
