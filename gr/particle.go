/*
 * $Id: particle.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gtr

/**
 * Sparks.
 */
type Spark struct{
	*LuminousActor

  pos, ppos Vector
  vel Vector
  r, g, b float32
  cnt int
}

func NewSpark(p  Vector, vx float32, vy float32 , r float32 , g float32 , b float32 , c int ) *Spark {
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
	cnt--
	if (cnt <= 0 || vel.dist() < 0.005) {
		exists = false
		return
	}
	ppos.x = pos.x
	ppos.y = pos.y
	pos += vel
	vel *= 0.96
}

func (this *Spark) draw() {
	ox := vel.x
	oy := vel.y
	setScreenColor(r, g, b, 1)
	ox *= 2
	oy *= 2
	gl.Vertex3(pos.x - ox, pos.y - oy, 0)
	ox *= 0.5
	oy *= 0.5
	setScreenColor(r * 0.5, g * 0.5, b * 0.5, 0)
	gl.Vertex3(pos.x - oy, pos.y + ox, 0)
	gl.Vertex3(pos.x + oy, pos.y - ox, 0)
}

func (this *Spark) drawLuminous() {
	ox := vel.x
	oy := vel.y
	setScreenColor(r, g, b, 1)
	ox *= 2
	oy *= 2
	gl.Vertex3(pos.x - ox, pos.y - oy, 0)
	ox *= 0.5
	oy *= 0.5
	setScreenColor(r * 0.5, g * 0.5, b * 0.5, 0)
	gl.Vertex3(pos.x - oy, pos.y + ox, 0)
	gl.Vertex3(pos.x + oy, pos.y - ox, 0)
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

  field Field
  pos, vel Vector3
  smokeType SmokeType
  cnt, startCnt int
  size, r, g, b, a float32
}

func NewSmoke(field Field) *Smoke{
	this = &Smoke{NewLuminousActor()}
	this.startCnt = 1
	this.size = 1
	this.field = field
	return this
}

func (this *Smoke) set(p Vector , mx float32 , my float32 , mz float32 , t int , c int /*= 60*/, sz float32 /* = 2 */) {
	set(p.x, p.y, mx, my, mz, t, c, sz)
}

func (this *Smoke) set(p Vector3 , mx float32 , my float32 , mz float32 , t int , c int /*= 60*/, sz float32 /* = 2*/) {
	set(p.x, p.y, mx, my, mz, t, c, sz)
	pos.z = p.z
}

func (this *Smoke) set(x float32 , y float32 , mx float32 , my float32 , mz float32 , t int , c int /* = 60 */, sz float32 /* = 2 */) {
	if (!field.checkInOuterField(x, y)) {
		return
	}
	pos.x = x
	pos.y = y
	pos.z = 0
	vel.x = mx
	vel.y = my
	vel.z = mz
	smokeType = t
	startCnt = cnt = c
	size = sz
	switch (smokeType) {
	case SmokeType.FIRE:
		r = rand.nextFloat(0.1) + 0.9
		g = rand.nextFloat(0.2) + 0.2
		b = 0
		a = 1
		break
	case SmokeType.EXPLOSION:
		r = rand.nextFloat(0.3) + 0.7
		g = rand.nextFloat(0.3) + 0.3
		b = 0
		a = 1
		break
	case SmokeType.SAND:
		r = 0.8
		g = 0.8
		b = 0.6
		a = 0.6
		break
	case SmokeType.SPARK:
		r = rand.nextFloat(0.3) + 0.7
		g = rand.nextFloat(0.5) + 0.5
		b = 0
		a = 1
		break
	case SmokeType.WAKE:
		r = 0.6
		g = 0.6
		b = 0.8
		a = 0.6
		break
	case SmokeType.SMOKE:
		r = rand.nextFloat(0.1) + 0.1
		g = rand.nextFloat(0.1) + 0.1
		b = 0.1
		a = 0.5
		break
	case SmokeType.LANCE_SPARK:
		r = 0.4
		g = rand.nextFloat(0.2) + 0.7
		b = rand.nextFloat(0.2) + 0.7
		a = 1
		break
	}
	exists = true
}

func (this *Smoke) move() {
	cnt--
	if (cnt <= 0 || !field.checkInOuterField(pos.x, pos.y)) {
		exists = false
		return
	}
	if (smokeType != SmokeType.WAKE) {
		vel.x += (windVel.x - vel.x) * 0.01
		vel.y += (windVel.y - vel.y) * 0.01
		vel.z += (windVel.z - vel.z) * 0.01
	}
	pos += vel
	pos.y -= field.lastScrollY
	switch (smokeType) {
	case SmokeType.FIRE:
	case SmokeType.EXPLOSION:
	case SmokeType.SMOKE:
		if (cnt < startCnt / 2) {
			r *= 0.95
			g *= 0.95
			b *= 0.95
		} else {
			a *= 0.97
		}
		size *= 1.01
		break
	case SmokeType.SAND:
		r *= 0.98
		g *= 0.98
		b *= 0.98
		a *= 0.98
		break
	case SmokeType.SPARK:
		r *= 0.92
		g *= 0.92
		a *= 0.95
		vel *= 0.9
		break
	case SmokeType.WAKE:
		a *= 0.98
		size *= 1.005
		break
	case SmokeType.LANCE_SPARK:
		a *= 0.95
		size *= 0.97
		break
	}
	if (size > 5) {
		size = 5
	}
	if (smokeType == SmokeType.EXPLOSION && pos.z < 0.01) {
		bl := field.getBlock(pos.x, pos.y)
		if (bl >= 1) {
			vel *= 0.8
		}
		if (cnt % 3 == 0 && bl < -1) {
			sp := sqrt(vel.x * vel.x + vel.y * vel.y)
			if (sp > 0.3) {
				d := atan2(vel.x, vel.y)
				assert(d <>= 0)
				wakePos.x = pos.x + sin(d + PI / 2) * size * 0.25
				wakePos.y = pos.y + cos(d + PI / 2) * size * 0.25
				Wake w = wakes.getInstanceForced()
				assert(wakePos.x <>= 0)
				assert(wakePos.y <>= 0)
				w.set(wakePos, d + PI - 0.2 + rand.nextSignedFloat(0.1), sp * 0.33,
							20 + rand.nextInt(12), size * (7.0 + rand.nextFloat(3)))
				wakePos.x = pos.x + sin(d - PI / 2) * size * 0.25
				wakePos.y = pos.y + cos(d - PI / 2) * size * 0.25
				w = wakes.getInstanceForced()
				assert(wakePos.x <>= 0)
				assert(wakePos.y <>= 0)
				w.set(wakePos, d + PI + 0.2 + rand.nextSignedFloat(0.1), sp * 0.33,
							20 + rand.nextInt(12), size * (7.0 + rand.nextFloat(3)))
			}
		}
	}
}

func (this *Smoke) draw() {
	quadSize := size / 2
	setScreenColor(r, g, b, a)
	gl.Vertex3(pos.x - quadSize, pos.y - quadSize, pos.z)
	gl.Vertex3(pos.x + quadSize, pos.y - quadSize, pos.z)
	gl.Vertex3(pos.x + quadSize, pos.y + quadSize, pos.z)
	gl.Vertex3(pos.x - quadSize, pos.y + quadSize, pos.z)
}

func (this *Smoke) drawLuminous() {
	if (r + g > 0.8 && b < 0.5) {
		quadSize := size / 2
		setScreenColor(r, g, b, a)
		gl.Vertex3(pos.x - quadSize, pos.y - quadSize, pos.z)
		gl.Vertex3(pos.x + quadSize, pos.y - quadSize, pos.z)
		gl.Vertex3(pos.x + quadSize, pos.y + quadSize, pos.z)
		gl.Vertex3(pos.x - quadSize, pos.y + quadSize, pos.z)
	}
}
}


/**
 * Fragments of destroyed enemies.
 */
var fragmentDisplayList *DisplayList

type Fragment struct {
  field Field
  pos, vel Vector3
  size, d2, md2 float32
}

func InitFragments() {
	fragmentDisplayList = new DisplayList(1)
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

func (this *Fragment) set( p Vector , mx float32 , my float32 , mz float32 , sz float32 /* = 1*/) {
	if (!field.checkInOuterField(p.x, p.y)) {
		return
	}
	pos.x = p.x
	pos.y = p.y
	pos.z = 0
	vel.x = mx
	vel.y = my
	vel.z = mz
	size = sz
	if (size > 5) {
		size = 5
	}
	d2 = rand.nextFloat(360)
	md2 = rand.nextSignedFloat(20)
	exists = true
}

func (this *Fragment) move() {
	if (!field.checkInOuterField(pos.x, pos.y)) {
		exists = false
		return
	}
	vel.x *= 0.96
	vel.y *= 0.96
	vel.z += (-0.04 - vel.z) * 0.01
	pos += vel
	if (pos.z < 0) {
		Smoke s = smokes.getInstanceForced()
		if (field.getBlock(pos.x, pos.y) < 0) {
			s.set(pos.x, pos.y, 0, 0, 0, Smoke.SmokeType.WAKE, 60, size * 0.66)
		} else {
			s.set(pos.x, pos.y, 0, 0, 0, Smoke.SmokeType.SAND, 60, size * 0.75)
		}
		exists = false
		return
	}
	pos.y -= field.lastScrollY
	d2 += md2
}

func (this *Fragment) draw() {
	gl.PushMatrix()
	Screen.gl.Translate(pos)
	gl.Rotatef(d2, 1, 0, 0)
	gl.Scalef(size, size, 1)
	displayList.call(0)
	gl.PopMatrix()
}


/**
 * Luminous fragments.
 */
var sparkFragmentdisplayList *DisplayList

type SparkFragment struct {
	*LuminousActor

  field Field
  pos, vel Vector3
  size, d2, md2 float32
  cnt int
  hasSmoke bool
}

func InitSparkFragments() {
	sparkFragmentDisplayList = new DisplayList(1)
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

func (this *SparkFragment) set(Vector p, float32 mx, float32 my, float32 mz, float32 sz = 1) {
	if (!field.checkInOuterField(p.x, p.y)) {
		return
	}
	pos.x = p.x
	pos.y = p.y
	pos.z = 0
	vel.x = mx
	vel.y = my
	vel.z = mz
	size = sz
	if (size > 5) {
		size = 5
	}
	d2 = rand.nextFloat(360)
	md2 = rand.nextSignedFloat(15)
	if (rand.nextInt(4) == 0) {
		hasSmoke = true
	} else {
		hasSmoke = false
	}
	cnt = 0
	exists = true
}

func (this *SparkFragment) move() {
	if (!field.checkInOuterField(pos.x, pos.y)) {
		exists = false
		return
	}
	vel.x *= 0.99
	vel.y *= 0.99
	vel.z += (-0.08 - vel.z) * 0.01
	pos += vel
	if (pos.z < 0) {
		Smoke s = smokes.getInstanceForced()
		if (field.getBlock(pos.x, pos.y) < 0) {
			s.set(pos.x, pos.y, 0, 0, 0, Smoke.SmokeType.WAKE, 60, size * 0.66)
		} else {
			s.set(pos.x, pos.y, 0, 0, 0, Smoke.SmokeType.SAND, 60, size * 0.75)
		}
		exists = false
		return
	}
	pos.y -= field.lastScrollY
	d2 += md2
	cnt++
	if (hasSmoke && cnt % 5 == 0) {
		Smoke s = smokes.getInstance()
		if (s) {
			s.set(pos, 0, 0, 0, Smoke.SmokeType.SMOKE, 90 + rand.nextInt(60), size * 0.5)
		}
	}
}

func (this *SparkFragment) draw() {
	gl.PushMatrix()
	setScreenColor(1, rand.nextFloat(1), 0, 0.8)
	Screen.gl.Translate(pos)
	gl.Rotatef(d2, 1, 0, 0)
	gl.Scalef(size, size, 1)
	displayList.call(0)
	gl.PopMatrix()
}

func (this *SparkFragment) drawLuminous() {
	gl.PushMatrix()
	setScreenColor(1, rand.nextFloat(1), 0, 0.8)
	Screen.gl.Translate(pos)
	gl.Rotatef(d2, 1, 0, 0)
	gl.Scalef(size, size, 1)
	displayList.call(0)
	gl.PopMatrix()
}
}


/**
 * Wakes of ships and smokes.
 */
type Wake struct {
  field Field
  pos, vel Vector
  deg, speed, size float32
  cnt int
  revShape bool
}

func NewWake(field Field) *Wake {
	this := new(Wake)
	this.size = 1
	this.field = field
	return this
}

func (this *Wake) set(p Vector , deg float32 , speed float32 speed, c int /*= 60*/, sz float32 /*= 1*/ , rs bool /* = false */) {
	if (!field.checkInOuterField(p.x, p.y)) {
		return
	}
	pos.x = p.x
	pos.y = p.y
	this.deg = deg
	this.speed = speed
	vel.x = sin(deg) * speed
	vel.y = cos(deg) * speed
	cnt = c
	size = sz
	revShape = rs
	exists = true
}

func (this *Wake) move() {
	cnt--
	if (cnt <= 0 || vel.dist() < 0.005 || !field.checkInOuterField(pos.x, pos.y)) {
		exists = false
		return
	}
	pos += vel
	pos.y -= field.lastScrollY
	vel *= 0.96
	size *= 1.02
}

func (this *Wake) draw() {
	ox := vel.x
	oy := vel.y
	setScreenColor(0.33, 0.33, 1)
	ox *= size
	oy *= size
	if (revShape) {
		gl.Vertex3(pos.x + ox, pos.y + oy, 0)
	} else {
		gl.Vertex3(pos.x - ox, pos.y - oy, 0)
	}
	ox *= 0.2
	oy *= 0.2
	setScreenColor(0.2, 0.2, 0.6, 0.5)
	gl.Vertex3(pos.x - oy, pos.y + ox, 0)
	gl.Vertex3(pos.x + oy, pos.y - ox, 0)
}

