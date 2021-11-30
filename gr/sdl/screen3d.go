package sdl

import (
	"fmt"

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

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG, gl.TRUE)

	// Create an OpenGL screen.
	var videoFlags uint32
	//if this._windowMode {
	videoFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_RESIZABLE
	//} else {
	// videoFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_FULLSCREEN
	//}

	window, err = sdl.CreateWindow("gunroar", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(this._width), int32(this._height), sdl.WINDOW_SHOWN|videoFlags)
	if err != nil {
		panic(err)
	}

	major, err := sdl.GLGetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION)
	if err != nil {
		panic(err)
	}
	minor, err := sdl.GLGetAttribute(sdl.GL_CONTEXT_MINOR_VERSION)
	if err != nil {
		panic(err)
	}
	fmt.Printf("opengl version %d.%d\n", major, minor)

	_, err = window.GLCreateContext()
	if err != nil {
		panic(err)
	}

	gl.Init()

	gl.Viewport(0, 0, int32(this.Width()), int32(this.Height()))
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	this.Resized(this._width, this._height)
	sdl.ShowCursor(sdl.DISABLE)
}

// Reset a viewport when the screen is resized.
func (this *Screen3D) ScreenResized() {
	gl.Viewport(0, 0, int32(this._width), int32(this._height))
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
	this.ScreenResized()
}

func (this *Screen3D) CloseSDL() {
	sdl.ShowCursor(sdl.ENABLE)
}

func (this *Screen3D) Flip() {
	this.handleError()
	window.GLSwap() // NOTE watch out for macos special issues.
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
	panic("error from open gl")
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

func GlVertex(v vector.Vector) {
	gl.Vertex3d(v.X, v.Y, 0)
}

func GlVertex3(v vector.Vector3) {
	gl.Vertex3d(v.X, v.Y, v.Z)
}

func GlTranslate(v vector.Vector) {
	gl.Translated(v.X, v.Y, 0)
}

func GlTranslate3(v vector.Vector3) {
	gl.Translated(v.X, v.Y, v.Z)
}

func SetColor(r, g, b, a float64 /* = 1 */) {
	gl.Color4d(r*_brightness, g*_brightness, b*_brightness, a)
}

func SetClearColor(r, g, b, a float64 /* = 1 */) {
	gl.ClearColor(float32(r*_brightness), float32(g*_brightness), float32(b*_brightness), float32(a))
}

func Brightness(v float64) float64 {
	_brightness = v
	return v
}
