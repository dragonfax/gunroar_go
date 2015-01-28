/*
 * $Id: twinstick.d,v 1.5 2006/03/18 02:42:09 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/veandco/go-sdl2/sdl"

const JOYSTICK_AXIS_MAX = 32768

var twinStick *TwinStick

/**
 * Twinstick input.
 */
type TwinStick struct {
	keys        []uint8
	rotate      float32
	reverse     float32
	enableAxis5 bool

	stick *sdl.Joystick
	state TwinStickState
}

func NewTwinStick() *TwinStick {
	this := new(TwinStick)
	this.state = NewTwinStickState()
	this.reverse = 1
	return this
}

func (this *TwinStick) openJoystick(st *sdl.Joystick /*= null*/) *sdl.Joystick {
	if st == nil {
		if sdl.InitSubSystem(sdl.INIT_JOYSTICK) < 0 {
			return nil
		}
		this.stick = sdl.JoystickOpen(0)
	} else {
		this.stick = st
	}
	return this.stick
}

func (this *TwinStick) update() {
	this.keys = sdl.GetKeyboardState()
}

func (this *TwinStick) getState() TwinStickState {
	if this.stick != nil {
		this.state.left.x = adjustAxis(int(this.stick.GetAxis(0)))
		this.state.left.y = -adjustAxis(int(this.stick.GetAxis(1)))
		var rx int = 0
		if this.enableAxis5 {
			rx = int(this.stick.GetAxis(4))
		} else {
			rx = int(this.stick.GetAxis(2))
		}
		var ry int = int(this.stick.GetAxis(3))
		if rx == 0 && ry == 0 {
			this.state.right.x = 0
			this.state.right.y = 0
		} else {
			ry = -ry
			var rd float32 = atan232(float32(rx), float32(ry))*this.reverse + this.rotate
			var rl float32 = sqrt32(float32(rx)*float32(rx) + float32(ry)*float32(ry))
			this.state.right.x = adjustAxis(int(Sin32(rd) * rl))
			this.state.right.y = adjustAxis(int(Cos32(rd) * rl))
		}
	} else {
		this.state.left.x = 0
		this.state.left.y = 0
		this.state.right.x = 0
		this.state.right.y = 0
	}
	if this.keys[sdl.SCANCODE_D] == sdl.PRESSED {
		this.state.left.x = 1
	}
	if this.keys[sdl.SCANCODE_L] == sdl.PRESSED {
		this.state.right.x = 1
	}
	if this.keys[sdl.SCANCODE_A] == sdl.PRESSED {
		this.state.left.x = -1
	}
	if this.keys[sdl.SCANCODE_J] == sdl.PRESSED {
		this.state.right.x = -1
	}
	if this.keys[sdl.SCANCODE_S] == sdl.PRESSED {
		this.state.left.y = -1
	}
	if this.keys[sdl.SCANCODE_K] == sdl.PRESSED {
		this.state.right.y = -1
	}
	if this.keys[sdl.SCANCODE_W] == sdl.PRESSED {
		this.state.left.y = 1
	}
	if this.keys[sdl.SCANCODE_I] == sdl.PRESSED {
		this.state.right.y = 1
	}
	return this.state
}

func adjustAxis(v int) float32 {
	var a float32 = 0
	if v > JOYSTICK_AXIS_MAX/3 {
		a = float32(v-JOYSTICK_AXIS_MAX/3) / (JOYSTICK_AXIS_MAX - JOYSTICK_AXIS_MAX/3)
		if a > 1 {
			a = 1
		}
	} else if v < -(JOYSTICK_AXIS_MAX / 3) {
		a = float32(v+JOYSTICK_AXIS_MAX/3) / (JOYSTICK_AXIS_MAX - JOYSTICK_AXIS_MAX/3)
		if a < -1 {
			a = -1
		}
	}
	return a
}

func (this *TwinStick) getNullState() TwinStickState {
	this.state.clear()
	return this.state
}

type TwinStickState struct {
	left, right Vector
}

func NewTwinStickState() TwinStickState {
	return TwinStickState{}
}

func (this *TwinStickState) dup(s TwinStickState) {
	this.left.x = s.left.x
	this.left.y = s.left.y
	this.right.x = s.right.x
	this.right.y = s.right.y
}

func (this *TwinStickState) clear() {
	this.left.x = 0
	this.left.y = 0
	this.right.x = 0
	this.right.y = 0
}

func (this *TwinStickState) equals(s TwinStickState) bool {
	return (this.left.x == s.left.x && this.left.y == s.left.y &&
		this.right.x == s.right.x && this.right.y == s.right.y)
}
