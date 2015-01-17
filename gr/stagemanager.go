/*
 * $Id: stagemanager.d,v 1.2 2005/07/03 07:05:22 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

/**
 * Manage an enemys' appearance, a rank(difficulty) and a field.
 */
const RANK_INC_BASE = 0.0018
const BLOCK_DENSITY_MIN = 0
const BLOCK_DENSITY_MAX = 3

type StageManager struct {
	rank, baseRank, addRank, rankVel, rankInc float32
	enemyApp                                  [3]*EnemyAppearance
	blockDensity                              int
	batteryNum                                int
	platformEnemySpec                         PlatformEnemySpec
	bossMode                                  bool
	bossAppCnt                                int
	bossAppTime, bossAppTimeBase              int
	bgmStartCnt                               int
}

func NewStageManager() *StageManager {
	this.ship = ship
	for i, _ := range this.enemyApp {
		this.enemyApp[i] = NewEnemyAppearance()
	}
	this.platformEnemySpec = NewPlatformEnemySpec()
	this.rank = 1
	this.baseRank = 1
	this.blockDensity = 2
}

func (this *StageManager) start(rankIncRatio float32) {
	this.rank = 1
	this.baseRank = 1
	this.addRank = 0
	this.rankVel = 0
	this.rankInc = RANK_INC_BASE * this.rankIncRatio
	this.blockDensity = Int(BLOCK_DENSITY_MAX-BLOCK_DENSITY_MIN+1) + BLOCK_DENSITY_MIN
	this.bossMode = false
	this.bossAppTimeBase = 60 * 1000
	this.resetBossMode()
	this.gotoNextBlockArea()
	this.bgmStartCnt = -1
}

func (this *StageManager) startBossMode() {
	this.bossMode = true
	this.bossAppCnt = 2
	fadeBgm()
	this.bgmStartCnt = 120
	this.rankVel = 0
}

func (this *StageManager) resetBossMode() {
	if this.bossMode {
		this.bossMode = false
		fadeBgm()
		this.bgmStartCnt = 120
		this.bossAppTimeBase += 30 * 1000
	}
	this.bossAppTime = this.bossAppTimeBase
}

func (this *StageManager) move() {
	this.bgmStartCnt--
	if this.bgmStartCnt == 0 {
		if this.bossMode {
			playBgm("gr0.ogg")
		} else {
			nextBgm()
		}
	}
	if _bossMode {
		this.addRank *= 0.999
		if !enemies.hasBoss() && this.bossAppCnt <= 0 {
			this.resetBossMode()
		}
	} else {
		rv := field.lastScrollY/this.ship.scrollSpeedBase - 2
		this.bossAppTime -= 17
		if this.bossAppTime <= 0 {
			this.bossAppTime = 0
			this.startBossMode()
		}
		if rv > 0 {
			this.rankVel += rv * rv * 0.0004 * this.baseRank
		} else {
			this.rankVel += rv * this.baseRank
			if this.rankVel < 0 {
				this.rankVel = 0
			}
		}
		this.addRank += this.rankInc * (this.rankVel + 1)
		this.addRank *= 0.999
		this.baseRank += this.rankInc + this.addRank*0.0001
	}
	this.rank = this.baseRank + this.addRank
	for _, ea := range this.enemyApp {
		ea.move()
	}
}

func (this *StageManager) shipDestroyed() {
	this.rankVel = 0
	if !this.bossMode {
		this.addRank = 0
	} else {
		this.addRank /= 2
	}
}

func (this *StageManager) gotoNextBlockArea() {
	if this.bossMode {
		this.bossAppCnt--
		if this.bossAppCnt == 0 {
			ses := NewShipEnemySpec()
			ses.setParam(rank, ShipEnemySpec.ShipClass.BOSS)
			en := NewEnemy()
			if ses.setFirstState(en.state, EnemyState.AppearanceType.CENTER) {
				en.set(ses)
			}
		}
		for _, ea := range this.enemyApp {
			ea.close()
		}
		return
	}
	var noSmallShip bool
	if this.blockDensity < BLOCK_DENSITY_MAX && Int(2) == 0 {
		noSmallShip = true
	} else {
		noSmallShip = false
	}
	this.blockDensity += SignedInt(1)
	if this.blockDensity < BLOCK_DENSITY_MIN {
		this.blockDensity = BLOCK_DENSITY_MIN
	} else if this.blockDensity > BLOCK_DENSITY_MAX {
		this.blockDensity = BLOCK_DENSITY_MAX
	}
	this.batteryNum = (this.blockDensity + SignedFloat(1)) * 0.75
	tr := this.rank
	largeShipNum := (2 - this.blockDensity + SignedFloat(1)) * 0.5
	if noSmallShip {
		largeShipNum *= 1.5
	} else {
		largeShipNum *= 0.5
	}
	appType := Int(2)
	if largeShipNum > 0 {
		lr := tr * (0.25 + nextFloat(0.15))
		if noSmallShip {
			lr *= 1.5
		}
		tr -= lr
		ses := NewShipEnemySpec()
		ses.setParam(lr/largeShipNum, ShipEnemySpec.ShipClass.LARGE)
		this.enemyApp[0].set(ses, largeShipNum, appType)
	} else {
		this.enemyApp[0].close()
	}
	if batteryNum > 0 {
		this.platformEnemySpec = NewPlatformEnemySpec()
		pr := tr * (0.3 + nextFloat(0.1))
		this.platformEnemySpec.setParam(pr / batteryNum)
	}
	appType = (appType + 1) % 2
	middleShipNum := (4 - _blockDensity + nextSignedFloat(1)) * 0.66
	if noSmallShip {
		middleShipNum *= 2
	}
	if middleShipNum > 0 {
		var mr float32
		if noSmallShip {
			mr = tr
		} else {
			mr = tr * (0.33 + nextFloat(0.33))
		}
		tr -= mr
		ses := NewShipEnemySpec()
		ses.setParam(mr/middleShipNum, ShipEnemySpec.ShipClass.MIDDLE)
		this.enemyApp[1].set(ses, middleShipNum, appType)
	} else {
		this.enemyApp[1].close()
	}
	if !noSmallShip {
		appType = EnemyState.AppearanceType.TOP
		smallShipNum := (sqrt(3+tr) * (1 + nextSignedFloat(0.5)) * 2) + 1
		if smallShipNum > 256 {
			smallShipNum = 256
		}
		sses := NewSmallShipEnemySpec()
		sses.setParam(tr / smallShipNum)
		enemyApp[2].set(sses, smallShipNum, appType)
	} else {
		enemyApp[2].unset()
	}
}

func (this *StageManager) addBatteries(platformPos []PlatformPos, platformPosNum int) {
	ppn := platformPosNum
	bn := batteryNum
	for i := 0; i < 100; i++ {
		if ppn <= 0 || bn <= 0 {
			break
		}
		ppi := Int(platformPosNum)
		for j := 0; j < platformPosNum; j++ {
			if !platformPos[ppi].used {
				break
			}
			ppi++
			if ppi >= platformPosNum {
				ppi = 0
			}
		}
		if platformPos[ppi].used {
			break
		}
		en := NewEnemy()
		platformPos[ppi].used = true
		ppn--
		p := field.convertToScreenPos(platformPos[ppi].pos.x, platformPos[ppi].pos.y)
		if !platformEnemySpec.setFirstState(en.state, p.x, p.y, platformPos[ppi].deg) {
			continue
		}
		for i := 0; i < platformPosNum; i++ {
			if fabs32(platformPos[ppi].pos.x-platformPos[i].pos.x) <= 1 &&
				fabs32(platformPos[ppi].pos.y-platformPos[i].pos.y) <= 1 &&
				!platformPos[i].used {
				platformPos[i].used = true
				ppn--
			}
		}
		en.set(platformEnemySpec)
		bn--
	}
}

func (this *StageManager) draw() {
	drawNum(this.rank*1000, 620, 10, 10, 0, 0, 33, 3)
	drawTime(this.bossAppTime, 120, 20, 7)
}

type EnemyAppearance struct {
	spec                             EnemySpec
	nextAppDist, nextAppDistInterval float32
	appType                          int
}

func NewEnemyAppearance(s EnemySpec, num int, appType int) *EnemyAppearance {
	this := new(EnemyAppearance)
	this.nextAppDistInterval = 1
	this.spec = s
	this.nextAppDistInterval = NEXT_BLOCK_AREA_SIZE / num
	this.nextAppDist = SignedFloat(nextAppDistInterval)
	this.appType = appType
	return this
}

func (this *EnemyAppearance) move() {
	if spec == nil {
		return
	}
	this.nextAppDist -= field.lastScrollY
	if this.nextAppDist <= 0 {
		this.nextAppDist += this.nextAppDistInterval
		this.appear()
	}
}

func (this *EnemyAppearance) appear() {
	if this.spec.setFirstState(en.state, this.appType) {
		NewEnemy(this.spec)
	}
}
