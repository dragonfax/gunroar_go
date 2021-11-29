package main

import (
	"math"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Game field.
 */

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

type PlatformPos struct {
	pos  vector.Vector
	deg  float64
	used bool
}

type Panel struct {
	x, y, z    float64
	ci         int
	or, og, ob float64
}

var baseColorTime = [TIME_COLOR_INDEX][6][3]float64{
	{{0.15, 0.15, 0.3}, {0.25, 0.25, 0.5}, {0.35, 0.35, 0.45},
		{0.6, 0.7, 0.35}, {0.45, 0.8, 0.3}, {0.2, 0.6, 0.1}},
	{{0.1, 0.1, 0.3}, {0.2, 0.2, 0.5}, {0.3, 0.3, 0.4},
		{0.5, 0.65, 0.35}, {0.4, 0.7, 0.3}, {0.1, 0.5, 0.1}},
	{{0.1, 0.1, 0.3}, {0.2, 0.2, 0.5}, {0.3, 0.3, 0.4},
		{0.5, 0.65, 0.35}, {0.4, 0.7, 0.3}, {0.1, 0.5, 0.1}},
	{{0.2, 0.15, 0.25}, {0.35, 0.2, 0.4}, {0.5, 0.35, 0.45},
		{0.7, 0.6, 0.3}, {0.6, 0.65, 0.25}, {0.2, 0.45, 0.1}},
	{{0.0, 0.0, 0.1}, {0.1, 0.1, 0.3}, {0.2, 0.2, 0.3},
		{0.2, 0.3, 0.15}, {0.2, 0.2, 0.1}, {0.0, 0.15, 0.0}},
}

const PANEL_WIDTH = 1.8
const PANEL_HEIGHT_BASE = 0.66

type Field struct {
	stageManager                          *StageManager
	ship                                  *Ship
	rand                                  *r.Rand
	_size, _outerSize                     vector.Vector
	block                                 [BLOCK_SIZE_X][BLOCK_SIZE_Y]int
	panel                                 [BLOCK_SIZE_X][BLOCK_SIZE_Y]Panel
	nextBlockY                            int
	screenY, blockCreateCnt, _lastScrollY float64
	screenPos                             vector.Vector
	platformPos                           [SCREEN_BLOCK_SIZE_X * NEXT_BLOCK_AREA_SIZE]PlatformPos
	platformPosNum                        int
	baseColor                             [6][3]float64
	time                                  float64
}

func NewField() *Field {
	this := &Field{}
	this.rand = r.New(r.NewSource(time.Now().Unix()))
	this._size = vector.Vector{SCREEN_BLOCK_SIZE_X / 2 * 0.9, SCREEN_BLOCK_SIZE_Y / 2 * 0.8}
	this._outerSize = vector.Vector{SCREEN_BLOCK_SIZE_X / 2, SCREEN_BLOCK_SIZE_Y / 2}
	return this
}

func (this *Field) setRandSeed(s int64) {
	this.rand = r.New(r.NewSource(s))
}

func (this *Field) setStageManager(sm *StageManager) {
	this.stageManager = sm
}

func (this *Field) setShip(sp *Ship) {
	this.ship = sp
}

func (this *Field) start() {
	this._lastScrollY = 0
	this.nextBlockY = 0
	this.screenY = NEXT_BLOCK_AREA_SIZE
	this.blockCreateCnt = 0
	for y := 0; y < BLOCK_SIZE_Y; y++ {
		for x := 0; x < BLOCK_SIZE_X; x++ {
			this.block[x][y] = -3
			this.createPanel(x, y)
		}
	}
	this.time = nextFloat(this.rand, TIME_COLOR_INDEX)
}

func (this *Field) createPanel(x, y int) {
	p := &(this.panel[x][y])
	p.x = nextFloat(this.rand, 1) - 0.75
	p.y = nextFloat(this.rand, 1) - 0.75
	p.z = float64(this.block[x][y])*PANEL_HEIGHT_BASE + nextFloat(this.rand, PANEL_HEIGHT_BASE)
	p.ci = this.block[x][y] + 3
	p.or = 1 + nextSignedFloat(this.rand, 0.1)
	p.og = 1 + nextSignedFloat(this.rand, 0.1)
	p.ob = 1 + nextSignedFloat(this.rand, 0.1)
	p.or *= 0.33
	p.og *= 0.33
	p.ob *= 0.33
}

func (this *Field) scroll(my float64, isDemo bool /* = false */) {
	this._lastScrollY = my
	this.screenY -= my
	if this.screenY < 0 {
		this.screenY += BLOCK_SIZE_Y
	}
	this.blockCreateCnt -= my
	if this.blockCreateCnt < 0 {
		this.stageManager.gotoNextBlockArea()
		var bd int
		if this.stageManager.bossMode() {
			bd = 0
		} else {
			bd = this.stageManager.blockDensity()
		}
		this.createBlocks(bd)
		if !isDemo {
			this.stageManager.addBatteries(this.platformPos[:], this.platformPosNum)
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
	typ := this.rand.Intn(3)
	for i := 0; i < groundDensity; i++ {
		this.addGround(typ)
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
				if pd >= -math.Pi*2 {
					this.platformPos[this.platformPosNum].pos.X = float64(bx)
					this.platformPos[this.platformPosNum].pos.Y = float64(by)
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

func (this *Field) addGround(typ int) {
	var cx int
	switch typ {
	case 0:
		cx = this.rand.Intn(int(BLOCK_SIZE_X*0.4)) + int(BLOCK_SIZE_X*0.1)
	case 1:
		cx = this.rand.Intn(int(BLOCK_SIZE_X*0.4)) + int(BLOCK_SIZE_X*0.5)
	case 2:
		if this.rand.Intn(2) == 0 {
			cx = this.rand.Intn(int(BLOCK_SIZE_X*0.4)) - int(BLOCK_SIZE_X*0.2)
		} else {
			cx = this.rand.Intn(int(BLOCK_SIZE_X*0.4)) + int(BLOCK_SIZE_X*0.8)
		}
	}
	cy := this.rand.Intn(12 /* int(NEXT_BLOCK_AREA_SIZE*0.6)) + int(NEXT_BLOCK_AREA_SIZE*0.2 */)
	cy += this.nextBlockY
	w := this.rand.Intn(12 /* int(BLOCK_SIZE_X*0.33)) + int(BLOCK_SIZE_X*0.33 */)
	h := this.rand.Intn(8 /* int(NEXT_BLOCK_AREA_SIZE*0.24)) + int(NEXT_BLOCK_AREA_SIZE*0.33 */)
	cx -= w / 2
	cy -= h / 2
	var wr, hr float64
	for y := this.nextBlockY; y < this.nextBlockY+NEXT_BLOCK_AREA_SIZE; y++ {
		by := y % BLOCK_SIZE_Y
		for bx := 0; bx < BLOCK_SIZE_X; bx++ {
			if bx >= cx && bx < cx+w && y >= cy && y < cy+h {
				var o, to float64
				wr = nextFloat(this.rand, 0.2) + 0.2
				hr = nextFloat(this.rand, 0.3) + 0.4
				o = float64(bx-cx)*wr + float64(y-cy)*hr
				wr = nextFloat(this.rand, 0.2) + 0.2
				hr = nextFloat(this.rand, 0.3) + 0.4
				to = float64(cx+w-1-bx)*wr + float64(y-cy)*hr
				if to < o {
					o = to
				}
				wr = nextFloat(this.rand, 0.2) + 0.2
				hr = nextFloat(this.rand, 0.3) + 0.4
				to = float64(bx-cx)*wr + float64(cy+h-1-y)*hr
				if to < o {
					o = to
				}
				wr = nextFloat(this.rand, 0.2) + 0.2
				hr = nextFloat(this.rand, 0.3) + 0.4
				to = float64(cx+w-1-bx)*wr + float64(cy+h-1-y)*hr
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

func (this *Field) getBlockVector(p vector.Vector) int {
	return this.getBlock(p.X, p.Y)
}

func (this *Field) getBlock(x, y float64) int {
	y -= this.screenY - math.Floor(this.screenY)
	var bx, by int
	bx = int((x + BLOCK_WIDTH*SCREEN_BLOCK_SIZE_X/2) / BLOCK_WIDTH)
	by = int(this.screenY) + int((-y+BLOCK_WIDTH*SCREEN_BLOCK_SIZE_Y/2)/BLOCK_WIDTH)
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

func (this *Field) convertToScreenPos(bx, y int) vector.Vector {
	oy := this.screenY - math.Floor(this.screenY)
	by := y - int(this.screenY)
	if by <= -BLOCK_SIZE_Y {
		by += BLOCK_SIZE_Y
	}
	if by > 0 {
		by -= BLOCK_SIZE_Y
	}
	this.screenPos.X = float64(bx)*BLOCK_WIDTH - BLOCK_WIDTH*SCREEN_BLOCK_SIZE_X/2 + BLOCK_WIDTH/2
	this.screenPos.Y = float64(by)*-BLOCK_WIDTH + BLOCK_WIDTH*SCREEN_BLOCK_SIZE_Y/2 + oy - BLOCK_WIDTH/2
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
	sdl.SetColor(0, 0, 0, 1)
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
	ci := int(this.time)
	nci := ci + 1
	if nci >= TIME_COLOR_INDEX {
		nci = 0
	}
	co := int(this.time) - ci
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			this.baseColor[i][j] = baseColorTime[ci][i][j]*(1-float64(co)) + baseColorTime[nci][i][j]*float64(co)
		}
	}
	by := int(this.screenY)
	oy := this.screenY - float64(by)
	var sx float64
	sy := BLOCK_WIDTH*SCREEN_BLOCK_SIZE_Y/2 + oy
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
			p := &(this.panel[bx][by])
			sdl.SetColor(this.baseColor[p.ci][0]*p.or*0.66,
				this.baseColor[p.ci][1]*p.og*0.66,
				this.baseColor[p.ci][2]*p.ob*0.66, 1)
			gl.Vertex3d(sx+p.x, sy-p.y, p.z)
			gl.Vertex3d(sx+p.x+PANEL_WIDTH, sy-p.y, p.z)
			gl.Vertex3d(sx+p.x+PANEL_WIDTH, sy-p.y-PANEL_WIDTH, p.z)
			gl.Vertex3d(sx+p.x, sy-p.y-PANEL_WIDTH, p.z)
			sdl.SetColor(this.baseColor[p.ci][0]*0.33,
				this.baseColor[p.ci][1]*0.33,
				this.baseColor[p.ci][2]*0.33, 1)
			gl.Vertex2d(sx, sy)
			gl.Vertex2d(sx+BLOCK_WIDTH, sy)
			gl.Vertex2d(sx+BLOCK_WIDTH, sy-BLOCK_WIDTH)
			gl.Vertex2d(sx, sy-BLOCK_WIDTH)
			sx += BLOCK_WIDTH
		}
		sy -= BLOCK_WIDTH
		by++
	}
	gl.End()
}

var degBlockOfs = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func (this *Field) calcPlatformDeg(x, y int) float64 {
	d := this.rand.Intn(4)
	for i := 0; i < 4; i++ {
		if !this.checkBlock(x+degBlockOfs[d][0], y+degBlockOfs[d][1], -1, true) {
			pd := float64(d) * math.Pi / 2
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
				pd -= math.Pi / 4
			}
			if b1 && !b2 {
				pd += math.Pi / 4
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

func (this *Field) countAroundBlock(x, y int, th int /* = 0 */) int {
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

func (this *Field) checkBlock(x, y int, th int /* = 0 */, outScreen bool /* = false */) bool {
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
	return this.block[x][by] >= th
}

func (this *Field) checkInFieldVector(p vector.Vector) bool {
	return this._size.ContainsVector(p, 1)
}

func (this *Field) checkInField(x, y float64) bool {
	return this._size.Contains(x, y, 1)
}

func (this *Field) checkInOuterFieldVector(p vector.Vector) bool {
	return this._outerSize.ContainsVector(p, 1)
}

func (this *Field) checkInOuterField(x, y float64) bool {
	return this._outerSize.Contains(x, y, 1)
}

func (this *Field) checkInOuterHeightField(p vector.Vector) bool {
	return p.X >= -this._size.X && p.X <= this._size.X && p.Y >= -this._outerSize.Y && p.Y <= this._outerSize.Y
}

func (this *Field) checkInFieldExceptTop(p vector.Vector) bool {
	return p.X >= -this._size.X && p.X <= this._size.X && p.Y >= -this._size.Y
}

func (this *Field) checkInOuterFieldExceptTop(p vector.Vector) bool {
	return p.X >= -this._outerSize.X && p.X <= this._outerSize.X && p.Y >= -this._outerSize.Y && p.Y <= this._outerSize.Y*2
}

func (this *Field) size() vector.Vector {
	return this._size
}

func (this *Field) outerSize() vector.Vector {
	return this._outerSize
}

func (this *Field) lastScrollY() float64 {
	return this._lastScrollY
}
