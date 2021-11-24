package main

import (
	"github.com/dragonfax/gunroar/gr/sdl/mouse"
	"github.com/dragonfax/gunroar/gr/sdl/screen"
)

const MOUSE_SCREEN_MAPPING_RATIO_X = 26.0
const MOUSE_SCREEN_MAPPING_RATIO_Y = 19.5

type RecordableMouse struct {
	*mouse.RecordableMouse
	screen screen.SizableScreen
}

func NewMouse(screen screen.SizableScreen) *RecordableMouse {
	this := &RecordableMouse{RecordableMouse: mouse.NewRecordableMouse()}
	this.screen = screen
	return this
}

func (this *RecordableMouse) AdjustPos(ms *mouse.MouseState) {
	ms.X = (ms.X - float64(this.screen.Width()/2)) * MOUSE_SCREEN_MAPPING_RATIO_X / float64(this.screen.Width())
	ms.Y = -(ms.Y - float64(this.screen.Height()/2)) * MOUSE_SCREEN_MAPPING_RATIO_Y / float64(this.screen.Height())
}
