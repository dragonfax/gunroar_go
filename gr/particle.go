/*
 * $Id: particle.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/go-gl/gl"
)

/**
 * Sparks.
 */
type Spark struct {
	pos, ppos Vector
	vel       Vector
	r, g, b   float32
	cnt       int
}

func NewSpark(p Vector, vx float32, vy float32, r float32, g float32, b float32, c int) *Spark {
	this := new(Spark)
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
	actors[this] = true
	return this
}

func (this *Spark) move() {
	this.cnt--
	if this.cnt <= 0 || this.vel.dist(0, 0) < 0.005 {
		this.close()
		return
	}
	this.ppos.x = this.pos.x
	this.ppos.y = this.pos.y
	this.pos.AddAssign(this.vel)
	this.vel.MulAssign(0.96)
}

func (this *Spark) close() {
	delete(actors, this)
}

func (this *Spark) draw() {
	ox := this.vel.x
	oy := this.vel.y
	setScreenColor(this.r, this.g, this.b, 1)
	ox *= 2
	oy *= 2
	gl.Vertex3f(this.pos.x-ox, this.pos.y-oy, 0)
	ox *= 0.5
	oy *= 0.5
	setScreenColor(this.r*0.5, this.g*0.5, this.b*0.5, 0)
	gl.Vertex3f(this.pos.x-oy, this.pos.y+ox, 0)
	gl.Vertex3f(this.pos.x+oy, this.pos.y-ox, 0)
}

func (this *Spark) drawLuminous() {
	ox := this.vel.x
	oy := this.vel.y
	setScreenColor(this.r, this.g, this.b, 1)
	ox *= 2
	oy *= 2
	gl.Vertex3f(this.pos.x-ox, this.pos.y-oy, 0)
	ox *= 0.5
	oy *= 0.5
	setScreenColor(this.r*0.5, this.g*0.5, this.b*0.5, 0)
	gl.Vertex3f(this.pos.x-oy, this.pos.y+ox, 0)
	gl.Vertex3f(this.pos.x+oy, this.pos.y-ox, 0)
}

/**
 * Smokes.
 */
type SmokeType int

const (
	SmokeTypeFIRE SmokeType = iota
	SmokeTypeEXPLOSION
	SmokeTypeSAND
	SmokeTypeSPARK
	SmokeTypeWAKE
	SmokeTypeSMOKE
	SmokeTypeLANCE_SPARK
)

var windVel Vector3 = Vector3{0.04, 0.04, 0.02}

type Smoke struct {
	pos, vel         Vector3
	smokeType        SmokeType
	cnt, startCnt    int
	size, r, g, b, a float32
}

func NewSmoke(x float32, y float32, z float32 /*=0*/, mx float32, my float32, mz float32, t SmokeType, c int /* = 60 */, sz float32 /* = 2 */) *Smoke {
	this := new(Smoke)
	this.startCnt = 1
	this.size = 1
	actors[this] = true
	if !field.checkInOuterField(x, y) {
		return nil
	}
	this.pos.x = x
	this.pos.y = y
	this.pos.z = z
	this.vel.x = mx
	this.vel.y = my
	this.vel.z = mz
	this.smokeType = t
	this.startCnt = c
	this.cnt = c
	this.size = sz
	switch this.smokeType {
	case SmokeTypeFIRE:
		this.r = nextFloat(0.1) + 0.9
		this.g = nextFloat(0.2) + 0.2
		this.b = 0
		this.a = 1
		break
	case SmokeTypeEXPLOSION:
		this.r = nextFloat(0.3) + 0.7
		this.g = nextFloat(0.3) + 0.3
		this.b = 0
		this.a = 1
		break
	case SmokeTypeSAND:
		this.r = 0.8
		this.g = 0.8
		this.b = 0.6
		this.a = 0.6
		break
	case SmokeTypeSPARK:
		this.r = nextFloat(0.3) + 0.7
		this.g = nextFloat(0.5) + 0.5
		this.b = 0
		this.a = 1
		break
	case SmokeTypeWAKE:
		this.r = 0.6
		this.g = 0.6
		this.b = 0.8
		this.a = 0.6
		break
	case SmokeTypeSMOKE:
		this.r = nextFloat(0.1) + 0.1
		this.g = nextFloat(0.1) + 0.1
		this.b = 0.1
		this.a = 0.5
		break
	case SmokeTypeLANCE_SPARK:
		this.r = 0.4
		this.g = nextFloat(0.2) + 0.7
		this.b = nextFloat(0.2) + 0.7
		this.a = 1
		break
	}
	return this
}

func (this *Smoke) move() {
	this.cnt--
	if this.cnt <= 0 || !field.checkInOuterField(this.pos.x, this.pos.y) {
		this.close()
		return
	}
	if this.smokeType != SmokeTypeWAKE {
		this.vel.x += (windVel.x - this.vel.x) * 0.01
		this.vel.y += (windVel.y - this.vel.y) * 0.01
		this.vel.z += (windVel.z - this.vel.z) * 0.01
	}
	this.pos.AddAssign(this.vel)
	this.pos.y -= field.lastScrollY
	switch this.smokeType {
	case SmokeTypeFIRE, SmokeTypeEXPLOSION, SmokeTypeSMOKE:
		if this.cnt < this.startCnt/2 {
			this.r *= 0.95
			this.g *= 0.95
			this.b *= 0.95
		} else {
			this.a *= 0.97
		}
		this.size *= 1.01
		break
	case SmokeTypeSAND:
		this.r *= 0.98
		this.g *= 0.98
		this.b *= 0.98
		this.a *= 0.98
		break
	case SmokeTypeSPARK:
		this.r *= 0.92
		this.g *= 0.92
		this.a *= 0.95
		this.vel.MulAssign(0.9)
		break
	case SmokeTypeWAKE:
		this.a *= 0.98
		this.size *= 1.005
		break
	case SmokeTypeLANCE_SPARK:
		this.a *= 0.95
		this.size *= 0.97
		break
	}
	if this.size > 5 {
		this.size = 5
	}
	if this.smokeType == SmokeTypeEXPLOSION && this.pos.z < 0.01 {
		bl := field.getBlock(this.pos.x, this.pos.y)
		if bl >= 1 {
			this.vel.MulAssign(0.8)
		}
		if this.cnt%3 == 0 && bl < -1 {
			sp := sqrt32(this.vel.x*this.vel.x + this.vel.y*this.vel.y)
			if sp > 0.3 {
				d := atan232(this.vel.x, this.vel.y)
				wakePos.x = this.pos.x + Sin32(d+Pi32/2)*this.size*0.25
				wakePos.y = this.pos.y + Cos32(d+Pi32/2)*this.size*0.25
				NewWake(wakePos, d+Pi32-0.2+nextSignedFloat(0.1), sp*0.33,
					20+nextInt(12), this.size*(7.0+nextFloat(3)), false)
				wakePos.x = this.pos.x + Sin32(d-Pi32/2)*this.size*0.25
				wakePos.y = this.pos.y + Cos32(d-Pi32/2)*this.size*0.25
				NewWake(wakePos, d+Pi32+0.2+nextSignedFloat(0.1), sp*0.33,
					20+nextInt(12), this.size*(7.0+nextFloat(3)), false)
			}
		}
	}
}

func (this *Smoke) draw() {
	quadSize := this.size / 2
	setScreenColor(this.r, this.g, this.b, this.a)
	gl.Vertex3f(this.pos.x-quadSize, this.pos.y-quadSize, this.pos.z)
	gl.Vertex3f(this.pos.x+quadSize, this.pos.y-quadSize, this.pos.z)
	gl.Vertex3f(this.pos.x+quadSize, this.pos.y+quadSize, this.pos.z)
	gl.Vertex3f(this.pos.x-quadSize, this.pos.y+quadSize, this.pos.z)
}

func (this *Smoke) drawLuminous() {
	if this.r+this.g > 0.8 && this.b < 0.5 {
		quadSize := this.size / 2
		setScreenColor(this.r, this.g, this.b, this.a)
		gl.Vertex3f(this.pos.x-quadSize, this.pos.y-quadSize, this.pos.z)
		gl.Vertex3f(this.pos.x+quadSize, this.pos.y-quadSize, this.pos.z)
		gl.Vertex3f(this.pos.x+quadSize, this.pos.y+quadSize, this.pos.z)
		gl.Vertex3f(this.pos.x-quadSize, this.pos.y+quadSize, this.pos.z)
	}
}

func (this *Smoke) close() {
	delete(actors, this)
}

/**
 * Fragments of destroyed enemies.
 */
var fragmentDisplayList *DisplayList

type Fragment struct {
	pos, vel      Vector3
	size, d2, md2 float32
}

func InitFragments() {
	fragmentDisplayList = NewDisplayList(1)
	fragmentDisplayList.beginSingleList()
	setScreenColor(0.7, 0.5, 0.5, 0.5)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex2f(-0.5, -0.25)
	gl.Vertex2f(0.5, -0.25)
	gl.Vertex2f(0.5, 0.25)
	gl.Vertex2f(-0.5, 0.25)
	gl.End()
	setScreenColor(0.7, 0.5, 0.5, 0.9)
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex2f(-0.5, -0.25)
	gl.Vertex2f(0.5, -0.25)
	gl.Vertex2f(0.5, 0.25)
	gl.Vertex2f(-0.5, 0.25)
	gl.End()
	fragmentDisplayList.endSingleList()
}

func CloseFragments() {
	fragmentDisplayList.close()
}

func NewFragment(p Vector, mx float32, my float32, mz float32, sz float32 /* = 1*/) *Fragment {
	this := new(Fragment)
	this.size = 1

	if !field.checkInOuterField(p.x, p.y) {
		return nil
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
	this.d2 = nextFloat(360)
	this.md2 = nextSignedFloat(20)

	actors[this] = true
	return this
}

func (this *Fragment) move() {
	if !field.checkInOuterField(this.pos.x, this.pos.y) {
		this.close()
		return
	}
	this.vel.x *= 0.96
	this.vel.y *= 0.96
	this.vel.z += (-0.04 - this.vel.z) * 0.01
	this.pos.AddAssign(this.vel)
	if this.pos.z < 0 {
		if field.getBlock(this.pos.x, this.pos.y) < 0 {
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, 0, SmokeTypeWAKE, 60, this.size*0.66)
		} else {
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, 0, SmokeTypeSAND, 60, this.size*0.75)
		}
		this.close()
		return
	}
	this.pos.y -= field.lastScrollY
	this.d2 += this.md2
}

func (this *Fragment) draw() {
	gl.PushMatrix()
	glTranslate3(this.pos)
	gl.Rotatef(this.d2, 1, 0, 0)
	gl.Scalef(this.size, this.size, 1)
	fragmentDisplayList.call(0)
	gl.PopMatrix()
}

func (this *Fragment) close() {
	delete(actors, this)
}

/**
 * Luminous fragments.
 */
var sparkFragmentDisplayList *DisplayList

type SparkFragment struct {
	pos, vel      Vector3
	size, d2, md2 float32
	cnt           int
	hasSmoke      bool
}

func InitSparkFragments() {
	sparkFragmentDisplayList = NewDisplayList(1)
	sparkFragmentDisplayList.beginSingleList()
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex2f(-0.25, -0.25)
	gl.Vertex2f(0.25, -0.25)
	gl.Vertex2f(0.25, 0.25)
	gl.Vertex2f(-0.25, 0.25)
	gl.End()
	sparkFragmentDisplayList.endSingleList()
}

func CloseSparkFragments() {
	sparkFragmentDisplayList.close()
}

func NewSparkFragment(p Vector, mx float32, my float32, mz float32, sz float32 /*= 1*/) *SparkFragment {
	this := new(SparkFragment)
	this.size = 1

	if !field.checkInOuterField(p.x, p.y) {
		return nil
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
	this.d2 = nextFloat(360)
	this.md2 = nextSignedFloat(15)
	if nextInt(4) == 0 {
		this.hasSmoke = true
	} else {
		this.hasSmoke = false
	}
	this.cnt = 0

	actors[this] = true
	return this
}

func (this *SparkFragment) move() {
	if !field.checkInOuterField(this.pos.x, this.pos.y) {
		this.close()
		return
	}
	this.vel.x *= 0.99
	this.vel.y *= 0.99
	this.vel.z += (-0.08 - this.vel.z) * 0.01
	this.pos.AddAssign(this.vel)
	if this.pos.z < 0 {
		if field.getBlock(this.pos.x, this.pos.y) < 0 {
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, 0, SmokeTypeWAKE, 60, this.size*0.66)
		} else {
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, 0, SmokeTypeSAND, 60, this.size*0.75)
		}
		this.close()
		return
	}
	this.pos.y -= field.lastScrollY
	this.d2 += this.md2
	this.cnt++
	if this.hasSmoke && this.cnt%5 == 0 {
		NewSmoke(this.pos.x, this.pos.y, this.pos.z, 0, 0, 0, SmokeTypeSMOKE, 90+nextInt(60), this.size*0.5)
	}
}

func (this *SparkFragment) draw() {
	gl.PushMatrix()
	setScreenColor(1, nextFloat(1), 0, 0.8)
	glTranslate3(this.pos)
	gl.Rotatef(this.d2, 1, 0, 0)
	gl.Scalef(this.size, this.size, 1)
	sparkFragmentDisplayList.call(0)
	gl.PopMatrix()
}

func (this *SparkFragment) drawLuminous() {
	gl.PushMatrix()
	setScreenColor(1, nextFloat(1), 0, 0.8)
	glTranslate3(this.pos)
	gl.Rotatef(this.d2, 1, 0, 0)
	gl.Scalef(this.size, this.size, 1)
	sparkFragmentDisplayList.call(0)
	gl.PopMatrix()
}

func (this *SparkFragment) close() {
	delete(actors, this)
}

/**
 * Wakes of ships and smokes.
 */
type Wake struct {
	pos, vel         Vector
	deg, speed, size float32
	cnt              int
	revShape         bool
}

func NewWake(p Vector, deg float32, speed float32, c int /*= 60*/, sz float32 /*= 1*/, rs bool /* = false */) *Wake {
	this := new(Wake)
	this.size = 1

	if !field.checkInOuterField(p.x, p.y) {
		return nil
	}
	this.pos.x = p.x
	this.pos.y = p.y
	this.deg = deg
	this.speed = speed
	this.vel.x = Sin32(deg) * speed
	this.vel.y = Cos32(deg) * speed
	this.cnt = c
	this.size = sz
	this.revShape = rs
	actors[this] = true

	return this
}

func (this *Wake) move() {
	this.cnt--
	if this.cnt <= 0 || this.vel.dist(0, 0) < 0.005 || !field.checkInOuterField(this.pos.x, this.pos.y) {
		this.close()
		return
	}
	this.pos.AddAssign(this.vel)
	this.pos.y -= field.lastScrollY
	this.vel.MulAssign(0.96)
	this.size *= 1.02
}

func (this *Wake) draw() {
	ox := this.vel.x
	oy := this.vel.y
	setScreenColor(0.33, 0.33, 1, 1)
	ox *= this.size
	oy *= this.size
	if this.revShape {
		gl.Vertex3f(this.pos.x+ox, this.pos.y+oy, 0)
	} else {
		gl.Vertex3f(this.pos.x-ox, this.pos.y-oy, 0)
	}
	ox *= 0.2
	oy *= 0.2
	setScreenColor(0.2, 0.2, 0.6, 0.5)
	gl.Vertex3f(this.pos.x-oy, this.pos.y+ox, 0)
	gl.Vertex3f(this.pos.x+oy, this.pos.y-ox, 0)
}

func (this *Wake) close() {
	delete(actors, this)
}
