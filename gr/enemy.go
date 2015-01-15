/*
 * $Id: enemy.d,v 1.2 2005/07/17 11:02:45 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Enemy ships.
 */
type Enemy struct {
	spec  EnemySpec
	state EnemyState
}

func NewEnemy(field Field, screen Screen, ship Ship, scoreReel ScoreReel) *Enemy {
	e := new(Enemy)
	e.state = NewEnemyState(field, screen, ship, scoreReel)
	return e
}

func (this *Enemy) setStageManager(stageManager StageManager) {
	this.state.setStageManager(stageManager)
}

func (this *Enemy) set(spec EnemySpec) {
	this.spec = spec
	this.exists = true
}

func (this *Enemy) move() {
	if !this.spec.move(this.state) {
		this.remove()
	}
}

func (this *Enemy) checkShotHit(p Vector, shape Collidable, shot Shot) {
	if this.state.destroyedCnt >= 0 {
		return
	}
	if this.spec.checkCollision(this.state, p.x, p.y, shape, shot) {
		if shot {
			shot.removeHitToEnemy(this.spec.isSmallEnemy)
		}
	}
}

func (this *Enemy) checkHitShip(x float32, y float32, largeOnly bool /*= false*/) bool {
	return this.spec.checkShipCollision(this.state, x, y, largeOnly)
}

func (this *Enemy) addDamage(n int) {
	this.state.addDamage(n)
}

func (this *Enemy) increaseMultiplier(m float32) {
	this.state.increaseMultiplier(m)
}

func (this *Enemy) addScore(s int) {
	this.state.addScore(s)
}

func (this *Enemy) remove() {
	this.state.removeTurrets()
	this.exists = false
}

func (this *Enemy) draw() {
	this.spec.draw(this.state)
}

func (this *Enemy) pos() Vector {
	return this.state.pos
}

func (this *Enemy) size() float32 {
	return this.spec.size
}

func (this *Enemy) ndex() int {
	return this.state.idx
}

func (this *Enemy) isBoss() bool {
	return this.spec.isBoss
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

type EnemyState struct {
	appType                                               int
	ppos, pos                                             Vector
	shield                                                int
	deg, velDeg, speed, turnWay, trgDeg                   float32
	turnCnt, state, cnt                                   int
	vel                                                   Vector
	turretGroup                                           [TURRET_GROUP_MAX]turretGroup
	movingTurretGroup                                     [MOVING_TURRET_GROUP_MAX]MovingTurretGroup
	damaged                                               bool
	damagedCnt, destroyedCnt, explodeCnt, explodeItv, idx int
	multiplier                                            float32
	spec                                                  EnemySpec

	field        Field
	screen       Screen
	ship         Ship
	enemy        Enemy
	stageManager StageManager
	scoreReel    ScoreReel
}

func NewEnemyState(field Field, screen Screen, ship Ship, scoreReel ScoreReel) EnemyState {
	this := new(EnemyState)
	this.idx = idxCount
	idxCount++
	this.field = field
	this.screen = screen
	this.bullets = bullets
	this.ship = ship
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.sparkFragments = sparkFragments
	this.numIndicators = numIndicators
	this.scoreReel = scoreReel
	turnWay = 1
	explodeItv = 1
	multiplier = 1
}

func (this *EnemyState) setEnemyAndPool(enemy Enemy) {
	this.enemy = enemy
	this.enemies = enemies
	for i, _ := range turretGroup {
		this.turretGroup[i] = NewTurretGroup(field, bullets, ship, sparks, smokes, fragments, enemy)
	}
	for i, _ := range movingTurretGroup {
		this.movingTurretGroup[i] = NewMovingTurretGroup(field, bullets, ship, sparks, smokes, fragments, enemy)
	}
}

func (this *EnemyState) setStageManager(stageManager StageManager) {
	this.stageManager = stageManager
}

func (this *EnemyState) setSpec(spec EnemySpec) {
	this.spec = spec
	this.shield = spec.shield
	for i := 0; i < spec.turretGroupNum; i++ {
		this.turretGroup[i].set(spec.turretGroupSpec[i])
	}
	for i := 0; i < spec.movingTurretGroupNum; i++ {
		this.movingTurretGroup[i].set(spec.movingTurretGroupSpec[i])
	}
	this.cnt = 0
	this.damaged = false
	this.damagedCnt = 0
	this.destroyedCnt = -1
	this.explodeCnt = 0
	this.explodeItv = 1
	this.multiplier = 1
}

func (this *EnemyState) setAppearancePos(field Field, ship Ship, appType int /*= AppearanceType.TOP*/) bool {
	this.appType = appType
	for i := 0; i < 8; i++ {
		switch appType {
		case AppearanceType.TOP:
			this.pos.x = rand.nextSignedFloat(field.size.x)
			this.pos.y = field.outerSize.y*0.99 + this.spec.size
			if this.pos.x < 0 {
				this.deg = Pi3232 - rand.nextFloat(0.5)
				this.velDeg = this.deg
			} else {
				this.deg = Pi3232 + rand.nextFloat(0.5)
				this.velDeg = this.deg
			}
			break
		case AppearanceType.SIDE:
			if rand.nextInt(2) == 0 {
				this.pos.x = -field.outerSize.x * 0.99
				this.deg = Pi3232/2 + rand.nextFloat(0.66)
				this.velDeg = this.deg
			} else {
				this.pos.x = field.outerSize.x * 0.99
				this.deg = -Pi32/2 - rand.nextFloat(0.66)
				this.velDeg = this.deg
			}
			this.pos.y = field.size.y + rand.nextFloat(field.size.y) + this.spec.size
			break
		case AppearanceType.CENTER:
			this.pos.x = 0
			this.pos.y = field.outerSize.y*0.99 + this.spec.size
			this.deg = 0
			this.velDeg = this.deg
			break
		}
		this.ppos.x = this.pos.x
		this.ppos.y = this.pos.y
		this.vel.y = 0
		this.vel.x = 0
		this.speed = 0
		if this.appType == AppearanceType.CENTER || this.checkFrontClear(true) {
			return true
		}
	}
	return false
}

func (this *EnemyState) checkFrontClear(checkCurrentPos bool /*= false*/) bool {
	var si = 1
	if this.checkCurrentPos() {
		si = 0
	}
	for i := si; i < 5; i++ {
		cx := this.pos.x + Sin32(deg)*i*this.spec.size
		cy := this.pos.y + Cos32(deg)*i*this.spec.size
		if this.field.getBlock(cx, cy) >= 0 {
			return false
		}
		if checkAllEnemiesHitShip(cx, cy, enemy, true) {
			return false
		}
	}
	return true
}

func (this *EnemyState) move() bool {
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
			this.explodeItv = this.explodeItv * (1.2 + rand.nextFloat(1))
			this.explodeCnt = this.explodeItv
			this.destroyedEdge(sqrt(this.spec.size) * 27.0 / (this.explodeItv*0.1 + 1))
		}
	}
	this.damaged = false
	if this.damagedCnt > 0 {
		this.damagedCnt--
	}
	alive := false
	for i := 0; i < this.spec.turretGroupNum; i++ {
		alive |= this.turretGroup[i].move(this.pos, this.deg)
	}
	for i := 0; i < this.spec.movingTurretGroupNum; i++ {
		this.movingTurretGroup[i].move(this.pos, this.deg)
	}
	if this.destroyedCnt < 0 && !alive {
		return this.destroyed()
	}
	return true
}

func (this *EnemyState) checkCollision(x float32, y float32, c Collidable, shot Shot) bool {
	ox := fabs32(this.pos.x - x)
	oy := fabs32(this.pos.y - y)
	if ox+oy > this.spec.size*2 {
		return false
	}
	for i := 0; i < this.spec.turretGroupNum; i++ {
		if this.turretGroup[i].checkCollision(x, y, c, shot) {
			return true
		}
	}
	if this.spec.bridgeShape.checkCollision(ox, oy, c) {
		this.addDamage(shot.damage, shot)
		return true
	}
	return false
}

func (this *EnemyState) increaseMultiplier(m float32) {
	this.multiplier += m
}

func (this *EnemyState) addScore(s int) {
	this.setScoreIndicator(s, 1)
}

func (this *EnemyState) addDamage(n int, shot Shot /*= null*/) {
	this.shield -= n
	if this.shield <= 0 {
		this.destroyed(shot)
	} else {
		this.damaged = true
		this.damagedCnt = 7
	}
}

func (this *EnemyState) destroyed(shot Shot /*= null*/) bool {
	var vz float32
	if shot != nil {
		this.explodeVel.x = Shot.SPEED * Sin32(shot.deg) / 2
		this.explodeVel.y = Shot.SPEED * Cos32(shot.deg) / 2
		vz = 0
	} else {
		this.explodeVel.x = 0
		this.explodeVel.y = 0
		vz = 0.05
	}
	ss := this.spec.size * 1.5
	if ss > 2 {
		ss = 2
	}
	var sn float32
	if this.spec.size < 1 {
		sn = this.spec.size
	} else {
		sn = sqrt32(spec.size)
	}
	if sn > 3 {
		sn = 3
	}
	for i := 0; i < sn*8; i++ {
		NewSmoke(this.pos, rand.nextSignedFloat(0.1)+this.explodeVel.x, rand.nextSignedFloat(0.1)+this.explodeVel.y, rand.nextFloat(vz), SmokeTypeEXPLOSION, 32+rand.nextInt(30), ss)
	}
	for i := 0; i < sn*36; i++ {
		NewSpark(this.pos, rand.nextSignedFloat(0.8)+this.explodeVel.x, rand.nextSignedFloat(0.8)+this.explodeVel.y, 0.5+rand.nextFloat(0.5), 0.5+rand.nextFloat(0.5), 0, 30+rand.nextInt(30))
	}
	for i := 0; i < sn*12; i++ {
		NewFragment(this.pos, rand.nextSignedFloat(0.33)+this.explodeVel.x, rand.nextSignedFloat(0.33)+this.explodeVel.y, 0.05+rand.nextFloat(0.1), 0.2+rand.nextFloat(0.33))
	}
	this.removeTurrets()
	sc := this.spec.score
	var r bool
	if this.spec.enemyType == EnemyTypeSMALL {
		playSe("small_destroyed.wav")
		r = false
	} else {
		playSe("destroyed.wav")
		bn := removeAllBulletsIndexedBullets(idx)
		this.destroyedCnt = 0
		this.explodeCnt = 1
		this.explodeItv = 3
		sc += bn * 10
		r = true
		if this.spec.isBoss {
			setScreenShake(45, 0.04)
		}
	}
	this.setScoreIndicator(sc, multiplier)
	return r
}

func (this *EnemyState) setScoreIndicator(sc int, mp float32) {
	ty := getTargetY()
	if mp > 1 {
		ni := NewNumIndicator(sc, NumIndicator.IndicatorType.SCORE, 0.5, this.pos)
		ni.addTarget(8, ty, FlyingToTypeRIGHT, 1, 0.5, sc, 40)
		ni.addTarget(11, ty, FlyingToTypeRIGHT, 0.5, 0.75,
			(sc * mp), 30)
		ni.addTarget(13, ty, FlyingToTypeRIGHT, 0.25, 1,
			(sc * mp * this.stageManager.rankMultiplier), 20)
		ni.addTarget(12, -8, FlyingToTypeBOTTOM, 0.5, 0.1,
			(sc * mp * this.stageManager.rankMultiplier), 40)
		ni.gotoNextTarget()

		mn := int(mp * 1000)
		ni = NewNumIndicator(mn, IndicatorTypeMULTIPLIER, 0.7, this.pos)
		ni.addTarget(10.5, ty, FlyingToTypeRIGHT, 0.5, 0.2, mn, 70)
		ni.gotoNextTarget()

		rn := int(this.stageManager.rankMultiplier * 1000)
		ni = NewNumIndicator(rn, IndicatorTypeMULTIPLIER, 0.4, 11, 8)
		ni.addTarget(13, ty, FlyingToTypeRIGHT, 0.5, 0.2, rn, 40)
		ni.gotoNextTarget()
		this.scoreReel.addActualScore(int(sc * mp * stageManager.rankMultiplier))
	} else {
		ni := NewNumIndicator(sc, IndicatorTypeSCORE, 0.3, this.pos)
		ni.addTarget(11, ty, FlyingToTypeRIGHT, 1.5, 0.2, sc, 40)
		ni.addTarget(13, ty, FlyingToTypeRIGHT, 0.25, 0.25, int(sc*this.stageManager.rankMultiplier), 20)
		ni.addTarget(12, -8, FlyingToTypeBOTTOM, 0.5, 0.1, int(sc*this.stageManager.rankMultiplier), 40)
		ni.gotoNextTarget()

		rn := int(this.stageManager.rankMultiplier * 1000)
		ni = NewNumIndicator(rn, IndicatorTypeMULTIPLIER, 0.4, 11, 8)
		ni.addTarget(13, ty, FlyingToTypeRIGHT, 0.5, 0.2, rn, 40)
		ni.gotoNextTarget()

		this.scoreReel.addActualScore(int(sc * this.stageManager.rankMultiplier))
	}
}

func (this *EnemyState) destroyedEdge(n int) {
	playSe("explode.wav")
	sn := n
	if sn > 48 {
		sn = 48
	}
	spp := this.spec.shape.shape.pointPos
	spd := this.spec.shape.shape.pointDeg
	i := rand.nextInt(spp.length)
	this.edgePos.x = spp[si].x*this.spec.size + this.pos.x
	this.edgePos.y = spp[si].y*this.spec.size + this.pos.y
	ss := this.spec.size * 0.5
	if ss > 1 {
		ss = 1
	}
	for i := 0; i < sn; i++ {
		sr := rand.nextFloat(0.5)
		sd := spd[si] + rand.nextSignedFloat(0.2)
		s := NewSmoke(this.edgePos, Sin32(sd)*sr, Cos32(sd)*sr, -0.004, SmokeTypeEXPLOSION, 75+rand.nextInt(25), ss)
		for j := 0; j < 2; j++ {
			NewSpark(this.edgePos, Sin32(sd)*sr*2, Cos32(sd)*sr*2, 0.5+rand.nextFloat(0.5), 0.5+rand.nextFloat(0.5), 0, 30+rand.nextInt(30))
		}
		if i%2 == 0 {
			NewSparkFragment(this.edgePos, Sin32(sd)*sr*0.5, Cos32(sd)*sr*0.5, 0.06+rand.nextFloat(0.07), (0.2 + rand.nextFloat(0.1)))
		}
	}
}

func (this *EnemyState) removeTurrets() {
	for i := 0; i < this.spec.turretGroupNum; i++ {
		this.turretGroup[i].remove()
	}
	for i := 0; i < this.spec.movingTurretGroupNum; i++ {
		this.movingTurretGroup[i].remove()
	}
}

func (this *EnemyState) draw() {
	glPushMatrix()
	if this.destroyedCnt < 0 && this.damagedCnt > 0 {
		this.damagedPos.x = this.pos.x + rand.nextSignedFloat(damagedCnt*0.01)
		this.damagedPos.y = this.pos.y + rand.nextSignedFloat(damagedCnt*0.01)
		glTranslate(this.damagedPos)
	} else {
		glTranslate(this.pos)
	}
	glRotatef(-this.deg*180/Pi32, 0, 0, 1)
	if this.destroyedCnt >= 0 {
		this.spec.destroyedShape.draw()
	} else if !this.damaged {
		this.spec.shape.draw()
	} else {
		this.spec.damagedShape.draw()
	}
	if destroyedCnt < 0 {
		this.spec.bridgeShape.draw()
	}
	glPopMatrix()
	if this.destroyedCnt >= 0 {
		return
	}
	for i := 0; i < this.spec.turretGroupNum; i++ {
		this.turretGroup[i].draw()
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
		drawNumSign(int(multiplier*1000), this.pos.x+ox, this.pos.y+oy, 0.33, 1, 33, 3)
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

type EnemySpec struct {
	field                                            Field
	ship                                             Ship
	shield                                           int
	size                                             float32
	distRatio                                        float32
	turretGroupSpec                                  [TURRET_GROUP_MAX]TurretGroupSpec
	turretGroupNum                                   int
	movingTurretGroupSpec                            [MOVING_TURRET_GROUP_MAX]MovingTurretGroupSpec
	movingTurretGroupNum                             int
	shape, damagedShape, destroyedShape, bridgeShape EnemyShape
	enemyType                                        EnemyType
}

func NewEnemySpec(field Field, ship Ship, enemyType EnemyType) *EnemySpec {
	this := new(EnemySpec)
	this.field = field
	this.ship = ship
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.wakes = wakes
	for i, _ := range this.turretGroupSpec {
		this.turretGroupSpec[i] = NewTurretGroupSpec()
	}
	for i, _ := range this.movingTurretGroupSpec {
		this.movingTurretGroupSpec[i] = NewMovingTurretGroupSpec()
	}
	this.shield = 1
	this.size = 1
	this.enemyType = enemyType
}

func (this *EnemySpec) getTurretGroupSpec() *TurretGroupSpec {
	this.turretGroupNum++
	this.turretGroupSpec[this.turretGroupNum-1].init()
	return this.turretGroupSpec[this.turretGroupNum-1]
}

func (this *EnemySpec) getMovingTurretGroupSpec() *MovingTurretGroupSpec {
	this.movingTurretGroupNum++
	this.movingTurretGroupSpec[this.movingTurretGroupNum-1].init()
	return this.movingTurretGroupSpec[this.movingTurretGroupNum-1]
}

func (this *EnemySpec) addMovingTurret(rank float32, bossMode bool /*= false*/) {
	mtn := int(rank * 0.2)
	if mtn > MOVING_TURRET_GROUP_MAX {
		mtn = MOVING_TURRET_GROUP_MAX
	}
	if mtn >= 2 {
		mtn = 1 + rand.nextInt(mtn-1)
	} else {
		mtn = 1
	}
	br := this.rank / mtn
	var moveType TurretMoveType
	if !this.bossMode {
		switch rand.nextInt(4) {
		case 0,1:
			moveType = TurretMoveTypeROLL
			break
		case 2:
			moveType = TurretMoveTypeSWING_FIX
			break
		case 3:
			moveType = TurretMoveTypeSWING_AIM
			break
		}
	} else {
		moveType = TurretMoveTypeROLL
	}
	rad := 0.9 + rand.nextFloat(0.4) - mtn*0.1
	radInc := 0.5 + rand.nextFloat(0.25)
	ad := Pi32 * 2
	var a, av, dv, s, sv float32
	switch moveType {
	case TurretMoveTypeROLL:
		a = 0.01 + rand.nextFloat(0.04)
		av = 0.01 + rand.nextFloat(0.03)
		dv = 0.01 + rand.nextFloat(0.04)
		break
	case TurretMoveTypeSWING_FIX:
		ad = Pi32/10 + rand.nextFloat(Pi32/15)
		s = 0.01 + rand.nextFloat(0.02)
		sv = 0.01 + rand.nextFloat(0.03)
		break
	case TurretMoveTypeSWING_AIM:
		ad = Pi32/10 + rand.nextFloat(Pi32/15)
		if rand.nextInt(5) == 0 {
			s = 0.01 + rand.nextFloat(0.01)
		} else {
			s = 0
		}
		sv = 0.01 + rand.nextFloat(0.02)
		break
	}
	for i := 0; i < mtn; i++ {
		tgs := this.getMovingTurretGroupSpec()
		tgs.moveType = moveType
		tgs.radiusBase = rad
		var sr float32
		switch moveType {
		case TurretMoveTypeROLL:
			tgs.alignDeg = ad
			tgs.num = 4 + rand.nextInt(6)
			if rand.nextInt(2) == 0 {
				if rand.nextInt(2) == 0 {
					tgs.setRoll(dv, 0, 0)
				} else {
					tgs.setRoll(-dv, 0, 0)
				}
			} else {
				if rand.nextInt(2) == 0 {
					tgs.setRoll(0, a, av)
				} else {
					tgs.setRoll(0, -a, av)
				}
			}
			if rand.nextInt(3) == 0 {
				tgs.setRadiusAmp(1+rand.nextFloat(1), 0.01+rand.nextFloat(0.03))
			}
			if rand.nextInt(2) == 0 {
				tgs.distRatio = 0.8 + rand.nextSignedFloat(0.3)
			}
			sr = br / tgs.num
			break
		case TurretMoveTypeSWING_FIX:
			tgs.num = 3 + rand.nextInt(5)
			tgs.alignDeg = ad * (tgs.num*0.1 + 0.3)
			if rand.nextInt(2) == 0 {
				tgs.setSwing(s, sv)
			} else {
				tgs.setSwing(-s, sv)
			}
			if rand.nextInt(6) == 0 {
				tgs.setRadiusAmp(1+rand.nextFloat(1), 0.01+rand.nextFloat(0.03))
			}
			if rand.nextInt(4) == 0 {
				tgs.setAlignAmp(0.25+rand.nextFloat(0.25), 0.01+rand.nextFloat(0.02))
			}
			sr = br / tgs.num
			sr *= 0.6
			break
		case TurretMoveTypeSWING_AIM:
			tgs.num = 3 + rand.nextInt(4)
			tgs.alignDeg = ad * (tgs.num*0.1 + 0.3)
			if rand.nextInt(2) == 0 {
				tgs.setSwing(s, sv, true)
			} else {
				tgs.setSwing(-s, sv, true)
			}
			if rand.nextInt(4) == 0 {
				tgs.setRadiusAmp(1+rand.nextFloat(1), 0.01+rand.nextFloat(0.03))
			}
			if rand.nextInt(5) == 0 {
				tgs.setAlignAmp(0.25+rand.nextFloat(0.25), 0.01+rand.nextFloat(0.02))
			}
			sr = br / tgs.num
			sr *= 0.4
			break
		}
		if rand.nextInt(4) == 0 {
			tgs.setXReverse(-1)
		}
		tgs.turretSpec.setParam(sr, TurretTypeMOVING, rand)
		if this.bossMode {
			tgs.turretSpec.setBossSpec()
		}
		rad += radInc
		ad *= 1 + rand.nextSignedFloat(0.2)
	}
}

func (this *EnemySpec) checkCollision(es EnemyState, x float32, y float32, c Collidable, shot Shot) bool {
	return es.checkCollision(x, y, c, shot)
}

func (this *EnemySpec) checkShipCollision(es EnemyState, x float32, y float32, largeOnly bool /*= false*/) bool {
	if es.destroyedCnt >= 0 || (largeOnly && this.enemyType != EnemyTypeLARGE) {
		return false
	}
	return this.shape.checkShipCollision(x-es.pos.x, y-es.pos.y, es.deg)
}

func (this *EnemySpec) move(es EnemyState) bool {
	return es.move()
}

func (this *EnemySpec) draw(es EnemyState) {
	es.draw()
}

func (this *EnemySpec) size(v float32) float32 {
	this.size = v
	if this.shape {
		this.shape.size = this.size
	}
	if this.damagedShape {
		this.damagedShape.size = this.size
	}
	if this.destroyedShape {
		this.destroyedShape.size = this.size
	}
	if this.bridgeShape {
		s := 0.9
		this.bridgeShape.size = s * (1 - this.distRatio)
	}
	return this.size
}

func (this *EnemySpec) isSmallEnemy() bool {
	return this.enemyType == EnemyTypeSMALL
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
	MoveStateSTAYIN MoveState = iotaG
	MoveStateMOVING
)

type SmallShipEnemySpec struct {
	*EnemySpec

	moveType                   int
	accel, maxSpeed, staySpeed float32
	moveDuration, stayDuration int
	speed, turnDeg             float32
}

func NewSmallShipEnemySpec(Field field, Ship ship) *SmallShipEnemySpec {
	this := SmallShipEnemySpec{NewEnemySpec(field, ship)}
	this.moveDuration = 1
	this.stayDuration = 1
	return this
}

func (this *SmallShipEnemySpec) setParam(rank float32) {
	this.set(EnemyType.SMALL)
	this.shape = NewEnemyShape(EnemyShapeTypeSMALL)
	this.damagedShape = NewEnemyShape(EnemyShapeTypeSMALL_DAMAGED)
	this.bridgeShape = NewEnemyShape(EnemyShapeTypeSMALL_BRIDGE)
	this.moveType = rand.nextInt(2)
	sr := rand.nextFloat(rank * 0.8)
	if sr > 25 {
		sr = 25
	}
	switch this.moveType {
	case MoveTypeSTOPANDGO:
		this.distRatio = 0.5
		this.size = 0.47 + rand.nextFloat(0.1)
		this.accel = 0.5 - 0.5/(2.0+rand.nextFloat(rank))
		this.maxSpeed = 0.05 * (1.0 + sr)
		this.staySpeed = 0.03
		this.moveDuration = 32 + rand.nextSignedInt(12)
		this.stayDuration = 32 + rand.nextSignedInt(12)
		break
	case MoveTypeCHASE:
		this.distRatio = 0.5
		this.size = 0.5 + rand.nextFloat(0.1)
		this.speed = 0.036 * (1.0 + sr)
		this.turnDeg = 0.02 + rand.nextSignedFloat(0.04)
		break
	}
	this.shield = 1
	tgs := this.getTurretGroupSpec()
	tgs.turretSpec.setParam(rank-sr*0.5, TurretTypeSMALL)
}

func (this *SmallShipEnemySpec) setFirstState(es EnemyState, appType int) bool {
	es.setSpec(this)
	if !es.setAppearancePos(this.field, this.ship, this.appType) {
		return false
	}
	switch this.moveType {
	case MoveTypeSTOPANDGO:
		es.speed = 0
		es.state = MoveStateMOVING
		es.cnt = this.moveDuration
		break
	case MoveTypeCHASE:
		es.speed = this.speed
		break
	}
	return true
}

func (this *SmallShipEnemySpec) move(es EnemyState) bool {
	if !super.move(es) {
		return false
	}
	switch this.moveType {
	case MoveTypeSTOPANDGO:
		es.pos.x += Sin32(es.velDeg) * es.speed
		es.pos.y += Cos32(es.velDeg) * es.speed
		es.pos.y -= this.field.lastScrollY
		if es.pos.y <= -this.field.outerSize.y {
			return false
		}
		if this.field.getBlock(es.pos) >= 0 || !this.field.checkInOuterHeightField(es.pos) {
			es.velDeg += Pi32
			es.pos.x += Sin32(es.velDeg) * es.speed * 2
			es.pos.y += Cos32(es.velDeg) * es.speed * 2
		}
		switch es.state {
		case MoveStateMOVING:
			es.speed += (maxSpeed - es.speed) * this.accel
			es.cnt--
			if es.cnt <= 0 {
				es.velDeg = rand.nextFloat(Pi32 * 2)
				es.cnt = this.stayDuration
				es.state = MoveStateSTAYING
			}
			break
		case MoveStateSTAYING:
			es.speed += (this.staySpeed - es.speed) * this.accel
			es.cnt--
			if es.cnt <= 0 {
				es.cnt = this.moveDuration
				es.state = MoveStateMOVING
			}
			break
		}
		break
	case MoveTypeCHASE:
		es.pos.x += Sin32(es.velDeg) * this.speed
		es.pos.y += Cos32(es.velDeg) * this.speed
		es.pos.y -= this.field.lastScrollY
		if es.pos.y <= -this.field.outerSize.y {
			return false
		}
		if this.field.getBlock(es.pos) >= 0 || !this.field.checkInOuterHeightField(es.pos) {
			es.velDeg += Pi32
			es.pos.x += Sin32(es.velDeg) * es.speed * 2
			es.pos.y += Cos32(es.velDeg) * es.speed * 2
		}
		var od float32
		shipPos := this.ship.nearPos(es.pos)
		if shipPos.dist(es.pos) < 0.1 {
			ad = 0
		} else {
			ad = atan2(shipPos.x-es.pos.x, shipPos.y-es.pos.y)
		}
		od = ad - es.velDeg
		normalizeDeg(od)
		if od <= this.turnDeg && od >= -this.turnDeg {
			es.velDeg = ad
		} else if od < 0 {
			es.velDeg -= this.turnDeg
		} else {
			es.velDeg += this.turnDeg
		}
		normalizeDeg(es.velDeg)
		es.cnt++
	}
	od := es.velDeg - es.deg
	normalizeDeg(od)
	es.deg += od * 0.05
	normalizeDeg(es.deg)
	if es.cnt%6 == 0 && es.speed >= 0.03 {
		this.shape.addWake(wakes, es.pos, es.deg, es.speed)
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
type ShipClass int

const (
	ShipClassMIDDLE ShipClass = iota
	LARGE
	BOSS
)

const SINK_INTERVAL = 120

type ShipEnemySpec struct {
	*EnemySpec

	speed, degVel float32
	shipClass     int
}

func NewShipEnemySpec(field Field, ship Ship) *ShipEnemySpec {
	return ShipEnemySpec{NewEnemySpec(field.ship)}
}

func (this *ShipClass) setParam(rank float32, cls int) {
	this.shipClass = cls
	this.set(EnemyTypeLARGE)
	this.shape = NewEnemyShape(EnemyShapeTypeMIDDLE)
	this.damagedShape = NewEnemyShape(EnemyShapeTypeMIDDLE_DAMAGED)
	this.destroyedShape = NewEnemyShape(EnemyShapeTypeMIDDLE_DESTROYED)
	this.bridgeShape = NewEnemyShape(EnemyShapeTypeMIDDLE_BRIDGE)
	this.distRatio = 0.7
	var mainTurretNum, subTurretNum int
	var movingTurretRatio float32
	rk := this.rank
	switch cls {
	case ShipClassMIDDLE:
		sz := 1.5 + this.rank/15 + rand.nextFloat(this.rank/15)
		ms := 2 + rand.nextFloat(0.5)
		if sz > ms {
			sz = ms
		}
		this.size = sz
		this.speed = 0.015 + rand.nextSignedFloat(0.005)
		this.degVel = 0.005 + rand.nextSignedFloat(0.003)
		switch rand.nextInt(3) {
		case 0:
			mainTurretNum = int(this.size*(1+rand.nextSignedFloat(0.25)) + 1)
			break
		case 1:
			subTurretNum = int(this.size*1.6*(1+rand.nextSignedFloat(0.5)) + 2)
			break
		case 2:
			mainTurretNum = int(this.size*(0.5+rand.nextSignedFloat(0.12)) + 1)
			movingTurretRatio = 0.5 + rand.nextFloat(0.25)
			rk = this.rank * (1 - movingTurretRatio)
			movingTurretRatio *= 2
			break
		}
		break
	case ShipClassLARGE:
		sz := 2.5 + this.rank/24 + rand.nextFloat(this.rank/24)
		ms := 3 + rand.nextFloat(1)
		if sz > ms {
			sz = ms
		}
		this.size = sz
		this.speed = 0.01 + rand.nextSignedFloat(0.005)
		this.degVel = 0.003 + rand.nextSignedFloat(0.002)
		mainTurretNum = int(this.size*(0.7+rand.nextSignedFloat(0.2)) + 1)
		subTurretNum = (this.size*1.6*(0.7+rand.nextSignedFloat(0.33)) + 2)
		movingTurretRatio = 0.25 + rand.nextFloat(0.5)
		rk = this.rank * (1 - movingTurretRatio)
		movingTurretRatio *= 3
		break
	case ShipClassBOSS:
		sz := 5 + this.rank/30 + rand.nextFloat(this.rank/30)
		ms := 9 + this.rand.nextFloat(3)
		if sz > ms {
			sz = ms
		}
		this.size = sz
		this.speed = this.ship.scrollSpeedBase + 0.0025 + rand.nextSignedFloat(0.001)
		this.degVel = 0.003 + rand.nextSignedFloat(0.002)
		mainTurretNum = int(size*0.8*(1.5+rand.nextSignedFloat(0.4)) + 2)
		subTurretNum = int(size*0.8*(2.4+rand.nextSignedFloat(0.6)) + 2)
		movingTurretRatio = 0.2 + rand.nextFloat(0.3)
		rk = this.rank * (1 - movingTurretRatio)
		movingTurretRatio *= 2.5
		break
	}
	this.shield = int(this.size * 10)
	if cls == ShipClassBOSS {
		this.shield *= 2.4
	}
	if mainTurretNum+subTurretNum <= 0 {
		tgs := this.getTurretGroupSpec()
		tgs.turretSpec.setParam(0, TurretTypeDUMMY)
	} else {
		subTurretRank := rk / (mainTurretNum*3 + subTurretNum)
		mainTurretRank := subTurretRank * 2.5
		if cls != ShipClassBOSS {
			frontMainTurretNum := int(mainTurretNum/2 + 0.99)
			rearMainTurretNum := mainTurretNum - frontMainTurretNum
			if frontMainTurretNum > 0 {
				tgs := this.getTurretGroupSpec()
				tgs.turretSpec.setParam(mainTurretRank, TurretTypeMAIN)
				tgs.num = frontMainTurretNum
				tgs.alignType = AlignTypeSTRAIGHT
				tgs.offset.y = -this.size * (0.9 + rand.nextSignedFloat(0.05))
			}
			if rearMainTurretNum > 0 {
				tgs := this.getTurretGroupSpec()
				tgs.turretSpec.setParam(mainTurretRank, TurretTypeMAIN)
				tgs.num = rearMainTurretNum
				tgs.alignType = AlignTypeSTRAIGHT
				tgs.offset.y = this.size * (0.9 + rand.nextSignedFloat(0.05))
			}
			var pts TurretSpec
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
						if rand.nextInt(2) == 0 {
							tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB)
						} else {
							tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB_DESTRUCTIVE)
						}
						pts = tgs.turretSpec
					} else {
						tgs.turretSpec.setParam(pts)
					}
					tgs.num = tn
					tgs.alignType = AlignTypeROUND
					tgs.alignDeg = ad
					ad += Pi32 / 2
					tgs.alignWidth = Pi32/6 + rand.nextFloat(Pi32/8)
					tgs.radius = this.size * 0.75
					tgs.distRatio = this.distRatio
				}
			}
		} else {
			mainTurretRank *= 2.5
			subTurretRank *= 2
			var pts TurretSpec
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
						tgs.turretSpec.setParam(pts)
					}
					tgs.num = tn
					tgs.alignType = AlignTypeROUND
					tgs.alignDeg = ad
					ad += Pi32 / 2
					tgs.alignWidth = Pi32/6 + rand.nextFloat(Pi32/8)
					tgs.radius = this.size * 0.45
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
						if rand.nextInt(2) == 0 {
							tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB)
						} else {
							tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB_DESTRUCTIVE)
						}
						pts = tgs.turretSpec
						pts.setBossSpec()
					} else {
						tgs.turretSpec.setParam(pts)
					}
					tgs.num = tn[idx]
					tgs.alignType = AlignTypeROUND
					tgs.alignDeg = ad[i]
					tgs.alignWidth = Pi32/7 + rand.nextFloat(Pi32/9)
					tgs.radius = this.size * 0.75
					tgs.distRatio = this.distRatio
				}
			}
		}
	}
	if movingTurretRatio > 0 {
		if cls == ShipClassBOSS {
			this.addMovingTurret(rank*movingTurretRatio, true)
		} else {
			this.addMovingTurret(rank * movingTurretRatio)
		}
	}
}

func (this *ShipClass) setFirstState(es EnemyState, appType int) bool {
	es.setSpec(this)
	if !es.setAppearancePos(this.field, this.ship, appType) {
		return false
	}
	es.speed = this.speed
	if es.pos.x < 0 {
		es.turnWay = -1
	} else {
		es.turnWay = 1
	}
	if isBoss {
		es.trgDeg = rand.nextFloat(0.1) + 0.1
		if rand.nextInt(2) == 0 {
			es.trgDeg *= -1
		}
		es.turnCnt = 250 + rand.nextInt(150)
	}
	return true
}

func (this *ShipClass) move(EnemyState es) bool {
	if es.destroyedCnt >= SINK_INTERVAL {
		return false
	}
	if !super.move(es) {
		return false
	}
	es.pos.x += Sin32(es.deg) * es.speed
	es.pos.y += Cos32(es.deg) * es.speed
	es.pos.y -= this.field.lastScrollY
	if es.pos.x <= -this.field.outerSize.x-this.size || es.pos.x >= this.field.outerSize.x+this.size ||
		es.pos.y <= -this.field.outerSize.y-this.size {
		return false
	}
	if es.pos.y > this.field.outerSize.y*2.2+this.size {
		es.pos.y = this.field.outerSize.y*2.2 + this.size
	}
	if isBoss {
		es.turnCnt--
		if es.turnCnt <= 0 {
			es.turnCnt = 250 + rand.nextInt(150)
			es.trgDeg = rand.nextFloat(0.1) + 0.2
			if es.pos.x > 0 {
				es.trgDeg *= -1
			}
		}
		es.deg += (es.trgDeg - es.deg) * 0.0025
		if this.ship.higherPos.y > es.pos.y {
			es.speed += (this.speed*2 - es.speed) * 0.005
		} else {
			es.speed += (this.speed - es.speed) * 0.01
		}
	} else {
		if !es.checkFrontClear() {
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
		this.shape.addWake(es.pos, es.deg, es.speed)
	}
	return true
}

func (this *ShipClass) draw(es EnemyState) {
	if es.destroyedCnt >= 0 {
		setScreenColor(
			EnemyShapeMIDDLE_COLOR_R*(1-float32(es.destroyedCnt)/SINK_INTERVAL)*0.5,
			EnemyShapeMIDDLE_COLOR_G*(1-float32(es.destroyedCnt)/SINK_INTERVAL)*0.5,
			EnemyShapeMIDDLE_COLOR_B*(1-float32(es.destroyedCnt)/SINK_INTERVAL)*0.5)
	}
	super.draw(es)
}

func (this *ShipClass) score() int {
	switch this.shipClass {
	case ShipClassMIDDLE:
		return 100
	case ShipClassLARGE:
		return 300
	case ShipClassBOSS:
		return 1000
	}
}

func (this *ShipClass) isBoss() bool {
	return this.shipClass == ShipClassBOSS
}

/**
 * Specification for a sea-based platform.
 */
type PlatformEnemySpec struct {
	*EnemySpec
}

func NewPlatformEnemySpec(field Field, ship Ship) *PlatformEnemySpec {
	return &PlatformEnemySpec{NewEnemySpec(field, ship)}
}

func (this *PlatformEnemySpec) setParam(rank float32) {
	this.set(EnemyType.PLATFORM)
	this.shape = NewEnemyShape(EnemyShapeTypePLATFORM)
	this.damagedShape = NewEnemyShape(EnemyShapeTypePLATFORM_DAMAGED)
	this.destroyedShape = NewEnemyShape(EnemyShapeTypePLATFORM_DESTROYED)
	this.bridgeShape = NewEnemyShape(EnemyShapeTypePLATFORM_BRIDGE)
	this.distRatio = 0
	this.size = 1 + this.rank/30 + rand.nextFloat(this.rank/30)
	ms := 1 + rand.nextFloat(0.25)
	if this.size > ms {
		this.size = ms
	}
	var mainTurretNum, frontTurretNum, sideTurretNum int
	rk := this.rank
	var movingTurretRatio float32
	switch rand.nextInt(3) {
	case 0:
		frontTurretNum = int(size*(2+rand.nextSignedFloat(0.5)) + 1)
		movingTurretRatio = 0.33 + rand.nextFloat(0.46)
		rk *= (1 - movingTurretRatio)
		movingTurretRatio *= 2.5
		break
	case 1:
		frontTurretNum = int(this.size*(0.5+rand.nextSignedFloat(0.2)) + 1)
		sideTurretNum = int(this.size*(0.5+rand.nextSignedFloat(0.2))+1) * 2
		break
	case 2:
		mainTurretNum = int(this.size*(1+rand.nextSignedFloat(0.33)) + 1)
		break
	}
	this.shield = int(this.size * 20)
	subTurretNum := frontTurretNum + sideTurretNum
	subTurretRank := rk / (mainTurretNum*3 + subTurretNum)
	mainTurretRank := subTurretRank * 2.5
	if mainTurretNum > 0 {
		tgs := this.getTurretGroupSpec()
		tgs.turretSpec.setParam(mainTurretRank, TurretTypeMAIN)
		tgs.num = mainTurretNum
		tgs.alignType = AlignTypeROUND
		tgs.alignDeg = 0
		tgs.alignWidth = Pi32*0.66 + rand.nextFloat(Pi32/2)
		tgs.radius = this.size * 0.7
		tgs.distRatio = this.distRatio
	}
	if frontTurretNum > 0 {
		tgs := this.getTurretGroupSpec()
		tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB)
		tgs.num = frontTurretNum
		tgs.alignType = AlignTypeROUND
		tgs.alignDeg = 0
		tgs.alignWidth = Pi32/5 + rand.nextFloat(Pi32/6)
		tgs.radius = this.size * 0.8
		tgs.distRatio = this.distRatio
	}
	sideTurretNum /= 2
	if sideTurretNum > 0 {
		var pts TurretSpec
		for i := 0; i < 2; i++ {
			tgs := this.getTurretGroupSpec()
			if i == 0 {
				tgs.turretSpec.setParam(subTurretRank, TurretTypeSUB)
				pts = tgs.turretSpec
			} else {
				tgs.turretSpec.setParam(pts)
			}
			tgs.num = sideTurretNum
			tgs.alignType = AlignTypeROUND
			tgs.alignDeg = Pi32/2 - Pi32*i
			tgs.alignWidth = Pi32/5 + rand.nextFloat(Pi32/6)
			tgs.radius = this.size * 0.75
			tgs.distRatio = this.distRatio
		}
	}
	if movingTurretRatio > 0 {
		this.addMovingTurret(this.rank * movingTurretRatio)
	}
}

func (this *PlatformEnemySpec) setFirstState(es EnemyState, x float32, y float32, d float32) bool {
	es.setSpec(this)
	es.pos.x = x
	es.pos.y = y
	es.deg = d
	es.speed = 0
	return es.checkFrontClear(true)
}

func (this *PlatformEnemySpec) move(es EnemyState) bool {
	if !super.move(es) {
		return false
	}
	es.pos.y -= this.field.lastScrollY
	return !(es.pos.y <= -this.field.outerSize.y)
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

func checkAllEnemiesShotHit(pos Vector, shape Collidable, shot Shot /*= null*/) {
	for a, _ := range actor {
		e, ok := a.(Enemy)
		if ok && e.exists {
			e.checkShotHit(pos, shape, shot)
		}
	}
}

func checkAllEnemiesHitShip(x float32, y float32, deselection Enemy /*= null*/, largeOnly bool /*= false*/) *Enemy {
	for a, _ := range actor {
		e, ok := a.(Enemy)
		if ok && e.exists && e != deselection {
			if e.checkHitShip(x, y, largeOnly) {
				return e
			}
		}
	}
	return null
}

func hasBoss() bool {
	for a, _ := range actor {
		e, ok := a.(Enemy)
		if ok && e.exists && e.isBoss() {
			return true
		}
	}
	return false
}
