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


func NewTurret(field Field, ship Ship, parent *Enemy, spec TurretSpec) *Turret{
	this := new(Turret)
	this.field = field
	this.ship = ship
	this.parent = parent
	this.bulletSpeed = 1
	this.spec = spec
	this.shield = spec.shield
	this.destroyedCnt = -1
	this.bulletSpeed = 1
	return this
}

func (this *Turret) move(x float32, y float32, d float32, bulletFireSpeed float32 /*= 0*/, bulletFireDeg float32 /*= -99999*/) bool {
	pos.x = x
	pos.y = y
	baseDeg = d
	if (destroyedCnt >= 0) {
		destroyedCnt++
		itv := 5 + destroyedCnt / 12
		if (itv < 60 && destroyedCnt % itv == 0) {
			NewSmoke(pos, 0, 0, 0.01 + rand.nextfloat32(0.01), SmokeType.FIRE, 90 + rand.nextInt(30), spec.size)
		}
		return false
	}
	td := baseDeg + deg
	shipPos := ship.nearPos(pos)
	shipVel := ship.nearVel(pos)
	ax := shipPos.x - pos.x
	ay := shipPos.y - pos.y
	if (spec.lookAheadRatio != 0) {
		rd := pos.dist(shipPos) / spec.speed * 1.2
		ax += shipVel.x * spec.lookAheadRatio * rd
		ay += shipVel.y * spec.lookAheadRatio * rd
	}
	var ad float32
	if (fabs32(ax) + fabs32(ay) < 0.1) {
		ad = 0
	} else {
		ad = atan2(ax, ay)
	}
	od := td - ad
	normalizeDeg(od)
	var ts float32
	if (cnt >= 0) {
		ts = spec.turnSpeed
	} else {
		ts = spec.turnSpeed * spec.burstTurnRatio
	}
	if (fabs32(od) <= ts) {
		deg = ad - baseDeg
	} else if (od > 0) {
		deg -= ts
	} else {
		deg += ts
	}
	normalizeDeg(deg)
	if (deg > spec.turnRange) {
		deg = spec.turnRange
	} else if (deg < -spec.turnRange) {
		deg = -spec.turnRange
	}
	cnt++
	if (field.checkInField(pos) || (parent.isBoss && cnt % 4 == 0)) {
		appCnt++
	}
	if (cnt >= spec.interval) {
		if (spec.blind || (fabs32(od) <= spec.turnSpeed && 
											 pos.dist(shipPos) < spec.maxRange * 1.1 &&
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
		bd := baseDeg + deg
		Smoke s = smokes.getInstance()
		if (s) {
			s.set(pos, sin(bd) * bulletSpeed, cos(bd) * bulletSpeed, 0,
						Smoke.SmokeType.SPARK, 20, spec.size * 2)
		}
		nw := spec.nway
		if (spec.nwayChange && burstCnt % 2 == 1) {
			nw--
		}
		bd -= spec.nwayAngle * (nw - 1) / 2
		for (int i = 0; i < nw; i++) {
			Bullet b = bullets.getInstance()
			if (!b) {
				break
			}
			b.set(parent.index,
						pos, bd, bulletSpeed, spec.size * 3, spec.bulletShape, spec.maxRange,
						bulletFireSpeed, bulletFireDeg, spec.bulletDestructive)
			bd += spec.nwayAngle
		}
		bulletSpeed += spec.speedAccel
		burstCnt++
	}
	damaged = false
	if (damagedCnt > 0) {
		damagedCnt--
	}
	startCnt++
	return true
}

func (this *Turret) draw() {
	if (spec.invisible) {
		return
	}
	glPushMatrix()
	if (destroyedCnt < 0 && damagedCnt > 0) { 
		damagedPos.x = pos.x + rand.nextSignedfloat32(damagedCnt * 0.015)
		damagedPos.y = pos.y + rand.nextSignedfloat32(damagedCnt * 0.015)
		Screen.glTranslate(damagedPos)
	} else {
		Screen.glTranslate(pos)
	}
	glRotatef(-(baseDeg + deg) * 180 / Pi32, 0, 0, 1)
	if (destroyedCnt >= 0) {
		spec.destroyedShape.draw()
	} else if (!damaged) {
		spec.shape.draw()
	} else {
		spec.damagedShape.draw()
	}
	glPopMatrix()
	if (destroyedCnt >= 0) {
		return
	}
	if (appCnt > 120) {
		return
	}
	a := 1 - float32(appCnt) / 120
	if (startCnt < 12) {
		a = float32(startCnt) / 12
	}
	td := baseDeg + deg
	if (spec.nway <= 1) { 
		glBegin(GL_LINE_STRIP)
		Screen.setColor(0.9, 0.1, 0.1, a)
		glVertex2(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
		Screen.setColor(0.9, 0.1, 0.1, a * 0.5)
		glVertex2(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
		glEnd()
	} else {
		td -= spec.nwayAngle * (spec.nway - 1) / 2
		glBegin(GL_LINE_STRIP)
		Screen.setColor(0.9, 0.1, 0.1, a * 0.75)
		glVertex2(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
		Screen.setColor(0.9, 0.1, 0.1, a * 0.25)
		glVertex2(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
		glEnd()
		glBegin(GL_QUADS)
		for (int i = 0; i < spec.nway - 1; i++) {
			Screen.setColor(0.9, 0.1, 0.1, a * 0.3)
			glVertex2(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
			Screen.setColor(0.9, 0.1, 0.1, a * 0.05)
			glVertex2(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
			td += spec.nwayAngle
			glVertex2(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
			Screen.setColor(0.9, 0.1, 0.1, a * 0.3)
			glVertex2(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
		}
		glEnd()
		glBegin(GL_LINE_STRIP)
		Screen.setColor(0.9, 0.1, 0.1, a * 0.75)
		glVertex2(pos.x + sin(td) * spec.minRange, pos.y + cos(td) * spec.minRange)
		Screen.setColor(0.9, 0.1, 0.1, a * 0.25)
		glVertex2(pos.x + sin(td) * spec.maxRange, pos.y + cos(td) * spec.maxRange)
		glEnd()
	}
}

func (this *Turret) checkCollision(x float32, y float32, c Collidable, shot Shot) bool {
	if (destroyedCnt >= 0 || spec.invisible) {
		return false
	}
	ox := fabs32(pos.x - x), oy = fabs32(pos.y - y)
	if (spec.shape.checkCollision(ox, oy, c)) {
		addDamage(shot.damage)
		return true
	}
	return false
}

func (this *Turret) addDamage(n int) {
	shield -= n
	if (shield <= 0) {
		destroyed()
	}
	damaged = true
	damagedCnt = 10
}

func (this *Turret) destroyed() {
	SoundManager.playSe("turret_destroyed.wav")
	destroyedCnt = 0
	for (int i = 0; i < 6; i++) {
		Smoke s = smokes.getInstanceForced()
		s.set(pos, rand.nextSignedfloat32(0.1), rand.nextSignedfloat32(0.1), rand.nextfloat32(0.04),
					Smoke.SmokeType.EXPLOSION, 30 + rand.nextInt(20), spec.size * 1.5)
	}
	for (int i = 0; i < 32; i++) {
		Spark sp = sparks.getInstanceForced()
		sp.set(pos, rand.nextSignedfloat32(0.5), rand.nextSignedfloat32(0.5),
					 0.5 + rand.nextfloat32(0.5), 0.5 + rand.nextfloat32(0.5), 0, 30 + rand.nextInt(30))
	}
	for (int i = 0; i < 7; i++) {
		Fragment f = fragments.getInstanceForced()
		f.set(pos, rand.nextSignedfloat32(0.25), rand.nextSignedfloat32(0.25), 0.05 + rand.nextfloat32(0.05),
					spec.size * (0.5 + rand.nextfloat32(0.5)))
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
	if (destroyedCnt < 0) {
		destroyedCnt = 999
	}
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

func NewTurretSpec() *TurretSpec {
	this := new(TurretSpec)
	this.shape = NewTurretShape(TurretShapeTypeNORMAL)
	this.damagedShape = NewTurretShape(TurretShapeTypeDAMAGED)
	this.destroyedShape = NewTurretShape(TurretShapeTypeDESTROYED)
	this.type = TurretTypeMAIN
	this.speed = 1
	this.maxRange = 99999
	this.turnSpeed = 99999
	this.turnRange = 99999
	this.burstNum = 1
	this.burstInterval = 99999
	this.nway = 1
	this.bulletShape = BulletShape.BulletShapeType.NORMAL
	this.shield = 99999
	this.size = 1
	return this
}	



func (this *TurretSpect) setParam(ts TurretSpec) {
	this.type = ts.type
	this.interval = ts.interval
	this.speed = ts.speed
	this.speedAccel = ts.speedAccel
	this.minRange = ts.minRange
	this.maxRange = ts.maxRange
	this.turnSpeed = ts.turnSpeed
	this.turnRange = ts.turnRange
	this.burstNum = ts.burstNum
	this.burstInterval = ts.burstInterval
	this.burstTurnRatio = ts.burstTurnRatio
	this.blind = ts.blind
	this.lookAheadRatio = ts.lookAheadRatio
	this.nway = ts.nway
	this.nwayAngle = ts.nwayAngle
	this.nwayChange = ts.nwayChange
	this.bulletShape = ts.bulletShape
	this.bulletDestructive = ts.bulletDestructive
	this.shield = ts.shield
	this.invisible = ts.invisible
	this.size = ts.size
}

func (this *TurretSpect) setParam(rank float32, type int) {
	init()
	this.type = type
	if (type == TurretType.DUMMY) {
		invisible = true
		return
	}
	rk := rank
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
		rk *= (10.0 / sqrt(maxRange))
		break
	default:
		maxRange = 9 + rand.nextfloat32(16)
		minRange = maxRange / (4 + rand.nextfloat32(0.5))
		if (type == TurretType.SUB || type == TurretType.SUB_DESTRUCTIVE) {
			maxRange *= 0.72
			minRange *= 0.9
		}
		rk *= (10.0 / sqrt(maxRange))
		if (rand.nextInt(4) == 0) {
			lar := rank * 0.1
			if (lar > 1) {
				lar = 1
			}
			lookAheadRatio = rand.nextfloat32(lar / 2) + lar / 2
			rk /= (1 + lookAheadRatio * 0.3)
		}
		if (rand.nextInt(3) == 0 && lookAheadRatio == 0) {
			blind = false
			rk *= 1.5
		} else {
			blind = true
		}
		turnRange = Pi32 / 4 + rand.nextfloat32(Pi32 / 4)
		turnSpeed = 0.005 + rand.nextfloat32(0.015)
		if (type == TurretType.MAIN) {
			turnRange *= 1.2
		}
		if (rand.nextInt(4) == 0) {
			burstTurnRatio = rand.nextfloat32(0.66) + 0.33
		}
		break
	}
	burstInterval = 6 + rand.nextInt(8)
	switch (type) {
	case TurretType.MAIN:
		size = 0.42 + rand.nextfloat32(0.05)
		br := (rk * 0.3) * (1 + rand.nextSignedfloat32(0.2))
		nr := (rk * 0.33) * rand.nextfloat32(1)
		ir := (rk * 0.1) * (1 + rand.nextSignedfloat32(0.2))
		burstNum = cast(int) br + 1
		nway = cast(int) (nr * 0.66 + 1)
		interval = cast(int) (120.0 / (ir * 2 + 1)) + 1
		sr := rk - burstNum + 1 - (nway - 1) / 0.66 - ir
		if (sr < 0) {
			sr = 0
		}
		speed = sqrt(sr * 0.6)
		speed *= 0.12
		shield = 20
		break
	case TurretType.SUB:
		size = 0.36 + rand.nextfloat32(0.025)
		br := (rk * 0.4) * (1 + rand.nextSignedfloat32(0.2))
		nr := (rk * 0.2) * rand.nextfloat32(1)
		ir := (rk * 0.2) * (1 + rand.nextSignedfloat32(0.2))
		burstNum = cast(int) br + 1
		nway = cast(int) (nr * 0.66 + 1)
		interval = cast(int) (120.0 / (ir * 2 + 1)) + 1
		sr := rk - burstNum + 1 - (nway - 1) / 0.66 - ir
		if (sr < 0) {
			sr = 0
		}
		speed = sqrt(sr * 0.7)
		speed *= 0.2
		shield = 12
		break
	case TurretType.SUB_DESTRUCTIVE:
		size = 0.36 + rand.nextfloat32(0.025)
		br := (rk * 0.4) * (1 + rand.nextSignedfloat32(0.2))
		nr := (rk * 0.2) * rand.nextfloat32(1)
		ir := (rk * 0.2) * (1 + rand.nextSignedfloat32(0.2))
		burstNum = cast(int) br * 2 + 1
		nway = cast(int) (nr * 0.66 + 1)
		interval = cast(int) (60.0 / (ir * 2 + 1)) + 1
		burstInterval *= 0.88
		bulletShape = BulletShape.BulletShapeType.DESTRUCTIVE
		bulletDestructive = true
		sr := rk - (burstNum - 1) / 2 - (nway - 1) / 0.66 - ir
		if (sr < 0) {
			sr = 0
		}
		speed = sqrt(sr * 0.7)
		speed *= 0.33
		shield = 12
		break
	case TurretType.SMALL:
		size = 0.33
		br := (rk * 0.33) * (1 + rand.nextSignedfloat32(0.2))
		ir := (rk * 0.2) * (1 + rand.nextSignedfloat32(0.2))
		burstNum = cast(int) br + 1
		nway = 1
		interval = cast(int) (120.0 / (ir * 2 + 1)) + 1
		sr := rk - burstNum + 1 - ir
		if (sr < 0) {
			sr = 0
		}
		speed = sqrt(sr)
		speed *= 0.24
		break
	case TurretType.MOVING:
		size = 0.36
		br := (rk * 0.3) * (1 + rand.nextSignedfloat32(0.2))
		nr := (rk * 0.1) * rand.nextfloat32(1)
		ir := (rk * 0.33) * (1 + rand.nextSignedfloat32(0.2))
		burstNum = cast(int) br + 1
		nway = cast(int) (nr * 0.66 + 1)
		interval = cast(int) (120.0 / (ir * 2 + 1)) + 1
		sr := rk - burstNum + 1 - (nway - 1) / 0.66 - ir
		if (sr < 0) {
			sr = 0
		}
		speed = sqrt(sr * 0.7)
		speed *= 0.2
		break
	}
	if (speed < 0.1) {
		speed = 0.1
	} else {
		speed = sqrt(speed * 10) / 10
	}
	if (burstNum > 2) {
		if (rand.nextInt(4) == 0) {
			speed *= 0.8
			burstInterval *= 0.7
			speedAccel = (speed * (0.4 + rand.nextfloat32(0.3))) / burstNum
			if (rand.nextInt(2) == 0) {
				speedAccel *= -1
			}
			speed -= speedAccel * burstNum / 2
		}
		if (rand.nextInt(5) == 0) {
			if (nway > 1) {
				nwayChange = true
			}
		}
	}
	nwayAngle = (0.1 + rand.nextfloat32(0.33)) / (1 + nway * 0.1)
}

func (this *TurretSpect) setBossSpec() {
	minRange = 0
	maxRange *= 1.5
	shield *= 2.1
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
  turret [TURRET_GROUPMAX_NUM]*Turret
  cnt int
}


func NewTurretGroup( field Field, ship Ship, parent Enemy, spec TurretGroupSpec) *TurretGroup {
	this := new(TurretGroup)
	this.ship = ship
	for i,_ := range this.turret {
		this.turret[i] = NewTurret(field, bullets, ship, sparks, smokes, fragments, parent)
	}
	this.spec = spec
	return this
}

func (this *TurretGroup) move(p Vector, deg float32) bool {
	alive := false
	centerPos.x = p.x
	centerPos.y = p.y
	var d, md, y, my float32
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
		var tbx, tby float32
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
			break
		}
		tbx *= (1 - spec.distRatio)
		bx := tbx * cos(-deg) - tby * sin(-deg)
		by := tbx * sin(-deg) + tby * cos(-deg)
		alive |= turret[i].move(centerPos.x + bx, centerPos.y + by, d + deg)
		if (spec.alignType == TurretGroupSpec.AlignType.ROUND) {
			d += md
		}
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
	col := false
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

func NewTurretGroupSpec() *TurretGroupSpec {
	this := new(TurretGroupSpec)
	this.num = 1
	this.alignType = AlignType.ROUND
	return this
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

func NewMovingTurretGroup(field Field, ship Ship, parent Enemy, spec MovingTurretGroupSpec) *MovingTurretGroup{
	this := new(MovingTurretGroup)
	this.ship = ship
	for i,_ := range this.turret {
		this.turret[i] = NewTurret(field, ship,  parent)
	}
	this.spec = spec
	this.radius = spec.radiusBase
	this.swingFixDeg = Pi32
	return this
}

func (this *MovingTurretGroupt) move(p Vector, od float32) {
	if (spec.moveType == MovingTurretGroupSpec.MoveType.SWING_FIX) {
		swingFixDeg = ed
	}
	centerPos.x = p.x
	centerPos.y = p.y
	if (spec.radiusAmp > 0) {
		radiusAmpCnt += spec.radiusAmpVel
		av := sin(radiusAmpCnt)
		radius = spec.radiusBase + spec.radiusAmp * av
	}
	if (spec.moveType == MovingTurretGroupSpec.MoveType.ROLL) {
		if (spec.rollAmp != 0) {
			rollAmpCnt += spec.rollAmpVel
			av := sin(rollAmpCnt)
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
			var od float32
			Vector shipPos = ship.nearPos(centerPos)
			if (shipPos.dist(centerPos) < 0.1) {
				od = 0
			} else {
				od = atan2(shipPos.x - centerPos.x, shipPos.y - centerPos.y)
			od += swingAmpDeg - deg
			normalizeDeg(od)
			deg += od * 0.1
		} else {
			od := swingFixDeg + swingAmpDeg - deg
			normalizeDeg(od)
			deg += od * 0.1
		}
	}
	float32 d, ad, md
	calcAlignDeg(d, ad, md)
	for (int i = 0; i < spec.num; i++) {
		d += md
		bx := sin(d) * radius * spec.xReverse
		by := cos(d) * radius * (1 - spec.distRatio)
		var fs, fd float32
		if (fabs32(bx) + fabs32(by) < 0.1) {
			fs = radius
			fd = d
		} else {
			fs = sqrt(bx * bx + by * by)
			fd = atan2(bx, by)
		}
		fs *= 0.06
		turret[i].move(centerPos.x, centerPos.y, d, fs, fd)
	}
	cnt++
}

func (this *MovingTurretGroupt) calcAlignDeg(d *float32, ad *float32, md *float32) {
	alignAmpCnt += spec.alignAmpVel
	ad = spec.alignDeg * (1 + sin(alignAmpCnt) * spec.alignAmp)
	if (spec.num > 1) {
		if (spec.moveType == MovingTurretGroupSpec.MoveType.ROLL) {
			md = ad / spec.num
		} else {
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

func NewMovingTurretGroupSpec() *MovingTurretGroupSpec {
	this := new(MovingTurretGroupSpec)
	this.num = 1
	this.initParam()
	this.num = 1
	this.alignDeg = Pi32 * 2
	this.radiusBase = 1
	this.moveType = TurretMoveType.SWING_FIX
	this.xReverse = 1
	return this
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
	if (aim) {
		moveType = MoveType.SWING_AIM
	} else {
		moveType = MoveType.SWING_FIX
	}
	swingDegVel = dv
	swingAmpVel = a
}

func (this *MovingTurretGroupSpec) setXReverse(xr float32) {
	xReverse = xr
}
