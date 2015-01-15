/*
 * $Id: turret.d,v 1.3 2005/07/17 11:02:46 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Turret mounted on a deck of an enemy ship.
 */

var turretDamagedPos Vector

type Turret struct {
	field                         Field
	ship                          Ship
	spec                          TurretSpec
	pos                           Vector
	deg, baseDeg                  float32
	cnt, appCnt, startCnt, shield int
	damaged                       bool
	destroyedCnt, damagedCnt      int
	bulletSpeed                   float32
	burstCnt                      int
	parent                        *Enemy
}

func NewTurret(field Field, ship Ship, parent *Enemy, spec TurretSpec) *Turret {
	this := new(Turret)
	field = field
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
	this.pos.x = x
	this.pos.y = y
	this.baseDeg = d
	if this.destroyedCnt >= 0 {
		this.destroyedCnt++
		itv := 5 + this.destroyedCnt/12
		if itv < 60 && this.destroyedCnt%itv == 0 {
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, 0.01+nextFloat(0.01), SmokeType.FIRE, 90+nextInt(30), this.spec.size)
		}
		return false
	}
	td := this.baseDeg + this.deg
	shipPos := this.ship.nearPos(this.pos)
	shipVel := this.ship.nearVel(this.pos)
	ax := shipPos.x - this.pos.x
	ay := shipPos.y - this.pos.y
	if this.spec.lookAheadRatio != 0 {
		rd := this.pos.dist(shipPos) / this.spec.speed * 1.2
		ax += shipVel.x * this.spec.lookAheadRatio * rd
		ay += shipVel.y * this.spec.lookAheadRatio * rd
	}
	var ad float32
	if fabs32(ax)+fabs32(ay) < 0.1 {
		ad = 0
	} else {
		ad = atan2(ax, ay)
	}
	od := td - ad
	normalizeDeg(od)
	var ts float32
	if this.cnt >= 0 {
		ts = this.spec.turnSpeed
	} else {
		ts = this.spec.turnSpeed * this.spec.burstTurnRatio
	}
	if fabs32(od) <= ts {
		this.deg = ad - this.baseDeg
	} else if od > 0 {
		this.deg -= ts
	} else {
		this.deg += ts
	}
	normalizeDeg(this.deg)
	if this.deg > this.spec.turnRange {
		this.deg = this.spec.turnRange
	} else if this.deg < -this.spec.turnRange {
		this.deg = -this.spec.turnRange
	}
	this.cnt++
	if field.checkInField(this.pos) || (this.parent.isBoss && this.cnt%4 == 0) {
		this.appCnt++
	}
	if this.cnt >= this.spec.interval {
		if this.spec.blind || (fabs32(od) <= this.spec.turnSpeed &&
			this.pos.dist(shipPos) < this.spec.maxRange*1.1 &&
			this.pos.dist(shipPos) > this.spec.minRange) {
			this.cnt = -(this.spec.burstNum - 1) * this.spec.burstInterval
			this.bulletSpeed = this.spec.speed
			this.burstCnt = 0
		}
	}
	if this.cnt <= 0 && -this.cnt%this.spec.burstInterval == 0 &&
		((this.spec.invisible && field.checkInField(this.pos)) ||
			(this.spec.invisible && this.parent.isBoss && field.checkInOuterField(this.pos)) ||
			(!this.spec.invisible && field.checkInFieldExceptTop(this.pos))) &&
		this.pos.dist(shipPos) > this.spec.minRange {
		bd := this.baseDeg + this.deg
		NewSmoke(this.pos.x, this.pos.y, 0, Sin32(bd)*this.bulletSpeed, Cos32(bd)*this.bulletSpeed, 0,
			Smoke.SmokeType.SPARK, 20, this.spec.size*2)
		nw := this.spec.nway
		if this.spec.nwayChange && this.burstCnt%2 == 1 {
			nw--
		}
		bd -= this.spec.nwayAngle * (nw - 1) / 2
		for i := 0; i < nw; i++ {
			NewBullet(this.parent.index,
				this.pos, bd, this.bulletSpeed, this.spec.size*3, this.spec.bulletShape, this.spec.maxRange,
				this.bulletFireSpeed, this.bulletFireDeg, this.spec.bulletDestructive)
			bd += this.spec.nwayAngle
		}
		this.bulletSpeed += this.spec.speedAccel
		this.burstCnt++
	}
	this.damaged = false
	if this.damagedCnt > 0 {
		this.damagedCnt--
	}
	this.startCnt++
	return true
}

func (this *Turret) draw() {
	if this.spec.invisible {
		return
	}
	gl.PushMatrix()
	if this.destroyedCnt < 0 && this.damagedCnt > 0 {
		turretDamagedPos.x = this.pos.x + nextSignedFloat(this.damagedCnt*0.015)
		turretDamagedPos.y = this.pos.y + nextSignedFloat(this.damagedCnt*0.015)
		gl.Translate(turretDamagedPos)
	} else {
		gl.Translate(this.pos)
	}
	gl.Rotatef(-(this.baseDeg+this.deg)*180/Pi32, 0, 0, 1)
	if this.destroyedCnt >= 0 {
		this.spec.destroyedShape.draw()
	} else if !this.damaged {
		this.spec.shape.draw()
	} else {
		this.spec.damagedShape.draw()
	}
	gl.PopMatrix()
	if this.destroyedCnt >= 0 {
		return
	}
	if this.appCnt > 120 {
		return
	}
	a := 1 - float32(this.appCnt)/120
	if this.startCnt < 12 {
		a = float32(this.startCnt) / 12
	}
	td := this.baseDeg + this.deg
	if this.spec.nway <= 1 {
		gl.Begin(gl.LINE_STRIP)
		setScreenColor(0.9, 0.1, 0.1, a)
		gl.Vertex2(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
		setScreenColor(0.9, 0.1, 0.1, a*0.5)
		gl.Vertex2(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
		gl.End()
	} else {
		td -= this.spec.nwayAngle * (this.spec.nway - 1) / 2
		gl.Begin(gl.LINE_STRIP)
		setScreenColor(0.9, 0.1, 0.1, a*0.75)
		gl.Vertex2(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
		setScreenColor(0.9, 0.1, 0.1, a*0.25)
		gl.Vertex2(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
		gl.End()
		gl.Begin(gl.QUADS)
		for i := 0; i < spec.nway-1; i++ {
			setScreenColor(0.9, 0.1, 0.1, a*0.3)
			gl.Vertex2(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
			setScreenColor(0.9, 0.1, 0.1, a*0.05)
			gl.Vertex2(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
			td += this.spec.nwayAngle
			gl.Vertex2(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
			setScreenColor(0.9, 0.1, 0.1, a*0.3)
			gl.Vertex2(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
		}
		gl.End()
		gl.Begin(gl.LINE_STRIP)
		setScreenColor(0.9, 0.1, 0.1, a*0.75)
		gl.Vertex2(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
		setScreenColor(0.9, 0.1, 0.1, a*0.25)
		gl.Vertex2(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
		gl.End()
	}
}

func (this *Turret) checkCollision(x float32, y float32, c Shape, shot Shot) bool {
	if this.destroyedCnt >= 0 || this.spec.invisible {
		return false
	}
	ox := fabs32(this.pos.x - x)
	oy := fabs32(this.pos.y - y)
	if this.spec.shape.checkCollision(ox, oy, c) {
		this.addDamage(shot.damage)
		return true
	}
	return false
}

func (this *Turret) addDamage(n int) {
	this.shield -= n
	if this.shield <= 0 {
		this.destroyed()
	}
	this.damaged = true
	this.damagedCnt = 10
}

func (this *Turret) destroyed() {
	playSe("turret_destroyed.wav")
	this.destroyedCnt = 0
	for i := 0; i < 6; i++ {
		NewSmoke(this.pos.x, this.pos.y, 0, nextSignedFloat(0.1), nextSignedFloat(0.1), nextFloat(0.04),
			Smoke.SmokeType.EXPLOSION, 30+nextInt(20), this.spec.size*1.5)
	}
	for i := 0; i < 32; i++ {
		NewSpark(this.pos, nextSignedFloat(0.5), nextSignedFloat(0.5),
			0.5+nextFloat(0.5), 0.5+nextFloat(0.5), 0, 30+nextInt(30))
	}
	for i := 0; i < 7; i++ {
		NewFragment(this.pos, nextSignedFloat(0.25), nextSignedFloat(0.25), 0.05+nextFloat(0.05),
			this.spec.size*(0.5+nextFloat(0.5)))
	}
	switch this.spec.enemyType {
	case TurretSpec.TurretType.MAIN:
		this.parent.increaseMultiplier(2)
		this.parent.addScore(40)
		break
	case TurretSpec.TurretType.SUB, TurretSpec.TurretType.SUB_DESTRUCTIVE:
		this.parent.increaseMultiplier(1)
		this.parent.addScore(20)
		break
	}
}

func (this *Turret) close() {
	if this.destroyedCnt < 0 {
		this.destroyedCnt = 999
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
	turretType                          TurretType
	interval                            int
	speed, speedAccel                   float32
	minRange, maxRange                  float32
	turnSpeed, turnRange                float32
	burstNum, burstInterval             int
	burstTurnRatio                      float32
	blind                               bool
	lookAheadRatio                      float32
	nway                                int
	nwayAngle                           float32
	nwayChange                          bool
	bulletShape                         int
	bulletDestructive                   bool
	shield                              int
	invisible                           bool
	shape, damagedShape, destroyedShape TurretShape
	size                                float32
}

func NewTurretSpec() *TurretSpec {
	this := new(TurretSpec)
	this.shape = NewTurretShape(TurretShapeTypeNORMAL)
	this.damagedShape = NewTurretShape(TurretShapeTypeDAMAGED)
	this.destroyedShape = NewTurretShape(TurretShapeTypeDESTROYED)
	this.turretType = TurretTypeMAIN
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

func (this *TurretSpec) setParamTurretSpec(ts TurretSpec) {
	this.turretType = ts.turretType
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

func (this *TurretSpec) setParam(rank float32, turretType TurretType) {
	this.turretType = turretType
	if turretType == TurretTypeDUMMY {
		this.invisible = true
		return
	}
	rk := this.rank
	switch this.turretType {
	case TurretType.SMALL:
		this.minRange = 8
		this.bulletShape = BulletShape.BulletShapeType.SMALL
		this.blind = true
		this.invisible = true
		break
	case TurretType.MOVING:
		this.minRange = 6
		this.bulletShape = BulletShape.BulletShapeType.MOVING_TURRET
		this.blind = true
		this.invisible = true
		this.turnSpeed = 0
		this.maxRange = 9 + nextFloat(12)
		rk *= (10.0 / sqrt(this.maxRange))
		break
	default:
		this.maxRange = 9 + nextFloat(16)
		this.minRange = this.maxRange / (4 + nextFloat(0.5))
		if this.turretType == TurretType.SUB || turretType == TurretType.SUB_DESTRUCTIVE {
			this.maxRange *= 0.72
			this.minRange *= 0.9
		}
		rk *= (10.0 / sqrt(this.maxRange))
		if nextInt(4) == 0 {
			lar := this.rank * 0.1
			if lar > 1 {
				lar = 1
			}
			this.lookAheadRatio = nextFloat(lar/2) + lar/2
			rk /= (1 + this.lookAheadRatio*0.3)
		}
		if nextInt(3) == 0 && this.lookAheadRatio == 0 {
			this.blind = false
			rk *= 1.5
		} else {
			this.blind = true
		}
		this.turnRange = Pi32/4 + nextFloat(Pi32/4)
		this.turnSpeed = 0.005 + nextFloat(0.015)
		if this.turretType == TurretType.MAIN {
			this.turnRange *= 1.2
		}
		if nextInt(4) == 0 {
			this.burstTurnRatio = nextFloat(0.66) + 0.33
		}
		break
	}
	this.burstInterval = 6 + nextInt(8)
	switch turretType {
	case TurretType.MAIN:
		this.size = 0.42 + nextFloat(0.05)
		br := (rk * 0.3) * (1 + nextSignedFloat(0.2))
		nr := (rk * 0.33) * nextFloat(1)
		ir := (rk * 0.1) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br) + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - this.burstNum + 1 - (this.nway-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt(sr * 0.6)
		this.speed *= 0.12
		this.shield = 20
		break
	case TurretType.SUB:
		this.size = 0.36 + nextFloat(0.025)
		br := (rk * 0.4) * (1 + nextSignedFloat(0.2))
		nr := (rk * 0.2) * nextFloat(1)
		ir := (rk * 0.2) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br) + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - this.burstNum + 1 - (this.nway-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt(sr * 0.7)
		this.speed *= 0.2
		this.shield = 12
		break
	case TurretType.SUB_DESTRUCTIVE:
		this.size = 0.36 + nextFloat(0.025)
		br := (rk * 0.4) * (1 + nextSignedFloat(0.2))
		nr := (rk * 0.2) * nextFloat(1)
		ir := (rk * 0.2) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br)*2 + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(60.0/(ir*2+1)) + 1
		this.burstInterval *= 0.88
		this.bulletShape = BulletShape.BulletShapeType.DESTRUCTIVE
		this.bulletDestructive = true
		sr := rk - (this.burstNum-1)/2 - (this.nway-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt(sr * 0.7)
		this.speed *= 0.33
		this.shield = 12
		break
	case TurretType.SMALL:
		this.size = 0.33
		br := (rk * 0.33) * (1 + nextSignedFloat(0.2))
		ir := (rk * 0.2) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br) + 1
		this.nway = 1
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - this.burstNum + 1 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt(sr)
		this.speed *= 0.24
		break
	case TurretType.MOVING:
		this.size = 0.36
		br := (rk * 0.3) * (1 + nextSignedFloat(0.2))
		nr := (rk * 0.1) * nextFloat(1)
		ir := (rk * 0.33) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br) + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - this.burstNum + 1 - (nway-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt(sr * 0.7)
		this.speed *= 0.2
		break
	}
	if this.speed < 0.1 {
		this.speed = 0.1
	} else {
		this.speed = sqrt(this.speed*10) / 10
	}
	if this.burstNum > 2 {
		if nextInt(4) == 0 {
			this.speed *= 0.8
			this.burstInterval *= 0.7
			this.speedAccel = (this.speed * (0.4 + nextFloat(0.3))) / this.burstNum
			if nextInt(2) == 0 {
				this.speedAccel *= -1
			}
			this.speed -= this.speedAccel * this.burstNum / 2
		}
		if nextInt(5) == 0 {
			if this.nway > 1 {
				this.nwayChange = true
			}
		}
	}
	this.nwayAngle = (0.1 + nextFloat(0.33)) / (1 + this.nway*0.1)
}

func (this *TurretSpec) setBossSpec() {
	this.minRange = 0
	this.maxRange *= 1.5
	this.shield *= 2.1
}

func (this *TurretSpec) sizes(v float32) float32 {
	this.size = v
	this.destroyedShape.size = v
	this.damagedShape.size = v
	this.shape.size = v
	return v
}

/**
 * Grouped turrets.
 */
const TURRET_GROUP_MAX_NUM = 16

type TurretGroup struct {
	ship      Ship
	spec      TurretGroupSpec
	centerPos Vector
	turret    [TURRET_GROUP_MAX_NUM]*Turret
	cnt       int
}

func NewTurretGroup(field Field, ship Ship, parent Enemy, spec TurretGroupSpec) *TurretGroup {
	this := new(TurretGroup)
	this.ship = ship
	for i, _ := range this.turret {
		this.turret[i] = NewTurret(field, bullets, ship, sparks, smokes, fragments, parent)
	}
	this.spec = spec
	return this
}

func (this *TurretGroup) move(p Vector, deg float32) bool {
	alive := false
	this.centerPos.x = p.x
	this.centerPos.y = p.y
	var d, md, y, my float32
	switch this.spec.alignType {
	case TurretGroupSpec.AlignType.ROUND:
		d = this.spec.alignDeg
		if this.spec.num > 1 {
			md = this.spec.alignWidth / (this.spec.num - 1)
			d -= this.spec.alignWidth / 2
		} else {
			md = 0
		}
		break
	case TurretGroupSpec.AlignType.STRAIGHT:
		y = 0
		my = this.spec.offset.y / (this.spec.num + 1)
		break
	}
	for i := 0; i < this.spec.num; i++ {
		var tbx, tby float32
		switch this.spec.alignType {
		case TurretGroupSpec.AlignType.ROUND:
			tbx = Sin32(d) * this.spec.radius
			tby = Cos32(d) * this.spec.radius
			break
		case TurretGroupSpec.AlignType.STRAIGHT:
			y += my
			tbx = this.spec.offset.x
			tby = y
			d = atan2(tbx, tby)
			break
		}
		tbx *= (1 - this.spec.distRatio)
		bx := tbx*Cos32(-deg) - tby*Sin32(-deg)
		by := tbx*Sin32(-deg) + tby*Cos32(-deg)
		alive |= this.turret[i].move(this.centerPos.x+bx, this.centerPos.y+by, d+deg)
		if this.spec.alignType == TurretGroupSpec.AlignType.ROUND {
			d += md
		}
	}
	this.cnt++
	return alive
}

func (this *TurretGroup) draw() {
	for i := 0; i < this.spec.num; i++ {
		this.turret[i].draw()
	}
}

func (this *TurretGroup) close() {
	for i := 0; i < this.spec.num; i++ {
		this.turret[i].close()
	}
}

func (this *TurretGroup) checkCollision(x float32, y float32, c Shape, shot Shot) bool {
	col := false
	for i := 0; i < this.spec.num; i++ {
		col |= this.turret[i].checkCollision(x, y, c, shot)
	}
	return col
}

type AlignType int

const (
	ROUND AlignType = iota
	STRAIGHT
)

type TurretGroupSpec struct {
	turretSpec                              TurretSpec
	num, alignType                          int
	alignDeg, alignWidth, radius, distRatio float32
	offset                                  Vector
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
	ship                                                                                                           Ship
	spec                                                                                                           MovingTurretGroupSpec
	radius, radiusAmpCnt, deg, rollAmpCnt, swingAmpCnt, swingAmpDeg, swingFixDeg, alignAmpCnt, distDeg, distAmpCnt float32
	cnt                                                                                                            int
	centerPos                                                                                                      Vector
	turret                                                                                                         [MOVING_TURRET_MAX_NUM]Turret
}

func NewMovingTurretGroup(field Field, ship Ship, parent Enemy, spec MovingTurretGroupSpec) *MovingTurretGroup {
	this := new(MovingTurretGroup)
	this.ship = ship
	for i, _ := range this.turret {
		this.turret[i] = NewTurret(field, ship, parent)
	}
	this.spec = spec
	this.radius = spec.radiusBase
	this.swingFixDeg = Pi32
	return this
}

func (this *MovingTurretGroup) move(p Vector, od float32) {
	if this.spec.moveType == MovingTurretGroupSpec.MoveType.SWING_FIX {
		this.swingFixDeg = ed
	}
	this.centerPos.x = p.x
	this.centerPos.y = p.y
	if this.spec.radiusAmp > 0 {
		this.radiusAmpCnt += this.spec.radiusAmpVel
		av := Sin32(this.radiusAmpCnt)
		this.radius = this.spec.radiusBase + this.spec.radiusAmp*av
	}
	if this.spec.moveType == MovingTurretGroupSpec.MoveType.ROLL {
		if this.spec.rollAmp != 0 {
			this.rollAmpCnt += this.spec.rollAmpVel
			av := Sin32(this.rollAmpCnt)
			this.deg += this.spec.rollDegVel + this.spec.rollAmp*av
		} else {
			this.deg += this.spec.rollDegVel
		}
	} else {
		this.swingAmpCnt += this.spec.swingAmpVel
		if Cos32(this.swingAmpCnt) > 0 {
			this.swingAmpDeg += this.spec.swingDegVel
		} else {
			this.swingAmpDeg -= this.spec.swingDegVel
		}
		if this.spec.moveType == MovingTurretGroupSpec.MoveType.SWING_AIM {
			var od float32
			shipPos := ship.nearPos(this.centerPos)
			if shipPos.dist(this.centerPos) < 0.1 {
				od = 0
			} else {
				od = atan2(shipPos.x-this.centerPos.x, shipPos.y-this.centerPos.y)
			}
			od += this.swingAmpDeg - this.deg
			normalizeDeg(od)
			deg += od * 0.1
		} else {
			od := this.swingFixDeg + this.swingAmpDeg - this.deg
			normalizeDeg(od)
			deg += od * 0.1
		}
	}
	var d, ad, md float32
	calcAlignDeg(d, ad, md)
	for i := 0; i < this.spec.num; i++ {
		d += md
		bx := Sin32(d) * this.radius * this.spec.xReverse
		by := Cos32(d) * this.radius * (1 - this.spec.distRatio)
		var fs, fd float32
		if fabs32(bx)+fabs32(by) < 0.1 {
			fs = this.radius
			fd = d
		} else {
			fs = sqrt(bx*bx + by*by)
			fd = atan2(bx, by)
		}
		fs *= 0.06
		this.turret[i].move(this.centerPos.x, this.centerPos.y, d, fs, fd)
	}
	cnt++
}

func (this *MovingTurretGroup) calcAlignDeg(d *float32, ad *float32, md *float32) {
	this.alignAmpCnt += this.spec.alignAmpVel
	ad = this.spec.alignDeg * (1 + Sin32(this.alignAmpCnt)*this.spec.alignAmp)
	if this.spec.num > 1 {
		if this.spec.moveType == MovingTurretGroupSpec.MoveType.ROLL {
			md = ad / this.spec.num
		} else {
			md = ad / (this.spec.num - 1)
		}
	} else {
		md = 0
	}
	d = this.deg - md - ad/2
}

func (this *MovingTurretGroup) draw() {
	for i := 0; i < this.spec.num; i++ {
		this.turret[i].draw()
	}
}

func (this *MovingTurretGroup) close() {
	for i := 0; i < this.spec.num; i++ {
		this.turret[i].close()
	}
}

type TurretMoveType int

const (
	TurretMoveTypeROLL MoveType = iota
	TurretMoveTypeSWING_FIX
	TurretMoveTypeSWING_AIM
)

type MovingTurretGroupSpec struct {
	turretSpec                                                                                                                                           TurretSpec
	num                                                                                                                                                  int
	moveType                                                                                                                                             TurretMoveType
	alignDeg, alignAmp, alignAmpVel, radiusBase, radiusAmp, radiusAmpVel, rollDegVel, rollAmp, rollAmpVel, swingDegVel, swingAmpVel, distRatio, xReverse float32
}

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
	this.alignAmp = a
	this.alignAmpVel = v
}

func (this *MovingTurretGroupSpec) setRadiusAmp(a float32, v float32) {
	this.radiusAmp = a
	this.radiusAmpVel = v
}

func (this *MovingTurretGroupSpec) setRoll(dv float32, a float32, v float32) {
	this.moveType = MoveType.ROLL
	this.rollDegVel = dv
	this.rollAmp = a
	this.rollAmpVel = v
}

func (this *MovingTurretGroupSpec) setSwing(dv float32, a float32, aim bool /*= false*/) {
	if aim {
		this.moveType = MoveType.SWING_AIM
	} else {
		this.moveType = MoveType.SWING_FIX
	}
	this.swingDegVel = dv
	this.swingAmpVel = a
}

func (this *MovingTurretGroupSpec) setXReverse(xr float32) {
	this.xReverse = xr
}
