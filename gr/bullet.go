package main

import (
	"math"

	"github.com/dragonfax/gunroar/gr/actor"
	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Enemy's bullets.
 */

var _ actor.Actor = &Bullet{}

type Bullet struct {
	actor.ExistsImpl

	gameManager                        *GameManager
	field                              *Field
	ship                               *Ship
	smokes                             *SmokePool
	wakes                              *WakePool
	crystals                           *CrystalPool
	pos                                vector.Vector
	ppos                               vector.Vector
	deg, speed, trgDeg, trgSpeed, size float64
	cnt                                int
	rang                               float64
	_destructive                       bool
	shape                              *BulletShape
	_enemyIdx                          int
}

func NewBullet() *Bullet {
	this := &Bullet{}
	this.shape = NewBulletShape()
	this.speed = 1
	this.trgSpeed = 1
	this.size = 1
	this.rang = 1
	return this
}

func (this *Bullet) Init(args []interface{}) {
	this.gameManager = args[0].(*GameManager)
	this.field = args[1].(*Field)
	this.ship = args[2].(*Ship)
	this.smokes = args[3].(*SmokePool)
	this.wakes = args[4].(*WakePool)
	this.crystals = args[5].(*CrystalPool)
}

func (this *Bullet) set(enemyIdx int,
	p vector.Vector, deg, speed, size float64, shapeType int,
	rang, startSpeed /* = 0 */, startDeg float64, /* = -99999 */
	destructive bool /* = false */) {
	if !this.field.checkInOuterFieldExceptTop(p) {
		return
	}
	this._enemyIdx = enemyIdx
	this.ppos.X = p.X
	this.pos.X = p.X
	this.ppos.Y = p.Y
	this.pos.Y = p.Y
	this.speed = startSpeed
	if startDeg == -99999 {
		this.deg = deg
	} else {
		this.deg = startDeg
	}
	this.trgDeg = deg
	this.trgSpeed = speed
	this.size = size
	this.rang = rang
	this._destructive = destructive
	this.shape.Set(shapeType)
	this.shape.SetSize(size)
	this.cnt = 0
	this.SetExists(true)
}

func (this *Bullet) Move() {
	this.ppos.X = this.pos.X
	this.ppos.Y = this.pos.Y
	if this.cnt < 30 {
		this.speed += (this.trgSpeed - this.speed) * 0.066
		md := this.trgDeg - this.deg
		md = normalizeDeg(md)
		this.deg += md * 0.066
		if this.cnt == 29 {
			this.speed = this.trgSpeed
			this.deg = this.trgDeg
		}
	}
	if this.field.checkInOuterFieldVector(this.pos) {
		this.gameManager.addSlowdownRatio(this.speed * 0.24)
	}
	mx := math.Sin(this.deg) * this.speed
	my := math.Cos(this.deg) * this.speed
	this.pos.X += mx
	this.pos.Y += my
	this.pos.Y -= this.field.lastScrollY()
	if this.ship.checkBulletHit(this.pos, this.ppos) || !this.field.checkInOuterFieldExceptTop(this.pos) {
		this.remove()
		return
	}
	this.cnt++
	this.rang -= this.speed
	if this.rang <= 0 {
		this.startDisappear()
	}
	if this.field.getBlockVector(this.pos) >= ON_BLOCK_THRESHOLD {
		this.startDisappear()
	}
}

func (this *Bullet) startDisappear() {
	if this.field.getBlockVector(this.pos) >= 0 {
		s := this.smokes.GetInstanceForced()
		s.setVector(this.pos, math.Sin(this.deg)*this.speed*0.2, math.Cos(this.deg)*this.speed*0.2, 0,
			SAND, 30, this.size*0.5)
	} else {
		w := this.wakes.GetInstanceForced()
		w.set(this.pos, this.deg, this.speed, 60, this.size*3, true)
	}
	this.remove()
}

func (this *Bullet) changeToCrystal() {
	c := this.crystals.GetInstance()
	if c != nil {
		c.set(this.pos)
	}
	this.remove()
}

func (this *Bullet) remove() {
	this.SetExists(false)
}

func (this *Bullet) Draw() {
	if !this.field.checkInOuterFieldVector(this.pos) {
		return
	}
	gl.PushMatrix()
	sdl.GlTranslate(this.pos)
	if this._destructive {
		gl.Rotated(float64(this.cnt)*13, 0, 0, 1)
	} else {
		gl.Rotated(-this.deg*180/math.Pi, 0, 0, 1)
		gl.Rotated(float64(this.cnt)*13, 0, 1, 0)
	}
	this.shape.Draw()
	gl.PopMatrix()
}

func (this *Bullet) checkShotHit(p vector.Vector, s sdl.Collidable, shot *Shot) {
	ox := math.Abs(this.pos.X - p.X)
	oy := math.Abs(this.pos.Y - p.Y)
	if ox+oy < 0.5 {
		shot.removeHitToBullet()
		s := this.smokes.GetInstance()
		if s != nil {
			s.setVector(this.pos, math.Sin(this.deg)*this.speed, math.Cos(this.deg)*this.speed, 0,
				SPARK, 30, this.size*0.5)
		}
		this.remove()
	}
}

func (this *Bullet) destructive() bool {
	return this._destructive
}

func (this *Bullet) enemyIdx() int {
	return this._enemyIdx
}

type BulletPool struct {
	actor.ActorPool
}

func NewBulletPool(n int, args []interface{}) *BulletPool {
	f := func() actor.Actor { return NewBullet() }
	this := &BulletPool{
		ActorPool: actor.NewActorPool(f, n, args),
	}
	return this
}

func (this *BulletPool) removeIndexedBullets(idx int) int {
	n := 0
	for _, a := range this.Actor {
		b := a.(*Bullet)
		if b.Exists() && b.enemyIdx() == idx {
			b.changeToCrystal()
			n++
		}
	}
	return n
}

func (this *BulletPool) checkShotHit(pos vector.Vector, shape sdl.Collidable, shot *Shot) {
	for _, a := range this.Actor {
		b := a.(*Bullet)
		if b.Exists() && b.destructive() {
			b.checkShotHit(pos, shape, shot)
		}
	}
}
