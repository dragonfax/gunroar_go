/*
 * $Id: reel.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

/**
 * Rolling reel that displays the score.
 */

import (
	"github.com/go-gl/gl"
)

const MAX_DIGIT = 16

var scoreReel *ScoreReel

type ScoreReel struct {
	score, targetScore int
	actualScore        int
	digit              int
	numReel            [MAX_DIGIT]*NumReel
}

func NewScoreReel() *ScoreReel {
	sr := new(ScoreReel)
	for i, _ := range sr.numReel {
		sr.numReel[i] = NewNumReel()
	}
	sr.digit = 1
	actorsLock.Lock()
	actors[sr] = true
	actorsLock.Unlock()
	return sr
}

func (sr *ScoreReel) clear(digit int /*= 9 */) {
	sr.score = 0
	sr.targetScore = 0
	sr.actualScore = 0
	sr.digit = digit
	for i := 0; i < digit; i++ {
		sr.numReel[i].close()
		sr.numReel[i] = NewNumReel()
	}
}

func (sr *ScoreReel) move() {
	for i := 0; i < sr.digit; i++ {
		sr.numReel[i].move()
	}
}

/* draw() to satisfy Actor */
func (sr *ScoreReel) draw() {
}

func (sr *ScoreReel) drawAtPos(x float32, y float32, s float32) {
	lx := x
	ly := y
	for i := 0; i < sr.digit; i++ {
		sr.numReel[i].drawAtPos(lx, ly, s)
		lx -= s * 2
	}
}

func (sr *ScoreReel) addReelScore(as int) {
	sr.targetScore += as
	var ts int = sr.targetScore
	for i := 0; i < sr.digit; i++ {
		sr.numReel[i].targetDeg = float32(ts) * 360 / 10
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

func (sr *ScoreReel) close() {
	actorsLock.Lock()
	delete(actors, sr)
	actorsLock.Unlock()
}

const VEL_MIN float32 = 5

type NumReel struct {
	deg       float32
	targetDeg float32
	ofs       float32
	velRatio  float32
}

func NewNumReel() *NumReel {
	nr := new(NumReel)
	nr.velRatio = 1
	actorsLock.Lock()
	actors[nr] = true
	actorsLock.Unlock()
	return nr
}

/* to satisfy Actor */
func (nr *NumReel) draw() {
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

func (nr *NumReel) drawAtPos(x float32, y float32, s float32) {
	var n int = int(Mod32(((nr.deg*10/360 + 0.99) + 1), 10))
	var d float32 = Mod32(nr.deg, 360)
	var od float32 = d - float32(n)*360/10
	od -= 15
	od = normalizeDeg360(od)
	od *= 1.5
	for i := 0; i < 3; i++ {
		gl.PushMatrix()
		if nr.ofs > 0.005 {
			gl.Translatef(x+nextSignedFloat(nr.ofs), y+nextSignedFloat(nr.ofs), 0)
		} else {
			gl.Translatef(x, y, 0)
		}
		gl.Rotatef(od, 1, 0, 0)
		gl.Translatef(0, 0, s*2.4)
		gl.Scalef(s, -s, s)
		a := float32(1 - fabs32((od+15)/(360/10*1.5))/2)
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

/* TODO this may be mis-used in place of a call to targetDeg */
func (nr *NumReel) targetDegUpdate(td float32) float32 {
	if (td - nr.targetDeg) > 1 {
		nr.ofs += 0.1
	}
	nr.targetDeg = td
	return nr.targetDeg
}

func (nr *NumReel) accelerate() {
	nr.velRatio = 4
}

func (nr *NumReel) close() {
	actorsLock.Lock()
	delete(actors, nr)
	actorsLock.Unlock()
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

func InitNumIndicators() {
	targetY = TARGET_Y_MIN
}

func InitTargetY() {
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

func NewNumIndicator(n int, t IndicatorType, size float32, x float32, y float32) *NumIndicator {
	ni := new(NumIndicator)
	ni.alpha = 1

	if ni.t == IndicatorTypeSCORE {
		if ni.target[ni.targetIdx].flyingTo == FlyingToTypeRIGHT {
			decTargetY()
		}
		scoreReel.addReelScore(ni.target[ni.targetNum-1].n)
	}
	ni.n = n
	ni.t = t
	ni.size = size
	ni.pos = Vector{x, y}
	ni.targetIdx = -1
	ni.alpha = 0.1
	actorsLock.Lock()
	actors[ni] = true
	actorsLock.Unlock()
	return ni
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
			scoreReel.addReelScore(ni.target[ni.targetIdx-1].n)
		}
		ni.close()
		return
	}
	switch ni.target[ni.targetIdx].flyingTo {
	case FlyingToTypeRIGHT:
		x := -0.3 + nextSignedFloat(0.05)
		y := nextSignedFloat(0.1)
		ni.vel = Vector{x, y}
	case FlyingToTypeBOTTOM:
		x := nextSignedFloat(0.1)
		y := 0.3 + nextSignedFloat(0.05)
		ni.vel = Vector{x, y}
		decTargetY()
	}
	ni.vel.MulAssign(ni.target[ni.targetIdx].initialVelRatio)
	ni.cnt = ni.target[ni.targetIdx].cnt
}

func (ni *NumIndicator) move() {
	if ni.targetIdx < 0 {
		return
	}
	tp := ni.target[ni.targetIdx].pos
	switch ni.target[ni.targetIdx].flyingTo {
	case FlyingToTypeRIGHT:
		x := (tp.x - ni.pos.x) * 0.0036
		ni.vel = Vector{ni.vel.x + x, ni.vel.y}
		y := (tp.y - ni.pos.y) * 0.1
		ni.pos = Vector{ni.pos.x, ni.pos.y + y}
		if fabs32(ni.pos.y-tp.y) < 0.5 {
			y := (tp.y - ni.pos.y) * 0.33
			ni.pos = Vector{ni.pos.x, ni.pos.y + y}
		}
		ni.alpha += (1 - ni.alpha) * 0.03
	case FlyingToTypeBOTTOM:
		/* I was here with the conversions */
		ni.pos = Vector{ni.pos.x + (tp.x-ni.pos.x)*0.1, ni.pos.y}
		ni.vel = Vector{ni.vel.x, ni.vel.y + (tp.y-ni.pos.y)*0.0036}
		ni.alpha *= 0.97
	}
	ni.vel.MulAssign(0.98)
	ni.size += (ni.target[ni.targetIdx].size - ni.size) * 0.025
	ni.pos.AddAssign(ni.vel)
	var vn int = int(float32(ni.target[ni.targetIdx].n-ni.n) * 0.2)
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
	case FlyingToTypeBOTTOM:
		if ni.pos.y < tp.y {
			ni.pos.y = tp.y
			ni.vel.y *= -0.05
		}
	}
	ni.cnt--
	if ni.cnt < 0 {
		ni.gotoNextTarget()
	}
}

func (ni *NumIndicator) draw() {
	setScreenColor(ni.alpha, ni.alpha, ni.alpha, 1)
	switch ni.t {
	case IndicatorTypeSCORE:
		drawNumSignOption(ni.n, ni.pos.x, ni.pos.y, ni.size, LETTER_LINE_COLOR, -1, -1)
	case IndicatorTypeMULTIPLIER:
		setScreenColor(ni.alpha, ni.alpha, ni.alpha, 1)
		drawNumSignOption(ni.n, ni.pos.x, ni.pos.y, ni.size, LETTER_LINE_COLOR, 33, 3)
	}
}

func (ni *NumIndicator) close() {
	actorsLock.Lock()
	delete(actors, ni)
	actorsLock.Unlock()
}
