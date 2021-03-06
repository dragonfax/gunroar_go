/*
 * $Id: luminous.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"image"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

/**
 * Luminous effect texture.
 */

const TEXTURE_SIZE_MIN = 0.02
const TEXTURE_SIZE_MAX = 0.98
const LUMINOUS_TEXTURE_WIDTH_MAX = 64
const LUMINOUS_TEXTURE_HEIGHT_MAX = 64

type LuminousScreen struct {
	luminousTexture           uint32
	td                        []uint8
	luminousTextureWidth      uint32
	luminousTextureHeight     uint32
	screenWidth, screenHeight uint32
	luminosity                float32
}

func (ls *LuminousScreen) Init(luminosity float32, width uint32, height uint32) {
	ls.luminousTextureWidth = LUMINOUS_TEXTURE_WIDTH_MAX
	ls.luminousTextureHeight = LUMINOUS_TEXTURE_HEIGHT_MAX
	ls.td = image.NewRGBA(image.Rect(0, 0, int(ls.luminousTextureHeight), int(ls.luminousTextureWidth))).Pix
	ls.makeLuminousTexture()
	ls.luminosity = luminosity
	ls.resized(width, height)
}

func (ls *LuminousScreen) makeLuminousTexture() {
	data := ls.td
	gl.GenTextures(1, &ls.luminousTexture)
	gl.BindTexture(gl.TEXTURE_2D, ls.luminousTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, 4, int32(ls.luminousTextureWidth), int32(ls.luminousTextureHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&data[0]))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
}

func (ls *LuminousScreen) resized(width uint32, height uint32) {
	ls.screenWidth = width
	ls.screenHeight = height
}

func (ls *LuminousScreen) close() {
	gl.DeleteTextures(1, &ls.luminousTexture)
}

func (ls *LuminousScreen) startRender() {
	gl.Viewport(0, 0, int32(ls.luminousTextureWidth), int32(ls.luminousTextureHeight))
}

func (ls *LuminousScreen) endRender() {
	gl.BindTexture(gl.TEXTURE_2D, ls.luminousTexture)
	gl.CopyTexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 0, 0, int32(ls.luminousTextureWidth), int32(ls.luminousTextureHeight), 0)
	gl.Viewport(0, 0, int32(ls.screenWidth), int32(ls.screenHeight))
}

func (ls *LuminousScreen) viewOrtho() {
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(0, float64(ls.screenWidth), float64(ls.screenHeight), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()
}

func (ls *LuminousScreen) viewPerspective() {
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.MatrixMode(gl.MODELVIEW)
	gl.PopMatrix()
}

var lmOfs = [][]float32{[]float32{-2, -1}, []float32{2, 1}}

const lmOfsBs = 3

func (ls *LuminousScreen) draw() {
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, ls.luminousTexture)
	ls.viewOrtho()
	gl.Color4f(1, 0.8, 0.9, ls.luminosity)
	gl.Begin(gl.QUADS)
	for i := 0; i < 2; i++ {
		gl.TexCoord2f(TEXTURE_SIZE_MIN, TEXTURE_SIZE_MAX)
		gl.Vertex2f(0+lmOfs[i][0]*lmOfsBs, 0+lmOfs[i][1]*lmOfsBs)
		gl.TexCoord2f(TEXTURE_SIZE_MIN, TEXTURE_SIZE_MIN)
		gl.Vertex2f(0+lmOfs[i][0]*lmOfsBs, float32(ls.screenHeight)+lmOfs[i][1]*lmOfsBs)
		gl.TexCoord2f(TEXTURE_SIZE_MAX, TEXTURE_SIZE_MIN)
		gl.Vertex2f(float32(ls.screenWidth)+lmOfs[i][0]*lmOfsBs, float32(ls.screenHeight)+lmOfs[i][0]*lmOfsBs)
		gl.TexCoord2f(TEXTURE_SIZE_MAX, TEXTURE_SIZE_MAX)
		gl.Vertex2f(float32(ls.screenWidth)+lmOfs[i][0]*lmOfsBs, 0+lmOfs[i][0]*lmOfsBs)
	}
	gl.End()
	viewPerspective()
	gl.Disable(gl.TEXTURE_2D)
}

type LuminousActor interface {
	drawLuminous()
}
