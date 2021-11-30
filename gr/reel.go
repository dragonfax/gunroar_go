package main

import (
	"math"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/actor"
	"github.com/dragonfax/gunroar/gr/letter"
	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

const MAX_DIGIT = 16

/**
 * Rolling reel that displays the score.
 */
type ScoreReel struct {
	score, targetScore int
	_actualScore       int
	digit              int
	numReel            [MAX_DIGIT]NumReel
}

func NewScoreReel() *ScoreReel {
	this := &ScoreReel{digit: 1}
	for i := range this.numReel {
		this.numReel[i] = NewNumReelInternal()
	}
	return this
}

func (this *ScoreReel) clear(digit int /* = 9 */) {
	this.score = 0
	this.targetScore = 0
	this._actualScore = 0
	this.digit = digit
	for i := 0; i < digit; i++ {
		this.numReel[i].clear()
	}
}

func (this *ScoreReel) move() {
	for i := 0; i < this.digit; i++ {
		this.numReel[i].move()
	}
}

func (this *ScoreReel) draw(x, y, s float64) {
	lx := x
	ly := y
	for i := 0; i < this.digit; i++ {
		this.numReel[i].draw(lx, ly, s)
		lx -= s * 2
	}
}

func (this *ScoreReel) addReelScore(as int) {
	this.targetScore += as
	ts := this.targetScore
	for i := 0; i < this.digit; i++ {
		this.numReel[i].setTargetDeg(float64(ts) * 360 / 10)
		ts /= 10
		if ts < 0 {
			break
		}
	}
}

func (this *ScoreReel) accelerate() {
	for i := 0; i < this.digit; i++ {
		this.numReel[i].accelerate()
	}
}

func (this *ScoreReel) addActualScore(as int) {
	this._actualScore += as
}

func (this *ScoreReel) actualScore() int {
	return this._actualScore
}

const VEL_MIN = 5

var numReelRand = r.New(r.NewSource(time.Now().Unix()))

func setNumReelRandSeed(seed int64) {
	numReelRand = r.New(r.NewSource(seed))
}

type NumReel struct {
	deg, _targetDeg, ofs, velRatio float64
}

func NewNumReelInternal() NumReel {
	this := NumReel{}
	this.init()
	return this
}

func (this *NumReel) init() {
	this.deg = 0
	this._targetDeg = 0
	this.ofs = 0
	this.velRatio = 1
}

func (this *NumReel) clear() {
	this.init()
}

func (this *NumReel) move() {
	vd := this._targetDeg - this.deg
	vd *= 0.05 * this.velRatio
	if vd < VEL_MIN*this.velRatio {
		vd = VEL_MIN * this.velRatio
	}
	this.deg += vd
	if this.deg > this._targetDeg {
		this.deg = this._targetDeg
	}
}

func nextSignedFloat(rand *r.Rand, n float64) float64 {
	return rand.Float64()*n*2 - n
}

func (this *NumReel) draw(x, y, s float64) {
	n := int((this.deg*10/360+0.99)+1) % 10
	d := math.Mod(this.deg, 360)
	od := d - float64(n)*360/10
	od -= 15
	od = normalizeDeg360(od)
	od *= 1.5
	for i := 0; i < 3; i++ {
		gl.PushMatrix()
		if this.ofs > 0.005 {
			gl.Translated(x+nextSignedFloat(numReelRand, 1)*this.ofs, y+nextSignedFloat(numReelRand, 1)*this.ofs, 0)
		} else {
			gl.Translated(x, y, 0)
		}
		gl.Rotated(od, 1, 0, 0)
		gl.Translated(0, 0, s*2.4)
		gl.Scaled(s, -s, s)
		a := 1 - math.Abs((od+15)/(360/10*1.5))/2
		if a < 0 {
			a = 0
		}
		sdl.SetColor(a, a, a, 1)
		letter.DrawLetterAsIs(n, 2)
		sdl.SetColor(a/2, a/2, a/2, 1)
		letter.DrawLetterAsIs(n, 3)
		gl.PopMatrix()
		n--
		if n < 0 {
			n = 9
		}
		od += 360 / 10 * 1.5
		od = normalizeDeg360(od)
	}
	this.ofs *= 0.95
}

func (this *NumReel) setTargetDeg(td float64) float64 {
	if (td - this._targetDeg) > 1 {
		this.ofs += 0.1
	}
	this._targetDeg = td
	return td
}

func (this *NumReel) accelerate() {
	this.velRatio = 4
}

/**
 * Flying indicator that shows the score and the multiplier.
 */

type IndicatorType int

const (
	SCORE IndicatorType = iota
	MULTIPLIER
)

type FlyingToType int

const (
	RIGHT FlyingToType = iota
	BOTTOM
)

var numIndicatorRand = r.New(r.NewSource(time.Now().Unix()))

const TARGET_Y_MIN = -7
const TARGET_Y_MAX = 7
const TARGET_Y_INTERVAL = 1

var targetY float64 = TARGET_Y_MIN

type Target struct {
	pos             vector.Vector
	flyingTo        FlyingToType
	initialVelRatio float64
	size            float64
	n               int
	cnt             int
}

var _ actor.Actor = &NumIndicator{}

type NumIndicator struct {
	actor.ExistsImpl

	scoreReel ScoreReel
	pos, vel  vector.Vector
	n         int
	typ       IndicatorType
	size      float64
	cnt       int
	alpha     float64
	target    [4]Target
	targetIdx int
	targetNum int
}

func setNumIndicatorRandSeed(seed int64) {
	numIndicatorRand = r.New(r.NewSource(seed))
}

func initTargetY() {
	targetY = TARGET_Y_MIN
}

func getTargetY() float64 {
	ty := targetY
	targetY += TARGET_Y_INTERVAL
	if targetY > TARGET_Y_MAX {
		targetY = TARGET_Y_MIN
	}
	return ty
}

func decTargetY() {
	targetY -= TARGET_Y_INTERVAL
	if targetY < TARGET_Y_MIN {
		targetY = TARGET_Y_MAX
	}
}

func NewNumIndicator() *NumIndicator {
	this := &NumIndicator{alpha: 1}
	return this
}

func (this *NumIndicator) Init(args []interface{}) {
	this.scoreReel = args[0].(ScoreReel)
}

func (this *NumIndicator) setWithVector(n int, typ IndicatorType, size float64, p vector.Vector) {
	this.set(n, typ, size, p.X, p.Y)
}

func (this *NumIndicator) set(n int, typ IndicatorType, size, x, y float64) {
	if this.Exists() && this.typ == SCORE {
		if this.target[this.targetIdx].flyingTo == RIGHT {
			decTargetY()
		}
		this.scoreReel.addReelScore(this.target[this.targetNum-1].n)
	}
	this.n = n
	this.typ = typ
	this.size = size
	this.pos.X = x
	this.pos.Y = y
	this.targetIdx = -1
	this.targetNum = 0
	this.alpha = 0.1
	this.SetExists(true)
}

func (this *NumIndicator) addTarget(x, y float64, flyingTo FlyingToType, initialVelRatio,
	size float64, n, cnt int) {
	this.target[this.targetNum].pos.X = x
	this.target[this.targetNum].pos.Y = y
	this.target[this.targetNum].flyingTo = flyingTo
	this.target[this.targetNum].initialVelRatio = initialVelRatio
	this.target[this.targetNum].size = size
	this.target[this.targetNum].n = n
	this.target[this.targetNum].cnt = cnt
	this.targetNum++
}

func (this *NumIndicator) gotoNextTarget() {
	this.targetIdx++
	if this.targetIdx > 0 {
		playSe("score_up.wav")
	}
	if this.targetIdx >= this.targetNum {
		if this.target[this.targetIdx-1].flyingTo == BOTTOM {
			this.scoreReel.addReelScore(this.target[this.targetIdx-1].n)
		}
		this.SetExists(false)
		return
	}
	switch this.target[this.targetIdx].flyingTo {
	case RIGHT:
		this.vel.X = -0.3 + nextSignedFloat(numIndicatorRand, 0.05)
		this.vel.Y = nextSignedFloat(numIndicatorRand, 0.1)
	case BOTTOM:
		this.vel.X = nextSignedFloat(numIndicatorRand, 0.1)
		this.vel.Y = 0.3 + nextSignedFloat(numIndicatorRand, 0.05)
		decTargetY()
	}
	this.vel.OpMulAssign(this.target[this.targetIdx].initialVelRatio)
	this.cnt = this.target[this.targetIdx].cnt
}

func (this *NumIndicator) Move() {
	if this.targetIdx < 0 {
		return
	}
	tp := this.target[this.targetIdx].pos
	switch this.target[this.targetIdx].flyingTo {
	case RIGHT:
		this.vel.X += (tp.X - this.pos.X) * 0.0036
		this.pos.Y += (tp.Y - this.pos.Y) * 0.1
		if math.Abs(this.pos.Y-tp.Y) < 0.5 {
			this.pos.Y += (tp.Y - this.pos.Y) * 0.33
		}
		this.alpha += (1 - this.alpha) * 0.03
	case BOTTOM:
		this.pos.X += (tp.X - this.pos.X) * 0.1
		this.vel.Y += (tp.Y - this.pos.Y) * 0.0036
		this.alpha *= 0.97
	}
	this.vel.OpMulAssign(0.98)
	this.size += (this.target[this.targetIdx].size - this.size) * 0.025
	this.pos.OpAddAssign(this.vel)
	vn := int(float64(this.target[this.targetIdx].n-this.n) * 0.2)
	if vn < 10 && vn > -10 {
		this.n = this.target[this.targetIdx].n
	} else {
		this.n += vn
	}
	switch this.target[this.targetIdx].flyingTo {
	case RIGHT:
		if this.pos.X > tp.X {
			this.pos.X = tp.X
			this.vel.X *= -0.05
		}
	case BOTTOM:
		if this.pos.Y < tp.Y {
			this.pos.Y = tp.Y
			this.vel.Y *= -0.05
		}
	}
	this.cnt--
	if this.cnt < 0 {
		this.gotoNextTarget()
	}
}

func (this *NumIndicator) Draw() {
	sdl.SetColor(this.alpha, this.alpha, this.alpha, 1)
	switch this.typ {
	case SCORE:
		letter.DrawNumSign(this.n, this.pos.X, this.pos.Y, this.size, letter.LINE_COLOR, -1, -1)
	case MULTIPLIER:
		sdl.SetColor(this.alpha, this.alpha, this.alpha, 1)
		letter.DrawNumSign(this.n, this.pos.X, this.pos.Y, this.size, letter.LINE_COLOR, 33, 3)
	}
}

type NumIndicatorPool struct {
	actor.ActorPool
}

func NewNumIndicatorPool(n int, args []interface{}) *NumIndicatorPool {
	var f actor.CreateActorFunc = func() actor.Actor { return NewNumIndicator() }
	this := &NumIndicatorPool{
		ActorPool: actor.NewActorPool(f, n, args),
	}

	return this
}

func (this *NumIndicatorPool) GetInstance() *NumIndicator {
	return this.ActorPool.GetInstance().(*NumIndicator)
}

func (this *NumIndicatorPool) GetInstanceForced() *NumIndicator {
	return this.ActorPool.GetInstance().(*NumIndicator)
}
