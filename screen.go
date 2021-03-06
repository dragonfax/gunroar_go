/*
 * $Id: screen.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"fmt"

	"github.com/dragonfax/glu"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

/**
 * Initialize an OpenGL and set the caption.
 * Handle a luminous screen and a viewpoint.
 */

const CAPTION = "Gunroar"

var screen *Screen

var lineWidthBase float32
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
	gl.LineWidth(floor32(lineWidthBase * float32(w)))
}

func (s *Screen) Clear() {
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
	glu.LookAt(float64(ex), float64(ey), float64(ez), float64(lx), float64(ly), float64(lz), 0, 1, 0)
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
	if s.windowMode {
		videoFlags = sdl.WINDOW_OPENGL
		window, err = sdl.CreateWindow("Title", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(s.width), int32(s.height), videoFlags)
	} else {
		videoFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_FULLSCREEN_DESKTOP
		window, err = sdl.CreateWindow("Title", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 0, 0, videoFlags)
	}
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
	err = gl.Init()
	if err != nil {
		panic(err)
	}
	gl.Viewport(0, 0, int32(s.width), int32(s.height))
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	surface, err := s.window.GetSurface()
	if err != nil {
		panic(err)
	}
	s.width = uint32(surface.W)
	s.height = uint32(surface.H)
	s.resized(s.width, s.height)
	sdl.ShowCursor(sdl.DISABLE)
	s.Init()
}

func (s *Screen) closeSDL() {
	s.Close()
	sdl.ShowCursor(sdl.ENABLE)
}

func (s *Screen) flip() {
	// s.handleError()
	s.window.GLSwap()
}

func (s *Screen) clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
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
