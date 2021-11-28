package letter

import (
	"math"

	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Letters.
 */
var displayList *sdl.DisplayList

const LETTER_WIDTH = 2.1
const LETTER_HEIGHT = 3.0
const LINE_COLOR = 2
const POLY_COLOR = 3
const COLOR_NUM = 4

var COLOR_RGB = [][]float64{{1, 1, 1}, {0.9, 0.7, 0.5}}

const LETTER_NUM = 44
const DISPLAY_LIST_NUM = LETTER_NUM * COLOR_NUM

type Letter struct {
}

func LetterInit() {
	displayList = sdl.NewDisplayList(DISPLAY_LIST_NUM)
	displayList.ResetList()
	for j := 0; j < COLOR_NUM; j++ {
		for i := 0; i < LETTER_NUM; i++ {
			displayList.NewList()
			setLetter(i, j)
			displayList.EndList()
		}
	}
}

func getWidth(n int, s float64) float64 {
	return float64(n) * s * LETTER_WIDTH
}

func getHeight(s float64) float64 {
	return s * LETTER_HEIGHT
}

func drawLetterAsIs(n, c int) {
	displayList.Call(n + c*LETTER_NUM)
}

func drawLetter(n int, x, y, s, d float64, c int) {
	gl.PushMatrix()
	gl.Translated(x, y, 0)
	gl.Scaled(s, s, s)
	gl.Rotated(d, 0, 0, 1)
	displayList.Call(n + c*LETTER_NUM)
	gl.PopMatrix()
}

func drawLetterRev(n int, x, y, s, d float64, c int) {
	gl.PushMatrix()
	gl.Translated(x, y, 0)
	gl.Scaled(s, -s, s)
	gl.Rotated(d, 0, 0, 1)
	displayList.Call(n + c*LETTER_NUM)
	gl.PopMatrix()
}

type Direction int

const TO_RIGHT Direction = 0
const TO_DOWN Direction = 1
const TO_LEFT Direction = 2
const TO_UP Direction = 3

func convertCharToInt(c rune) int {
	var idx int
	if c >= '0' && c <= '9' {
		idx = int(c) - int('0')
	} else if c >= 'A' && c <= 'Z' {
		idx = int(c) - int('A') + 10
	} else if c >= 'a' && c <= 'z' {
		idx = int(c) - int('a') + 10
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

func drawString(str string, lx, y, s float64,
	d Direction /* = TO_RIGHT */, cl int, /* = 0 */
	rev bool /* = false */, od float64 /* = 0 */) {
	lx += LETTER_WIDTH * s / 2
	y += LETTER_HEIGHT * s / 2
	x := lx
	var idx int
	var ld float64
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
	for _, c := range str {
		if c != ' ' {
			idx = convertCharToInt(c)
			if rev {
				drawLetterRev(idx, x, y, s, ld, cl)
			} else {
				drawLetter(idx, x, y, s, ld, cl)
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
			x += math.Cos(ld*math.Pi/180) * s * LETTER_WIDTH
			y += math.Sin(ld*math.Pi/180) * s * LETTER_WIDTH
		}
	}
}

func drawNum(num int, lx, y, s float64,
	cl /* = 0 */, dg int, /* = 0 */
	headChar /* = -1 */, floatDigit int /* = -1 */) {
	lx += LETTER_WIDTH * s / 2
	y += LETTER_HEIGHT * s / 2
	n := num
	x := lx
	ld := 0.0
	digit := dg
	fd := floatDigit
	for {
		if fd <= 0 {
			drawLetter(n%10, x, y, s, ld, cl)
			x -= s * LETTER_WIDTH
		} else {
			drawLetter(n%10, x, y+s*LETTER_WIDTH*0.25, s*0.5, ld, cl)
			x -= s * LETTER_WIDTH * 0.5
		}
		n /= 10
		digit--
		fd--
		if n <= 0 && digit <= 0 && fd < 0 {
			break
		}
		if fd == 0 {
			drawLetter(36, x, y+s*LETTER_WIDTH*0.25, s*0.5, ld, cl)
			x -= s * LETTER_WIDTH * 0.5
		}
	}
	if headChar >= 0 {
		drawLetter(headChar, x+s*LETTER_WIDTH*0.2, y+s*LETTER_WIDTH*0.2,
			s*0.6, ld, cl)
	}
}

func DrawNumSign(num int, lx, ly, s float64, cl int, /* = 0 */
	headChar int /* = -1 */, floatDigit int /* = -1 */) {
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
		drawLetterRev(headChar, x+s*LETTER_WIDTH*0.2, y-s*LETTER_WIDTH*0.2,
			s*0.6, 0, cl)
	}
}

func drawTime(time int, lx, y, s float64, cl int /* = 0 */) {
	n := time
	if n < 0 {
		n = 0
	}
	x := lx
	for i := 0; i < 7; i++ {
		if i != 4 {
			drawLetter(n%10, x, y, s, float64(TO_RIGHT), cl)
			n /= 10
		} else {
			drawLetter(n%6, x, y, s, float64(TO_RIGHT), cl)
			n /= 6
		}
		if (i&1) == 1 || i == 0 {
			switch i {
			case 3:
				drawLetter(41, x+s*1.16, y, s, float64(TO_RIGHT), cl)
			case 5:
				drawLetter(40, x+s*1.16, y, s, float64(TO_RIGHT), cl)
			default:
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

func setLetter(idx, c int) {
	var x, y, length, size, deg float64
	for i := 0; ; i++ {
		deg = spData[idx][i][4]
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
		deg = math.Mod(deg, 180)
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

func setBox(x, y, width, height, deg, r, g, b float64) {
	gl.PushMatrix()
	gl.Translated(x-width/2, y-height/2, 0)
	gl.Rotated(deg, 0, 0, 1)
	sdl.SetColor(r, g, b, 0.5)
	gl.Begin(gl.TRIANGLE_FAN)
	setBoxPart(width, height)
	gl.End()
	sdl.SetColor(r, g, b, 1)
	gl.Begin(gl.LINE_LOOP)
	setBoxPart(width, height)
	gl.End()
	gl.PopMatrix()
}

func setBoxLine(x, y, width, height, deg float64) {
	gl.PushMatrix()
	gl.Translated(x-width/2, y-height/2, 0)
	gl.Rotated(deg, 0, 0, 1)
	gl.Begin(gl.LINE_LOOP)
	setBoxPart(width, height)
	gl.End()
	gl.PopMatrix()
}

func setBoxPoly(x, y, width, height, deg float64) {
	gl.PushMatrix()
	gl.Translated(x-width/2, y-height/2, 0)
	gl.Rotated(deg, 0, 0, 1)
	gl.Begin(gl.TRIANGLE_FAN)
	setBoxPart(width, height)
	gl.End()
	gl.PopMatrix()
}

func setBoxPart(width, height float64) {
	gl.Vertex3d(-width/2, 0, 0)
	gl.Vertex3d(-width/3*1, -height/2, 0)
	gl.Vertex3d(width/3*1, -height/2, 0)
	gl.Vertex3d(width/2, 0, 0)
	gl.Vertex3d(width/3*1, height/2, 0)
	gl.Vertex3d(-width/3*1, height/2, 0)
}

var spData = [][][5]float64{{
	{0, 1.15, 0.65, 0.3, 0},
	{-0.6, 0.55, 0.65, 0.3, 90}, {0.6, 0.55, 0.65, 0.3, 90},
	{-0.6, -0.55, 0.65, 0.3, 90}, {0.6, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{0.5, 0.55, 0.65, 0.3, 90},
	{0.5, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{0.65, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{0.65, 0.55, 0.65, 0.3, 90},
	{0.65, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, { //A
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{-0.18, 1.15, 0.45, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.45, 0.55, 0.65, 0.3, 90},
	{-0.18, 0, 0.45, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90},
	{-0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{-0.15, 1.15, 0.45, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.45, 0.45, 0.65, 0.3, 90},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, { //F
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90},
	{0.05, 0, 0.3, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0, 0.55, 0.65, 0.3, 90},
	{0, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0.65, 0.55, 0.65, 0.3, 90},
	{0.65, -0.55, 0.65, 0.3, 90}, {-0.7, -0.7, 0.3, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, { //K
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.4, 0.55, 0.65, 0.3, 100},
	{-0.25, 0, 0.45, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.6, -0.55, 0.65, 0.3, 80},
	{0, 0, 0, 0, 99999},
}, {
	{-0.65, 0.55, 0.65, 0.3, 90},
	{-0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{-0.5, 1.15, 0.3, 0.3, 0}, {0.1, 1.15, 0.3, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, 0.55, 0.65, 0.3, 90},
	{0, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, { //P
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0.05, -0.55, 0.45, 0.3, 60},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{-0.2, 0, 0.45, 0.3, 0},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.45, -0.55, 0.65, 0.3, 80},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{-0.65, 0.55, 0.65, 0.3, 90},
	{0, 0, 0.65, 0.3, 0},
	{0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{-0.5, 1.15, 0.55, 0.3, 0}, {0.5, 1.15, 0.55, 0.3, 0},
	{0.1, 0.55, 0.65, 0.3, 90},
	{0.1, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, { //U
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{-0.5, -0.55, 0.65, 0.3, 90}, {0.5, -0.55, 0.65, 0.3, 90},
	{-0.1, -1.15, 0.45, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, {
	{-0.65, 0.55, 0.65, 0.3, 90}, {0.65, 0.55, 0.65, 0.3, 90},
	{-0.65, -0.55, 0.65, 0.3, 90}, {0.65, -0.55, 0.65, 0.3, 90},
	{-0.5, -1.15, 0.3, 0.3, 0}, {0.1, -1.15, 0.3, 0.3, 0},
	{0, 0.55, 0.65, 0.3, 90},
	{0, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{-0.4, 0.6, 0.85, 0.3, 360 - 120},
	{0.4, 0.6, 0.85, 0.3, 360 - 60},
	{-0.4, -0.6, 0.85, 0.3, 360 - 240},
	{0.4, -0.6, 0.85, 0.3, 360 - 300},
	{0, 0, 0, 0, 99999},
}, {
	{-0.4, 0.6, 0.85, 0.3, 360 - 120},
	{0.4, 0.6, 0.85, 0.3, 360 - 60},
	{-0.1, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, {
	{0, 1.15, 0.65, 0.3, 0},
	{0.3, 0.4, 0.65, 0.3, 120},
	{-0.3, -0.4, 0.65, 0.3, 120},
	{0, -1.15, 0.65, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, { //.
	{0, -1.15, 0.3, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, { //_
	{0, -1.15, 0.8, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, { //-
	{0, 0, 0.9, 0.3, 0},
	{0, 0, 0, 0, 99999},
}, { //+
	{-0.5, 0, 0.45, 0.3, 0}, {0.45, 0, 0.45, 0.3, 0},
	{0.1, 0.55, 0.65, 0.3, 90},
	{0.1, -0.55, 0.65, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, { //'
	{0, 1.0, 0.4, 0.2, 90},
	{0, 0, 0, 0, 99999},
}, { //''
	{-0.19, 1.0, 0.4, 0.2, 90},
	{0.2, 1.0, 0.4, 0.2, 90},
	{0, 0, 0, 0, 99999},
}, { //!
	{0.56, 0.25, 1.1, 0.3, 90},
	{0, -1.0, 0.3, 0.3, 90},
	{0, 0, 0, 0, 99999},
}, { // /
	{0.8, 0, 1.75, 0.3, 120},
	{0, 0, 0, 0, 99999},
}}
