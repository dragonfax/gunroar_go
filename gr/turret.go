/*
 * $Id: turret.d,v 1.3 2005/07/17 11:02:46 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Turret mounted on a deck of an enemy ship.
 */

var damagedPos Vector

type Turret struct {
  field Field
  ship Ship
  spec TurretSpec
  pos Vector
  deg, baseDeg float32
  cnt, appCnt, startCnt, shield int
  damaged bool
  destroyedCnt, damagedCnt int
  bulletSpeed float32
  burstCnt int
  parent *Enemy
}


this(field Field, ship Ship, parent *Enemy) {
	this.field = field
	this.bullets = bullets
	this.ship = ship
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.parent = parent
	pos = new Vector
	deg = baseDeg = 0
	bulletSpeed = 1
}

func (this *Turret) start(spec TurretSpec) {
	this.spec = spec
	shield = spec.shield
	appCnt = cnt = startCnt = 0
	deg = baseDeg = 0
	damaged = false
	damagedCnt = 0
	destroyedCnt = -1
	bulletSpeed = 1
	burstCnt = 0
}

func (this *Turret) move(x float32, y float32, d float32, bulletFireSpeed float32 /*= 0*/, bulletFireDeg float32 /*= -99999*/) bool {
	pos.x = x
	pos.y = y
	baseDeg = d
	if (destroyedCnt >= 0) {
		destroyedCnt++
		int itv = 5 + destroyedCnt / 12
		if (itv < 60 && destroyedCnt % itv == 0) {
			Smoke s = smokes.getInstance()
			if (s)
				s.set(pos, 0, 0, 0.01f + rand.nextfloat32(0.01f), Smoke.SmokeType.FIRE, 90 + rand.nextInt(30), spec.size)
		}
		return false
	}
	float32 td = baseDeg + deg
	Vector shipPos = ship.nearPos(pos)
	Vector shipVel = ship.nearVel(pos)
	float32 ax = shipPos.x - pos.x
	float32 ay = shipPos.y - pos.y
	if (spec.lookAheadRatio != 0) {
		float32 rd = pos.dist(shipPos) / spec.speed * 1.2f
		ax += shipVel.x * spec.lookAheadRatio * rd
		ay += shipVel.y * spec.lookAheadRatio * rd
	}
	float32 ad
	if (fabs(ax) + fabs(ay) < 0.1f)
		ad = 0
	else
		ad = atan2(ax, ay)
	assert(ad <>= 0)
	float32 od = td - ad
	Math.normalizeDeg(od)
	float32 ts
	if (cnt >= 0)
		ts = spec.turnSpeed
	else
		ts = spec.turnSpeed * spec.burstTurnRatio
	if (fabs(od) <= ts)
		deg = ad - baseDeg
	else if (od > 0)
		deg -= ts
	else
		deg += ts
	Math.normalizeDeg(deg)
	if (deg > spec.turnRange)
		deg = spec.turnRange
	else if (deg < -spec.turnRange)
		deg = -spec.turnRange
	cnt++
	if (field.checkInField(pos) || (parent.isBoss && cnt % 4 == 0))
		appCnt++
	if (cnt >= spec.interval) {
		if (spec.blind || (fabs(od) <= spec.turnSpeed &&
											 pos.dist(shipPos) < spec.maxRange * 1.1f &&
											 pos.dist(shipPos) > spec.minRange)) {
			cnt = -(spec.burstNum - 1) * spec.burstInterval
			bulletSpeed = spec.speed
			burstCnt = 0
		}
	}
	if (cnt <= 0 && -cnt % spec.burstInterval == 0 &&
			((spec.invisible && field.checkInField(pos)) ||
			 (spec.invisible && parent.isBoss && field.checkInOuterField(pos)) ||
			 (!spec.invisible && field.checkInFieldExceptTop(pos))) &&
			pos.dist(shipPos) > spec.minRange) {
		float32 bd = baseDeg + deg
		Smoke s = smokes.getInstance()
		if (s)
			s.set(pos, sin(bd) * bulletSpeed, cos(bd) * bulletSpeed, 0,
						Smoke.SmokeType.SPARK, 20, spec.size * 2)
		int nw = spec.nway
		if (spec.nwayChange && burstCnt % 2 == 1)
			nw--
		bd -= spec.nwayAngle * (nw - 1) / 2
		for (int i = 0; i < nw; i++) {
			Bullet b = bullets.getInstance()
			if (!b)
				break
			b.set(parent.index,
						pos, bd, bulletSpeed, spec.size * 3, spec.bulletShape, spec.maxRange,
						bulletFireSpeed, bulletFireDeg, spec.bulletDestructive)
			bd += spec.nwayAngle
		}
		bulletSpeed += spec.speedAccel
		burstCnt++
	}
	damaged = false
	if (damagedCnt > 0)
		damagedCnt--
	startCnt++
	return true
}

func (this *Turret) draw() {
	if (spec.invisible)
		return
	glPushMatrix()
	if (destroyedCnt < 0 && damagedCnt > 0) {
		damagedPos.x = pos.x + rand.nextSignedfloat32(damagedCnt * 0.015f)
		damagedPos.y = pos.y + rand.nextSignedfloat32(damagedCnt * 0.015f)
		Screen.glTranslate(damagedPos)
	} else {
		Screen.glTranslate(pos)
	}
	glRotatef(-(baseDeg + deg) * 180 / PI, 0, 0, 1)
	if (destroyedCnt >= 0)
		spec.destroyedShape.draw()
	else if (!damaged)
		spec.shape.draw()
	else
		spec.damagedShape.draw()
	glPopMatrix()
	if (destroyedCnt >= 0)
		return
	if (appCnt > 120)
		return
	float32 a = 1 - cast(float32) appCnt / 120
	if (startCnt < 12)
		a = cast(float32) startCnt / 12
	float32 td = baseDeg + deg
	if (spec.nway <= 1) {
		glBegin(GL_LINE_STRIP)
		Screen.setColor(0.9f, 0.1f, 0.1f, a)
		glVertex2f(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
		Screen.setColor(0.9f, 0.1f, 0.1f, a * 0.5f)
		glVertex2f(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
		glEnd()
	} else {
		td -= spec.nwayAngle * (spec.nway - 1) / 2
		glBegin(GL_LINE_STRIP)
		Screen.setColor(0.9f, 0.1f, 0.1f, a * 0.75f)
		glVertex2f(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
		Screen.setColor(0.9f, 0.1f, 0.1f, a * 0.25f)
		glVertex2f(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
		glEnd()
		glBegin(GL_QUADS)
		for (int i = 0; i < spec.nway - 1; i++) {
			Screen.setColor(0.9f, 0.1f, 0.1f, a * 0.3f)
			glVertex2f(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
			Screen.setColor(0.9f, 0.1f, 0.1f, a * 0.05f)
			glVertex2f(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
			td += spec.nwayAngle
			glVertex2f(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
			Screen.setColor(0.9f, 0.1f, 0.1f, a * 0.3f)
			glVertex2f(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
		}
		glEnd()
		glBegin(GL_LINE_STRIP)
		Screen.setColor(0.9f, 0.1f, 0.1f, a * 0.75f)
		glVertex2f(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
		Screen.setColor(0.9f, 0.1f, 0.1f, a * 0.25f)
		glVertex2f(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
		glEnd()
	}
}

func (this *Turret) checkCollision(x float32, y float32, c Collidable, shot Shot) bool {
	if (destroyedCnt >= 0 || spec.invisible)
		return false
	float32 ox = fabs(pos.x - x), oy = fabs(pos.y - y)
	if (spec.shape.checkCollision(ox, oy, c)) {
		addDamage(shot.damage)
		return true
	}
	return false
}

func (this *Turret) addDamage(n int) {
	shield -= n
	if (shield <= 0)
		destroyed()
	damaged = true
	damagedCnt = 10
}

func (this *Turret) destroyed() {
	SoundManager.playSe("turret_destroyed.wav")
	destroyedCnt = 0
	for (int i = 0; i < 6; i++) {
		Smoke s = smokes.getInstanceForced()
		s.set(pos, rand.nextSignedfloat32(0.1f), rand.nextSignedfloat32(0.1f), rand.nextfloat32(0.04f),
					Smoke.SmokeType.EXPLOSION, 30 + rand.nextInt(20), spec.size * 1.5f)
	}
	for (int i = 0; i < 32; i++) {
		Spark sp = sparks.getInstanceForced()
		sp.set(pos, rand.nextSignedfloat32(0.5f), rand.nextSignedfloat32(0.5f),
					 0.5f + rand.nextfloat32(0.5f), 0.5f + rand.nextfloat32(0.5f), 0, 30 + rand.nextInt(30))
	}
	for (int i = 0; i < 7; i++) {
		Fragment f = fragments.getInstanceForced()
		f.set(pos, rand.nextSignedfloat32(0.25f), rand.nextSignedfloat32(0.25f), 0.05f + rand.nextfloat32(0.05f),
					spec.size * (0.5f + rand.nextfloat32(0.5f)))
	}
	switch (spec.type) {
	case TurretSpec.TurretType.MAIN:
		parent.increaseMultiplier(2)
		parent.addScore(40)
		break
	case TurretSpec.TurretType.SUB:
	case TurretSpec.TurretType.SUB_DESTRUCTIVE:
		parent.increaseMultiplier(1)
		parent.addScore(20)
		break
	}
}

func (this *Turret) remove() {
	if (destroyedCnt < 0)
		destroyedCnt = 999
}

/**
 * Turret specification changing according to a rank(difficulty).
 */
type TurretType int

const (
  TurretTypeMAIN TurretType = iota
  TurretTypeSUB 
  TurretTypeSUB_DESTRUCTIVE
  TurretTypeSMALL
  TurretTypeMOVING
  TurretTypeDUMMY
)

type TurretSpec struct {
  turretType TurretType
  interval int
  speed, speedAccel float32
  minRange, maxRange float32
  turnSpeed, turnRange float32
  burstNum, burstInterval int
  burstTurnRatio float32
  blind blool
  lookAheadRatio float32
  nway int
  nwayAngle float32
  nwayChange bool
  bulletShape int
  bulletDestructive bool
  shield int
  invisible bool
  shape, damagedShape, destroyedShape TurretShape
  size float32
}

this() {
	shape = new TurretShape(TurretShape.TurretShapeType.NORMAL)
	damagedShape = new TurretShape(TurretShape.TurretShapeType.DAMAGED)
	destroyedShape = new TurretShape(TurretShape.TurretShapeType.DESTROYED)
	init()
}

func (this *TurretSpect) init() {
	type = 0
	interval = 99999
	speed = 1
	speedAccel = 0
	minRange = 0
	maxRange = 99999
	turnSpeed = 99999
	turnRange = 99999
	burstNum = 1
	burstInterval = 99999
	burstTurnRatio = 0
	blind = false
	lookAheadRatio = 0
	nway = 1
	nwayAngle = 0
	nwayChange = false
	bulletShape = BulletShape.BulletShapeType.NORMAL
	bulletDestructive = false
	shield = 99999
	invisible = false
	_size = 1
}

func (this *TurretSpect) setParam(ts TurretSpec) {
	type = ts.type
	interval = ts.interval
	speed = ts.speed
	speedAccel = ts.speedAccel
	minRange = ts.minRange
	maxRange = ts.maxRange
	turnSpeed = ts.turnSpeed
	turnRange = ts.turnRange
	burstNum = ts.burstNum
	burstInterval = ts.burstInterval
	burstTurnRatio = ts.burstTurnRatio
	blind = ts.blind
	lookAheadRatio = ts.lookAheadRatio
	nway = ts.nway
	nwayAngle = ts.nwayAngle
	nwayChange = ts.nwayChange
	bulletShape = ts.bulletShape
	bulletDestructive = ts.bulletDestructive
	shield = ts.shield
	invisible = ts.invisible
	size = ts.size
}

func (this *TurretSpect) setParam(rank float32, type int) {
	init()
	this.type = type
	if (type == TurretType.DUMMY) {
		invisible = true
		return
	}
	float32 rk = rank
	switch (type) {
	case TurretType.SMALL:
		minRange = 8
		bulletShape = BulletShape.BulletShapeType.SMALL
		blind = true
		invisible = true
		break
	case TurretType.MOVING:
		minRange = 6
		bulletShape = BulletShape.BulletShapeType.MOVING_TURRET
		blind = true
		invisible = true
		turnSpeed = 0
		maxRange = 9 + rand.nextfloat32(12)
		rk *= (10.0f / sqrt(maxRange))
		break
	default:
		maxRange = 9 + rand.nextfloat32(16)
		minRange = maxRange / (4 + rand.nextfloat32(0.5f))
		if (type == TurretType.SUB || type == TurretType.SUB_DESTRUCTIVE) {
			maxRange *= 0.72f
			minRange *= 0.9f
		}
		rk *= (10.0f / sqrt(maxRange))
		if (rand.nextInt(4) == 0) {
			float32 lar = rank * 0.1f
			if (lar > 1)
				lar = 1
			lookAheadRatio = rand.nextfloat32(lar / 2) + lar / 2
			rk /= (1 + lookAheadRatio * 0.3f)
		}
		if (rand.nextInt(3) == 0 && lookAheadRatio == 0) {
			blind = false
			rk *= 1.5f
		} else {
			blind = true
		}
		turnRange = PI / 4 + rand.nextfloat32(PI / 4)
		turnSpeed = 0.005f + rand.nextfloat32(0.015f)
		if (type == TurretType.MAIN)
			turnRange *= 1.2f
		if (rand.nextInt(4) == 0)
			burstTurnRatio = rand.nextfloat32(0.66f) + 0.33f
		break
	}
	burstInterval = 6 + rand.nextInt(8)
	switch (type) {
	case TurretType.MAIN:
		size = 0.42f + rand.nextfloat32(0.05f)
		float32 br = (rk * 0.3f) * (1 + rand.nextSignedfloat32(0.2f))
		float32 nr = (rk * 0.33f) * rand.nextfloat32(1)
		float32 ir = (rk * 0.1f) * (1 + rand.nextSignedfloat32(0.2f))
		burstNum = cast(int) br + 1
		nway = cast(int) (nr * 0.66f + 1)
		interval = cast(int) (120.0f / (ir * 2 + 1)) + 1
		float32 sr = rk - burstNum + 1 - (nway - 1) / 0.66f - ir
		if (sr < 0)
			sr = 0
		speed = sqrt(sr * 0.6f)
		assert(speed <>= 0)
		speed *= 0.12f
		shield = 20
		break
	case TurretType.SUB:
		size = 0.36f + rand.nextfloat32(0.025f)
		float32 br = (rk * 0.4f) * (1 + rand.nextSignedfloat32(0.2f))
		float32 nr = (rk * 0.2f) * rand.nextfloat32(1)
		float32 ir = (rk * 0.2f) * (1 + rand.nextSignedfloat32(0.2f))
		burstNum = cast(int) br + 1
		nway = cast(int) (nr * 0.66f + 1)
		interval = cast(int) (120.0f / (ir * 2 + 1)) + 1
		float32 sr = rk - burstNum + 1 - (nway - 1) / 0.66f - ir
		if (sr < 0)
			sr = 0
		speed = sqrt(sr * 0.7f)
		assert(speed <>= 0)
		speed *= 0.2f
		shield = 12
		break
	case TurretType.SUB_DESTRUCTIVE:
		size = 0.36f + rand.nextfloat32(0.025f)
		float32 br = (rk * 0.4f) * (1 + rand.nextSignedfloat32(0.2f))
		float32 nr = (rk * 0.2f) * rand.nextfloat32(1)
		float32 ir = (rk * 0.2f) * (1 + rand.nextSignedfloat32(0.2f))
		burstNum = cast(int) br * 2 + 1
		nway = cast(int) (nr * 0.66f + 1)
		interval = cast(int) (60.0f / (ir * 2 + 1)) + 1
		burstInterval *= 0.88f
		bulletShape = BulletShape.BulletShapeType.DESTRUCTIVE
		bulletDestructive = true
		float32 sr = rk - (burstNum - 1) / 2 - (nway - 1) / 0.66f - ir
		if (sr < 0)
			sr = 0
		speed = sqrt(sr * 0.7f)
		assert(speed <>= 0)
		speed *= 0.33f
		shield = 12
		break
	case TurretType.SMALL:
		size = 0.33f
		float32 br = (rk * 0.33f) * (1 + rand.nextSignedfloat32(0.2f))
		float32 ir = (rk * 0.2f) * (1 + rand.nextSignedfloat32(0.2f))
		burstNum = cast(int) br + 1
		nway = 1
		interval = cast(int) (120.0f / (ir * 2 + 1)) + 1
		float32 sr = rk - burstNum + 1 - ir
		if (sr < 0)
			sr = 0
		speed = sqrt(sr)
		assert(speed <>= 0)
		speed *= 0.24f
		break
	case TurretType.MOVING:
		size = 0.36f
		float32 br = (rk * 0.3f) * (1 + rand.nextSignedfloat32(0.2f))
		float32 nr = (rk * 0.1f) * rand.nextfloat32(1)
		float32 ir = (rk * 0.33f) * (1 + rand.nextSignedfloat32(0.2f))
		burstNum = cast(int) br + 1
		nway = cast(int) (nr * 0.66f + 1)
		interval = cast(int) (120.0f / (ir * 2 + 1)) + 1
		float32 sr = rk - burstNum + 1 - (nway - 1) / 0.66f - ir
		if (sr < 0)
			sr = 0
		speed = sqrt(sr * 0.7f)
		assert(speed <>= 0)
		speed *= 0.2f
		break
	}
	if (speed < 0.1f)
		speed = 0.1f
	else
		speed = sqrt(speed * 10) / 10
	assert(speed <>= 0)
	if (burstNum > 2) {
		if (rand.nextInt(4) == 0) {
			speed *= 0.8f
			burstInterval *= 0.7f
			speedAccel = (speed * (0.4f + rand.nextfloat32(0.3f))) / burstNum
			if (rand.nextInt(2) == 0)
				speedAccel *= -1
			speed -= speedAccel * burstNum / 2
		}
		if (rand.nextInt(5) == 0) {
			if (nway > 1)
				nwayChange = true
		}
	}
	nwayAngle = (0.1f + rand.nextfloat32(0.33f)) / (1 + nway * 0.1f)
}

func (this *TurretSpect) setBossSpec() {
	minRange = 0
	maxRange *= 1.5f
	shield *= 2.1f
}

func (this *TurretSpect) float32 size() {
	return _size
}

func (this *TurretSpect) float32 size(float32 v) {
	_size = v
	shape.size = damagedShape.size = destroyedShape.size = _size
	return _size
}

/**
 * Grouped turrets.
 */
const TURRET_GROUP_MAX_NUM = 16
type TurretGroup  struct {
  ship Ship
  spec TurretGroupSpec
  centerPos Vector
  turret [TURRET_GROUPMAX_NUM]Turret
  cnt int
}


this(field Field, ship Ship, parent Enemy) {
	this.ship = ship
	centerPos = new Vector
	foreach (inout Turret t; turret)
		t = new Turret(field, bullets, ship, sparks, smokes, fragments, parent)
}

func (this *TurretGroup) set(spec TurretGroupSpec) {
	this.spec = spec
	for (int i = 0; i < spec.num; i++)
		turret[i].start(spec.turretSpec)
	cnt = 0
}

func (this *TurretGroup) move(p Vector, deg float32) bool {
	bool alive = false
	centerPos.x = p.x
	centerPos.y = p.y
	float32 d, md, y, my
	switch (spec.alignType) {
	case TurretGroupSpec.AlignType.ROUND:
		d = spec.alignDeg
		if (spec.num > 1) {
			md = spec.alignWidth / (spec.num - 1)
			d -= spec.alignWidth / 2
		} else {
			md = 0
		}
		break
	case TurretGroupSpec.AlignType.STRAIGHT:
		y = 0
		my = spec.offset.y / (spec.num + 1)
		break
	}
	for (int i = 0; i < spec.num; i++) {
		float32 tbx, tby
		switch (spec.alignType) {
		case TurretGroupSpec.AlignType.ROUND:
			tbx = sin(d) * spec.radius
			tby = cos(d) * spec.radius
			break
		case TurretGroupSpec.AlignType.STRAIGHT:
			y += my
			tbx = spec.offset.x
			tby = y
			d = atan2(tbx, tby)
			assert(d <>= 0)
			break
		}
		tbx *= (1 - spec.distRatio)
		float32 bx = tbx * cos(-deg) - tby * sin(-deg)
		float32 by = tbx * sin(-deg) + tby * cos(-deg)
		alive |= turret[i].move(centerPos.x + bx, centerPos.y + by, d + deg)
		if (spec.alignType == TurretGroupSpec.AlignType.ROUND)
			d += md
	}
	cnt++
	return alive
}

func (this *TurretGroup) draw() {
	for (int i = 0; i < spec.num; i++)
		turret[i].draw()
}

func (this *TurretGroup) remove() {
	for (int i = 0; i < spec.num; i++)
		turret[i].remove()
}

func (this *TurretGroup) checkCollision( x float32, y float32, c Collidable, shot Shot) bool {
	bool col = false
	for (int i = 0; i < spec.num; i++)
		col |= turret[i].checkCollision(x, y, c, shot)
	return col
}

type AlignType int

const(
	ROUND AlignType = iota
	STRAIGHT
)

type TurretGroupSpec struct {
  turretSpec TurretSpec
  num, alignType int
  alignDeg, alignWidth, radius, distRatio float32
  offset Vector
}

this() {
	turretSpec = new TurretSpec
	offset = new Vector
	num = 1
	alignDeg = alignWidth = 0
	radius = 0
	distRatio = 0
}

func (this *TurretGroupSpec) init() {
	num = 1
	alignType = AlignType.ROUND
	alignDeg = alignWidth = radius = distRatio = 0
	offset.x = offset.y = 0
}

/**
 * Turrets moving around a bridge.
 */

const MOVING_TURRET_MAX_NUM = 16
type MovingTurretGroup struct {
  ship Ship
  spec MovingTurretGroupSpec
  radius, radiusAmpCnt, deg, rollAmpCnt, swingAmpCnt, swingAmpDeg, swingFixDeg, alignAmpCnt, distDeg, distAmpCnt float32
  cnt int
  centerPos Vector
  turret [MOVING_TURRET_MAX_NUM]Turret
}

this(field Field, ship Ship, parent Enemy) {
	this.ship = ship
	centerPos = new Vector
	foreach (inout Turret t; turret)
		t = new Turret(field, bullets, ship, sparks, smokes, fragments, parent)
	radius = radiusAmpCnt = 0
	deg = 0
	rollAmpCnt = swingAmpCnt = swingAmpDeg = swingFixDeg = alignAmpCnt = 0
	distDeg = distAmpCnt = 0
}

func (this *MovingTurretGroupt) set(spec MovingTurretGroupSpec) {
	this.spec = spec
	radius = spec.radiusBase
	radiusAmpCnt = 0
	deg = 0
	rollAmpCnt = swingAmpCnt = swingAmpDeg = alignAmpCnt = 0
	distDeg = distAmpCnt = 0
	swingFixDeg = PI
	for (int i = 0; i < spec.num; i++)
		turret[i].start(spec.turretSpec)
	cnt = 0
}

func (this *MovingTurretGroupt) move(p Vector, od float32) {
	if (spec.moveType == MovingTurretGroupSpec.MoveType.SWING_FIX)
		swingFixDeg = ed
	centerPos.x = p.x
	centerPos.y = p.y
	if (spec.radiusAmp > 0) {
		radiusAmpCnt += spec.radiusAmpVel
		float32 av = sin(radiusAmpCnt)
		radius = spec.radiusBase + spec.radiusAmp * av
	}
	if (spec.moveType == MovingTurretGroupSpec.MoveType.ROLL) {
		if (spec.rollAmp != 0) {
			rollAmpCnt += spec.rollAmpVel
			float32 av = sin(rollAmpCnt)
			deg += spec.rollDegVel + spec.rollAmp * av
		} else {
			deg += spec.rollDegVel
		}
	} else {
		swingAmpCnt += spec.swingAmpVel
		if (cos(swingAmpCnt) > 0) {
			swingAmpDeg += spec.swingDegVel
		} else {
			swingAmpDeg -= spec.swingDegVel
		}
		if (spec.moveType == MovingTurretGroupSpec.MoveType.SWING_AIM) {
			float32 od
			Vector shipPos = ship.nearPos(centerPos)
			if (shipPos.dist(centerPos) < 0.1f)
				od = 0
			else
				od = atan2(shipPos.x - centerPos.x, shipPos.y - centerPos.y)
			assert(od <>= 0)
			od += swingAmpDeg - deg
			Math.normalizeDeg(od)
			deg += od * 0.1f
		} else {
			float32 od = swingFixDeg + swingAmpDeg - deg
			Math.normalizeDeg(od)
			deg += od * 0.1f
		}
	}
	float32 d, ad, md
	calcAlignDeg(d, ad, md)
	for (int i = 0; i < spec.num; i++) {
		d += md
		float32 bx = sin(d) * radius * spec.xReverse
		float32 by = cos(d) * radius * (1 - spec.distRatio)
		float32 fs, fd
		if (fabs(bx) + fabs(by) < 0.1f) {
			fs = radius
			fd = d
		} else {
			fs = sqrt(bx * bx + by * by)
			fd = atan2(bx, by)
			assert(fd <>= 0)
		}
		fs *= 0.06f
		turret[i].move(centerPos.x, centerPos.y, d, fs, fd)
	}
	cnt++
}

func (this *MovingTurretGroupt) calcAlignDeg(d *float32, ad *float32, md *float32) {
	alignAmpCnt += spec.alignAmpVel
	ad = spec.alignDeg * (1 + sin(alignAmpCnt) * spec.alignAmp)
	if (spec.num > 1) {
		if (spec.moveType == MovingTurretGroupSpec.MoveType.ROLL)
			md = ad / spec.num
		else
			md = ad / (spec.num - 1)
	} else {
		md = 0
	}
	d = deg - md - ad / 2
}

func (this *MovingTurretGroupt) draw() {
	for (int i = 0; i < spec.num; i++)
		turret[i].draw()
}

func (this *MovingTurretGroupt) remove() {
	for (int i = 0; i < spec.num; i++)
		turret[i].remove()
}

type TurretMoveType int 

const(
	TurretMoveTypeROLL MoveType = iota
	TurretMoveTypeSWING_FIX
	TurretMoveTypeSWING_AIM
)

type MovingTurretGroupSpec struct {
  TurretSpec turretSpec
  num int
	moveType TurretMoveType
  alignDeg, alignAmp, alignAmpVel, radiusBase, radiusAmp, radiusAmpVel, rollDegVel, rollAmp, rollAmpVel, swingDegVel, swingAmpVel, distRatio, xReverse float32
},

this() {
	turretSpec = new TurretSpec
	num = 1
	initParam()
}

func (this *MovingTurretGroupSpec) initParam() {
	num = 1
	alignDeg = PI * 2
	alignAmp = alignAmpVel = 0
	radiusBase = 1
	radiusAmp = radiusAmpVel = 0
	moveType = MoveType.SWING_FIX
	rollDegVel = rollAmp = rollAmpVel = 0
	swingDegVel = swingAmpVel = 0
	distRatio = 0
	xReverse = 1
}

func (this *MovingTurretGroupSpec) init() {
	initParam()
}

func (this *MovingTurretGroupSpec) setAlignAmp(a float32, v float32) {
	alignAmp = a
	alignAmpVel = v
}

func (this *MovingTurretGroupSpec) setRadiusAmp(a float32, v float32) {
	radiusAmp = a
	radiusAmpVel = v
}

func (this *MovingTurretGroupSpec) setRoll(dv float32, a float32, v float32) {
	moveType = MoveType.ROLL
	rollDegVel = dv
	rollAmp = a
	rollAmpVel = v
}

func (this *MovingTurretGroupSpec) setSwing(dv float32, a float32, aim bool /*= false*/) {
	if (aim)
		moveType = MoveType.SWING_AIM
	else
		moveType = MoveType.SWING_FIX
	swingDegVel = dv
	swingAmpVel = a
}

func (this *MovingTurretGroupSpec) setXReverse(xr float32) {
	xReverse = xr
}
