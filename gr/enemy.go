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
  spec EnemySpec
  state EnemyState
}

func NewEnemy(field Field, screen Screen, ship Ship, scoreReel ScoreReel) *Enemy {
	e := new(Enemy)
	e.state = NewEnemyState (field, screen, ship, scoreReel)
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
	if (!this.spec.move(this.state)) {
		this.remove()
	}
}

func (this *Enemy) checkShotHit(p Vector, shape Collidable, shot Shot) {
	if (this.state.destroyedCnt >= 0) {
		return
	}
	if (this.spec.checkCollision(this.state, p.x, p.y, shape, shot)) {
		if (shot) {
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
  appType int
  ppos, pos Vector
  shield int
  deg, velDeg, speed, turnWay, trgDeg float32
  turnCnt, state, cnt int
  vel Vector
  turretGroup [TURRET_GROUP_MAX]turretGroup
  movingTurretGroup [MOVING_TURRET_GROUP_MAX]MovingTurretGroup
  damaged bool
  damagedCnt, destroyedCnt, explodeCnt, explodeItv, idx int
  multiplier float32
  spec EnemySpec

  field Field
  screen Screen
  ship Ship
  enemy Enemy
  stageManager StageManager
  scoreReel ScoreReel
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

func (this *EnemyState) setAppearancePos(field Field, ship Ship, appType int /*= AppearanceType.TOP*/)  bool {
	this.appType = appType
	for i := 0 ; i < 8 ; i++ {
		switch (appType) {
		case AppearanceType.TOP:
			this.pos.x = rand.nextSignedfloat32(field.size.x)
			this.pos.y = field.outerSize.y * 0.99 + this.spec.size
			if (this.pos.x < 0) {
				this.deg = Pi3232 - rand.nextfloat32(0.5)
				this.velDeg = this.deg
			} else {
				this.deg = Pi3232 + rand.nextfloat32(0.5)
				this.velDeg = this.deg
			}
			break
		case AppearanceType.SIDE:
			if (rand.nextInt(2) == 0) {
				this.pos.x = -field.outerSize.x * 0.99
				this.deg = Pi3232 / 2 + rand.nextfloat32(0.66)
				this.velDeg = this.deg
			} else {
				this.pos.x = field.outerSize.x * 0.99
				this.deg = -Pi32 / 2 - rand.nextfloat32(0.66)
				this.velDeg = this.deg
			}
			this.pos.y = field.size.y + rand.nextfloat32(field.size.y) + this.spec.size
			break
		case AppearanceType.CENTER:
			this.pos.x = 0
			this.pos.y = field.outerSize.y * 0.99 + this.spec.size
			this.deg = 0
			this.velDeg = this.deg
			break
		}
		this.ppos.x = this.pos.x
		this.ppos.y = this.pos.y
		this.vel.y = 0
		this.vel.x = 0
		this.speed = 0
		if (this.appType == AppearanceType.CENTER || this.checkFrontClear(true)) {
			return true
		}
	}
	return false
}

func (this *EnemyState) checkFrontClear( checkCurrentPos bool /*= false*/) bool {
	var si = 1
	if (this.checkCurrentPos()) {
		si = 0
	}
	for i := si; i < 5; i++ {
		cx := this.pos.x + Sin32(deg) * i * this.spec.size
		cy := this.pos.y + Cos32(deg) * i * this.spec.size
		if (this.field.getBlock(cx, cy) >= 0) {
			return false
		}
		if (checkAllEnemiesHitShip(cx, cy, enemy, true)) {
			return false
		}
	}
	return true
}

func (this *EnemyState) move() bool {
	this.ppos.x = this.pos.x
	this.ppos.y = this.pos.y
	this.multiplier -= MULTIPLIER_DECREASE_RATIO
	if (this.multiplier < 1) {
		this.multiplier = 1
	}
	if (this.destroyedCnt >= 0) {
		this.destroyedCnt++
		this.explodeCnt--
		if (this.explodeCnt < 0) {
			this.explodeItv += 2
			this.explodeItv = this.explodeItv * (1.2 + rand.nextfloat32(1))
			this.explodeCnt = this.explodeItv
			this.destroyedEdge(sqrt(this.spec.size) * 27.0 / (this.explodeItv * 0.1 + 1))
		}
	}
	this.damaged = false
	if (this.damagedCnt > 0) {
		this.damagedCnt--
	}
	alive := false
	for i := 0; i < this.spec.turretGroupNum; i++ {
		alive |= this.turretGroup[i].move(this.pos, this.deg)
	}
	for i := 0; i < this.spec.movingTurretGroupNum; i++ {
		this.movingTurretGroup[i].move(this.pos, this.deg)
	}
	if (this.destroyedCnt < 0 && !alive) {
		return this.destroyed()
	}
	return true
}

func (this *EnemyState) checkCollision(x float32, y float32, c Collidable, shot Shot) bool {
	ox := fabs32(pos.x - x), oy = fabs32(pos.y - y)
	if (ox + oy > spec.size * 2) {
		return false
	}
	for i := 0; i < spec.turretGroupNum; i++ {
		if (turretGroup[i].checkCollision(x, y, c, shot)) {
			return true
		}
	}
	if (spec.bridgeShape.checkCollision(ox, oy, c)) {
		addDamage(shot.damage, shot)
		return true
	}
	return false
}

func (this *EnemyState) increaseMultiplier(m float32) {
	multiplier += m
}

func (this *EnemyState) addScore(s int) {
	setScoreIndicator(s, 1)
}

func (this *EnemyState) addDamage(n int, shot Shot /*= null*/) {
	shield -= n
	if (shield <= 0) {
		destroyed(shot)
	} else {
		damaged = true
		damagedCnt = 7
	}
}

func (this *EnemyState) destroyed(shot Shot /*= null*/) bool {
	float32 vz
	if (shot) func (this *EnemyState) {
		explodeVel.x = Shot.SPEED * Sin32(shot.deg) / 2
		explodeVel.y = Shot.SPEED * Cos32(shot.deg) / 2
		vz = 0
	} else {
		explodeVel.x = explodeVel.y = 0
		vz = 0.05
	}
	float32 ss = spec.size * 1.5
	if (ss > 2) {
		ss = 2
	}
	float32 sn
	if (spec.size < 1) {
		sn = spec.size
	}
	else
		sn = sqrt(spec.size)
	assert(sn <>= 0)
	if (sn > 3) {
		sn = 3
	}
	for i := 0; i < sn * 8; i++ {
		Smoke s = smokes.getInstanceForced()
		s.set(pos, rand.nextSignedfloat32(0.1) + explodeVel.x, rand.nextSignedfloat32(0.1) + explodeVel.y,
					rand.nextfloat32(vz),
					Smoke.SmokeType.EXPLOSION, 32 + rand.nextInt(30), ss)
	}
	for i := 0; i < sn * 36; i++ {
		Spark sp = sparks.getInstanceForced()
		sp.set(pos, rand.nextSignedfloat32(0.8) + explodeVel.x, rand.nextSignedfloat32(0.8) + explodeVel.y,
					 0.5 + rand.nextfloat32(0.5), 0.5 + rand.nextfloat32(0.5), 0, 30 + rand.nextInt(30))
	}
	for i := 0; i < sn * 12; i++ {
		Fragment f = fragments.getInstanceForced()
		f.set(pos, rand.nextSignedfloat32(0.33) + explodeVel.x, rand.nextSignedfloat32(0.33) + explodeVel.y,
					0.05 + rand.nextfloat32(0.1),
					0.2 + rand.nextfloat32(0.33))
	}
	removeTurrets()
	int sc = spec.score
	bool r
	if (spec.type == EnemySpec.EnemyType.SMALL) {
		SoundManager.playSe("small_destroyed.wav")
		r = false
	} else {
		SoundManager.playSe("destroyed.wav")
		int bn = bullets.removeIndexedBullets(idx)
		destroyedCnt = 0
		explodeCnt = 1
		explodeItv = 3
		sc += bn * 10
		r = true
		if (spec.isBoss) {
			screen.setScreenShake(45, 0.04)
		}
	}
	setScoreIndicator(sc, multiplier)
	return r
}

func (this *EnemyState) setScoreIndicator(sc int, mp float32) {
	float32 ty = NumIndicator.getTargetY()
	if (mp > 1) {
		NumIndicator ni = numIndicators.getInstanceForced()
		ni.set(sc, NumIndicator.IndicatorType.SCORE, 0.5, pos)
		ni.addTarget(8, ty, NumIndicator.FlyingToType.RIGHT, 1, 0.5, sc, 40)
		ni.addTarget(11, ty, NumIndicator.FlyingToType.RIGHT, 0.5, 0.75,
								 cast(int) (sc * mp), 30)
		ni.addTarget(13, ty, NumIndicator.FlyingToType.RIGHT, 0.25, 1,
								 cast(int) (sc * mp * stageManager.rankMultiplier), 20)
		ni.addTarget(12, -8, NumIndicator.FlyingToType.BOTTOM, 0.5, 0.1,
								 cast(int) (sc * mp * stageManager.rankMultiplier), 40)
		ni.gotoNextTarget()
		ni = numIndicators.getInstanceForced()
		int mn = cast(int) (mp * 1000)
		ni.set(mn, NumIndicator.IndicatorType.MULTIPLIER, 0.7, pos)
		ni.addTarget(10.5, ty, NumIndicator.FlyingToType.RIGHT, 0.5, 0.2, mn, 70)
		ni.gotoNextTarget()
		ni = numIndicators.getInstanceForced()
		int rn = cast(int) (stageManager.rankMultiplier * 1000)
		ni.set(rn, NumIndicator.IndicatorType.MULTIPLIER, 0.4, 11, 8)
		ni.addTarget(13, ty, NumIndicator.FlyingToType.RIGHT, 0.5, 0.2, rn, 40)
		ni.gotoNextTarget()
		scoreReel.addActualScore(cast(int) (sc * mp * stageManager.rankMultiplier))
	} else {
		NumIndicator ni = numIndicators.getInstanceForced()
		ni.set(sc, NumIndicator.IndicatorType.SCORE, 0.3, pos)
		ni.addTarget(11, ty, NumIndicator.FlyingToType.RIGHT, 1.5, 0.2, sc, 40)
		ni.addTarget(13, ty, NumIndicator.FlyingToType.RIGHT, 0.25, 0.25,
								 cast(int) (sc * stageManager.rankMultiplier), 20)
		ni.addTarget(12, -8, NumIndicator.FlyingToType.BOTTOM, 0.5, 0.1,
								 cast(int) (sc * stageManager.rankMultiplier), 40)
		ni.gotoNextTarget()
		ni = numIndicators.getInstanceForced()
		int rn = cast(int) (stageManager.rankMultiplier * 1000)
		ni.set(rn, NumIndicator.IndicatorType.MULTIPLIER, 0.4, 11, 8)
		ni.addTarget(13, ty, NumIndicator.FlyingToType.RIGHT, 0.5, 0.2, rn, 40)
		ni.gotoNextTarget()
		scoreReel.addActualScore(cast(int) (sc * stageManager.rankMultiplier))
	}
}

func (this *EnemyState) destroyedEdge(n int) {
	SoundManager.playSe("explode.wav")
	int sn = n
	if (sn > 48) {
		sn = 48
	}
	Vector[] spp = (cast(BaseShape) spec.shape.shape).pointPos
	float32[] spd = (cast(BaseShape)spec.shape.shape).pointDeg
	int si = rand.nextInt(spp.length)
	edgePos.x = spp[si].x * spec.size + pos.x
	edgePos.y = spp[si].y * spec.size + pos.y
	float32 ss = spec.size * 0.5
	if (ss > 1) {
		ss = 1
	}
	for i := 0; i < sn; i++ {
		Smoke s = smokes.getInstanceForced()
		float32 sr = rand.nextfloat32(0.5)
		float32 sd = spd[si] + rand.nextSignedfloat32(0.2)
		assert(sd <>= 0)
		s.set(edgePos, Sin32(sd) * sr, Cos32(sd) * sr, -0.004,
					Smoke.SmokeType.EXPLOSION, 75 + rand.nextInt(25), ss)
		for j := 0; j < 2; j++ {
			Spark sp = sparks.getInstanceForced()
			sp.set(edgePos, Sin32(sd) * sr * 2, Cos32(sd) * sr * 2,
						 0.5 + rand.nextfloat32(0.5), 0.5 + rand.nextfloat32(0.5), 0, 30 + rand.nextInt(30))
		}
		if (i % 2 == 0) {
			SparkFragment sf = sparkFragments.getInstanceForced()
			sf.set(edgePos, Sin32(sd) * sr * 0.5, Cos32(sd) * sr * 0.5, 0.06 + rand.nextfloat32(0.07),
						 (0.2 + rand.nextfloat32(0.1)))
		}
	}
}

func (this *EnemyState) removeTurrets() {
	for i := 0; i < spec.turretGroupNum; i++ {
		turretGroup[i].remove()
	}
	for i := 0; i < spec.movingTurretGroupNum; i++ {
		movingTurretGroup[i].remove()
	}
}

func (this *EnemyState) draw() {
	glPushMatrix()
	if (destroyedCnt < 0 && damagedCnt > 0) {
		damagedPos.x = pos.x + rand.nextSignedfloat32(damagedCnt * 0.01)
		damagedPos.y = pos.y + rand.nextSignedfloat32(damagedCnt * 0.01)
		Screen.glTranslate(damagedPos)
	} else {
		Screen.glTranslate(pos)
	}
	glRotatef(-deg * 180 / Pi32, 0, 0, 1)
	if (destroyedCnt >= 0) {
		spec.destroyedShape.draw()
	} else if (!damaged) {
		spec.shape.draw()
	} else {
		spec.damagedShape.draw()
	}
	if (destroyedCnt < 0) {
		spec.bridgeShape.draw()
	}
	glPopMatrix()
	if (destroyedCnt >= 0) {
		return
	}
	for i := 0; i < spec.turretGroupNum; i++ {
		turretGroup[i].draw()
	}
	if (multiplier > 1) {
		float32 ox, oy
		if (multiplier < 10) {
			ox = 2.1
		} else {
			ox = 1.4
		}
		oy = 1.25
		if(spec.isBoss) {
			ox += 4
			oy -= 1.25
		}
		Letter.drawNumSign(cast(int) (multiplier * 1000),
											 pos.x + ox, pos.y + oy, 0.33, 1, 33, 3)
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
  field Field
  ship Ship
  shield int
  size float32
  distRatio float32
  turretGroupSpec [EnemyState.TURRET_GROUP_MAX]TurretGroupSpec
  turretGroupNum int
  movingTurretGroupSpec [EnemyState.MOVING_TURRET_GROUP_MAX]MovingTurretGroupSpec
  movingTurretGroupNum int
  shape, damagedShape, destroyedShape, bridgeShape EnemyShape
  enemyType int
}

this(field Field, ship Ship ) {
	this.field = field
	this.ship = ship
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.wakes = wakes
	foreach (inout TurretGroupSpec tgs; turretGroupSpec)
		tgs = new TurretGroupSpec
	foreach (inout MovingTurretGroupSpec tgs; movingTurretGroupSpec)
		tgs = new MovingTurretGroupSpec
	distRatio = 0
	shield = 1
	_size = 1
}

func (this *EnemySpec) set(enemyType int) {
	this.enemyType = enemyType
	_size = 1
	distRatio = 0
	turretGroupNum = movingTurretGroupNum = 0
}

func (this *EnemySpec) getTurretGroupSpec() *TurretGroupSpec {
	turretGroupNum++
	turretGroupSpec[turretGroupNum - 1].init()
	return turretGroupSpec[turretGroupNum - 1]
}

func (this *EnemySpec) getMovingTurretGroupSpec() *MovingTurretGroupSpec {
	movingTurretGroupNum++
	movingTurretGroupSpec[movingTurretGroupNum - 1].init()
	return movingTurretGroupSpec[movingTurretGroupNum - 1]
}

func (this *EnemySpec) addMovingTurret(rank float32, bossMode bool /*= false*/) {
	int mtn = cast(int) (rank * 0.2)
	if (mtn > EnemyState.MOVING_TURRET_GROUP_MAX) {
		mtn = EnemyState.MOVING_TURRET_GROUP_MAX
	}
	if (mtn >= 2) {
		mtn = 1 + rand.nextInt(mtn - 1)
	} else {
		mtn = 1
	}
	float32 br = rank / mtn
	int moveType
	if (!bossMode) {
		switch (rand.nextInt(4)) {
		case 0:
		case 1:
			moveType = MovingTurretGroupSpec.MoveType.ROLL
			break
		case 2:
			moveType = MovingTurretGroupSpec.MoveType.SWING_FIX
			break
		case 3:
			moveType = MovingTurretGroupSpec.MoveType.SWING_AIM
			break
		}
	} else {
		moveType = MovingTurretGroupSpec.MoveType.ROLL
	}
	float32 rad = 0.9 + rand.nextfloat32(0.4) - mtn * 0.1
	float32 radInc = 0.5 + rand.nextfloat32(0.25)
	float32 ad = Pi32 * 2
	float32 a, av, dv, s, sv
	switch (moveType) {
	case MovingTurretGroupSpec.MoveType.ROLL:
		a = 0.01 + rand.nextfloat32(0.04)
		av = 0.01 + rand.nextfloat32(0.03)
		dv = 0.01 + rand.nextfloat32(0.04)
		break
	case MovingTurretGroupSpec.MoveType.SWING_FIX:
		ad = Pi32 / 10 + rand.nextfloat32(Pi32 / 15)
		s = 0.01 + rand.nextfloat32(0.02)
		sv = 0.01 + rand.nextfloat32(0.03)
		break
	case MovingTurretGroupSpec.MoveType.SWING_AIM:
		ad = Pi32 / 10 + rand.nextfloat32(Pi32 / 15)
		if (rand.nextInt(5) == 0) {
			s = 0.01 + rand.nextfloat32(0.01)
		} else {
			s = 0
		}
		sv = 0.01 + rand.nextfloat32(0.02)
		break
	}
	for i := 0; i < mtn; i++ {
		MovingTurretGroupSpec tgs = getMovingTurretGroupSpec()
		tgs.moveType = moveType
		tgs.radiusBase = rad
		float32 sr
		switch (moveType) {
		case MovingTurretGroupSpec.MoveType.ROLL:
			tgs.alignDeg = ad
			tgs.num = 4 + rand.nextInt(6)
			if (rand.nextInt(2) == 0) {
				if (rand.nextInt(2) == 0) {
					tgs.setRoll(dv, 0, 0)
				} else {
					tgs.setRoll(-dv, 0, 0)
				}
			} else {
				if (rand.nextInt(2) == 0) {
					tgs.setRoll(0, a, av)
				} else {
					tgs.setRoll(0, -a, av)
				}
			}
			if (rand.nextInt(3) == 0) {
				tgs.setRadiusAmp(1 + rand.nextfloat32(1), 0.01 + rand.nextfloat32(0.03))
			}
			if (rand.nextInt(2) == 0) {
				tgs.distRatio = 0.8 + rand.nextSignedfloat32(0.3)
			}
			sr = br / tgs.num
			break
		case MovingTurretGroupSpec.MoveType.SWING_FIX:
			tgs.num = 3 + rand.nextInt(5)
			tgs.alignDeg = ad * (tgs.num * 0.1 + 0.3)
			if (rand.nextInt(2) == 0) {
				tgs.setSwing(s, sv)
			} else {
				tgs.setSwing(-s, sv)
			}
			if (rand.nextInt(6) == 0) {
				tgs.setRadiusAmp(1 + rand.nextfloat32(1), 0.01 + rand.nextfloat32(0.03))
			}
			if (rand.nextInt(4) == 0) {
				tgs.setAlignAmp(0.25 + rand.nextfloat32(0.25), 0.01 + rand.nextfloat32(0.02))
			}
			sr = br / tgs.num
			sr *= 0.6
			break
		case MovingTurretGroupSpec.MoveType.SWING_AIM:
			tgs.num = 3 + rand.nextInt(4)
			tgs.alignDeg = ad * (tgs.num * 0.1 + 0.3)
			if (rand.nextInt(2) == 0) {
				tgs.setSwing(s, sv, true)
			} else {
				tgs.setSwing(-s, sv, true)
			}
			if (rand.nextInt(4) == 0) {
				tgs.setRadiusAmp(1 + rand.nextfloat32(1), 0.01 + rand.nextfloat32(0.03))
			}
			if (rand.nextInt(5) == 0) {
				tgs.setAlignAmp(0.25 + rand.nextfloat32(0.25), 0.01 + rand.nextfloat32(0.02))
			}
			sr = br / tgs.num
			sr *= 0.4
			break
		}
		if (rand.nextInt(4) == 0) {
			tgs.setXReverse(-1)
		}
		tgs.turretSpec.setParam(sr, TurretSpec.TurretType.MOVING, rand)
		if (bossMode) {
			tgs.turretSpec.setBossSpec()
		}
		rad += radInc
		ad *= 1 + rand.nextSignedfloat32(0.2)
	}
}

func (this *EnemySpec) checkCollision(es EnemyState, x float32, y float32, c Collidable, shot Shot) bool {
	return es.checkCollision(x, y, c, shot)
}

func (this *EnemySpec) checkShipCollision(es EnemyState, x float32, y float32 y,largeOnly bool /*= false*/) bool {
	if (es.destroyedCnt >= 0 || (largeOnly && enemyType != EnemyType.LARGE)) {
		return false
	}
	return shape.checkShipCollision(x - es.pos.x, y - es.pos.y, es.deg); 
}

func (this *EnemySpec) move(es EnemyState) bool {
	return es.move()
}

func (this *EnemySpec) draw(es EnemyState) {
	es.draw()
}

func (this *EnemySpec) size(v float32) float32 {
	_size = v
	if (shape) {
		shape.size = _size
	}
	if (damagedShape) {
		damagedShape.size = _size
	}
	if (destroyedShape) {
		destroyedShape.size = _size
	}
	if (bridgeShape) {
		float32 s = 0.9
		bridgeShape.size = s * (1 - distRatio)
	}
	return _size
}

func (this *EnemySpec) isSmallEnemy() bool {
	return enemyType == EnemyType.SMALL
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

  moveType int
  accel, maxSpeed, staySpeed float32
  moveDuration, stayDuration int
  speed, turnDeg float32
}

this(Field field, Ship ship){
	super(field, ship)
	moveDuration = stayDuration = 1
}

func (this *SmallShipEnemySpec) setParam(rank float32) {
	set(EnemyType.SMALL)
	shape = new EnemyShape(EnemyShape.EnemyShapeType.SMALL)
	damagedShape = new EnemyShape(EnemyShape.EnemyShapeType.SMALL_DAMAGED)
	bridgeShape = new EnemyShape(EnemyShape.EnemyShapeType.SMALL_BRIDGE)
	moveType = rand.nextInt(2)
	float32 sr = rand.nextfloat32(rank * 0.8)
	if (sr > 25) {
		sr = 25
	}
	switch (moveType) {
	case MoveType.STOPANDGO:
		distRatio = 0.5
		size = 0.47 + rand.nextfloat32(0.1)
		accel = 0.5 - 0.5 / (2.0 + rand.nextfloat32(rank))
		maxSpeed = 0.05 * (1.0 + sr)
		staySpeed = 0.03
		moveDuration = 32 + rand.nextSignedInt(12)
		stayDuration = 32 + rand.nextSignedInt(12)
		break
	case MoveType.CHASE:
		distRatio = 0.5
		size = 0.5 + rand.nextfloat32(0.1)
		speed = 0.036 * (1.0 + sr)
		turnDeg = 0.02 + rand.nextSignedfloat32(0.04)
		break
	}
	shield = 1
	TurretGroupSpec tgs = getTurretGroupSpec()
	tgs.turretSpec.setParam(rank - sr * 0.5, TurretSpec.TurretType.SMALL, rand)
}

func (this *SmallShipEnemySpec) setFirstState(es EnemyState, appType int) bool {
	es.setSpec(this)
	if (!es.setAppearancePos(field, ship, rand, appType)) {
		return false
	}
	switch (moveType) {
	case MoveType.STOPANDGO:
		es.speed = 0
		es.state = MoveState.MOVING
		es.cnt = moveDuration
		break
	case MoveType.CHASE:
		es.speed = speed
		break
	}
	return true
}

func (this *SmallShipEnemySpec) move(es EnemyState) bool {
	if (!super.move(es)) {
		return false
	}
	switch (moveType) {
	case MoveType.STOPANDGO:
		es.pos.x += Sin32(es.velDeg) * es.speed
		es.pos.y += Cos32(es.velDeg) * es.speed
		es.pos.y -= field.lastScrollY
		if  (es.pos.y <= -field.outerSize.y) {
			return false
		}
		if (field.getBlock(es.pos) >= 0 || !field.checkInOuterHeightField(es.pos)) {
			es.velDeg += Pi32
			es.pos.x += Sin32(es.velDeg) * es.speed * 2
			es.pos.y += Cos32(es.velDeg) * es.speed * 2
		}
		switch (es.state) {
		case MoveState.MOVING:
			es.speed += (maxSpeed - es.speed) * accel
			es.cnt--
			if (es.cnt <= 0) {
				es.velDeg = rand.nextfloat32(Pi32 * 2)
				es.cnt = stayDuration
				es.state = MoveState.STAYING
			}
			break
		case MoveState.STAYING:
			es.speed += (staySpeed - es.speed) * accel
			es.cnt--
			if (es.cnt <= 0) {
				es.cnt = moveDuration
				es.state = MoveState.MOVING
			}
			break
		}
		break
	case MoveType.CHASE:
		es.pos.x += Sin32(es.velDeg) * speed
		es.pos.y += Cos32(es.velDeg) * speed
		es.pos.y -= field.lastScrollY
		if  (es.pos.y <= -field.outerSize.y) {
			return false
		}
		if (field.getBlock(es.pos) >= 0 || !field.checkInOuterHeightField(es.pos)) {
			es.velDeg += Pi32
			es.pos.x += Sin32(es.velDeg) * es.speed * 2
			es.pos.y += Cos32(es.velDeg) * es.speed * 2
		}
		float32 ad
		Vector shipPos = ship.nearPos(es.pos)
		if (shipPos.dist(es.pos) < 0.1) {
			ad = 0
		} else {
			ad = atan2(shipPos.x - es.pos.x, shipPos.y - es.pos.y)
		}
		assert(ad <>= 0)
		float32 od = ad - es.velDeg
		Math.normalizeDeg(od)
		if (od <= turnDeg && od >= -turnDeg) {
			es.velDeg = ad
		} else if (od < 0) {
			es.velDeg -= turnDeg
		} else {
			es.velDeg += turnDeg
		}
		Math.normalizeDeg(es.velDeg)
		es.cnt++
	}
	float32 od = es.velDeg - es.deg
	Math.normalizeDeg(od)
	es.deg += od * 0.05
	Math.normalizeDeg(es.deg)
	if (es.cnt % 6 == 0 && es.speed >= 0.03) {
		shape.addWake(wakes, es.pos, es.deg, es.speed)
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

type ShipEnemySpec struct{
	*EnemySpec

  speed, degVel float32
  shipClass int
}


this(Field field, Ship ship ) {
	super(field, ship)
}

func (this *ShipClass) setParam(rank float32, cls int) {
	shipClass = cls
	set(EnemyType.LARGE)
	shape = new EnemyShape(EnemyShape.EnemyShapeType.MIDDLE)
	damagedShape = new EnemyShape(EnemyShape.EnemyShapeType.MIDDLE_DAMAGED)
	destroyedShape = new EnemyShape(EnemyShape.EnemyShapeType.MIDDLE_DESTROYED)
	bridgeShape = new EnemyShape(EnemyShape.EnemyShapeType.MIDDLE_BRIDGE)
	distRatio = 0.7
	int mainTurretNum = 0, subTurretNum = 0
	float32 movingTurretRatio = 0
	float32 rk = rank
	switch (cls) {
	case ShipClass.MIDDLE:
		float32 sz = 1.5 + rank / 15 + rand.nextfloat32(rank / 15)
		float32 ms = 2 + rand.nextfloat32(0.5)
		if (sz > ms) {
			sz = ms
		}
		size = sz
		speed = 0.015 + rand.nextSignedfloat32(0.005)
		degVel = 0.005 + rand.nextSignedfloat32(0.003)
		switch (rand.nextInt(3)) {
		case 0:
			mainTurretNum = cast(int) (size * (1 + rand.nextSignedfloat32(0.25)) + 1)
			break
		case 1:
			subTurretNum = cast(int) (size * 1.6 * (1 + rand.nextSignedfloat32(0.5)) + 2)
			break
		case 2:
			mainTurretNum = cast(int) (size * (0.5 + rand.nextSignedfloat32(0.12)) + 1)
			movingTurretRatio = 0.5 + rand.nextfloat32(0.25)
			rk = rank * (1 - movingTurretRatio)
			movingTurretRatio *= 2
			break
		}
		break
	case ShipClass.LARGE:
		float32 sz = 2.5 + rank / 24 + rand.nextfloat32(rank / 24)
		float32 ms = 3 + rand.nextfloat32(1)
		if (sz > ms) {
			sz = ms
		}
		size = sz
		speed = 0.01 + rand.nextSignedfloat32(0.005)
		degVel = 0.003 + rand.nextSignedfloat32(0.002)
		mainTurretNum = cast(int) (size * (0.7 + rand.nextSignedfloat32(0.2)) + 1)
		subTurretNum = cast(int) (size * 1.6 * (0.7 + rand.nextSignedfloat32(0.33)) + 2)
		movingTurretRatio = 0.25 + rand.nextfloat32(0.5)
		rk = rank * (1 - movingTurretRatio)
		movingTurretRatio *= 3
		break
	case ShipClass.BOSS:
		float32 sz = 5 + rank / 30 + rand.nextfloat32(rank / 30)
		float32 ms = 9 + rand.nextfloat32(3)
		if (sz > ms) {
			sz = ms
		}
		size = sz
		speed = ship.scrollSpeedBase + 0.0025 + rand.nextSignedfloat32(0.001)
		degVel = 0.003 + rand.nextSignedfloat32(0.002)
		mainTurretNum = cast(int) (size * 0.8 * (1.5 + rand.nextSignedfloat32(0.4)) + 2)
		subTurretNum = cast(int) (size * 0.8 * (2.4 + rand.nextSignedfloat32(0.6)) + 2)
		movingTurretRatio = 0.2 + rand.nextfloat32(0.3)
		rk = rank * (1 - movingTurretRatio)
		movingTurretRatio *= 2.5
		break
	}
	shield = cast(int) (size * 10)
	if (cls == ShipClass.BOSS) {
		shield *= 2.4
	}
	if (mainTurretNum + subTurretNum <= 0) {
		TurretGroupSpec tgs = getTurretGroupSpec()
		tgs.turretSpec.setParam(0, TurretSpec.TurretType.DUMMY, rand)
	} else {
		float32 subTurretRank = rk / (mainTurretNum * 3 + subTurretNum)
		float32 mainTurretRank = subTurretRank * 2.5
		if (cls != ShipClass.BOSS) {
			int frontMainTurretNum = cast(int) (mainTurretNum / 2 + 0.99)
			int rearMainTurretNum = mainTurretNum - frontMainTurretNum
			if (frontMainTurretNum > 0) {
				TurretGroupSpec tgs = getTurretGroupSpec()
				tgs.turretSpec.setParam(mainTurretRank, TurretSpec.TurretType.MAIN, rand)
				tgs.num = frontMainTurretNum
				tgs.alignType = TurretGroupSpec.AlignType.STRAIGHT
				tgs.offset.y = -size * (0.9 + rand.nextSignedfloat32(0.05))
			}
			if (rearMainTurretNum > 0) {
				TurretGroupSpec tgs = getTurretGroupSpec()
				tgs.turretSpec.setParam(mainTurretRank, TurretSpec.TurretType.MAIN, rand)
				tgs.num = rearMainTurretNum
				tgs.alignType = TurretGroupSpec.AlignType.STRAIGHT
				tgs.offset.y = size * (0.9 + rand.nextSignedfloat32(0.05))
			} 
			TurretSpec pts
			if (subTurretNum > 0) {
				int frontSubTurretNum = (subTurretNum + 2) / 4
				int rearSubTurretNum = (subTurretNum - frontSubTurretNum * 2) / 2
				int tn = frontSubTurretNum
				float32 ad = -Pi32 / 4
				for i := 0; i < 4; i++ {
					if (i == 2) {
						tn = rearSubTurretNum
					}
					if (tn <= 0) {
						continue
					}
					TurretGroupSpec tgs = getTurretGroupSpec()
					if (i == 0 || i == 2) {
						if (rand.nextInt(2) == 0) {
							tgs.turretSpec.setParam(subTurretRank, TurretSpec.TurretType.SUB, rand)
						}
						else
							tgs.turretSpec.setParam(subTurretRank, TurretSpec.TurretType.SUB_DESTRUCTIVE, rand)
						pts = tgs.turretSpec
					} else {
						tgs.turretSpec.setParam(pts)
					}
					tgs.num = tn
					tgs.alignType = TurretGroupSpec.AlignType.ROUND
					tgs.alignDeg = ad
					ad += Pi32 / 2
					tgs.alignWidth = Pi32 / 6 + rand.nextfloat32(Pi32 / 8)
					tgs.radius = size * 0.75
					tgs.distRatio = distRatio
				}
			}
		} else {
			mainTurretRank *= 2.5
			subTurretRank *= 2
			TurretSpec pts
			if (mainTurretNum > 0) {
				int frontMainTurretNum = (mainTurretNum + 2) / 4
				int rearMainTurretNum = (mainTurretNum - frontMainTurretNum * 2) / 2
				int tn = frontMainTurretNum
				float32 ad = -Pi32 / 4
				for i := 0; i < 4; i++ {
					if (i == 2) {
						tn = rearMainTurretNum
					}
					if (tn <= 0) {
						continue
					}
					TurretGroupSpec tgs = getTurretGroupSpec()
					if (i == 0 || i == 2) {
						tgs.turretSpec.setParam(mainTurretRank, TurretSpec.TurretType.MAIN, rand)
						pts = tgs.turretSpec
						pts.setBossSpec()
					} else {
						tgs.turretSpec.setParam(pts)
					}
					tgs.num = tn
					tgs.alignType = TurretGroupSpec.AlignType.ROUND
					tgs.alignDeg = ad
					ad += Pi32 / 2
					tgs.alignWidth = Pi32 / 6 + rand.nextfloat32(Pi32 / 8)
					tgs.radius = size * 0.45
					tgs.distRatio = distRatio
				}
			}
			if (subTurretNum > 0) {
				int[3] tn
				tn[0] = (subTurretNum + 2) / 6
				tn[1] = (subTurretNum - tn[0] * 2) / 4
				tn[2] = (subTurretNum - tn[0] * 2 - tn[1] * 2) / 2
				static const float32[] ad = [Pi32 / 4, -Pi32 / 4, Pi32 / 2, -Pi32 / 2, Pi32 / 4 * 3, -Pi32 / 4 * 3]
				for i := 0; i < 6; i++ {
					int idx = i / 2
					if (tn[idx] <= 0) {
						continue
					}
					TurretGroupSpec tgs = getTurretGroupSpec()
					if (i == 0 || i == 2 || i == 4) {
						if (rand.nextInt(2) == 0) {
							tgs.turretSpec.setParam(subTurretRank, TurretSpec.TurretType.SUB, rand)
						} else {
							tgs.turretSpec.setParam(subTurretRank, TurretSpec.TurretType.SUB_DESTRUCTIVE, rand)
						}
						pts = tgs.turretSpec
						pts.setBossSpec()
					} else {
						tgs.turretSpec.setParam(pts)
					}
					tgs.num = tn[idx]
					tgs.alignType = TurretGroupSpec.AlignType.ROUND
					tgs.alignDeg = ad[i]
					tgs.alignWidth = Pi32 / 7 + rand.nextfloat32(Pi32 / 9)
					tgs.radius = size * 0.75
					tgs.distRatio = distRatio
				}
			}
		}
	}
	if (movingTurretRatio > 0) {
		if (cls == ShipClass.BOSS) {
			addMovingTurret(rank * movingTurretRatio, true)
		} else {
			addMovingTurret(rank * movingTurretRatio)
		}
	}
}

func (this *ShipClass) setFirstState(es EnemyState, appType int) bool {
	es.setSpec(this)
	if (!es.setAppearancePos(field, ship, rand, appType)) {
		return false
	}
	es.speed = speed
	if (es.pos.x < 0) {
		es.turnWay = -1
	} else {
		es.turnWay = 1
	}
	if (isBoss) {
		es.trgDeg = rand.nextfloat32(0.1) + 0.1
		if (rand.nextInt(2) == 0) {
			es.trgDeg *= -1
		}
		es.turnCnt = 250 + rand.nextInt(150)
	}
	return true
}

func (this *ShipClass) move(EnemyState es) bool {
	if (es.destroyedCnt >= SINK_INTERVAL) {
		return false
	}
	if (!super.move(es)) {
		return false
	}
	es.pos.x += Sin32(es.deg) * es.speed
	es.pos.y += Cos32(es.deg) * es.speed
	es.pos.y -= field.lastScrollY
	if  (es.pos.x <= -field.outerSize.x - size || es.pos.x >= field.outerSize.x + size ||
			 es.pos.y <= -field.outerSize.y - size) {
		return false
	}
	if (es.pos.y > field.outerSize.y * 2.2 + size) {
		es.pos.y = field.outerSize.y * 2.2 + size
	}
	if (isBoss) {
		es.turnCnt--
		if (es.turnCnt <= 0) {
			es.turnCnt = 250 + rand.nextInt(150)
			es.trgDeg = rand.nextfloat32(0.1) + 0.2
			if (es.pos.x > 0) {
				es.trgDeg *= -1
			}
		}
		es.deg += (es.trgDeg - es.deg) * 0.0025
		if (ship.higherPos.y > es.pos.y) {
			es.speed += (speed * 2 - es.speed) * 0.005
		} else {
			es.speed += (speed - es.speed) * 0.01
		}
	} else {
		if (!es.checkFrontClear()) {
			es.deg += degVel * es.turnWay
			es.speed *= 0.98
		} else {
			if (es.destroyedCnt < 0) {
				es.speed += (speed - es.speed) * 0.01
			} else {
				es.speed *= 0.98
			}
		}
	}
	es.cnt++
	if (es.cnt % 6 == 0 && es.speed >= 0.01 && es.destroyedCnt < SINK_INTERVAL / 2) {
		shape.addWake(wakes, es.pos, es.deg, es.speed)
	}
	return true
}

func (this *ShipClass) draw(es EnemyState) {
	if (es.destroyedCnt >= 0) {
		Screen.setColor(
			EnemyShape.MIDDLE_COLOR_R * (1 - cast(float32) es.destroyedCnt / SINK_INTERVAL) * 0.5,
			EnemyShape.MIDDLE_COLOR_G * (1 - cast(float32) es.destroyedCnt / SINK_INTERVAL) * 0.5,
			EnemyShape.MIDDLE_COLOR_B * (1 - cast(float32) es.destroyedCnt / SINK_INTERVAL) * 0.5)
	}
	super.draw(es)
}

func (this *ShipClass) score() int {
	switch (shipClass) {
	case ShipClass.MIDDLE:
		return 100
	case ShipClass.LARGE:
		return 300
	case ShipClass.BOSS:
		return 1000
	}
}

func (this *ShipClass) isBoss() bool {
	return shipClass == ShipClass.BOSS
}

/**
 * Specification for a sea-based platform.
 */
type PlatformEnemySpec struct{
	*EnemySpec
}

this(field Field, ship Ship ) {
	super(field, ship)
}

func (this *PlatformEnemySpec) setParam(rank float32) {
	set(EnemyType.PLATFORM)
	shape = new EnemyShape(EnemyShape.EnemyShapeType.PLATFORM)
	damagedShape = new EnemyShape(EnemyShape.EnemyShapeType.PLATFORM_DAMAGED)
	destroyedShape = new EnemyShape(EnemyShape.EnemyShapeType.PLATFORM_DESTROYED)
	bridgeShape = new EnemyShape(EnemyShape.EnemyShapeType.PLATFORM_BRIDGE)
	distRatio = 0
	size = 1 + rank / 30 + rand.nextfloat32(rank / 30)
	float32 ms = 1 + rand.nextfloat32(0.25)
	if (size > ms) {
		size = ms
	}
	int mainTurretNum = 0, frontTurretNum = 0, sideTurretNum = 0
	float32 rk = rank
	float32 movingTurretRatio = 0
	switch (rand.nextInt(3)) {
	case 0:
		frontTurretNum = cast(int) (size * (2 + rand.nextSignedfloat32(0.5)) + 1)
		movingTurretRatio = 0.33 + rand.nextfloat32(0.46)
		rk *= (1 - movingTurretRatio)
		movingTurretRatio *= 2.5
		break
	case 1:
		frontTurretNum = cast(int) (size * (0.5 + rand.nextSignedfloat32(0.2)) + 1)
		sideTurretNum = cast(int) (size * (0.5 + rand.nextSignedfloat32(0.2)) + 1) * 2
		break
	case 2:
		mainTurretNum = cast(int) (size * (1 + rand.nextSignedfloat32(0.33)) + 1)
		break
	}
	shield = cast(int) (size * 20)
	int subTurretNum = frontTurretNum + sideTurretNum
	float32 subTurretRank = rk / (mainTurretNum * 3 + subTurretNum)
	float32 mainTurretRank = subTurretRank * 2.5
	if (mainTurretNum > 0) {
		TurretGroupSpec tgs = getTurretGroupSpec()
		tgs.turretSpec.setParam(mainTurretRank, TurretSpec.TurretType.MAIN, rand)
		tgs.num = mainTurretNum
		tgs.alignType = TurretGroupSpec.AlignType.ROUND
		tgs.alignDeg = 0
		tgs.alignWidth = Pi32 * 0.66 + rand.nextfloat32(Pi32 / 2)
		tgs.radius = size * 0.7
		tgs.distRatio = distRatio
	}
	if (frontTurretNum > 0) {
		TurretGroupSpec tgs = getTurretGroupSpec()
		tgs.turretSpec.setParam(subTurretRank, TurretSpec.TurretType.SUB, rand)
		tgs.num = frontTurretNum
		tgs.alignType = TurretGroupSpec.AlignType.ROUND
		tgs.alignDeg = 0
		tgs.alignWidth = Pi32 / 5 + rand.nextfloat32(Pi32 / 6)
		tgs.radius = size * 0.8
		tgs.distRatio = distRatio
	}
	sideTurretNum /= 2
	if (sideTurretNum > 0) {
		TurretSpec pts
		for i := 0; i < 2; i++ {
			TurretGroupSpec tgs = getTurretGroupSpec()
			if (i == 0) {
				tgs.turretSpec.setParam(subTurretRank, TurretSpec.TurretType.SUB, rand)
				pts = tgs.turretSpec
			} else {
				tgs.turretSpec.setParam(pts)
			}
			tgs.num = sideTurretNum
			tgs.alignType = TurretGroupSpec.AlignType.ROUND
			tgs.alignDeg = Pi32 / 2 - Pi32 * i
			tgs.alignWidth = Pi32 / 5 + rand.nextfloat32(Pi32 / 6)
			tgs.radius = size * 0.75
			tgs.distRatio = distRatio
		}
	}
	if (movingTurretRatio > 0) {
		addMovingTurret(rank * movingTurretRatio)
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
	if (!super.move(es)) {
		return false
	}
	es.pos.y -= field.lastScrollY
	return ! ( es.pos.y <= -field.outerSize.y) 
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
	for a,_ := range actor {
		e, ok := a.(Enemy)
		if (ok && e.exists) {
			e.checkShotHit(pos, shape, shot)
		}
	}
}

func checkAllEnemiesHitShip(x float32, y float32, deselection Enemy /*= null*/, largeOnly bool /*= false*/) *Enemy {
	for a,_ := range actor {
		e, ok := a.(Enemy)
		if (ok && e.exists && e != deselection) {
			if (e.checkHitShip(x, y, largeOnly)) {
				return e
			}
		}
	}
	return null
}

func hasBoss() bool {
	foreach (Enemy e; actor) {
		if (e.exists && e.isBoss) {
			return true
		}
	}
	return false
}
