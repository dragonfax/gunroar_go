/*
 * $Id: mouseandpad.d,v 1.1 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type MouseAndPadState struct {
	mouseState MouseState
	padState   PadState
}

type MouseAndPad struct {
	state MouseAndPadState
	mouse *Mouse
	pad   *Pad
}

func NewMouseAndPad(mouse *Mouse, pad *Pad) *MouseAndPad {
	return &MouseAndPad{MouseAndPadState{}, mouse, pad}
}

func (this *MouseAndPad) getState() MouseAndPadState {
	this.state.mouseState = this.mouse.getState()
	this.state.padState = this.pad.getState()
	return this.state
}

func (this *MouseAndPad) handleEvent(event sdl.Event) {
	this.mouse.handleEvent(event)
	this.pad.handleEvent(event)
}
