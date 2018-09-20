/*
 * $Id: turret.d,v 1.3 2005/07/17 11:02:46 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/go-gl/gl/v2.1/gl"

/**
 * Turret mounted on a deck of an enemy ship.
 */

var turretDamagedPos Vector

type Turret struct {
	spec                          *TurretSpec
	pos                           Vector
	deg, baseDeg                  float32
	cnt, appCnt, startCnt, shield int
	damaged                       bool
	destroyedCnt, damagedCnt      int
	bulletSpeed                   float32
	burstCnt                      int
	isBoss                        bool
	enemyIndex                    int
	multiplier                    *float32
	addScore                      func(int)
}

func NewTurret(spec *TurretSpec, isBoss bool, enemyIndex int, multiplier *float32, addScore func(int)) *Turret {
	if spec.shape == nil {
		panic("turret spec shape nil")
	}
	this := new(Turret)
	this.bulletSpeed = 1
	this.spec = spec
	this.shield = spec.shield
	this.destroyedCnt = -1
	this.bulletSpeed = 1
	this.enemyIndex = enemyIndex
	this.multiplier = multiplier
	this.addScore = addScore
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
			NewSmoke(this.pos.x, this.pos.y, 0, 0, 0, 0.01+nextFloat(0.01), SmokeTypeFIRE, 90+nextInt(30), this.spec.size)
		}
		return false
	}
	td := this.baseDeg + this.deg
	shipPos := ship.nearPos(this.pos)
	shipVel := ship.nearVel(this.pos)
	ax := shipPos.x - this.pos.x
	ay := shipPos.y - this.pos.y
	if this.spec.lookAheadRatio != 0 {
		rd := this.pos.distVector(shipPos) / this.spec.speed * 1.2
		ax += shipVel.x * this.spec.lookAheadRatio * rd
		ay += shipVel.y * this.spec.lookAheadRatio * rd
	}
	var ad float32
	if fabs32(ax)+fabs32(ay) < 0.1 {
		ad = 0
	} else {
		ad = atan232(ax, ay)
	}
	od := td - ad
	od = normalizeDeg(od)
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
	this.deg = normalizeDeg(this.deg)
	if this.deg > this.spec.turnRange {
		this.deg = this.spec.turnRange
	} else if this.deg < -this.spec.turnRange {
		this.deg = -this.spec.turnRange
	}
	this.cnt++
	if field.checkInFieldVector(this.pos) || (this.isBoss && this.cnt%4 == 0) {
		this.appCnt++
	}
	if this.cnt >= this.spec.interval {
		if this.spec.blind || (fabs32(od) <= this.spec.turnSpeed &&
			this.pos.distVector(shipPos) < this.spec.maxRange*1.1 &&
			this.pos.distVector(shipPos) > this.spec.minRange) {
			this.cnt = -(this.spec.burstNum - 1) * this.spec.burstInterval
			this.bulletSpeed = this.spec.speed
			this.burstCnt = 0
		}
	}
	if this.cnt <= 0 && -this.cnt%this.spec.burstInterval == 0 &&
		((this.spec.invisible && field.checkInFieldVector(this.pos)) ||
			(this.spec.invisible && this.isBoss && field.checkInOuterFieldVector(this.pos)) ||
			(!this.spec.invisible && field.checkInFieldExceptTop(this.pos))) &&
		this.pos.distVector(shipPos) > this.spec.minRange {
		bd := this.baseDeg + this.deg
		NewSmoke(this.pos.x, this.pos.y, 0, Sin32(bd)*this.bulletSpeed, Cos32(bd)*this.bulletSpeed, 0,
			SmokeTypeSPARK, 20, this.spec.size*2)
		nw := this.spec.nway
		if this.spec.nwayChange && this.burstCnt%2 == 1 {
			nw--
		}
		bd -= this.spec.nwayAngle * (float32(nw) - 1) / 2
		for i := 0; i < nw; i++ {
			NewBullet(this.enemyIndex,
				this.pos, bd, this.bulletSpeed, this.spec.size*3, this.spec.bulletShape, this.spec.maxRange,
				bulletFireSpeed, bulletFireDeg, this.spec.bulletDestructive)
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
		turretDamagedPos.x = this.pos.x + nextSignedFloat(float32(this.damagedCnt)*0.015)
		turretDamagedPos.y = this.pos.y + nextSignedFloat(float32(this.damagedCnt)*0.015)
		glTranslate(turretDamagedPos)
	} else {
		glTranslate(this.pos)
	}
	gl.Rotatef(-(this.baseDeg+this.deg)*180/Pi32, 0, 0, 1)
	if this.destroyedCnt >= 0 {
		this.spec.destroyedShape.draw()
	} else if !this.damaged {
		if this.spec.shape == nil {
			panic("turret spec shape nil")
		}
		this.spec.shape.draw() // this is the bad turret
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
		gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
		setScreenColor(0.9, 0.1, 0.1, a*0.5)
		gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
		gl.End()
	} else {
		td -= this.spec.nwayAngle * (float32(this.spec.nway) - 1) / 2
		gl.Begin(gl.LINE_STRIP)
		setScreenColor(0.9, 0.1, 0.1, a*0.75)
		gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
		setScreenColor(0.9, 0.1, 0.1, a*0.25)
		gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
		gl.End()
		gl.Begin(gl.QUADS)
		for i := 0; i < this.spec.nway-1; i++ {
			setScreenColor(0.9, 0.1, 0.1, a*0.3)
			gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
			setScreenColor(0.9, 0.1, 0.1, a*0.05)
			gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
			td += this.spec.nwayAngle
			gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
			setScreenColor(0.9, 0.1, 0.1, a*0.3)
			gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
		}
		gl.End()
		gl.Begin(gl.LINE_STRIP)
		setScreenColor(0.9, 0.1, 0.1, a*0.75)
		gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.minRange, this.pos.y+Cos32(td)*this.spec.minRange)
		setScreenColor(0.9, 0.1, 0.1, a*0.25)
		gl.Vertex2f(this.pos.x+Sin32(td)*this.spec.maxRange, this.pos.y+Cos32(td)*this.spec.maxRange)
		gl.End()
	}
}

func (this *Turret) checkCollision(x float32, y float32, c Shape, shot *Shot) bool {
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
			SmokeTypeEXPLOSION, 30+nextInt(20), this.spec.size*1.5)
	}
	for i := 0; i < 32; i++ {
		NewSpark(this.pos, nextSignedFloat(0.5), nextSignedFloat(0.5),
			0.5+nextFloat(0.5), 0.5+nextFloat(0.5), 0, 30+nextInt(30))
	}
	for i := 0; i < 7; i++ {
		NewFragment(this.pos, nextSignedFloat(0.25), nextSignedFloat(0.25), 0.05+nextFloat(0.05),
			this.spec.size*(0.5+nextFloat(0.5)))
	}
	switch this.spec.turretType {
	case TurretTypeMAIN:
		*(this.multiplier) += 2
		this.addScore(40)
	case TurretTypeSUB, TurretTypeSUB_DESTRUCTIVE:
		*(this.multiplier) += 1
		this.addScore(20)
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
	bulletShape                         BulletShapeType
	bulletDestructive                   bool
	shield                              int
	invisible                           bool
	shape, damagedShape, destroyedShape *TurretShape
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
	this.bulletShape = BulletShapeTypeNORMAL
	this.shield = 99999
	this.sizes(1)
	return this
}

func (this *TurretSpec) setParamTurretSpec(ts *TurretSpec) {
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
	this.sizes(ts.size)
}

func (this *TurretSpec) setParam(rank float32, turretType TurretType) {
	this.turretType = turretType
	if turretType == TurretTypeDUMMY {
		this.invisible = true
		return
	}
	rk := rank
	switch this.turretType {
	case TurretTypeSMALL:
		this.minRange = 8
		this.bulletShape = BulletShapeTypeSMALL
		this.blind = true
		this.invisible = true
	case TurretTypeMOVING:
		this.minRange = 6
		this.bulletShape = BulletShapeTypeMOVING_TURRET
		this.blind = true
		this.invisible = true
		this.turnSpeed = 0
		this.maxRange = 9 + nextFloat(12)
		rk *= (10.0 / sqrt32(this.maxRange))
	default:
		this.maxRange = 9 + nextFloat(16)
		this.minRange = this.maxRange / (4 + nextFloat(0.5))
		if this.turretType == TurretTypeSUB || turretType == TurretTypeSUB_DESTRUCTIVE {
			this.maxRange *= 0.72
			this.minRange *= 0.9
		}
		rk *= (10.0 / sqrt32(this.maxRange))
		if nextInt(4) == 0 {
			lar := rank * 0.1
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
		if this.turretType == TurretTypeMAIN {
			this.turnRange *= 1.2
		}
		if nextInt(4) == 0 {
			this.burstTurnRatio = nextFloat(0.66) + 0.33
		}
	}
	this.burstInterval = 6 + nextInt(8)
	switch turretType {
	case TurretTypeMAIN:
		this.sizes(0.42 + nextFloat(0.05))
		br := (rk * 0.3) * (1 + nextSignedFloat(0.2))
		nr := (rk * 0.33) * nextFloat(1)
		ir := (rk * 0.1) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br) + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - float32(this.burstNum) + 1 - (float32(this.nway)-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt32(sr * 0.6)
		this.speed *= 0.12
		this.shield = 20
	case TurretTypeSUB:
		this.sizes(0.36 + nextFloat(0.025))
		br := (rk * 0.4) * (1 + nextSignedFloat(0.2))
		nr := (rk * 0.2) * nextFloat(1)
		ir := (rk * 0.2) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br) + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - float32(this.burstNum) + 1 - (float32(this.nway)-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt32(sr * 0.7)
		this.speed *= 0.2
		this.shield = 12
	case TurretTypeSUB_DESTRUCTIVE:
		this.sizes(0.36 + nextFloat(0.025))
		br := (rk * 0.4) * (1 + nextSignedFloat(0.2))
		nr := (rk * 0.2) * nextFloat(1)
		ir := (rk * 0.2) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br)*2 + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(60.0/(ir*2+1)) + 1
		this.burstInterval = int(float32(this.burstInterval) * 0.88)
		this.bulletShape = BulletShapeTypeDESTRUCTIVE
		this.bulletDestructive = true
		sr := rk - (float32(this.burstNum)-1)/2 - (float32(this.nway)-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt32(sr * 0.7)
		this.speed *= 0.33
		this.shield = 12
	case TurretTypeSMALL:
		this.sizes(0.33)
		br := (rk * 0.33) * (1 + nextSignedFloat(0.2))
		ir := (rk * 0.2) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br) + 1
		this.nway = 1
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - float32(this.burstNum) + 1 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt32(sr)
		this.speed *= 0.24
	case TurretTypeMOVING:
		this.sizes(0.36)
		br := (rk * 0.3) * (1 + nextSignedFloat(0.2))
		nr := (rk * 0.1) * nextFloat(1)
		ir := (rk * 0.33) * (1 + nextSignedFloat(0.2))
		this.burstNum = int(br) + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - float32(this.burstNum) + 1 - (float32(this.nway)-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = sqrt32(sr * 0.7)
		this.speed *= 0.2
	}
	if this.speed < 0.1 {
		this.speed = 0.1
	} else {
		this.speed = sqrt32(this.speed*10) / 10
	}
	if this.burstNum > 2 {
		if nextInt(4) == 0 {
			this.speed *= 0.8
			this.burstInterval = int(float32(this.burstInterval) * 0.7)
			this.speedAccel = (this.speed * (0.4 + nextFloat(0.3))) / float32(this.burstNum)
			if nextInt(2) == 0 {
				this.speedAccel *= -1
			}
			this.speed -= this.speedAccel * float32(this.burstNum) / 2
		}
		if nextInt(5) == 0 {
			if this.nway > 1 {
				this.nwayChange = true
			}
		}
	}
	this.nwayAngle = (0.1 + nextFloat(0.33)) / (1 + float32(this.nway)*0.1)
}

func (this *TurretSpec) setBossSpec() {
	this.minRange = 0
	this.maxRange *= 1.5
	this.shield = int(float32(this.shield) * 2.1)
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
	spec      *TurretGroupSpec
	centerPos Vector
	turret    []*Turret
	cnt       int
	parent    *Enemy
}

func NewTurretGroup(parent *Enemy, spec *TurretGroupSpec) *TurretGroup {
	this := new(TurretGroup)
	this.parent = parent
	if spec.turretSpec.shape == nil {
		panic("turret spec shape nil")
	}
	this.turret = make([]*Turret, 0, TURRET_GROUP_MAX_NUM)
	this.spec = spec
	this.addTurrets()
	return this
}

func (this *TurretGroup) addTurret() {
	this.turret = append(this.turret, NewTurret(this.spec.turretSpec, this.parent.isBoss(), this.parent.index(), &(this.parent.multiplier), this.parent.addScore))
}

func (this *TurretGroup) addTurrets() {
	for len(this.turret) < this.spec.num {
		this.addTurret()
	}
}

func (this *TurretGroup) move(p Vector, deg float32) bool {
	alive := false
	this.centerPos.x = p.x
	this.centerPos.y = p.y
	var d, md, y, my float32
	switch this.spec.alignType {
	case AlignTypeROUND:
		d = this.spec.alignDeg
		if this.spec.num > 1 {
			md = this.spec.alignWidth / (float32(this.spec.num) - 1)
			d -= this.spec.alignWidth / 2
		} else {
			md = 0
		}
	case AlignTypeSTRAIGHT:
		y = 0
		my = this.spec.offset.y / (float32(this.spec.num) + 1)
	}
	for _, t := range this.turret {
		var tbx, tby float32
		switch this.spec.alignType {
		case AlignTypeROUND:
			tbx = Sin32(d) * this.spec.radius
			tby = Cos32(d) * this.spec.radius
		case AlignTypeSTRAIGHT:
			y += my
			tbx = this.spec.offset.x
			tby = y
			d = atan232(tbx, tby)
		}
		tbx *= (1 - this.spec.distRatio)
		bx := tbx*Cos32(-deg) - tby*Sin32(-deg)
		by := tbx*Sin32(-deg) + tby*Cos32(-deg)
		a := t.move(this.centerPos.x+bx, this.centerPos.y+by, d+deg, 0, -99999)
		alive = alive || a
		if this.spec.alignType == AlignTypeROUND {
			d += md
		}
	}
	this.cnt++
	return alive
}

func (this *TurretGroup) draw() {
	for _, t := range this.turret {
		t.draw()
	}
}

func (this *TurretGroup) close() {
	for _, t := range this.turret {
		t.close()
	}
}

func (this *TurretGroup) checkCollision(x float32, y float32, c Shape, shot *Shot) bool {
	col := false
	for _, t := range this.turret {
		col = col || t.checkCollision(x, y, c, shot)
	}
	return col
}

type AlignType int

const (
	AlignTypeROUND AlignType = iota
	AlignTypeSTRAIGHT
)

type TurretGroupSpec struct {
	turretSpec                              *TurretSpec
	num                                     int
	alignType                               AlignType
	alignDeg, alignWidth, radius, distRatio float32
	offset                                  Vector
}

func NewTurretGroupSpec() *TurretGroupSpec {
	this := new(TurretGroupSpec)
	this.turretSpec = NewTurretSpec()
	this.num = 1
	this.alignType = AlignTypeROUND
	return this
}

/**
 * Turrets moving around a bridge.
 */

const MOVING_TURRET_MAX_NUM = 16

type MovingTurretGroup struct {
	spec                                  *MovingTurretGroupSpec
	radius, radiusAmpCnt, deg, rollAmpCnt float32
	swingAmpCnt, swingAmpDeg, swingFixDeg float32
	alignAmpCnt, distDeg, distAmpCnt      float32
	cnt                                   int
	centerPos                             Vector
	turret                                []*Turret
	parent                                *Enemy
}

func NewMovingTurretGroup(parent *Enemy, spec *MovingTurretGroupSpec) *MovingTurretGroup {
	this := new(MovingTurretGroup)
	this.turret = make([]*Turret, 0, MOVING_TURRET_MAX_NUM)
	this.spec = spec
	this.radius = spec.radiusBase
	this.swingFixDeg = Pi32
	this.parent = parent
	this.addTurrets()
	return this
}

func (this *MovingTurretGroup) addTurret() {
	this.turret = append(this.turret, NewTurret(this.spec.turretSpec, this.parent.isBoss(), this.parent.index(), &(this.parent.multiplier), this.parent.addScore))
}

func (this *MovingTurretGroup) addTurrets() {
	for len(this.turret) < this.spec.num {
		this.addTurret()
	}
}

func (this *MovingTurretGroup) move(p Vector, od float32) {
	if this.spec.moveType == TurretMoveTypeSWING_FIX {
		this.swingFixDeg = od
	}
	this.centerPos.x = p.x
	this.centerPos.y = p.y
	if this.spec.radiusAmp > 0 {
		this.radiusAmpCnt += this.spec.radiusAmpVel
		av := Sin32(this.radiusAmpCnt)
		this.radius = this.spec.radiusBase + this.spec.radiusAmp*av
	}
	if this.spec.moveType == TurretMoveTypeROLL {
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
		if this.spec.moveType == TurretMoveTypeSWING_AIM {
			var od float32
			shipPos := ship.nearPos(this.centerPos)
			if shipPos.distVector(this.centerPos) < 0.1 {
				od = 0
			} else {
				od = atan232(shipPos.x-this.centerPos.x, shipPos.y-this.centerPos.y)
			}
			od += this.swingAmpDeg - this.deg
			od = normalizeDeg(od)
			this.deg += od * 0.1
		} else {
			od := this.swingFixDeg + this.swingAmpDeg - this.deg
			od = normalizeDeg(od)
			this.deg += od * 0.1
		}
	}
	var d, ad, md float32
	this.calcAlignDeg(&d, &ad, &md)
	for _, t := range this.turret {
		d += md
		bx := Sin32(d) * this.radius * this.spec.xReverse
		by := Cos32(d) * this.radius * (1 - this.spec.distRatio)
		var fs, fd float32
		if fabs32(bx)+fabs32(by) < 0.1 {
			fs = this.radius
			fd = d
		} else {
			fs = sqrt32(bx*bx + by*by)
			fd = atan232(bx, by)
		}
		fs *= 0.06
		t.move(this.centerPos.x, this.centerPos.y, d, fs, fd)
	}
	this.cnt++
}

func (this *MovingTurretGroup) calcAlignDeg(d *float32, ad *float32, md *float32) {
	this.alignAmpCnt += this.spec.alignAmpVel
	*ad = this.spec.alignDeg * (1 + Sin32(this.alignAmpCnt)*this.spec.alignAmp)
	if this.spec.num > 1 {
		if this.spec.moveType == TurretMoveTypeROLL {
			*md = *ad / float32(this.spec.num)
		} else {
			*md = *ad / (float32(this.spec.num) - 1)
		}
	} else {
		*md = 0
	}
	*d = this.deg - *md - *ad/2
}

func (this *MovingTurretGroup) draw() {
	for _, t := range this.turret {
		t.draw()
	}
}

func (this *MovingTurretGroup) close() {
	for _, t := range this.turret {
		t.close()
	}
}

type TurretMoveType int

const (
	TurretMoveTypeROLL TurretMoveType = iota
	TurretMoveTypeSWING_FIX
	TurretMoveTypeSWING_AIM
)

type MovingTurretGroupSpec struct {
	turretSpec                                   *TurretSpec
	num                                          int
	moveType                                     TurretMoveType
	alignDeg, alignAmp, alignAmpVel, radiusBase  float32
	radiusAmp, radiusAmpVel, rollDegVel, rollAmp float32
	rollAmpVel, swingDegVel, swingAmpVel         float32
	distRatio, xReverse                          float32
}

func NewMovingTurretGroupSpec() *MovingTurretGroupSpec {
	this := new(MovingTurretGroupSpec)
	this.turretSpec = NewTurretSpec()
	this.num = 1
	this.alignDeg = Pi32 * 2
	this.radiusBase = 1
	this.moveType = TurretMoveTypeSWING_FIX
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
	this.moveType = TurretMoveTypeROLL
	this.rollDegVel = dv
	this.rollAmp = a
	this.rollAmpVel = v
}

func (this *MovingTurretGroupSpec) setSwing(dv float32, a float32, aim bool /*= false*/) {
	if aim {
		this.moveType = TurretMoveTypeSWING_AIM
	} else {
		this.moveType = TurretMoveTypeSWING_FIX
	}
	this.swingDegVel = dv
	this.swingAmpVel = a
}

func (this *MovingTurretGroupSpec) setXReverse(xr float32) {
	this.xReverse = xr
}
