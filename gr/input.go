package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Input interface {
	handleEvent(event *sdl.Event)
}
