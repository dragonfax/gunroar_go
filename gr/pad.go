/*
 * $Id: pad.d,v 1.2 2005/07/03 07:05:23 kenta Exp $
 *
 * Copyright 2004 Kenta Cho. Some rights reserved.
 */
package gr

import (
	"github.com/veandco/go-sdl2/sdl"
)

const JOYSTICK_AXIS = 16384

/**
 * Joystick and keyboard input.
 */
type Pad struct {
	keys           []uint8
	buttonReversed bool

	stick *sdl.Joystick
	state PadState
}

func (pad *Pad) openJoystick(st *sdl.Joystick /* = null */) *sdl.Joystick {
	if st == nil {
		if sdl.InitSubSystem(sdl.INIT_JOYSTICK) < 0 {
			return nil
		}
		pad.stick = sdl.JoystickOpen(0)
	} else {
		pad.stick = st
	}
	return pad.stick
}

func (pad *Pad) handleEvent(event *sdl.Event) {
	pad.keys = sdl.GetKeyboardState()
}

func (pad *Pad) getState() PadState {
	var x, y int
	pad.state.dir = 0
	if pad.stick != nil {
		x = sdl.JoystickGetAxis(pad.stick, 0)
		y = sdl.JoystickGetAxis(pad.stick, 1)
	}
	if pad.keys[sdl.K_RIGHT] == sdl.PRESSED || pad.keys[sdl.K_KP6] == sdl.PRESSED ||
		pad.keys[sdl.K_d] == sdl.PRESSED || pad.keys[sdl.K_l] == sdl.PRESSED ||
		x > JOYSTICK_AXIS {
		pad.state.dir |= RIGHT
	}
	if pad.keys[sdl.K_LEFT] == sdl.PRESSED || pad.keys[sdl.K_KP4] == sdl.PRESSED ||
		pad.keys[sdl.K_a] == sdl.PRESSED || pad.keys[sdl.K_j] == sdl.PRESSED ||
		x < -JOYSTICK_AXIS {
		pad.state.dir |= LEFT
	}
	if pad.keys[sdl.K_DOWN] == sdl.PRESSED || pad.keys[sdl.K_KP2] == sdl.PRESSED ||
		pad.keys[sdl.K_s] == sdl.PRESSED || pad.keys[sdl.K_k] == sdl.PRESSED ||
		y > JOYSTICK_AXIS {
		pad.state.dir |= DOWN
	}
	if pad.keys[sdl.K_UP] == sdl.PRESSED || pad.keys[sdl.K_KP8] == sdl.PRESSED ||
		pad.keys[sdl.K_w] == sdl.PRESSED || pad.keys[sdl.K_i] == sdl.PRESSED ||
		y < -JOYSTICK_AXIS {
		pad.state.dir |= UP
	}
	pad.state.button = 0
	var btn1, btn2 int
	var leftTrigger, rightTrigger float = 0
	if pad.stick {
		btn1 = sdl.JoystickGetButton(pad.stick, 0) + sdl.JoystickGetButton(pad.stick, 3) +
			sdl.JoystickGetButton(pad.stick, 4) + sdl.JoystickGetButton(pad.stick, 7) +
			sdl.JoystickGetButton(pad.stick, 8) + sdl.JoystickGetButton(pad.stick, 11)
		btn2 = sdl.JoystickGetButton(pad.stick, 1) + sdl.JoystickGetButton(pad.stick, 2) +
			sdl.JoystickGetButton(pad.stick, 5) + sdl.JoystickGetButton(pad.stick, 6) +
			sdl.JoystickGetButton(pad.stick, 9) + sdl.JoystickGetButton(pad.stick, 10)
	}
	if pad.keys[sdl.K_z] == sdl.PRESSED || pad.keys[sdl.K_PERIOD] == sdl.PRESSED ||
		pad.keys[sdl.K_LCTRL] == sdl.PRESSED || pad.keys[sdl.K_RCTRL] == sdl.PRESSED ||
		btn1 {
		if !pad.buttonReversed {
			pad.state.button |= PadState.Button.A
		} else {
			pad.state.button |= PadState.Button.B
		}
	}
	if pad.keys[sdl.K_x] == sdl.PRESSED || pad.keys[sdl.K_SLASH] == sdl.PRESSED ||
		pad.keys[sdl.K_LALT] == sdl.PRESSED || pad.keys[sdl.K_RALT] == sdl.PRESSED ||
		pad.keys[sdl.K_LSHIFT] == sdl.PRESSED || pad.keys[sdl.K_RSHIFT] == sdl.PRESSED ||
		pad.keys[sdl.K_RETURN] == sdl.PRESSED ||
		btn2 {
		if !buttonReversed {
			pad.state.button |= B
		} else {
			pad.state.button |= A
		}
	}
	return pad.state
}

type Dir int

const (
	UP    Dir = 1
	DOWN      = 2
	LEFT      = 4
	RIGHT     = 8
)

type Button int

const (
	A   Button = 16
	B   Button = 32
	ANY Button = 48
)

type PadState struct {
	dir, button int
}
