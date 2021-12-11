package main

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
)

func checkGLError() {
	n := gl.GetError()
	if n != 0 {
		panic(fmt.Sprintf("gl error %d", n))
	}
}
