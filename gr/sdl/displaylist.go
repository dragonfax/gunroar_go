/*
 * $Id: letter.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package sdl

import (
	"github.com/jackyb/go-gl/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type DisplayList struct {
	registered bool
	num        int
	idx        int
	enumIdx    int
}

func NewDisplayList(int num) *DisplayList {
	dl := &DisplayList{}
	dl.num = num
	dl.idx = glGenLists(num)
	return dl
}

func (dp *DisplayList) BeginNewList() {
	dp.ResetList()
	dp.NewList()
}

func (dp *DisplayList) NextNewList() error {
	glEndList()
	dp.enumIdx++
	if dp.enumIdx >= dp.idx+dp.num || dp.enumIdx < dp.idx {
		return errors.error("Can't create new list. Index out of bound.")
	}
	glNewList(dp.enumIdx, GL_COMPILE)
	return nil
}

func (dp *DisplayList) EndNewList() {
	glEndList()
	dp.registered = true
}

func (dp *DisplayList) ResetList() {
	dp.enumIdx = dp.idx
}

func (dp *DisplayList) NewList() {
	glNewList(dp.enumIdx, GL_COMPILE)
}

func (dp *DisplayList) EndList() {
	glEndList()
	dp.enumIdx++
	dp.registered = true
}

func (dp *DisplayList) Call(int i) { // default value should be 0
	glCallList(dp.idx + i)
}

func (dp *DisplayList) Close() {
	if !dp.registered {
		return
	}
	glDeleteLists(dp.idx, dp.num)
}
