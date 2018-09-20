/*
 * $Id: letter.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

const LETTER_WIDTH = 2.1
const LETTER_HEIGHT = 3.0
const LETTER_LINE_COLOR = 2
const POLY_COLOR = 3
const COLOR_NUM = 4

var COLOR_RGB = [][]float32{[]float32{1, 1, 1}, []float32{0.9, 0.7, 0.5}}

const LETTER_NUM = 44
const DISPLAY_LIST_NUM = LETTER_NUM * COLOR_NUM

var displayList *DisplayList

func InitLetter() {
	displayList = NewDisplayList(DISPLAY_LIST_NUM)
	displayList.ResetLists()
	for j := 0; j < COLOR_NUM; j++ {
		for i := 0; i < LETTER_NUM; i++ {
			displayList.NewList()
			setLetter(i, j)
			displayList.EndList()
		}
	}
}

func CloseLetter() {
	displayList.close()
}

func getLetterWidth(n int, s float32) float32 {
	return float32(n) * s * LETTER_WIDTH
}

func getLetterHeight(s float32) float32 {
	return s * LETTER_HEIGHT
}

func drawLetter(n uint32, c uint32) {
	displayList.call(n + c*uint32(LETTER_NUM))
}

func drawLetterOption(n uint32, x float32, y float32, s float32, d float32, c uint32) {
	gl.PushMatrix()
	gl.Translatef(x, y, 0)
	gl.Scalef(s, s, s)
	gl.Rotatef(float32(d), 0, 0, 1)
	displayList.call(n + c*uint32(LETTER_NUM))
	gl.PopMatrix()
}

func drawLetterRev(n uint32, x float32, y float32, s float32, d float32, c uint32) {
	gl.PushMatrix()
	gl.Translatef(x, y, 0)
	gl.Scalef(s, -s, s)
	gl.Rotatef(float32(d), 0, 0, 1)
	displayList.call(n + c*uint32(LETTER_NUM))
	gl.PopMatrix()
}

type Direction int

const ( // Direction
	TO_RIGHT Direction = iota
	TO_DOWN
	TO_LEFT
	TO_UP
)

func convertCharToInt(c rune) uint32 {
	var idx uint32
	if c >= '0' && c <= '9' {
		idx = uint32(c - '0')
	} else if c >= 'A' && c <= 'Z' {
		idx = uint32(c-'A') + 10
	} else if c >= 'a' && c <= 'z' {
		idx = uint32(c-'a') + 10
	} else if c == '.' {
		idx = 36
	} else if c == '-' {
		idx = 38
	} else if c == '+' {
		idx = 39
	} else if c == '_' {
		idx = 37
	} else if c == '!' {
		idx = 42
	} else if c == '/' {
		idx = 43
	}
	return idx
}

func drawString(str string, lx float32, y float32, s float32) {
	drawStringOption(str, lx, y, s, TO_RIGHT, 0, false, 0)
}

func drawStringOption(str string, lx float32, y float32, s float32, d Direction, cl uint32, rev bool, od float32) {
	lx += LETTER_WIDTH * s / 2
	y += LETTER_HEIGHT * s / 2
	x := lx
	var idx uint32
	var ld float32
	switch d {
	case TO_RIGHT:
		ld = 0
	case TO_DOWN:
		ld = 90
	case TO_LEFT:
		ld = 180
	case TO_UP:
		ld = 270
	}
	ld += od
	var c rune
	for _, c = range str {
		if c != ' ' {
			idx = convertCharToInt(c)
			if rev {
				drawLetterRev(idx, x, y, s, ld, cl)
			} else {
				drawLetterOption(idx, x, y, s, ld, cl)
			}
		}
		if od == 0 {
			switch d {
			case TO_RIGHT:
				x += s * LETTER_WIDTH
			case TO_DOWN:
				y += s * LETTER_WIDTH
			case TO_LEFT:
				x -= s * LETTER_WIDTH
			case TO_UP:
				y -= s * LETTER_WIDTH
			}
		} else {
			x += Cos32(ld*math.Pi/180) * s * LETTER_WIDTH
			y += Sin32(ld*math.Pi/180) * s * LETTER_WIDTH
		}
	}
}

func drawNum(num uint32, lx float32, y float32, s float32) {
	drawNumOption(num, lx, y, s, 0, 0, -1, -1)
}

func drawNumOption(num uint32, lx float32, y float32, s float32, cl uint32, dg uint32, headChar int, floatDigit int) {
	lx += LETTER_WIDTH * s / 2
	y += LETTER_HEIGHT * s / 2
	n := num
	x := lx
	var ld float32 = 0 // TO_RIGHT
	digit := dg
	var fd int = floatDigit
	for {
		if fd <= 0 {
			drawLetterOption(n%10, x, y, s, ld, cl)
			x -= s * LETTER_WIDTH
		} else {
			drawLetterOption(n%10, x, y+s*LETTER_WIDTH*0.25, s*0.5, ld, cl)
			x -= s * LETTER_WIDTH * 0.5
		}
		n /= 10
		digit--
		fd--
		if n <= 0 && digit <= 0 && fd < 0 {
			break
		}
		if fd == 0 {
			drawLetterOption(36, x, y+s*LETTER_WIDTH*0.25, s*0.5, ld, cl)
			x -= s * LETTER_WIDTH * 0.5
		}
	}
	if headChar >= 0 {
		drawLetterOption(uint32(headChar), x+s*LETTER_WIDTH*0.2, y+s*LETTER_WIDTH*0.2, s*0.6, ld, cl)
	}
}

func drawNumSign(num uint32, lx float32, ly float32, s float32) {
	drawNumSignOption(num, lx, ly, s, 0, -1, -1)
}

func drawNumSignOption(num uint32, lx float32, ly float32, s float32, cl uint32, headChar int, floatDigit int) {
	x := lx
	y := ly
	n := num
	fd := floatDigit
	for {
		if fd <= 0 {
			drawLetterRev(n%10, x, y, s, 0, cl)
			x -= s * LETTER_WIDTH
		} else {
			drawLetterRev(n%10, x, y-s*LETTER_WIDTH*0.25, s*0.5, 0, cl)
			x -= s * LETTER_WIDTH * 0.5
		}
		n /= 10
		if n <= 0 {
			break
		}
		fd--
		if fd == 0 {
			drawLetterRev(36, x, y-s*LETTER_WIDTH*0.25, s*0.5, 0, cl)
			x -= s * LETTER_WIDTH * 0.5
		}
	}
	if headChar >= 0 {
		drawLetterRev(uint32(headChar), x+s*LETTER_WIDTH*0.2, y-s*LETTER_WIDTH*0.2, s*0.6, 0, cl)
	}
}

func drawTime(time uint32, lx float32, y float32, s float32, cl uint32 /* default 0 */) {
	n := time
	if n < 0 {
		n = 0
	}
	var x float32 = lx
	for i := 0; i < 7; i++ {
		if i != 4 {
			drawLetterOption(n%10, x, y, s, 0, cl)
			n /= 10
		} else {
			drawLetterOption(n%6, x, y, s, 0, cl)
			n /= 6
		}
		if (i&1) == 1 || i == 0 {
			switch i {
			case 3:
				drawLetterOption(41, x+s*1.16, y, s, 0, cl)
			case 5:
				drawLetterOption(40, x+s*1.16, y, s, 0, cl)
			}
			x -= s * LETTER_WIDTH
		} else {
			x -= s * LETTER_WIDTH * 1.3
		}
		if n <= 0 {
			break
		}
	}
}

func setLetter(idx int, c int) {
	var x, y, length, size, deg float32
	for i := 0; ; i++ {
		deg = floor32(spData[idx][i][4])
		if deg > 99990 {
			break
		}
		x = -spData[idx][i][0]
		y = -spData[idx][i][1]
		size = spData[idx][i][2]
		length = spData[idx][i][3]
		y *= 0.9
		size *= 1.4
		length *= 1.05
		x = -x
		// y = y
		deg = Mod32(deg, 180)
		if c == LETTER_LINE_COLOR {
			setBoxLine(x, y, size, length, deg)
		} else if c == POLY_COLOR {
			setBoxPoly(x, y, size, length, deg)
		} else {
			setBox(x, y, size, length, deg,
				COLOR_RGB[c][0], COLOR_RGB[c][1], COLOR_RGB[c][2])
		}
	}
}

func Mod32(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}

func setBox(x float32, y float32, width float32, height float32, deg float32, r float32, g float32, b float32) {
	gl.PushMatrix()
	gl.Translatef(x-width/2, y-height/2, 0)
	gl.Rotatef(deg, 0, 0, 1)
	setLetterColorAlpha(r, g, b, 0.5)
	gl.Begin(gl.TRIANGLE_FAN)
	setBoxPart(width, height)
	gl.End()
	setLetterColor(r, g, b)
	gl.Begin(gl.LINE_LOOP)
	setBoxPart(width, height)
	gl.End()
	gl.PopMatrix()
}

func setLetterColor(red float32, green float32, blue float32) {
	gl.Color3f(red, green, blue)
}

func setLetterColorAlpha(red float32, green float32, blue float32, alpha float32) {
	gl.Color4f(red, green, blue, alpha)
}

func setBoxLine(x float32, y float32, width float32, height float32, deg float32) {
	gl.PushMatrix()
	gl.Translatef(x-width/2, y-height/2, 0)
	gl.Rotatef(deg, 0, 0, 1)
	gl.Begin(gl.LINE_LOOP)
	setBoxPart(width, height)
	gl.End()
	gl.PopMatrix()
}

func setBoxPoly(x float32, y float32, width float32, height float32, deg float32) {
	gl.PushMatrix()
	gl.Translatef(x-width/2, y-height/2, 0)
	gl.Rotatef(deg, 0, 0, 1)
	gl.Begin(gl.TRIANGLE_FAN)
	setBoxPart(width, height)
	gl.End()
	gl.PopMatrix()
}

func setBoxPart(width float32, height float32) {
	gl.Vertex3f(-width/2, 0, 0)
	gl.Vertex3f(-width/3*1, -height/2, 0)
	gl.Vertex3f(width/3*1, -height/2, 0)
	gl.Vertex3f(width/2, 0, 0)
	gl.Vertex3f(width/3*1, height/2, 0)
	gl.Vertex3f(-width/3*1, height/2, 0)
}

var spData = [][][]float32{
	[][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.6, 0.55, 0.65, 0.3, 90}, []float32{0.6, 0.55, 0.65, 0.3, 90},
		[]float32{-0.6, -0.55, 0.65, 0.3, 90}, []float32{0.6, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0.5, 0.55, 0.65, 0.3, 90},
		[]float32{0.5, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //A
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.18, 1.15, 0.45, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.45, 0.55, 0.65, 0.3, 90},
		[]float32{-0.18, 0, 0.45, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.15, 1.15, 0.45, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.45, 0.45, 0.65, 0.3, 90},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //F
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0.05, 0, 0.3, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 0.55, 0.65, 0.3, 90},
		[]float32{0, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0.65, -0.55, 0.65, 0.3, 90}, []float32{-0.7, -0.7, 0.3, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //K
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.4, 0.55, 0.65, 0.3, 100},
		[]float32{-0.25, 0, 0.45, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.6, -0.55, 0.65, 0.3, 80},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.5, 1.15, 0.3, 0.3, 0}, []float32{0.1, 1.15, 0.3, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0.55, 0.65, 0.3, 90},
		[]float32{0, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //P
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0.05, -0.55, 0.45, 0.3, 60},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.2, 0, 0.45, 0.3, 0},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.45, -0.55, 0.65, 0.3, 80},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{-0.65, 0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0.65, 0.3, 0},
		[]float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.5, 1.15, 0.55, 0.3, 0}, []float32{0.5, 1.15, 0.55, 0.3, 0},
		[]float32{0.1, 0.55, 0.65, 0.3, 90},
		[]float32{0.1, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //U
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.5, -0.55, 0.65, 0.3, 90}, []float32{0.5, -0.55, 0.65, 0.3, 90},
		[]float32{-0.1, -1.15, 0.45, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.65, 0.55, 0.65, 0.3, 90}, []float32{0.65, 0.55, 0.65, 0.3, 90},
		[]float32{-0.65, -0.55, 0.65, 0.3, 90}, []float32{0.65, -0.55, 0.65, 0.3, 90},
		[]float32{-0.5, -1.15, 0.3, 0.3, 0}, []float32{0.1, -1.15, 0.3, 0.3, 0},
		[]float32{0, 0.55, 0.65, 0.3, 90},
		[]float32{0, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.4, 0.6, 0.85, 0.3, 360 - 120},
		[]float32{0.4, 0.6, 0.85, 0.3, 360 - 60},
		[]float32{-0.4, -0.6, 0.85, 0.3, 360 - 240},
		[]float32{0.4, -0.6, 0.85, 0.3, 360 - 300},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{-0.4, 0.6, 0.85, 0.3, 360 - 120},
		[]float32{0.4, 0.6, 0.85, 0.3, 360 - 60},
		[]float32{-0.1, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{
		[]float32{0, 1.15, 0.65, 0.3, 0},
		[]float32{0.3, 0.4, 0.65, 0.3, 120},
		[]float32{-0.3, -0.4, 0.65, 0.3, 120},
		[]float32{0, -1.15, 0.65, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //.
		[]float32{0, -1.15, 0.3, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //_
		[]float32{0, -1.15, 0.8, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //-
		[]float32{0, 0, 0.9, 0.3, 0},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //+
		[]float32{-0.5, 0, 0.45, 0.3, 0}, []float32{0.45, 0, 0.45, 0.3, 0},
		[]float32{0.1, 0.55, 0.65, 0.3, 90},
		[]float32{0.1, -0.55, 0.65, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //'
		[]float32{0, 1.0, 0.4, 0.2, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //''
		[]float32{-0.19, 1.0, 0.4, 0.2, 90},
		[]float32{0.2, 1.0, 0.4, 0.2, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ //!
		[]float32{0.56, 0.25, 1.1, 0.3, 90},
		[]float32{0, -1.0, 0.3, 0.3, 90},
		[]float32{0, 0, 0, 0, 99999},
	}, [][]float32{ // /
		[]float32{0.8, 0, 1.75, 0.3, 120},
		[]float32{0, 0, 0, 0, 99999},
	}}
