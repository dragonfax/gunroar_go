package main

import (
	"math"
	"math/rand"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/letter"
)

/**
 * Manage an enemys' appearance, a rank(difficulty) and a field.
 */

const RANK_INC_BASE = 0.0018
const BLOCK_DENSITY_MIN = 0
const BLOCK_DENSITY_MAX = 3

type StageManager struct {
	field                                                                            *Field
	enemies                                                                          *EnemyPool
	ship                                                                             *Ship
	bullets                                                                          *BulletPool
	sparks                                                                           *SparkPool
	smokes                                                                           *SmokePool
	fragments                                                                        *FragmentPool
	wakes                                                                            *WakePool
	rand                                                                             *r.Rand
	rank, baseRank, addRank, rankVel, rankInc                                        float64
	enemyApp                                                                         [3]*EnemyAppearance
	platformEnemySpec                                                                PlatformEnemySpec
	_bossMode                                                                        bool
	bossAppCnt, bossAppTime, bossAppTimeBase, bgmStartCnt, _blockDensity, batteryNum int
}

func NewStageManager(field *Field, enemies *EnemyPool, ship *Ship, bullets *BulletPool, sparks *SparkPool,
	smokes *SmokePool, fragments *FragmentPool, wakes *WakePool) *StageManager {

	this := &StageManager{}
	this.field = field
	this.enemies = enemies
	this.ship = ship
	this.bullets = bullets
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.wakes = wakes
	this.rand = r.New(r.NewSource(time.Now().Unix()))
	for i := range this.enemyApp {
		this.enemyApp[i] = NewEnemyAppearance()
	}
	this.platformEnemySpec = NewPlatformEnemySpec(field, ship, sparks, smokes, fragments, wakes)
	this.rank = 1
	this.baseRank = 1
	this._blockDensity = 2
	return this
}

func (this *StageManager) setRandSeed(seed int64) {
	this.rand = r.New(r.NewSource(seed))
}

func (this *StageManager) start(rankIncRatio float64) {
	this.rank = 1
	this.baseRank = 1
	this.addRank = 0
	this.rankVel = 0
	this.rankInc = RANK_INC_BASE * rankIncRatio
	this._blockDensity = this.rand.Intn(BLOCK_DENSITY_MAX-BLOCK_DENSITY_MIN+1) + BLOCK_DENSITY_MIN
	this._bossMode = false
	this.bossAppTimeBase = 60 * 1000
	this.resetBossMode()
	this.gotoNextBlockArea()
	this.bgmStartCnt = -1
}

func (this *StageManager) startBossMode() {
	this._bossMode = true
	this.bossAppCnt = 2
	fadeBgm()
	this.bgmStartCnt = 120
	this.rankVel = 0
}

func (this *StageManager) resetBossMode() {
	if this._bossMode {
		this._bossMode = false
		fadeBgm()
		this.bgmStartCnt = 120
		this.bossAppTimeBase += 30 * 1000
	}
	this.bossAppTime = this.bossAppTimeBase
}

func (this *StageManager) move() {
	this.bgmStartCnt--
	if this.bgmStartCnt == 0 {
		if this._bossMode {
			playBgmWithName("gr0.ogg")
		} else {
			nextBgm()
		}
	}
	if this._bossMode {
		this.addRank *= 0.999
		if !this.enemies.hasBoss() && this.bossAppCnt <= 0 {
			this.resetBossMode()
		}
	} else {
		rv := this.field.lastScrollY()/this.ship.scrollSpeedBase() - 2
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
		ea.move(this.enemies, this.field)
	}
}

func (this *StageManager) shipDestroyed() {
	this.rankVel = 0
	if !this._bossMode {
		this.addRank = 0
	} else {
		this.addRank /= 2
	}
}

func nextSignedInt(rand *r.Rand, max int) int {
	return rand.Intn(max*2) - max // NOTES skips some numbers
}

func (this *StageManager) gotoNextBlockArea() {
	if this._bossMode {
		this.bossAppCnt--
		if this.bossAppCnt == 0 {
			ses := NewShipEnemySpec(this.field, this.ship, this.sparks, this.smokes, this.fragments, this.wakes)
			ses.setParam(this.rank, BOSS, this.rand)
			en := this.enemies.GetInstance()
			if en != nil {
				if ses.setFirstState(en.state(), CENTER) {
					en.set(ses)
				}
			} else {
				this.resetBossMode()
			}
		}
		for _, ea := range this.enemyApp {
			ea.unset()
		}
		return
	}
	var noSmallShip bool
	if this._blockDensity < BLOCK_DENSITY_MAX && this.rand.Intn(2) == 0 {
		noSmallShip = true
	} else {
		noSmallShip = false
	}
	this._blockDensity += nextSignedInt(this.rand, 1)
	if this._blockDensity < BLOCK_DENSITY_MIN {
		this._blockDensity = BLOCK_DENSITY_MIN
	} else if this._blockDensity > BLOCK_DENSITY_MAX {
		this._blockDensity = BLOCK_DENSITY_MAX
	}
	this.batteryNum = int((float64(this._blockDensity) + nextSignedFloat(this.rand, 1)) * 0.75)
	tr := this.rank
	largeShipNum := int((2 - float64(this._blockDensity) + nextSignedFloat(this.rand, 1)) * 0.5)
	if noSmallShip {
		largeShipNum = int(float64(largeShipNum) * 1.5)
	} else {
		largeShipNum = int(float64(largeShipNum) * 0.5)
	}
	var appType = AppearanceType(this.rand.Intn(2))
	if largeShipNum > 0 {
		lr := tr * (0.25 + nextFloat(this.rand, 0.15))
		if noSmallShip {
			lr *= 1.5
		}
		tr -= lr
		ses := NewShipEnemySpec(this.field, this.ship, this.sparks, this.smokes, this.fragments, this.wakes)
		ses.setParam(lr/float64(largeShipNum), LARGE, this.rand)
		this.enemyApp[0].set(ses, largeShipNum, appType, this.rand)
	} else {
		this.enemyApp[0].unset()
	}
	if this.batteryNum > 0 {
		this.platformEnemySpec = NewPlatformEnemySpec(this.field, this.ship, this.sparks, this.smokes, this.fragments, this.wakes)
		pr := tr * (0.3 + nextFloat(this.rand, 0.1))
		this.platformEnemySpec.setParam(pr/float64(this.batteryNum), this.rand)
	}
	appType = (appType + 1) % 2
	middleShipNum := int((4 - float64(this._blockDensity) + nextSignedFloat(this.rand, 1)) * 0.66)
	if noSmallShip {
		middleShipNum *= 2
	}
	if middleShipNum > 0 {
		var mr float64
		if noSmallShip {
			mr = tr
		} else {
			mr = tr * (0.33 + nextFloat(this.rand, 0.33))
		}
		tr -= mr
		ses := NewShipEnemySpec(this.field, this.ship, this.sparks, this.smokes, this.fragments, this.wakes)
		ses.setParam(mr/float64(middleShipNum), MIDDLE, this.rand)
		this.enemyApp[1].set(ses, middleShipNum, appType, this.rand)
	} else {
		this.enemyApp[1].unset()
	}
	if !noSmallShip {
		appType = TOP
		smallShipNum := int(math.Sqrt(3+tr)*(1+nextSignedFloat(this.rand, 0.5))*2) + 1
		if smallShipNum > 256 {
			smallShipNum = 256
		}
		sses := NewSmallShipEnemySpec(this.field, this.ship, this.sparks, this.smokes, this.fragments, this.wakes)
		sses.setParam(tr/float64(smallShipNum), this.rand)
		this.enemyApp[2].set(sses, smallShipNum, appType, this.rand)
	} else {
		this.enemyApp[2].unset()
	}
}

func (this *StageManager) addBatteries(platformPos []PlatformPos, platformPosNum int) {
	ppn := platformPosNum
	bn := this.batteryNum
	for i := 0; i < 100; i++ {
		if ppn <= 0 || bn <= 0 {
			break
		}
		ppi := rand.Intn(platformPosNum)
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
		en := this.enemies.GetInstance()
		if en == nil {
			break
		}
		platformPos[ppi].used = true
		ppn--
		p := this.field.convertToScreenPos(int(platformPos[ppi].pos.X), int(platformPos[ppi].pos.Y))
		if !this.platformEnemySpec.setFirstState(en.state(), p.X, p.Y, platformPos[ppi].deg) {
			continue
		}
		for i := 0; i < platformPosNum; i++ {
			if math.Abs(platformPos[ppi].pos.X-platformPos[i].pos.X) <= 1 &&
				math.Abs(platformPos[ppi].pos.Y-platformPos[i].pos.Y) <= 1 &&
				!platformPos[i].used {
				platformPos[i].used = true
				ppn--
			}
		}
		en.set(this.platformEnemySpec)
		bn--
	}
}

func (this *StageManager) blockDensity() int {
	return this._blockDensity
}

func (this *StageManager) draw() {
	letter.DrawNum(int(this.rank)*1000, 620, 10, 10, 0, 0, 33, 3)
	letter.DrawTime(this.bossAppTime, 120, 20, 7, 0)
}

func (this *StageManager) rankMultiplier() float64 {
	return this.rank
}

func (this *StageManager) bossMode() bool {
	return this._bossMode
}

type EnemyAppearance struct {
	spec                             EnemySpec
	nextAppDist, nextAppDistInterval float64
	appType                          AppearanceType
}

func NewEnemyAppearance() *EnemyAppearance {
	this := &EnemyAppearance{}
	this.nextAppDistInterval = 1
	return this
}

func (this *EnemyAppearance) set(s EnemySpec, num int, appType AppearanceType, rand *r.Rand) {
	this.spec = s
	this.nextAppDistInterval = NEXT_BLOCK_AREA_SIZE / float64(num)
	this.nextAppDist = nextFloat(rand, this.nextAppDistInterval)
	this.appType = appType
}

func (this *EnemyAppearance) unset() {
	this.spec = nil
}

func (this *EnemyAppearance) move(enemies *EnemyPool, field *Field) {
	if this.spec == nil {
		return
	}
	this.nextAppDist -= field.lastScrollY()
	if this.nextAppDist <= 0 {
		this.nextAppDist += this.nextAppDistInterval
		this.appear(enemies)
	}
}

func (this *EnemyAppearance) appear(enemies *EnemyPool) {
	en := enemies.GetInstance()
	if en != nil {
		if this.spec.(HasAppearType).setFirstState(en.state(), this.appType) {
			en.set(this.spec)
		}
	}
}
