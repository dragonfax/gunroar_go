/*
 * $Id: crystal.d,v 1.2 2005/07/17 11:02:45 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/go-gl/gl/v3.3-compatibility/gl"
)

/**
 * Bonus crystals.
 */

const COUNT = 60
const PULLIN_COUNT = 48 // floor(COUNT * 0.8)

type Crystal struct {
	shape *CrystalShape
	pos   Vector
	vel   Vector
	cnt   int
}

func NewCrystal(p Vector) *Crystal {
	c := new(Crystal)
	c.shape = crystalShape
	c.pos = p
	c.cnt = COUNT
	c.vel = Vector{0, 0.1}
	actors[c] = true
	return c
}

func (c *Crystal) close() {
	c.shape.close()
	delete(actors, c)
}

func (c *Crystal) move() {
	c.cnt--
	dist := c.pos.distVector(ship.midstPos())
	if dist < 0.1 {
		dist = 0.1
	}
	if c.cnt < PULLIN_COUNT {
		midstPos := ship.midstPos()
		c.vel.x += (midstPos.x - c.pos.x) / dist * 0.07
		c.vel.y += (midstPos.y - c.pos.y) / dist * 0.07
		if c.cnt < 0 || dist < 2 {
			c.close()
			return
		}
	}
	c.vel.MulAssign(0.95)
	c.pos.AddAssign(c.vel)
}

func (c *Crystal) draw() {
	var r float32 = 0.25
	d := float32(c.cnt) * 0.1
	if c.cnt > PULLIN_COUNT {
		r *= (COUNT - float32(c.cnt)) / (COUNT - PULLIN_COUNT)
	}
	for i := 0; i < 4; i++ {
		gl.PushMatrix()
		gl.Translatef(c.pos.x+Sin32(d)*r, c.pos.y+Cos32(d)*r, 0)
		c.shape.draw()
		gl.PopMatrix()
		d += Pi32 / 2
	}
}
