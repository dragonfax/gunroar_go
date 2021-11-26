package main

import (
	"math"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"

	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

const POINT_NUM = 16
const PILLAR_POINT_NUM = 8

var rand *r.Rand
var wakePos vector.Vector

type ShapeType int

const (
	SHIP ShapeType = iota + 1
	SHIP_ROUNDTAIL
	SHIP_SHADOW
	PLATFORM
	TURRET
	BRIDGE
	SHIP_DAMAGED
	SHIP_DESTROYED
	PLATFORM_DAMAGED
	PLATFORM_DESTROYED
	TURRET_DAMAGED
	TURRET_DESTROYED
)

/**
 * Shape of a ship/platform/turret/bridge.
 */
type BaseShape struct {
	sdl.DrawableShape

	size, distRatio, spinyRatio float64
	typ                         ShapeType
	r, g, b                     float64
	pillarPos                   []vector.Vector
	_pointPos                   []vector.Vector
	_pointDeg                   []float64
}

func BaseShapeInit() {
	rand = r.New(r.NewSource(time.Now().Unix()))
}

func SetRandSeed(seed int64) {
	rand = r.New(r.NewSource(seed))
}

func NewBaseShapeInternal(size, distRatio, spinyRatio float64, typ ShapeType, r, g, b float64) BaseShape {
	this := BaseShape{
		size:       size,
		distRatio:  distRatio,
		spinyRatio: spinyRatio,
		typ:        typ,
		r:          r,
		g:          g,
		b:          g,
	}
	this.DrawableShape = sdl.NewDrawableShapeInternal(this)
	return this
}

func (this *BaseShape) CreateDisplayList() {
	height := this.size * 0.5
	var z float64
	sz := 1.0
	if this.typ == BRIDGE {
		z += height
	}
	if this.typ != SHIP_DESTROYED {
		Screen.SetColor(r, g, b)
	}
	gl.Begin(gl.LINE_LOOP)
	if this.typ != BRIDGE {
		this.createLoop(sz, z, false, true)
	} else {
		this.createSquareLoop(sz, z, false, true)
	}
	gl.End()
	if this.typ != SHIP_SHADOW && this.typ != SHIP_DESTROYED &&
		this.typ != PLATFORM_DESTROYED && this.typ != TURRET_DESTROYED {
		Screen.SetColor(r*0.4, g*0.4, b*0.4)
		gl.Begin(gl.TRIANGLE_FAN)
		this.createLoop(sz, z, true)
		gl.End()
	}
	switch this.typ {
	case SHIP, SHIP_ROUNDTAIL, SHIP_SHADOW, SHIP_DAMAGED, SHIP_DESTROYED:
		if this.typ != SHIP_DESTROYED {
			Screen.SetColor(r*0.4, g*0.4, b*0.4)
		}
		for i := 0; i < 3; i++ {
			z -= height / 4
			sz -= 0.2
			gl.Begin(gl.LINE_LOOP)
			this.createLoop(sz, z)
			gl.End()
		}
	case PLATFORM, PLATFORM_DAMAGED, PLATFORM_DESTROYED:
		Screen.SetColor(r*0.4, g*0.4, b*0.4)
		for i := 0; i < 3; i++ {
			z -= height / 3
			for _, pp := range this.pillarPos {
				gl.Begin(gl.LINE_LOOP)
				this.createPillar(pp, this.size*0.2, z)
				gl.End()
			}
		}
	case BRIDGE, TURRET, TURRET_DAMAGED:
		Screen.SetColor(r*0.6, g*0.6, b*0.6)
		z += height
		sz -= 0.33
		gl.Begin(gl.LINE_LOOP)
		if this.typ == BRIDGE {
			this.createSquareLoop(sz, z)
		} else {
			this.createSquareLoop(sz, z/2, false, 3)
		}
		gl.End()
		Screen.SetColor(r*0.25, g*0.25, b*0.25)
		gl.Begin(gl.TRIANGLE_FAN)
		if this.typ == BRIDGE {
			this.createSquareLoop(sz, z, true)
		} else {
			this.createSquareLoop(sz, z/2, true, 3)
		}
		gl.End()
	}
}

func (this *BaseShape) createLoop(s, z float64, backToFirst bool /* = false */, record bool /* = false */) {
	var d float64
	var pn int
	firstPoint := true
	var fpx, fpy float64
	for i := 0; i < POINT_NUM; i++ {
		if this.typ != SHIP && this.typ != SHIP_DESTROYED && this.typ != SHIP_DAMAGED &&
			i > POINT_NUM*2/5 && i <= POINT_NUM*3/5 {
			continue
		}
		if (this.typ == TURRET || this.typ == TURRET_DAMAGED || this.typ == TURRET_DESTROYED) &&
			(i <= POINT_NUM/5 || i > POINT_NUM*4/5) {
			continue
		}
		d = math.Pi * 2 * float64(i) / POINT_NUM
		cx := math.Sin(d) * this.size * s * (1 - this.distRatio)
		cy := math.Cos(d) * this.size * s
		var sx, sy float64
		if i == POINT_NUM/4 || i == POINT_NUM/4*3 {
			sy = 0
		} else {
			sy = 1 / (1 + math.Abs(math.Tan(d)))
		}
		sx = 1 - sy
		if i >= POINT_NUM/2 {
			sx *= -1
		}
		if i >= POINT_NUM/4 && i <= POINT_NUM/4*3 {
			sy *= -1
		}
		sx *= this.size * s * (1 - this.distRatio)
		sy *= this.size * s
		px := cx*(1-this.spinyRatio) + sx*this.spinyRatio
		py := cy*(1-this.spinyRatio) + sy*this.spinyRatio
		gl.Vertex3f(px, py, z)
		if backToFirst && firstPoint {
			fpx = px
			fpy = py
			firstPoint = false
		}
		if record {
			if i == POINT_NUM/8 || i == POINT_NUM/8*3 ||
				i == POINT_NUM/8*5 || i == POINT_NUM/8*7 {
				this.pillarPos = append(this.pillarPos, vector.Vector{px * 0.8, py * 0.8})
			}
			this._pointPos = append(this._pointPos, vector.Vector{px, py})
			this._pointDeg = append(this._pointDeg, d)
		}
	}
	if backToFirst {
		gl.Vertex3f(fpx, fpy, z)
	}
}

func (this *BaseShape) createSquareLoop(s, z float64, backToFirst bool /* = false */, yRatio float64 /* = 1 */) {
	var d float64
	var pn int
	if backToFirst {
		pn = 4
	} else {
		pn = 3
	}
	for i := 0; i <= pn; i++ {
		d = math.Pi*2*float64(i)/4 + math.Pi/4
		px := math.Sin(d) * this.size * s
		py := math.Cos(d) * this.size * s
		if py > 0 {
			py *= yRatio
		}
		gl.Vertex3f(px, py, z)
	}
}

func (this *BaseShape) createPillar(p vector.Vector, s, z float64) {
	var d float64
	for i := 0; i < PILLAR_POINT_NUM; i++ {
		d := math.Pi * 2 * float64(i) / PILLAR_POINT_NUM
		gl.Vertex3f(math.Sin(d)*s+p.X, math.Cos(d)*s+p.Y, z)
	}
}

func (this *BaseShape) addWake(wakes WakePool, pos vector.Vector, deg float64, spd float64, sr float64 /* = 1 */) {
	sp := spd
	if sp > 0.1 {
		sp = 0.1
	}
	sz := this.size
	if sz > 10 {
		sz = 10
	}
	wakePos.X = pos.X + math.Sin(deg+math.Pi/2+0.7)*this.size*0.5*sr
	wakePos.Y = pos.Y + math.Cos(deg+math.Pi/2+0.7)*this.size*0.5*sr
	w := wakes.getInstanceForced()
	w.Set(wakePos, deg+math.Pi-0.2+rand.nextSignedFloat(0.1), sp, 40, sz*32*sr)
	wakePos.X = pos.X + math.Sin(deg-math.Pi/2-0.7)*this.size*0.5*sr
	wakePos.Y = pos.Y + math.Cos(deg-math.Pi/2-0.7)*this.size*0.5*sr
	w = wakes.getInstanceForced()
	w.Set(wakePos, deg+math.Pi+0.2+rand.nextSignedFloat(0.1), sp, 40, sz*32*sr)
}

func (this *BaseShape) pointPos() []vector.Vector {
	return this._pointPos
}

func (this *BaseShape) pointDeg() []float64 {
	return this._pointDeg
}

func (this *BaseShape) checkShipCollision(x, y, deg float64, sr float64 /* = 1 */) bool {
	cs := this.size * (1 - this.distRatio) * 1.1 * sr
	if this.dist(x, y, 0, 0) < cs {
		return true
	}
	var ofs float64
	for {
		ofs += cs
		cs *= this.distRatio
		if cs < 0.2 {
			return false
		}
		if this.dist(x, y, math.Sin(deg)*ofs, math.Cos(deg)*ofs) < cs ||
			this.dist(x, y, -math.Sin(deg)*ofs, -math.Cos(deg)*ofs) < cs {
			return true
		}
	}
}

func (this *BaseShape) dist(x, y, px, py float64) float64 {
	ax := math.Abs(x - px)
	ay := math.Abs(y - py)
	if ax > ay {
		return ax + ay/2
	} else {
		return ay + ax/2
	}
}

var _ sdl.Collidable = &CollidableBaseShape{}

type CollidableBaseShape struct {
	BaseShape
	sdl.CollidableImpl

	_collision vector.Vector
}

func NewCollidableBaseShapeInternal(size, distRatio, spinyRatio float64, typ ShapeType, r, g, b float64) CollidableBaseShape {
	this := CollidableBaseShape{
		BaseShape:  NewBaseShapeInternal(size, distRatio, spinyRatio, typ, r, g, b),
		_collision: vector.Vector{size / 2, size / 2},
	}
	this.CollidableImpl = sdl.NewCollidableInternal(&this)
	return this
}

func (this *CollidableBaseShape) Collision() *vector.Vector {
	v := this._collision
	return &v
}