/*
 * $Id: mouse.d,v 1.1 2005/09/11 00:47:41 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package sdl

/**
 * Mouse input.
 */
type Mouse struct {
  screen SizableScreen
  state MouseState
}

func (m *Mouse) getState() MouseState {
	mx, my, btn := sdl.GetMouseState()
	m.state.X = mx
	m.state.Y = my
	m.state.Button = NONE
	if btn & sdl.Button(LEFT) {
		state.Button |= LEFT
	}
	if btn & sdl.Button(RIGHT) {
		state.Button |= RIGHT
	}
	adjustPos(state)
	return state
}

type Button int 

const (
	NONE Button = 0,
	LEFT Button = 1,
	RIGHT Button = 2
)


type MouseState struct {
  X, Y float32
  Button Button
}
