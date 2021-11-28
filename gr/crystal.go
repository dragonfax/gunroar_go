package main

import (
	"math"

	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Bonus crystals.
 */

const COUNT = 60
const PULLIN_COUNT = COUNT * 0.8

var _ Actor = &Crystal{}

type Crystal struct {
	ExistsImpl

	ship Ship
	pos  vector.Vector
	vel  vector.Vector
	cnt  int
}

var _crystalShape CrystalShape

func crystalInit() {
	_crystalShape = NewCrystalShape()
}

func NewCrystalShape() *Crystal {
	this := &Crystal{
		ExistsImpl: NewExistsImpl(),
	}
	return this
}

func (this *Crystal) init(args []interface{}) {
	this.ship = args[0].(*Ship)
}

func (this *Crystal) set(p Vector) {
	this.pos.X = p.x
	this.pos.Y = p.y
	this.cnt = COUNT
	this.vel.X = 0
	this.vel.Y = 0.1
	this.SetExists(true)
}

func (this *Crystal) move() {
	this.cnt--
	dist := this.pos.dist(ship.midstPos)
	if dist < 0.1 {
		dist = 0.1
	}
	if this.cnt < PULLIN_COUNT {
		this.vel.X += (this.ship.midstPos.X - this.pos.X) / dist * 0.07
		this.vel.Y += (this.ship.midstPos.Y - this.pos.Y) / dist * 0.07
		if this.cnt < 0 || dist < 2 {
			this.SetExists(false)
			return
		}
	}
	this.vel.OpMulAssign(0.95)
	this.pos.OpAddAssign(this.vel)
}

func (this *Crystal) draw() {
	r := 0.25
	d := float64(this.cnt) * 0.1
	if this.cnt > PULLIN_COUNT {
		r *= (COUNT - float64(this.cnt)) / (COUNT - PULLIN_COUNT)
	}
	for i := 0; i < 4; i++ {
		gl.PushMatrix()
		gl.Translatef(this.pos.X+math.Sin(d)*r, this.pos.Y+math.Cos(d)*r, 0)
		_crystalShape.draw()
		gl.PopMatrix()
		d += math.Pi / 2
	}
}

type CrystalPool struct {
	ActorPool
}

func NewCrystalPool(n int, args []interface{}) *CrystalPool {
	return &CrystalPool{NewActorPool(func(args []interface{}) Actor { return NewCrystal(args) }, n, args)}
}
