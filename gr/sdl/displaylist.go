package sdl

import "github.com/go-gl/gl/v4.1-compatibility/gl"

/**
 * Manage a display list.
 */
type DisplayList struct {
	registered bool
	num        uint32
	idx        uint32
	enumIdx    uint32
}

func NewDisplayList(num uint32) *DisplayList {
	this := &DisplayList{
		num: num,
		idx: gl.GenLists(int32(num)),
	}
	return this
}

func (this *DisplayList) BeginNewList() {
	this.ResetList()
	this.NewList()
}

func (this *DisplayList) NextNewList() {
	gl.EndList()
	this.enumIdx++
	if this.enumIdx >= this.idx+this.num || this.enumIdx < this.idx {
		panic("Can't create new list. Index out of bound.")
	}
	gl.NewList(this.enumIdx, gl.COMPILE)
}

func (this *DisplayList) EndNewList() {
	gl.EndList()
	this.registered = true
}

func (this *DisplayList) ResetList() {
	this.enumIdx = this.idx
}

func (this *DisplayList) NewList() {
	gl.NewList(this.enumIdx, gl.COMPILE)
}

func (this *DisplayList) EndList() {
	gl.EndList()
	this.enumIdx++
	this.registered = true
}

func (this *DisplayList) Call(i int /* = 0 */) {
	gl.CallList(this.idx + uint32(i))
}

func (this *DisplayList) Close() {
	if !this.registered {
		return
	}
	gl.DeleteLists(this.idx, int32(this.num))
}
