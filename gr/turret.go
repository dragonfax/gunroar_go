package main

import (
	"math"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Turret mounted on a deck of an enemy ship.
 */

var turretRand = r.New(r.NewSource(time.Now().Unix()))
var turretDamagedPos vector.Vector

type Turret struct {
	field                         *Field
	bullets                       *BulletPool
	ship                          *Ship
	sparks                        *SparkPool
	smokes                        *SmokePool
	fragments                     *FragmentPool
	spec                          TurretSpec
	pos                           vector.Vector
	deg, baseDeg                  float64
	cnt, appCnt, startCnt, shield int
	damaged                       bool
	destroyedCnt, damagedCnt      int
	bulletSpeed                   float64
	burstCnt                      int
	parent                        *Enemy
}

func setTurretRandSeed(seed int64) {
	turretRand = r.New(r.NewSource(seed))
}

func NewTurret(field *Field, bullets *BulletPool, ship *Ship,
	sparks *SparkPool, smokes *SmokePool, fragments *FragmentPool,
	parent *Enemy) Turret {
	this := Turret{}
	this.field = field
	this.bullets = bullets
	this.ship = ship
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.parent = parent
	this.bulletSpeed = 1
	return this
}

func (this *Turret) start(spec TurretSpec) {
	this.spec = spec
	this.shield = spec.shield
	this.appCnt = 0
	this.cnt = 0
	this.startCnt = 0
	this.deg = 0
	this.baseDeg = 0
	this.damaged = false
	this.damagedCnt = 0
	this.destroyedCnt = -1
	this.bulletSpeed = 1
	this.burstCnt = 0
}

func (this *Turret) move(x, y, d float64, bulletFireSpeed float64 /* = 0 */, bulletFireDeg float64 /* = -99999 */) bool {
	this.pos.X = x
	this.pos.Y = y
	this.baseDeg = d
	if this.destroyedCnt >= 0 {
		this.destroyedCnt++
		itv := 5 + this.destroyedCnt/12
		if itv < 60 && this.destroyedCnt%itv == 0 {
			s := this.smokes.GetInstance()
			if s != nil {
				s.setVector(this.pos, 0, 0, 0.01+nextFloat(turretRand, 0.01), FIRE, 90+turretRand.Intn(30), this.spec.size())
			}
		}
		return false
	}
	td := this.baseDeg + this.deg
	shipPos := this.ship.nearPos(this.pos)
	shipVel := this.ship.nearVel(this.pos)
	ax := shipPos.X - this.pos.X
	ay := shipPos.Y - this.pos.Y
	if this.spec.lookAheadRatio != 0 {
		rd := this.pos.DistVector(shipPos) / this.spec.speed * 1.2
		ax += shipVel.X * this.spec.lookAheadRatio * rd
		ay += shipVel.Y * this.spec.lookAheadRatio * rd
	}
	var ad float64
	if math.Abs(ax)+math.Abs(ay) < 0.1 {
		ad = 0
	} else {
		ad = math.Atan2(ax, ay)
	}
	od := td - ad
	od = normalizeDeg(od)
	var ts float64
	if this.cnt >= 0 {
		ts = this.spec.turnSpeed
	} else {
		ts = this.spec.turnSpeed * this.spec.burstTurnRatio
	}
	if math.Abs(od) <= ts {
		this.deg = ad - this.baseDeg
	} else if od > 0 {
		this.deg -= ts
	} else {
		this.deg += ts
	}
	this.deg = normalizeDeg(this.deg)
	if this.deg > this.spec.turnRange {
		this.deg = this.spec.turnRange
	} else if this.deg < -this.spec.turnRange {
		this.deg = -this.spec.turnRange
	}
	this.cnt++
	if this.field.checkInFieldVector(this.pos) || (this.parent.isBoss() && this.cnt%4 == 0) {
		this.appCnt++
	}
	if this.cnt >= this.spec.interval {
		if this.spec.blind || (math.Abs(od) <= this.spec.turnSpeed &&
			this.pos.DistVector(shipPos) < this.spec.maxRange*1.1 &&
			this.pos.DistVector(shipPos) > this.spec.minRange) {
			this.cnt = -(this.spec.burstNum - 1) * this.spec.burstInterval
			this.bulletSpeed = this.spec.speed
			this.burstCnt = 0
		}
	}
	if this.cnt <= 0 && -this.cnt%this.spec.burstInterval == 0 &&
		((this.spec.invisible && this.field.checkInFieldVector(this.pos)) ||
			(this.spec.invisible && this.parent.isBoss() && this.field.checkInOuterFieldVector(this.pos)) ||
			(!this.spec.invisible && this.field.checkInFieldExceptTop(this.pos))) &&
		this.pos.DistVector(shipPos) > this.spec.minRange {
		bd := this.baseDeg + this.deg
		s := this.smokes.GetInstance()
		if s != nil {
			s.setVector(this.pos, math.Sin(bd)*this.bulletSpeed, math.Cos(bd)*this.bulletSpeed, 0,
				SPARK, 20, this.spec.size()*2)
		}
		nw := this.spec.nway
		if this.spec.nwayChange && this.burstCnt%2 == 1 {
			nw--
		}
		bd -= this.spec.nwayAngle * float64(nw-1) / 2
		for i := 0; i < nw; i++ {
			b := this.bullets.GetInstance()
			if b == nil {
				break
			}
			b.set(this.parent.index(),
				this.pos, bd, this.bulletSpeed, this.spec.size()*3, this.spec.bulletShape, this.spec.maxRange,
				bulletFireSpeed, bulletFireDeg, this.spec.bulletDestructive)
			bd += this.spec.nwayAngle
		}
		this.bulletSpeed += this.spec.speedAccel
		this.burstCnt++
	}
	this.damaged = false
	if this.damagedCnt > 0 {
		this.damagedCnt--
	}
	this.startCnt++
	return true
}

func (this *Turret) draw() {
	if this.spec.invisible {
		return
	}
	gl.PushMatrix()
	if this.destroyedCnt < 0 && this.damagedCnt > 0 {
		turretDamagedPos.X = this.pos.X + nextSignedFloat(turretRand, float64(this.damagedCnt)*0.015)
		turretDamagedPos.Y = this.pos.Y + nextSignedFloat(turretRand, float64(this.damagedCnt)*0.015)
		sdl.GlTranslate(turretDamagedPos)
	} else {
		sdl.GlTranslate(this.pos)
	}
	gl.Rotated(-(this.baseDeg+this.deg)*180/math.Pi, 0, 0, 1)
	if this.destroyedCnt >= 0 {
		this.spec.destroyedShape.Draw()
	} else if !this.damaged {
		this.spec.shape.Draw()
	} else {
		this.spec.damagedShape.Draw()
	}
	gl.PopMatrix()
	if this.destroyedCnt >= 0 {
		return
	}
	if this.appCnt > 120 {
		return
	}
	a := 1 - float64(this.appCnt)/120
	if this.startCnt < 12 {
		a = float64(this.startCnt) / 12
	}
	td := this.baseDeg + this.deg
	if this.spec.nway <= 1 {
		gl.Begin(gl.LINE_STRIP)
		sdl.SetColor(0.9, 0.1, 0.1, a)
		gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.minRange, this.pos.Y+math.Cos(td)*this.spec.minRange)
		sdl.SetColor(0.9, 0.1, 0.1, a*0.5)
		gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.maxRange, this.pos.Y+math.Cos(td)*this.spec.maxRange)
		gl.End()
	} else {
		td -= this.spec.nwayAngle * float64(this.spec.nway-1) / 2
		gl.Begin(gl.LINE_STRIP)
		sdl.SetColor(0.9, 0.1, 0.1, a*0.75)
		gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.minRange, this.pos.Y+math.Cos(td)*this.spec.minRange)
		sdl.SetColor(0.9, 0.1, 0.1, a*0.25)
		gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.maxRange, this.pos.Y+math.Cos(td)*this.spec.maxRange)
		gl.End()
		gl.Begin(gl.QUADS)
		for i := 0; i < this.spec.nway-1; i++ {
			sdl.SetColor(0.9, 0.1, 0.1, a*0.3)
			gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.minRange, this.pos.Y+math.Cos(td)*this.spec.minRange)
			sdl.SetColor(0.9, 0.1, 0.1, a*0.05)
			gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.maxRange, this.pos.Y+math.Cos(td)*this.spec.maxRange)
			td += this.spec.nwayAngle
			gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.maxRange, this.pos.Y+math.Cos(td)*this.spec.maxRange)
			sdl.SetColor(0.9, 0.1, 0.1, a*0.3)
			gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.minRange, this.pos.Y+math.Cos(td)*this.spec.minRange)
		}
		gl.End()
		gl.Begin(gl.LINE_STRIP)
		sdl.SetColor(0.9, 0.1, 0.1, a*0.75)
		gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.minRange, this.pos.Y+math.Cos(td)*this.spec.minRange)
		sdl.SetColor(0.9, 0.1, 0.1, a*0.25)
		gl.Vertex2d(this.pos.X+math.Sin(td)*this.spec.maxRange, this.pos.Y+math.Cos(td)*this.spec.maxRange)
		gl.End()
	}
}

func (this *Turret) checkCollision(x, y float64, c sdl.Collidable, shot *Shot) bool {
	if this.destroyedCnt >= 0 || this.spec.invisible {
		return false
	}
	ox := math.Abs(this.pos.X - x)
	oy := math.Abs(this.pos.Y - y)
	if this.spec.shape.CheckCollision(ox, oy, c) {
		this.addDamage(shot.damage())
		return true
	}
	return false
}

func (this *Turret) addDamage(n int) {
	this.shield -= n
	if this.shield <= 0 {
		this.destroyed()
	}
	this.damaged = true
	this.damagedCnt = 10
}

func (this *Turret) destroyed() {
	playSe("turret_destroyed.wav")
	this.destroyedCnt = 0
	for i := 0; i < 6; i++ {
		s := this.smokes.GetInstanceForced()
		s.setVector(this.pos, nextSignedFloat(turretRand, 0.1), nextSignedFloat(turretRand, 0.1), nextFloat(turretRand, 0.04),
			EXPLOSION, 30+turretRand.Intn(20), this.spec.size()*1.5)
	}
	for i := 0; i < 32; i++ {
		sp := this.sparks.GetInstanceForced()
		sp.set(this.pos, nextSignedFloat(turretRand, 0.5), nextSignedFloat(turretRand, 0.5),
			0.5+nextFloat(turretRand, 0.5), 0.5+nextFloat(turretRand, 0.5), 0, 30+turretRand.Intn(30))
	}
	for i := 0; i < 7; i++ {
		f := this.fragments.GetInstanceForced()
		f.set(this.pos, nextSignedFloat(turretRand, 0.25), nextSignedFloat(turretRand, 0.25), 0.05+nextFloat(turretRand, 0.05),
			this.spec.size()*(0.5+nextFloat(turretRand, 0.5)))
	}
	switch this.spec.typ {
	case TurretMAIN:
		this.parent.increaseMultiplier(2)
		this.parent.addScore(40)
	case TurretSUB, TurretSUB_DESTRUCTIVE:
		this.parent.increaseMultiplier(1)
		this.parent.addScore(20)
	}
}

func (this *Turret) remove() {
	if this.destroyedCnt < 0 {
		this.destroyedCnt = 999
	}
}

/**
 * Turret specification changing according to a rank(difficulty).
 */

type TurretType int

const (
	TurretMAIN TurretType = iota
	TurretSUB
	TurretSUB_DESTRUCTIVE
	TurretSMALL
	TurretMOVING
	TurretDUMMY
)

type TurretSpec struct {
	typ                                                         TurretType
	interval                                                    int
	speed, speedAccel, minRange, maxRange, turnSpeed, turnRange float64
	burstNum, burstInterval                                     int
	burstTurnRatio                                              float64
	blind                                                       bool
	lookAheadRatio                                              float64
	nway                                                        int
	nwayAngle                                                   float64
	nwayChange                                                  bool
	bulletShape                                                 BulletShapeType
	bulletDestructive                                           bool
	shield                                                      int
	invisible                                                   bool
	shape, damagedShape, destroyedShape                         *TurretShape
	_size                                                       float64
}

func NewTurretSpec() TurretSpec {
	this := TurretSpec{}
	this.shape = NewTurretShape(TurretNORMAL)
	this.damagedShape = NewTurretShape(TurretDAMAGED)
	this.destroyedShape = NewTurretShape(TurretDESTROYED)
	this.init()
	return this
}

func (this *TurretSpec) init() {
	this.typ = 0
	this.interval = 99999
	this.speed = 1
	this.speedAccel = 0
	this.minRange = 0
	this.maxRange = 99999
	this.turnSpeed = 99999
	this.turnRange = 99999
	this.burstNum = 1
	this.burstInterval = 99999
	this.burstTurnRatio = 0
	this.blind = false
	this.lookAheadRatio = 0
	this.nway = 1
	this.nwayAngle = 0
	this.nwayChange = false
	this.bulletShape = BulletNORMAL
	this.bulletDestructive = false
	this.shield = 99999
	this.invisible = false
	this._size = 1
}

func (this *TurretSpec) setParamTurretSpec(ts TurretSpec) {
	this.typ = ts.typ
	this.interval = ts.interval
	this.speed = ts.speed
	this.speedAccel = ts.speedAccel
	this.minRange = ts.minRange
	this.maxRange = ts.maxRange
	this.turnSpeed = ts.turnSpeed
	this.turnRange = ts.turnRange
	this.burstNum = ts.burstNum
	this.burstInterval = ts.burstInterval
	this.burstTurnRatio = ts.burstTurnRatio
	this.blind = ts.blind
	this.lookAheadRatio = ts.lookAheadRatio
	this.nway = ts.nway
	this.nwayAngle = ts.nwayAngle
	this.nwayChange = ts.nwayChange
	this.bulletShape = ts.bulletShape
	this.bulletDestructive = ts.bulletDestructive
	this.shield = ts.shield
	this.invisible = ts.invisible
	this.setSize(ts.size())
}

func (this *TurretSpec) setParam(rank float64, typ TurretType, rand *r.Rand) {
	this.init()
	this.typ = typ
	if this.typ == TurretDUMMY {
		this.invisible = true
		return
	}
	rk := rank
	switch this.typ {
	case TurretSMALL:
		this.minRange = 8
		this.bulletShape = BulletSMALL
		this.blind = true
		this.invisible = true
	case TurretMOVING:
		this.minRange = 6
		this.bulletShape = BulletMOVING_TURRET
		this.blind = true
		this.invisible = true
		this.turnSpeed = 0
		this.maxRange = 9 + nextFloat(rand, 12)
		rk *= (10.0 / math.Sqrt(this.maxRange))
	default:
		this.maxRange = 9 + nextFloat(rand, 16)
		this.minRange = this.maxRange / (4 + nextFloat(rand, 0.5))
		if this.typ == TurretSUB || this.typ == TurretSUB_DESTRUCTIVE {
			this.maxRange *= 0.72
			this.minRange *= 0.9
		}
		rk *= (10.0 / math.Sqrt(this.maxRange))
		if rand.Intn(4) == 0 {
			lar := rank * 0.1
			if lar > 1 {
				lar = 1
			}
			this.lookAheadRatio = nextFloat(rand, lar/2) + lar/2
			rk /= (1 + this.lookAheadRatio*0.3)
		}
		if rand.Intn(3) == 0 && this.lookAheadRatio == 0 {
			this.blind = false
			rk *= 1.5
		} else {
			this.blind = true
		}
		this.turnRange = math.Pi/4 + nextFloat(rand, math.Pi/4)
		this.turnSpeed = 0.005 + nextFloat(rand, 0.015)
		if this.typ == TurretMAIN {
			this.turnRange *= 1.2
		}
		if rand.Intn(4) == 0 {
			this.burstTurnRatio = nextFloat(rand, 0.66) + 0.33
		}
	}
	this.burstInterval = 6 + rand.Intn(8)
	switch this.typ {
	case TurretMAIN:
		this.setSize(0.42 + nextFloat(rand, 0.05))
		br := (rk * 0.3) * (1 + nextSignedFloat(rand, 0.2))
		nr := (rk * 0.33) * nextFloat(rand, 1)
		ir := (rk * 0.1) * (1 + nextSignedFloat(rand, 0.2))
		this.burstNum = int(br) + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - float64(this.burstNum) + 1 - float64(this.nway-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = math.Sqrt(sr * 0.6)
		this.speed *= 0.12
		this.shield = 20
	case TurretSUB:
		this.setSize(0.36 + nextFloat(rand, 0.025))
		br := (rk * 0.4) * (1 + nextSignedFloat(rand, 0.2))
		nr := (rk * 0.2) * nextFloat(rand, 1)
		ir := (rk * 0.2) * (1 + nextSignedFloat(rand, 0.2))
		this.burstNum = int(br) + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - float64(this.burstNum) + 1 - float64(this.nway-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = math.Sqrt(sr * 0.7)
		this.speed *= 0.2
		this.shield = 12
	case TurretSUB_DESTRUCTIVE:
		this.setSize(0.36 + nextFloat(rand, 0.025))
		br := (rk * 0.4) * (1 + nextSignedFloat(rand, 0.2))
		nr := (rk * 0.2) * nextFloat(rand, 1)
		ir := (rk * 0.2) * (1 + nextSignedFloat(rand, 0.2))
		this.burstNum = int(br)*2 + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(60.0/(ir*2+1)) + 1
		this.burstInterval = int(float64(this.burstInterval) * 0.88)
		this.bulletShape = BulletDESTRUCTIVE
		this.bulletDestructive = true
		sr := rk - float64(this.burstNum-1)/2 - float64(this.nway-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = math.Sqrt(sr * 0.7)
		this.speed *= 0.33
		this.shield = 12
	case TurretSMALL:
		this.setSize(0.33)
		br := (rk * 0.33) * (1 + nextSignedFloat(rand, 0.2))
		ir := (rk * 0.2) * (1 + nextSignedFloat(rand, 0.2))
		this.burstNum = int(br) + 1
		this.nway = 1
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - float64(this.burstNum) + 1 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = math.Sqrt(sr)
		this.speed *= 0.24
	case TurretMOVING:
		this.setSize(0.36)
		br := (rk * 0.3) * (1 + nextSignedFloat(rand, 0.2))
		nr := (rk * 0.1) * nextFloat(rand, 1)
		ir := (rk * 0.33) * (1 + nextSignedFloat(rand, 0.2))
		this.burstNum = int(br) + 1
		this.nway = int(nr*0.66 + 1)
		this.interval = int(120.0/(ir*2+1)) + 1
		sr := rk - float64(this.burstNum) + 1 - float64(this.nway-1)/0.66 - ir
		if sr < 0 {
			sr = 0
		}
		this.speed = math.Sqrt(sr * 0.7)
		this.speed *= 0.2
	}
	if this.speed < 0.1 {
		this.speed = 0.1
	} else {
		this.speed = math.Sqrt(this.speed*10) / 10
	}
	if this.burstNum > 2 {
		if rand.Intn(4) == 0 {
			this.speed *= 0.8
			this.burstInterval = int(float64(this.burstInterval) * 0.7)
			this.speedAccel = (this.speed * (0.4 + nextFloat(rand, 0.3))) / float64(this.burstNum)
			if rand.Intn(2) == 0 {
				this.speedAccel *= -1
			}
			this.speed -= this.speedAccel * float64(this.burstNum) / 2
		}
		if rand.Intn(5) == 0 {
			if this.nway > 1 {
				this.nwayChange = true
			}
		}
	}
	this.nwayAngle = (0.1 + nextFloat(rand, 0.33)) / (1 + float64(this.nway)*0.1)
}

func (this *TurretSpec) setBossSpec() {
	this.minRange = 0
	this.maxRange *= 1.5
	this.shield = int(float64(this.shield) * 2.1)
}

func (this *TurretSpec) size() float64 {
	return this._size
}

func (this *TurretSpec) setSize(v float64) float64 {
	this._size = v
	this.shape.SetSize(v)
	this.damagedShape.SetSize(v)
	this.destroyedShape.SetSize(v)
	return this._size
}

/**
 * Grouped turrets.
 */

const TURRET_GROUP_MAX_NUM = 16

type TurretGroup struct {
	ship      *Ship
	sparks    *SparkPool
	smokes    *SmokePool
	fragments *FragmentPool
	spec      TurretGroupSpec
	centerPos vector.Vector
	turret    [TURRET_GROUP_MAX_NUM]Turret
	cnt       int
}

func NewTurretGroup(field *Field, bullets *BulletPool, ship *Ship,
	sparks *SparkPool, smokes *SmokePool, fragments *FragmentPool,
	parent *Enemy) TurretGroup {
	this := TurretGroup{}
	this.ship = ship
	for i := range this.turret {
		this.turret[i] = NewTurret(field, bullets, ship, sparks, smokes, fragments, parent)
	}
	return this
}

func (this *TurretGroup) set(spec TurretGroupSpec) {
	this.spec = spec
	for i := 0; i < this.spec.num; i++ {
		this.turret[i].start(this.spec.turretSpec)
	}
	this.cnt = 0
}

func (this *TurretGroup) move(p vector.Vector, deg float64) bool {
	alive := false
	this.centerPos.X = p.X
	this.centerPos.Y = p.Y
	var d, md, y, my float64
	switch this.spec.alignType {
	case ROUND:
		d = this.spec.alignDeg
		if this.spec.num > 1 {
			md = this.spec.alignWidth / float64(this.spec.num-1)
			d -= this.spec.alignWidth / 2
		} else {
			md = 0
		}
	case STRAIGHT:
		y = 0
		my = this.spec.offset.Y / float64(this.spec.num+1)
	}
	for i := 0; i < this.spec.num; i++ {
		var tbx, tby float64
		switch this.spec.alignType {
		case ROUND:
			tbx = math.Sin(d) * this.spec.radius
			tby = math.Cos(d) * this.spec.radius
		case STRAIGHT:
			y += my
			tbx = this.spec.offset.X
			tby = y
			d = math.Atan2(tbx, tby)
		}
		tbx *= (1 - this.spec.distRatio)
		bx := tbx*math.Cos(-deg) - tby*math.Sin(-deg)
		by := tbx*math.Sin(-deg) + tby*math.Cos(-deg)
		alive = alive || this.turret[i].move(this.centerPos.X+bx, this.centerPos.Y+by, d+deg, 0, -99999)
		if this.spec.alignType == ROUND {
			d += md
		}
	}
	this.cnt++
	return alive
}

func (this *TurretGroup) draw() {
	for i := 0; i < this.spec.num; i++ {
		this.turret[i].draw()
	}
}

func (this *TurretGroup) remove() {
	for i := 0; i < this.spec.num; i++ {
		this.turret[i].remove()
	}
}

func (this *TurretGroup) checkCollision(x, y float64, c sdl.Collidable, shot *Shot) bool {
	col := false
	for i := 0; i < this.spec.num; i++ {
		col = col || this.turret[i].checkCollision(x, y, c, shot)
	}
	return col
}

type AlignType int

const (
	ROUND AlignType = iota
	STRAIGHT
)

type TurretGroupSpec struct {
	turretSpec                              TurretSpec
	num                                     int
	alignType                               AlignType
	alignDeg, alignWidth, radius, distRatio float64
	offset                                  vector.Vector
}

func NewTurretGroupSpec() TurretGroupSpec {
	this := TurretGroupSpec{}
	this.turretSpec = NewTurretSpec()
	this.num = 1
	return this
}

func (this *TurretGroupSpec) init() {
	this.num = 1
	this.alignType = ROUND
	this.alignDeg = 0
	this.alignWidth = 0
	this.radius = 0
	this.distRatio = 0
	this.offset.X = 0
	this.offset.Y = 0
}

/**
 * Turrets moving around a bridge.
 */

const MOVING_TURRET_MAX_NUM = 16

type MovingTurretGroup struct {
	ship                                                                                                           *Ship
	spec                                                                                                           MovingTurretGroupSpec
	radius, radiusAmpCnt, deg, rollAmpCnt, swingAmpCnt, swingAmpDeg, swingFixDeg, alignAmpCnt, distDeg, distAmpCnt float64
	cnt                                                                                                            int
	centerPos                                                                                                      vector.Vector
	turret                                                                                                         [MOVING_TURRET_MAX_NUM]Turret
}

func NewMovingTurretGroup(field *Field, bullets *BulletPool, ship *Ship,
	sparks *SparkPool, smokes *SmokePool, fragments *FragmentPool,
	parent *Enemy) MovingTurretGroup {
	this := MovingTurretGroup{}
	this.ship = ship
	for i := range this.turret {
		this.turret[i] = NewTurret(field, bullets, ship, sparks, smokes, fragments, parent)
	}
	return this
}

func (this *MovingTurretGroup) set(spec MovingTurretGroupSpec) {
	this.spec = spec
	this.radius = spec.radiusBase
	this.radiusAmpCnt = 0
	this.deg = 0
	this.rollAmpCnt = 0
	this.swingAmpCnt = 0
	this.swingAmpDeg = 0
	this.alignAmpCnt = 0
	this.distDeg = 0
	this.distAmpCnt = 0
	this.swingFixDeg = math.Pi
	for i := 0; i < spec.num; i++ {
		this.turret[i].start(spec.turretSpec)
	}
	this.cnt = 0
}

func (this *MovingTurretGroup) move(p vector.Vector, ed float64) {
	if this.spec.moveType == SWING_FIX {
		this.swingFixDeg = ed
	}
	this.centerPos.X = p.X
	this.centerPos.Y = p.Y
	if this.spec.radiusAmp > 0 {
		this.radiusAmpCnt += this.spec.radiusAmpVel
		av := math.Sin(this.radiusAmpCnt)
		this.radius = this.spec.radiusBase + this.spec.radiusAmp*av
	}
	if this.spec.moveType == ROLL {
		if this.spec.rollAmp != 0 {
			this.rollAmpCnt += this.spec.rollAmpVel
			av := math.Sin(this.rollAmpCnt)
			this.deg += this.spec.rollDegVel + this.spec.rollAmp*av
		} else {
			this.deg += this.spec.rollDegVel
		}
	} else {
		this.swingAmpCnt += this.spec.swingAmpVel
		if math.Cos(this.swingAmpCnt) > 0 {
			this.swingAmpDeg += this.spec.swingDegVel
		} else {
			this.swingAmpDeg -= this.spec.swingDegVel
		}
		if this.spec.moveType == SWING_AIM {
			var od float64
			shipPos := this.ship.nearPos(this.centerPos)
			if shipPos.DistVector(this.centerPos) < 0.1 {
				od = 0
			} else {
				od = math.Atan2(shipPos.X-this.centerPos.X, shipPos.Y-this.centerPos.Y)
			}
			od += this.swingAmpDeg - this.deg
			od = normalizeDeg(od)
			this.deg += od * 0.1
		} else {
			od := this.swingFixDeg + this.swingAmpDeg - this.deg
			od = normalizeDeg(od)
			this.deg += od * 0.1
		}
	}
	var d, ad, md float64
	d, ad, md = this.calcAlignDeg(d, ad, md)
	for i := 0; i < this.spec.num; i++ {
		d += md
		bx := math.Sin(d) * this.radius * this.spec.xReverse
		by := math.Cos(d) * this.radius * (1 - this.spec.distRatio)
		var fs, fd float64
		if math.Abs(bx)+math.Abs(by) < 0.1 {
			fs = this.radius
			fd = d
		} else {
			fs = math.Sqrt(bx*bx + by*by)
			fd = math.Atan2(bx, by)
		}
		fs *= 0.06
		this.turret[i].move(this.centerPos.X, this.centerPos.Y, d, fs, fd)
	}
	this.cnt++
}

func (this MovingTurretGroup) calcAlignDeg(d, ad, md float64) (float64, float64, float64) {
	this.alignAmpCnt += this.spec.alignAmpVel
	ad = this.spec.alignDeg * (1 + math.Sin(this.alignAmpCnt)*this.spec.alignAmp)
	if this.spec.num > 1 {
		if this.spec.moveType == ROLL {
			md = ad / float64(this.spec.num)
		} else {
			md = ad / float64(this.spec.num-1)
		}
	} else {
		md = 0
	}
	d = this.deg - md - ad/2
	return d, ad, md
}

func (this MovingTurretGroup) draw() {
	for i := 0; i < this.spec.num; i++ {
		this.turret[i].draw()
	}
}

func (this MovingTurretGroup) remove() {
	for i := 0; i < this.spec.num; i++ {
		this.turret[i].remove()
	}
}

type BulletMoveType int

const (
	ROLL BulletMoveType = iota
	SWING_FIX
	SWING_AIM
)

type MovingTurretGroupSpec struct {
	turretSpec TurretSpec
	num        int
	moveType   BulletMoveType
	alignDeg, alignAmp, alignAmpVel, radiusBase,
	radiusAmp, radiusAmpVel, rollDegVel, rollAmp,
	rollAmpVel, swingDegVel, swingAmpVel, distRatio, xReverse float64
}

func NewMovingTurretGroupSpec() MovingTurretGroupSpec {
	this := MovingTurretGroupSpec{}
	this.turretSpec = NewTurretSpec()
	this.num = 1
	this.initParam()
	return this
}

func (this *MovingTurretGroupSpec) initParam() {
	this.num = 1
	this.alignDeg = math.Pi * 2
	this.alignAmp = 0
	this.alignAmpVel = 0
	this.radiusBase = 1
	this.radiusAmp = 0
	this.radiusAmpVel = 0
	this.moveType = SWING_FIX
	this.rollDegVel = 0
	this.rollAmp = 0
	this.rollAmpVel = 0
	this.swingDegVel = 0
	this.swingAmpVel = 0
	this.distRatio = 0
	this.xReverse = 1
}

func (this *MovingTurretGroupSpec) init() {
	this.initParam()
}

func (this *MovingTurretGroupSpec) setAlignAmp(a, v float64) {
	this.alignAmp = a
	this.alignAmpVel = v
}

func (this *MovingTurretGroupSpec) setRadiusAmp(a, v float64) {
	this.radiusAmp = a
	this.radiusAmpVel = v
}

func (this *MovingTurretGroupSpec) setRoll(dv, a, v float64) {
	this.moveType = ROLL
	this.rollDegVel = dv
	this.rollAmp = a
	this.rollAmpVel = v
}

func (this *MovingTurretGroupSpec) setSwing(dv, a float64, aim bool /*= false */) {
	if aim {
		this.moveType = SWING_AIM
	} else {
		this.moveType = SWING_FIX
	}
	this.swingDegVel = dv
	this.swingAmpVel = a
}

func (this *MovingTurretGroupSpec) setXReverse(xr float64) {
	this.xReverse = xr
}
