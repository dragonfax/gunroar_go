package main

import (
	"math"

	"github.com/dragonfax/gunroar/gr/actor"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Bonus crystals.
 */

const COUNT = 60
const PULLIN_COUNT = COUNT * 0.8

var _ actor.Actor = &Crystal{}

type Crystal struct {
	actor.ExistsImpl

	ship *Ship
	pos  vector.Vector
	vel  vector.Vector
	cnt  int
}

var _crystalShape *CrystalShape

func crystalInit() {
	_crystalShape = NewCrystalShape()
}

func NewCrystal() *Crystal {
	this := &Crystal{}
	return this
}

func (this *Crystal) Init(args []interface{}) {
	this.ship = args[0].(*Ship)
}

func (this *Crystal) set(p vector.Vector) {
	this.pos.X = p.X
	this.pos.Y = p.Y
	this.cnt = COUNT
	this.vel.X = 0
	this.vel.Y = 0.1
	this.SetExists(true)
}

func (this *Crystal) Move() {
	this.cnt--
	dist := this.pos.DistVector(this.ship.midstPos())
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

func (this *Crystal) Draw() {
	r := 0.25
	d := float64(this.cnt) * 0.1
	if this.cnt > PULLIN_COUNT {
		r *= (COUNT - float64(this.cnt)) / (COUNT - PULLIN_COUNT)
	}
	for i := 0; i < 4; i++ {
		gl.PushMatrix()
		gl.Translated(this.pos.X+math.Sin(d)*r, this.pos.Y+math.Cos(d)*r, 0)
		_crystalShape.Draw()
		gl.PopMatrix()
		d += math.Pi / 2
	}
}

type CrystalPool struct {
	actor.ActorPool
}

func NewCrystalPool(n int, args []interface{}) *CrystalPool {
	return &CrystalPool{actor.NewActorPool(func() actor.Actor { return NewCrystal() }, n, args)}
}
