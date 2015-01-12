/*
 * $Id: reel.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Rolling reel that displays the score.
 */

import (
	"github.com/go-gl/gl"
	"math/rand"
)

const MAX_DIGIT = 16

type ScoreReel struct {
	score, targetScore int
	actualScore        int
	digit              int
	numReel            [MAX_DIGIT]NumReel
}

func (sr *ScoreReel) Init() {
	for i, _ := range sr.numReel {
		sr.numReel[i].Init()
	}
	sr.digit = 1
}

func (sr *ScoreReel) clear(digit int /*= 9 */) {
	sr.score = 0
	sr.targetScore = 0
	sr.actualScore = 0
	sr.digit = digit
	for i := 0; i < digit; i++ {
		sr.numReel[i].clear()
	}
}

func (sr *ScoreReel) move() {
	for i := 0; i < sr.digit; i++ {
		sr.numReel[i].move()
	}
}

func (sr *ScoreReel) draw(x float32, y float32, s float32) {
	lx := x
	ly := y
	for i := 0; i < sr.digit; i++ {
		sr.numReel[i].draw(lx, ly, s)
		lx -= s * 2
	}
}

func (sr *ScoreReel) addReelScore(as int) {
	sr.targetScore += as
	ts := sr.targetScore
	for i := 0; i < sr.digit; i++ {
		sr.numReel[i].targetDeg = float32(ts * 360 / 10)
		ts /= 10
		if ts < 0 {
			break
		}
	}
}

func (sr *ScoreReel) accelerate() {
	for i := 0; i < sr.digit; i++ {
		sr.numReel[i].accelerate()
	}
}

func (sr *ScoreReel) addActualScore(as int) {
	sr.actualScore += as
}

const VEL_MIN float32 = 5

type NumReel struct {
	deg       float32
	targetDeg float32
	ofs       float32
	velRatio  float32
}

func (nr *NumReel) Init() {
	nr.deg = 0
	nr.ofs = 0
	nr.velRatio = 1
}

func (nr *NumReel) clear() {
	nr.Init()
}

func (nr *NumReel) move() {
	vd := nr.targetDeg - nr.deg
	vd *= 0.05 * nr.velRatio
	if vd < VEL_MIN*nr.velRatio {
		vd = VEL_MIN * nr.velRatio
	}
	nr.deg += vd
	if nr.deg > nr.targetDeg {
		nr.deg = nr.targetDeg
	}
}

func (nr *NumReel) draw(x float32, y float32, s float32) {
	n := Mod32(((nr.deg*10/360 + 0.99) + 1), 10)
	d := Mod32(nr.deg, 360)
	od := d - n*360/10
	od -= 15
	od = normalizeDeg360(od)
	od *= 1.5
	for i := 0; i < 3; i++ {
		gl.PushMatrix()
		if nr.ofs > 0.005 {
			gl.Translatef(x+rand.Float32()*nr.ofs, y+rand.Float32()*nr.ofs, 0)
		} else {
			gl.Translatef(x, y, 0)
		}
		gl.Rotatef(od, 1, 0, 0)
		gl.Translatef(0, 0, s*2.4)
		gl.Scalef(s, -s, s)
		a := float32(1 - fabs((od+15)/(360/10*1.5))/2)
		if a < 0 {
			a = 0
		}
		setScreenColor(a, a, a, 1)
		drawLetter(int(n), 2)
		setScreenColor(a/2, a/2, a/2, 1)
		drawLetter(int(n), 3)
		gl.PopMatrix()
		n--
		if n < 0 {
			n = 9
		}
		od += 360 / 10 * 1.5
		od = normalizeDeg360(od)
	}
	nr.ofs *= 0.95
}

func (nr *NumReel) targetDeg(td float32) float32 {
	if (td - nr.targetDeg) > 1 {
		nr.ofs += 0.1
	}
	nr.targetDeg = td
	return nr.targetDeg
}

func (nr *NumReel) accelerate() {
	nr.velRatio = 4
}

/**
 * Flying indicator that shows the score and the multiplier.
 */

type IndicatorType int

const (
	IndicatorTypeSCORE      IndicatorType = 1
	IndicatorTypeMULTIPLIER               = 2
)

type FlyingToType int

const (
	FlyingToTypeRIGHT  FlyingToType = 1
	FlyingToTypeBOTTOM              = 2
)

const TARGET_Y_MIN = -7
const TARGET_Y_MAX = 7
const TARGET_Y_INTERVAL = 1

var targetY float32

type Target struct {
	pos             Vector
	flyingTo        FlyingToType
	initialVelRatio float32
	size            float32
	n               int
	cnt             int
}

type NumIndicator struct {
	ActorImpl

	scoreReel *ScoreReel
	pos, vel  Vector
	n         int
	t         IndicatorType
	size      float32
	cnt       int
	alpha     float32
	target    [4]Target
	targetIdx int
	targetNum int
}

func InitNumIndicator() {
	targetY = TARGET_Y_MIN
}

func initTargetY() {
	targetY = TARGET_Y_MIN
}

func getTargetY() float32 {
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

func (ni *NumIndicator) Init() {
	ni.pos = Vector{}
	ni.vel = Vector{}
	for _, t := range ni.target {
		t.pos = Vector{}
		t.initialVelRatio = 0
		t.size = 0
	}
	ni.targetIdx = 0
	ni.targetNum = 0
	ni.alpha = 1
}

func (ni *NumIndicator) set(n int, t IndicatorType, size float32, x float32, y float32) {
	if ni.Exists() && ni.t == IndicatorTypeSCORE {
		if ni.target[ni.targetIdx].flyingTo == FlyingToTypeRIGHT {
			decTargetY()
		}
		ni.scoreReel.addReelScore(ni.target[ni.targetNum-1].n)
	}
	ni.n = n
	ni.t = t
	ni.size = size
	ni.pos = Vector{x, y}
	ni.targetIdx = -1
	ni.targetNum = 0
	ni.alpha = 0.1
	ni.SetExists(true)
}

func (ni *NumIndicator) addTarget(x float32, y float32, flyingTo FlyingToType, initialVelRatio float32,
	size float32, n int, cnt int) {
	ni.target[ni.targetNum].pos = Vector{x, y}
	ni.target[ni.targetNum].flyingTo = flyingTo
	ni.target[ni.targetNum].initialVelRatio = initialVelRatio
	ni.target[ni.targetNum].size = size
	ni.target[ni.targetNum].n = n
	ni.target[ni.targetNum].cnt = cnt
	ni.targetNum++
}

func (ni *NumIndicator) gotoNextTarget() {
	ni.targetIdx++
	if ni.targetIdx > 0 {
		playSe("score_up.wav")
	}
	if ni.targetIdx >= ni.targetNum {
		if ni.target[ni.targetIdx-1].flyingTo == FlyingToTypeBOTTOM {
			ni.scoreReel.addReelScore(ni.target[ni.targetIdx-1].n)
		}
		ni.SetExists(false)
		return
	}
	switch ni.target[ni.targetIdx].flyingTo {
	case FlyingToTypeRIGHT:
		x := -0.3 + rand.Float32()*0.05
		y := rand.Float32() * 0.1
		ni.vel = Vector{x, y}
		break
	case FlyingToTypeBOTTOM:
		x := rand.Float32() * 0.1
		y := 0.3 + rand.Float32()*0.05
		ni.vel = Vector{x, y}
		decTargetY()
		break
	}
	ni.vel = ni.vel.Mul(ni.target[ni.targetIdx].initialVelRatio)
	ni.cnt = ni.target[ni.targetIdx].cnt
}

func (ni *NumIndicator) move() {
	if ni.targetIdx < 0 {
		return
	}
	tp := ni.target[ni.targetIdx].pos
	switch ni.target[ni.targetIdx].flyingTo {
	case FlyingToTypeRIGHT:
		x := (tp.X() - ni.pos.X()) * 0.0036
		ni.vel = Vector{ni.vel.x + x, ni.vel.y}
		y := (tp.Y() - ni.pos.Y()) * 0.1
		ni.pos = Vector{ni.pos.x, ni.pos.y + y}
		if fabs(ni.pos.Y()-tp.Y()) < 0.5 {
			y := (tp.Y() - ni.pos.Y()) * 0.33
			ni.pos = Vector{ni.pos.x, ni.pos.y + y}
		}
		ni.alpha += (1 - ni.alpha) * 0.03
		break
	case FlyingToTypeBOTTOM:
		/* I was here with the conversions */
		ni.pos.x += (tp.X() - ni.pos.X()) * 0.1
		ni.vel.y += (tp.Y() - ni.pos.Y()) * 0.0036
		ni.alpha *= 0.97
		break
	}
	ni.vel *= 0.98
	ni.size += (ni.target[ni.targetIdx].size - ni.size) * 0.025
	ni.pos += ni.vel
	vn := int((ni.target[ni.targetIdx].n - ni.n) * 0.2)
	if vn < 10 && vn > -10 {
		ni.n = ni.target[ni.targetIdx].n
	} else {
		ni.n += vn
	}
	switch ni.target[ni.targetIdx].flyingTo {
	case FlyingToTypeRIGHT:
		if ni.pos.x > tp.x {
			ni.pos.x = tp.x
			ni.vel.x *= -0.05
		}
		break
	case FlyingToTypeBOTTOM:
		if ni.pos.y < tp.y {
			ni.pos.y = tp.y
			ni.vel.y *= -0.05
		}
		break
	}
	ni.cnt--
	if ni.cnt < 0 {
		ni.gotoNextTarget()
	}
}

func (ni *NumIndicator) draw() {
	setScrenColor(ni.alpha, ni.alpha, ni.alpha, 1)
	switch ni.t {
	case IndicatorTypeSCORE:
		Letter.drawNumSign(ni.n, ni.pos.x, ni.pos.y, ni.size, Letter.LINE_COLOR)
		break
	case IndicatorTypeMULTIPLIER:
		setScreenColor(ni.alpha, ni.alpha, ni.alpha, 1)
		Letter.drawNumSign(ni.n, ni.pos.x, ni.pos.y, ni.size, Letter.LINE_COLOR, 33, 3)
		break
	}
}
