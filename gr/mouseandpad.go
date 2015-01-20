/*
 * $Id: mouseandpad.d,v 1.1 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

type MouseAndPadState struct {
	mouseState MouseState
	padState   PadState
}

type MouseAndPad struct {
	state MouseAndPadState
}

func NewMouseAndPad() *MouseAndPad {
	return &MouseAndPad{MouseAndPadState{}}
}

func (this *MouseAndPad) getState() MouseAndPadState {
	this.state.mouseState = mouse.getState()
	this.state.padState = pad.getState()
	return this.state
}

func (this *MouseAndPad) update() {
	pad.update()
}
