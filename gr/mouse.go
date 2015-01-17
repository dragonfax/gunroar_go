/*
 * $Id: mouse.d,v 1.1 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const MOUSE_SCREEN_MAPPING_RATIO_X = 26.0
const MOUSE_SCREEN_MAPPING_RATIO_Y = 19.5

func (m *Mouse) adjustPos(ms *MouseState) {
	ms.x = (ms.x - float32(m.screen.width)/2) * MOUSE_SCREEN_MAPPING_RATIO_X / float32(m.screen.width)
	ms.y = -(ms.y - float32(m.screen.height)/2) * MOUSE_SCREEN_MAPPING_RATIO_Y / float32(m.screen.height)
}

/**
 * Mouse input.
 */
type Mouse struct {
	screen Screen
	state  MouseState
}

func NewMouse() *Mouse {
	return new(Mouse)
}

func (m *Mouse) getNullState() MouseState {
	return MouseState{}
}

func (m *Mouse) getState() MouseState {
	mx, my, btn := sdl.GetMouseState()
	m.state.x = float32(mx)
	m.state.y = float32(my)
	m.state.button = MouseButtonNONE
	if btn&sdl.Button(MouseButtonLEFT) != 0 {
		m.state.button |= MouseButtonLEFT
	}
	if btn&sdl.Button(MouseButtonRIGHT) != 0 {
		m.state.button |= MouseButtonRIGHT
	}
	m.adjustPos(&m.state)
	return m.state
}

type MouseButton int

const (
	MouseButtonNONE  MouseButton = 0
	MouseButtonLEFT              = 1
	MouseButtonRIGHT             = 2
)

type MouseState struct {
	x, y   float32
	button MouseButton
}
