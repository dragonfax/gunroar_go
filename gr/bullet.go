/*
 * $Id: bullet.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/go-gl/gl"

type Bullet struct {
	pos              Vector
	ppos             Vector
	deg, speed       float32
	trgDeg, trgSpeed float32
	size             float32
	cnt              int
	rng              float32
	destructive      bool
	shape            *BulletShape
	enemyIdx         int
	stopMovingC      chan bool
}

func NewBullet(enemyIdx int,
	p Vector, deg float32,
	speed float32, size float32, shapeType BulletShapeType, rng float32,
	startSpeed float32 /*= 0*/, startDeg float32 /*= -99999 */, destructive bool /*= false*/) *Bullet {
	b := new(Bullet)
	b.shape = NewBulletShape(shapeType)
	b.speed = 1
	b.trgSpeed = 1
	b.size = 1
	b.rng = 1

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
	actorsLock.Lock()
	actors[b] = true
	actorsLock.Unlock()

	b.stopMovingC = make(chan bool)
	go func() {
		limit := NewFrameLimiter()
		for {
			select {
			case <-b.stopMovingC:
				close(b.stopMovingC)
				return
			default:
				b.moveG()
			}
			limit.cycle()
		}
	}()

	return b
}

func (this *Bullet) move() {
	// to satisfy interface
}

func (this *Bullet) moveG() {
	this.ppos.x = this.pos.x
	this.ppos.y = this.pos.y
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
	mx := Sin32(this.deg) * this.speed
	my := Cos32(this.deg) * this.speed
	this.pos.x += mx
	this.pos.y += my
	this.pos.y -= field.lastScrollY
	if ship.checkBulletHit(this.pos, this.ppos) || !field.checkInOuterFieldExceptTop(this.pos) {
		this.close()
		return
	}
	this.cnt++
	this.rng -= this.speed
	if this.rng <= 0 {
		this.startDisappear()
	}
	if field.getBlockVector(this.pos) >= ON_BLOCK_THRESHOLD {
		this.startDisappear()
	}
}

func (this *Bullet) startDisappear() {
	if field.getBlockVector(this.pos) >= 0 {
		NewSmoke(this.pos.x, this.pos.y, 0, Sin32(this.deg)*this.speed*0.2, Cos32(this.deg)*this.speed*0.2, 0, SmokeTypeSAND, 30, this.size*0.5)
	} else {
		NewWake(this.pos, this.deg, this.speed, 60, this.size*3, true)
	}
	this.close()
}

func (this *Bullet) changeToCrystal() {
	NewCrystal(this.pos)
	this.close()
}

func (this *Bullet) draw() {
	if !field.checkInOuterFieldVector(this.pos) {
		return
	}
	gl.PushMatrix()
	glTranslate(this.pos)
	if this.destructive {
		gl.Rotatef(float32(this.cnt)*13, 0, 0, 1)
	} else {
		gl.Rotatef(-this.deg*180/Pi32, 0, 0, 1)
		gl.Rotatef(float32(this.cnt)*13, 0, 1, 0)
	}
	this.shape.draw()
	gl.PopMatrix()
}

func (this *Bullet) checkShotHit(p Vector, s Shape, shot *Shot) {
	ox := fabs32(this.pos.x - p.x)
	oy := fabs32(this.pos.y - p.y)
	if ox+oy < 0.5 {
		shot.removeHitToBullet()
		NewSmoke(this.pos.x, this.pos.y, 0, Sin32(this.deg)*this.speed, Cos32(this.deg)*this.speed, 0, SmokeTypeSPARK, 30, this.size*0.5)
		this.close()
	}
}

func (this *Bullet) close() {
	actorsLock.Lock()
	delete(actors, this)
	actorsLock.Unlock()
	this.stopMovingC <- true
}

/* operations against the set of all bullets */

func removeAllIndexedBullets(idx int) int {
	n := 0
	for a := range actors {
		b, ok := a.(*Bullet)
		if ok && b.enemyIdx == idx {
			b.changeToCrystal()
			n++
		}
	}
	return n
}

func checkAllBulletsShotHit(pos Vector, shape Shape, shot *Shot) {
	for a := range actors {
		b, ok := a.(*Bullet)
		if ok && b.destructive {
			b.checkShotHit(pos, shape, shot)
		}
	}
}
