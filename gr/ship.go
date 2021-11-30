package main

import (
	"math"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/sdl/record"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Player's ship.
 */

const SHIP_SCROLL_SPEED_BASE = 0.01
const SCROLL_SPEED_MAX = 0.1
const SCROLL_START_Y = 2.5

type Ship struct {
	field                                                *Field
	boat                                                 [2]Boat
	gameMode                                             GameMode
	boatNum                                              int
	gameState                                            *InGameState
	scrollSpeed, _scrollSpeedBase                        float64
	_midstPos, _higherPos, _lowerPos, _nearPos, _nearVel vector.Vector
	bridgeShape                                          BaseShape
}

func NewShip(twinStick *sdl.RecordableTwinStick, field *Field, screen *Screen,
	sparks *SparkPool, smokes *SmokePool, fragments *FragmentPool, wakes *WakePool) *Ship {
	this := &Ship{}
	this.field = field
	for i := range this.boat {
		this.boat[i] = NewBoat(i, this, twinStick,
			field, screen, sparks, smokes, fragments, wakes)
		i++
	}
	this.boatNum = 1
	this.scrollSpeed = SHIP_SCROLL_SPEED_BASE
	this._scrollSpeedBase = SHIP_SCROLL_SPEED_BASE
	this.bridgeShape = NewBaseShapeInternal(0.3, 0.2, 0.1, BRIDGE, 0.3, 0.7, 0.7)
	return this
}

func (this *Ship) setRandSeed(seed int64) {
	setBoatRandSeed(seed)
}

func (this *Ship) setShots(shots *ShotPool) {
	for _, b := range this.boat {
		b.setShots(shots)
	}
}

func (this *Ship) setEnemies(enemies *EnemyPool) {
	for _, b := range this.boat {
		b.setEnemies(enemies)
	}
}

func (this *Ship) setStageManager(stageManager *StageManager) {
	for _, b := range this.boat {
		b.setStageManager(stageManager)
	}
}

func (this *Ship) setGameState(gameState *InGameState) {
	this.gameState = gameState
	for _, b := range this.boat {
		b.setGameState(gameState)
	}
}

func (this *Ship) start(gameMode GameMode) {
	this.gameMode = gameMode
	if gameMode == DOUBLE_PLAY {
		this.boatNum = 2
	} else {
		this.boatNum = 1
	}
	this._scrollSpeedBase = SHIP_SCROLL_SPEED_BASE
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].start(gameMode)
	}
	this._midstPos.X = 0
	this._midstPos.Y = 0
	this._higherPos.X = 0
	this._higherPos.Y = 0
	this._lowerPos.X = 0
	this._lowerPos.Y = 0
	this._nearPos.X = 0
	this._nearPos.Y = 0
	this._nearVel.X = 0
	this._nearVel.Y = 0
	this.restart()
}

func (this *Ship) restart() {
	this.scrollSpeed = this._scrollSpeedBase
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].restart()
	}
}

func (this *Ship) move() {
	this.field.scroll(this.scrollSpeed, false)
	sf := false
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].move()
		if this.boat[i].hasCollision() && this.boat[i].pos().X > this.field.size().X/3 && this.boat[i].pos().Y < -this.field.size().Y/4*3 {
			sf = true
		}
	}
	if sf {
		this.gameState.shrinkScoreReel()
	}
	if this.higherPos().Y >= SCROLL_START_Y {
		this.scrollSpeed += (SCROLL_SPEED_MAX - this.scrollSpeed) * 0.1
	} else {
		this.scrollSpeed += (this._scrollSpeedBase - this.scrollSpeed) * 0.1
	}
	this._scrollSpeedBase += (SCROLL_SPEED_MAX - this._scrollSpeedBase) * 0.00001
}

func (this *Ship) checkBulletHit(p, pp vector.Vector) bool {
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].checkBulletHit(p, pp) {
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
	if this.gameMode == DOUBLE_PLAY && this.boat[0].hasCollision() {
		sdl.SetColor(0.5, 0.5, 0.9, 0.8)
		gl.Begin(gl.LINE_STRIP)
		gl.Vertex2d(this.boat[0].pos().X, this.boat[0].pos().Y)
		sdl.SetColor(0.5, 0.5, 0.9, 0.3)
		gl.Vertex2d(this.midstPos().X, this.midstPos().Y)
		sdl.SetColor(0.5, 0.5, 0.9, 0.8)
		gl.Vertex2d(this.boat[1].pos().X, this.boat[1].pos().Y)
		gl.End()
		gl.PushMatrix()
		sdl.GlTranslate(this.midstPos())
		gl.Rotated(-this.degAmongBoats()*180/math.Pi, 0, 0, 1)
		this.bridgeShape.Draw()
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
	this._midstPos.X = 0
	this._midstPos.Y = 0
	for i := 0; i < this.boatNum; i++ {
		this._midstPos.X += this.boat[i].pos().X
		this._midstPos.Y += this.boat[i].pos().Y
	}
	this._midstPos.OpDivAssign(float64(this.boatNum))
	return this._midstPos
}

func (this *Ship) higherPos() vector.Vector {
	this._higherPos.Y = -99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos().Y > this._higherPos.Y {
			this._higherPos.X = this.boat[i].pos().X
			this._higherPos.Y = this.boat[i].pos().Y
		}
	}
	return this._higherPos
}

func (this *Ship) lowerPos() vector.Vector {
	this._lowerPos.Y = 99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos().Y < this._lowerPos.Y {
			this._lowerPos.X = this.boat[i].pos().X
			this._lowerPos.Y = this.boat[i].pos().Y
		}
	}
	return this._lowerPos
}

func (this *Ship) nearPos(p vector.Vector) vector.Vector {
	dist := 99999.0
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos().DistVector(p) < dist {
			dist = this.boat[i].pos().DistVector(p)
			this._nearPos.X = this.boat[i].pos().X
			this._nearPos.Y = this.boat[i].pos().Y
		}
	}
	return this._nearPos
}

func (this *Ship) nearVel(p vector.Vector) vector.Vector {
	dist := 99999.0
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos().DistVector(p) < dist {
			dist = this.boat[i].pos().DistVector(p)
			this._nearVel.X = this.boat[i].vel().X
			this._nearVel.Y = this.boat[i].vel().Y
		}
	}
	return this._nearVel
}

func (this *Ship) distAmongBoats() float64 {
	return this.boat[0].pos().DistVector(this.boat[1].pos())
}

func (this *Ship) degAmongBoats() float64 {
	if this.distAmongBoats() < 0.1 {
		return 0
	} else {
		return math.Atan2(this.boat[0].pos().X-this.boat[1].pos().X, this.boat[0].pos().Y-this.boat[1].pos().Y)
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
var stickInput sdl.TwinStickState

type Boat struct {
	twinStick                *sdl.RecordableTwinStick
	field                    *Field
	screen                   *Screen
	shots                    *ShotPool
	sparks                   *SparkPool
	smokes                   *SmokePool
	fragments                *FragmentPool
	wakes                    *WakePool
	enemies                  *EnemyPool
	stageManager             *StageManager
	gameState                *InGameState
	_pos, firePos            vector.Vector
	deg, speed, turnRatio    float64
	_shape                   *BaseShape
	bridgeShape              *BaseShape
	fireCnt, fireSprCnt      int
	fireInterval, fireSprDeg float64
	fireLanceCnt             int
	fireDeg                  float64
	aPressed, bPressed       bool
	cnt                      int
	onBlock                  bool
	_vel, refVel             vector.Vector
	shieldCnt                int
	shieldShape              *ShieldShape
	_replayMode              bool
	turnSpeed                float64
	reverseFire              bool
	gameMode                 GameMode
	vx, vy                   float64
	idx                      int
	ship                     *Ship
}

func setBoatRandSeed(seed int64) {
	boatRand = r.New(r.NewSource(seed))
}

func NewBoat(idx int, ship *Ship,
	twinStick *sdl.RecordableTwinStick,
	field *Field, screen *Screen,
	sparks *SparkPool, smokes *SmokePool, fragments *FragmentPool, wakes *WakePool) Boat {
	this := Boat{}
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
		this._shape = NewBaseShape(0.7, 0.6, 0.6, SHIP_ROUNDTAIL, 0.5, 0.7, 0.5)
		this.bridgeShape = NewBaseShape(0.3, 0.6, 0.6, BRIDGE, 0.3, 0.7, 0.3)
	case 1:
		this._shape = NewBaseShape(0.7, 0.6, 0.6, SHIP_ROUNDTAIL, 0.4, 0.3, 0.8)
		this.bridgeShape = NewBaseShape(0.3, 0.6, 0.6, BRIDGE, 0.2, 0.3, 0.6)
	}
	this.turnSpeed = 1
	this.fireInterval = FIRE_INTERVAL
	this.shieldShape = NewShieldShape()
	return this
}

func (this *Boat) setShots(shots *ShotPool) {
	this.shots = shots
}

func (this *Boat) setEnemies(enemies *EnemyPool) {
	this.enemies = enemies
}

func (this *Boat) setStageManager(stageManager *StageManager) {
	this.stageManager = stageManager
}

func (this *Boat) setGameState(gameState *InGameState) {
	this.gameState = gameState
}

func (this *Boat) start(gameMode GameMode) {
	this.gameMode = gameMode
	if gameMode == DOUBLE_PLAY {
		switch this.idx {
		case 0:
			this._pos.X = -this.field.size().X * 0.5
		case 1:
			this._pos.X = this.field.size().X * 0.5
		}
	} else {
		this._pos.X = 0
	}
	this._pos.Y = -this.field.size().Y * 0.8
	this.firePos.X = 0
	this.firePos.Y = 0
	this._vel.X = 0
	this._vel.Y = 0
	this.deg = 0
	this.speed = SPEED_BASE
	this.turnRatio = TURN_RATIO_BASE
	this.cnt = -INVINCIBLE_CNT
	this.aPressed = true
	this.bPressed = true
	stickInput = this.twinStick.GetNullState()
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
	if this.field.getBlockVector(this._pos) >= 0 {
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
	this.vx = 0
	this.vy = 0
	switch this.gameMode {
	case TWIN_STICK:
		this.moveTwinStick()
	case DOUBLE_PLAY:
		this.moveDoublePlay()
	}
	if this.gameState.isGameOver {
		this.clearBullets()
		if this.cnt < -INVINCIBLE_CNT {
			this.cnt = -RESTART_CNT
		}
	} else if this.cnt < -INVINCIBLE_CNT {
		this.clearBullets()
	}
	this.vx *= this.speed
	this.vy *= this.speed
	this.vx += this.refVel.X
	this.vy += this.refVel.Y
	this.refVel.OpMulAssign(0.9)
	if this.field.checkInField(this._pos.X, this._pos.Y-this.field.lastScrollY()) {
		this._pos.Y -= this.field.lastScrollY()
	}
	if (this.onBlock || this.field.getBlock(this._pos.X+this.vx, this._pos.Y) < 0) &&
		this.field.checkInField(this._pos.X+this.vx, this._pos.Y) {
		this._pos.X += this.vx
		this._vel.X = this.vx
	} else {
		this._vel.X = 0
		this.refVel.X = 0
	}
	// srf := false
	if (this.onBlock || this.field.getBlock(px, this._pos.Y+this.vy) < 0) &&
		this.field.checkInField(this._pos.X, this._pos.Y+this.vy) {
		this._pos.Y += this.vy
		this._vel.Y = this.vy
	} else {
		this._vel.Y = 0
		this.refVel.Y = 0
	}
	if this.field.getBlock(this._pos.X, this._pos.Y) >= 0 {
		if !this.onBlock {
			if this.cnt <= 0 {
				this.onBlock = true
			} else {
				if this.field.checkInField(this._pos.X, this._pos.Y-this.field.lastScrollY()) {
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
		if this.vx != 0 || this.vy != 0 {
			sp = 0.4
		} else {
			sp = 0.2
		}
		sp *= 1 + nextSignedFloat(boatRand, 0.33)
		sp *= SPEED_BASE
		this._shape.addWake(this.wakes, this._pos, this.deg, sp, 1)
	}
	he := this.enemies.checkHitShip(this.pos().X, this.pos().Y, nil, false)
	if he != nil {
		var rd float64
		if this.pos().DistVector(he.pos()) < 0.1 {
			rd = 0
		} else {
			rd = math.Atan2(this._pos.X-he.pos().X, this._pos.Y-he.pos().Y)
		}
		sz := he.size()
		this.refVel.X = math.Sin(rd) * sz * 0.1
		this.refVel.Y = math.Cos(rd) * sz * 0.1
		rs := this.refVel.VctSize()
		if rs > 1 {
			this.refVel.X /= rs
			this.refVel.Y /= rs
		}
	}
	if this.shieldCnt > 0 {
		this.shieldCnt--
	}
}

func (this *Boat) moveTwinStick() {
	if !this._replayMode {
		stickInput = this.twinStick.GetState()
	} else {
		var err error
		stickInput, err = twinStick.Replay()
		if err != nil {
			if err == record.EndRecordingErr {
				this.gameState.isGameOver = true
				stickInput = this.twinStick.GetNullState()
			} else {
				panic(err)
			}
		}
	}
	if this.gameState.isGameOver || this.cnt < -INVINCIBLE_CNT {
		stickInput.Clear()
	}
	this.vx = stickInput.Left.X
	this.vy = stickInput.Left.Y
	if this.vx != 0 || this.vy != 0 {
		ad := math.Atan2(this.vx, this.vy)
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
			stickInput = this.twinStick.GetState()
		} else {
			var err error
			stickInput, err = twinStick.Replay()
			if err != nil {
				this.gameState.isGameOver = true
				stickInput = twinStick.GetNullState()
			}
		}
		if this.gameState.isGameOver || this.cnt < -INVINCIBLE_CNT {
			stickInput.Clear()
		}
		this.vx = stickInput.Left.X
		this.vy = stickInput.Left.Y
	case 1:
		this.vx = stickInput.Right.X
		this.vy = stickInput.Right.Y
	}
	if this.vx != 0 || this.vy != 0 {
		ad := math.Atan2(this.vx, this.vy)
		ad = normalizeDeg(ad)
		ad -= this.deg
		ad = normalizeDeg(ad)
		this.deg += ad * this.turnRatio * this.turnSpeed
		this.deg = normalizeDeg(this.deg)
	}
}

func (this *Boat) fireTwinStick() {
	if math.Abs(stickInput.Right.X)+math.Abs(stickInput.Right.Y) > 0.01 {
		this.fireDeg = math.Atan2(stickInput.Right.X, stickInput.Right.Y)
		if this.fireCnt <= 0 {
			playSe("shot.wav")
			foc := (this.fireSprCnt%2)*2 - 1
			rsd := stickInput.Right.VctSize()
			if rsd > 1 {
				rsd = 1
			}
			this.fireSprDeg = 1 - rsd + 0.05
			this.firePos.X = this._pos.X + math.Cos(this.fireDeg+math.Pi)*0.2*float64(foc)
			this.firePos.Y = this._pos.Y - math.Sin(this.fireDeg+math.Pi)*0.2*float64(foc)
			this.fireCnt = int(this.fireInterval)
			var td float64
			switch foc {
			case -1:
				td = this.fireSprDeg * float64(this.fireSprCnt/2%4+1) * 0.2
			case 1:
				td = -this.fireSprDeg * float64(this.fireSprCnt/2%4+1) * 0.2
			}
			this.fireSprCnt++
			s := this.shots.GetInstance()
			if s != nil {
				s.set(this.firePos, this.fireDeg+td/2, false, 2)
			}
			s = this.shots.GetInstance()
			if s != nil {
				s.set(this.firePos, this.fireDeg+td, false, 2)
			}
			sm := this.smokes.GetInstanceForced()
			sd := this.fireDeg + td/2
			sm.setVector(this.firePos, math.Sin(sd)*SPEED*0.33, math.Cos(sd)*SPEED*0.33, 0, SPARK, 10, 0.33)
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
	dist := this.ship.distAmongBoats()
	this.fireInterval = FIRE_INTERVAL + 10.0/(dist+0.005)
	if dist < 2 {
		this.fireInterval = 99999
	} else if dist < 4 {
		this.fireInterval *= 3
	} else if dist < 6 {
		this.fireInterval *= 1.6
	}
	if float64(this.fireCnt) > this.fireInterval {
		this.fireCnt = int(this.fireInterval)
	}
	if this.fireCnt <= 0 {
		playSe("shot.wav")
		foc := (this.fireSprCnt%2)*2 - 1
		this.fireDeg = 0
		this.firePos.X = this._pos.X + math.Cos(this.fireDeg+math.Pi)*0.2*float64(foc)
		this.firePos.Y = this._pos.Y - math.Sin(this.fireDeg+math.Pi)*0.2*float64(foc)
		s := this.shots.GetInstance()
		if s != nil {
			s.set(this.firePos, this.fireDeg, false, 2)
		}
		this.fireCnt = int(this.fireInterval)
		sm := this.smokes.GetInstanceForced()
		sd := this.fireDeg
		sm.setVector(this.firePos, math.Sin(sd)*SPEED*0.33, math.Cos(sd)*SPEED*0.33, 0, SPARK, 10, 0.33)
		if this.idx == 0 {
			fd := this.ship.degAmongBoats() + math.Pi/2
			var td float64
			switch foc {
			case -1:
				td = this.fireSprDeg * float64(this.fireSprCnt/2%4+1) * 0.15
			case 1:
				td = -this.fireSprDeg * float64(this.fireSprCnt/2%4+1) * 0.15
			}
			this.firePos.X = this.ship.midstPos().X + math.Cos(fd+math.Pi)*0.2*float64(foc)
			this.firePos.Y = this.ship.midstPos().Y - math.Sin(fd+math.Pi)*0.2*float64(foc)
			s = this.shots.GetInstance()
			if s != nil {
				s.set(this.firePos, fd, false, 2)
			}
			s = this.shots.GetInstance()
			if s != nil {
				s.set(this.firePos, fd+td, false, 2)
			}
			sm = this.smokes.GetInstanceForced()
			sm.setVector(this.firePos, math.Sin(fd+td/2)*SPEED*0.33, math.Cos(fd+td/2)*SPEED*0.33, 0,
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
	bmvx = pp.X
	bmvy = pp.Y
	bmvx -= p.X
	bmvy -= p.Y
	inaa = bmvx*bmvx + bmvy*bmvy
	if inaa > 0.00001 {
		var sofsx, sofsy, inab, hd float64
		sofsx = this._pos.X
		sofsy = this._pos.Y
		sofsx -= p.X
		sofsy -= p.Y
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
		sp := this.sparks.GetInstanceForced()
		sp.set(this.pos(), nextSignedFloat(boatRand, 1), nextSignedFloat(boatRand, 1),
			0.5+nextFloat(boatRand, 0.5), 0.5+nextFloat(boatRand, 0.5), 0,
			40+boatRand.Intn(40))
	}
	playSe("ship_shield_lost.wav")
	screen.setScreenShake(30, 0.02)
	this.shieldCnt = 0
	this.cnt = -INVINCIBLE_CNT / 2
}

func (this *Boat) destroyedBoat() {
	for i := 0; i < 128; i++ {
		sp := this.sparks.GetInstanceForced()
		sp.set(this.pos(), nextSignedFloat(boatRand, 1), nextSignedFloat(boatRand, 1),
			0.5+nextFloat(boatRand, 0.5), 0.5+nextFloat(boatRand, 0.5), 0,
			40+boatRand.Intn(40))
	}
	playSe("ship_destroyed.wav")
	for i := 0; i < 64; i++ {
		s := this.smokes.GetInstanceForced()
		s.setVector(this.pos(), nextSignedFloat(boatRand, 0.2), nextSignedFloat(boatRand, 0.2),
			nextFloat(boatRand, 0.1),
			EXPLOSION, 50+boatRand.Intn(30), 1)
	}
	screen.setScreenShake(60, 0.05)
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
		sdl.SetColor(0.5, 0.9, 0.7, 0.4)
		gl.Begin(gl.LINE_STRIP)
		gl.Vertex2d(this._pos.X, this._pos.Y)
		sdl.SetColor(0.5, 0.9, 0.7, 0.8)
		gl.Vertex2d(this._pos.X+math.Sin(this.fireDeg)*20, this._pos.Y+math.Cos(this.fireDeg)*20)
		gl.End()
	}
	if this.cnt < 0 && (-this.cnt%32) < 16 {
		return
	}
	gl.PushMatrix()
	sdl.GlTranslate(this.pos())
	gl.Rotated(-this.deg*180/math.Pi, 0, 0, 1)
	this._shape.Draw()
	this.bridgeShape.Draw()
	if this.shieldCnt > 0 {
		ss := 0.66
		if this.shieldCnt < 120 {
			ss *= float64(this.shieldCnt) / 120
		}
		gl.Scaled(ss, ss, ss)
		gl.Rotated(float64(this.shieldCnt)*5, 0, 0, 1)
		this.shieldShape.Draw()
	}
	gl.PopMatrix()
}

func (this *Boat) drawFront() {
	if this.cnt < -INVINCIBLE_CNT {
		return
	}
}

func (this *Boat) drawShape() {
	this._shape.Draw()
	this.bridgeShape.Draw()
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
	this.turnSpeed = shipTurnSpeed
	this.reverseFire = shipReverseFire
}

func (this *Boat) replayMode() bool {
	return this._replayMode
}
