package main

import (
	"fmt"
	"math"
	"time"

	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"
)

/**
 * Player's ship.
 */

const SCROLL_SPEED_BASE = 0.01
const SCROLL_SPEED_MAX = 0.1
const SCROLL_START_Y = 2.5

type Ship struct {
	field                                                Field
	boat                                                 [2]Boat
	gameMode                                             int
	boatNum                                              int
	InGameState                                          gameState
	scrollSpeed, _scrollSpeedBase                        float64
	_midstPos, _higherPos, _lowerPos, _nearPos, _nearVel vector.Vector
	bridgeShape                                          BaseShape
}

func NewShip(twinStick TwinStick, field Field, screen Screen,
	sparks SparkPool, smokes SmokePool, fragments FragmentPool, wakes WakePool) {
	this := &Ship{}
	this.field = field
	Boat.init()
	for i := range this.boat {
		this.boat[i] = NewBoat(i, this, twinStick,
			field, screen, sparks, smokes, fragments, wakes)
		i++
	}
	this.boatNum = 1
	this.scrollSpeed = SCROLL_SPEED_BASE
	this._scrollSpeedBase = SCROLL_SPEED_BASE
	this.bridgeShape = NewBaseShape(0.3, 0.2, 0.1, BRIDGE, 0.3, 0.7, 0.7)
	return ship
}

func (this *Ship) setRandSeed(seed int64) {
	Boat.setRandSeed(seed)
}

func (this *Ship) setShots(shots *ShotPool) {
	for _, b := range this.boat {
		b.setShots(shots)
	}
}

func (this *Ship) setEnemies(enemies EnemyPool) {
	for _, b := range this.boat {
		b.setEnemies(enemies)
	}
}

func (this *Ship) setStageManager(stageManager StageManager) {
	for _, b := range this.boat {
		b.setStageManager(stageManager)
	}
}

func (this *Ship) setGameState(gameState InGameState) {
	this.gameState = gameState
	for _, b := range this.boat {
		b.setGameState(gameState)
	}
}

func (this *Ship) start(gameMode int) {
	this.gameMode = gameMode
	if gameMode == InGameState.GameMode.DOUBLE_PLAY {
		this.boatNum = 2
	} else {
		this.boatNum = 1
	}
	this._scrollSpeedBase = SCROLL_SPEED_BASE
	for i := 0; i < boatNum; i++ {
		this.boat[i].start(gameMode)
	}
	this._midstPos.x = 0
	this._midstPos.y = 0
	this._higherPos.x = 0
	this._higherPos.y = 0
	this._lowerPos.x = 0
	this._lowerPos.y = 0
	this._nearPos.x = 0
	this._nearPos.y = 0
	this._nearVel.x = 0
	this._nearVel.y = 0
	this.restart()
}

func (this *Ship) restart() {
	this.scrollSpeed = this._scrollSpeedBase
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].restart()
	}
}

func (this *Ship) move() {
	this.field.scroll(scrollSpeed)
	sf := false
	for i := 0; i < boatNum; i++ {
		this.boat[i].move()
		if this.boat[i].hasCollision && this.boat[i].pos.X > this.field.size.X/3 && this.boat[i].pos.Y < -this.field.size.Y/4*3 {
			sf = true
		}
	}
	if sf {
		this.gameState.shrinkScoreReel()
	}
	if this.higherPos.y >= SCROLL_START_Y {
		this.scrollSpeed += (SCROLL_SPEED_MAX - this.scrollSpeed) * 0.1
	} else {
		this.scrollSpeed += (this._scrollSpeedBase - this.scrollSpeed) * 0.1
	}
	this._scrollSpeedBase += (SCROLL_SPEED_MAX - this._scrollSpeedBase) * 0.00001
}

func (this *Ship) checkBulletHit(p, pp vector.Vector) bool {
	for i := 0; i < boatNum; i++ {
		if boat[i].checkBulletHit(p, pp) {
			return true
		}
	}
	return false
}

func (this *Ship) clearBullets() {
	this.gameState.clearBullets()
}

func (this *Ship) destroyed() {
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].destroyedBoat()
	}
}

func (this *Ship) draw() {
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].draw()
	}
	if this.gameMode == InGameState.GameMode.DOUBLE_PLAY && this.boat[0].hasCollision {
		sdl.setColor(0.5, 0.5, 0.9, 0.8)
		gl.Begin(gl.LINE_STRIP)
		gl.Vertex2f(this.boat[0].pos.X, this.boat[0].pos.Y)
		sdl.SetColor(0.5, 0.5, 0.9, 0.3)
		gl.Vertex2f(this.midstPos.X, this.midstPos.Y)
		sdl.SetColor(0.5, 0.5, 0.9, 0.8)
		gl.Vertex2f(this.boat[1].pos.X, this.boat[1].pos.Y)
		gl.End()
		gl.PushMatrix()
		sdl.glTranslate(this.midstPos)
		gl.Rotatef(-this.degAmongBoats*180/math.Pi, 0, 0, 1)
		this.bridgeShape.draw()
		gl.PopMatrix()
	}
}

func (this *Ship) drawFront() {
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].drawFront()
	}
}

func (this *Ship) drawShape() {
	this.boat[0].drawShape()
}

func (this *Ship) scrollSpeedBase() float64 {
	return this._scrollSpeedBase
}

func (this *Ship) setReplayMode(turnSpeed float64, reverseFire bool) {
	for _, b := range this.boat {
		b.setReplayMode(turnSpeed, reverseFire)
	}
}

func (this *Ship) unsetReplayMode() {
	for _, b := range this.boat {
		b.unsetReplayMode()
	}
}

func (this *Ship) replayMode() bool {
	return this.boat[0].replayMode()
}

func (this *Ship) midstPos() vector.Vector {
	this._midstPos.x = 0
	this._midstPos.y = 0
	for i := 0; i < this.boatNum; i++ {
		this._midstPos.X += this.boat[i].pos.X
		this._midstPos.Y += this.boat[i].pos.Y
	}
	this._midstPos /= this.boatNum
	return this._midstPos
}

func (this *Ship) higherPos() vector.Vector {
	this._higherPos.Y = -99999
	for i := 0; i < boatNum; i++ {
		if this.boat[i].pos.Y > this._higherPos.Y {
			this._higherPos.X = this.boat[i].pos.X
			this._higherPos.Y = this.boat[i].pos.Y
		}
	}
	return this._higherPos
}

func (this *Ship) lowerPos() vector.Vector {
	this._lowerPos.Y = 99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos.Y < this._lowerPos.Y {
			this._lowerPos.X = this.boat[i].pos.X
			this._lowerPos.Y = this.boat[i].pos.Y
		}
	}
	return this._lowerPos
}

func (this *Ship) nearPos(p vector.Vector) vector.Vector {
	dist := 99999.0
	for i := 0; i < boatNum; i++ {
		if this.boat[i].pos.dist(p) < dist {
			dist = this.boat[i].pos.dist(p)
			this._nearPos.X = this.boat[i].pos.X
			this._nearPos.Y = this.boat[i].pos.Y
		}
	}
	return _nearPos
}

func (this *Ship) nearVel(p Vector) vector.Vector {
	dist := 99999.0
	for i := 0; i < boatNum; i++ {
		if this.boat[i].pos.dist(p) < dist {
			dist = this.boat[i].pos.dist(p)
			this._nearVel.x = this.boat[i].vel.X
			this._nearVel.y = this.boat[i].vel.Y
		}
	}
	return this._nearVel
}

func (this *Ship) distAmongBoats() float64 {
	return this.boat[0].pos.dist(this.boat[1].pos)
}

func (this *Ship) degAmongBoats() float64 {
	if this.distAmongBoats < 0.1 {
		return 0
	} else {
		return math.Atan2(this.boat[0].pos.X-this.boat[1].pos.X, this.boat[0].pos.Y-this.boat[1].pos.Y)
	}
}

const RESTART_CNT = 300
const INVINCIBLE_CNT = 228
const HIT_WIDTH = 0.02
const FIRE_INTERVAL = 2
const FIRE_INTERVAL_MAX = 4
const FIRE_LANCE_INTERVAL = 15
const SPEED_BASE = 0.15
const TURN_RATIO_BASE = 0.2
const SLOW_TURN_RATIO = 0.0
const TURN_CHANGE_RATIO = 0.5

var boatRand = r.New(r.NewSource(time.Now().Unix()))
var stickInput TwinStickState

type Boat struct {
	twinStick                RecordableTwinStick
	field                    Field
	screen                   Screen
	shots                    ShotPool
	sparks                   SparkPool
	smokes                   SmokePool
	fragments                FragmentPool
	wakes                    WakePool
	enemies                  EnemyPool
	stageManager             StageManager
	gameState                InGameState
	_pos, firePos            vector.Vector
	deg, speed, turnRatio    float64
	_shape                   BaseShape
	bridgeShape              BaseShape
	fireCnt, fireSprCnt      int
	fireInterval, fireSprDeg float64
	fireLanceCnt             int
	fireDeg                  float64
	aPressed, bPressed       bool
	cnt                      int
	onBlock                  bool
	_vel, refVel             vector.Vector
	shieldCnt                int
	shieldShape              ShieldShape
	_replayMode              bool
	turnSpeed                float64
	reverseFire              bool
	gameMode                 int
	vx, vy                   float64
	idx                      int
	ship                     *Ship
}

func setBoatRandSeed(seed int64) {
	boatRand = r.New(r.NewSoure(seed))
}

func NewBoat(idx int, ship *Ship,
	twinStick TwinStick,
	field Field, screen Screen,
	sparks SparkPool, smokes SmokePool, fragments FragmentPool, wakes WakePool) *Boat {
	this := &Boat{}
	this.idx = idx
	this.ship = ship
	this.twinStick = twinStick
	this.field = field
	this.screen = screen
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.wakes = wakes
	switch idx {
	case 0:
		this._shape = NewBaseShape(0.7, 0.6, 0.6, BaseShape.ShapeType.SHIP_ROUNDTAIL, 0.5, 0.7, 0.5)
		this.bridgeShape = NewBaseShape(0.3, 0.6, 0.6, BaseShape.ShapeType.BRIDGE, 0.3, 0.7, 0.3)
	case 1:
		this._shape = NewBaseShape(0.7, 0.6, 0.6, BaseShape.ShapeType.SHIP_ROUNDTAIL, 0.4, 0.3, 0.8)
		this.bridgeShape = NewBaseShape(0.3, 0.6, 0.6, BaseShape.ShapeType.BRIDGE, 0.2, 0.3, 0.6)
	}
	this.turnSpeed = 1
	this.fireInterval = FIRE_INTERVAL
	this.shieldShape = NewShieldShape()
	return this
}

func (this *Boat) setShots(shots ShotPool) {
	this.shots = shots
}

func (this *Boat) setEnemies(enemies EnemyPool) {
	this.enemies = enemies
}

func (this *Boat) setStageManager(stageManager StageManager) {
	this.stageManager = stageManager
}

func (this *Boat) setGameState(gameState InGameState) {
	this.gameState = gameState
}

func (this *Boat) start(gameMode int) {
	this.gameMode = gameMode
	if gameMode == InGameState.GameMode.DOUBLE_PLAY {
		switch idx {
		case 0:
			this._pos.X = -this.field.size.X * 0.5
		case 1:
			this._pos.X = this.field.size.X * 0.5
		}
	} else {
		this._pos.X = 0
	}
	this._pos.Y = -this.field.size.Y * 0.8
	this.firePos.X = 0
	this.firePos.Y = 0
	this._vel.x = 0
	this._vel.y = 0
	this.deg = 0
	this.speed = SPEED_BASE
	this.turnRatio = TURN_RATIO_BASE
	this.cnt = -INVINCIBLE_CNT
	this.aPressed = true
	this.bPressed = true
	this.stickInput = this.twinStick.getNullState()
}

func (this *Boat) restart() {
	switch this.gameMode {
	case TWIN_STICK, DOUBLE_PLAY:
		this.fireCnt = 0
		this.fireInterval = FIRE_INTERVAL
	}
	this.fireSprCnt = 0
	this.fireSprDeg = 0.5
	this.fireLanceCnt = 0
	if this.field.getBlock(this._pos) >= 0 {
		this.onBlock = true
	} else {
		this.onBlock = false
	}
	this.refVel.X = 0
	this.refVel.Y = 0
	this.shieldCnt = 20 * 60
}

func (this *Boat) move() {
	px := this._pos.X
	py := this._pos.Y
	this.cnt++
	vx := 0
	vy := 0
	switch this.gameMode {
	case TWIN_STICK:
		this.moveTwinStick()
	case DOUBLE_PLAY:
		this.moveDoublePlay()
	}
	if this.gameState.isGameOver {
		this.clearBullets()
		if this.cnt < -INVINCIBLE_CNT {
			cnt = -RESTART_CNT
		}
	} else if this.cnt < -INVINCIBLE_CNT {
		this.clearBullets()
	}
	vx *= speed
	vy *= speed
	vx += refVel.x
	vy += refVel.y
	refVel *= 0.9
	if this.field.checkInField(this._pos.X, this._pos.Y-this.field.lastScrollY) {
		_pos.y -= field.lastScrollY
	}
	if (this.onBlock || this.field.getBlock(this._pos.X+vx, this._pos.Y) < 0) &&
		this.field.checkInField(this._pos.X+vx, _pos.y) {
		this._pos.X += vx
		this._vel.X = vx
	} else {
		this._vel.X = 0
		this.refVel.X = 0
	}
	srf := false
	if (this.onBlock || this.field.getBlock(px, this._pos.Y+vy) < 0) &&
		this.field.checkInField(this._pos.X, this._pos.Y+vy) {
		this._pos.Y += vy
		this._vel.Y = vy
	} else {
		this._vel.Y = 0
		this.refVel.Y = 0
	}
	if this.field.getBlock(this._pos.X, this._pos.Y) >= 0 {
		if !this.onBlock {
			if cnt <= 0 {
				onBlock = true
			} else {
				if this.field.checkInField(this._pos.X, this._pos.Y-this.field.lastScrollY) {
					this._pos.X = px
					this._pos.Y = py
				} else {
					this.destroyed()
				}
			}
		}
	} else {
		this.onBlock = false
	}
	switch this.gameMode {
	case TWIN_STICK:
		this.fireTwinStick()
	case DOUBLE_PLAY:
		this.fireDoublePlay()
	}
	if this.cnt%3 == 0 && this.cnt >= -INVINCIBLE_CNT {
		var sp float64
		if vx != 0 || vy != 0 {
			sp = 0.4
		} else {
			sp = 0.2
		}
		sp *= 1 + rand.nextSignedFloat(0.33)
		sp *= SPEED_BASE
		this._shape.addWake(this.wakes, this._pos, this.deg, sp)
	}
	he := enemies.checkHitShip(this.pos.X, this.pos.Y)
	if he != nil {
		var rd float
		if this.pos.dist(he.pos) < 0.1 {
			rd = 0
		} else {
			rd = math.Atan2(this._pos.X-he.pos.X, _pos.Y-he.pos.Y)
		}
		sz := he.size
		this.refVel.X = math.Sin(rd) * sz * 0.1
		this.refVel.Y = math.Cos(rd) * sz * 0.1
		rs := this.refVel.vctSize
		if rs > 1 {
			this.refVel.X /= rs
			this.refVel.Y /= rs
		}
	}
	if this.shieldCnt > 0 {
		shieldCnt--
	}
}

func (this *Boat) moveTwinStick() {
	if !this._replayMode {
		this.stickInput = this.twinStick.getState()
	} else {
		i, err := twinStick.replay()
		this.stickInput = i
		if err != nil {
			fmt.Printf("warn : %s", err.Error())
			this.gameState.isGameOver = true
			this.stickInput = this.twinStick.getNullState()
		}
	}
	if this.gameState.isGameOver || this.cnt < -INVINCIBLE_CNT {
		stickInput.clear()
	}
	vx = this.stickInput.left.X
	vy = this.stickInput.left.Y
	if vx != 0 || vy != 0 {
		ad := math.Atan2(vx, vy)
		ad = normalizeDeg(ad)
		ad -= this.deg
		ad = normalizeDeg(ad)
		this.deg += ad * this.turnRatio * this.turnSpeed
		this.deg = normalizeDeg(this.deg)
	}
}

func (this *Boat) moveDoublePlay() {
	switch this.idx {
	case 0:
		if !this._replayMode {
			this.stickInput = this.twinStick.getState()
		} else {
			i, err := twinStick.replay()
			this.stickInput = i
			if err != nil {
				this.gameState.isGameOver = true
				this.stickInput = twinStick.getNullState()
			}
		}
		if this.gameState.isGameOver || this.cnt < -INVINCIBLE_CNT {
			stickInput.clear()
		}
		vx = this.stickInput.left.X
		vy = this.stickInput.left.Y
	case 1:
		vx = this.stickInput.right.X
		vy = this.stickInput.right.Y
	}
	if vx != 0 || vy != 0 {
		ad := math.Atan2(vx, vy)
		ad = normalizeDeg(ad)
		ad -= this.deg
		ad = normalizeDeg(ad)
		this.deg += ad * this.turnRatio * this.turnSpeed
		this.deg = normalizeDeg(this.deg)
	}
}

func (this *Boat) fireTwinStick() {
	if math.Abs(this.stickInput.right.X)+math.Abs(this.stickInput.right.Y) > 0.01 {
		this.fireDeg = math.Atan2(this.stickInput.right.X, this.stickInput.right.Y)
		if this.fireCnt <= 0 {
			playSe("shot.wav")
			foc := (fireSprCnt%2)*2 - 1
			rsd := this.stickInput.right.vctSize
			if rsd > 1 {
				rsd = 1
			}
			this.fireSprDeg = 1 - rsd + 0.05
			this.firePos.X = this._pos.X + math.Cos(this.fireDeg+math.Pi)*0.2*foc
			this.firePos.Y = this._pos.Y - math.Sin(this.fireDeg+math.Pi)*0.2*foc
			this.fireCnt = int(fireInterval)
			var td float
			switch foc {
			case -1:
				td = this.fireSprDeg * (this.fireSprCnt/2%4 + 1) * 0.2
			case 1:
				td = -this.fireSprDeg * (this.fireSprCnt/2%4 + 1) * 0.2
			}
			this.fireSprCnt++
			s := shots.getInstance()
			if s != nil {
				s.set(this.firePos, this.fireDeg+td/2, false, 2)
			}
			s = shots.getInstance()
			if s != nil {
				s.set(this.firePos, this.fireDeg+td, false, 2)
			}
			sm := smokes.getInstanceForced()
			sd := this.fireDeg + td/2
			sm.set(this.firePos, math.Sin(sd)*SPEED*0.33, math.Cos(sd)*SPEED*0.33, 0, SPARK, 10, 0.33)
		}
	} else {
		this.fireDeg = 99999
	}
	this.fireCnt--
}

func (this *Boat) fireDoublePlay() {
	if this.gameState.isGameOver || this.cnt < -INVINCIBLE_CNT {
		return
	}
	dist := ship.distAmongBoats()
	this.fireInterval = FIRE_INTERVAL + 10.0/(dist+0.005)
	if dist < 2 {
		this.fireInterval = 99999
	} else if dist < 4 {
		this.fireInterval *= 3
	} else if dist < 6 {
		this.fireInterval *= 1.6
	}
	if this.fireCnt > this.fireInterval {
		this.fireCnt = int(fireInterval)
	}
	if this.fireCnt <= 0 {
		playSe("shot.wav")
		foc := int(math.Mod(fireSprCnt, 2))*2 - 1
		this.fireDeg = 0
		this.firePos.X = this._pos.X + math.Cos(this.fireDeg+math.Pi)*0.2*foc
		this.firePos.Y = this._pos.Y - math.Sin(this.fireDeg+math.Pi)*0.2*foc
		s := shots.getInstance()
		if s != nil {
			s.set(this.firePos, this.fireDeg, false, 2)
		}
		this.fireCnt = int(fireInterval)
		sm := this.smokes.getInstanceForced()
		sd := this.fireDeg
		sm.set(this.firePos, math.Sin(sd)*SPEED*0.33, math.Cos(sd)*SPEED*0.33, 0, SPARK, 10, 0.33)
		if this.idx == 0 {
			fd := ship.degAmongBoats() + math.Pi/2
			var td float
			switch foc {
			case -1:
				td = this.fireSprDeg * (this.fireSprCnt/math.Mod(2, 4) + 1) * 0.15
			case 1:
				td = -this.fireSprDeg * (this.fireSprCnt/math.Mod(2, 4) + 1) * 0.15
			}
			this.firePos.x = this.ship.midstPos.X + math.Cos(fd+math.Pi)*0.2*foc
			this.firePos.y = this.ship.midstPos.Y - math.Sin(fd+math.Pi)*0.2*foc
			s = shots.getInstance()
			if s != nil {
				s.set(this.firePos, fd, false, 2)
			}
			s = shots.getInstance()
			if s != nil {
				s.set(this.firePos, fd+td, false, 2)
			}
			sm = smokes.getInstanceForced()
			sm.set(this.firePos, math.Sin(fd+td/2)*SPEED*0.33, math.Cos(fd+td/2)*SPEED*0.33, 0,
				SPARK, 10, 0.33)
		}
		this.fireSprCnt++
	}
	this.fireCnt--
}

func (this *Boat) checkBulletHit(p, pp vector.Vector) bool {
	if this.cnt <= 0 {
		return false
	}
	var bmvx, bmvy, inaa float64
	bmvx = pp.x
	bmvy = pp.y
	bmvx -= p.x
	bmvy -= p.y
	inaa = bmvx*bmvx + bmvy*bmvy
	if inaa > 0.00001 {
		var sofsx, sofsy, inab, hd float64
		sofsx = _pos.x
		sofsy = _pos.y
		sofsx -= p.x
		sofsy -= p.y
		inab = bmvx*sofsx + bmvy*sofsy
		if inab >= 0 && inab <= inaa {
			hd = sofsx*sofsx + sofsy*sofsy - inab*inab/inaa
			if hd >= 0 && hd <= HIT_WIDTH {
				this.destroyed()
				return true
			}
		}
	}
	return false
}

func (this *Boat) destroyed() {
	if this.cnt <= 0 {
		return
	}
	if this.shieldCnt > 0 {
		this.destroyedBoatShield()
		return
	}
	this.ship.destroyed()
	this.gameState.shipDestroyed()
}

func (this *Boat) destroyedBoatShield() {
	for i := 0; i < 100; i++ {
		sp := this.sparks.getInstanceForced()
		sp.set(this.pos, rand.nextSignedFloat(1), rand.nextSignedFloat(1),
			0.5+rand.nextFloat(0.5), 0.5+rand.nextFloat(0.5), 0,
			40+rand.nextInt(40))
	}
	playSe("ship_shield_lost.wav")
	sdl.setScreenShake(30, 0.02)
	this.shieldCnt = 0
	this.cnt = -INVINCIBLE_CNT / 2
}

func (this *Boat) destroyedBoat() {
	for i := 0; i < 128; i++ {
		sp := sparks.getInstanceForced()
		sp.set(this.pos, rand.nextSignedFloat(1), rand.nextSignedFloat(1),
			0.5+rand.nextFloat(0.5), 0.5+rand.nextFloat(0.5), 0,
			40+rand.nextInt(40))
	}
	playSe("ship_destroyed.wav")
	for i := 0; i < 64; i++ {
		s := smokes.getInstanceForced()
		s.set(this.pos, rand.nextSignedFloat(0.2), rand.nextSignedFloat(0.2),
			rand.nextFloat(0.1),
			EXPLOSION, 50+rand.nextInt(30), 1)
	}
	sdl.setScreenShake(60, 0.05)
	this.restart()
	this.cnt = -RESTART_CNT
}

func (this *Boat) hasCollision() bool {
	return !(this.cnt < -INVINCIBLE_CNT)
}

func (this *Boat) draw() {
	if this.cnt < -INVINCIBLE_CNT {
		return
	}
	if this.fireDeg < 99999 {
		sdl.setColor(0.5, 0.9, 0.7, 0.4)
		gl.Begin(gl.LINE_STRIP)
		gl.Vertex2(this._pos.X, this._pos.Y)
		sdl.SetColor(0.5, 0.9, 0.7, 0.8)
		gl.Vertex2(this._pos.X+math.Sin(this.fireDeg)*20, this._pos.Y+math.Cos(this.fireDeg)*20)
		gl.End()
	}
	if this.cnt < 0 && (-this.cnt%32) < 16 {
		return
	}
	gl.PushMatrix()
	sdl.glTranslate(pos)
	gl.Rotatef(-deg*180/PI, 0, 0, 1)
	this._shape.draw()
	this.bridgeShape.draw()
	if this.shieldCnt > 0 {
		ss := 0.66
		if this.shieldCnt < 120 {
			ss *= float64(shieldCnt) / 120
		}
		gl.Scalef(ss, ss, ss)
		gl.Rotatef(shieldCnt*5, 0, 0, 1)
		this.shieldShape.draw()
	}
	gl.PopMatrix()
}

func (this *Boat) drawFront() {
	if this.cnt < -INVINCIBLE_CNT {
		return
	}
}

func (this *Boat) drawShape() {
	this._shape.draw()
	this.bridgeShape.draw()
}

func (this *Boat) clearBullets() {
	this.gameState.clearBullets()
}

func (this *Boat) pos() vector.Vector {
	return this._pos
}

func (this *Boat) vel() vector.Vector {
	return this._vel
}

func (this *Boat) setReplayMode(turnSpeed float64, reverseFire bool) {
	this._replayMode = true
	this.turnSpeed = turnSpeed
	this.reverseFire = reverseFire
}

func (this *Boat) unsetReplayMode() {
	this._replayMode = false
	this.turnSpeed = GameManager.shipTurnSpeed
	this.reverseFire = GameManager.shipReverseFire
}

func (this *Boat) replayMode() bool {
	return this._replayMode
}
