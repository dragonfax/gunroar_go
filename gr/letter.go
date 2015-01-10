/*
 * $Id: letter.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

import (
	"./sdl"
	"github.com/go-gl/gl"
	// "github.com/veandco/go-sdl2/sdl"
)

const LETTER_WIDTH = 2.1
const LETTER_HEIGHT = 3.0
const LINE_COLOR = 2
const POLY_COLOR = 3
const COLOR_NUM = 4

var COLOR_RGB = [][]float32{[]float32{1, 1, 1}, []float32{0.9, 0.7, 0.5}}

const LETTER_NUM int = 44
const DISPLAY_LIST_NUM int = LETTER_NUM * COLOR_NUM

type Letter struct {
	DisplayList sdl.DisplayList
}

func (l *Letter) Init() {
	l.DisplayList = NewDisplayList(DISPLAY_LIST_NUM)
	l.DisplayList.ResetList()
	for j := 0; j < COLOR_NUM; j++ {
		for i := 0; i < LETTER_NUM; i++ {
			l.DisplayList.NewList()
			l.SetLetter(i, j)
			l.DisplayList.EndList()
		}
	}
}

func (l *Letter) close() {
	l.DisplayList.Close()
}

func getWidth(n int, s float32) float32 {
	return n * s * LETTER_WIDTH
}

func getHeight(s float32) float32 {
	return s * LETTER_HEIGHT
}

func (l *Letter) drawLetter(n int, c int) {
	l.DisplayList.Call(n + c*LETTER_NUM)
}

func (l *Letter) drawLetterOption(n int, x float32, y float32, s float32, d float32, c int) {
	gl.PushMatrix()
	gl.Translatef(x, y, 0)
	gl.Scalef(s, s, s)
	gl.Rotatef(d, 0, 0, 1)
	l.DisplayList.Call(n + c*LETTER_NUM)
	gl.PopMatrix()
}

func (l *Letter) drawLetterRev(n int, x float32, y float32, s float32, d float32, c int) {
	gl.PushMatrix()
	gl.Translatef(x, y, 0)
	gl.Scalef(s, -s, s)
	gl.Rotatef(d, 0, 0, 1)
	l.DisplayList.Call(n + c*LETTER_NUM)
	gl.PopMatrix()
}

type Direction int

const ( // Direction
	TO_RIGHT Direction = iota
	TO_DOWN
	TO_LEFT
	TO_UP
)

func ConvertCharToInt(c char) int {
	var idx int
	if c >= '0' && c <= '9' {
		idx = c - '0'
	} else if c >= 'A' && c <= 'Z' {
		idx = c - 'A' + 10
	} else if c >= 'a' && c <= 'z' {
		idx = c - 'a' + 10
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

func DrawString(str []char, lx float32, y float32, s float32) {
	DrawStringOption(str, lx, y, s, TO_RIGHT, 0, false, 0)
}

func DrawStringOption(str []char, lx float32, y float32, s float32,
	d Direction,
	cl int,
	rev bool,
	od float32) {
	lx += LETTER_WIDTH * s / 2
	y += LETTER_HEIGHT * s / 2
	x := lx
	var int idx
	var float32 ld
	switch d {
	case TO_RIGHT:
		ld = 0
		break
	case TO_DOWN:
		ld = 90
		break
	case TO_LEFT:
		ld = 180
		break
	case TO_UP:
		ld = 270
		break
	}
	ld += od
	for char = range c {
		if c != ' ' {
			idx = convertCharToInt(c)
			if rev {
				DrawLetterRev(idx, x, y, s, ld, cl)
			} else {
				DrawLetter(idx, x, y, s, ld, cl)
			}
		}
		if od == 0 {
			switch d {
			case TO_RIGHT:
				x += s * LETTER_WIDTH
				break
			case TO_DOWN:
				y += s * LETTER_WIDTH
				break
			case TO_LEFT:
				x -= s * LETTER_WIDTH
				break
			case TO_UP:
				y -= s * LETTER_WIDTH
				break
			}
		} else {
			x += cos(ld*PI/180) * s * LETTER_WIDTH
			y += sin(ld*PI/180) * s * LETTER_WIDTH
		}
	}
}

func DrawNum(num int, lx float32, y float32, s float32) {
	DrawNumOption(num, lx, y, s, 0, 0, -1, -1)
}

func DrawNumOption(num int, lx float32, y float32, s float32,
	cl int, dg int,
	headChar int, float32Digit int) {
	lx += LETTER_WIDTH * s / 2
	y += LETTER_HEIGHT * s / 2
	n := num
	x := lx
	var ld flaot = 0
	digit := dg
	var fd int = float32Digit
	for {
		if fd <= 0 {
			DrawLetterOption(n%10, x, y, s, ld, cl)
			x -= s * LETTER_WIDTH
		} else {
			DrawLetterOption(n%10, x, y+s*LETTER_WIDTH*0.25, s*0.5, ld, cl)
			x -= s * LETTER_WIDTH * 0.5
		}
		n /= 10
		digit--
		fd--
		if n <= 0 && digit <= 0 && fd < 0 {
			break
		}
		if fd == 0 {
			DrawLetter(36, x, y+s*LETTER_WIDTH*0.25, s*0.5, ld, cl)
			x -= s * LETTER_WIDTH * 0.5
		}
	}
	if headChar >= 0 {
		drawLetter(headChar, x+s*LETTER_WIDTH*0.2, y+s*LETTER_WIDTH*0.2, s*0.6, ld, cl)
	}
}

func DrawNumSign(num int, lx float32, ly float32, s float32) {
	DrawNumSignOption(num, lx, ly, s, 0, -1, -1)
}

func DrawNumSignOption(num int, lx float32, ly float32, s float32, int cl, int headChar, int float32Digit) {
	x := lx
	y := ly
	n := num
	fd := float32Digit
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
		drawLetterRev(headChar, x+s*LETTER_WIDTH*0.2, y-s*LETTER_WIDTH*0.2, s*0.6, 0, cl)
	}
}

func drawTime(time int, lx float32, y float32, s float32, cl int /* default 0 */) {
	n := time
	if n < 0 {
		n = 0
	}
	var x float32 = lx
	for i := 0; i < 7; i++ {
		if i != 4 {
			drawLetter(n%10, x, y, s, Direction.TO_RIGHT, cl)
			n /= 10
		} else {
			drawLetter(n%6, x, y, s, Direction.TO_RIGHT, cl)
			n /= 6
		}
		if (i&1) == 1 || i == 0 {
			switch i {
			case 3:
				drawLetter(41, x+s*1.16, y, s, Direction.TO_RIGHT, cl)
				break
			case 5:
				drawLetter(40, x+s*1.16, y, s, Direction.TO_RIGHT, cl)
				break
			default:
				break
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
	var x, y, length, size, t, deg float32
	for i := 0; ; i++ {
		deg = int(spData[idx][i][4])
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
		y = y
		deg %= 180
		if c == LINE_COLOR {
			setBoxLine(x, y, size, length, deg)
		} else if c == POLY_COLOR {
			setBoxPoly(x, y, size, length, deg)
		} else {
			setBox(x, y, size, length, deg,
				COLOR_RGB[c][0], COLOR_RGB[c][1], COLOR_RGB[c][2])
		}
	}
}

func setBox(x float32, y float32, width float32, height float32, deg float32, r float32, g float32, b float32) {
	gl.PushMatrix()
	gl.Translatef(x-width/2, y-height/2, 0)
	gl.Rotatef(deg, 0, 0, 1)
	Screen.setColor(r, g, b, 0.5)
	gl.Begin(gl_TRIANGLE_FAN)
	setBoxPart(width, height)
	gl.End()
	Screen.setColor(r, g, b)
	gl.Begin(gl.LINE_LOOP)
	setBoxPart(width, height)
	gl.End()
	gl.PopMatrix()
}

func setBoxLine(x float32, y float32, width flaot, height float32, deg float32) {
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

const spData = [][][]float32{
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
