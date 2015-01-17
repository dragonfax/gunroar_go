/*
 * $Id: letter.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"errors"

	"github.com/go-gl/gl"
)

type DisplayList struct {
	registered bool
	num        uint
	idx        uint
	enumIdx    uint
}

func NewDisplayList(num uint) *DisplayList {
	dl := &DisplayList{}
	dl.num = num
	dl.idx = gl.GenLists(int(num))
	return dl
}

func (dp *DisplayList) beginNewList() {
	dp.ResetList()
	dp.NewList()
}

func (dp *DisplayList) nextNewList() error {
	gl.EndList()
	dp.enumIdx++
	if dp.enumIdx >= dp.idx+dp.num || dp.enumIdx < dp.idx {
		return errors.New("Can't create new list. Index out of bound.")
	}
	gl.NewList(dp.enumIdx, gl.COMPILE)
	return nil
}

func (dp *DisplayList) endNewList() {
	gl.EndList()
	dp.registered = true
}

func (dp *DisplayList) ResetList() {
	dp.enumIdx = dp.idx
}

func (dp *DisplayList) NewList() {
	gl.NewList(dp.enumIdx, gl.COMPILE)
}

func (dp *DisplayList) EndList() {
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
