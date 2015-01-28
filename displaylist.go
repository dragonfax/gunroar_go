/*
 * $Id: letter.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/go-gl/gl"

type DisplayList struct {
	registered bool
	num        uint
	idx        uint
	enumIdx    uint
}

var displayListsFinalized = false

func NewDisplayList(num uint) *DisplayList {
	verifyNotFinalized()
	dl := new(DisplayList)
	dl.num = num
	dl.idx = gl.GenLists(int(num))
	return dl
}

func verifyNotFinalized() {
	if displayListsFinalized {
		panic("illegal method. display lists already finalized")
	}
}

func (dp *DisplayList) beginSingleList() {
	verifyNotFinalized()
	if dp.num > 1 {
		panic("can't use for multi lists")
	}
	if dp.registered {
		panic("already registered this list")
	}
	dp.ResetLists()
	dp.NewList()
}

func (dp *DisplayList) endSingleList() {
	verifyNotFinalized()
	if dp.num > 1 {
		panic("can't use for multi lists")
	}
	if dp.registered {
		panic("already registered this list")
	}
	gl.EndList()
	dp.registered = true
}

func (dp *DisplayList) ResetLists() {
	verifyNotFinalized()
	dp.enumIdx = dp.idx
}

func (dp *DisplayList) NewList() {
	verifyNotFinalized()
	gl.NewList(dp.enumIdx, gl.COMPILE)
}

func (dp *DisplayList) EndList() {
	verifyNotFinalized()
	gl.EndList()
	dp.enumIdx++
	dp.registered = true
}

func (dp *DisplayList) call(i uint /* = 0 */) {
	gl.CallList(dp.idx + i)
}

func (dp *DisplayList) close() {
	if !dp.registered {
		return
	}
	gl.DeleteLists(dp.idx, int(dp.num))
}
