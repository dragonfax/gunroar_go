/*
 * $Id: screen.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

import (
	"math/rand"
)

/**
 * Initialize an OpenGL and set the caption.
 * Handle a luminous screen and a viewpoint.
 */

const CAPTION = "Gunroar"

var lineWidthBase float32
var brightness float32 = 1

type Screen struct {
	luminousScreen     LuminousScreen
	luminosity         float32
	screenShakeCnt     int
	screenShakeIntense float32
	farPlane           float32
	nearPlane          float32
	width              int
	height             int
	windowMode         bool
}

func (s *Screen) Init() {
	s.farPlane = 1000
	s.nearPlane = 0.1
	s.width = 640
	s.height = 480
	s.setCaption(CAPTION)
	gl.LineWidth(1)
	gl.BlendFunc(gl.SRC_ALPHA, GL_ONE)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.LINE_SMOOTH)
	gl.Disable(gl.TEXTURE_2D)
	gl.Disable(gl.COLOR_MATERIAL)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	s.setClearColor(0, 0, 0, 1)
	if s.luminosity > 0 {
		s.luminousScreen = NewLuminousScreen()
		s.luminousScreen.Init(s.luminosity, s.width, s.height)
	} else {
		s.luminousScreen = nil
	}
	s.screenResized()
}

func (s *Screen) Close() {
	if s.luminousScreen != nil {
		s.luminousScreen.Close()
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

func (s *Screen) resized(width int, height int) {
	if s.luminousScreen != nil {
		s.luminousScreen.resized(width, height)
	}

	s.width = width
	s.height = height
	screenResized()
}

func (s *Screen) screenResized() {
	gl.Viewport(0, 0, s.width, s.height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-s.nearPlane,
		s.nearPlane,
		-s.nearPlane*float32(s.height)/float32(s.width),
		s.nearPlane*float32(s.height)/float32(s.width),
		0.1, s.farPlane)
	gl.MatrixMode(gl.MODELVIEW)

	lw := (float32(width)/640 + float32(height)/480) / 2
	if lw < 1 {
		lw = 1
	} else if lw > 4 {
		lw = 4
	}
	lineWidthBase = lw
	lineWidth(1)
}

func lineWidth(w int) {
	gl.LineWidth(int(lineWidthBase * w))
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
		mx = rand.Float32 * (s.screenShakeIntense * (s.screenShakeCnt + 4))
		my = rand.Float32 * (s.screenShakeIntense * (s.screenShakeCnt + 4))
		ex += mx
		ey += my
		lx += mx
		ly += my
	}
	glu.LookAt(ex, ey, ez, lx, ly, lz, 0, 1, 0)
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

func (s *Screen) luminosity(v float32) float32 {
	s.luminosity = v
	return s.luminosity
}

func setColorForced(r float, g float, b float, a float /* = 1 */) {
	gl.Color4f(r, g, b, a)
}

func (s *Screen) initSDL() {
	// Initialize SDL.
	if sdl.Init(sdl.INIT_VIDEO) < 0 {
		panic(" SDLInitFailedException( Unable to initialize SDL: " + sdl.GetError())
	}
	// Create an OpenGL screen.
	var videoFlags uint32
	if s.windowMode {
		videoFlags = sdl.OPENGL | sdl.RESIZABLE
	} else {
		videoFlags = sdl.OPENGL | sdl.FULLSCREEN
	}
	if sdl.SetVideoMode(s.width, s.height, 0, videoFlags) == nil {
		panic("SDLInitFailedException (Unable to create SDL screen: " + sdl.GetError())
	}
	gl.Viewport(0, 0, s.width, s.height)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	s.resized(s.width, s.height)
	sdl.ShowCursor(sdl.DISABLE)
	s.Init()
}

func (s *Screen) closeSDL() {
	s.close()
	sdl.ShowCursor(sdl.ENABLE)
}

func (s *Screen) flip() {
	s.handleError()
	gl.SwapBuffers()
}

func (s *Screen) clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (s *Screen) handleError() {
	error := glGetError()
	if error == gl.NO_ERROR {
		return
	}
	s.closeSDL()
	panic("OpenGL error(" + error + ")")
}

func (s *Screen) setCaption(name string) {
	sdl.WM_SetCaption(name, nil)
}

func glVertex(v Vector) {
	gl.Vertex3f(v.X(), v.Y(), 0)
}

func glVertex3(v Vector3) {
	gl.Vertex3f(v.X(), v.Y(), v.Z())
}

func glTranslate(v Vector) {
	gl.Translatef(v.X(), v.Y(), 0)
}

func glTranslate3(v Vector3) {
	gl.Translatef(v.X(), v.Y(), v.Z())
}

func setColor(r float32, g float32, b float32, a float32 /* = 1 */) {
	gl.Color4f(r*brightness, g*brightness, b*brightness, a)
}

func setClearColor(r float32, g float32, b float32, a float32 /*= 1*/) {
	gl.ClearColor(r*brightness, g*brightness, b*brightness, a)
}
