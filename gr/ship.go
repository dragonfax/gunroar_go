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
  field Field
  boat [2]*Boat
  gameMode int
  boatNum int
  gameState InGameState
  scrollSpeed, scrollSpeedBase float32
  midstPos, higherPos, lowerPos, nearPos, nearVel Vector
  bridgeShape BaseShape
}

func NewShip(Pad pad, /*TwinStick twinStick, */Mouse mouse, MouseAndPad mouseAndPad, Field field, Screen screen ) *Ship {
	this := new(Ship)
	this.field = field
	Boat.init()
	for i, _ := range this.boat {
		boat[i] = NewBoat(i, this, pad, /*twinStick, */mouse, mouseAndPad, field, screen, sparks, smokes, fragments, wakes)
	}
	this.boatNum = 1
	this.scrollSpeed = SCROLL_SPEED_BASE
	this.scrollSpeedBase = SCROLL_SPEED_BASE
	this.bridgeShape = NewBaseShape(0.3, 0.2, 0.1, ShapeTypeBRIDGE, 0.3, 0.7, 0.7)
}

func (this *Ship) close() {
	for _, b := range this.boat {
		b.close()
	}
}

func (this *Ship) setStageManager(stageManager StageManager) {
	for _,b:= range this.boat {
		b.setStageManager(stageManager)
	}
}

func (this *Ship) setGameState(gameState InGameState) {
	this.gameState = gameState
	for _,b:= range this.boat {
		b.setGameState(gameState)
	}
}

func (this *Ship) start(gameMode GameMode) {
	this.gameMode = gameMode
	if (gameMode == InGameState.GameMode.DOUBLE_PLAY) {
		boatNum = 2
	} else {
		boatNum = 1
	}	
	_scrollSpeedBase = SCROLL_SPEED_BASE
	for i := 0; i < boatNum; i++ {
		boat[i].start(gameMode)
	}
	_midstPos.x = _midstPos.y = 0
	_higherPos.x = _higherPos.y = 0
	_lowerPos.x = _lowerPos.y = 0
	_nearPos.x = _nearPos.y = 0
	_nearVel.x = _nearVel.y = 0
	restart()
}

func (this *Ship) restart() {
	scrollSpeed = _scrollSpeedBase
	for i := 0; i < boatNum; i++ {
		boat[i].restart()
	}
}

func (this *Ship) move() {
	field.scroll(scrollSpeed)
	float sf = false
	for i := 0; i < boatNum; i++ {
		boat[i].move()
		if (boat[i].hasCollision &&
				boat[i].pos.x > field.size.x / 3 && boat[i].pos.y < -field.size.y / 4 * 3) {
			sf = true
		}
	}
	if (sf){
		gameState.shrinkScoreReel()
	}
	if (higherPos.y >= SCROLL_START_Y){
		scrollSpeed += (SCROLL_SPEED_MAX - scrollSpeed) * 0.1
	} else {
		scrollSpeed += (_scrollSpeedBase - scrollSpeed) * 0.1
	}
	_scrollSpeedBase += (SCROLL_SPEED_MAX - _scrollSpeedBase) * 0.00001
}

func (this *Ship) bool checkBulletHit(Vector p, Vector pp) {
	for i := 0; i < boatNum; i++ {
		if (boat[i].checkBulletHit(p, pp)) {
			return true
		}
	}
	return false
}

func (this *Ship) clearBullets() {
	gameState.clearBullets()
}

func (this *Ship) destroyed() {
	for i := 0; i < boatNum; i++ {
		boat[i].destroyedBoat()
	}
}

func (this *Ship) draw() {
	for i := 0; i < boatNum; i++ {
		boat[i].draw()
	}
	if (gameMode == InGameState.GameMode.DOUBLE_PLAY && boat[0].hasCollision) {
		Screen.setColor(0.5, 0.5, 0.9, 0.8)
		glBegin(GL_LINE_STRIP)
		glVertex2(boat[0].pos.x, boat[0].pos.y)
		Screen.setColor(0.5, 0.5, 0.9, 0.3)
		glVertex2(midstPos.x, midstPos.y)
		Screen.setColor(0.5, 0.5, 0.9, 0.8)
		glVertex2(boat[1].pos.x, boat[1].pos.y)
		glEnd()
		glPushMatrix()
		Screen.glTranslate(midstPos)
		glRotatef(-degAmongBoats * 180 / Pi32, 0, 0, 1)
		bridgeShape.draw()
		glPopMatrix()
	}
}

func (this *Ship) drawFront() {
	for i := 0; i < boatNum; i++ {
		boat[i].drawFront()
	}
}

func (this *Ship) drawShape() {
	boat[0].drawShape()
}

func (this *Ship) float scrollSpeedBase() {
	return _scrollSpeedBase
}

func (this *Ship) Vector midstPos() {
	_midstPos.x = _midstPos.y = 0
	for i := 0; i < boatNum; i++ {
		_midstPos.x += boat[i].pos.x
		_midstPos.y += boat[i].pos.y
	}
	_midstPos /= boatNum
	return _midstPos
}

func (this *Ship) Vector higherPos() {
	_higherPos.y = -99999
	for i := 0; i < boatNum; i++ {
		if (boat[i].pos.y > _higherPos.y) {
			_higherPos.x = boat[i].pos.x
			_higherPos.y = boat[i].pos.y
		}
	}
	return _higherPos
}

func (this *Ship) Vector lowerPos() {
	_lowerPos.y = 99999
	for i := 0; i < boatNum; i++ {
		if (boat[i].pos.y < _lowerPos.y) {
			_lowerPos.x = boat[i].pos.x
			_lowerPos.y = boat[i].pos.y
		}
	}
	return _lowerPos
}

func (this *Ship) Vector nearPos(Vector p) {
	float dist = 99999
	for i := 0; i < boatNum; i++ {
		if (boat[i].pos.dist(p) < dist) {
			dist = boat[i].pos.dist(p)
			_nearPos.x = boat[i].pos.x
			_nearPos.y = boat[i].pos.y
		}
	}
	return _nearPos
}

func (this *Ship) Vector nearVel(Vector p) {
	float dist = 99999
	for i := 0; i < boatNum; i++ {
		if (boat[i].pos.dist(p) < dist) {
			dist = boat[i].pos.dist(p)
			_nearVel.x = boat[i].vel.x
			_nearVel.y = boat[i].vel.y
		}
	}
	return _nearVel
}

func (this *Ship) float distAmongBoats() {
	return boat[0].pos.dist(boat[1].pos)
}

func (this *Ship) float degAmongBoats() {
	if (distAmongBoats < 0.1) {
		return 0
	} else {
		return atan2(boat[0].pos.x - boat[1].pos.x, boat[0].pos.y - boat[1].pos.y)
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
Rand rand
PadState padInput
// static TwinStickState stickInput
MouseState mouseInput

type Boat struct {
  pad Pad
  // TwinStick twinStick
  mouse Mouse
  mouseAndPad MouseAndPad
  field Field
  screen Screen
  stageManager StageManager
  gameState InGameState
  pos Vector
  firePos Vector
  deg, speed, turnRatio float32
  shape BaseShape
  bridgeShape BaseShape
  fireCnt, fireSprCnt int
  fireIntervalt, fireSprDeg float32
  fireLanceCnt int
  fireDeg float32
  aPressed, bPressed bool
  cnt int
  onBlock bool
  vel Vector
  refVel Vector
  shieldCnt int
  shieldShape ShieldShape
  turnSpeed float32
  reverseFire bool
  gameMode GameMode
  vx, vy float32
  idx int
  ship Ship
}

  this(int idx, Ship ship,
              Pad pad, /*TwinStick twinStick, */Mouse mouse, MouseAndPad mouseAndPad,
              Field field, Screen screen ) {
    this.idx = idx
    this.ship = ship
    this.pad = cast(Pad) pad
    //this.twinStick = cast(TwinStick) twinStick
    this.mouse = cast(Mouse) mouse
    this.mouseAndPad = mouseAndPad
    this.field = field
    this.screen = screen
    this.sparks = sparks
    this.smokes = smokes
    this.fragments = fragments
    this.wakes = wakes
    _pos = new Vector
    firePos = new Vector
    _vel = new Vector
    refVel = new Vector
    switch (idx) {
    case 0:
      _shape = new BaseShape(0.7, 0.6, 0.6, BaseShape.ShapeType.SHIP_ROUNDTAIL, 0.5, 0.7, 0.5)
      bridgeShape = new BaseShape(0.3, 0.6, 0.6, BaseShape.ShapeType.BRIDGE, 0.3, 0.7, 0.3)
      break
    case 1:
      _shape = new BaseShape(0.7, 0.6, 0.6, BaseShape.ShapeType.SHIP_ROUNDTAIL, 0.4, 0.3, 0.8)
      bridgeShape = new BaseShape(0.3, 0.6, 0.6, BaseShape.ShapeType.BRIDGE, 0.2, 0.3, 0.6)
      break
    }
    deg = 0
    speed = 0
    turnRatio = 0
    turnSpeed = 1
    fireInterval = FIRE_INTERVAL
    fireSprDeg = 0
    cnt = 0
    shieldCnt = 0
    shieldShape = new ShieldShape
  }

  func (this *Boat) close() {
    _shape.close()
    bridgeShape.close()
    shieldShape.close()
  }

  func (this *Boat) start(int gameMode) {
    this.gameMode = gameMode
    if (gameMode == InGameState.GameMode.DOUBLE_PLAY) {
      switch (idx) {
      case 0:
        _pos.x = -field.size.x * 0.5
        break
      case 1:
        _pos.x = field.size.x * 0.5
        break
      }
    } else {
      _pos.x = 0
    }
    _pos.y = -field.size.y * 0.8
    firePos.x = firePos.y = 0
    _vel.x = _vel.y = 0
    deg = 0
    speed = SPEED_BASE
    turnRatio = TURN_RATIO_BASE
    cnt = -INVINCIBLE_CNT
    aPressed = bPressed = true
    padInput = pad.getNullState()
    //stickInput = twinStick.getNullState()
    mouseInput = mouse.getNullState()
  }

  func (this *Boat) restart() {
    switch (gameMode) {
    case InGameState.GameMode.NORMAL:
      fireCnt = 99999
      fireInterval = 99999
      break
      /*
    case InGameState.GameMode.TWIN_STICK:
    case InGameState.GameMode.DOUBLE_PLAY:
    */
    case InGameState.GameMode.MOUSE:
      fireCnt = 0
      fireInterval = FIRE_INTERVAL
      break
    }
    fireSprCnt = 0
    fireSprDeg = 0.5
    fireLanceCnt = 0
    if (field.getBlock(_pos) >= 0) {
      onBlock = true
		} else {
      onBlock = false
		}
    refVel.x = refVel.y = 0
    shieldCnt = 20 * 60
  }

  func (this *Boat) move() {
    float px = _pos.x, py = _pos.y
    cnt++
    vx = vy = 0
    switch (gameMode) {
    case InGameState.GameMode.NORMAL:
      moveNormal()
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
      moveMouse()
      break
    }
    if (gameState.isGameOver) {
      clearBullets()
      if (cnt < -INVINCIBLE_CNT) {
        cnt = -RESTART_CNT
			}
    } else if (cnt < -INVINCIBLE_CNT) {
      clearBullets()
    }
    vx *= speed
    vy *= speed
    vx += refVel.x
    vy += refVel.y
    refVel *= 0.9
    if (field.checkInField(_pos.x, _pos.y - field.lastScrollY)) {
      _pos.y -= field.lastScrollY
		}
    if ((onBlock || field.getBlock(_pos.x + vx, _pos.y) < 0) &&
        field.checkInField(_pos.x + vx, _pos.y)) {
      _pos.x += vx
      _vel.x = vx
    } else {
      _vel.x = 0
      refVel.x = 0
    }
    bool srf = false
    if ((onBlock || field.getBlock(px, _pos.y + vy) < 0) &&
        field.checkInField(_pos.x, _pos.y + vy)) {
      _pos.y += vy
      _vel.y = vy
    } else {
      _vel.y = 0
      refVel.y = 0
    }
    if (field.getBlock(_pos.x, _pos.y) >= 0) {
      if (!onBlock) {
        if (cnt <= 0) {
          onBlock = true
				} else {
          if (field.checkInField(_pos.x, _pos.y - field.lastScrollY)) {
            _pos.x = px
            _pos.y = py
          } else {
            destroyed()
          }
        }
			}
    } else {
      onBlock = false
    }
    switch (gameMode) {
    case InGameState.GameMode.NORMAL:
      fireNormal()
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
      fireMouse()
      break
    }
    if (cnt % 3 == 0 && cnt >= -INVINCIBLE_CNT) {
      float sp
      if (vx != 0 || vy != 0) {
        sp = 0.4
			} else {
        sp = 0.2
			}
      sp *= 1 + rand.nextSignedFloat(0.33)
      sp *= SPEED_BASE
      _shape.addWake(wakes, _pos, deg, sp)
    }
    Enemy he = enemies.checkHitShip(pos.x, pos.y)
    if (he) {
      float rd
      if (pos.dist(he.pos) < 0.1) {
        rd = 0
			} else {
        rd = atan2(_pos.x - he.pos.x, _pos.y - he.pos.y)
			}
      assert(rd <>= 0)
      float sz = he.size
      refVel.x = Sin32(rd) * sz * 0.1
      refVel.y = Cos32(rd) * sz * 0.1
      float rs = refVel.vctSize
      if (rs > 1) {
        refVel.x /= rs
        refVel.y /= rs
      }
    }
    if (shieldCnt > 0) {
      shieldCnt--
		}
  }

  func (this *Boat) moveNormal() {
    padInput = pad.getState()
    if (gameState.isGameOver || cnt < -INVINCIBLE_CNT) {
      padInput.clear()
		}
    if (padInput.dir & PadState.Dir.UP) {
      vy = 1
		}
    if (padInput.dir & PadState.Dir.DOWN) {
      vy = -1
		}
    if (padInput.dir & PadState.Dir.RIGHT) {
      vx = 1
		}
    if (padInput.dir & PadState.Dir.LEFT) {
      vx = -1
		}
    if (vx != 0 && vy != 0) {
      vx *= 0.7
      vy *= 0.7
    }
    if (vx != 0 || vy != 0) {
      float ad = atan2(vx, vy)
      assert(ad <>= 0)
      Math.normalizeDeg(ad)
      ad -= deg
      Math.normalizeDeg(ad)
      deg += ad * turnRatio * turnSpeed
      Math.normalizeDeg(deg)
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
      float ad = atan2(vx, vy)
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
      float ad = atan2(vx, vy)
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
      MouseAndPadState mps = mouseAndPad.getState()
      padInput = mps.padState
      mouseInput = mps.mouseState
    if (gameState.isGameOver || cnt < -INVINCIBLE_CNT) {
      padInput.clear()
      mouseInput.clear()
    }
    if (padInput.dir & PadState.Dir.UP) {
      vy = 1
		}
    if (padInput.dir & PadState.Dir.DOWN) {
      vy = -1
		}
    if (padInput.dir & PadState.Dir.RIGHT) {
      vx = 1
		}
    if (padInput.dir & PadState.Dir.LEFT) {
      vx = -1
		}
    if (vx != 0 && vy != 0) {
      vx *= 0.7
      vy *= 0.7
    }
    if (vx != 0 || vy != 0) {
      float ad = atan2(vx, vy)
      assert(ad <>= 0)
      Math.normalizeDeg(ad)
      ad -= deg
      Math.normalizeDeg(ad)
      deg += ad * turnRatio * turnSpeed
      Math.normalizeDeg(deg)
    }
  }

  func (this *Boat) fireNormal() {
    if (padInput.button & PadState.Button.A) {
      turnRatio += (SLOW_TURN_RATIO - turnRatio) * TURN_CHANGE_RATIO
      fireInterval = FIRE_INTERVAL
      if (!aPressed) {
        fireCnt = 0
        aPressed = true
      }
    } else {
      turnRatio += (TURN_RATIO_BASE - turnRatio) * TURN_CHANGE_RATIO
      aPressed = false
      fireInterval *= 1.033
      if (fireInterval > FIRE_INTERVAL_MAX){
        fireInterval = 99999
			}
    }
    fireDeg = deg
    if (reverseFire) {
      fireDeg += Pi32
		}
    if (fireCnt <= 0) {
      SoundManager.playSe("shot.wav")
      Shot s = shots.getInstance()
      int foc = (fireSprCnt % 2) * 2 - 1
      firePos.x = _pos.x + Cos32(fireDeg + Pi32) * 0.2 * foc
      firePos.y = _pos.y - Sin32(fireDeg + Pi32) * 0.2 * foc
      if (s) {
        s.set(firePos, fireDeg)
			}
      fireCnt = cast(int) fireInterval
      float td
      switch (foc) {
      case -1:
        td = fireSprDeg * (fireSprCnt / 2 % 4 + 1) * 0.2
        break
      case 1:
        td = -fireSprDeg * (fireSprCnt / 2 % 4 + 1) * 0.2
        break
      }
      fireSprCnt++
      s = shots.getInstance()
      if (s) {
        s.set(firePos, fireDeg + td)
			}
      Smoke sm = smokes.getInstanceForced()
      float sd = fireDeg + td / 2
      sm.set(firePos, Sin32(sd) * Shot.SPEED * 0.33, Cos32(sd) * Shot.SPEED * 0.33, 0,
             Smoke.SmokeType.SPARK, 10, 0.33)
    }
    fireCnt--
    if (padInput.button & PadState.Button.B) {
      if (!bPressed && fireLanceCnt <= 0 && !shots.existsLance()) {
        SoundManager.playSe("lance.wav")
        float fd = deg
        if (reverseFire) {
          fd += Pi32
				}
        Shot s = shots.getInstance()
        if (s) {
          s.set(pos, fd, true)
				}
				for i := 0; i < 4; i++ {
          Smoke sm = smokes.getInstanceForced()
          float sd = fd + rand.nextSignedFloat(1)
          sm.set(pos,
                 Sin32(sd) * Shot.LANCE_SPEED * i * 0.2,
                 Cos32(sd) * Shot.LANCE_SPEED * i * 0.2,
                 0, Smoke.SmokeType.SPARK, 15, 0.5)
        }
        fireLanceCnt = FIRE_LANCE_INTERVAL
      }
      bPressed = true
    } else {
      bPressed = false
    }
    fireLanceCnt--
  }

  /*
  fireTwinStick() {
    if (fabs(stickInput.right.x) + fabs(stickInput.right.y) > 0.01) {
      fireDeg = atan2(stickInput.right.x, stickInput.right.y)
      assert(fireDeg <>= 0)
      if (fireCnt <= 0) {
        SoundManager.playSe("shot.wav")
        int foc = (fireSprCnt % 2) * 2 - 1
        float rsd = stickInput.right.vctSize
        if (rsd > 1)
          rsd = 1
        fireSprDeg = 1 - rsd + 0.05
        firePos.x = _pos.x + Cos32(fireDeg + Pi32) * 0.2 * foc
        firePos.y = _pos.y - Sin32(fireDeg + Pi32) * 0.2 * foc
        fireCnt = cast(int) fireInterval
        float td
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
        float sd = fireDeg + td / 2
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
    float dist = ship.distAmongBoats()
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
      firePos.x = _pos.x + Cos32(fireDeg + Pi32) * 0.2 * foc
      firePos.y = _pos.y - Sin32(fireDeg + Pi32) * 0.2 * foc
      Shot s = shots.getInstance()
      if (s)
        s.set(firePos, fireDeg, false , 2)
      fireCnt = cast(int) fireInterval
      Smoke sm = smokes.getInstanceForced()
      float sd = fireDeg
      sm.set(firePos, Sin32(sd) * Shot.SPEED * 0.33, Cos32(sd) * Shot.SPEED * 0.33, 0,
             Smoke.SmokeType.SPARK, 10, 0.33)
      if (idx == 0) {
        float fd = ship.degAmongBoats() + Pi32 / 2
        float td
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
    float fox = mouseInput.x - _pos.x
    float foy = mouseInput.y - _pos.y
    if (fabs(fox) < 0.01) {
      fox = 0.01
		}
    if (fabs(foy) < 0.01) {
      foy = 0.01
		}
    fireDeg = atan2(fox, foy)
    assert(fireDeg <>= 0)
    if (mouseInput.button & (MouseState.Button.LEFT | MouseState.Button.RIGHT)) {
      if (fireCnt <= 0) {
        SoundManager.playSe("shot.wav")
        int foc = (fireSprCnt % 2) * 2 - 1
        float rsd = stickInput.right.vctSize
        float fstd = 0.05
        if (mouseInput.button & MouseState.Button.RIGHT) {
          fstd += 0.5
				}
        fireSprDeg += (fstd - fireSprDeg) * 0.16
        firePos.x = _pos.x + Cos32(fireDeg + Pi32) * 0.2 * foc
        firePos.y = _pos.y - Sin32(fireDeg + Pi32) * 0.2 * foc
        fireCnt = cast(int) fireInterval
        float td
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
        if (s) {
          s.set(firePos, fireDeg + td / 2, false, 2)
				}
        s = shots.getInstance()
        if (s){
          s.set(firePos, fireDeg + td, false, 2)
				}
        Smoke sm = smokes.getInstanceForced()
        float sd = fireDeg + td / 2
        sm.set(firePos, Sin32(sd) * Shot.SPEED * 0.33, Cos32(sd) * Shot.SPEED * 0.33, 0,
               Smoke.SmokeType.SPARK, 10, 0.33)
      }
    }
    fireCnt--
  }

  func (this *Boat) bool checkBulletHit(Vector p, Vector pp) {
    if (cnt <= 0) {
      return false
		}
    float bmvx, bmvy, inaa
    bmvx = pp.x
    bmvy = pp.y
    bmvx -= p.x
    bmvy -= p.y
    inaa = bmvx * bmvx + bmvy * bmvy
    if (inaa > 0.00001) {
      float sofsx, sofsy, inab, hd
      sofsx = _pos.x
      sofsy = _pos.y
      sofsx -= p.x
      sofsy -= p.y
      inab = bmvx * sofsx + bmvy * sofsy
      if (inab >= 0 && inab <= inaa) {
        hd = sofsx * sofsx + sofsy * sofsy - inab * inab / inaa
        if (hd >= 0 && hd <= HIT_WIDTH) {
          destroyed()
          return true
        }
      }
    }
    return false
  }

  func (this *Boat) destroyed() {
    if (cnt <= 0) {
      return
		}
    if (shieldCnt > 0) {
      destroyedBoatShield()
      return
    }
    ship.destroyed()
    gameState.shipDestroyed()
  }

  func (this *Boat) destroyedBoatShield() {
		for i := 0; i < 100; i++ {
      Spark sp = sparks.getInstanceForced()
      sp.set(pos, rand.nextSignedFloat(1), rand.nextSignedFloat(1),
             0.5 + rand.nextFloat(0.5), 0.5 + rand.nextFloat(0.5), 0,
             40 + rand.nextInt(40))
    }
    SoundManager.playSe("ship_shield_lost.wav")
    screen.setScreenShake(30, 0.02)
    shieldCnt = 0
    cnt = -INVINCIBLE_CNT / 2
  }

  func (this *Boat) destroyedBoat() {
		for i := 0; i < 128; i++ {
      Spark sp = sparks.getInstanceForced()
      sp.set(pos, rand.nextSignedFloat(1), rand.nextSignedFloat(1),
             0.5 + rand.nextFloat(0.5), 0.5 + rand.nextFloat(0.5), 0,
             40 + rand.nextInt(40))
    }
    SoundManager.playSe("ship_destroyed.wav")
		for i := 0; i < 64; i++ {
      Smoke s = smokes.getInstanceForced()
      s.set(pos, rand.nextSignedFloat(0.2), rand.nextSignedFloat(0.2),
            rand.nextFloat(0.1),
            Smoke.SmokeType.EXPLOSION, 50 + rand.nextInt(30), 1)
    }
    screen.setScreenShake(60, 0.05)
    restart()
    cnt = -RESTART_CNT
  }

  func (this *Boat) bool hasCollision() {
    rerturn ! (cnt < -INVINCIBLE_CNT)
  }

  func (this *Boat) draw() {
    if (cnt < -INVINCIBLE_CNT) {
      return
		}
    if (fireDeg < 99999) {
      Screen.setColor(0.5, 0.9, 0.7, 0.4)
      glBegin(GL_LINE_STRIP)
      glVertex2(_pos.x, _pos.y)
      Screen.setColor(0.5, 0.9, 0.7, 0.8)
      glVertex2(_pos.x + Sin32(fireDeg) * 20, _pos.y + Cos32(fireDeg) * 20)
      glEnd()
    }
    if (cnt < 0 && (-cnt % 32) < 16) {
      return
		}
    glPushMatrix()
    Screen.glTranslate(pos)
    glRotatef(-deg * 180 / Pi32, 0, 0, 1)
    _shape.draw()
    bridgeShape.draw()
    if (shieldCnt > 0) {
      float ss = 0.66
      if (shieldCnt < 120)
        ss *= cast(float) shieldCnt / 120
      glScalef(ss, ss, ss)
      glRotatef(shieldCnt * 5, 0, 0, 1)
      shieldShape.draw()
    }
    glPopMatrix()
  }

  func (this *Boat) drawFront() {
    if (cnt < -INVINCIBLE_CNT)
      return
    if (gameMode == InGameState.GameMode.MOUSE) {
      Screen.setColor(0.7, 0.9, 0.8, 1.0)
      Screen.lineWidth(2)
      drawSight(mouseInput.x, mouseInput.y, 0.3)
      float ss = 0.9 - 0.8 * ((cnt + 1024) % 32) / 32
      Screen.setColor(0.5, 0.9, 0.7, 0.8)
      drawSight(mouseInput.x, mouseInput.y, ss)
      Screen.lineWidth(1)
    }
  }

  func (this *Boat) drawSight(float x, float y, float size) {
    glBegin(GL_LINE_STRIP)
    glVertex2(x - size, y - size * 0.5)
    glVertex2(x - size, y - size)
    glVertex2(x - size * 0.5, y - size)
    glEnd()
    glBegin(GL_LINE_STRIP)
    glVertex2(x + size, y - size * 0.5)
    glVertex2(x + size, y - size)
    glVertex2(x + size * 0.5, y - size)
    glEnd()
    glBegin(GL_LINE_STRIP)
    glVertex2(x + size, y + size * 0.5)
    glVertex2(x + size, y + size)
    glVertex2(x + size * 0.5, y + size)
    glEnd()
    glBegin(GL_LINE_STRIP)
    glVertex2(x - size, y + size * 0.5)
    glVertex2(x - size, y + size)
    glVertex2(x - size * 0.5, y + size)
    glEnd()
  }

  func (this *Boat) drawShape() {
    _shape.draw()
    bridgeShape.draw()
  }

  func (this *Boat) clearBullets() {
    gameState.clearBullets()
  }



}
