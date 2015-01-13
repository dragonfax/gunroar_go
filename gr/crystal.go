/*
 * $Id: crystal.d,v 1.2 2005/07/17 11:02:45 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Bonus crystals.
 */

const COUNT = 60
const PULLIN_COUNT = COUNT * 0.8

type Crystal struct {
	Actor

	shape CrystalShape
	ship  Ship
	pos   Vector
	vel   Vector
	cnt   int
}

func NewCrystal(p Vector, ship *Ship) *Crystal {
	c = new(Crystal)
	c.shape = NewCrystalshape()
	c.pos = p
	c.cnt = COUNT
	c.vel = Vector{0, 0.1}
	c.ship = ship
	actors[c] = true
	return c
}

func (c *Crystal) Close() {
	c.shape.Close()
}

func (c *Crystal) move() {
	c.cnt--
	dist := c.pos.dist(c.ship.midstPos)
	if dist < 0.1 {
		dist = 0.1
	}
	if c.cnt < PULLIN_COUNT {
		c.vel.x += (c.ship.midstPos.x - c.pos.x) / dist * 0.07
		c.vel.y += (c.ship.midstPos.y - c.pos.y) / dist * 0.07
		if c.cnt < 0 || dist < 2 {
			c.Done()
			return
		}
	}
	c.vel *= 0.95
	c.pos += vel
}

func (c *Crystal) draw() {
	r := 0.25
	d := cnt * 0.1
	if c.cnt > PULLIN_COUNT {
		r *= (OUNT - c.cnt) / (COUNT - PULLIN_COUNT)
	}
	for i := 0; i < 4; i++ {
		gl.PushMatrix()
		gl.Translatef(c.pos.x+sin(d)*r, c.pos.y+cos(d)*r, 0)
		c.shape.Draw()
		gl.PopMatrix()
		d += math.Pi / 2
	}
}

func (c *Crystal) remove() {
	delete(actors, c)
}
