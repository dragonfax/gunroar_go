package sdl

import (
	"github.com/dragonfax/gunroar/gr/sdl/screen"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var _ screen.Screen = &Screen3D{}
var _ screen.SizableScreen = &Screen3D{}

var _brightness = 1.0

var window *sdl.Window

/**
 * SDL screen handler(3D, OpenGL).
 */
type Screen3D struct {
	_farPlane   float64
	_nearPlane  float64
	_width      int
	_height     int
	_windowMode bool
}

func NewScreen3D() *Screen3D {
	this := &Screen3D{
		_farPlane:  1000,
		_nearPlane: 0.1,
		_width:     640,
		_height:    480,
	}
	return this
}

func (this *Screen3D) InitSDL() {

	// Initialize SDL.
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic("Unable to initialize SDL: " + err.Error())
	}

	// Create an OpenGL screen.
	var videoFlags uint32
	if this._windowMode {
		videoFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_RESIZABLE
	} else {
		videoFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_FULLSCREEN
	}

	window, err = sdl.CreateWindow("gunroar", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(this._width), int32(this._height), sdl.WINDOW_SHOWN|videoFlags)
	if err != nil {
		panic(err)
	}

	gl.Viewport(0, 0, this.Width(), this.Height())
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	this.resized(this._width, this._height)
	sdl.ShowCursor(sdl.DISABLE)
}

// Reset a viewport when the screen is resized.
func (this *Screen3D) ScreenResized() {
	gl.Viewport(0, 0, this._width, this._height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-this._nearPlane,
		this._nearPlane,
		-this._nearPlane*float64(this._height)/float64(this._width),
		this._nearPlane*float64(this._height)/float64(this._width),
		0.1, this._farPlane)
	gl.MatrixMode(gl.MODELVIEW)
}

func (this *Screen3D) Resized(w, h int) {
	this._width = w
	this._height = h
	this.screenResized()
}

func (this *Screen3D) CloseSDL() {
	sdl.ShowCursor(sdl.ENABLE)
}

func (this *Screen3D) Flip() {
	this.handleError()
	gl.SwapBuffers()
}

func (this *Screen3D) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (this *Screen3D) handleError() {
	err := gl.GetError()
	if err == gl.NO_ERROR {
		return
	}
	this.CloseSDL()
	panic("OpenGL error(" + err.Error() + ")")
}

func (this *Screen3D) SetCaption(name string) {
	window.SetTitle(name)
}

func (this *Screen3D) SetWindowMode(v bool) bool {
	this._windowMode = v
	return v
}

func (this *Screen3D) WindowMode() bool {
	return this._windowMode
}

func (this *Screen3D) SetWidth(v int) int {
	this._width = v
	return v
}

func (this *Screen3D) Width() int {
	return this._width
}

func (this *Screen3D) SetHeight(v int) int {
	this._height = v
	return v
}

func (this *Screen3D) Height() int {
	return this._height
}

func glVertex(v vector.Vector) {
	gl.Vertex3f(v.X, v.Y, 0)
}

func glVertex3(v vector.Vector3) {
	gl.Vertex3f(v.X, v.Y, v.Z)
}

func glTranslate(v vector.Vector) {
	gl.Translatef(v.X, v.Y, 0)
}

func glTranslate3(v vector.Vector3) {
	gl.Translatef(v.X, v.Y, v.Z)
}

func setColor(r, g, b, a float64 /* = 1 */) {
	gl.Color4f(r*_brightness, g*_brightness, b*_brightness, a)
}

func SetClearColor(r, g, b, a float64 /* = 1 */) {
	gl.ClearColor(r*_brightness, g*_brightness, b*_brightness, a)
}

func brightness(v float64) float64 {
	_brightness = v
	return v
}
