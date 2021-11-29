package main

import (
	"math"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/actor"
	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Player's shot.
 */

var _ actor.Actor = &Shot{}

const SPEED = 0.6
const LANCE_SPEED = 0.5

var shotShape ShotShape
var lanceShape LanceShape
var shotRand *r.Rand

type Shot struct {
	field   *Field
	enemeis *EnemyPool
	sparks  *SparkPool
	smokes  *SmokePool
	bullets *BulletPool
	pos     vector.Vector
	cnt     int
	hitCnt  int
	_deg    float64
	_damage int
	lance   bool
}

func shotInit() {
	shotShape = NewShotShape()
	lanceShape = NewLanceShape()
	shotRand = r.New(r.NewSource(time.Now().Unix()))
}

func setShotRandSeed(seed int64) {
	shotRand = r.New(r.NewSource(seed))
}

func NewShot() *Shot {
	this := &Shot{}
	this._damage = 1
	return this
}

func (this *Shot) init(args []interface{}) {
	this.field = args[0].(Field)
	this.enemies = args[1].(*EnemyPool)
	this.sparks = args[2].(*SparkPool)
	this.smokes = args[3].(*SmokePool)
	this.bullets = args[4].(*BulletPool)
}

func (this *Shot) set(p Vector, d float64, lance bool /* = false */, dmg int /* = -1 */) {
	this.pos.X = p.X
	this.pos.Y = p.Y
	this.cnt = 0
	this.hitCnt = 0
	this._deg = d
	this.lance = lance
	if lance {
		this._damage = 10
	} else {
		thi._damage = 1
	}
	if dmg >= 0 {
		this._damage = dmg
	}
	this.SetExists(true)
}

func (this *Shot) move() {
	this.cnt++
	if this.hitCnt > 0 {
		this.hitCnt++
		if this.hitCnt > 30 {
			this.remove()
		}
		return
	}
	var sp float64
	if !lance {
		sp = SPEED
	} else {
		if this.cnt < 10 {
			sp = LANCE_SPEED * this.cnt / 10
		} else {
			sp = LANCE_SPEED
		}
	}
	this.pos.x += math.Sin(this._deg) * sp
	this.pos.y += math.cos(this._deg) * sp
	this.pos.y -= this.field.lastScrollY
	if this.field.getBlock(this.pos) >= Field.ON_BLOCK_THRESHOLD ||
		!this.field.checkInOuterField(this.pos) || this.pos.Y > this.field.size.Y {
		this.remove()
	}
	if lance {
		this.enemies.checkShotHit(this.pos, this.lanceShape, this)
	} else {
		this.bullets.checkShotHit(this.pos, this.shape, this)
		this.enemies.checkShotHit(this.pos, this.shape, this)
	}
}

func (this *Shot) remove() {
	if this.lance && this.hitCnt <= 0 {
		this.hitCnt = 1
		return
	}
	this.SetExists(false)
}

func (this *Shot) removeHitToBullet() {
	this.removeHit()
}

func (this *Shot) removeHitToEnemy(isSmallEnemy bool /* = false */) {
	if isSmallEnemy && this.lance {
		return
	}
	playSe("hit.wav")
	this.removeHit()
}

func (this *Shot) removeHit() {
	this.remove()
	var sn int
	if this.lance {
		for i := 0; i < 10; i++ {
			s := this.smokes.getInstanceForced()
			d := this._deg + rand.nextSignedFloat(0.1)
			sp := shotRand.nextFloat(LANCE_SPEED)
			s.set(this.pos, math.Sin(d)*sp, math.Cos(d)*sp, 0,
				Smoke.SmokeType.LANCE_SPARK, 30+shotRand.nextInt(30), 1)
			s = this.smokes.getInstanceForced()
			d = this._deg + shotRand.nextSignedFloat(0.1)
			sp = shotRand.nextFloat(LANCE_SPEED)
			s.set(this.pos, -math.Sin(d)*sp, -math.Cos(d)*sp, 0,
				Smoke.SmokeType.LANCE_SPARK, 30+shotRand.nextInt(30), 1)
		}
	} else {
		s := this.sparks.getInstanceForced()
		d := this._deg + shotRand.nextSignedFloat(0.5)
		s.set(this.pos, math.Sin(d)*SPEED, math.Cos(d)*SPEED,
			0.6+shotRand.nextSignedFloat(0.4), 0.6+shotRand.nextSignedFloat(0.4), 0.1, 20)
		s = this.sparks.getInstanceForced()
		d = this._deg + shotRand.nextSignedFloat(0.5)
		s.set(this.pos, -math.Sin(d)*SPEED, -math.Cos(d)*SPEED,
			0.6+shotRand.nextSignedFloat(0.4), 0.6+shotRand.nextSignedFloat(0.4), 0.1, 20)
	}
}

func (this *Shot) draw() {
	if this.lance {
		x := this.pos.X
		y := this.pos.Y
		size := 0.25
		a := 0.6
		hc := this.hitCnt
		for i := 0; i < cnt/4+1; i++ {
			size *= 0.9
			a *= 0.8
			if hc > 0 {
				hc--
				continue
			}
			d := float64(i*13 + this.cnt*3)
			for j := 0; j < 6; j++ {
				gl.PushMatrix()
				gl.Translatef(x, y, 0)
				gl.Rotatef(-this._deg*180/math.Pi, 0, 0, 1)
				gl.Rotatef(d, 0, 1, 0)
				Screen.setColor(0.4, 0.8, 0.8, a)
				gl.Begin(gl.LINE_LOOP)
				gl.Vertex3f(-size, LANCE_SPEED, size/2)
				gl.Vertex3f(size, LANCE_SPEED, size/2)
				gl.Vertex3f(size, -LANCE_SPEED, size/2)
				gl.Vertex3f(-size, -LANCE_SPEED, size/2)
				gl.End()
				Screen.setColor(0.2, 0.5, 0.5, a/2)
				gl.Begin(gl.TRIANGLE_FAN)
				gl.Vertex3f(-size, LANCE_SPEED, size/2)
				gl.Vertex3f(size, LANCE_SPEED, size/2)
				gl.Vertex3f(size, -LANCE_SPEED, size/2)
				gl.Vertex3f(-size, -LANCE_SPEED, size/2)
				gl.End()
				gl.PopMatrix()
				d += 60
			}
			x -= math.Sin(this.deg) * LANCE_SPEED * 2
			y -= math.Cos(this.deg) * LANCE_SPEED * 2
		}
	} else {
		gl.PushMatrix()
		Screen.glTranslate(this.pos)
		gl.Rotatef(-this._deg*180/math.Pi, 0, 0, 1)
		gl.Rotatef(this.cnt*31, 0, 1, 0)
		this.shape.draw()
		gl.PopMatrix()
	}
}

func (this *Shot) deg() float64 {
	return this._deg
}

func (this *Shot) damage() int {
	return this._damage
}

func (this *Shot) removed() bool {
	return this.hitCnt > 0
}

type ShotPool struct {
	ActorPool
}

func NewShotPool(n int, args []interface{}) *ShotPool {
	f := func() Actor { return NewShot() }
	this := &ShotPool{
		ActorPool: NewActorPool(f, n, args),
	}
	return this
}

func (this *ShotPool) existsLance() bool {
	for _, a := range this.actor {
		s := a.(*Shot)
		if s.Exists() && (s.lance && !s.removed) {
			return true
		}
	}
	return false
}

type ShotShape struct {
	sdl.CollidableDrawable
}

func NewShotShape() *ShotShape {
	this := &ShotShape{}
	this.CollidableDrawable = sdl.NewCollidableDrawable(this, this, this)
	return this
}

func (this *ShotShape) CreateDisplayList() {
	sdl.SetColor(0.1, 0.33, 0.1, 1)
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
}

func (this *ShotShape) SetCollision() vector.Vector {
	this._collision = vector.Vector{0.33, 0.33}
	return this._collision
}

var _ Collidable = &LanceShape{}

type LanceShape struct {
	CollidableImpl
	_collision vector.Vector
}

func NewLanceShape() *LanceShape {
	this := &LanceShape{NewCollidableImpl()}
	this._collision = Vector{0.66, 0.66}
	return this
}

func (this *LanceShape) collision() vector.Vector {
	return this._collision
}
