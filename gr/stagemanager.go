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
	bossMode                                  bool
	bossAppCnt                                int
	bossAppTime, bossAppTimeBase              int
	bgmStartCnt                               int
	platformRank                              float32
}

func NewStageManager() *StageManager {
	this := new(StageManager)
	/*for i, _ := range this.enemyApp {
		this.enemyApp[i] = NewEnemyAppearance()
	} */
	this.rank = 1
	this.baseRank = 1
	this.blockDensity = 2
	return this
}

func (this *StageManager) start(rankIncRatio float32) {
	this.rank = 1
	this.baseRank = 1
	this.addRank = 0
	this.rankVel = 0
	this.rankInc = RANK_INC_BASE * rankIncRatio
	this.blockDensity = nextInt(BLOCK_DENSITY_MAX-BLOCK_DENSITY_MIN+1) + BLOCK_DENSITY_MIN
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
			playBgmByName("gr0.ogg")
		} else {
			nextBgm()
		}
	}
	if this.bossMode {
		this.addRank *= 0.999
		if !hasBoss() && this.bossAppCnt <= 0 {
			this.resetBossMode()
		}
	} else {
		rv := field.lastScrollY/ship.scrollSpeedBase - 2
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
		if ea != nil { // TODO ensure this is only nil when it should be
			ea.move()
		}
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
			ses.setParam(this.rank, ShipClassBOSS)
			en := NewEnemy(ses)
			if !ses.setFirstState(en.state, AppearanceTypeCENTER, 0, 0, 0) {
				en.close()
			}
		}
		for i, _ := range this.enemyApp {
			// ea.close()
			this.enemyApp[i] = nil
		}
		return
	}
	var noSmallShip bool
	if this.blockDensity < BLOCK_DENSITY_MAX && nextInt(2) == 0 {
		noSmallShip = true
	} else {
		noSmallShip = false
	}
	this.blockDensity += nextSignedInt(1)
	if this.blockDensity < BLOCK_DENSITY_MIN {
		this.blockDensity = BLOCK_DENSITY_MIN
	} else if this.blockDensity > BLOCK_DENSITY_MAX {
		this.blockDensity = BLOCK_DENSITY_MAX
	}
	this.batteryNum = int((float32(this.blockDensity) + nextSignedFloat(1)) * 0.75)
	tr := this.rank
	largeShipNum := (2 - float32(this.blockDensity) + nextSignedFloat(1)) * 0.5
	if noSmallShip {
		largeShipNum *= 1.5
	} else {
		largeShipNum *= 0.5
	}
	appType := AppearanceType(nextInt(2))
	if largeShipNum > 0 {
		lr := tr * (0.25 + nextFloat(0.15))
		if noSmallShip {
			lr *= 1.5
		}
		tr -= lr
		ses := NewShipEnemySpec()
		ses.setParam(lr/largeShipNum, ShipClassLARGE)
		this.enemyApp[0] = NewEnemyAppearance(ses, int(largeShipNum), appType)
	} else {
		this.enemyApp[0] = nil
	}
	if this.batteryNum > 0 {
		pr := tr * (0.3 + nextFloat(0.1))
		this.platformRank = pr / float32(this.batteryNum)
	}
	appType = (appType + 1) % 2
	middleShipNum := (4 - float32(this.blockDensity) + nextSignedFloat(1)) * 0.66
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
		ses.setParam(mr/middleShipNum, ShipClassMIDDLE)
		this.enemyApp[1] = NewEnemyAppearance(ses, int(middleShipNum), appType)
	} else {
		this.enemyApp[1] = nil
	}
	if !noSmallShip {
		appType = AppearanceTypeTOP
		smallShipNum := (sqrt32(3+tr) * (1 + nextSignedFloat(0.5)) * 2) + 1
		if smallShipNum > 256 {
			smallShipNum = 256
		}
		sses := NewSmallShipEnemySpec()
		sses.setParam(tr / smallShipNum)
		this.enemyApp[2] = NewEnemyAppearance(sses, int(smallShipNum), appType)
	} else {
		this.enemyApp[2] = nil
	}
}

func (this *StageManager) addBatteries(platformPos []PlatformPos, platformPosNum int) {
	ppn := platformPosNum
	bn := this.batteryNum
	for i := 0; i < 100; i++ {
		if ppn <= 0 || bn <= 0 {
			break
		}
		ppi := nextInt(platformPosNum)
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
		en := NewEnemy(NewPlatformEnemySpec(this.platformRank))
		platformPos[ppi].used = true
		ppn--
		p := field.convertToScreenPos(int(platformPos[ppi].pos.x), int(platformPos[ppi].pos.y))
		if !en.spec.setFirstState(en.state, 0, p.x, p.y, platformPos[ppi].deg) {
			en.close()
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
		bn--
	}
}

func (this *StageManager) draw() {
	drawNumOption(int(this.rank*1000), 620, 10, 10, 0, 0, 33, 3)
	drawTime(this.bossAppTime, 120, 20, 7, 0)
}

type EnemyAppearance struct {
	spec                             EnemySpec
	nextAppDist, nextAppDistInterval float32
	appType                          AppearanceType
}

func NewEnemyAppearance(s EnemySpec, num int, appType AppearanceType) *EnemyAppearance {
	if num == 0 {
		num = 1
	}
	this := new(EnemyAppearance)
	this.nextAppDistInterval = 1
	this.spec = s
	this.nextAppDistInterval = float32(NEXT_BLOCK_AREA_SIZE / num)
	this.nextAppDist = nextSignedFloat(this.nextAppDistInterval)
	this.appType = appType
	return this
}

func (this *EnemyAppearance) move() {
	if this.spec == nil {
		return
	}
	this.nextAppDist -= field.lastScrollY
	if this.nextAppDist <= 0 {
		this.nextAppDist += this.nextAppDistInterval
		this.appear()
	}
}

func (this *EnemyAppearance) appear() {
	en := NewEnemy(this.spec)
	if !this.spec.setFirstState(en.state, this.appType, 0, 0, 0) {
		en.close()
	}
}
