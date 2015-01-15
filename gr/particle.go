/*
 * $Id: particle.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gtr

/**
 * Sparks.
 */
type Spark struct {
	*LuminousActor

	pos, ppos Vector
	vel       Vector
	r, g, b   float32
	cnt       int
}

func NewSpark(p Vector, vx float32, vy float32, r float32, g float32, b float32, c int) *Spark {
	this := &Spark{NewLuminousActor()}
	this.pos.x = p.x
	this.ppos.x = p.x
	this.pos.y = p.y
	this.ppos.y = p.y
	this.vel.x = vx
	this.vel.y = vy
	this.r = r
	this.g = g
	this.b = b
	this.cnt = c
	return this
}

func (this *Spark) move() {
	this.cnt--
	if this.cnt <= 0 || this.vel.dist() < 0.005 {
		this.exists = false
		return
	}
	this.ppos.x = this.pos.x
	this.ppos.y = this.pos.y
	this.pos += this.vel
	this.vel *= 0.96
}

func (this *Spark) draw() {
	ox := this.vel.x
	oy := this.vel.y
	setScreenColor(r, g, b, 1)
	ox *= 2
	oy *= 2
	gl.Vertex3(this.pos.x-ox, this.pos.y-oy, 0)
	ox *= 0.5
	oy *= 0.5
	setScreenColor(r*0.5, g*0.5, b*0.5, 0)
	gl.Vertex3(this.pos.x-oy, this.pos.y+ox, 0)
	gl.Vertex3(this.pos.x+oy, this.pos.y-ox, 0)
}

func (this *Spark) drawLuminous() {
	ox := this.vel.x
	oy := this.vel.y
	setScreenColor(r, g, b, 1)
	ox *= 2
	oy *= 2
	gl.Vertex3(this.pos.x-ox, this.pos.y-oy, 0)
	ox *= 0.5
	oy *= 0.5
	setScreenColor(r*0.5, g*0.5, b*0.5, 0)
	gl.Vertex3(this.pos.x-oy, this.pos.y+ox, 0)
	gl.Vertex3(this.pos.x+oy, this.pos.y-ox, 0)
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

var windVel Vector3 = Vector3{0.04, 0.04, 0.02}
var wakePos Vector = Vector{}

type Smoke struct {
	*LuminousActor

	field            Field
	pos, vel         Vector3
	smokeType        SmokeType
	cnt, startCnt    int
	size, r, g, b, a float32
}

func NewSmoke(field Field) *Smoke {
	this = &Smoke{NewLuminousActor()}
	this.startCnt = 1
	this.size = 1
	this.field = field
	return this
}

func (this *Smoke) set(p Vector, mx float32, my float32, mz float32, t int, c int /*= 60*/, sz float32 /* = 2 */) {
	set(p.x, p.y, mx, my, mz, t, c, sz)
}

func (this *Smoke) set(p Vector3, mx float32, my float32, mz float32, t int, c int /*= 60*/, sz float32 /* = 2*/) {
	set(p.x, p.y, mx, my, mz, t, c, sz)
	this.pos.z = p.z
}

func (this *Smoke) set(x float32, y float32, mx float32, my float32, mz float32, t int, c int /* = 60 */, sz float32 /* = 2 */) {
	if !this.field.checkInOuterField(x, y) {
		return
	}
	this.pos.x = x
	this.pos.y = y
	this.pos.z = 0
	this.vel.x = mx
	this.vel.y = my
	this.vel.z = mz
	this.smokeType = t
	this.startCnt = c
	this.cnt = c
	this.size = sz
	switch this.smokeType {
	case SmokeType.FIRE:
		this.r = rand.nextFloat(0.1) + 0.9
		this.g = rand.nextFloat(0.2) + 0.2
		this.b = 0
		this.a = 1
		break
	case SmokeType.EXPLOSION:
		this.r = rand.nextFloat(0.3) + 0.7
		this.g = rand.nextFloat(0.3) + 0.3
		this.b = 0
		this.a = 1
		break
	case SmokeType.SAND:
		this.r = 0.8
		this.g = 0.8
		this.b = 0.6
		this.a = 0.6
		break
	case SmokeType.SPARK:
		this.r = rand.nextFloat(0.3) + 0.7
		this.g = rand.nextFloat(0.5) + 0.5
		this.b = 0
		this.a = 1
		break
	case SmokeType.WAKE:
		this.r = 0.6
		this.g = 0.6
		this.b = 0.8
		this.a = 0.6
		break
	case SmokeType.SMOKE:
		this.r = rand.nextFloat(0.1) + 0.1
		this.g = rand.nextFloat(0.1) + 0.1
		this.b = 0.1
		this.a = 0.5
		break
	case SmokeType.LANCE_SPARK:
		this.r = 0.4
		this.g = rand.nextFloat(0.2) + 0.7
		this.b = rand.nextFloat(0.2) + 0.7
		this.a = 1
		break
	}
	this.exists = true
}

func (this *Smoke) move() {
	this.cnt--
	if this.cnt <= 0 || !this.field.checkInOuterField(this.pos.x, this.pos.y) {
		this.exists = false
		return
	}
	if this.smokeType != SmokeType.WAKE {
		this.vel.x += (this.windVel.x - this.vel.x) * 0.01
		this.vel.y += (this.windVel.y - this.vel.y) * 0.01
		this.vel.z += (this.windVel.z - this.vel.z) * 0.01
	}
	this.pos += this.vel
	this.pos.y -= this.field.lastScrollY
	switch this.smokeType {
	case SmokeType.FIRE:
	case SmokeType.EXPLOSION:
	case SmokeType.SMOKE:
		if this.cnt < this.startCnt/2 {
			this.r *= 0.95
			this.g *= 0.95
			this.b *= 0.95
		} else {
			this.a *= 0.97
		}
		this.size *= 1.01
		break
	case SmokeType.SAND:
		this.r *= 0.98
		this.g *= 0.98
		this.b *= 0.98
		this.a *= 0.98
		break
	case SmokeType.SPARK:
		this.r *= 0.92
		this.g *= 0.92
		this.a *= 0.95
		this.vel *= 0.9
		break
	case SmokeType.WAKE:
		this.a *= 0.98
		this.size *= 1.005
		break
	case SmokeType.LANCE_SPARK:
		this.a *= 0.95
		this.size *= 0.97
		break
	}
	if this.size > 5 {
		this.size = 5
	}
	if this.smokeType == SmokeType.EXPLOSION && this.pos.z < 0.01 {
		bl := this.field.getBlock(this.pos.x, this.pos.y)
		if bl >= 1 {
			this.vel *= 0.8
		}
		if this.cnt%3 == 0 && bl < -1 {
			sp := sqrt(this.vel.x*this.vel.x + this.vel.y*this.vel.y)
			if sp > 0.3 {
				d := atan2(this.vel.x, this.vel.y)
				this.wakePos.x = this.pos.x + sin(d+PI/2)*this.size*0.25
				this.wakePos.y = this.pos.y + cos(d+PI/2)*this.size*0.25
				NewWake(this.wakePos, d+PI-0.2+rand.nextSignedFloat(0.1), sp*0.33,
					20+rand.nextInt(12), this.size*(7.0+rand.nextFloat(3)))
				this.wakePos.x = this.pos.x + sin(d-PI/2)*this.size*0.25
				this.wakePos.y = this.pos.y + cos(d-PI/2)*this.size*0.25
				NewWake(this.wakePos, d+PI+0.2+rand.nextSignedFloat(0.1), sp*0.33,
					20+rand.nextInt(12), this.size*(7.0+rand.nextFloat(3)))
			}
		}
	}
}

func (this *Smoke) draw() {
	quadSize := this.size / 2
	setScreenColor(this.r, this.g, this.b, this.a)
	gl.Vertex3(this.pos.x-quadSize, this.pos.y-quadSize, this.pos.z)
	gl.Vertex3(this.pos.x+quadSize, this.pos.y-quadSize, this.pos.z)
	gl.Vertex3(this.pos.x+quadSize, this.pos.y+quadSize, this.pos.z)
	gl.Vertex3(this.pos.x-quadSize, this.pos.y+quadSize, this.pos.z)
}

func (this *Smoke) drawLuminous() {
	if this.r+this.g > 0.8 && this.b < 0.5 {
		quadSize := this.size / 2
		setScreenColor(this.r, this.g, this.b, this.a)
		gl.Vertex3(this.pos.x-quadSize, this.pos.y-quadSize, this.pos.z)
		gl.Vertex3(this.pos.x+quadSize, this.pos.y-quadSize, this.pos.z)
		gl.Vertex3(this.pos.x+quadSize, this.pos.y+quadSize, this.pos.z)
		gl.Vertex3(this.pos.x-quadSize, this.pos.y+quadSize, this.pos.z)
	}
}

/**
 * Fragments of destroyed enemies.
 */
var fragmentDisplayList *DisplayList

type Fragment struct {
	field         Field
	pos, vel      Vector3
	size, d2, md2 float32
}

func InitFragments() {
	fragmentDisplayList = NewDisplayList(1)
	fragmentDisplayList.beginNewList()
	setScreenColor(0.7, 0.5, 0.5, 0.5)
	gl.Begin(gl.TRIANgl.E_FAN)
	gl.Vertex2(-0.5, -0.25)
	gl.Vertex2(0.5, -0.25)
	gl.Vertex2(0.5, 0.25)
	gl.Vertex2(-0.5, 0.25)
	gl.End()
	setScreenColor(0.7, 0.5, 0.5, 0.9)
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex2(-0.5, -0.25)
	gl.Vertex2(0.5, -0.25)
	gl.Vertex2(0.5, 0.25)
	gl.Vertex2(-0.5, 0.25)
	gl.End()
	fragmentDisplayList.endNewList()
}

func CloseFragments() {
	fragmentDisplayList.close()
}

func NewFragment(field Field) {
	this := new(Fragment)
	this.size = 1
	this.field = field
	return this
}

func (this *Fragment) set(p Vector, mx float32, my float32, mz float32, sz float32 /* = 1*/) {
	if !this.field.checkInOuterField(p.x, p.y) {
		return
	}
	this.pos.x = p.x
	this.pos.y = p.y
	this.pos.z = 0
	this.vel.x = mx
	this.vel.y = my
	this.vel.z = mz
	this.size = sz
	if this.size > 5 {
		this.size = 5
	}
	d2 = rand.nextFloat(360)
	md2 = rand.nextSignedFloat(20)
	this.exists = true
}

func (this *Fragment) move() {
	if !this.field.checkInOuterField(pos.x, pos.y) {
		this.exists = false
		return
	}
	this.vel.x *= 0.96
	this.vel.y *= 0.96
	this.vel.z += (-0.04 - this.vel.z) * 0.01
	this.pos += this.vel
	if this.pos.z < 0 {
		if this.field.getBlock(this.pos.x, this.pos.y) < 0 {
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, Smoke.SmokeType.WAKE, 60, this.size*0.66)
		} else {
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, Smoke.SmokeType.SAND, 60, this.size*0.75)
		}
		this.exists = false
		return
	}
	this.pos.y -= this.field.lastScrollY
	d2 += md2
}

func (this *Fragment) draw() {
	gl.PushMatrix()
	Screen.gl.Translate(this.pos)
	gl.Rotatef(d2, 1, 0, 0)
	gl.Scalef(this.size, this.size, 1)
	fragmentDisplayList.call(0)
	gl.PopMatrix()
}

/**
 * Luminous fragments.
 */
var sparkFragmentdisplayList *DisplayList

type SparkFragment struct {
	*LuminousActor

	field         Field
	pos, vel      Vector3
	size, d2, md2 float32
	cnt           int
	hasSmoke      bool
}

func InitSparkFragments() {
	sparkFragmentDisplayList = NewDisplayList(1)
	sparkFragmentDisplayList.beginNewList()
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex2(-0.25, -0.25)
	gl.Vertex2(0.25, -0.25)
	gl.Vertex2(0.25, 0.25)
	gl.Vertex2(-0.25, 0.25)
	gl.End()
	sparkFragmentDisplayList.endNewList()
}

func CloseSparkFragments() {
	sparkFragmentDisplayList.close()
}

func NewSparkFragment(field Field) *SparkFragment {
	this := &SparkFragment{NewLuminousActor()}
	this.size = 1
	this.field = field
	return this
}

func (this *SparkFragment) set(p Vector, mx float32, my float32, mz float32, sz float32 /*= 1*/) {
	if !this.field.checkInOuterField(p.x, p.y) {
		return
	}
	this.pos.x = p.x
	this.pos.y = p.y
	this.pos.z = 0
	this.vel.x = mx
	this.vel.y = my
	this.vel.z = mz
	this.size = sz
	if this.size > 5 {
		this.size = 5
	}
	d2 = rand.nextFloat(360)
	md2 = rand.nextSignedFloat(15)
	if rand.nextInt(4) == 0 {
		this.hasSmoke = true
	} else {
		this.hasSmoke = false
	}
	this.cnt = 0
	this.exists = true
}

func (this *SparkFragment) move() {
	if !this.field.checkInOuterField(this.pos.x, this.pos.y) {
		this.exists = false
		return
	}
	this.vel.x *= 0.99
	this.vel.y *= 0.99
	this.vel.z += (-0.08 - this.vel.z) * 0.01
	this.pos += vel
	if this.pos.z < 0 {
		if this.field.getBlock(this.pos.x, this.pos.y) < 0 {
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, Smoke.SmokeType.WAKE, 60, this.size*0.66)
		} else {
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, Smoke.SmokeType.SAND, 60, this.size*0.75)
		}
		this.exists = false
		return
	}
	this.pos.y -= this.field.lastScrollY
	d2 += md2
	this.cnt++
	if this.hasSmoke && this.cnt%5 == 0 {
		NewSmoke(this.pos, 0, 0, 0, Smoke.SmokeType.SMOKE, 90+rand.nextInt(60), this.size*0.5)
	}
}

func (this *SparkFragment) draw() {
	gl.PushMatrix()
	setScreenColor(1, rand.nextFloat(1), 0, 0.8)
	Screen.gl.Translate(this.pos)
	gl.Rotatef(d2, 1, 0, 0)
	gl.Scalef(this.size, this.size, 1)
	sparkFragmentDisplayList.call(0)
	gl.PopMatrix()
}

func (this *SparkFragment) drawLuminous() {
	gl.PushMatrix()
	setScreenColor(1, rand.nextFloat(1), 0, 0.8)
	Screen.gl.Translate(this.pos)
	gl.Rotatef(d2, 1, 0, 0)
	gl.Scalef(size, size, 1)
	sparkFragmentDisplayList.call(0)
	gl.PopMatrix()
}

/**
 * Wakes of ships and smokes.
 */
type Wake struct {
	field            Field
	pos, vel         Vector
	deg, speed, size float32
	cnt              int
	revShape         bool
}

func NewWake(field Field) *Wake {
	this := new(Wake)
	this.size = 1
	this.field = field
	return this
}

func (this *Wake) set(p Vector, deg float32, speed float32, c int /*= 60*/, sz float32 /*= 1*/, rs bool /* = false */) {
	if !this.field.checkInOuterField(p.x, p.y) {
		return
	}
	this.pos.x = p.x
	this.pos.y = p.y
	this.deg = deg
	this.speed = speed
	this.vel.x = sin(deg) * speed
	this.vel.y = cos(deg) * speed
	this.cnt = c
	this.size = sz
	this.revShape = rs
	this.exists = true
}

func (this *Wake) move() {
	this.cnt--
	if this.cnt <= 0 || this.vel.dist() < 0.005 || !this.field.checkInOuterField(this.pos.x, this.pos.y) {
		this.exists = false
		return
	}
	this.pos += this.vel
	this.pos.y -= this.field.lastScrollY
	this.vel *= 0.96
	this.size *= 1.02
}

func (this *Wake) draw() {
	ox := this.vel.x
	oy := this.vel.y
	setScreenColor(0.33, 0.33, 1)
	ox *= this.size
	oy *= this.size
	if this.revShape {
		gl.Vertex3(this.pos.x+ox, this.pos.y+oy, 0)
	} else {
		gl.Vertex3(this.pos.x-ox, this.pos.y-oy, 0)
	}
	ox *= 0.2
	oy *= 0.2
	setScreenColor(0.2, 0.2, 0.6, 0.5)
	gl.Vertex3(this.pos.x-oy, this.pos.y+ox, 0)
	gl.Vertex3(this.pos.x+oy, this.pos.y-ox, 0)
}
