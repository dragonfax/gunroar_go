package main

import (
	"math"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/actor"
	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Sparks.
 */

var sparkRand = r.New(r.NewSource(time.Now().Unix()))

var _ sdl.LuminousActor = &Spark{}

type Spark struct {
	actor.ExistsImpl

	pos, ppos, vel vector.Vector
	r, g, b        float64
	cnt            int
}

func setSparkRandSeed(seed int64) {
	sparkRand = r.New(r.NewSource(seed))
}

func NewSpark() *Spark {
	return &Spark{}
}

func (*Spark) Init(args []interface{}) {
}

func (this *Spark) set(p vector.Vector, vx, vy, r, g, b float64, c int) {
	this.ppos.X = p.X
	this.pos.X = p.X
	this.ppos.Y = p.Y
	this.pos.Y = p.Y
	this.vel.X = vx
	this.vel.Y = vy
	this.r = r
	this.g = g
	this.b = b
	this.cnt = c
	this.SetExists(true)
}

func (this *Spark) Move() {
	this.cnt--
	if this.cnt <= 0 || this.vel.Dist(0, 0) < 0.005 {
		this.SetExists(false)
		return
	}
	this.ppos.X = this.pos.X
	this.ppos.Y = this.pos.Y
	this.pos.OpAddAssign(this.vel)
	this.vel.OpMulAssign(0.96)
}

func (this *Spark) Draw() {
	ox := this.vel.X
	oy := this.vel.Y
	sdl.SetColor(this.r, this.g, this.b, 1)
	ox *= 2
	oy *= 2
	gl.Vertex3d(this.pos.X-ox, this.pos.Y-oy, 0)
	ox *= 0.5
	oy *= 0.5
	sdl.SetColor(this.r*0.5, this.g*0.5, this.b*0.5, 0)
	gl.Vertex3d(this.pos.X-oy, this.pos.Y+ox, 0)
	gl.Vertex3d(this.pos.X+oy, this.pos.Y-ox, 0)
}

func (this *Spark) DrawLuminous() {
	ox := this.vel.X
	oy := this.vel.Y
	sdl.SetColor(this.r, this.g, this.b, 1)
	ox *= 2
	oy *= 2
	gl.Vertex3d(this.pos.X-ox, this.pos.Y-oy, 0)
	ox *= 0.5
	oy *= 0.5
	sdl.SetColor(this.r*0.5, this.g*0.5, this.b*0.5, 0)
	gl.Vertex3d(this.pos.X-oy, this.pos.Y+ox, 0)
	gl.Vertex3d(this.pos.X+oy, this.pos.Y-ox, 0)
}

type SparkPool struct {
	actor.ActorPool
}

func NewSparkPool(n int, args []interface{}) *SparkPool {
	f := func() actor.Actor { return NewSpark() }
	this := &SparkPool{ActorPool: actor.NewActorPool(f, n, args)}
	return this
}

/**
 * Smokes.
 */

type SmokeType int

const (
	FIRE SmokeType = iota
	EXPLOSION
	SAND
	SPARK
	WAKE
	SMOKE
	LANCE_SPARK
)

var smokeRand = r.New(r.NewSource(time.Now().Unix()))
var windVel = vector.Vector3{0.04, 0.04, 0.02}
var smokeWakePos vector.Vector

var _ sdl.LuminousActor = &Smoke{}

type Smoke struct {
	actor.ExistsImpl

	field            *Field
	wakes            *WakePool
	pos, vel         vector.Vector3
	typ              SmokeType
	cnt, startCnt    int
	size, r, g, b, a float64
}

func setSmokeRandSeed(seed int64) {
	smokeRand = r.New(r.NewSource(seed))
}

func NewSmoke() *Smoke {
	this := &Smoke{
		startCnt: 1,
		size:     1,
	}
	return this
}

func (this *Smoke) Init(args []interface{}) {
	this.field = args[0].(*Field)
	this.wakes = args[1].(*WakePool)
}

func (this *Smoke) setVector(p vector.Vector, mx, my, mz float64, t SmokeType, c int /* = 60 */, sz float64 /* = 2 */) {
	this.set(p.X, p.Y, mx, my, mz, t, c, sz)
}

func (this *Smoke) setVector3(p vector.Vector3, mx, my, mz float64, t SmokeType, c int /* = 60 */, sz float64 /* = 2 */) {
	this.set(p.X, p.Y, mx, my, mz, t, c, sz)
	this.pos.Z = p.Z
}

func (this *Smoke) set(x, y, mx, my, mz float64, t SmokeType, c int /* = 60 */, sz float64 /* = 2 */) {
	if !this.field.checkInOuterField(x, y) {
		return
	}
	this.pos.X = x
	this.pos.Y = y
	this.pos.Z = 0
	this.vel.X = mx
	this.vel.Y = my
	this.vel.Z = mz
	this.typ = t
	this.startCnt = c
	this.cnt = c
	this.size = sz
	switch this.typ {
	case FIRE:
		this.r = nextFloat(rand, 0.1) + 0.9
		this.g = nextFloat(rand, 0.2) + 0.2
		this.b = 0
		this.a = 1
	case EXPLOSION:
		this.r = nextFloat(rand, 0.3) + 0.7
		this.g = nextFloat(rand, 0.3) + 0.3
		this.b = 0
		this.a = 1
	case SAND:
		this.r = 0.8
		this.g = 0.8
		this.b = 0.6
		this.a = 0.6
	case SPARK:
		this.r = nextFloat(rand, 0.3) + 0.7
		this.g = nextFloat(rand, 0.5) + 0.5
		this.b = 0
		this.a = 1
	case WAKE:
		this.r = 0.6
		this.g = 0.6
		this.b = 0.8
		this.a = 0.6
	case SMOKE:
		this.r = nextFloat(rand, 0.1) + 0.1
		this.g = nextFloat(rand, 0.1) + 0.1
		this.b = 0.1
		this.a = 0.5
	case LANCE_SPARK:
		this.r = 0.4
		this.g = nextFloat(rand, 0.2) + 0.7
		this.b = nextFloat(rand, 0.2) + 0.7
		this.a = 1
	}
	this.SetExists(true)
}

func (this *Smoke) Move() {
	this.cnt--
	if this.cnt <= 0 || !this.field.checkInOuterField(this.pos.X, this.pos.Y) {
		this.SetExists(false)
		return
	}
	if this.typ != WAKE {
		this.vel.X += (windVel.X - this.vel.X) * 0.01
		this.vel.Y += (windVel.Y - this.vel.Y) * 0.01
		this.vel.Z += (windVel.Z - this.vel.Z) * 0.01
	}
	this.pos.OpAddAssign(this.vel)
	this.pos.Y -= this.field.lastScrollY()
	switch this.typ {
	case FIRE, EXPLOSION, SMOKE:
		if this.cnt < this.startCnt/2 {
			this.r *= 0.95
			this.g *= 0.95
			this.b *= 0.95
		} else {
			this.a *= 0.97
		}
		this.size *= 1.01
	case SAND:
		this.r *= 0.98
		this.g *= 0.98
		this.b *= 0.98
		this.a *= 0.98
	case SPARK:
		this.r *= 0.92
		this.g *= 0.92
		this.a *= 0.95
		this.vel.OpMulAssign(0.9)
	case WAKE:
		this.a *= 0.98
		this.size *= 1.005
	case LANCE_SPARK:
		this.a *= 0.95
		this.size *= 0.97
	}
	if this.size > 5 {
		this.size = 5
	}
	if this.typ == EXPLOSION && this.pos.Z < 0.01 {
		bl := this.field.getBlock(this.pos.X, this.pos.Y)
		if bl >= 1 {
			this.vel.OpMulAssign(0.8)
		}
		if this.cnt%3 == 0 && bl < -1 {
			sp := math.Sqrt(this.vel.X*this.vel.X + this.vel.Y*this.vel.Y)
			if sp > 0.3 {
				d := math.Atan2(this.vel.X, this.vel.Y)
				smokeWakePos.X = this.pos.X + math.Sin(d+math.Pi/2)*this.size*0.25
				smokeWakePos.Y = this.pos.Y + math.Cos(d+math.Pi/2)*this.size*0.25
				w := this.wakes.GetInstanceForced()
				w.set(smokeWakePos, d+math.Pi-0.2+nextSignedFloat(rand, 0.1), sp*0.33,
					20+rand.Intn(12), this.size*(7.0+nextFloat(rand, 3)), false)
				smokeWakePos.X = this.pos.X + math.Sin(d-math.Pi/2)*this.size*0.25
				smokeWakePos.Y = this.pos.Y + math.Cos(d-math.Pi/2)*this.size*0.25
				w = this.wakes.GetInstanceForced()
				w.set(smokeWakePos, d+math.Pi+0.2+nextSignedFloat(rand, 0.1), sp*0.33,
					20+rand.Intn(12), this.size*(7.0+nextFloat(rand, 3)), false)
			}
		}
	}
}

func (this *Smoke) Draw() {
	quadSize := this.size / 2
	sdl.SetColor(this.r, this.g, this.b, this.a)
	gl.Vertex3d(this.pos.X-quadSize, this.pos.Y-quadSize, this.pos.Z)
	gl.Vertex3d(this.pos.X+quadSize, this.pos.Y-quadSize, this.pos.Z)
	gl.Vertex3d(this.pos.X+quadSize, this.pos.Y+quadSize, this.pos.Z)
	gl.Vertex3d(this.pos.X-quadSize, this.pos.Y+quadSize, this.pos.Z)
}

func (this *Smoke) DrawLuminous() {
	if this.r+this.g > 0.8 && this.b < 0.5 {
		quadSize := this.size / 2
		sdl.SetColor(this.r, this.g, this.b, this.a)
		gl.Vertex3d(this.pos.X-quadSize, this.pos.Y-quadSize, this.pos.Z)
		gl.Vertex3d(this.pos.X+quadSize, this.pos.Y-quadSize, this.pos.Z)
		gl.Vertex3d(this.pos.X+quadSize, this.pos.Y+quadSize, this.pos.Z)
		gl.Vertex3d(this.pos.X-quadSize, this.pos.Y+quadSize, this.pos.Z)
	}
}

type SmokePool struct {
	sdl.LuminousActorPool
}

func NewSmokePool(n int, args []interface{}) *SmokePool {
	f := func() actor.Actor { return NewSmoke() }
	this := &SmokePool{LuminousActorPool: sdl.NewLuminousActorPool(f, n, args)}
	return this
}

func (this *SmokePool) GetInstanceForced() *Smoke {
	return this.LuminousActorPool.GetInstanceForced().(*Smoke)
}

func (this *SmokePool) GetInstance() *Smoke {
	return this.LuminousActorPool.GetInstance().(*Smoke)
}

/**
 * Fragments of destroyed enemies.
 */
var _ actor.Actor = &Fragment{}

var fragmentDisplayList *sdl.DisplayList
var fragmentRand = r.New(r.NewSource(time.Now().Unix()))

type Fragment struct {
	actor.ExistsImpl

	field         *Field
	smokes        *SmokePool
	pos, vel      vector.Vector3
	size, d2, md2 float64
}

func fragmentInit() {
	fragmentDisplayList = sdl.NewDisplayList(1)
	fragmentDisplayList.BeginNewList()
	sdl.SetColor(0.7, 0.5, 0.5, 0.5)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex2d(-0.5, -0.25)
	gl.Vertex2d(0.5, -0.25)
	gl.Vertex2d(0.5, 0.25)
	gl.Vertex2d(-0.5, 0.25)
	gl.End()
	sdl.SetColor(0.7, 0.5, 0.5, 0.9)
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex2d(-0.5, -0.25)
	gl.Vertex2d(0.5, -0.25)
	gl.Vertex2d(0.5, 0.25)
	gl.Vertex2d(-0.5, 0.25)
	gl.End()
	fragmentDisplayList.EndNewList()
}

func setFragmentRandSeed(seed int64) {
	fragmentRand = r.New(r.NewSource(seed))
}

func NewFragment() *Fragment {
	this := &Fragment{}
	this.size = 1
	return this
}

func (this *Fragment) Init(args []interface{}) {
	this.field = args[0].(*Field)
	this.smokes = args[1].(*SmokePool)
}

func (this *Fragment) set(p vector.Vector, mx, my, mz float64, sz float64 /* = 1 */) {
	if !this.field.checkInOuterField(p.X, p.Y) {
		return
	}
	this.pos.X = p.X
	this.pos.Y = p.Y
	this.pos.Z = 0
	this.vel.X = mx
	this.vel.Y = my
	this.vel.Z = mz
	this.size = sz
	if this.size > 5 {
		this.size = 5
	}
	this.d2 = nextFloat(rand, 360)
	this.md2 = nextSignedFloat(rand, 20)
	this.SetExists(true)
}

func (this *Fragment) Move() {
	if !this.field.checkInOuterField(this.pos.X, this.pos.Y) {
		this.SetExists(false)
		return
	}
	this.vel.X *= 0.96
	this.vel.Y *= 0.96
	this.vel.Z += (-0.04 - this.vel.Z) * 0.01
	this.pos.OpAddAssign(this.vel)
	if this.pos.Z < 0 {
		s := this.smokes.GetInstanceForced()
		if this.field.getBlock(this.pos.X, this.pos.Y) < 0 {
			s.set(this.pos.X, this.pos.Y, 0, 0, 0, WAKE, 60, this.size*0.66)
		} else {
			s.set(this.pos.X, this.pos.Y, 0, 0, 0, SAND, 60, this.size*0.75)
		}
		this.SetExists(false)
		return
	}
	this.pos.Y -= this.field.lastScrollY()
	this.d2 += this.md2
}

func (this *Fragment) Draw() {
	gl.PushMatrix()
	sdl.GlTranslate3(this.pos)
	gl.Rotated(this.d2, 1, 0, 0)
	gl.Scaled(this.size, this.size, 1)
	fragmentDisplayList.Call(0)
	gl.PopMatrix()
}

type FragmentPool struct {
	actor.ActorPool
}

func NewFragmentPool(n int, args []interface{}) *FragmentPool {
	f := func() actor.Actor { return NewFragment() }
	return &FragmentPool{actor.NewActorPool(f, n, args)}
}

/**
 * Luminous fragments.
 */

var sparkFragmentDisplayList *sdl.DisplayList
var sparkFragmentRand = r.New(r.NewSource(time.Now().Unix()))

var _ sdl.LuminousActor = &SparkFragment{}

type SparkFragment struct {
	actor.ExistsImpl

	field         *Field
	smokes        *SmokePool
	pos, vel      vector.Vector3
	size, d2, md2 float64
	cnt           int
	hasSmoke      bool
}

func sparkFragmentInit() {
	sparkFragmentDisplayList = sdl.NewDisplayList(1)
	sparkFragmentDisplayList.BeginNewList()
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex2d(-0.25, -0.25)
	gl.Vertex2d(0.25, -0.25)
	gl.Vertex2d(0.25, 0.25)
	gl.Vertex2d(-0.25, 0.25)
	gl.End()
	sparkFragmentDisplayList.EndNewList()
}

func setSparkFragmentRandSeed(seed int64) {
	rand = r.New(r.NewSource(seed))
}

func NewSparkFragment() *SparkFragment {
	this := &SparkFragment{}
	this.size = 1
	return this
}

func (this *SparkFragment) Init(args []interface{}) {
	this.field = args[0].(*Field)
	this.smokes = args[1].(*SmokePool)
}

func (this *SparkFragment) set(p vector.Vector, mx, my, mz float64, sz float64 /* = 1 */) {
	if !this.field.checkInOuterField(p.X, p.Y) {
		return
	}
	this.pos.X = p.X
	this.pos.Y = p.Y
	this.pos.Z = 0
	this.vel.X = mx
	this.vel.Y = my
	this.vel.Z = mz
	this.size = sz
	if this.size > 5 {
		this.size = 5
	}
	this.d2 = nextFloat(rand, 360)
	this.md2 = nextSignedFloat(rand, 15)
	this.hasSmoke = rand.Intn(4) == 0
	this.cnt = 0
	this.SetExists(true)
}

func (this *SparkFragment) Move() {
	if !this.field.checkInOuterField(this.pos.X, this.pos.Y) {
		this.SetExists(false)
		return
	}
	this.vel.X *= 0.99
	this.vel.Y *= 0.99
	this.vel.Z += (-0.08 - this.vel.Z) * 0.01
	this.pos.OpAddAssign(this.vel)
	if this.pos.Z < 0 {
		s := this.smokes.GetInstanceForced()
		if this.field.getBlock(this.pos.X, this.pos.Y) < 0 {
			s.set(this.pos.X, this.pos.Y, 0, 0, 0, WAKE, 60, this.size*0.66)
		} else {
			s.set(this.pos.X, this.pos.Y, 0, 0, 0, SAND, 60, this.size*0.75)
		}
		this.SetExists(false)
		return
	}
	this.pos.Y -= this.field.lastScrollY()
	this.d2 += this.md2
	this.cnt++
	if this.hasSmoke && this.cnt%5 == 0 {
		s := this.smokes.GetInstance()
		if s != nil {
			s.setVector3(this.pos, 0, 0, 0, SMOKE, 90+rand.Intn(60), this.size*0.5)
		}
	}
}

func (this *SparkFragment) Draw() {
	gl.PushMatrix()
	sdl.SetColor(1, nextFloat(rand, 1), 0, 0.8)
	sdl.GlTranslate3(this.pos)
	gl.Rotated(this.d2, 1, 0, 0)
	gl.Scaled(this.size, this.size, 1)
	sparkFragmentDisplayList.Call(0)
	gl.PopMatrix()
}

func (this *SparkFragment) DrawLuminous() {
	gl.PushMatrix()
	sdl.SetColor(1, nextFloat(rand, 1), 0, 0.8)
	sdl.GlTranslate3(this.pos)
	gl.Rotated(this.d2, 1, 0, 0)
	gl.Scaled(this.size, this.size, 1)
	sparkFragmentDisplayList.Call(0)
	gl.PopMatrix()
}

type SparkFragmentPool struct {
	actor.ActorPool
}

func NewSparkFragmentPool(n int, args []interface{}) *SparkFragmentPool {
	f := func() actor.Actor { return NewSparkFragment() }
	return &SparkFragmentPool{actor.NewActorPool(f, n, args)}
}

/**
 * Wakes of ships and smokes.
 */
var _ actor.Actor = &Wake{}

type Wake struct {
	actor.ExistsImpl

	field            *Field
	pos, vel         vector.Vector
	deg, speed, size float64
	cnt              int
	revShape         bool
}

func NewWake() *Wake {
	this := &Wake{}
	this.size = 1
	return this
}

func (this *Wake) Init(args []interface{}) {
	this.field = args[0].(*Field)
}

func (this *Wake) set(p vector.Vector, deg, speed float64, c int /* = 60 */, sz float64 /* = 1 */, rs bool /* = false */) {
	if !this.field.checkInOuterField(p.X, p.Y) {
		return
	}
	this.pos.X = p.X
	this.pos.Y = p.Y
	this.deg = deg
	this.speed = speed
	this.vel.X = math.Sin(deg) * speed
	this.vel.Y = math.Cos(deg) * speed
	this.cnt = c
	this.size = sz
	this.revShape = rs
	this.SetExists(true)
}

func (this *Wake) Move() {
	this.cnt--
	if this.cnt <= 0 || this.vel.Dist(0, 0) < 0.005 || !this.field.checkInOuterField(this.pos.X, this.pos.Y) {
		this.SetExists(false)
		return
	}
	this.pos.OpAddAssign(this.vel)
	this.pos.Y -= this.field.lastScrollY()
	this.vel.OpMulAssign(0.96)
	this.size *= 1.02
}

func (this *Wake) Draw() {
	ox := this.vel.X
	oy := this.vel.Y
	sdl.SetColor(0.33, 0.33, 1, 1)
	ox *= this.size
	oy *= this.size
	if this.revShape {
		gl.Vertex3d(this.pos.X+ox, this.pos.Y+oy, 0)
	} else {
		gl.Vertex3d(this.pos.X-ox, this.pos.Y-oy, 0)
	}
	ox *= 0.2
	oy *= 0.2
	sdl.SetColor(0.2, 0.2, 0.6, 0.5)
	gl.Vertex3d(this.pos.X-oy, this.pos.Y+ox, 0)
	gl.Vertex3d(this.pos.X+oy, this.pos.Y-ox, 0)
}

type WakePool struct {
	actor.ActorPool
}

func NewWakePool(n int, args []interface{}) *WakePool {
	f := func() actor.Actor { return NewWake() }
	return &WakePool{actor.NewActorPool(f, n, args)}
}

func (this *WakePool) GetInstanceForced() *Wake {
	return this.ActorPool.GetInstanceForced().(*Wake)
}
