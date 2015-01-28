/*
 * $Id: ship.d,v 1.4 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/go-gl/gl"
)

/**
 * Player's ship.
 */

const SCROLL_SPEED_BASE = 0.01
const SCROLL_SPEED_MAX = 0.1
const SCROLL_START_Y = 2.5

var ship *Ship

type Ship struct {
	livesLeft                                            int
	boat                                                 [2]*Boat
	gameMode                                             GameMode
	boatNum                                              int
	scrollSpeed, scrollSpeedBase                         float32
	_midstPos, _higherPos, _lowerPos, _nearPos, _nearVel Vector
	bridgeShape                                          *ComplexShape
}

var shipBridgeShape *ComplexShape

func InitShip() {
	shipBridgeShape = NewComplexShape(0.3, 0.2, 0.1, BRIDGE, 0.3, 0.7, 0.7, false)
}

func NewShip() *Ship {
	this := new(Ship)
	this.livesLeft = 2
	for i, _ := range this.boat {
		this.boat[i] = NewBoat(i)
	}
	this.boatNum = 1
	this.scrollSpeed = SCROLL_SPEED_BASE
	this.scrollSpeedBase = SCROLL_SPEED_BASE
	this.bridgeShape = shipBridgeShape
	actors[this] = true
	return this
}

func (this *Ship) close() {
	for _, b := range this.boat {
		b.close()
	}
	delete(actors, this)
}

func (this *Ship) start(gameMode GameMode) {
	this.gameMode = gameMode
	if gameMode == GameModeDOUBLE_PLAY {
		this.boatNum = 2
	} else {
		this.boatNum = 1
	}
	this.scrollSpeedBase = SCROLL_SPEED_BASE
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].start(this.gameMode)
	}
	this._midstPos.y = 0
	this._midstPos.x = 0
	this._higherPos.x = 0
	this._higherPos.y = 0
	this._lowerPos.x = 0
	this._lowerPos.y = 0
	this._nearPos.y = 0
	this._nearPos.x = 0
	this._nearVel.y = 0
	this._nearVel.x = 0
	this.restart()
}

func (this *Ship) restart() {
	this.scrollSpeed = this.scrollSpeedBase
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].restart()
	}
}

func (this *Ship) move() {
	field.scroll(this.scrollSpeed, false)
	sf := false
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].move()
		if this.boat[i].hasCollision() &&
			this.boat[i].pos.x > field.size.x/3 && this.boat[i].pos.y < -field.size.y/4*3 {
			sf = true
		}
	}
	if sf {
		inGameState.shrinkScoreReel()
	}
	if this.higherPos().y >= SCROLL_START_Y {
		this.scrollSpeed += (SCROLL_SPEED_MAX - this.scrollSpeed) * 0.1
	} else {
		this.scrollSpeed += (this.scrollSpeedBase - this.scrollSpeed) * 0.1
	}
	this.scrollSpeedBase += (SCROLL_SPEED_MAX - this.scrollSpeedBase) * 0.00001
}

func (this *Ship) checkBulletHit(p Vector, pp Vector) bool {
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].checkBulletHit(p, pp) {
			return true
		}
	}
	return false
}

func (this *Ship) clearBullets() {
	clearBullets()
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
	if this.gameMode == GameModeDOUBLE_PLAY && this.boat[0].hasCollision() {
		setScreenColor(0.5, 0.5, 0.9, 0.8)
		gl.Begin(gl.LINE_STRIP)
		gl.Vertex2f(this.boat[0].pos.x, this.boat[0].pos.y)
		setScreenColor(0.5, 0.5, 0.9, 0.3)
		gl.Vertex2f(this.midstPos().x, this.midstPos().y)
		setScreenColor(0.5, 0.5, 0.9, 0.8)
		gl.Vertex2f(this.boat[1].pos.x, this.boat[1].pos.y)
		gl.End()
		gl.PushMatrix()
		glTranslate(this.midstPos())
		gl.Rotatef(-this.degAmongBoats()*180/Pi32, 0, 0, 1)
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

func (this *Ship) midstPos() Vector {
	this._midstPos.x = 0
	this._midstPos.y = 0
	for i := 0; i < this.boatNum; i++ {
		this._midstPos.x += this.boat[i].pos.x
		this._midstPos.y += this.boat[i].pos.y
	}
	this._midstPos.DivAssign(float32(this.boatNum))
	return this._midstPos
}

func (this *Ship) higherPos() Vector {
	this._higherPos.y = -99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos.y > this._higherPos.y {
			this._higherPos.x = this.boat[i].pos.x
			this._higherPos.y = this.boat[i].pos.y
		}
	}
	return this._higherPos
}

func (this *Ship) lowerPos() Vector {
	this._lowerPos.y = 99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos.y < this._lowerPos.y {
			this._lowerPos.x = this.boat[i].pos.x
			this._lowerPos.y = this.boat[i].pos.y
		}
	}
	return this._lowerPos
}

func (this *Ship) nearPos(p Vector) Vector {
	var dist float32 = 99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos.distVector(p) < dist {
			dist = this.boat[i].pos.distVector(p)
			this._nearPos.x = this.boat[i].pos.x
			this._nearPos.y = this.boat[i].pos.y
		}
	}
	return this._nearPos
}

func (this *Ship) nearVel(p Vector) Vector {
	var dist float32 = 99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos.distVector(p) < dist {
			dist = this.boat[i].pos.distVector(p)
			this._nearVel.x = this.boat[i].vel.x
			this._nearVel.y = this.boat[i].vel.y
		}
	}
	return this._nearVel
}

func (this *Ship) distAmongBoats() float32 {
	return this.boat[0].pos.distVector(this.boat[1].pos)
}

func (this *Ship) degAmongBoats() float32 {
	if this.distAmongBoats() < 0.1 {
		return 0
	} else {
		return atan232(this.boat[0].pos.x-this.boat[1].pos.x, this.boat[0].pos.y-this.boat[1].pos.y)
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
const SLOW_TURN_RATIO = 0
const TURN_CHANGE_RATIO = 0.5

var padInput PadState
var mouseInput MouseState
var stickInput TwinStickState

type Boat struct {
	pos                      Vector
	firePos                  Vector
	deg, speed, turnRatio    float32
	shape                    *ComplexShape
	bridgeShape              *ComplexShape
	fireCnt, fireSprCnt      int
	fireInterval, fireSprDeg float32
	fireLanceCnt             int
	fireDeg                  float32
	aPressed, bPressed       bool
	cnt                      int
	onBlock                  bool
	vel                      Vector
	refVel                   Vector
	shieldCnt                int
	shieldShape              *ShieldShape
	turnSpeed                float32
	reverseFire              bool
	gameMode                 GameMode
	vx, vy                   float32
	idx                      int
}

var boatShapes []*ComplexShape
var shieldShape *ShieldShape

func InitBoats() {
	boatShapes = make([]*ComplexShape, 4, 4)
	boatShapes[0] = NewComplexShape(0.7, 0.6, 0.6, SHIP_ROUNDTAIL, 0.5, 0.7, 0.5, false)
	boatShapes[1] = NewComplexShape(0.3, 0.6, 0.6, BRIDGE, 0.3, 0.7, 0.3, false)
	boatShapes[2] = NewComplexShape(0.7, 0.6, 0.6, SHIP_ROUNDTAIL, 0.4, 0.3, 0.8, false)
	boatShapes[3] = NewComplexShape(0.3, 0.6, 0.6, BRIDGE, 0.2, 0.3, 0.6, false)
	shieldShape = NewShieldShape()
}

func NewBoat(idx int) *Boat {
	this := new(Boat)
	this.idx = idx
	switch idx {
	case 0:
		this.shape = boatShapes[0]
		this.bridgeShape = boatShapes[1]
	case 1:
		this.shape = boatShapes[2]
		this.bridgeShape = boatShapes[3]
	}
	this.turnSpeed = 1
	this.fireInterval = FIRE_INTERVAL
	this.shieldShape = shieldShape
	return this
}

func (this *Boat) close() {
}

func closeBoats() {
	shipBridgeShape.close()
	shieldShape.close()
	for _, s := range boatShapes {
		s.close()
	}
}

func (this *Boat) start(gameMode GameMode) {
	this.gameMode = gameMode
	if gameMode == GameModeDOUBLE_PLAY {
		switch this.idx {
		case 0:
			this.pos.x = -field.size.x * 0.5
		case 1:
			this.pos.x = field.size.x * 0.5
		}
	} else {
		this.pos.x = 0
	}
	this.pos.y = -field.size.y * 0.8
	this.firePos.x = 0
	this.firePos.y = 0
	this.vel.x = 0
	this.vel.y = 0
	this.deg = 0
	this.speed = SPEED_BASE
	this.turnRatio = TURN_RATIO_BASE
	this.cnt = -INVINCIBLE_CNT
	this.aPressed = true
	this.bPressed = true
	padInput = pad.getNullState()
	stickInput = twinStick.getNullState()
	mouseInput = mouse.getNullState()
}

func (this *Boat) restart() {
	switch this.gameMode {
	case GameModeNORMAL:
		this.fireCnt = 99999
		this.fireInterval = 99999
	case GameModeTWIN_STICK, GameModeDOUBLE_PLAY, GameModeMOUSE:
		this.fireCnt = 0
		this.fireInterval = FIRE_INTERVAL
	}
	this.fireSprCnt = 0
	this.fireSprDeg = 0.5
	this.fireLanceCnt = 0
	if field.getBlockVector(this.pos) >= 0 {
		this.onBlock = true
	} else {
		this.onBlock = false
	}
	this.refVel.x = 0
	this.refVel.y = 0
	this.shieldCnt = 20 * 60
}

func (this *Boat) move() {
	px := this.pos.x
	py := this.pos.y
	this.cnt++
	this.vx = 0
	this.vy = 0
	switch this.gameMode {
	case GameModeNORMAL:
		this.moveNormal()
	case GameModeTWIN_STICK:
		this.moveTwinStick()
	case GameModeDOUBLE_PLAY:
		this.moveDoublePlay()
	case GameModeMOUSE:
		this.moveMouse()
	}
	this.handleGameOver()
	this.vx *= this.speed
	this.vy *= this.speed
	this.vx += this.refVel.x
	this.vy += this.refVel.y
	this.refVel.MulAssign(0.9)
	if field.checkInField(this.pos.x, this.pos.y-field.lastScrollY) {
		this.pos.y -= field.lastScrollY
	}
	if (this.onBlock || field.getBlock(this.pos.x+this.vx, this.pos.y) < 0) &&
		field.checkInField(this.pos.x+this.vx, this.pos.y) {
		this.pos.x += this.vx
		this.vel.x = this.vx
	} else {
		this.vel.x = 0
		this.refVel.x = 0
	}
	if (this.onBlock || field.getBlock(px, this.pos.y+this.vy) < 0) &&
		field.checkInField(this.pos.x, this.pos.y+this.vy) {
		this.pos.y += this.vy
		this.vel.y = this.vy
	} else {
		this.vel.y = 0
		this.refVel.y = 0
	}
	if field.getBlock(this.pos.x, this.pos.y) >= 0 {
		if !this.onBlock {
			if this.cnt <= 0 {
				this.onBlock = true
			} else {
				if field.checkInField(this.pos.x, this.pos.y-field.lastScrollY) {
					this.pos.x = px
					this.pos.y = py
				} else {
					this.destroyed()
				}
			}
		}
	} else {
		this.onBlock = false
	}
	this.fire()
	this.addWake()
	this.checkForEnemyHit()
	this.decreaseShield()
}

func (this *Boat) handleGameOver() {
	if isGameOver {
		this.clearBullets()
		if this.cnt < -INVINCIBLE_CNT {
			this.cnt = -RESTART_CNT
		}
	} else if this.cnt < -INVINCIBLE_CNT {
		this.clearBullets()
	}
}

func (this *Boat) fire() {
	switch this.gameMode {
	case GameModeNORMAL:
		this.fireNormal()
	case GameModeTWIN_STICK:
		this.fireTwinStick()
	case GameModeDOUBLE_PLAY:
		this.fireDobulePlay()
	case GameModeMOUSE:
		this.fireMouse()
	}
}

func (this *Boat) addWake() {
	if this.cnt%3 == 0 && this.cnt >= -INVINCIBLE_CNT {
		var sp float32
		if this.vx != 0 || this.vy != 0 {
			sp = 0.4
		} else {
			sp = 0.2
		}
		sp *= 1 + nextSignedFloat(0.33)
		sp *= SPEED_BASE
		this.shape.addWake(this.pos, this.deg, sp, 1)
	}
}

func (this *Boat) decreaseShield() {
	if this.shieldCnt > 0 {
		this.shieldCnt--
	}
}
func (this *Boat) checkForEnemyHit() {
	he := checkAllEnemiesHitShip(this.pos.x, this.pos.y, nil, false)
	if he != nil {
		var rd float32
		if this.pos.distVector(he.pos) < 0.1 {
			rd = 0
		} else {
			rd = atan232(this.pos.x-he.pos.x, this.pos.y-he.pos.y)
		}
		sz := he.size()
		this.refVel.x = Sin32(rd) * sz * 0.1
		this.refVel.y = Cos32(rd) * sz * 0.1
		rs := this.refVel.vctSize()
		if rs > 1 {
			this.refVel.x /= rs
			this.refVel.y /= rs
		}
	}
}

func (this *Boat) moveNormal() {
	padInput = pad.getState()
	if isGameOver || this.cnt < -INVINCIBLE_CNT {
		padInput = pad.getNullState()
	}
	if padInput.dir&PadDirUP != 0 {
		this.vy = 1
	}
	if padInput.dir&PadDirDOWN != 0 {
		this.vy = -1
	}
	if padInput.dir&PadDirRIGHT != 0 {
		this.vx = 1
	}
	if padInput.dir&PadDirLEFT != 0 {
		this.vx = -1
	}
	if this.vx != 0 && this.vy != 0 {
		this.vx *= 0.7
		this.vy *= 0.7
	}
	if this.vx != 0 || this.vy != 0 {
		ad := atan232(this.vx, this.vy)
		ad = normalizeDeg(ad)
		ad -= this.deg
		ad = normalizeDeg(ad)
		this.deg += ad * this.turnRatio * this.turnSpeed
		this.deg = normalizeDeg(this.deg)
	}
}

func (this *Boat) moveTwinStick() {
	stickInput = twinStick.getState()
	if isGameOver || this.cnt < -INVINCIBLE_CNT {
		stickInput.clear()
	}
	this.vx = stickInput.left.x
	this.vy = stickInput.left.y
	if this.vx != 0 || this.vy != 0 {
		ad := atan232(this.vx, this.vy)
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
		stickInput = twinStick.getState()
		if isGameOver || this.cnt < -INVINCIBLE_CNT {
			stickInput.clear()
		}
		this.vx = stickInput.left.x
		this.vy = stickInput.left.y
	case 1:
		this.vx = stickInput.right.x
		this.vy = stickInput.right.y
	}
	if this.vx != 0 || this.vy != 0 {
		ad := atan232(this.vx, this.vy)
		ad = normalizeDeg(ad)
		ad -= this.deg
		ad = normalizeDeg(ad)
		this.deg += ad * this.turnRatio * this.turnSpeed
		this.deg = normalizeDeg(this.deg)
	}
}

func (this *Boat) moveMouse() {
	mps := mouseAndPad.getState()
	padInput = mps.padState
	mouseInput = mps.mouseState
	if isGameOver || this.cnt < -INVINCIBLE_CNT {
		padInput = pad.getNullState()
		mouseInput = mouse.getNullState()
	}
	if padInput.dir&PadDirUP != 0 {
		this.vy = 1
	}
	if padInput.dir&PadDirDOWN != 0 {
		this.vy = -1
	}
	if padInput.dir&PadDirRIGHT != 0 {
		this.vx = 1
	}
	if padInput.dir&PadDirLEFT != 0 {
		this.vx = -1
	}
	if this.vx != 0 && this.vy != 0 {
		this.vx *= 0.7
		this.vy *= 0.7
	}
	if this.vx != 0 || this.vy != 0 {
		ad := atan232(this.vx, this.vy)
		ad = normalizeDeg(ad)
		ad -= this.deg
		ad = normalizeDeg(ad)
		this.deg += ad * this.turnRatio * this.turnSpeed
		this.deg = normalizeDeg(this.deg)
	}
}

func (this *Boat) fireNormal() {
	if padInput.button&PadButtonA != 0 {
		this.turnRatio += (SLOW_TURN_RATIO - this.turnRatio) * TURN_CHANGE_RATIO
		this.fireInterval = FIRE_INTERVAL
		if !this.aPressed {
			this.fireCnt = 0
			this.aPressed = true
		}
	} else {
		this.turnRatio += (TURN_RATIO_BASE - this.turnRatio) * TURN_CHANGE_RATIO
		this.aPressed = false
		this.fireInterval *= 1.033
		if this.fireInterval > FIRE_INTERVAL_MAX {
			this.fireInterval = 99999
		}
	}
	this.fireDeg = this.deg
	if this.reverseFire {
		this.fireDeg += Pi32
	}
	if this.fireCnt <= 0 {
		playSe("shot.wav")
		foc := (this.fireSprCnt%2)*2 - 1
		this.firePos.x = this.pos.x + Cos32(this.fireDeg+Pi32)*0.2*float32(foc)
		this.firePos.y = this.pos.y - Sin32(this.fireDeg+Pi32)*0.2*float32(foc)
		NewShot(this.firePos, this.fireDeg, false, -1)

		this.fireCnt = int(this.fireInterval)
		var td float32
		switch foc {
		case -1:
			td = this.fireSprDeg * float32(this.fireSprCnt/2%4+1) * 0.2
		case 1:
			td = -this.fireSprDeg * float32(this.fireSprCnt/2%4+1) * 0.2
		}
		this.fireSprCnt++
		NewShot(this.firePos, this.fireDeg+td, false, -1)

		sd := this.fireDeg + td/2
		NewSmoke(this.firePos.x, this.firePos.y, 0, Sin32(sd)*SHOT_SPEED*0.33, Cos32(sd)*SHOT_SPEED*0.33, 0, SmokeTypeSPARK, 10, 0.33)
	}
	this.fireCnt--
	if padInput.button&PadButtonB != 0 {
		if !this.bPressed && this.fireLanceCnt <= 0 && !existsLance() {
			playSe("lance.wav")
			fd := this.deg
			if this.reverseFire {
				fd += Pi32
			}
			NewShot(this.pos, fd, true, -1)
			for i := 0; i < 4; i++ {
				sd := fd + nextSignedFloat(1)
				NewSmoke(this.pos.x, this.pos.y, 0,
					Sin32(sd)*LANCE_SPEED*float32(i)*0.2,
					Cos32(sd)*LANCE_SPEED*float32(i)*0.2,
					0, SmokeTypeSPARK, 15, 0.5)
			}
			this.fireLanceCnt = FIRE_LANCE_INTERVAL
		}
		this.bPressed = true
	} else {
		this.bPressed = false
	}
	this.fireLanceCnt--
}

func (this *Boat) fireTwinStick() {
	if fabs32(stickInput.right.x)+fabs32(stickInput.right.y) > 0.01 {
		this.fireDeg = atan232(stickInput.right.x, stickInput.right.y)
		if this.fireCnt <= 0 {
			playSe("shot.wav")
			foc := (this.fireSprCnt%2)*2 - 1
			rsd := stickInput.right.vctSize()
			if rsd > 1 {
				rsd = 1
			}
			this.fireSprDeg = 1 - rsd + 0.05
			this.firePos.x = this.pos.x + Cos32(this.fireDeg+Pi32)*0.2*float32(foc)
			this.firePos.y = this.pos.y - Sin32(this.fireDeg+Pi32)*0.2*float32(foc)
			this.fireCnt = int(this.fireInterval)
			var td float32
			switch foc {
			case -1:
				td = this.fireSprDeg * (Mod32(float32(this.fireSprCnt)/2, 4) + 1) * 0.2
			case 1:
				td = -this.fireSprDeg * (Mod32(float32(this.fireSprCnt)/2, 4) + 1) * 0.2
			}
			this.fireSprCnt++
			NewShot(this.firePos, this.fireDeg+td/2, false, 2)
			NewShot(this.firePos, this.fireDeg+td, false, 2)
			var sd float32 = this.fireDeg + td/2
			NewSmoke(this.firePos.x, this.firePos.y, 0, Sin32(sd)*SHOT_SPEED*0.33, Cos32(sd)*SHOT_SPEED*0.33, 0, SmokeTypeSPARK, 10, 0.33)
		}
	} else {
		this.fireDeg = 99999
	}
	this.fireCnt--
}

func (this *Boat) fireDobulePlay() {
	if isGameOver || this.cnt < -INVINCIBLE_CNT {
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
	if float32(this.fireCnt) > this.fireInterval {
		this.fireCnt = int(this.fireInterval)
	}
	if this.fireCnt <= 0 {
		playSe("shot.wav")
		var foc int = (this.fireSprCnt%2)*2 - 1
		this.fireDeg = 0 //ship.degAmongBoats() + Pi32 / 2
		this.firePos.x = this.pos.x + Cos32(this.fireDeg+Pi32)*0.2*float32(foc)
		this.firePos.y = this.pos.y - Sin32(this.fireDeg+Pi32)*0.2*float32(foc)
		NewShot(this.firePos, this.fireDeg, false, 2)
		this.fireCnt = int(this.fireInterval)
		var sd float32 = this.fireDeg
		NewSmoke(this.firePos.x, this.firePos.y, 0, Sin32(sd)*SHOT_SPEED*0.33, Cos32(sd)*SHOT_SPEED*0.33, 0, SmokeTypeSPARK, 10, 0.33)
		if this.idx == 0 {
			var fd float32 = ship.degAmongBoats() + Pi32/2
			var td float32
			switch foc {
			case -1:
				td = this.fireSprDeg * (Mod32(float32(this.fireSprCnt)/2, 4) + 1) * 0.15
			case 1:
				td = -this.fireSprDeg * (Mod32(float32(this.fireSprCnt)/2, 4) + 1) * 0.15
			}
			this.firePos.x = ship.midstPos().x + Cos32(fd+Pi32)*0.2*float32(foc)
			this.firePos.y = ship.midstPos().y - Sin32(fd+Pi32)*0.2*float32(foc)
			NewShot(this.firePos, fd, false, 2)
			NewShot(this.firePos, fd+td, false, 2)
			NewSmoke(this.firePos.x, this.firePos.y, 0, Sin32(fd+td/2)*SHOT_SPEED*0.33, Cos32(fd+td/2)*SHOT_SPEED*0.33, 0, SmokeTypeSPARK, 10, 0.33)
		}
		this.fireSprCnt++
	}
	this.fireCnt--
}

func (this *Boat) fireMouse() {
	fox := mouseInput.x - this.pos.x
	foy := mouseInput.y - this.pos.y
	if fabs32(fox) < 0.01 {
		fox = 0.01
	}
	if fabs32(foy) < 0.01 {
		foy = 0.01
	}
	this.fireDeg = atan232(fox, foy)
	if mouseInput.button&(MouseButtonLEFT|MouseButtonRIGHT) != 0 {
		if this.fireCnt <= 0 {
			playSe("shot.wav")
			foc := (this.fireSprCnt%2)*2 - 1
			// rsd := stickInput.right.vctSize()
			var fstd float32 = 0.05
			if mouseInput.button&MouseButtonRIGHT != 0 {
				fstd += 0.5
			}
			this.fireSprDeg += (fstd - this.fireSprDeg) * 0.16
			this.firePos.x = this.pos.x + Cos32(this.fireDeg+Pi32)*0.2*float32(foc)
			this.firePos.y = this.pos.y - Sin32(this.fireDeg+Pi32)*0.2*float32(foc)
			this.fireCnt = int(this.fireInterval)
			var td float32
			switch foc {
			case -1:
				td = this.fireSprDeg * float32(this.fireSprCnt/2%4+1) * 0.2
			case 1:
				td = -this.fireSprDeg * float32(this.fireSprCnt/2%4+1) * 0.2
			}
			this.fireSprCnt++
			NewShot(this.firePos, this.fireDeg+td/2, false, 2)
			NewShot(this.firePos, this.fireDeg+td, false, 2)
			sd := this.fireDeg + td/2
			NewSmoke(this.firePos.x, this.firePos.y, 0, Sin32(sd)*SHOT_SPEED*0.33, Cos32(sd)*SHOT_SPEED*0.33, 0,
				SmokeTypeSPARK, 10, 0.33)
		}
	}
	this.fireCnt--
}

func (this *Boat) checkBulletHit(p Vector, pp Vector) bool {
	if this.cnt <= 0 {
		return false
	}
	var bmvx, bmvy, inaa float32
	bmvx = pp.x
	bmvy = pp.y
	bmvx -= p.x
	bmvy -= p.y
	inaa = bmvx*bmvx + bmvy*bmvy
	if inaa > 0.00001 {
		var sofsx, sofsy, inab, hd float32
		sofsx = this.pos.x
		sofsy = this.pos.y
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
	ship.destroyed()
	inGameState.shipDestroyed()
}

func (this *Boat) destroyedBoatShield() {
	for i := 0; i < 100; i++ {
		NewSpark(this.pos, nextSignedFloat(1), nextSignedFloat(1),
			0.5+nextFloat(0.5), 0.5+nextFloat(0.5), 0,
			40+nextInt(40))
	}
	playSe("ship_shield_lost.wav")
	screen.setScreenShake(30, 0.02)
	this.shieldCnt = 0
	this.cnt = -INVINCIBLE_CNT / 2
}

func (this *Boat) destroyedBoat() {
	for i := 0; i < 128; i++ {
		NewSpark(this.pos, nextSignedFloat(1), nextSignedFloat(1),
			0.5+nextFloat(0.5), 0.5+nextFloat(0.5), 0,
			40+nextInt(40))
	}
	playSe("ship_destroyed.wav")
	for i := 0; i < 64; i++ {
		NewSmoke(this.pos.x, this.pos.y, 0, nextSignedFloat(0.2), nextSignedFloat(0.2),
			nextFloat(0.1),
			SmokeTypeEXPLOSION, 50+nextInt(30), 1)
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
		setScreenColor(0.5, 0.9, 0.7, 0.4)
		gl.Begin(gl.LINE_STRIP)
		gl.Vertex2f(this.pos.x, this.pos.y)
		setScreenColor(0.5, 0.9, 0.7, 0.8)
		gl.Vertex2f(this.pos.x+Sin32(this.fireDeg)*20, this.pos.y+Cos32(this.fireDeg)*20)
		gl.End()
	}
	if this.cnt < 0 && (-this.cnt%32) < 16 {
		return
	}
	gl.PushMatrix()
	glTranslate(this.pos)
	gl.Rotatef(-this.deg*180/Pi32, 0, 0, 1)
	this.shape.draw()
	this.bridgeShape.draw()
	if this.shieldCnt > 0 {
		var ss float32 = 0.66
		if this.shieldCnt < 120 {
			ss *= float32(this.shieldCnt) / 120
		}
		gl.Scalef(ss, ss, ss)
		gl.Rotatef(float32(this.shieldCnt)*5, 0, 0, 1)
		this.shieldShape.draw()
	}
	gl.PopMatrix()
}

func (this *Boat) drawFront() {
	if this.cnt < -INVINCIBLE_CNT {
		return
	}
	if this.gameMode == GameModeMOUSE {
		setScreenColor(0.7, 0.9, 0.8, 1.0)
		lineWidth(2)
		this.drawSight(mouseInput.x, mouseInput.y, 0.3)
		ss := 0.9 - 0.8*float32((this.cnt+1024)%32)/32
		setScreenColor(0.5, 0.9, 0.7, 0.8)
		this.drawSight(mouseInput.x, mouseInput.y, ss)
		lineWidth(1)
	}
}

func (this *Boat) drawSight(x float32, y float32, size float32) {
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(x-size, y-size*0.5)
	gl.Vertex2f(x-size, y-size)
	gl.Vertex2f(x-size*0.5, y-size)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(x+size, y-size*0.5)
	gl.Vertex2f(x+size, y-size)
	gl.Vertex2f(x+size*0.5, y-size)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(x+size, y+size*0.5)
	gl.Vertex2f(x+size, y+size)
	gl.Vertex2f(x+size*0.5, y+size)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(x-size, y+size*0.5)
	gl.Vertex2f(x-size, y+size)
	gl.Vertex2f(x-size*0.5, y+size)
	gl.End()
}

func (this *Boat) drawShape() {
	this.shape.draw()
	this.bridgeShape.draw()
}

func (this *Boat) clearBullets() {
	clearBullets()
}
