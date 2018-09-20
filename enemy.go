/*
 * $Id: enemy.d,v 1.2 2005/07/17 11:02:45 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/go-gl/gl/v2.1/gl"

/**
 * Enemy ships.
 */
type Enemy struct {
	spec EnemySpec

	appType                                               AppearanceType
	ppos, pos                                             Vector
	shield                                                int
	deg, velDeg, speed, turnWay, trgDeg                   float32
	turnCnt, cnt                                          int
	vel                                                   Vector
	state                                                 MoveState
	turretGroup                                           []*TurretGroup
	movingTurretGroup                                     []*MovingTurretGroup
	damaged                                               bool
	damagedCnt, destroyedCnt, explodeCnt, explodeItv, idx int
	multiplier                                            float32
}

func NewEnemy(spec EnemySpec) *Enemy {
	this := new(Enemy)
	this.spec = spec
	this.idx = idxCount
	idxCount++
	this.turnWay = 1
	this.explodeItv = 1
	this.multiplier = 1

	this.turretGroup = make([]*TurretGroup, 0, TURRET_GROUP_MAX)
	this.movingTurretGroup = make([]*MovingTurretGroup, 0, MOVING_TURRET_GROUP_MAX)

	this.spec = spec
	this.shield = spec.shield()
	this.damaged = false
	this.destroyedCnt = -1
	this.explodeItv = 1
	this.multiplier = 1
	actors[this] = true
	return this
}

func (this *Enemy) addTurretGroup(i int) {
	if len(this.turretGroup) <= i {
		for j := len(this.turretGroup); j <= i; j++ {
			tg := NewTurretGroup(this, this.spec.turretGroupSpec()[j])
			this.turretGroup = append(this.turretGroup, tg)
		}
	}
}

func (this *Enemy) addMovingTurretGroup(i int) {
	if len(this.movingTurretGroup) <= i {
		for j := len(this.movingTurretGroup); j <= i; j++ {
			tg := NewMovingTurretGroup(this, this.spec.movingTurretGroupSpec()[j])
			this.movingTurretGroup = append(this.movingTurretGroup, tg)
		}
	}
}

func (this *Enemy) index() int {
	return this.idx
}

func (this *Enemy) move() {
	if !this.spec.move(this) {
		this.close()
	}
}

func (this *Enemy) size() float32 {
	return this.spec.size()
}

func (this *Enemy) checkShotHit(p Vector, shape Shape, shot *Shot) {
	if this.destroyedCnt >= 0 {
		return
	}
	if this.spec.checkCollision(this, p.x, p.y, shape, shot) {
		if shot != nil {
			shot.removeHitToEnemy(this.spec.isSmallEnemy())
		}
	}
}

func (this *Enemy) checkHitShip(x float32, y float32, largeOnly bool /*= false*/) bool {
	return this.spec.checkShipCollision(this, x, y, largeOnly)
}

func (this *Enemy) close() {
	this.removeTurrets()
	delete(actors, this)
}

func (this *Enemy) ndex() int {
	return this.idx
}

func (this *Enemy) isBoss() bool {
	return this.spec.isBoss()
}

/**
 * Enemy status (position, direction, velocity, turrets, etc).
 */
type AppearanceType int

const (
	AppearanceTypeTOP AppearanceType = iota
	AppearanceTypeSIDE
	AppearanceTypeCENTER
)

const TURRET_GROUP_MAX = 10
const MOVING_TURRET_GROUP_MAX = 4
const MULTIPLIER_DECREASE_RATIO = 0.005

var edgePos, explodeVel, damagedPos Vector
var idxCount int = 0

func (this *Enemy) setAppearancePos(appType AppearanceType /*= AppearanceTypeTOP*/) bool {
	this.appType = appType
	for i := 0; i < 8; i++ {
		switch appType {
		case AppearanceTypeTOP:
			this.pos.x = nextSignedFloat(field.size.x)
			this.pos.y = field.outerSize.y*0.99 + this.spec.size()
			if this.pos.x < 0 {
				this.deg = Pi32 - nextFloat(0.5)
				this.velDeg = this.deg
			} else {
				this.deg = Pi32 + nextFloat(0.5)
				this.velDeg = this.deg
			}
		case AppearanceTypeSIDE:
			if nextInt(2) == 0 {
				this.pos.x = -field.outerSize.x * 0.99
				this.deg = Pi32/2 + nextFloat(0.66)
				this.velDeg = this.deg
			} else {
				this.pos.x = field.outerSize.x * 0.99
				this.deg = -Pi32/2 - nextFloat(0.66)
				this.velDeg = this.deg
			}
			this.pos.y = field.size.y + nextFloat(field.size.y) + this.spec.size()
		case AppearanceTypeCENTER:
			this.pos.x = 0
			this.pos.y = field.outerSize.y*0.99 + this.spec.size()
			this.deg = 0
			this.velDeg = this.deg
		}
		this.ppos.x = this.pos.x
		this.ppos.y = this.pos.y
		this.vel.y = 0
		this.vel.x = 0
		this.speed = 0
		if this.appType == AppearanceTypeCENTER || this.checkFrontClear(true) {
			return true
		}
	}
	return false
}

func (this *Enemy) checkFrontClear(checkCurrentPos bool /*= false*/) bool {
	var si = 1
	if checkCurrentPos {
		si = 0
	}
	for i := si; i < 5; i++ {
		cx := this.pos.x + Sin32(this.deg)*float32(i)*this.spec.size()
		cy := this.pos.y + Cos32(this.deg)*float32(i)*this.spec.size()
		if field.getBlock(cx, cy) >= 0 {
			return false
		}
		if checkAllEnemiesHitShip(cx, cy, this, true) != nil {
			return false
		}
	}
	return true
}

func (this *Enemy) stateMove() bool {
	this.ppos.x = this.pos.x
	this.ppos.y = this.pos.y
	this.multiplier -= MULTIPLIER_DECREASE_RATIO
	if this.multiplier < 1 {
		this.multiplier = 1
	}
	if this.destroyedCnt >= 0 {
		this.destroyedCnt++
		this.explodeCnt--
		if this.explodeCnt < 0 {
			this.explodeItv += 2
			this.explodeItv = int(float32(this.explodeItv) * (1.2 + nextFloat(1)))
			this.explodeCnt = this.explodeItv
			this.destroyedEdge(int(sqrt32(this.spec.size()) * 27.0 / (float32(this.explodeItv)*0.1 + 1)))
		}
	}
	this.damaged = false
	if this.damagedCnt > 0 {
		this.damagedCnt--
	}
	alive := false
	this.addTurretGroup(this.spec.turretGroupNum() - 1)
	for _, tg := range this.turretGroup {
		a := tg.move(this.pos, this.deg)
		alive = alive || a
	}
	this.addMovingTurretGroup(this.spec.movingTurretGroupNum() - 1)
	for _, mtg := range this.movingTurretGroup {
		mtg.move(this.pos, this.deg)
	}
	if this.destroyedCnt < 0 && !alive {
		return this.destroyed(nil)
	}
	return true
}

func (this *Enemy) checkCollision(x float32, y float32, c Shape, shot *Shot) bool {
	ox := fabs32(this.pos.x - x)
	oy := fabs32(this.pos.y - y)
	if ox+oy > this.spec.size()*2 {
		return false
	}
	for _, t := range this.turretGroup {
		if t.checkCollision(x, y, c, shot) {
			return true
		}
	}
	if this.spec.bridgeShape().checkCollision(ox, oy, c) {
		this.addDamage(shot.damage, shot)
		return true
	}
	return false
}

func (this *Enemy) addScore(s uint32) {
	this.setScoreIndicator(s, 1)
}

func (this *Enemy) addDamage(n int, shot *Shot /*= null*/) {
	this.shield -= n
	if this.shield <= 0 {
		this.destroyed(shot)
	} else {
		this.damaged = true
		this.damagedCnt = 7
	}
}

func (this *Enemy) destroyed(shot *Shot /*= null*/) bool {
	var vz float32
	if shot != nil {
		explodeVel.x = SHOT_SPEED * Sin32(shot.deg) / 2
		explodeVel.y = SHOT_SPEED * Cos32(shot.deg) / 2
		vz = 0
	} else {
		explodeVel.x = 0
		explodeVel.y = 0
		vz = 0.05
	}
	ss := this.spec.size() * 1.5
	if ss > 2 {
		ss = 2
	}
	var sn float32
	if this.spec.size() < 1 {
		sn = this.spec.size()
	} else {
		sn = sqrt32(this.spec.size())
	}
	if sn > 3 {
		sn = 3
	}
	for i := 0; i < int(sn)*8; i++ {
		NewSmoke(this.pos.x, this.pos.y, 0, nextSignedFloat(0.1)+explodeVel.x, nextSignedFloat(0.1)+explodeVel.y, nextFloat(vz), SmokeTypeEXPLOSION, 32+nextInt(30), ss)
	}
	for i := 0; i < int(sn)*36; i++ {
		NewSpark(this.pos, nextSignedFloat(0.8)+explodeVel.x, nextSignedFloat(0.8)+explodeVel.y, 0.5+nextFloat(0.5), 0.5+nextFloat(0.5), 0, 30+nextInt(30))
	}
	for i := 0; i < int(sn)*12; i++ {
		NewFragment(this.pos, nextSignedFloat(0.33)+explodeVel.x, nextSignedFloat(0.33)+explodeVel.y, 0.05+nextFloat(0.1), 0.2+nextFloat(0.33))
	}
	this.removeTurrets()
	sc := this.spec.score()
	var r bool
	if this.spec.enemyType() == EnemyTypeSMALL {
		playSe("small_destroyed.wav")
		r = false
	} else {
		playSe("destroyed.wav")
		bn := removeAllIndexedBullets(this.idx)
		this.destroyedCnt = 0
		this.explodeCnt = 1
		this.explodeItv = 3
		sc += bn * 10
		r = true
		if this.spec.isBoss() {
			screen.setScreenShake(45, 0.04)
		}
	}
	this.setScoreIndicator(uint32(sc), this.multiplier)
	return r
}

func (this *Enemy) setScoreIndicator(sc uint32, mp float32) {
	ty := getTargetY()
	if mp > 1 {
		ni := NewNumIndicator(sc, IndicatorTypeSCORE, 0.5, this.pos.x, this.pos.y)
		ni.addTarget(8, ty, FlyingToTypeRIGHT, 1, 0.5, sc, 40)
		ni.addTarget(11, ty, FlyingToTypeRIGHT, 0.5, 0.75,
			uint32(float32(sc)*mp), 30)
		ni.addTarget(13, ty, FlyingToTypeRIGHT, 0.25, 1,
			uint32(float32(sc)*mp*stageManager.rank), 20)
		ni.addTarget(12, -8, FlyingToTypeBOTTOM, 0.5, 0.1,
			uint32(float32(sc)*mp*stageManager.rank), 40)
		ni.gotoNextTarget()

		mn := uint32(mp * 1000)
		ni = NewNumIndicator(mn, IndicatorTypeMULTIPLIER, 0.7, this.pos.x, this.pos.y)
		ni.addTarget(10.5, ty, FlyingToTypeRIGHT, 0.5, 0.2, mn, 70)
		ni.gotoNextTarget()

		rn := uint32(stageManager.rank * 1000)
		ni = NewNumIndicator(rn, IndicatorTypeMULTIPLIER, 0.4, 11, 8)
		ni.addTarget(13, ty, FlyingToTypeRIGHT, 0.5, 0.2, rn, 40)
		ni.gotoNextTarget()
		scoreReel.addActualScore(uint32(float32(sc) * mp * stageManager.rank))
	} else {
		ni := NewNumIndicator(sc, IndicatorTypeSCORE, 0.3, this.pos.x, this.pos.y)
		ni.addTarget(11, ty, FlyingToTypeRIGHT, 1.5, 0.2, sc, 40)
		ni.addTarget(13, ty, FlyingToTypeRIGHT, 0.25, 0.25, sc*uint32(stageManager.rank), 20)
		ni.addTarget(12, -8, FlyingToTypeBOTTOM, 0.5, 0.1, sc*uint32(stageManager.rank), 40)
		ni.gotoNextTarget()

		rn := uint32(stageManager.rank * 1000)
		ni = NewNumIndicator(rn, IndicatorTypeMULTIPLIER, 0.4, 11, 8)
		ni.addTarget(13, ty, FlyingToTypeRIGHT, 0.5, 0.2, rn, 40)
		ni.gotoNextTarget()

		scoreReel.addActualScore(uint32(float32(sc) * stageManager.rank))
	}
}

func (this *Enemy) destroyedEdge(n int) {
	playSe("explode.wav")
	sn := n
	if sn > 48 {
		sn = 48
	}
	spp := this.spec.shape().shape.getPointPos()
	spd := this.spec.shape().shape.getPointDeg()
	si := nextInt(len(spp))
	edgePos.x = spp[si].x*this.spec.size() + this.pos.x
	edgePos.y = spp[si].y*this.spec.size() + this.pos.y
	ss := this.spec.size() * 0.5
	if ss > 1 {
		ss = 1
	}
	for i := 0; i < sn; i++ {
		sr := nextFloat(0.5)
		sd := spd[si] + nextSignedFloat(0.2)
		NewSmoke(edgePos.x, edgePos.y, 0, Sin32(sd)*sr, Cos32(sd)*sr, -0.004, SmokeTypeEXPLOSION, 75+nextInt(25), ss)
		for j := 0; j < 2; j++ {
			NewSpark(edgePos, Sin32(sd)*sr*2, Cos32(sd)*sr*2, 0.5+nextFloat(0.5), 0.5+nextFloat(0.5), 0, 30+nextInt(30))
		}
		if i%2 == 0 {
			NewSparkFragment(edgePos, Sin32(sd)*sr*0.5, Cos32(sd)*sr*0.5, 0.06+nextFloat(0.07), (0.2 + nextFloat(0.1)))
		}
	}
}

func (this *Enemy) removeTurrets() {
	for _, t := range this.turretGroup {
		t.close()
	}
	for _, t := range this.movingTurretGroup {
		t.close()
	}
}

func (this *Enemy) draw() {
	gl.PushMatrix()
	if this.destroyedCnt < 0 && this.damagedCnt > 0 {
		damagedPos.x = this.pos.x + nextSignedFloat(float32(this.damagedCnt)*0.01)
		damagedPos.y = this.pos.y + nextSignedFloat(float32(this.damagedCnt)*0.01)
		glTranslate(damagedPos)
	} else {
		glTranslate(this.pos)
	}
	gl.Rotatef(-this.deg*180/Pi32, 0, 0, 1)
	if this.destroyedCnt >= 0 {
		this.spec.destroyedShape().draw()
	} else if !this.damaged {
		this.spec.shape().draw()
	} else {
		this.spec.damagedShape().draw()
	}
	if this.destroyedCnt < 0 {
		if this.spec.bridgeShape() != nil {
			this.spec.bridgeShape().draw()
		}
	}
	gl.PopMatrix()
	if this.destroyedCnt >= 0 {
		return
	}
	for _, t := range this.turretGroup {
		t.draw()
	}
	if this.multiplier > 1 {
		var ox, oy float32
		if this.multiplier < 10 {
			ox = 2.1
		} else {
			ox = 1.4
		}
		oy = 1.25
		if this.spec.isBoss() {
			ox += 4
			oy -= 1.25
		}
		drawNumSignOption(uint32(this.multiplier*1000), this.pos.x+ox, this.pos.y+oy, 0.33, 1, 33, 3)
	}
}

/**
 * Base class for a specification of an enemy.
 */
type EnemyType int

const (
	EnemyTypeSMALL EnemyType = iota
	EnemyTypeLARGE
	EnemyTypePLATFORM
)

type ShipClass int

const (
	ShipClassNONE ShipClass = iota
	ShipClassMIDDLE
	ShipClassLARGE
	ShipClassBOSS
)

type EnemySpec interface {
	move(es *Enemy) bool
	checkCollision(es *Enemy, x float32, y float32, c Shape, shot *Shot) bool
	isSmallEnemy() bool
	checkShipCollision(es *Enemy, x float32, y float32, largeOnly bool /*= false*/) bool
	draw(es *Enemy)
	size() float32
	isBoss() bool
	turretGroupSpec() []*TurretGroupSpec
	turretGroupNum() int
	movingTurretGroupSpec() []*MovingTurretGroupSpec
	movingTurretGroupNum() int
	shield() int
	score() int
	enemyType() EnemyType
	shape() *EnemyShape
	bridgeShape() *EnemyShape
	destroyedShape() *EnemyShape
	damagedShape() *EnemyShape
	setFirstState(es *Enemy, appType AppearanceType, x float32, y float32, d float32) bool
}

type EnemySpecBase struct {
	_shield                                              int
	_size                                                float32
	distRatio                                            float32
	_turretGroupSpec                                     []*TurretGroupSpec
	_movingTurretGroupSpec                               []*MovingTurretGroupSpec
	_movingTurretGroupNum                                int
	_shape, _damagedShape, _destroyedShape, _bridgeShape *EnemyShape
	_enemyType                                           EnemyType
	_shipClass                                           ShipClass
}

func NewEnemySpecBase(enemyType EnemyType) *EnemySpecBase {
	this := new(EnemySpecBase)
	this._turretGroupSpec = make([]*TurretGroupSpec, 0, TURRET_GROUP_MAX)
	this._movingTurretGroupSpec = make([]*MovingTurretGroupSpec, 0, MOVING_TURRET_GROUP_MAX)
	this._shield = 1
	this.sizes(1)
	this._enemyType = enemyType
	return this
}

func (this *EnemySpecBase) movingTurretGroupSpec() []*MovingTurretGroupSpec {
	return this._movingTurretGroupSpec
}

func (this *EnemySpecBase) turretGroupSpec() []*TurretGroupSpec {
	return this._turretGroupSpec
}

func (this *EnemySpecBase) turretGroupNum() int {
	return len(this._turretGroupSpec)
}

func (this *EnemySpecBase) movingTurretGroupNum() int {
	return len(this._movingTurretGroupSpec)
}

func (this *EnemySpecBase) enemyType() EnemyType {
	return this._enemyType
}

func (this *EnemySpecBase) shipClass() ShipClass {
	return this._shipClass
}

func (this *EnemySpecBase) shield() int {
	return this._shield
}

func (this *EnemySpecBase) size() float32 {
	return this._size
}

func (this *EnemySpecBase) destroyedShape() *EnemyShape {
	return this._destroyedShape
}

func (this *EnemySpecBase) damagedShape() *EnemyShape {
	return this._damagedShape
}

func (this *EnemySpecBase) shape() *EnemyShape {
	return this._shape
}

func (this *EnemySpecBase) bridgeShape() *EnemyShape {
	return this._bridgeShape
}

func (this *EnemySpecBase) score() int {
	return 0
}

func (this *EnemySpecBase) isBoss() bool {
	return this._shipClass == ShipClassBOSS
}

func (this *EnemySpecBase) getTurretGroupSpec() *TurretGroupSpec {
	tgs := NewTurretGroupSpec()
	this._turretGroupSpec = append(this._turretGroupSpec, tgs)
	return tgs
}

func (this *EnemySpecBase) getMovingTurretGroupSpec() *MovingTurretGroupSpec {
	tgs := NewMovingTurretGroupSpec()
	this._movingTurretGroupSpec = append(this._movingTurretGroupSpec, tgs)
	return tgs
}

func (this *EnemySpecBase) addMovingTurret(rank float32, bossMode bool /*= false*/) {
	mtn := int(rank * 0.2)
	if mtn > MOVING_TURRET_GROUP_MAX {
		mtn = MOVING_TURRET_GROUP_MAX
	}
	if mtn >= 2 {
		mtn = 1 + nextInt(mtn-1)
	} else {
		mtn = 1
	}
	br := rank / float32(mtn)
	var moveType TurretMoveType
	if !bossMode {
		switch nextInt(4) {
		case 0, 1:
			moveType = TurretMoveTypeROLL
		case 2:
			moveType = TurretMoveTypeSWING_FIX
		case 3:
			moveType = TurretMoveTypeSWING_AIM
		}
	} else {
		moveType = TurretMoveTypeROLL
	}
	rad := 0.9 + nextFloat(0.4) - float32(mtn)*0.1
	radInc := 0.5 + nextFloat(0.25)
	ad := Pi32 * 2
	var a, av, dv, s, sv float32
	switch moveType {
	case TurretMoveTypeROLL:
		a = 0.01 + nextFloat(0.04)
		av = 0.01 + nextFloat(0.03)
		dv = 0.01 + nextFloat(0.04)
	case TurretMoveTypeSWING_FIX:
		ad = Pi32/10 + nextFloat(Pi32/15)
		s = 0.01 + nextFloat(0.02)
		sv = 0.01 + nextFloat(0.03)
	case TurretMoveTypeSWING_AIM:
		ad = Pi32/10 + nextFloat(Pi32/15)
		if nextInt(5) == 0 {
			s = 0.01 + nextFloat(0.01)
		} else {
			s = 0
		}
		sv = 0.01 + nextFloat(0.02)
	}
	for i := 0; i < mtn; i++ {
		tgs := this.getMovingTurretGroupSpec()
		tgs.moveType = moveType
		tgs.radiusBase = rad
		var sr float32
		switch moveType {
		case TurretMoveTypeROLL:
			tgs.alignDeg = ad
			tgs.num = 4 + nextInt(6)
			if nextInt(2) == 0 {
				if nextInt(2) == 0 {
					tgs.setRoll(dv, 0, 0)
				} else {
					tgs.setRoll(-dv, 0, 0)
				}
			} else {
				if nextInt(2) == 0 {
					tgs.setRoll(0, a, av)
				} else {
					tgs.setRoll(0, -a, av)
				}
			}
			if nextInt(3) == 0 {
				tgs.setRadiusAmp(1+nextFloat(1), 0.01+nextFloat(0.03))
			}
			if nextInt(2) == 0 {
				tgs.distRatio = 0.8 + nextSignedFloat(0.3)
			}
			sr = br / float32(tgs.num)
		case TurretMoveTypeSWING_FIX:
			tgs.num = 3 + nextInt(5)
			tgs.alignDeg = ad * (float32(tgs.num)*0.1 + 0.3)
			if nextInt(2) == 0 {
				tgs.setSwing(s, sv, false)
			} else {
				tgs.setSwing(-s, sv, false)
			}
			if nextInt(6) == 0 {
				tgs.setRadiusAmp(1+nextFloat(1), 0.01+nextFloat(0.03))
			}
			if nextInt(4) == 0 {
				tgs.setAlignAmp(0.25+nextFloat(0.25), 0.01+nextFloat(0.02))
			}
			sr = br / float32(tgs.num)
			sr *= 0.6
		case TurretMoveTypeSWING_AIM:
			tgs.num = 3 + nextInt(4)
			tgs.alignDeg = ad * (float32(tgs.num)*0.1 + 0.3)
			if nextInt(2) == 0 {
				tgs.setSwing(s, sv, true)
			} else {
				tgs.setSwing(-s, sv, true)
			}
			if nextInt(4) == 0 {
				tgs.setRadiusAmp(1+nextFloat(1), 0.01+nextFloat(0.03))
			}
			if nextInt(5) == 0 {
				tgs.setAlignAmp(0.25+nextFloat(0.25), 0.01+nextFloat(0.02))
			}
			sr = br / float32(tgs.num)
			sr *= 0.4
		}
		if nextInt(4) == 0 {
			tgs.setXReverse(-1)
		}
		tgs.turretSpec.setParam(sr, TurretTypeMOVING)
		if bossMode {
			tgs.turretSpec.setBossSpec()
		}
		rad += radInc
		ad *= 1 + nextSignedFloat(0.2)
	}
}

func (this *EnemySpecBase) checkCollision(es *Enemy, x float32, y float32, c Shape, shot *Shot) bool {
	return es.checkCollision(x, y, c, shot)
}

func (this *EnemySpecBase) checkShipCollision(es *Enemy, x float32, y float32, largeOnly bool /*= false*/) bool {
	if es.destroyedCnt >= 0 || (largeOnly && this._enemyType != EnemyTypeLARGE) {
		return false
	}
	return this._shape.checkShipCollision(x-es.pos.x, y-es.pos.y, es.deg)
}

func (this *EnemySpecBase) move(es *Enemy) bool {
	return es.stateMove()
}

func (this *EnemySpecBase) draw(es *Enemy) {
	es.draw()
}

func (this *EnemySpecBase) sizes(v float32) float32 {
	this._size = v
	if this._shape != nil {
		this._shape.size = this._size
	}
	if this._damagedShape != nil {
		this._damagedShape.size = this._size
	}
	if this._destroyedShape != nil {
		this._destroyedShape.size = this._size
	}
	if this._bridgeShape != nil {
		var s float32 = 0.9
		this._bridgeShape.size = s * (1 - this.distRatio)
	}
	return this._size
}

func (this *EnemySpecBase) isSmallEnemy() bool {
	return this._enemyType == EnemyTypeSMALL
}

/**
 * Specification for a small class ship.
 */
type MoveType int

const (
	MoveTypeSTOPANDGO MoveType = iota
	MoveTypeCHASE
)

type MoveState int

const (
	MoveStateSTAYING MoveState = iota
	MoveStateMOVING
)

type SmallShipEnemySpec struct {
	*EnemySpecBase

	moveType                   MoveType
	accel, maxSpeed, staySpeed float32
	moveDuration, stayDuration int
	speed, turnDeg             float32
}

func NewSmallShipEnemySpec(rank float32) *SmallShipEnemySpec {
	this := new(SmallShipEnemySpec)
	this.EnemySpecBase = NewEnemySpecBase(EnemyTypeSMALL)
	this.moveDuration = 1
	this.stayDuration = 1
	this._shape = NewEnemyShape(EnemyShapeTypeSMALL)
	this._damagedShape = NewEnemyShape(EnemyShapeTypeSMALL_DAMAGED)
	this._bridgeShape = NewEnemyShape(EnemyShapeTypeSMALL_BRIDGE)
	this.moveType = MoveType(nextInt(2))
	sr := nextFloat(rank * 0.8)
	if sr > 25 {
		sr = 25
	}
	switch this.moveType {
	case MoveTypeSTOPANDGO:
		this.distRatio = 0.5
		this.sizes(0.47 + nextFloat(0.1))
		this.accel = 0.5 - 0.5/(2.0+nextFloat(rank))
		this.maxSpeed = 0.05 * (1.0 + sr)
		this.staySpeed = 0.03
		this.moveDuration = 32 + nextSignedInt(12)
		this.stayDuration = 32 + nextSignedInt(12)
	case MoveTypeCHASE:
		this.distRatio = 0.5
		this.sizes(0.5 + nextFloat(0.1))
		this.speed = 0.036 * (1.0 + sr)
		this.turnDeg = 0.02 + nextSignedFloat(0.04)
	}
	this._shield = 1
	tgs := this.getTurretGroupSpec()
	tgs.turretSpec.setParam(rank-sr*0.5, TurretTypeSMALL)
	return this
}

func (this *SmallShipEnemySpec) setFirstState(es *Enemy, appType AppearanceType, x float32, y float32, d float32) bool {
	if !es.setAppearancePos(appType) {
		return false
	}
	switch this.moveType {
	case MoveTypeSTOPANDGO:
		es.speed = 0
		es.state = MoveStateMOVING
		es.cnt = this.moveDuration
	case MoveTypeCHASE:
		es.speed = this.speed
	}
	return true
}

func (this *SmallShipEnemySpec) move(es *Enemy) bool {
	if !this.EnemySpecBase.move(es) {
		return false
	}
	switch this.moveType {
	case MoveTypeSTOPANDGO:
		es.pos.x += Sin32(es.velDeg) * es.speed
		es.pos.y += Cos32(es.velDeg) * es.speed
		es.pos.y -= field.lastScrollY
		if es.pos.y <= -field.outerSize.y {
			return false
		}
		if field.getBlockVector(es.pos) >= 0 || !field.checkInOuterHeightField(es.pos) {
			es.velDeg += Pi32
			es.pos.x += Sin32(es.velDeg) * es.speed * 2
			es.pos.y += Cos32(es.velDeg) * es.speed * 2
		}
		switch es.state {
		case MoveStateMOVING:
			es.speed += (this.maxSpeed - es.speed) * this.accel
			es.cnt--
			if es.cnt <= 0 {
				es.velDeg = nextFloat(Pi32 * 2)
				es.cnt = this.stayDuration
				es.state = MoveStateSTAYING
			}
		case MoveStateSTAYING:
			es.speed += (this.staySpeed - es.speed) * this.accel
			es.cnt--
			if es.cnt <= 0 {
				es.cnt = this.moveDuration
				es.state = MoveStateMOVING
			}
		}
	case MoveTypeCHASE:
		es.pos.x += Sin32(es.velDeg) * this.speed
		es.pos.y += Cos32(es.velDeg) * this.speed
		es.pos.y -= field.lastScrollY
		if es.pos.y <= -field.outerSize.y {
			return false
		}
		if field.getBlockVector(es.pos) >= 0 || !field.checkInOuterHeightField(es.pos) {
			es.velDeg += Pi32
			es.pos.x += Sin32(es.velDeg) * es.speed * 2
			es.pos.y += Cos32(es.velDeg) * es.speed * 2
		}
		var ad float32
		ship.nearPos(es.pos)
		shipPos := ship._nearPos
		if shipPos.distVector(es.pos) < 0.1 {
			ad = 0
		} else {
			ad = atan232(shipPos.x-es.pos.x, shipPos.y-es.pos.y)
		}
		od := ad - es.velDeg
		od = normalizeDeg(od)
		if od <= this.turnDeg && od >= -this.turnDeg {
			es.velDeg = ad
		} else if od < 0 {
			es.velDeg -= this.turnDeg
		} else {
			es.velDeg += this.turnDeg
		}
		es.velDeg = normalizeDeg(es.velDeg)
		es.cnt++
	}
	od := es.velDeg - es.deg
	od = normalizeDeg(od)
	es.deg += od * 0.05
	es.deg = normalizeDeg(es.deg)
	if es.cnt%6 == 0 && es.speed >= 0.03 {
		this._shape.addWake(es.pos, es.deg, es.speed)
	}
	return true
}

func (this *SmallShipEnemySpec) score() int {
	return 50
}

func (this *SmallShipEnemySpec) isBoss() bool {
	return false
}

/**
 * Specification for a large/middle class ship.
 */

const SINK_INTERVAL = 120

type ShipEnemySpec struct {
	*EnemySpecBase

	speed, degVel float32
}

func NewShipEnemySpec(rank float32, cls ShipClass) *ShipEnemySpec {
	this := new(ShipEnemySpec)
	this.EnemySpecBase = NewEnemySpecBase(EnemyTypeLARGE)
	this._shipClass = cls
	this._shape = NewEnemyShape(EnemyShapeTypeMIDDLE)
	this._damagedShape = NewEnemyShape(EnemyShapeTypeMIDDLE_DAMAGED)
	this._destroyedShape = NewEnemyShape(EnemyShapeTypeMIDDLE_DESTROYED)
	this._bridgeShape = NewEnemyShape(EnemyShapeTypeMIDDLE_BRIDGE)
	this.distRatio = 0.7
	var mainTurretNum, subTurretNum int
	var movingTurretRatio float32
	rk := rank
	switch cls {
	case ShipClassMIDDLE:
		sz := 1.5 + rank/15 + nextFloat(rank/15)
		ms := 2 + nextFloat(0.5)
		if sz > ms {
			sz = ms
		}
		this.sizes(sz)
		this.speed = 0.015 + nextSignedFloat(0.005)
		this.degVel = 0.005 + nextSignedFloat(0.003)
		switch nextInt(3) {
		case 0:
			mainTurretNum = int(this._size*(1+nextSignedFloat(0.25)) + 1)
		case 1:
			subTurretNum = int(this._size*1.6*(1+nextSignedFloat(0.5)) + 2)
		case 2:
			mainTurretNum = int(this._size*(0.5+nextSignedFloat(0.12)) + 1)
			movingTurretRatio = 0.5 + nextFloat(0.25)
			rk = rank * (1 - movingTurretRatio)
			movingTurretRatio *= 2
		}
	case ShipClassLARGE:
		sz := 2.5 + rank/24 + nextFloat(rank/24)
		ms := 3 + nextFloat(1)
		if sz > ms {
			sz = ms
		}
		this.sizes(sz)
		this.speed = 0.01 + nextSignedFloat(0.005)
		this.degVel = 0.003 + nextSignedFloat(0.002)
		mainTurretNum = int(this._size*(0.7+nextSignedFloat(0.2)) + 1)
		subTurretNum = int(this._size*1.6*(0.7+nextSignedFloat(0.33)) + 2)
		movingTurretRatio = 0.25 + nextFloat(0.5)
		rk = rank * (1 - movingTurretRatio)
		movingTurretRatio *= 3
	case ShipClassBOSS:
		sz := 5 + rank/30 + nextFloat(rank/30)
		ms := 9 + nextFloat(3)
		if sz > ms {
			sz = ms
		}
		this.sizes(sz)
		this.speed = ship.scrollSpeedBase + 0.0025 + nextSignedFloat(0.001)
		this.degVel = 0.003 + nextSignedFloat(0.002)
		mainTurretNum = int(this._size*0.8*(1.5+nextSignedFloat(0.4)) + 2)
		subTurretNum = int(this._size*0.8*(2.4+nextSignedFloat(0.6)) + 2)
		movingTurretRatio = 0.2 + nextFloat(0.3)
		rk = rank * (1 - movingTurretRatio)
		movingTurretRatio *= 2.5
	}
	this._shield = int(this._size * 10)
	if cls == ShipClassBOSS {
		this._shield = int(float32(this._shield) * 2.4)
	}
	if mainTurretNum+subTurretNum <= 0 {
		tgs := this.getTurretGroupSpec()
		tgs.turretSpec.setParam(0, TurretTypeDUMMY)
	} else {
		var subTurretRank float32 = rk / float32(mainTurretNum*3+subTurretNum)
		var mainTurretRank float32 = subTurretRank * 2.5
		if cls != ShipClassBOSS {
			frontMainTurretNum := int(float32(mainTurretNum)/2 + 0.99)
			rearMainTurretNum := mainTurretNum - frontMainTurretNum
			if frontMainTurretNum > 0 {
				tgs := this.getTurretGroupSpec()
				tgs.turretSpec.setParam(mainTurretRank, TurretTypeMAIN)
				tgs.num = frontMainTurretNum
				tgs.alignType = AlignTypeSTRAIGHT
				tgs.offset.y = -this._size * (0.9 + nextSignedFloat(0.05))
			}
			if rearMainTurretNum > 0 {
				tgs := this.getTurretGroupSpec()
				tgs.turretSpec.setParam(mainTurretRank, TurretTypeMAIN)
				tgs.num = rearMainTurretNum
				tgs.alignType = AlignTypeSTRAIGHT
				tgs.offset.y = this._size * (0.9 + nextSignedFloat(0.05))
			}
			var pts *TurretSpec
			if subTurretNum > 0 {
				frontSubTurretNum := (subTurretNum + 2) / 4
				rearSubTurretNum := (subTurretNum - frontSubTurretNum*2) / 2
				tn := frontSubTurretNum
				ad := -Pi32 / 4
				for i := 0; i < 4; i++ {
					if i == 2 {
						tn = rearSubTurretNum
					}
					if tn <= 0 {
						continue
					}
					tgs := this.getTurretGroupSpec()
					if i == 0 || i == 2 {
						if nextInt(2) == 0 {
							tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB)
						} else {
							tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB_DESTRUCTIVE)
						}
						pts = tgs.turretSpec
					} else {
						tgs.turretSpec.setParamTurretSpec(pts)
					}
					tgs.num = tn
					tgs.alignType = AlignTypeROUND
					tgs.alignDeg = ad
					ad += Pi32 / 2
					tgs.alignWidth = Pi32/6 + nextFloat(Pi32/8)
					tgs.radius = this._size * 0.75
					tgs.distRatio = this.distRatio
				}
			}
		} else {
			mainTurretRank *= 2.5
			subTurretRank *= 2
			var pts *TurretSpec
			if mainTurretNum > 0 {
				frontMainTurretNum := (mainTurretNum + 2) / 4
				rearMainTurretNum := (mainTurretNum - frontMainTurretNum*2) / 2
				tn := frontMainTurretNum
				ad := -Pi32 / 4
				for i := 0; i < 4; i++ {
					if i == 2 {
						tn = rearMainTurretNum
					}
					if tn <= 0 {
						continue
					}
					tgs := this.getTurretGroupSpec()
					if i == 0 || i == 2 {
						tgs.turretSpec.setParam(mainTurretRank, TurretTypeMAIN)
						pts = tgs.turretSpec
						pts.setBossSpec()
					} else {
						tgs.turretSpec.setParamTurretSpec(pts)
					}
					tgs.num = tn
					tgs.alignType = AlignTypeROUND
					tgs.alignDeg = ad
					ad += Pi32 / 2
					tgs.alignWidth = Pi32/6 + nextFloat(Pi32/8)
					tgs.radius = this._size * 0.45
					tgs.distRatio = this.distRatio
				}
			}
			if subTurretNum > 0 {
				var tn [3]int
				tn[0] = (subTurretNum + 2) / 6
				tn[1] = (subTurretNum - tn[0]*2) / 4
				tn[2] = (subTurretNum - tn[0]*2 - tn[1]*2) / 2
				ad := []float32{Pi32 / 4, -Pi32 / 4, Pi32 / 2, -Pi32 / 2, Pi32 / 4 * 3, -Pi32 / 4 * 3}
				for i := 0; i < 6; i++ {
					idx := i / 2
					if tn[idx] <= 0 {
						continue
					}
					tgs := this.getTurretGroupSpec()
					if i == 0 || i == 2 || i == 4 {
						if nextInt(2) == 0 {
							tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB)
						} else {
							tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB_DESTRUCTIVE)
						}
						pts = tgs.turretSpec
						pts.setBossSpec()
					} else {
						tgs.turretSpec.setParamTurretSpec(pts)
					}
					tgs.num = tn[idx]
					tgs.alignType = AlignTypeROUND
					tgs.alignDeg = ad[i]
					tgs.alignWidth = Pi32/7 + nextFloat(Pi32/9)
					tgs.radius = this._size * 0.75
					tgs.distRatio = this.distRatio
				}
			}
		}
	}
	if movingTurretRatio > 0 {
		if cls == ShipClassBOSS {
			this.addMovingTurret(rank*movingTurretRatio, true)
		} else {
			this.addMovingTurret(rank*movingTurretRatio, false)
		}
	}
	return this
}

func (this *ShipEnemySpec) setFirstState(es *Enemy, appType AppearanceType, x float32, y float32, d float32) bool {
	if !es.setAppearancePos(appType) {
		return false
	}
	es.speed = this.speed
	if es.pos.x < 0 {
		es.turnWay = -1
	} else {
		es.turnWay = 1
	}
	if this.isBoss() {
		es.trgDeg = nextFloat(0.1) + 0.1
		if nextInt(2) == 0 {
			es.trgDeg *= -1
		}
		es.turnCnt = 250 + nextInt(150)
	}
	return true
}

func (this *ShipEnemySpec) move(es *Enemy) bool {
	if es.destroyedCnt >= SINK_INTERVAL {
		return false
	}
	if !this.EnemySpecBase.move(es) {
		return false
	}
	es.pos.x += Sin32(es.deg) * es.speed
	es.pos.y += Cos32(es.deg) * es.speed
	es.pos.y -= field.lastScrollY
	if es.pos.x <= -field.outerSize.x-this._size || es.pos.x >= field.outerSize.x+this._size ||
		es.pos.y <= -field.outerSize.y-this._size {
		return false
	}
	if es.pos.y > field.outerSize.y*2.2+this._size {
		es.pos.y = field.outerSize.y*2.2 + this._size
	}
	if this.isBoss() {
		es.turnCnt--
		if es.turnCnt <= 0 {
			es.turnCnt = 250 + nextInt(150)
			es.trgDeg = nextFloat(0.1) + 0.2
			if es.pos.x > 0 {
				es.trgDeg *= -1
			}
		}
		es.deg += (es.trgDeg - es.deg) * 0.0025
		if ship.higherPos().y > es.pos.y {
			es.speed += (this.speed*2 - es.speed) * 0.005
		} else {
			es.speed += (this.speed - es.speed) * 0.01
		}
	} else {
		if !es.checkFrontClear(false) {
			es.deg += this.degVel * es.turnWay
			es.speed *= 0.98
		} else {
			if es.destroyedCnt < 0 {
				es.speed += (this.speed - es.speed) * 0.01
			} else {
				es.speed *= 0.98
			}
		}
	}
	es.cnt++
	if es.cnt%6 == 0 && es.speed >= 0.01 && es.destroyedCnt < SINK_INTERVAL/2 {
		this._shape.addWake(es.pos, es.deg, es.speed)
	}
	return true
}

func (this *ShipEnemySpec) draw(es *Enemy) {
	if es.destroyedCnt >= 0 {
		setScreenColor(
			MIDDLE_COLOR_R*(1-float32(es.destroyedCnt)/SINK_INTERVAL)*0.5,
			MIDDLE_COLOR_G*(1-float32(es.destroyedCnt)/SINK_INTERVAL)*0.5,
			MIDDLE_COLOR_B*(1-float32(es.destroyedCnt)/SINK_INTERVAL)*0.5, 1)
	}
	this.EnemySpecBase.draw(es)
}

func (this *ShipEnemySpec) score() int {
	switch this._shipClass {
	case ShipClassMIDDLE:
		return 100
	case ShipClassLARGE:
		return 300
	case ShipClassBOSS:
		return 1000
	}
	return 0
}

/**
 * Specification for a sea-based platform.
 */
type PlatformEnemySpec struct {
	*EnemySpecBase
}

func NewPlatformEnemySpec(rank float32) *PlatformEnemySpec {
	this := new(PlatformEnemySpec)
	this.EnemySpecBase = NewEnemySpecBase(EnemyTypePLATFORM)
	this._shape = NewEnemyShape(EnemyShapeTypePLATFORM)
	this._damagedShape = NewEnemyShape(EnemyShapeTypePLATFORM_DAMAGED)
	this._destroyedShape = NewEnemyShape(EnemyShapeTypePLATFORM_DESTROYED)
	this._bridgeShape = NewEnemyShape(EnemyShapeTypePLATFORM_BRIDGE)
	this.distRatio = 0
	this.sizes(1 + rank/30 + nextFloat(rank/30))
	ms := 1 + nextFloat(0.25)
	if this.size() > ms {
		this.sizes(ms)
	}
	var mainTurretNum, frontTurretNum, sideTurretNum int
	rk := rank
	var movingTurretRatio float32
	switch nextInt(3) {
	case 0:
		frontTurretNum = int(this._size*(2+nextSignedFloat(0.5)) + 1)
		movingTurretRatio = 0.33 + nextFloat(0.46)
		rk *= (1 - movingTurretRatio)
		movingTurretRatio *= 2.5
	case 1:
		frontTurretNum = int(this._size*(0.5+nextSignedFloat(0.2)) + 1)
		sideTurretNum = int(this._size*(0.5+nextSignedFloat(0.2))+1) * 2
	case 2:
		mainTurretNum = int(this._size*(1+nextSignedFloat(0.33)) + 1)
	}
	this._shield = int(this._size * 20)
	subTurretNum := frontTurretNum + sideTurretNum
	subTurretRank := int(rk) / (mainTurretNum*3 + subTurretNum)
	mainTurretRank := int(float32(subTurretRank) * 2.5)
	if mainTurretNum > 0 {
		tgs := this.getTurretGroupSpec()
		tgs.turretSpec.setParam(float32(mainTurretRank), TurretTypeMAIN)
		tgs.num = mainTurretNum
		tgs.alignType = AlignTypeROUND
		tgs.alignDeg = 0
		tgs.alignWidth = Pi32*0.66 + nextFloat(Pi32/2)
		tgs.radius = this._size * 0.7
		tgs.distRatio = this.distRatio
	}
	if frontTurretNum > 0 {
		tgs := this.getTurretGroupSpec()
		tgs.turretSpec.setParam(float32(subTurretRank), TurretTypeSUB)
		tgs.num = frontTurretNum
		tgs.alignType = AlignTypeROUND
		tgs.alignDeg = 0
		tgs.alignWidth = Pi32/5 + nextFloat(Pi32/6)
		tgs.radius = this._size * 0.8
		tgs.distRatio = this.distRatio
	}
	sideTurretNum /= 2
	if sideTurretNum > 0 {
		var pts *TurretSpec
		for i := 0; i < 2; i++ {
			tgs := this.getTurretGroupSpec()
			if i == 0 {
				tgs.turretSpec.setParam(float32(subTurretRank), TurretTypeSUB)
				pts = tgs.turretSpec
			} else {
				tgs.turretSpec.setParamTurretSpec(pts)
			}
			tgs.num = sideTurretNum
			tgs.alignType = AlignTypeROUND
			tgs.alignDeg = Pi32/2 - Pi32*float32(i)
			tgs.alignWidth = Pi32/5 + nextFloat(Pi32/6)
			tgs.radius = this._size * 0.75
			tgs.distRatio = this.distRatio
		}
	}
	if movingTurretRatio > 0 {
		this.addMovingTurret(rank*movingTurretRatio, false)
	}
	return this
}

func (this *PlatformEnemySpec) setFirstState(es *Enemy, appType AppearanceType, x float32, y float32, d float32) bool {
	es.pos.x = x
	es.pos.y = y
	es.deg = d
	es.speed = 0
	return es.checkFrontClear(true)
}

func (this *PlatformEnemySpec) move(es *Enemy) bool {
	if !this.EnemySpecBase.move(es) {
		return false
	}
	es.pos.y -= field.lastScrollY
	return !(es.pos.y <= -field.outerSize.y)
}

func (this *PlatformEnemySpec) score() int {
	return 100
}

func (this *PlatformEnemySpec) isBoss() bool {
	return false
}

/* Actor Pool Functions
 *
 * functions that run across all enemies
 */

func checkAllEnemiesShotHit(pos Vector, shape Shape, shot *Shot /*= null*/) {
	for a, _ := range actors {
		e, ok := a.(*Enemy)
		if ok {
			e.checkShotHit(pos, shape, shot)
		}
	}
}

func checkAllEnemiesHitShip(x float32, y float32, deselection *Enemy /*= null*/, largeOnly bool /*= false*/) *Enemy {
	for a, _ := range actors {
		e, ok := a.(*Enemy)
		if ok && e != deselection {
			if e.checkHitShip(x, y, largeOnly) {
				return e
			}
		}
	}
	return nil
}

func hasBoss() bool {
	for a, _ := range actors {
		e, ok := a.(*Enemy)
		if ok && e.isBoss() {
			return true
		}
	}
	return false
}
