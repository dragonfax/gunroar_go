package sdl

import "github.com/veandco/go-sdl2/sdl"

/**
 * Input device interface.
 */
type Input interface {
	HandleEvent(event *sdl.Event)
}

type MultipleInputDevice struct {
	Inputs []Input
}

func (this *MultipleInputDevice) HandleEvent(event *sdl.Event) {
	for _, i := range this.Inputs {
		i.HandleEvent(event)
	}
}

var _ Input = &MultipleInputDevice{}