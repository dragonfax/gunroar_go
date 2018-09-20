/*
 * $Id: letter.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import "github.com/go-gl/gl/v2.1/gl"

type DisplayList struct {
	registered bool
	num        uint32
	idx        uint32
	enumIdx    uint32
}

var displayListsFinalized = false

func NewDisplayList(num uint32) *DisplayList {
	verifyNotFinalized()
	dl := new(DisplayList)
	dl.num = num
	dl.idx = gl.GenLists(int32(num))
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

func (dp *DisplayList) call(i uint32 /* = 0 */) {
	gl.CallList(dp.idx + i)
}

func (dp *DisplayList) close() {
	if !dp.registered {
		return
	}
	gl.DeleteLists(dp.idx, int32(dp.num))
}
