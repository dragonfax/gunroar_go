/*
 * $Id: pad.d,v 1.2 2005/07/03 07:05:23 kenta Exp $
 *
 * Copyright 2004 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/veandco/go-sdl2/sdl"

const JOYSTICK_AXIS = 16384

var pad *Pad

/**
 * Joystick and keyboard input.
 */
type Pad struct {
	keys           []uint8
	buttonReversed bool

	stick *sdl.Joystick
	state PadState
}

func NewPad() *Pad {
	return new(Pad)
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

func (pad *Pad) update() {
	pad.keys = sdl.GetKeyboardState()
}

func (pad *Pad) getNullState() PadState {
	return PadState{}
}

func (pad *Pad) getState() PadState {
	var x, y int16
	pad.state.dir = 0
	if pad.stick != nil {
		x = pad.stick.GetAxis(0)
		y = pad.stick.GetAxis(1)
	}
	if pad.keys[sdl.SCANCODE_RIGHT] == sdl.PRESSED || pad.keys[sdl.SCANCODE_KP_6] == sdl.PRESSED ||
		pad.keys[sdl.SCANCODE_D] == sdl.PRESSED || pad.keys[sdl.SCANCODE_L] == sdl.PRESSED ||
		x > JOYSTICK_AXIS {
		pad.state.dir |= PadDirRIGHT
	}
	if pad.keys[sdl.SCANCODE_LEFT] == sdl.PRESSED || pad.keys[sdl.SCANCODE_KP_4] == sdl.PRESSED ||
		pad.keys[sdl.SCANCODE_A] == sdl.PRESSED || pad.keys[sdl.SCANCODE_J] == sdl.PRESSED ||
		x < -JOYSTICK_AXIS {
		pad.state.dir |= PadDirLEFT
	}
	if pad.keys[sdl.SCANCODE_DOWN] == sdl.PRESSED || pad.keys[sdl.SCANCODE_KP_2] == sdl.PRESSED ||
		pad.keys[sdl.SCANCODE_S] == sdl.PRESSED || pad.keys[sdl.SCANCODE_K] == sdl.PRESSED ||
		y > JOYSTICK_AXIS {
		pad.state.dir |= PadDirDOWN
	}
	if pad.keys[sdl.SCANCODE_UP] == sdl.PRESSED || pad.keys[sdl.SCANCODE_KP_8] == sdl.PRESSED ||
		pad.keys[sdl.SCANCODE_W] == sdl.PRESSED || pad.keys[sdl.SCANCODE_I] == sdl.PRESSED ||
		y < -JOYSTICK_AXIS {
		pad.state.dir |= PadDirUP
	}
	pad.state.button = 0
	var btn1, btn2 byte
	// var leftTrigger float32 = 0
	// var rightTrigger float32 = 0
	if pad.stick != nil {
		btn1 = pad.stick.GetButton(0) + pad.stick.GetButton(3) +
			pad.stick.GetButton(4) + pad.stick.GetButton(7) +
			pad.stick.GetButton(8) + pad.stick.GetButton(11)
		btn2 = pad.stick.GetButton(1) + pad.stick.GetButton(2) +
			pad.stick.GetButton(5) + pad.stick.GetButton(6) +
			pad.stick.GetButton(9) + pad.stick.GetButton(10)
	}
	if pad.keys[sdl.SCANCODE_Z] == sdl.PRESSED || pad.keys[sdl.SCANCODE_PERIOD] == sdl.PRESSED ||
		pad.keys[sdl.SCANCODE_LCTRL] == sdl.PRESSED || pad.keys[sdl.SCANCODE_RCTRL] == sdl.PRESSED ||
		btn1 != 0 {
		if !pad.buttonReversed {
			pad.state.button |= PadButtonA
		} else {
			pad.state.button |= PadButtonB
		}
	}
	if pad.keys[sdl.SCANCODE_X] == sdl.PRESSED || pad.keys[sdl.SCANCODE_SLASH] == sdl.PRESSED ||
		pad.keys[sdl.SCANCODE_LALT] == sdl.PRESSED || pad.keys[sdl.SCANCODE_RALT] == sdl.PRESSED ||
		pad.keys[sdl.SCANCODE_LSHIFT] == sdl.PRESSED || pad.keys[sdl.SCANCODE_RSHIFT] == sdl.PRESSED ||
		pad.keys[sdl.SCANCODE_RETURN] == sdl.PRESSED ||
		btn2 != 0 {
		if !pad.buttonReversed {
			pad.state.button |= PadButtonB
		} else {
			pad.state.button |= PadButtonA
		}
	}
	return pad.state
}

type PadDir int

const (
	PadDirUP    PadDir = 1
	PadDirDOWN         = 2
	PadDirLEFT         = 4
	PadDirRIGHT        = 8
)

type PadButton int

const (
	PadButtonA   PadButton = 16
	PadButtonB   PadButton = 32
	PadButtonANY PadButton = 48
)

type PadState struct {
	dir    PadDir
	button PadButton
}
