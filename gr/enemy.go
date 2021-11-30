package main

import (
	"math"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/actor"
	"github.com/dragonfax/gunroar/gr/letter"
	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v2.1/gl"
)

/**
 * Enemy ships.
 */

var _ actor.Actor = &Enemy{}

type Enemy struct {
	actor.ExistsImpl

	spec   EnemySpec
	_state EnemyState
}

func NewEnemy() *Enemy {
	this := &Enemy{}
	return this
}

func (this *Enemy) Init(args []interface{}) {
	this._state = NewEnemyState(
		args[0].(*Field), args[1].(*Screen),
		args[2].(*BulletPool), args[3].(*Ship),
		args[4].(*SparkPool), args[5].(*SmokePool),
		args[6].(*FragmentPool), args[7].(*SparkFragmentPool),
		args[8].(*NumIndicatorPool), args[9].(*ScoreReel))
}

func (this *Enemy) setEnemyPool(enemies *EnemyPool) {
	this._state.setEnemyAndPool(this, enemies)
}

func (this *Enemy) setStageManager(stageManager *StageManager) {
	this._state.setStageManager(stageManager)
}

func (this *Enemy) set(spec EnemySpec) {
	this.spec = spec
	this.SetExists(true)
}

func (this *Enemy) Move() {
	if !this.spec.move(this.state()) {
		this.remove()
	}
}

func (this *Enemy) checkShotHit(p vector.Vector, shape sdl.Collidable, shot *Shot) {
	if this._state.destroyedCnt >= 0 {
		return
	}
	if this.spec.checkCollision(this._state, p.X, p.Y, shape, shot) {
		if shot != nil {
			shot.removeHitToEnemy(this.spec.isSmallEnemy())
		}
	}
}

func (this *Enemy) checkHitShip(x float64, y float64, largeOnly bool /* = false */) bool {
	return this.spec.checkShipCollision(this._state, x, y, largeOnly)
}

/* unused
func (this *Enemy) addDamage(n int) {
	this._state.addDamage(n, nil)
}
*/

func (this *Enemy) increaseMultiplier(m float64) {
	this._state.increaseMultiplier(m)
}

func (this *Enemy) addScore(s int) {
	this._state.addScore(s)
}

func (this *Enemy) remove() {
	this._state.removeTurrets()
	this.SetExists(false)
}

func (this *Enemy) Draw() {
	this.spec.draw(this._state)
}

func (this *Enemy) state() EnemyState {
	return this._state
}

func (this *Enemy) pos() vector.Vector {
	return this._state.pos
}

func (this *Enemy) size() float64 {
	return this.spec.Size()
}

func (this *Enemy) index() int {
	return this._state.idx
}

func (this *Enemy) isBoss() bool {
	return this.spec.isBoss()
}

/**
 * Enemy status (position, direction, velocity, turrets, etc).
 */

type AppearanceType int

const (
	TOP AppearanceType = iota
	SIDE
	CENTER
)

const TURRET_GROUP_MAX = 10
const MOVING_TURRET_GROUP_MAX = 4
const MULTIPLIER_DECREASE_RATIO = 0.005

var enemyStateRand = r.New(r.NewSource(time.Now().Unix()))
var edgePos, explodeVel, enemyDamagedPos vector.Vector
var idxCount = 0

func setEnemyStateRandSeed(seed int64) {
	enemyStateRand = r.New(r.NewSource(seed))
}

type EnemyState struct {
	appType                                               AppearanceType
	pos, ppos                                             vector.Vector
	shield                                                int
	deg, velDeg, speed, turnWay, trgDeg                   float64
	turnCnt, cnt                                          int
	state                                                 MoveState
	vel                                                   vector.Vector
	turretGroup                                           [TURRET_GROUP_MAX]TurretGroup
	movingTurretGroup                                     [MOVING_TURRET_GROUP_MAX]MovingTurretGroup
	damaged                                               bool
	damagedCnt, destroyedCnt, explodeCnt, explodeItv, idx int
	multiplier                                            float64
	spec                                                  EnemySpec

	field          *Field
	screen         *Screen
	bullets        *BulletPool
	ship           *Ship
	sparks         *SparkPool
	smokes         *SmokePool
	fragments      *FragmentPool
	sparkFragments *SparkFragmentPool
	numIndicators  *NumIndicatorPool
	enemy          *Enemy
	enemies        *EnemyPool
	stageManager   *StageManager
	scoreReel      *ScoreReel
}

func NewEnemyState(field *Field, screen *Screen, bullets *BulletPool, ship *Ship,
	sparks *SparkPool, smokes *SmokePool,
	fragments *FragmentPool, sparkFragments *SparkFragmentPool,
	numIndicators *NumIndicatorPool, scoreReel *ScoreReel) EnemyState {
	this := EnemyState{}
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
	this.turnWay = 1
	this.explodeItv = 1
	this.multiplier = 1
	return this
}

func (this *EnemyState) setEnemyAndPool(enemy *Enemy, enemies *EnemyPool) {
	this.enemy = enemy
	this.enemies = enemies
	for i := range this.turretGroup {
		this.turretGroup[i] = NewTurretGroup(this.field, this.bullets, this.ship, this.sparks, this.smokes, this.fragments, this.enemy)
	}
	for i := range this.movingTurretGroup {
		this.movingTurretGroup[i] = NewMovingTurretGroup(this.field, this.bullets, this.ship, this.sparks, this.smokes, this.fragments, this.enemy)
	}
}

func (this *EnemyState) setStageManager(stageManager *StageManager) {
	this.stageManager = stageManager
}

func (this *EnemyState) setSpec(spec EnemySpec) {
	this.spec = spec
	this.shield = spec.Shield()
	for i := 0; i < spec.TurretGroupNum(); i++ {
		this.turretGroup[i].set(spec.TurretGroupSpec(i))
	}
	for i := 0; i < spec.MovingTurretGroupNum(); i++ {
		this.movingTurretGroup[i].set(spec.MovingTurretGroupSpec(i))
	}
	this.cnt = 0
	this.damaged = false
	this.damagedCnt = 0
	this.destroyedCnt = -1
	this.explodeCnt = 0
	this.explodeItv = 1
	this.multiplier = 1
}

func nextFloat(rand *r.Rand, n float64) float64 {
	return rand.Float64() * n
}

func (this *EnemyState) setAppearancePos(field *Field, ship *Ship, rand *r.Rand, appType AppearanceType /* = TOP */) bool {
	this.appType = appType
	for i := 0; i < 8; i++ {
		switch appType {
		case TOP:
			this.pos.X = nextSignedFloat(rand, field.size().X)
			this.pos.Y = field.outerSize().Y*0.99 + this.spec.Size()
			if this.pos.X < 0 {
				this.velDeg = math.Pi - nextFloat(rand, 0.5)
				this.deg = math.Pi - nextFloat(rand, 0.5)
			} else {
				this.velDeg = math.Pi + nextFloat(rand, 0.5)
				this.deg = math.Pi + nextFloat(rand, 0.5)
			}
		case SIDE:
			if rand.Intn(2) == 0 {
				this.pos.X = -field.outerSize().X * 0.99
				this.velDeg = math.Pi/2 + nextFloat(rand, 0.66)
				this.deg = this.velDeg
			} else {
				this.pos.X = field.outerSize().X * 0.99
				this.velDeg = -math.Pi/2 - nextFloat(rand, 0.66)
				this.deg = this.velDeg
			}
			this.pos.Y = field.size().Y + nextFloat(rand, field.size().Y) + this.spec.Size()
		case CENTER:
			this.pos.X = 0
			this.pos.Y = field.outerSize().Y*0.99 + this.spec.Size()
			this.velDeg = 0
			this.deg = 0
		}
		this.ppos.X = this.pos.X
		this.ppos.Y = this.pos.Y
		this.vel.X = 0
		this.vel.Y = 0
		this.speed = 0
		if appType == CENTER || this.checkFrontClear(true) {
			return true
		}
	}
	return false
}

func (this *EnemyState) checkFrontClear(checkCurrentPos bool /* = false */) bool {
	si := 1
	if checkCurrentPos {
		si = 0
	}
	for i := si; i < 5; i++ {
		cx := this.pos.X + math.Sin(this.deg)*float64(i)*this.spec.Size()
		cy := this.pos.Y + math.Cos(this.deg)*float64(i)*this.spec.Size()
		if this.field.getBlock(cx, cy) >= 0 {
			return false
		}
		if this.enemies.checkHitShip(cx, cy, this.enemy, true) != nil {
			return false
		}
	}
	return true
}

func (this *EnemyState) move() bool {
	this.ppos.X = this.pos.X
	this.ppos.Y = this.pos.Y
	this.multiplier -= MULTIPLIER_DECREASE_RATIO
	if this.multiplier < 1 {
		this.multiplier = 1
	}
	if this.destroyedCnt >= 0 {
		this.destroyedCnt++
		this.explodeCnt--
		if this.explodeCnt < 0 {
			this.explodeItv += 2
			this.explodeItv = int(float64(this.explodeItv) * (1.2 + nextFloat(enemyStateRand, 1)))
			this.explodeCnt = this.explodeItv
			this.destroyedEdge(int(math.Sqrt(this.spec.Size()) * 27.0 / (float64(this.explodeItv)*0.1 + 1)))
		}
	}
	this.damaged = false
	if this.damagedCnt > 0 {
		this.damagedCnt--
	}
	alive := false
	for i := 0; i < this.spec.TurretGroupNum(); i++ {
		alive = alive || this.turretGroup[i].move(this.pos, this.deg)
	}
	for i := 0; i < this.spec.MovingTurretGroupNum(); i++ {
		this.movingTurretGroup[i].move(this.pos, this.deg)
	}
	if this.destroyedCnt < 0 && !alive {
		return this.destroyed(nil)
	}
	return true
}

func (this *EnemyState) checkCollision(x, y float64, c sdl.Collidable, shot *Shot) bool {
	ox := math.Abs(this.pos.X - x)
	oy := math.Abs(this.pos.Y - y)
	if ox+oy > this.spec.Size()*2 {
		return false
	}
	for i := 0; i < this.spec.TurretGroupNum(); i++ {
		if this.turretGroup[i].checkCollision(x, y, c, shot) {
			return true
		}
	}
	if this.spec.BridgeShape().CheckCollision(ox, oy, c) {
		this.addDamage(shot.damage(), shot)
		return true
	}
	return false
}

func (this *EnemyState) increaseMultiplier(m float64) {
	this.multiplier += m
}

func (this *EnemyState) addScore(s int) {
	this.setScoreIndicator(s, 1)
}

func (this *EnemyState) addDamage(n int, shot *Shot /* = null */) {
	this.shield -= n
	if this.shield <= 0 {
		this.destroyed(shot)
	} else {
		this.damaged = true
		this.damagedCnt = 7
	}
}

func (this *EnemyState) destroyed(shot *Shot /* = null */) bool {
	var vz float64
	if shot != nil {
		explodeVel.X = SPEED * math.Sin(shot.deg()) / 2
		explodeVel.Y = SPEED * math.Cos(shot.deg()) / 2
		vz = 0
	} else {
		explodeVel.X = 0
		explodeVel.Y = 0
		vz = 0.05
	}
	ss := this.spec.Size() * 1.5
	if ss > 2 {
		ss = 2
	}
	var sn float64
	if this.spec.Size() < 1 {
		sn = this.spec.Size()
	} else {
		sn = math.Sqrt(this.spec.Size())
	}
	if sn > 3 {
		sn = 3
	}
	for i := 0; i < int(sn*8); i++ {
		s := this.smokes.GetInstanceForced()
		s.setVector(this.pos, nextSignedFloat(enemyStateRand, 0.1)+explodeVel.X, nextSignedFloat(enemyStateRand, 0.1)+explodeVel.Y,
			nextFloat(enemyStateRand, vz),
			EXPLOSION, 32+enemyStateRand.Intn(30), ss)
	}
	for i := 0; i < int(sn)*36; i++ {
		sp := this.sparks.GetInstanceForced()
		sp.set(this.pos, nextSignedFloat(enemyStateRand, 0.8)+explodeVel.X, nextSignedFloat(enemyStateRand, 0.8)+explodeVel.Y,
			0.5+nextFloat(enemyStateRand, 0.5), 0.5+nextFloat(enemyStateRand, 0.5), 0, 30+enemyStateRand.Intn(30))
	}
	for i := 0; i < int(sn)*12; i++ {
		f := this.fragments.GetInstanceForced()
		f.set(this.pos, nextSignedFloat(enemyStateRand, 0.33)+explodeVel.X, nextSignedFloat(enemyStateRand, 0.33)+explodeVel.Y,
			0.05+nextFloat(enemyStateRand, 0.1),
			0.2+nextFloat(enemyStateRand, 0.33))
	}
	this.removeTurrets()
	sc := this.spec.score()
	var r bool
	if this.spec.Typ() == EnemySMALL {
		playSe("small_destroyed.wav")
		r = false
	} else {
		playSe("destroyed.wav")
		bn := this.bullets.removeIndexedBullets(this.idx)
		this.destroyedCnt = 0
		this.explodeCnt = 1
		this.explodeItv = 3
		sc += bn * 10
		r = true
		if this.spec.isBoss() {
			this.screen.setScreenShake(45, 0.04)
		}
	}
	this.setScoreIndicator(sc, this.multiplier)
	return r
}

func (this *EnemyState) setScoreIndicator(sc int, mp float64) {
	ty := getTargetY()
	if mp > 1 {
		ni := this.numIndicators.GetInstanceForced()
		ni.setWithVector(sc, SCORE, 0.5, this.pos)
		ni.addTarget(8, ty, RIGHT, 1, 0.5, sc, 40)
		ni.addTarget(11, ty, RIGHT, 0.5, 0.75,
			sc*int(mp), 30)
		ni.addTarget(13, ty, RIGHT, 0.25, 1,
			sc*int(mp*this.stageManager.rankMultiplier()), 20)
		ni.addTarget(12, -8, BOTTOM, 0.5, 0.1,
			sc*int(mp*this.stageManager.rankMultiplier()), 40)
		ni.gotoNextTarget()
		ni = this.numIndicators.GetInstanceForced()
		mn := int(mp * 1000)
		ni.setWithVector(mn, MULTIPLIER, 0.7, this.pos)
		ni.addTarget(10.5, ty, RIGHT, 0.5, 0.2, mn, 70)
		ni.gotoNextTarget()
		ni = this.numIndicators.GetInstanceForced()
		rn := int(this.stageManager.rankMultiplier() * 1000)
		ni.set(rn, MULTIPLIER, 0.4, 11, 8)
		ni.addTarget(13, ty, RIGHT, 0.5, 0.2, rn, 40)
		ni.gotoNextTarget()
		this.scoreReel.addActualScore(sc * int(mp*this.stageManager.rankMultiplier()))
	} else {
		ni := this.numIndicators.GetInstanceForced()
		ni.setWithVector(sc, SCORE, 0.3, this.pos)
		ni.addTarget(11, ty, RIGHT, 1.5, 0.2, sc, 40)
		ni.addTarget(13, ty, RIGHT, 0.25, 0.25,
			sc*int(this.stageManager.rankMultiplier()), 20)
		ni.addTarget(12, -8, BOTTOM, 0.5, 0.1,
			sc*int(this.stageManager.rankMultiplier()), 40)
		ni.gotoNextTarget()
		ni = this.numIndicators.GetInstanceForced()
		rn := int(this.stageManager.rankMultiplier() * 1000)
		ni.set(rn, MULTIPLIER, 0.4, 11, 8)
		ni.addTarget(13, ty, RIGHT, 0.5, 0.2, rn, 40)
		ni.gotoNextTarget()
		this.scoreReel.addActualScore(sc * int(this.stageManager.rankMultiplier()))
	}
}

func (this *EnemyState) destroyedEdge(n int) {
	playSe("explode.wav")
	sn := n
	if sn > 48 {
		sn = 48
	}
	spp := this.spec.Shape().Shape().(*BaseShape).pointPos()
	spd := this.spec.Shape().Shape().(*BaseShape).pointDeg()
	si := enemyStateRand.Intn(len(spp))
	edgePos.X = spp[si].X*this.spec.Size() + this.pos.X
	edgePos.Y = spp[si].Y*this.spec.Size() + this.pos.Y
	ss := this.spec.Size() * 0.5
	if ss > 1 {
		ss = 1
	}
	for i := 0; i < sn; i++ {
		s := this.smokes.GetInstanceForced()
		sr := nextFloat(enemyStateRand, 0.5)
		sd := spd[si] + nextSignedFloat(enemyStateRand, 0.2)
		s.setVector(edgePos, math.Sin(sd)*sr, math.Cos(sd)*sr, -0.004,
			EXPLOSION, 75+enemyStateRand.Intn(25), ss)
		for j := 0; j < 2; j++ {
			sp := this.sparks.GetInstanceForced()
			sp.set(edgePos, math.Sin(sd)*sr*2, math.Cos(sd)*sr*2,
				0.5+nextFloat(enemyStateRand, 0.5), 0.5+nextFloat(enemyStateRand, 0.5), 0, 30+enemyStateRand.Intn(30))
		}
		if i%2 == 0 {
			sf := this.sparkFragments.GetInstanceForced()
			sf.set(edgePos, math.Sin(sd)*sr*0.5, math.Cos(sd)*sr*0.5, 0.06+nextFloat(enemyStateRand, 0.07),
				(0.2 + nextFloat(enemyStateRand, 0.1)))
		}
	}
}

func (this *EnemyState) removeTurrets() {
	for i := 0; i < this.spec.TurretGroupNum(); i++ {
		this.turretGroup[i].remove()
	}
	for i := 0; i < this.spec.MovingTurretGroupNum(); i++ {
		this.movingTurretGroup[i].remove()
	}
}

func (this *EnemyState) draw() {
	gl.PushMatrix()
	if this.destroyedCnt < 0 && this.damagedCnt > 0 {
		enemyDamagedPos.X = this.pos.X + nextSignedFloat(enemyStateRand, float64(this.damagedCnt)*0.01)
		enemyDamagedPos.Y = this.pos.Y + nextSignedFloat(enemyStateRand, float64(this.damagedCnt)*0.01)
		sdl.GlTranslate(enemyDamagedPos)
	} else {
		sdl.GlTranslate(this.pos)
	}
	gl.Rotated(-this.deg*180/math.Pi, 0, 0, 1)
	if this.destroyedCnt >= 0 {
		this.spec.DestroyedShape().Draw()
	} else if !this.damaged {
		this.spec.Shape().Draw()
	} else {
		this.spec.DamagedShape().Draw()
	}
	if this.destroyedCnt < 0 {
		this.spec.BridgeShape().Draw()
	}
	gl.PopMatrix()
	if this.destroyedCnt >= 0 {
		return
	}
	for i := 0; i < this.spec.TurretGroupNum(); i++ {
		this.turretGroup[i].draw()
	}
	if this.multiplier > 1 {
		var ox, oy float64
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
		letter.DrawNumSign(int(this.multiplier*1000), this.pos.X+ox, this.pos.Y+oy, 0.33, 1, 33, 3)
	}
}

/**
 * Base class for a specification of an enemy.
 */

type EnemyType int

const (
	EnemySMALL EnemyType = iota
	EnemyLARGE
	EnemyPLATFORM
)

var enemySpecRand = r.New(r.NewSource(time.Now().Unix()))

type EnemySpec interface {
	score() int
	isBoss() bool
	draw(EnemyState)
	move(EnemyState) bool
	checkCollision(es EnemyState, x, y float64, c sdl.Collidable, shot *Shot) bool
	checkShipCollision(es EnemyState, x, y float64, largeOnly bool /* = false */) bool
	isSmallEnemy() bool
	Shield() int
	TurretGroupNum() int
	MovingTurretGroupNum() int
	Size() float64
	BridgeShape() *EnemyShape
	Typ() EnemyType
	TurretGroupSpec(int) TurretGroupSpec
	MovingTurretGroupSpec(int) MovingTurretGroupSpec
	Shape() *EnemyShape
	DamagedShape() *EnemyShape
	DestroyedShape() *EnemyShape
}

type EnemySpecBase struct {
	field                                            *Field
	ship                                             *Ship
	sparks                                           *SparkPool
	smokes                                           *SmokePool
	fragments                                        *FragmentPool
	wakes                                            *WakePool
	shield                                           int
	_size, distRatio                                 float64
	turretGroupSpec                                  [TURRET_GROUP_MAX]TurretGroupSpec
	turretGroupNum                                   int
	movingTurretGroupSpec                            [MOVING_TURRET_GROUP_MAX]MovingTurretGroupSpec
	movingTurretGroupNum                             int
	shape, damagedShape, destroyedShape, bridgeShape *EnemyShape
	typ                                              EnemyType
}

func (this EnemySpecBase) DamagedShape() *EnemyShape {
	return this.damagedShape
}

func (this EnemySpecBase) DestroyedShape() *EnemyShape {
	return this.destroyedShape
}

func (this EnemySpecBase) Shape() *EnemyShape {
	return this.shape
}

func (this EnemySpecBase) TurretGroupSpec(i int) TurretGroupSpec {
	return this.turretGroupSpec[i]
}
func (this EnemySpecBase) MovingTurretGroupSpec(i int) MovingTurretGroupSpec {
	return this.movingTurretGroupSpec[i]
}

func (this EnemySpecBase) Shield() int {
	return this.shield
}

func (this EnemySpecBase) TurretGroupNum() int {
	return this.turretGroupNum
}

func (this EnemySpecBase) MovingTurretGroupNum() int {
	return this.movingTurretGroupNum
}

func (this EnemySpecBase) Size() float64 {
	return this.size()
}

func (this EnemySpecBase) BridgeShape() *EnemyShape {
	return this.bridgeShape
}

func (this EnemySpecBase) Typ() EnemyType {
	return this.typ
}

func setEnemySpecRandSeed(seed int64) {
	enemySpecRand = r.New(r.NewSource(seed))
}

func NewEnemySpecBase(field *Field, ship *Ship,
	sparks *SparkPool, smokes *SmokePool, fragments *FragmentPool, wakes *WakePool) EnemySpecBase {
	this := EnemySpecBase{}
	this.field = field
	this.ship = ship
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.wakes = wakes
	for i := range this.turretGroupSpec {
		this.turretGroupSpec[i] = NewTurretGroupSpec()
	}
	for i := range this.movingTurretGroupSpec {
		this.movingTurretGroupSpec[i] = NewMovingTurretGroupSpec()
	}
	this.shield = 1
	this._size = 1
	return this
}

func (this *EnemySpecBase) set(typ EnemyType) {
	this.typ = typ
	this._size = 1
	this.distRatio = 0
	this.turretGroupNum = 0
	this.movingTurretGroupNum = 0
}

func (this EnemySpecBase) getTurretGroupSpec() TurretGroupSpec {
	this.turretGroupNum++
	this.turretGroupSpec[this.turretGroupNum-1].init()
	return this.turretGroupSpec[this.turretGroupNum-1]
}

func (this EnemySpecBase) getMovingTurretGroupSpec() MovingTurretGroupSpec {
	this.movingTurretGroupNum++
	this.movingTurretGroupSpec[this.movingTurretGroupNum-1].init()
	return this.movingTurretGroupSpec[this.movingTurretGroupNum-1]
}

func (this *EnemySpecBase) addMovingTurret(rank float64, bossMode bool /* = false */) {
	mtn := int(rank * 0.2)
	if mtn > MOVING_TURRET_GROUP_MAX {
		mtn = MOVING_TURRET_GROUP_MAX
	}
	if mtn >= 2 {
		mtn = 1 + enemySpecRand.Intn(mtn-1)
	} else {
		mtn = 1
	}
	br := rank / float64(mtn)
	var typ BulletMoveType
	if !bossMode {
		switch enemySpecRand.Intn(4) {
		case 0, 1:
			typ = ROLL
		case 2:
			typ = SWING_FIX
		case 3:
			typ = SWING_AIM
		}
	} else {
		typ = ROLL
	}
	rad := 0.9 + nextFloat(enemySpecRand, 0.4) - float64(mtn)*0.1
	radInc := 0.5 + nextFloat(enemySpecRand, 0.25)
	ad := math.Pi * 2
	var a, av, dv, s, sv float64
	switch typ {
	case ROLL:
		a = 0.01 + nextFloat(enemySpecRand, 0.04)
		av = 0.01 + nextFloat(enemySpecRand, 0.03)
		dv = 0.01 + nextFloat(enemySpecRand, 0.04)
	case SWING_FIX:
		ad = math.Pi/10 + nextFloat(enemySpecRand, math.Pi/15)
		s = 0.01 + nextFloat(enemySpecRand, 0.02)
		sv = 0.01 + nextFloat(enemySpecRand, 0.03)
	case SWING_AIM:
		ad = math.Pi/10 + nextFloat(enemySpecRand, math.Pi/15)
		if enemySpecRand.Intn(5) == 0 {
			s = 0.01 + nextFloat(enemySpecRand, 0.01)
		} else {
			s = 0
		}
		sv = 0.01 + nextFloat(enemySpecRand, 0.02)
	}
	for i := 0; i < mtn; i++ {
		tgs := this.getMovingTurretGroupSpec()
		tgs.moveType = typ
		tgs.radiusBase = rad
		var sr float64
		switch typ {
		case ROLL:
			tgs.alignDeg = ad
			tgs.num = 4 + enemySpecRand.Intn(6)
			if enemySpecRand.Intn(2) == 0 {
				if enemySpecRand.Intn(2) == 0 {
					tgs.setRoll(dv, 0, 0)
				} else {
					tgs.setRoll(-dv, 0, 0)
				}
			} else {
				if enemySpecRand.Intn(2) == 0 {
					tgs.setRoll(0, a, av)
				} else {
					tgs.setRoll(0, -a, av)
				}
			}
			if enemySpecRand.Intn(3) == 0 {
				tgs.setRadiusAmp(1+nextFloat(enemySpecRand, 1), 0.01+nextFloat(enemySpecRand, 0.03))
			}
			if enemySpecRand.Intn(2) == 0 {
				tgs.distRatio = 0.8 + nextSignedFloat(enemySpecRand, 0.3)
			}
			sr = br / float64(tgs.num)
		case SWING_FIX:
			tgs.num = 3 + enemySpecRand.Intn(5)
			tgs.alignDeg = ad * (float64(tgs.num)*0.1 + 0.3)
			if enemySpecRand.Intn(2) == 0 {
				tgs.setSwing(s, sv, false)
			} else {
				tgs.setSwing(-s, sv, false)
			}
			if enemySpecRand.Intn(6) == 0 {
				tgs.setRadiusAmp(1+nextFloat(enemySpecRand, 1), 0.01+nextFloat(enemySpecRand, 0.03))
			}
			if enemySpecRand.Intn(4) == 0 {
				tgs.setAlignAmp(0.25+nextFloat(enemySpecRand, 0.25), 0.01+nextFloat(enemySpecRand, 0.02))
			}
			sr = br / float64(tgs.num)
			sr *= 0.6
		case SWING_AIM:
			tgs.num = 3 + enemySpecRand.Intn(4)
			tgs.alignDeg = ad * (float64(tgs.num)*0.1 + 0.3)
			if enemySpecRand.Intn(2) == 0 {
				tgs.setSwing(s, sv, true)
			} else {
				tgs.setSwing(-s, sv, true)
			}
			if enemySpecRand.Intn(4) == 0 {
				tgs.setRadiusAmp(1+nextFloat(enemySpecRand, 1), 0.01+nextFloat(enemySpecRand, 0.03))
			}
			if enemySpecRand.Intn(5) == 0 {
				tgs.setAlignAmp(0.25+nextFloat(enemySpecRand, 0.25), 0.01+nextFloat(enemySpecRand, 0.02))
			}
			sr = br / float64(tgs.num)
			sr *= 0.4
		}
		if enemySpecRand.Intn(4) == 0 {
			tgs.setXReverse(-1)
		}
		tgs.turretSpec.setParam(sr, TurretMOVING, enemySpecRand)
		if bossMode {
			tgs.turretSpec.setBossSpec()
		}
		rad += radInc
		ad *= 1 + nextSignedFloat(enemySpecRand, 0.2)
	}
}

func (this EnemySpecBase) checkCollision(es EnemyState, x, y float64, c sdl.Collidable, shot *Shot) bool {
	return es.checkCollision(x, y, c, shot)
}

func (this EnemySpecBase) checkShipCollision(es EnemyState, x, y float64, largeOnly bool /* = false */) bool {
	if es.destroyedCnt >= 0 || (largeOnly && this.typ != EnemyLARGE) {
		return false
	}
	return this.shape.checkShipCollision(x-es.pos.X, y-es.pos.Y, es.deg)
}

func (this EnemySpecBase) move(es EnemyState) bool {
	return es.move()
}

func (this EnemySpecBase) draw(es EnemyState) {
	es.draw()
}

func (this EnemySpecBase) size() float64 {
	return this._size
}

func (this *EnemySpecBase) setSize(v float64) float64 {
	this._size = v
	if this.shape != nil {
		this.shape.SetSize(this._size)
	}
	if this.damagedShape != nil {
		this.damagedShape.SetSize(this._size)
	}
	if this.destroyedShape != nil {
		this.destroyedShape.SetSize(this._size)
	}
	if this.bridgeShape != nil {
		s := 0.9
		this.bridgeShape.SetSize(s * (1 - this.distRatio))
	}
	return this._size
}

func (this EnemySpecBase) isSmallEnemy() bool {
	return this.typ == EnemySMALL
}

type HasAppearType interface {
	setFirstState(es EnemyState, appType AppearanceType) bool
}

/**
 * Specification for a small class ship.
 */

var _ EnemySpec = &SmallShipEnemySpec{}
var _ HasAppearType = &SmallShipEnemySpec{}

type EnemyMoveType int

const (
	STOPANDGO EnemyMoveType = iota
	CHASE
)

type MoveState int

const (
	EnemySTAYING MoveState = iota
	EnemyMOVING
)

type SmallShipEnemySpec struct {
	EnemySpecBase

	typ                        EnemyMoveType
	accel, maxSpeed, staySpeed float64
	moveDuration, stayDuration int
	speed, turnDeg             float64
}

func NewSmallShipEnemySpec(field *Field, ship *Ship,
	sparks *SparkPool, smokes *SmokePool, fragments *FragmentPool, wakes *WakePool) SmallShipEnemySpec {
	this := SmallShipEnemySpec{
		EnemySpecBase: NewEnemySpecBase(field, ship, sparks, smokes, fragments, wakes),
	}
	this.moveDuration = 1
	this.stayDuration = 1
	return this
}

func (this *SmallShipEnemySpec) setParam(rank float64, rand *r.Rand) {
	this.set(EnemySMALL)
	this.shape = NewEnemyShape(EnemyShapeSMALL)
	this.damagedShape = NewEnemyShape(EnemyShapeSMALL_DAMAGED)
	this.bridgeShape = NewEnemyShape(EnemyShapeSMALL_BRIDGE)
	this.typ = EnemyMoveType(rand.Intn(2))
	sr := nextFloat(rand, rank*0.8)
	if sr > 25 {
		sr = 25
	}
	switch this.typ {
	case STOPANDGO:
		this.distRatio = 0.5
		this.setSize(0.47 + nextFloat(rand, 0.1))
		this.accel = 0.5 - 0.5/(2.0+nextFloat(rand, rank))
		this.maxSpeed = 0.05 * (1.0 + sr)
		this.staySpeed = 0.03
		this.moveDuration = 32 + nextSignedInt(rand, 12)
		this.stayDuration = 32 + nextSignedInt(rand, 12)
	case CHASE:
		this.distRatio = 0.5
		this.setSize(0.5 + nextFloat(rand, 0.1))
		this.speed = 0.036 * (1.0 + sr)
		this.turnDeg = 0.02 + nextSignedFloat(rand, 0.04)
	}
	this.shield = 1
	tgs := this.getTurretGroupSpec()
	tgs.turretSpec.setParam(rank-sr*0.5, TurretSMALL, rand)
}

func (this *SmallShipEnemySpec) setFirstState(es EnemyState, appType AppearanceType) bool {
	es.setSpec(this)
	if !es.setAppearancePos(this.field, this.ship, enemySpecRand, appType) {
		return false
	}
	switch this.typ {
	case STOPANDGO:
		es.speed = 0
		es.state = EnemyMOVING
		es.cnt = this.moveDuration
	case CHASE:
		es.speed = this.speed
	}
	return true
}

func (this SmallShipEnemySpec) move(es EnemyState) bool {
	if !this.EnemySpecBase.move(es) {
		return false
	}
	switch this.typ {
	case STOPANDGO:
		es.pos.X += math.Sin(es.velDeg) * es.speed
		es.pos.Y += math.Cos(es.velDeg) * es.speed
		es.pos.Y -= this.field.lastScrollY()
		if es.pos.Y <= -this.field.outerSize().Y {
			return false
		}
		if this.field.getBlockVector(es.pos) >= 0 || !this.field.checkInOuterHeightField(es.pos) {
			es.velDeg += math.Pi
			es.pos.X += math.Sin(es.velDeg) * es.speed * 2
			es.pos.Y += math.Cos(es.velDeg) * es.speed * 2
		}
		switch es.state {
		case EnemyMOVING:
			es.speed += (this.maxSpeed - es.speed) * this.accel
			es.cnt--
			if es.cnt <= 0 {
				es.velDeg = nextFloat(enemySpecRand, math.Pi*2)
				es.cnt = this.stayDuration
				es.state = EnemySTAYING
			}
		case EnemySTAYING:
			es.speed += (this.staySpeed - es.speed) * this.accel
			es.cnt--
			if es.cnt <= 0 {
				es.cnt = this.moveDuration
				es.state = EnemyMOVING
			}
		}
	case CHASE:
		es.pos.X += math.Sin(es.velDeg) * this.speed
		es.pos.Y += math.Cos(es.velDeg) * this.speed
		es.pos.Y -= this.field.lastScrollY()
		if es.pos.Y <= -this.field.outerSize().Y {
			return false
		}
		if this.field.getBlockVector(es.pos) >= 0 || !this.field.checkInOuterHeightField(es.pos) {
			es.velDeg += math.Pi
			es.pos.X += math.Sin(es.velDeg) * es.speed * 2
			es.pos.Y += math.Cos(es.velDeg) * es.speed * 2
		}
		var ad float64
		shipPos := this.ship.nearPos(es.pos)
		if shipPos.DistVector(es.pos) < 0.1 {
			ad = 0
		} else {
			ad = math.Atan2(shipPos.X-es.pos.X, shipPos.Y-es.pos.Y)
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
		this.shape.addWake(this.wakes, es.pos, es.deg, es.speed)
	}
	return true
}

func (this SmallShipEnemySpec) score() int {
	return 50
}

func (this SmallShipEnemySpec) isBoss() bool {
	return false
}

/**
 * Specification for a large/middle class ship.
 */
var _ EnemySpec = &ShipEnemySpec{}
var _ HasAppearType = &ShipEnemySpec{}

type ShipClass int

const (
	ShipMIDDLE ShipClass = iota
	ShipLARGE
	ShipBOSS
)

const SINK_INTERVAL = 120

type ShipEnemySpec struct {
	EnemySpecBase

	speed, degVel float64
	shipClass     ShipClass
}

func NewShipEnemySpec(field *Field, ship *Ship,
	sparks *SparkPool, smokes *SmokePool, fragments *FragmentPool, wakes *WakePool) ShipEnemySpec {
	this := ShipEnemySpec{
		EnemySpecBase: NewEnemySpecBase(field, ship, sparks, smokes, fragments, wakes),
	}
	return this
}

var constad = [6]float64{math.Pi / 4, -math.Pi / 4, math.Pi / 2, -math.Pi / 2, math.Pi / 4 * 3, -math.Pi / 4 * 3}

func (this *ShipEnemySpec) setParam(rank float64, cls ShipClass, rand *r.Rand) {
	this.shipClass = cls
	this.set(EnemyLARGE)
	this.shape = NewEnemyShape(EnemyShapeMIDDLE)
	this.damagedShape = NewEnemyShape(EnemyShapeMIDDLE_DAMAGED)
	this.destroyedShape = NewEnemyShape(EnemyShapeMIDDLE_DESTROYED)
	this.bridgeShape = NewEnemyShape(EnemyShapeMIDDLE_BRIDGE)
	this.distRatio = 0.7
	mainTurretNum := 0
	subTurretNum := 0
	movingTurretRatio := 0.0
	rk := rank
	switch cls {
	case ShipMIDDLE:
		sz := 1.5 + rank/15 + nextFloat(rand, rank/15)
		ms := 2 + nextFloat(rand, 0.5)
		if sz > ms {
			sz = ms
		}
		this.setSize(sz)
		this.speed = 0.015 + nextSignedFloat(rand, 0.005)
		this.degVel = 0.005 + nextSignedFloat(rand, 0.003)
		switch rand.Intn(3) {
		case 0:
			mainTurretNum = int(this.size()*(1+nextSignedFloat(rand, 0.25)) + 1)
		case 1:
			subTurretNum = int(this.size()*1.6*(1+nextSignedFloat(rand, 0.5)) + 2)
		case 2:
			mainTurretNum = int(this.size()*(0.5+nextSignedFloat(rand, 0.12)) + 1)
			movingTurretRatio = 0.5 + nextFloat(rand, 0.25)
			rk = rank * (1 - movingTurretRatio)
			movingTurretRatio *= 2
		}
	case ShipLARGE:
		sz := 2.5 + rank/24 + nextFloat(rand, rank/24)
		ms := 3 + nextFloat(rand, 1)
		if sz > ms {
			sz = ms
		}
		this.setSize(sz)
		this.speed = 0.01 + nextSignedFloat(rand, 0.005)
		this.degVel = 0.003 + nextSignedFloat(rand, 0.002)
		mainTurretNum = int(this.size()*(0.7+nextSignedFloat(rand, 0.2)) + 1)
		subTurretNum = int(this.size()*1.6*(0.7+nextSignedFloat(rand, 0.33)) + 2)
		movingTurretRatio = 0.25 + nextFloat(rand, 0.5)
		rk = rank * (1 - movingTurretRatio)
		movingTurretRatio *= 3
	case ShipBOSS:
		sz := 5 + rank/30 + nextFloat(rand, rank/30)
		ms := 9 + nextFloat(rand, 3)
		if sz > ms {
			sz = ms
		}
		this.setSize(sz)
		this.speed = this.ship.scrollSpeedBase() + 0.0025 + nextSignedFloat(rand, 0.001)
		this.degVel = 0.003 + nextSignedFloat(rand, 0.002)
		mainTurretNum = int(this.size()*0.8*(1.5+nextSignedFloat(rand, 0.4)) + 2)
		subTurretNum = int(this.size()*0.8*(2.4+nextSignedFloat(rand, 0.6)) + 2)
		movingTurretRatio = 0.2 + nextFloat(rand, 0.3)
		rk = rank * (1 - movingTurretRatio)
		movingTurretRatio *= 2.5
	}
	this.shield = int(this.size() * 10)
	if cls == ShipBOSS {
		this.shield = int(float64(this.shield) * 2.4)
	}
	if mainTurretNum+subTurretNum <= 0 {
		tgs := this.getTurretGroupSpec()
		tgs.turretSpec.setParam(0, TurretDUMMY, rand)
	} else {
		subTurretRank := rk / float64(mainTurretNum*3+subTurretNum)
		mainTurretRank := subTurretRank * 2.5
		if cls != ShipBOSS {
			frontMainTurretNum := int(float64(mainTurretNum)/2 + 0.99)
			rearMainTurretNum := mainTurretNum - frontMainTurretNum
			if frontMainTurretNum > 0 {
				tgs := this.getTurretGroupSpec()
				tgs.turretSpec.setParam(mainTurretRank, TurretMAIN, rand)
				tgs.num = frontMainTurretNum
				tgs.alignType = STRAIGHT
				tgs.offset.Y = -this.size() * (0.9 + nextSignedFloat(rand, 0.05))
			}
			if rearMainTurretNum > 0 {
				tgs := this.getTurretGroupSpec()
				tgs.turretSpec.setParam(mainTurretRank, TurretMAIN, rand)
				tgs.num = rearMainTurretNum
				tgs.alignType = STRAIGHT
				tgs.offset.Y = this.size() * (0.9 + nextSignedFloat(rand, 0.05))
			}
			var pts TurretSpec
			if subTurretNum > 0 {
				frontSubTurretNum := (subTurretNum + 2) / 4
				rearSubTurretNum := (subTurretNum - frontSubTurretNum*2) / 2
				tn := frontSubTurretNum
				ad := -math.Pi / 4
				for i := 0; i < 4; i++ {
					if i == 2 {
						tn = rearSubTurretNum
					}
					if tn <= 0 {
						continue
					}
					tgs := this.getTurretGroupSpec()
					if i == 0 || i == 2 {
						if rand.Intn(2) == 0 {
							tgs.turretSpec.setParam(subTurretRank, TurretSUB, rand)
						} else {
							tgs.turretSpec.setParam(subTurretRank, TurretSUB_DESTRUCTIVE, rand)
						}
						pts = tgs.turretSpec
					} else {
						tgs.turretSpec.setParamTurretSpec(pts)
					}
					tgs.num = tn
					tgs.alignType = ROUND
					tgs.alignDeg = ad
					ad += math.Pi / 2
					tgs.alignWidth = math.Pi/6 + nextFloat(rand, math.Pi/8)
					tgs.radius = this.size() * 0.75
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
				ad := -math.Pi / 4
				for i := 0; i < 4; i++ {
					if i == 2 {
						tn = rearMainTurretNum
					}
					if tn <= 0 {
						continue
					}
					tgs := this.getTurretGroupSpec()
					if i == 0 || i == 2 {
						tgs.turretSpec.setParam(mainTurretRank, TurretMAIN, rand)
						pts = tgs.turretSpec
						pts.setBossSpec()
					} else {
						tgs.turretSpec.setParamTurretSpec(pts)
					}
					tgs.num = tn
					tgs.alignType = ROUND
					tgs.alignDeg = ad
					ad += math.Pi / 2
					tgs.alignWidth = math.Pi/6 + nextFloat(rand, math.Pi/8)
					tgs.radius = this.size() * 0.45
					tgs.distRatio = this.distRatio
				}
			}
			if subTurretNum > 0 {
				var tn [3]int
				tn[0] = (subTurretNum + 2) / 6
				tn[1] = (subTurretNum - tn[0]*2) / 4
				tn[2] = (subTurretNum - tn[0]*2 - tn[1]*2) / 2
				for i := 0; i < 6; i++ {
					idx := i / 2
					if tn[idx] <= 0 {
						continue
					}
					tgs := this.getTurretGroupSpec()
					if i == 0 || i == 2 || i == 4 {
						if rand.Intn(2) == 0 {
							tgs.turretSpec.setParam(subTurretRank, TurretSUB, rand)
						} else {
							tgs.turretSpec.setParam(subTurretRank, TurretSUB_DESTRUCTIVE, rand)
						}
						pts = tgs.turretSpec
						pts.setBossSpec()
					} else {
						tgs.turretSpec.setParamTurretSpec(pts)
					}
					tgs.num = tn[idx]
					tgs.alignType = ROUND
					tgs.alignDeg = constad[i]
					tgs.alignWidth = math.Pi/7 + nextFloat(rand, math.Pi/9)
					tgs.radius = this.size() * 0.75
					tgs.distRatio = this.distRatio
				}
			}
		}
	}
	if movingTurretRatio > 0 {
		if cls == ShipBOSS {
			this.addMovingTurret(rank*movingTurretRatio, true)
		} else {
			this.addMovingTurret(rank*movingTurretRatio, false)
		}
	}
}

func (this *ShipEnemySpec) setFirstState(es EnemyState, appType AppearanceType) bool {
	es.setSpec(this)
	if !es.setAppearancePos(this.field, this.ship, enemySpecRand, appType) {
		return false
	}
	es.speed = this.speed
	if es.pos.X < 0 {
		es.turnWay = -1
	} else {
		es.turnWay = 1
	}
	if this.isBoss() {
		es.trgDeg = nextFloat(enemySpecRand, 0.1) + 0.1
		if enemySpecRand.Intn(2) == 0 {
			es.trgDeg *= -1
		}
		es.turnCnt = 250 + enemySpecRand.Intn(150)
	}
	return true
}

func (this ShipEnemySpec) move(es EnemyState) bool {
	if es.destroyedCnt >= SINK_INTERVAL {
		return false
	}
	if !this.EnemySpecBase.move(es) {
		return false
	}
	es.pos.X += math.Sin(es.deg) * es.speed
	es.pos.Y += math.Cos(es.deg) * es.speed
	es.pos.Y -= this.field.lastScrollY()
	if es.pos.X <= -this.field.outerSize().X-this.size() || es.pos.X >= this.field.outerSize().X+this.size() ||
		es.pos.Y <= -this.field.outerSize().Y-this.size() {
		return false
	}
	if es.pos.Y > this.field.outerSize().Y*2.2+this.size() {
		es.pos.Y = this.field.outerSize().Y*2.2 + this.size()
	}
	if this.isBoss() {
		es.turnCnt--
		if es.turnCnt <= 0 {
			es.turnCnt = 250 + enemySpecRand.Intn(150)
			es.trgDeg = nextFloat(enemySpecRand, 0.1) + 0.2
			if es.pos.X > 0 {
				es.trgDeg *= -1
			}
		}
		es.deg += (es.trgDeg - es.deg) * 0.0025
		if this.ship.higherPos().Y > es.pos.Y {
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
		this.shape.addWake(this.wakes, es.pos, es.deg, es.speed)
	}
	return true
}

func (this ShipEnemySpec) draw(es EnemyState) {
	if es.destroyedCnt >= 0 {
		sdl.SetColor(
			MIDDLE_COLOR_R*(1-float64(es.destroyedCnt)/SINK_INTERVAL)*0.5,
			MIDDLE_COLOR_G*(1-float64(es.destroyedCnt)/SINK_INTERVAL)*0.5,
			MIDDLE_COLOR_B*(1-float64(es.destroyedCnt)/SINK_INTERVAL)*0.5, 1)
	}
	this.EnemySpecBase.draw(es)
}

func (this ShipEnemySpec) score() int {
	switch this.shipClass {
	case ShipMIDDLE:
		return 100
	case ShipLARGE:
		return 300
	case ShipBOSS:
		return 1000
	}
	return 0
}

func (this ShipEnemySpec) isBoss() bool {
	return this.shipClass == ShipBOSS
}

/**
 * Specification for a sea-based platform.
 */
var _ EnemySpec = &PlatformEnemySpec{}

type PlatformEnemySpec struct {
	EnemySpecBase
}

func NewPlatformEnemySpec(field *Field, ship *Ship,
	sparks *SparkPool, smokes *SmokePool, fragments *FragmentPool, wakes *WakePool) PlatformEnemySpec {
	this := PlatformEnemySpec{NewEnemySpecBase(field, ship, sparks, smokes, fragments, wakes)}
	return this
}

func (this *PlatformEnemySpec) setParam(rank float64, rand *r.Rand) {
	this.set(EnemyPLATFORM)
	this.shape = NewEnemyShape(EnemyShapePLATFORM)
	this.damagedShape = NewEnemyShape(EnemyShapePLATFORM_DAMAGED)
	this.destroyedShape = NewEnemyShape(EnemyShapePLATFORM_DESTROYED)
	this.bridgeShape = NewEnemyShape(EnemyShapePLATFORM_BRIDGE)
	this.distRatio = 0
	this.setSize(1 + rank/30 + nextFloat(rand, rank/30))
	ms := 1 + nextFloat(rand, 0.25)
	if this.size() > ms {
		this.setSize(ms)
	}
	mainTurretNum := 0
	frontTurretNum := 0
	sideTurretNum := 0
	rk := rank
	movingTurretRatio := 0.0
	switch rand.Intn(3) {
	case 0:
		frontTurretNum = int(this.size()*(2+nextSignedFloat(rand, 0.5)) + 1)
		movingTurretRatio = 0.33 + nextFloat(rand, 0.46)
		rk *= (1 - movingTurretRatio)
		movingTurretRatio *= 2.5
	case 1:
		frontTurretNum = int(this.size()*(0.5+nextSignedFloat(rand, 0.2)) + 1)
		sideTurretNum = int(this.size()*(0.5+nextSignedFloat(rand, 0.2))+1) * 2
	case 2:
		mainTurretNum = int(this.size()*(1+nextSignedFloat(rand, 0.33)) + 1)
	}
	this.shield = int(this.size() * 20)
	subTurretNum := frontTurretNum + sideTurretNum
	subTurretRank := rk / float64(mainTurretNum*3+subTurretNum)
	mainTurretRank := subTurretRank * 2.5
	if mainTurretNum > 0 {
		tgs := this.getTurretGroupSpec()
		tgs.turretSpec.setParam(mainTurretRank, TurretMAIN, rand)
		tgs.num = mainTurretNum
		tgs.alignType = ROUND
		tgs.alignDeg = 0
		tgs.alignWidth = math.Pi*0.66 + nextFloat(rand, math.Pi/2)
		tgs.radius = this.size() * 0.7
		tgs.distRatio = this.distRatio
	}
	if frontTurretNum > 0 {
		tgs := this.getTurretGroupSpec()
		tgs.turretSpec.setParam(subTurretRank, TurretSUB, rand)
		tgs.num = frontTurretNum
		tgs.alignType = ROUND
		tgs.alignDeg = 0
		tgs.alignWidth = math.Pi/5 + nextFloat(rand, math.Pi/6)
		tgs.radius = this.size() * 0.8
		tgs.distRatio = this.distRatio
	}
	sideTurretNum /= 2
	if sideTurretNum > 0 {
		var pts TurretSpec
		for i := 0; i < 2; i++ {
			tgs := this.getTurretGroupSpec()
			if i == 0 {
				tgs.turretSpec.setParam(subTurretRank, TurretSUB, rand)
				pts = tgs.turretSpec
			} else {
				tgs.turretSpec.setParamTurretSpec(pts)
			}
			tgs.num = sideTurretNum
			tgs.alignType = ROUND
			tgs.alignDeg = math.Pi/2 - math.Pi*float64(i)
			tgs.alignWidth = math.Pi/5 + nextFloat(rand, math.Pi/6)
			tgs.radius = this.size() * 0.75
			tgs.distRatio = this.distRatio
		}
	}
	if movingTurretRatio > 0 {
		this.addMovingTurret(rank*movingTurretRatio, false)
	}
}

func (this *PlatformEnemySpec) setFirstState(es EnemyState, x, y, d float64) bool {
	es.setSpec(this)
	es.pos.X = x
	es.pos.Y = y
	es.deg = d
	es.speed = 0
	return es.checkFrontClear(true)
}

func (this PlatformEnemySpec) move(es EnemyState) bool {
	if !this.EnemySpecBase.move(es) {
		return false
	}
	es.pos.Y -= this.field.lastScrollY()
	return es.pos.Y > -this.field.outerSize().Y
}

func (this PlatformEnemySpec) score() int {
	return 100
}

func (this PlatformEnemySpec) isBoss() bool {
	return false
}

type EnemyPool struct {
	actor.ActorPool
}

func NewEnemyPool(n int, args []interface{}) *EnemyPool {
	f := func() actor.Actor { return NewEnemy() }
	this := &EnemyPool{
		ActorPool: actor.NewActorPool(f, n, args),
	}
	for _, a := range this.Actor {
		e := a.(*Enemy)
		e.setEnemyPool(this)
	}
	return this
}

func (this *EnemyPool) GetInstance() *Enemy {
	return this.ActorPool.GetInstance().(*Enemy)
}

func (this *EnemyPool) GetInstanceForced() *Enemy {
	return this.ActorPool.GetInstanceForced().(*Enemy)
}

func (this *EnemyPool) setStageManager(stageManager *StageManager) {
	for _, a := range this.Actor {
		e := a.(*Enemy)
		e.setStageManager(stageManager)
	}
}

func (this *EnemyPool) checkShotHit(pos vector.Vector, shape sdl.Collidable, shot *Shot /* = null */) {
	for _, a := range this.Actor {
		e := a.(*Enemy)
		if e.Exists() {
			e.checkShotHit(pos, shape, shot)
		}
	}
}

func (this *EnemyPool) checkHitShip(x, y float64, deselection *Enemy /* = null */, largeOnly bool /* = false */) *Enemy {
	for _, a := range this.Actor {
		e := a.(*Enemy)
		if e.Exists() && e != deselection {
			if e.checkHitShip(x, y, largeOnly) {
				return e
			}
		}
	}
	return nil
}

func (this *EnemyPool) hasBoss() bool {
	for _, a := range this.Actor {
		e := a.(*Enemy)
		if e.Exists() && e.isBoss() {
			return true
		}
	}
	return false
}
