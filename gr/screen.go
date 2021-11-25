package main

import (
	r "math/rand"

	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

const CAPTION = "Gunroar"

var rand *r.Rand = r.New(r.NewSource(0)) // TODO should the seed be random?
var lineWidthBase float64

/**
 * Initialize an OpenGL and set the caption.
 * Handle a luminous screen and a viewpoint.
 */
type Screen struct {
	*sdl.Screen3D

	luminousScreen     LuminousScreen
	_luminosity        float64
	screenShakeCnt     int
	screenShakeIntense float64
}

func NewScreen() *Screen {
	return &Screen{Screen3D: sdl.NewScreen3D()}
}

func setRandSeed(seed int64) {
	rand = r.New(r.NewSource(seed))
}

func (this *Screen) init() {
	this.Screen3D.InitSDL()

	this.SetCaption(CAPTION)
	gl.LineWidth(1)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.LINE_SMOOTH)
	gl.Disable(gl.TEXTURE_2D)
	gl.Disable(gl.COLOR_MATERIAL)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	sdl.SetClearColor(0, 0, 0, 1)
	if this._luminosity > 0 {
		this.luminousScreen = NewLuminousScreen()
		this.luminousScreen.init(this._luminosity, this.Width(), this.Height())
	} else {
		this.luminousScreen = nil
	}
	this.screenResized()
}

func (this *Screen) close() {
	if this.luminousScreen != nil {
		this.luminousScreen.close()
	}
	this.Screen3D.CloseSDL()
}

func (this *Screen) startRenderToLuminousScreen() bool {
	if this.luminousScreen == nil {
		return false
	}
	this.luminousScreen.startRender()
	return true
}

func (this *Screen) endRenderToLuminousScreen() {
	if this.luminousScreen != nil {
		this.luminousScreen.endRender()
	}
}

func (this *Screen) drawLuminous() {
	if this.luminousScreen != nil {
		this.luminousScreen.draw()
	}
}

func (this *Screen) resized(width, height int) {
	if this.luminousScreen != nil {
		this.luminousScreen.resized(width, height)
	}
	this.Screen3D.Resized(width, height)
}

func (this *Screen) screenResized() {
	this.Screen3D.ScreenResized()
	lw := (this.Width()/640 + this.Height()/480) / 2
	if lw < 1 {
		lw = 1
	} else if lw > 4 {
		lw = 4
	}
	lineWidthBase = float64(lw)
	lineWidth(1)
}

func lineWidth(w int) {
	gl.LineWidth(float32(lineWidthBase) * float32(w))
}

func (this *Screen) clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func viewOrthoFixed() {
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(0, 640, 480, 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()
}

func viewPerspective() {
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.MatrixMode(gl.MODELVIEW)
	gl.PopMatrix()
}

func (this *Screen) setEyepos() {
	var ex, ey, ez float64
	var lx, ly, lz float64
	ez = 13.0
	if this.screenShakeCnt > 0 {
		mx := rand.nextSignedFloat(this.screenShakeIntense * (this.screenShakeCnt + 4))
		my := rand.nextSignedFloat(this.screenShakeIntense * (this.screenShakeCnt + 4))
		ex += mx
		ey += my
		lx += mx
		ly += my
	}
	gluLookAt(ex, ey, ez, lx, ly, lz, 0, 1, 0)
}

func (this *Screen) setScreenShake(cnt int, its float64) {
	this.screenShakeCnt = cnt
	this.screenShakeIntense = its
}

func (this *Screen) move() {
	if this.screenShakeCnt > 0 {
		this.screenShakeCnt--
	}
}

func (this *Screen) luminosity(v float64) float64 {
	this._luminosity = v
	return v
}

func setColorForced(r, g, b, a float64 /* = 1 */) {
	gl.Color4f(r, g, b, a)
}
