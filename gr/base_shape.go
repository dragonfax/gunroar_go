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

var shapeRand *r.Rand
var shapeWakePos vector.Vector

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

type ShapeI interface {
	Draw()
}

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

func NewBaseShape(size, distRatio, spinyRatio float64, typ ShapeType, r, g, b float64) *BaseShape {
	this := NewBaseShapeInternal(size, distRatio, spinyRatio, typ, r, g, b)
	return &this
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
	this.DrawableShape = sdl.NewDrawableShapeInternal(&this)
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
		sdl.SetColor(this.r, this.g, this.b, 1)
	}
	gl.Begin(gl.LINE_LOOP)
	if this.typ != BRIDGE {
		this.createLoop(sz, z, false, true)
	} else {
		this.createSquareLoop(sz, z, false, 1)
	}
	gl.End()
	if this.typ != SHIP_SHADOW && this.typ != SHIP_DESTROYED &&
		this.typ != PLATFORM_DESTROYED && this.typ != TURRET_DESTROYED {
		sdl.SetColor(this.r*0.4, this.g*0.4, this.b*0.4, 1)
		gl.Begin(gl.TRIANGLE_FAN)
		this.createLoop(sz, z, true, false)
		gl.End()
	}
	switch this.typ {
	case SHIP, SHIP_ROUNDTAIL, SHIP_SHADOW, SHIP_DAMAGED, SHIP_DESTROYED:
		if this.typ != SHIP_DESTROYED {
			sdl.SetColor(this.r*0.4, this.g*0.4, this.b*0.4, 1)
		}
		for i := 0; i < 3; i++ {
			z -= height / 4
			sz -= 0.2
			gl.Begin(gl.LINE_LOOP)
			this.createLoop(sz, z, false, false)
			gl.End()
		}
	case PLATFORM, PLATFORM_DAMAGED, PLATFORM_DESTROYED:
		sdl.SetColor(this.r*0.4, this.g*0.4, this.b*0.4, 1)
		for i := 0; i < 3; i++ {
			z -= height / 3
			for _, pp := range this.pillarPos {
				gl.Begin(gl.LINE_LOOP)
				this.createPillar(pp, this.size*0.2, z)
				gl.End()
			}
		}
	case BRIDGE, TURRET, TURRET_DAMAGED:
		sdl.SetColor(this.r*0.6, this.g*0.6, this.b*0.6, 1)
		z += height
		sz -= 0.33
		gl.Begin(gl.LINE_LOOP)
		if this.typ == BRIDGE {
			this.createSquareLoop(sz, z, false, 1)
		} else {
			this.createSquareLoop(sz, z/2, false, 3)
		}
		gl.End()
		sdl.SetColor(this.r*0.25, this.g*0.25, this.b*0.25, 1)
		gl.Begin(gl.TRIANGLE_FAN)
		if this.typ == BRIDGE {
			this.createSquareLoop(sz, z, true, 1)
		} else {
			this.createSquareLoop(sz, z/2, true, 3)
		}
		gl.End()
	}
}

func (this *BaseShape) createLoop(s, z float64, backToFirst bool /* = false */, record bool /* = false */) {
	var d float64
	// var pn int
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
		gl.Vertex3d(px, py, z)
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
		gl.Vertex3d(fpx, fpy, z)
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
		gl.Vertex3d(px, py, z)
	}
}

func (this *BaseShape) createPillar(p vector.Vector, s, z float64) {
	for i := 0; i < PILLAR_POINT_NUM; i++ {
		d := math.Pi * 2 * float64(i) / PILLAR_POINT_NUM
		gl.Vertex3d(math.Sin(d)*s+p.X, math.Cos(d)*s+p.Y, z)
	}
}

func (this *BaseShape) addWake(wakes *WakePool, pos vector.Vector, deg float64, spd float64, sr float64 /* = 1 */) {
	sp := spd
	if sp > 0.1 {
		sp = 0.1
	}
	sz := this.size
	if sz > 10 {
		sz = 10
	}
	shapeWakePos.X = pos.X + math.Sin(deg+math.Pi/2+0.7)*this.size*0.5*sr
	shapeWakePos.Y = pos.Y + math.Cos(deg+math.Pi/2+0.7)*this.size*0.5*sr
	w := wakes.GetInstanceForced()
	w.set(shapeWakePos, deg+math.Pi-0.2+nextSignedFloat(rand, 0.1), sp, 40, sz*32*sr, false)
	shapeWakePos.X = pos.X + math.Sin(deg-math.Pi/2-0.7)*this.size*0.5*sr
	shapeWakePos.Y = pos.Y + math.Cos(deg-math.Pi/2-0.7)*this.size*0.5*sr
	w = wakes.GetInstanceForced()
	w.set(shapeWakePos, deg+math.Pi+0.2+nextSignedFloat(rand, 0.1), sp, 40, sz*32*sr, false)
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

func NewCollidableBaseShape(size, distRatio, spinyRatio float64, typ ShapeType, r, g, b float64) *CollidableBaseShape {
	this := NewCollidableBaseShapeInternal(size, distRatio, spinyRatio, typ, r, g, b)
	return &this
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
