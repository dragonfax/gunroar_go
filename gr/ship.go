/*
 * $Id: ship.d,v 1.4 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Player's ship.
 */

const SCROLL_SPEED_BASE = 0.01
const SCROLL_SPEED_MAX = 0.1
const SCROLL_START_Y = 2.5

type Ship struct {
	field                                           Field
	boat                                            [2]*Boat
	gameMode                                        int
	boatNum                                         int
	gameState                                       InGameState
	scrollSpeed, scrollSpeedBase                    float32
	midstPos, higherPos, lowerPos, nearPos, nearVel Vector
	bridgeShape                                     BaseShape
}

func NewShip(Pad pad /*TwinStick twinStick, */, Mouse mouse, MouseAndPad mouseAndPad, Field field, Screen screen) *Ship {
	this := new(Ship)
	this.field = field
	Boat.init()
	for i, _ := range this.boat {
		boat[i] = NewBoat(i, this, pad /*twinStick, */, mouse, mouseAndPad, field, screen, sparks, smokes, fragments, wakes)
	}
	this.boatNum = 1
	this.scrollSpeed = SCROLL_SPEED_BASE
	this.scrollSpeedBase = SCROLL_SPEED_BASE
	this.bridgeShape = NewBaseShape(0.3, 0.2, 0.1, ShapeTypeBRIDGE, 0.3, 0.7, 0.7)
	actors[this] = true
}

func (this *Ship) close() {
	for _, b := range this.boat {
		b.close()
	}
	delete(actors, this)
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

func (this *Ship) start(gameMode GameMode) {
	this.gameMode = gameMode
	if gameMode == InGameState.GameMode.DOUBLE_PLAY {
		this.boatNum = 2
	} else {
		this.boatNum = 1
	}
	this.scrollSpeedBase = SCROLL_SPEED_BASE
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].start(this.gameMode)
	}
	this.midstPos.y = 0
	this.midstPos.x = 0
	this.higherPos.x = 0
	this.higherPos.y = 0
	this.lowerPos.x = 0
	this.lowerPos.y = 0
	this.nearPos.y = 0
	this.nearPos.x = 0
	this.nearVel.y = 0
	this.nearVel.x = 0
	this.restart()
}

func (this *Ship) restart() {
	this.scrollSpeed = this.scrollSpeedBase
	for i := 0; i < this.boatNum; i++ {
		this.boat[i].restart()
	}
}

func (this *Ship) move() {
	this.field.scroll(scrollSpeed)
	sf := false
	for i := 0; i < boatNum; i++ {
		this.boat[i].move()
		if this.boat[i].hasCollision &&
			this.boat[i].pos.x > this.field.size.x/3 && this.boat[i].pos.y < -this.field.size.y/4*3 {
			sf = true
		}
	}
	if sf {
		this.gameState.shrinkScoreReel()
	}
	if this.higherPos.y >= SCROLL_START_Y {
		this.scrollSpeed += (SCROLL_SPEED_MAX - this.scrollSpeed) * 0.1
	} else {
		this.scrollSpeed += (this.scrollSpeedBase - this.scrollSpeed) * 0.1
	}
	this.scrollSpeedBase += (SCROLL_SPEED_MAX - this.scrollSpeedBase) * 0.00001
}

func (this *Ship) checkBulletHit(Vector p, Vector pp) bool {
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
	if this.gameMode == InGameState.GameMode.DOUBLE_PLAY && this.boat[0].hasCollision() {
		setScreenColor(0.5, 0.5, 0.9, 0.8)
		gl.Begin(gl.LINE_STRIP)
		gl.Vertex2(this.boat[0].pos.x, this.boat[0].pos.y)
		setScreenColor(0.5, 0.5, 0.9, 0.3)
		gl.Vertex2(this.midstPos.x, this.midstPos.y)
		setScreenColor(0.5, 0.5, 0.9, 0.8)
		gl.Vertex2(this.boat[1].pos.x, this.boat[1].pos.y)
		gl.End()
		gl.PushMatrix()
		glTranslate(this.midstPos)
		gl.Rotatef(-degAmongBoats*180/Pi32, 0, 0, 1)
		this.bridgeShape.draw()
		gl.PopMatrix()
	}
}

func (this *Ship) drawFront() {
	for i := 0; i < boatNum; i++ {
		this.boat[i].drawFront()
	}
}

func (this *Ship) drawShape() {
	this.boat[0].drawShape()
}

func (this *Ship) scrollSpeedBase() float32 {
	return this.scrollSpeedBase
}

func (this *Ship) midstPos() Vector {
	this.midstPos.x = 0
	this.midstPos.y = 0
	for i := 0; i < this.boatNum; i++ {
		this.midstPos.x += this.boat[i].pos.x
		this.midstPos.y += this.boat[i].pos.y
	}
	this.midstPos /= this.boatNum
	return this.midstPos
}

func (this *Ship) higherPos() Vector {
	this.higherPos.y = -99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos.y > this.higherPos.y {
			this.higherPos.x = this.boat[i].pos.x
			this.higherPos.y = this.boat[i].pos.y
		}
	}
	return this.higherPos
}

func (this *Ship) lowerPos() Vector {
	this.lowerPos.y = 99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos.y < this.lowerPos.y {
			this.lowerPos.x = this.boat[i].pos.x
			this.lowerPos.y = this.boat[i].pos.y
		}
	}
	return this.lowerPos
}

func (this *Ship) nearPos(Vector p) Vector {
	var dist float32 = 99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos.dist(p) < dist {
			dist = this.boat[i].pos.dist(p)
			this.nearPos.x = this.boat[i].pos.x
			this.nearPos.y = this.boat[i].pos.y
		}
	}
	return this.nearPos
}

func (this *Ship) nearVel(Vector p) Vector {
	var dist float = 99999
	for i := 0; i < this.boatNum; i++ {
		if this.boat[i].pos.dist(p) < dist {
			dist = this.boat[i].pos.dist(p)
			this.nearVel.x = this.boat[i].vel.x
			this.nearVel.y = this.boat[i].vel.y
		}
	}
	return this.nearVel
}

func (this *Ship) distAmongBoats() float32 {
	return this.boat[0].pos.dist(this.boat[1].pos)
}

func (this *Ship) degAmongBoats() float32 {
	if this.distAmongBoats < 0.1 {
		return 0
	} else {
		return atan2(this.boat[0].pos.x-this.boat[1].pos.x, this.boat[0].pos.y-this.boat[1].pos.y)
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

var padState PadInput

// static TwinStickState stickInput
var mouseState MouseInput

type Boat struct {
	pad Pad
	// TwinStick twinStick
	mouse                     Mouse
	mouseAndPad               MouseAndPad
	field                     Field
	screen                    Screen
	stageManager              StageManager
	gameState                 InGameState
	pos                       Vector
	firePos                   Vector
	deg, speed, turnRatio     float32
	shape                     BaseShape
	bridgeShape               BaseShape
	fireCnt, fireSprCnt       int
	fireIntervalt, fireSprDeg float32
	fireLanceCnt              int
	fireDeg                   float32
	aPressed, bPressed        bool
	cnt                       int
	onBlock                   bool
	vel                       Vector
	refVel                    Vector
	shieldCnt                 int
	shieldShape               ShieldShape
	turnSpeed                 float32
	reverseFire               bool
	gameMode                  GameMode
	vx, vy                    float32
	idx                       int
	ship                      Ship
}

func NewBoat(idx int, ship Ship, pad Pad /*TwinStick twinStick, */, mouse Mouse, mouseAndPad MouseAndPad, field Field, screen Screen) {
	this := new(Boat)
	this.idx = idx
	this.ship = ship
	this.pad = pad
	//this.twinStick = cast(TwinStick) twinStick
	this.mouse = mouse
	this.mouseAndPad = mouseAndPad
	this.field = field
	this.screen = screen
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.wakes = wakes
	switch idx {
	case 0:
		this.shape = NewBaseShape(0.7, 0.6, 0.6, BaseShape.ShapeType.SHIP_ROUNDTAIL, 0.5, 0.7, 0.5)
		this.bridgeShape = NewBaseShape(0.3, 0.6, 0.6, BaseShape.ShapeType.BRIDGE, 0.3, 0.7, 0.3)
		break
	case 1:
		this.shape = NewBaseShape(0.7, 0.6, 0.6, BaseShape.ShapeType.SHIP_ROUNDTAIL, 0.4, 0.3, 0.8)
		this.bridgeShape = NewBaseShape(0.3, 0.6, 0.6, BaseShape.ShapeType.BRIDGE, 0.2, 0.3, 0.6)
		break
	}
	this.turnSpeed = 1
	this.fireInterval = FIRE_INTERVAL
	this.shieldShape = NewShieldShape()
}

func (this *Boat) close() {
	this.shape.close()
	this.bridgeShape.close()
	this.shieldShape.close()
}

func (this *Boat) start(int gameMode) {
	this.gameMode = gameMode
	if gameMode == InGameState.GameMode.DOUBLE_PLAY {
		switch idx {
		case 0:
			this.pos.x = -field.size.x * 0.5
			break
		case 1:
			this.pos.x = field.size.x * 0.5
			break
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
	this.padInput = pad.getNullState()
	//stickInput = twinStick.getNullState()
	this.mouseInput = mouse.getNullState()
}

func (this *Boat) restart() {
	switch this.gameMode {
	case InGameState.GameMode.NORMAL:
		this.fireCnt = 99999
		this.fireInterval = 99999
		break
		/*
			case InGameState.GameMode.TWIN_STICK:
			case InGameState.GameMode.DOUBLE_PLAY:
		*/
	case InGameState.GameMode.MOUSE:
		this.fireCnt = 0
		this.fireInterval = FIRE_INTERVAL
		break
	}
	this.fireSprCnt = 0
	this.fireSprDeg = 0.5
	this.fireLanceCnt = 0
	if this.field.getBlock(this.pos) >= 0 {
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
	case InGameState.GameMode.NORMAL:
		this.moveNormal()
		break
		/*
			case InGameState.GameMode.TWIN_STICK:
				moveTwinStick()
				break
			case InGameState.GameMode.DOUBLE_PLAY:
				moveDoublePlay()
				break
		*/
	case InGameState.GameMode.MOUSE:
		this.moveMouse()
		break
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
	this.vx += this.refVel.x
	this.vy += this.refVel.y
	this.refVel *= 0.9
	if this.field.checkInField(this.pos.x, this.pos.y-this.field.lastScrollY) {
		this.pos.y -= this.field.lastScrollY
	}
	if (this.onBlock || this.field.getBlock(this.pos.x+this.vx, this.pos.y) < 0) &&
		this.field.checkInField(this.pos.x+this.vx, this.pos.y) {
		this.pos.x += this.vx
		this.vel.x = this.vx
	} else {
		this.vel.x = 0
		this.refVel.x = 0
	}
	srf := false
	if (this.onBlock || this.field.getBlock(this.px, this.pos.y+this.vy) < 0) &&
		this.field.checkInField(this.pos.x, this.pos.y+this.vy) {
		this.pos.y += this.vy
		this.vel.y = this.vy
	} else {
		this.vel.y = 0
		this.refVel.y = 0
	}
	if this.field.getBlock(this.pos.x, this.pos.y) >= 0 {
		if !this.onBlock {
			if this.cnt <= 0 {
				this.onBlock = true
			} else {
				if this.field.checkInField(this.pos.x, this.pos.y-this.field.lastScrollY) {
					this.pos.x = this.px
					this.pos.y = this.py
				} else {
					this.destroyed()
				}
			}
		}
	} else {
		this.onBlock = false
	}
	switch this.gameMode {
	case InGameState.GameMode.NORMAL:
		this.fireNormal()
		break
		/*
			case InGameState.GameMode.TWIN_STICK:
				fireTwinStick()
				break
			case InGameState.GameMode.DOUBLE_PLAY:
				fireDobulePlay()
				break
		*/
	case InGameState.GameMode.MOUSE:
		this.fireMouse()
		break
	}
	if this.cnt%3 == 0 && this.cnt >= -INVINCIBLE_CNT {
		var sp float32
		if this.vx != 0 || this.vy != 0 {
			sp = 0.4
		} else {
			sp = 0.2
		}
		sp *= 1 + rand.nextSignedFloat(0.33)
		sp *= SPEED_BASE
		this.shape.addWake(this.pos, this.deg, sp)
	}
	he := checkAllEnemiesHitShip(this.pos.x, this.pos.y)
	var rd float32
	if this.pos.dist(he.pos) < 0.1 {
		rd = 0
	} else {
		rd = atan2(this.pos.x-he.pos.x, this.pos.y-he.pos.y)
	}
	sz := he.size
	this.refVel.x = Sin32(rd) * sz * 0.1
	this.refVel.y = Cos32(rd) * sz * 0.1
	rs := this.refVel.vctSize
	if rs > 1 {
		this.refVel.x /= rs
		this.refVel.y /= rs
	}
	if this.shieldCnt > 0 {
		this.shieldCnt--
	}
}

func (this *Boat) moveNormal() {
	this.padInput = this.pad.getState()
	if this.gameState.isGameOver || this.cnt < -INVINCIBLE_CNT {
		this.padInput.clear()
	}
	if this.padInput.dir & PadStateDirUP {
		this.vy = 1
	}
	if this.padInput.dir & PadStateDirDOWN {
		this.vy = -1
	}
	if this.padInput.dir & PadStateDirRIGHT {
		this.vx = 1
	}
	if this.padInput.dir & PadStateDirLEFT {
		this.vx = -1
	}
	if this.vx != 0 && this.vy != 0 {
		this.vx *= 0.7
		this.vy *= 0.7
	}
	if this.vx != 0 || this.vy != 0 {
		ad := atan2(this.vx, this.vy)
		normalizeDeg(ad)
		ad -= this.deg
		normalizeDeg(ad)
		this.deg += ad * this.turnRatio * this.turnSpeed
		normalizeDeg(this.deg)
	}
}

/*
moveTwinStick() {
		stickInput = twinStick.getState()
	if (gameState.isGameOver || cnt < -INVINCIBLE_CNT)
		stickInput.clear()
	vx = stickInput.left.x
	vy = stickInput.left.y
	if (vx != 0 || vy != 0) {
		float32 ad = atan2(vx, vy)
		assert(ad <>= 0)
		Math.normalizeDeg(ad)
		ad -= deg
		Math.normalizeDeg(ad)
		deg += ad * turnRatio * turnSpeed
		Math.normalizeDeg(deg)
	}
}

moveDoublePlay() {
	switch (idx) {
	case 0:
			stickInput = twinStick.getState()
		if (gameState.isGameOver || cnt < -INVINCIBLE_CNT)
			stickInput.clear()
		vx = stickInput.left.x
		vy = stickInput.left.y
		break
	case 1:
		vx = stickInput.right.x
		vy = stickInput.right.y
		break
	}
	if (vx != 0 || vy != 0) {
		float32 ad = atan2(vx, vy)
		assert(ad <>= 0)
		Math.normalizeDeg(ad)
		ad -= deg
		Math.normalizeDeg(ad)
		deg += ad * turnRatio * turnSpeed
		Math.normalizeDeg(deg)
	}
}
*/

func (this *Boat) moveMouse() {
	mps := mouseAndPad.getState()
	this.padInput = mps.padState
	this.mouseInput = mps.mouseState
	if this.gameState.isGameOver || this.cnt < -INVINCIBLE_CNT {
		this.padInput.clear()
		this.mouseInput.clear()
	}
	if this.padInput.dir & PadStateDirUP {
		this.vy = 1
	}
	if this.padInput.dir & PadStateDirDOWN {
		this.vy = -1
	}
	if this.padInput.dir & PadStateDirRIGHT {
		this.vx = 1
	}
	if this.padInput.dir & PadStateDirLEFT {
		this.vx = -1
	}
	if this.vx != 0 && this.vy != 0 {
		this.vx *= 0.7
		this.vy *= 0.7
	}
	if this.vx != 0 || this.vy != 0 {
		ad := atan2(this.vx, this.vy)
		normalizeDeg(ad)
		ad -= this.deg
		normalizeDeg(ad)
		this.deg += ad * this.turnRatio * this.turnSpeed
		normalizeDeg(this.deg)
	}
}

func (this *Boat) fireNormal() {
	if this.padInput.button & PadStateButtonA {
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
		this.firePos.x = this.pos.x + Cos32(this.fireDeg+Pi32)*0.2*this.foc
		this.firePos.y = this.pos.y - Sin32(this.fireDeg+Pi32)*0.2*this.foc
		NewShot(this.firePos, this.fireDeg)

		this.fireCnt = int(this.fireInterval)
		var td float32
		switch this.foc {
		case -1:
			td = this.fireSprDeg * (this.fireSprCnt/2%4 + 1) * 0.2
			break
		case 1:
			td = -this.fireSprDeg * (this.fireSprCnt/2%4 + 1) * 0.2
			break
		}
		this.fireSprCnt++
		NewShot(this.firePos, this.fireDeg+td)

		sd := this.fireDeg + td/2
		NewSmoke(this.firePos, Sin32(sd)*SPEED*0.33, Cos32(sd)*SPEED*0.33, 0, SmokeTypeSPARK, 10, 0.33)
	}
	this.fireCnt--
	if this.padInput.button & PadStateButtonB {
		if !this.bPressed && this.fireLanceCnt <= 0 && !existsLance() {
			playSe("lance.wav")
			fd := this.deg
			if this.reverseFire {
				fd += Pi32
			}
			NewShot(this.pos, fd, true)
			for i := 0; i < 4; i++ {
				sd := fd + rand.nextSignedFloat(1)
				NewSmoke(this.pos,
					Sin32(sd)*LANCE_SPEED*i*0.2,
					Cos32(sd)*LANCE_SPEED*i*0.2,
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

/*
fireTwinStick() {
	if (fabs(stickInput.right.x) + fabs(stickInput.right.y) > 0.01) {
		fireDeg = atan2(stickInput.right.x, stickInput.right.y)
		assert(fireDeg <>= 0)
		if (fireCnt <= 0) {
			SoundManager.playSe("shot.wav")
			int foc = (fireSprCnt % 2) * 2 - 1
			float32 rsd = stickInput.right.vctSize
			if (rsd > 1)
				rsd = 1
			fireSprDeg = 1 - rsd + 0.05
			firePos.x = this.pos.x + Cos32(fireDeg + Pi32) * 0.2 * foc
			firePos.y = this.pos.y - Sin32(fireDeg + Pi32) * 0.2 * foc
			fireCnt = cast(int) fireInterval
			float32 td
			switch (foc) {
			case -1:
				td = fireSprDeg * (fireSprCnt / 2 % 4 + 1) * 0.2
				break
			case 1:
				td = -fireSprDeg * (fireSprCnt / 2 % 4 + 1) * 0.2
				break
			}
			fireSprCnt++
			Shot s = shots.getInstance()
			if (s)
				s.set(firePos, fireDeg + td / 2, false, 2)
			s = shots.getInstance()
			if (s)
				s.set(firePos, fireDeg + td, false, 2)
			Smoke sm = smokes.getInstanceForced()
			float32 sd = fireDeg + td / 2
			sm.set(firePos, Sin32(sd) * Shot.SPEED * 0.33, Cos32(sd) * Shot.SPEED * 0.33, 0,
						 Smoke.SmokeType.SPARK, 10, 0.33)
		}
	} else {
		fireDeg = 99999
	}
	fireCnt--
}

fireDobulePlay() {
	if (gameState.isGameOver || cnt < -INVINCIBLE_CNT)
		return
	float32 dist = ship.distAmongBoats()
	fireInterval = FIRE_INTERVAL + 10.0 / (dist + 0.005)
	if (dist < 2)
		fireInterval = 99999
	else if (dist < 4)
		fireInterval *= 3
	else if (dist < 6)
		fireInterval *= 1.6
	if (fireCnt > fireInterval)
		fireCnt = cast(int) fireInterval
	if (fireCnt <= 0) {
		SoundManager.playSe("shot.wav")
		int foc = (fireSprCnt % 2) * 2 - 1
		fireDeg = 0;//ship.degAmongBoats() + Pi32 / 2
		firePos.x = this.pos.x + Cos32(fireDeg + Pi32) * 0.2 * foc
		firePos.y = this.pos.y - Sin32(fireDeg + Pi32) * 0.2 * foc
		Shot s = shots.getInstance()
		if (s)
			s.set(firePos, fireDeg, false , 2)
		fireCnt = cast(int) fireInterval
		Smoke sm = smokes.getInstanceForced()
		float32 sd = fireDeg
		sm.set(firePos, Sin32(sd) * Shot.SPEED * 0.33, Cos32(sd) * Shot.SPEED * 0.33, 0,
					 Smoke.SmokeType.SPARK, 10, 0.33)
		if (idx == 0) {
			float32 fd = ship.degAmongBoats() + Pi32 / 2
			float32 td
			switch (foc) {
			case -1:
				td = fireSprDeg * (fireSprCnt / 2 % 4 + 1) * 0.15
				break
			case 1:
				td = -fireSprDeg * (fireSprCnt / 2 % 4 + 1) * 0.15
				break
			}
			firePos.x = ship.midstPos.x + Cos32(fd + Pi32) * 0.2 * foc
			firePos.y = ship.midstPos.y - Sin32(fd + Pi32) * 0.2 * foc
			s = shots.getInstance()
			if (s)
				s.set(firePos, fd, false, 2)
			s = shots.getInstance()
			if (s)
				s.set(firePos, fd + td, false , 2)
			sm = smokes.getInstanceForced()
			sm.set(firePos, Sin32(fd + td / 2) * Shot.SPEED * 0.33, Cos32(fd + td / 2) * Shot.SPEED * 0.33, 0,
						 Smoke.SmokeType.SPARK, 10, 0.33)
		}
		fireSprCnt++
	}
	fireCnt--
}
*/

func (this *Boat) fireMouse() {
	fox := this.mouseInput.x - this.pos.x
	foy := this.mouseInput.y - this.pos.y
	if fabs32(fox) < 0.01 {
		fox = 0.01
	}
	if fabs32(foy) < 0.01 {
		foy = 0.01
	}
	this.fireDeg = atan2(fox, foy)
	if this.mouseInput.button & (MouseStateButtonLEFT | MouseStateButtonRIGHT) {
		if this.fireCnt <= 0 {
			playSe("shot.wav")
			foc := (this.fireSprCnt%2)*2 - 1
			//rsd := this.stickInput.right.vctSize
			fstd := 0.05
			if this.mouseInput.button & MouseStateButtonRIGHT {
				fstd += 0.5
			}
			this.fireSprDeg += (fstd - this.fireSprDeg) * 0.16
			this.firePos.x = this.pos.x + Cos32(this.fireDeg+Pi32)*0.2*foc
			this.firePos.y = this.pos.y - Sin32(this.fireDeg+Pi32)*0.2*foc
			this.fireCnt = int(this.fireInterval)
			var td float32
			switch foc {
			case -1:
				td = this.fireSprDeg * (this.fireSprCnt/2%4 + 1) * 0.2
				break
			case 1:
				td = -this.fireSprDeg * (this.fireSprCnt/2%4 + 1) * 0.2
				break
			}
			this.fireSprCnt++
			NewShot(this.firePos, this.fireDeg+td/2, false, 2)
			NewShot(this.firePos, this.fireDeg+td, false, 2)
			sd := this.fireDeg + td/2
			NewSmoke(this.firePos, Sin32(sd)*SPEED*0.33, Cos32(sd)*SPEED*0.33, 0,
				SmokeTypeSPARK, 10, 0.33)
		}
	}
	this.fireCnt--
}

func (this *Boat) checkBulletHit(Vector p, Vector pp) bool {
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
	this.ship.destroyed()
	this.gameState.shipDestroyed()
}

func (this *Boat) destroyedBoatShield() {
	for i := 0; i < 100; i++ {
		NewSpark(pos, rand.nextSignedFloat(1), rand.nextSignedFloat(1),
			0.5+rand.nextFloat(0.5), 0.5+rand.nextFloat(0.5), 0,
			40+rand.nextInt(40))
	}
	playSe("ship_shield_lost.wav")
	setScreenShake(30, 0.02)
	this.shieldCnt = 0
	this.cnt = -INVINCIBLE_CNT / 2
}

func (this *Boat) destroyedBoat() {
	for i := 0; i < 128; i++ {
		NewSpark(pos, rand.nextSignedFloat(1), rand.nextSignedFloat(1),
			0.5+rand.nextFloat(0.5), 0.5+rand.nextFloat(0.5), 0,
			40+rand.nextInt(40))
	}
	playSe("ship_destroyed.wav")
	for i := 0; i < 64; i++ {
		NewSmoke(pos, rand.nextSignedFloat(0.2), rand.nextSignedFloat(0.2),
			rand.nextFloat(0.1),
			SmokeTypeEXPLOSION, 50+rand.nextInt(30), 1)
	}
	setScreenShake(60, 0.05)
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
		gl.Vertex2(this.pos.x, this.pos.y)
		setScreenColor(0.5, 0.9, 0.7, 0.8)
		gl.Vertex2(this.pos.x+Sin32(this.fireDeg)*20, this.pos.y+Cos32(this.fireDeg)*20)
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
		ss := 0.66
		if this.shieldCnt < 120 {
			ss *= float32(this.shieldCnt) / 120
		}
		gl.Scalef(ss, ss, ss)
		gl.Rotatef(this.shieldCnt*5, 0, 0, 1)
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
		this.drawSight(this.mouseInput.x, this.mouseInput.y, 0.3)
		ss := 0.9 - 0.8*((this.cnt+1024)%32)/32
		setScreenColor(0.5, 0.9, 0.7, 0.8)
		this.drawSight(this.mouseInput.x, this.mouseInput.y, ss)
		lineWidth(1)
	}
}

func (this *Boat) drawSight(x float32, x float32, size float32) {
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2(x-size, y-size*0.5)
	gl.Vertex2(x-size, y-size)
	gl.Vertex2(x-size*0.5, y-size)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2(x+size, y-size*0.5)
	gl.Vertex2(x+size, y-size)
	gl.Vertex2(x+size*0.5, y-size)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2(x+size, y+size*0.5)
	gl.Vertex2(x+size, y+size)
	gl.Vertex2(x+size*0.5, y+size)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2(x-size, y+size*0.5)
	gl.Vertex2(x-size, y+size)
	gl.Vertex2(x-size*0.5, y+size)
	gl.End()
}

func (this *Boat) drawShape() {
	this.shape.draw()
	this.bridgeShape.draw()
}

func (this *Boat) clearBullets() {
	this.gameState.clearBullets()
}
