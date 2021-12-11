/*
 * $Id: screen.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"fmt"

	"github.com/dragonfax/gunroar_go/glu"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

/**
 * Initialize an OpenGL and set the caption.
 * Handle a luminous screen and a viewpoint.
 */

const CAPTION = "Gunroar"

var screen *Screen

var lineWidthBase float32 = 1
var brightness float32 = 1

type Screen struct {
	luminousScreen     *LuminousScreen
	luminosity         float32
	screenShakeCnt     int
	screenShakeIntense float32
	farPlane           float32
	nearPlane          float32
	width              uint32
	height             uint32
	windowMode         bool
	window             *sdl.Window
	context            sdl.GLContext
}

func NewScreen() *Screen {
	this := new(Screen)
	this.width = 640
	this.height = 480
	return this
}

// called by InitSDL()
func (s *Screen) Init() {
	fmt.Println("initing screen")
	s.setCaption(CAPTION)
	gl.LineWidth(1)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.LINE_SMOOTH)
	gl.Disable(gl.TEXTURE_2D)
	gl.Disable(gl.COLOR_MATERIAL)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	setClearColor(0, 0, 0, 1)
	if s.luminosity > 0 {
		s.luminousScreen = &LuminousScreen{}
		s.luminousScreen.Init(s.luminosity, s.width, s.height)
	} else {
		s.luminousScreen = nil
	}
	s.screenResized()
	checkGLError()
}

func (s *Screen) Close() {
	if s.luminousScreen != nil {
		s.luminousScreen.close()
	}
}

func (s *Screen) startRenderToLuminousScreen() bool {
	if s.luminousScreen == nil {
		return false
	}
	s.luminousScreen.startRender()
	return true
}

func (s *Screen) endRenderToLuminousScreen() {
	if s.luminousScreen != nil {
		s.luminousScreen.endRender()
	}
}

func (s *Screen) drawLuminous() {
	if s.luminousScreen != nil {
		s.luminousScreen.draw()
	}
}

func (s *Screen) resized(width uint32, height uint32) {
	if s.luminousScreen != nil {
		s.luminousScreen.resized(width, height)
	}

	s.width = width
	s.height = height
	s.screenResized()
}

func (s *Screen) screenResized() {
	gl.Viewport(0, 0, int32(s.width), int32(s.height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-float64(s.nearPlane),
		float64(s.nearPlane),
		-float64(s.nearPlane*float32(s.height)/float32(s.width)),
		float64(s.nearPlane*float32(s.height)/float32(s.width)),
		0.1, float64(s.farPlane))
	gl.MatrixMode(gl.MODELVIEW)

	lw := (float32(s.width)/640 + float32(s.height)/480) / 2
	if lw < 1 {
		lw = 1
	} else if lw > 4 {
		lw = 4
	}
	lineWidthBase = lw
	lineWidth(1)
}

func lineWidth(w int) {
	newWidth := lineWidthBase * float32(w)
	if newWidth < 1 {
		panic("line too small")
	}
	gl.LineWidth(newWidth)
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

func (s *Screen) setEyepos() {
	var ex, ey, ez float32
	var lx, ly, lz float32
	ez = 13.0
	if s.screenShakeCnt > 0 {
		mx := nextSignedFloat(s.screenShakeIntense * float32(s.screenShakeCnt+4))
		my := nextSignedFloat(s.screenShakeIntense * float32(s.screenShakeCnt+4))
		ex += mx
		ey += my
		lx += mx
		ly += my
	}
	checkGLError()
	glu.LookAt(float64(ex), float64(ey), float64(ez), float64(lx), float64(ly), float64(lz), 0, 1, 0)
	checkGLError()
}

func (s *Screen) setScreenShake(cnt int, its float32) {
	s.screenShakeCnt = cnt
	s.screenShakeIntense = its
}

func (s *Screen) move() {
	if s.screenShakeCnt > 0 {
		s.screenShakeCnt--
	}
}

func setScreenColorForced(r float32, g float32, b float32, a float32 /* = 1 */) {
	gl.Color4f(r, g, b, a)
}

func (s *Screen) initSDL() {
	fmt.Println("initing SDL")
	s.farPlane = 1000
	s.nearPlane = 0.1
	// Initialize SDL.
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Sprintf(" SDLInitFailedException( Unable to initialize SDL: %s )", err))
	}
	// Create an OpenGL screen.
	var videoFlags uint32
	var window *sdl.Window
	var err error
	//if s.windowMode {
	videoFlags = sdl.WINDOW_OPENGL
	window, err = sdl.CreateWindow("Title", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(s.width), int32(s.height), videoFlags)
	/* } else {
		videoFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_FULLSCREEN_DESKTOP
		window, err = sdl.CreateWindow("Title", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 0, 0, videoFlags)
	} */
	if err != nil {
		panic(fmt.Sprintf("SDLInitFailedException (Unable to create SDL screen: %v", sdl.GetError()))
	}
	s.window = window
	// sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	// sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)
	// sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	s.context, err = s.window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	err = s.window.GLMakeCurrent(s.context)
	if err != nil {
		panic(err)
	}

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	{
		major, err := sdl.GLGetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION)
		if err != nil {
			panic(err)
		}
		checkGLError()
		minor, err := sdl.GLGetAttribute(sdl.GL_CONTEXT_MINOR_VERSION)
		if err != nil {
			panic(err)
		}
		checkGLError()
		fmt.Printf("opengl version %d.%d\n", major, minor)
	}

	gl.Clear(gl.COLOR_BUFFER_BIT)
	checkGLError()

	gl.Viewport(0, 0, int32(s.width), int32(s.height))
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	gl.Clear(gl.COLOR_BUFFER_BIT)
	checkGLError()

	w, h := s.window.GLGetDrawableSize()
	s.width = uint32(w)
	s.height = uint32(h)

	gl.Clear(gl.COLOR_BUFFER_BIT)
	checkGLError()

	s.resized(s.width, s.height)
	_, err = sdl.ShowCursor(sdl.DISABLE)
	if err != nil {
		panic(err)
	}

	gl.Clear(gl.COLOR_BUFFER_BIT)
	checkGLError()

	s.Init()
	checkGLError()

	gl.Clear(gl.COLOR_BUFFER_BIT)
	checkGLError()
}

func (s *Screen) closeSDL() {
	s.Close()
	sdl.ShowCursor(sdl.ENABLE)
}

func (s *Screen) flip() {
	checkGLError()
	s.handleError()
	checkGLError()
	s.window.GLSwap()
	checkGLError()
}

func (s *Screen) clear() {
	checkGLError()

	/*
		r := gl.CheckFramebufferStatus(gl.DRAW_FRAMEBUFFER)
		if r == gl.FRAMEBUFFER_UNDEFINED {
			fmt.Printf("frame buffer %d not defined\n", gl.DRAW_FRAMEBUFFER)
		} else if r != gl.FRAMEBUFFER_COMPLETE {
			panic(fmt.Sprintf("frame buffer %d not ready: %d", gl.DRAW_FRAMEBUFFER, r))
		}

		r = gl.CheckFramebufferStatus(gl.READ_FRAMEBUFFER)
		if r == gl.FRAMEBUFFER_UNDEFINED {
			fmt.Printf("frame buffer %d not defined\n", gl.READ_FRAMEBUFFER)
		} else if r != gl.FRAMEBUFFER_COMPLETE {
			panic(fmt.Sprintf("frame buffer %d not ready: %d", gl.READ_FRAMEBUFFER, r))
		}
	*/

	gl.Clear(gl.COLOR_BUFFER_BIT)
	checkGLError()
}

func (s *Screen) handleError() {
	error := gl.GetError()
	if error == gl.NO_ERROR {
		return
	}
	s.closeSDL()
	panic(fmt.Sprintf("OpenGL error( %v )", error))
}

func (s *Screen) setCaption(name string) {
	s.window.SetTitle(name)
}

func glVertex(v Vector) {
	gl.Vertex3f(v.x, v.y, 0)
}

func glVertex3(v Vector3) {
	gl.Vertex3f(v.x, v.y, v.z)
}

func glTranslate(v Vector) {
	gl.Translatef(v.x, v.y, 0)
}

func glTranslate3(v Vector3) {
	gl.Translatef(v.x, v.y, v.z)
}

func setScreenColor(r float32, g float32, b float32, a float32 /* = 1 */) {
	gl.Color4f(r*brightness, g*brightness, b*brightness, a)
}

func setClearColor(r float32, g float32, b float32, a float32 /*= 1*/) {
	gl.ClearColor(r*brightness, g*brightness, b*brightness, a)
}
