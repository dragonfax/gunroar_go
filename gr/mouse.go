/*
 * $Id: mouse.d,v 1.1 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

const MOUSE_SCREEN_MAPPING_RATIO_X = 26.0
const MOUSE_SCREEN_MAPPING_RATIO_Y = 19.5

func (m *Mouse) adjustPos(ms *MouseState) {
	ms.x = (ms.x - screen.width/2) * MOUSE_SCREEN_MAPPING_RATIO_X / screen.width
	ms.y = -(ms.y - screen.height/2) * MOUSE_SCREEN_MAPPING_RATIO_Y / screen.height
}

/**
 * Mouse input.
 */
type Mouse struct {
	screen SizableScreen
	state  MouseState
}

func (m *Mouse) getState() MouseState {
	mx, my, btn := sdl.GetMouseState()
	m.state.X = mx
	m.state.Y = my
	m.state.Button = MouseButtonNONE
	if btn & sdl.Button(MouseButtonLEFT) {
		state.Button |= MouseButtonLEFT
	}
	if btn & sdl.Button(MouseButtonRIGHT) {
		state.Button |= MouseButtonRIGHT
	}
	adjustPos(state)
	return state
}

type MouseButton int

const (
	MouseButtonNONE  MouseButton = 0
	MouseButtonLEFT              = 1
	MouseButtonRIGHT             = 2
)

type MouseState struct {
	X, Y   float32
	Button MouseButton
}
