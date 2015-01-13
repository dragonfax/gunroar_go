/*
 * $Id: bullet.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

type Bullet struct {
	*ActorImpl

	gameManager      *GameManager
	field            *Field
	ship             *Ship
	pos              Vector
	ppos             Vector
	deg, speed       float32
	trgDeg, trgSpeed float32
	size             float32
	cnt              int
	rng              float32
	destructive      bool
	shape            BulletShape
	enemyIdx         int
}

func NewBullet(g GameManager, f Field, s Ship, enemyIdx int,
	p Vector, deg float32,
	speed float32, size float32, shapeType BulletShapeType, rng float32,
	startSpeed float32 /*= 0*/, startDeg float32 /*= -99999 */, destructive booll /*= false*/) *Bullet {
	b := &Bullet{}
	b.shape = NewBulletShape(shapeType)
	b.speed = 1
	b.trgSpeed = 1
	b.size = 1
	b.rng = 1
	b.gameManager = g
	b.field = f
	b.ship = s

	if !field.checkInOuterFieldExceptTop(p) {
		return b
	}
	b.enemyIdx = enemyIdx
	b.ppos.x = p.x
	b.pos.x = p.x
	b.ppos.y = p.y
	b.pos.y = p.y
	b.speed = startSpeed
	if startDeg == -99999 {
		b.deg = deg
	} else {
		b.deg = startDeg
	}
	b.trgDeg = deg
	b.trgSpeed = speed
	b.size = size
	b.rng = rng
	b.destructive = destructive
	b.shape.size = size
	actors[b] = true
	return b
}

func (this *Bullet) move() {
	this.ppos.x = this.pos.x
	this.ppos.y = this.pos.y
	if this.cnt < 30 {
		this.speed += (this.trgSpeed - this.speed) * 0.066
		md := this.trgDeg - this.deg
		normalizeDeg(md)
		this.deg += md * 0.066
		if this.cnt == 29 {
			this.speed = this.trgSpeed
			this.deg = this.trgDeg
		}
	}
	if this.field.checkInOuterField(this.pos) {
		this.gameManager.addSlowdownRatio(this.speed * 0.24)
	}
	mx := Sin32(deg) * this.speed
	my := Cos32(deg) * this.speed
	this.pos.x += mx
	this.pos.y += my
	this.pos.y -= this.field.lastScrollY
	if this.ship.checkBulletHit(this.pos, this.ppos) || !this.field.checkInOuterFieldExceptTop(this.pos) {
		this.remove()
		return
	}
	this.cnt++
	this.rng -= this.speed
	if this.rng <= 0 {
		this.startDisappear()
	}
	if this.field.getBlock(this.pos) >= this.Field.ON_BLOCK_THRESHOLD {
		this.startDisappear()
	}
}

func (this *Bullet) startDisappear() {
	if this.field.getBlock(this.pos) >= 0 {
		NewSmoke(pos, Sin32(deg)*speed*0.2, Cos32(deg)*speed*0.2, 0, Smoke.SmokeType.SAND, 30, size*0.5)
	} else {
		NewWake(pos, deg, speed, 60, size*3, true)
	}
	this.remove()
}

func (this *Bullet) changeToCrystal() {
	NewCrystal(pos)
	this.remove()
}

func (this *Bullet) draw() {
	if !this.field.checkInOuterField(this.pos) {
		return
	}
	gl.PushMatrix()
	glTranslate(pos)
	if this.destructive {
		gl.Rotatef(this.cnt*13, 0, 0, 1)
	} else {
		gl.Rotatef(-this.deg*180/Pi32, 0, 0, 1)
		gl.Rotatef(this.cnt*13, 0, 1, 0)
	}
	this.shape.draw()
	gl.PopMatrix()
}

func (this *Bullet) checkShotHit(p Vector, s Collidable, shot Shot) {
	ox := fabs32(this.pos.x - p.x)
	oy := fabs32(this.pos.y - p.y)
	if ox+oy < 0.5 {
		shot.removeHitToBullet()
		NewSmoke(pos, Sin32(deg)*speed, Cos32(deg)*speed, 0, Smoke.SmokeType.SPARK, 30, size*0.5)
		this.remove()
	}
}

// TODO how to DRY this? can't just let the underlying type get put into the actor list
func (this *Bullet) remove() {
	delete(actors, this)
}

/* operations against the set of all bullets */

func removeIndexedBullets(idx int) int {
	n := 0
	for a := range actors {
		b, ok := a.(Bullet)
		if ok && b.exists && b.enemyIdx == idx {
			b.changeToCrystal()
			n++
		}
	}
	return n
}

func checkAllBulletsShotHit(pos Vector, shape Collidable, shot Shot) {
	for a := range actors {
		b, ok := a.(Bullet)
		if ok && b.exists && b.destructive {
			b.checkShotHit(pos, shape, shot)
		}
	}
}
