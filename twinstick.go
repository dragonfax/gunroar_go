/*
 * $Id: twinstick.d,v 1.5 2006/03/18 02:42:09 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

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
		err := sdl.InitSubSystem(sdl.INIT_JOYSTICK)
		if err != nil {
			panic(err)
		}
		if n := sdl.NumJoysticks(); n <= 0 {
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
		this.state.left.x = adjustAxis(float32(this.stick.Axis(0)))
		this.state.left.y = -adjustAxis(float32(this.stick.Axis(1)))
		var rx int16 = 0
		if this.enableAxis5 {
			rx = this.stick.Axis(4)
		} else {
			rx = this.stick.Axis(2)
		}
		ry := this.stick.Axis(3)
		if rx == 0 && ry == 0 {
			this.state.right.x = 0
			this.state.right.y = 0
		} else {
			if ry != math.MinInt16 {
				ry = -ry
			} else {
				ry = math.MaxInt16
			}
			// apply any configured rotation, and reversal of joystick axis
			var rd float32 = atan232(float32(rx), float32(ry))*this.reverse + this.rotate
			var rl float32 = sqrt32(float32(rx)*float32(rx) + float32(ry)*float32(ry))
			this.state.right.x = adjustAxis(Sin32(rd) * rl)
			temp := adjustAxis(Cos32(rd) * rl)
			this.state.right.y = temp
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

/* axis defaults to 0
 * If moved past 1/3 of range,
 * then it becomes a ratio of the remaining range (2/3 of original range).
 * with a max of 1
 */
func adjustAxis(v float32) float32 {
	var a float32 = 0
	if v > JOYSTICK_AXIS_MAX/3 {
		a = (v - JOYSTICK_AXIS_MAX/3) / (JOYSTICK_AXIS_MAX - JOYSTICK_AXIS_MAX/3)
		if a > 1 {
			a = 1
		}
	} else if v < -(JOYSTICK_AXIS_MAX / 3) {
		a = (v + JOYSTICK_AXIS_MAX/3) / (JOYSTICK_AXIS_MAX - JOYSTICK_AXIS_MAX/3)
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
