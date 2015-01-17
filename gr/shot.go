/*
 * $Id: shot.d,v 1.2 2005/07/03 07:05:22 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/go-gl/gl"
)

/**
 * Player's shot.
 */

const SHOT_SPEED = 0.6
const LANCE_SPEED = 0.5

var shotShape *ShotShape
var lanceShape *LanceShape

type Shot struct {
	pos    Vector
	cnt    int
	hitCnt int
	deg    float32
	damage int
	lance  bool
}

func initShots() {
	shotShape = NewShotShape()
	lanceShape = NewLanceShape()
}

func closeShots() {
	shotShape.close()
}

func NewShot(p Vector, d float32, lance bool /*= false*/, dmg int /*= -1*/) *Shot {
	s := new(Shot)
	s.damage = 1
	s.pos.x = p.x
	s.pos.y = p.y
	s.deg = d
	s.lance = lance
	if lance {
		s.damage = 10
	} else {
		s.damage = 1
	}
	if dmg >= 0 {
		s.damage = dmg
	}
	actors[s] = true
	return s
}

func (s *Shot) move() {
	s.cnt++
	if s.hitCnt > 0 {
		s.hitCnt++
		if s.hitCnt > 30 {
			s.close()
		}
		return
	}
	var sp float32
	if !s.lance {
		sp = SHOT_SPEED
	} else {
		if s.cnt < 10 {
			sp = LANCE_SPEED * float32(s.cnt) / 10
		} else {
			sp = LANCE_SPEED
		}
	}
	s.pos.x += Sin32(s.deg) * sp
	s.pos.y += Cos32(s.deg) * sp
	s.pos.y -= field.lastScrollY
	if field.getBlockVector(s.pos) >= ON_BLOCK_THRESHOLD ||
		!field.checkInOuterFieldVector(s.pos) || s.pos.y > field.size.y {
		s.close()
	}
	if s.lance {
		checkAllEnemiesShotHit(s.pos, lanceShape, s)
	} else {
		checkAllBulletsShotHit(s.pos, shotShape, s)
		checkAllEnemiesShotHit(s.pos, shotShape, s)
	}
}

func (s *Shot) close() {
	delete(actors, s)
	if s.lance && s.hitCnt <= 0 {
		s.hitCnt = 1
		return
	}
}

func (s *Shot) removeHitToBullet() {
	s.removeHit()
}

func (s *Shot) removeHitToEnemy(isSmallEnemy bool /*= false*/) {
	if isSmallEnemy && s.lance {
		return
	}
	playSe("hit.wav")
	s.removeHit()
}

func (this *Shot) removeHit() {
	this.close()
	if this.lance {
		for i := 0; i < 10; i++ {
			d := this.deg + nextSignedFloat(0.1)
			sp := nextSignedFloat(LANCE_SPEED)
			NewSmoke(this.pos.x, this.pos.y, 0, Sin32(d)*sp, Cos32(d)*sp, 0, SmokeTypeLANCE_SPARK, 30+int(30), 1)

			d = this.deg + nextSignedFloat(0.1)
			sp = nextFloat(LANCE_SPEED)
			NewSmoke(this.pos.x, this.pos.y, 0, -Sin32(d)*sp, -Cos32(d)*sp, 0, SmokeTypeLANCE_SPARK, 30+int(30), 1)
		}
	} else {
		d := this.deg + nextSignedFloat(0.5)
		NewSpark(this.pos, Sin32(d)*SHOT_SPEED, Cos32(d)*SHOT_SPEED, 0.6+nextSignedFloat(0.4), 0.6+nextSignedFloat(0.4), 0.1, 20)

		d = this.deg + nextSignedFloat(0.5)
		NewSpark(this.pos, -Sin32(d)*SHOT_SPEED, -Cos32(d)*SHOT_SPEED, 0.6+nextSignedFloat(0.4), 0.6+nextSignedFloat(0.4), 0.1, 20)
	}
}

func (this *Shot) draw() {
	if this.lance {
		x := this.pos.x
		y := this.pos.y
		var size float32 = 0.25
		var a float32 = 0.6
		hc := this.hitCnt
		for i := 0; i < this.cnt/4+1; i++ {
			size *= 0.9
			a *= 0.8
			if hc > 0 {
				hc--
				continue
			}
			d := i*13 + this.cnt*3
			for j := 0; j < 6; j++ {
				gl.PushMatrix()
				gl.Translatef(x, y, 0)
				gl.Rotatef(-this.deg*180/Pi32, 0, 0, 1)
				gl.Rotatef(float32(d), 0, 1, 0)
				setScreenColor(0.4, 0.8, 0.8, a)
				gl.Begin(gl.LINE_LOOP)
				gl.Vertex3f(-size, LANCE_SPEED, size/2)
				gl.Vertex3f(size, LANCE_SPEED, size/2)
				gl.Vertex3f(size, -LANCE_SPEED, size/2)
				gl.Vertex3f(-size, -LANCE_SPEED, size/2)
				gl.End()
				setScreenColor(0.2, 0.5, 0.5, a/2)
				gl.Begin(gl.TRIANGLE_FAN)
				gl.Vertex3f(-size, LANCE_SPEED, size/2)
				gl.Vertex3f(size, LANCE_SPEED, size/2)
				gl.Vertex3f(size, -LANCE_SPEED, size/2)
				gl.Vertex3f(-size, -LANCE_SPEED, size/2)
				gl.End()
				gl.PopMatrix()
				d += 60
			}
			x -= Sin32(this.deg) * LANCE_SPEED * 2
			y -= Cos32(this.deg) * LANCE_SPEED * 2
		}
	} else {
		gl.PushMatrix()
		glTranslate(this.pos)
		gl.Rotatef(-this.deg*180/Pi32, 0, 0, 1)
		gl.Rotatef(float32(this.cnt)*31, 0, 1, 0)
		shotShape.draw()
		gl.PopMatrix()
	}
}

func (this *Shot) removed() bool {
	return this.hitCnt > 0
}

func existsLance() bool {
	for a := range actors {
		s, ok := a.(*Shot)
		if ok && s.lance && !s.removed() {
			return true
		}
	}
	return false
}

type ShotShape struct {
	*SimpleShape
}

func NewShotShape() *ShotShape {
	ss := new(ShotShape)
	ss.startDisplayList()
	setScreenColor(0.1, 0.33, 0.1, 1)
	gl.Begin(gl.QUADS)
	gl.Vertex3f(0, 0.3, 0.1)
	gl.Vertex3f(0.066, 0.3, -0.033)
	gl.Vertex3f(0.1, -0.3, -0.05)
	gl.Vertex3f(0, -0.3, 0.15)
	gl.Vertex3f(0.066, 0.3, -0.033)
	gl.Vertex3f(-0.066, 0.3, -0.033)
	gl.Vertex3f(-0.1, -0.3, -0.05)
	gl.Vertex3f(0.1, -0.3, -0.05)
	gl.Vertex3f(-0.066, 0.3, -0.033)
	gl.Vertex3f(0, 0.3, 0.1)
	gl.Vertex3f(0, -0.3, 0.15)
	gl.Vertex3f(-0.1, -0.3, -0.05)
	gl.End()
	ss.endDisplayList()
	ss.collision = &Vector{0.33, 0.33}
	return ss
}

type LanceShape struct {
	*SimpleShape
}

func NewLanceShape() *LanceShape {
	ls := new(LanceShape)
	ls.startDisplayList()
	// no display for this shape.
	ls.endDisplayList()
	ls.collision = &Vector{0.66, 0.66}
	return ls
}
