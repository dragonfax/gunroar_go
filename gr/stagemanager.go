/*
 * $Id: stagemanager.d,v 1.2 2005/07/03 07:05:22 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Manage an enemys' appearance, a rank(difficulty) and a field.
 */
const RANK_INC_BASE = 0.0018
const BLOCK_DENSITY_MIN = 0
const BLOCK_DENSITY_MAX = 3

type StageManager struct {
  field Field
  ship Ship
  rank, baseRank, addRank, rankVel, rankInc float32
  enemyApp [3]*EnemyAppearance
  blockDensity int
  batteryNum int
  platformEnemySpec PlatformEnemySpec
  bossMode bool
  bossAppCnt int
  bossAppTime, bossAppTimeBase int
  bgmStartCnt int
}

func NewStateManager(field Field, ship Ship) {
	this.field = field
	this.ship = ship
	for i,_ := range this.enemyApp {
		this.enemyApp[i] = NewEnemyAppearance()
	}
	this.platformEnemySpec = NewPlatformEnemySpec(field, ship)
	this.rank =  1
	this.baseRank = 1
	this.blockDensity = 2
}

func (this *StateManager)  start(rankIncRatio float32) {
	this.rank = 1
	this.baseRank = 1
	this.addRank = 0
	this.rankVel = 0
	this.rankInc = RANK_INC_BASE * this.rankIncRatio
	this.blockDensity = rand.Int(BLOCK_DENSITY_MAX - BLOCK_DENSITY_MIN + 1) + BLOCK_DENSITY_MIN
	this.bossMode = false
	this.bossAppTimeBase = 60 * 1000
	this.resetBossMode()
	this.gotoNextBlockArea()
	this.bgmStartCnt = -1
}

func (this *StateManager)  startBossMode() {
	this.bossMode = true
	this.bossAppCnt = 2
	fadeBgm()
	this.bgmStartCnt = 120
	this.rankVel = 0
}

func (this *StateManager)  resetBossMode() {
	if (this.bossMode) {
		this.bossMode = false
		fadeBgm()
		this.bgmStartCnt = 120
		this.bossAppTimeBase += 30 * 1000
	}
	this.bossAppTime = this.bossAppTimeBase
}

func (this *StageManager) move() {
	this.bgmStartCnt--
	if (this.bgmStartCnt == 0) {
		if (this.bossMode) {
			playBgm("gr0.ogg")
		}
		else {
			nextBgm()
		}
	}
	if (_bossMode) {
		this.addRank *= 0.999
		if (!enemies.hasBoss() && this.bossAppCnt <= 0) {
			this.resetBossMode()
		}
	} else {
		rv := this.field.lastScrollY / this.ship.scrollSpeedBase - 2
		this.bossAppTime -= 17
		if (this.bossAppTime <= 0) {
			this.bossAppTime = 0
			this.startBossMode()
		}
		if (rv > 0) {
			this.rankVel += rv * rv * 0.0004 * this.baseRank
		} else {
			this.rankVel += rv * this.baseRank
			if (this.rankVel < 0) {
				this.rankVel = 0
			}
		}
		this.addRank += this.rankInc * (this.rankVel + 1)
		this.addRank *= 0.999
		this.baseRank += this.rankInc + this.addRank * 0.0001
	}
	this.rank = this.baseRank + this.addRank
	for _,ea : range this.enemyApp {
		ea.move(this.field)
	}
}

func (this *StageManager) shipDestroyed() {
	this.rankVel = 0
	if (!this.bossMode) {
		this.addRank = 0
	}
	else {
		this.addRank /= 2
	}
}

func (this *StageManager)  gotoNextBlockArea() {
	if (this.bossMode) {
		this.bossAppCnt--
		if (this.bossAppCnt == 0) {
			ses := NewShipEnemySpec(this.field, this.ship)
			ses.setParam(rank, ShipEnemySpec.ShipClass.BOSS, rand)
			en = NewEnemy()
			if ((cast(HasAppearType) ses).setFirstState(en.state, EnemyState.AppearanceType.CENTER)) {
				en.set(ses)
			}
		}
		for _,ea := range this.enemyApp) {
			ea.remove()
		}
		return
	}
	noSmallShip bool
	if (this.blockDensity < BLOCK_DENSITY_MAX && rand.Int(2) == 0) {
		noSmallShip = true
	} else {
		noSmallShip = false
	}
	this.blockDensity += rand.SignedInt(1)
	if (this.blockDensity < BLOCK_DENSITY_MIN) {
		this.blockDensity = BLOCK_DENSITY_MIN
	} else if (this.blockDensity > BLOCK_DENSITY_MAX) {
		this.blockDensity = BLOCK_DENSITY_MAX
	}
	this.batteryNum = (this.blockDensity + rand.SignedFloat(1)) * 0.75
	tr := this.rank
	largeShipNum := (2 - this.blockDensity + rand.SignedFloat(1)) * 0.5
	if (noSmallShip) {
		largeShipNum *= 1.5
	} else {
		largeShipNum *= 0.5
	}
	appType := rand.Int(2)
	if (largeShipNum > 0) {
		lr := tr * (0.25 + rand.nextFloat(0.15))
		if (noSmallShip) {
			lr *= 1.5
		}
		tr -= lr
		ses := NewShipEnemySpec(field, ship)
		ses.setParam(lr / largeShipNum, ShipEnemySpec.ShipClass.LARGE, rand)
		this.enemyApp[0].set(ses, largeShipNum, appType, rand)
	} else {
		this.enemyApp[0].remove()
	}
	if (batteryNum > 0) {
		this.platformEnemySpec = NewPlatformEnemySpec(field, ship, sparks, smokes, fragments, wakes)
		pr := tr * (0.3f + rand.nextFloat(0.1))
		this.platformEnemySpec.setParam(pr / batteryNum, rand)
	}
	appType = (appType + 1) % 2
	middleShipNum := ((4 - _blockDensity + rand.nextSignedFloat(1)) * 0.66
	if (noSmallShip) {
		middleShipNum *= 2
	}
	if (middleShipNum > 0) {
		var mr float32
		if (noSmallShip) {
			mr = tr
		} else {
			mr = tr * (0.33f + rand.nextFloat(0.33f))
		}
		tr -= mr
		ses := NewShipEnemySpec(field, ship)
		ses.setParam(mr / middleShipNum, ShipEnemySpec.ShipClass.MIDDLE, rand)
		this.enemyApp[1].set(ses, middleShipNum, appType, rand)
	} else {
		this.enemyApp[1].remove()
	}
	if (!noSmallShip) {
		appType = EnemyState.AppearanceType.TOP
		smallShipNum :=  (sqrt(3 + tr) * (1 + rand.nextSignedFloat(0.5f)) * 2) + 1
		if (smallShipNum > 256) {
			smallShipNum = 256
		}
		SmallShipEnemySpec sses = new SmallShipEnemySpec(field, ship, sparks, smokes, fragments, wakes)
		sses.setParam(tr / smallShipNum, rand)
		enemyApp[2].set(sses, smallShipNum, appType, rand)
	} else {
		enemyApp[2].unset()
	}
}

func (this *StageManager)  addBatteries(platformPos []PlatformPos, platformPosNum int) {
	ppn := platformPosNum
	bn := batteryNum
	for i := 0; i < 100; i++ {
		if (ppn <= 0 || bn <= 0) {
			break
		}
		ppi := rand.Int(platformPosNum)
		for j := 0; j < platformPosNum; j++ {
			if (!platformPos[ppi].used) {
				break
			}
			ppi++
			if (ppi >= platformPosNum) {
				ppi = 0
			}
		}
		if (platformPos[ppi].used) {
			break
		}
		en := NewEnemy()
		platformPos[ppi].used = true
		ppn--
		p := this.field.convertToScreenPos(platformPos[ppi].pos.x, platformPos[ppi].pos.y)
		if (!platformEnemySpec.setFirstState(en.state, p.x, p.y, platformPos[ppi].deg)) {
			continue
		}
		for i := 0; i < platformPosNum; i++ {
			if (fabs32(platformPos[ppi].pos.x - platformPos[i].pos.x) <= 1 &&
					fabs32(platformPos[ppi].pos.y - platformPos[i].pos.y) <= 1 &&
					!platformPos[i].used) {
				platformPos[i].used = true
				ppn--
			}
		}
		en.set(platformEnemySpec)
		bn--
	}
}

func (this *StageManager) draw() {
	drawNum(this.rank * 1000, 620, 10, 10, 0, 0, 33, 3)
	drawTime(this.bossAppTime, 120, 20, 7)
}


type EnemyAppearance struct {
	spec EnemySpec
	nextAppDist, nextAppDistInterval float32
	appType int
}

func NewEnemyAppearance( s EnemySpec, num int, appType int) *EnemyAppearance {
	this := new(EnemyAppearance)
	this.nextAppDistInterval = 1
	this.spec = s
	this.nextAppDistInterval = NEXT_BLOCK_AREA_SIZE / num
	this.nextAppDist = rand.SignedFloat(nextAppDistInterval)
	this.appType = appType
	return this
}

func (this *EnemyAppearance)  move(field Field) {
	if (spec == nil ) {
		return
	}
	this.nextAppDist -= this.field.lastScrollY
	if (this.nextAppDist <= 0) {
		this.nextAppDist += this.nextAppDistInterval
		this.appear()
	}
}

func (this *EnemyAppearance) appear() {
	if this.spec.setFirstState(en.state, this.appType) {
		NewEnemy(this.spec)
	}
}