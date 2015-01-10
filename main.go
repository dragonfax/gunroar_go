package main

import (
	"./gr"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"os"
	"runtime"
)

func main() {

	runtime.LockOSThread()

	if !glfw.Init() {
		fmt.Fprintf(os.Stderr, "Can't open GLFW")
		return
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True) // needed for macs

	window, err := glfw.CreateWindow(1024, 768, "Tutorial 1", nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	window.MakeContextCurrent()

	gl.Init()
	gl.GetError() // Ignore error
	window.SetInputMode(glfw.StickyKeys, 1)

	gl.ClearColor(0., 0., 0.4, 0.)

	l := new(gr.Letter)
	l.Init()

	l.DrawString("lets do this", 1, 1, 1)

	// Equivalent to a do... while
	for ok := true; ok; ok = (window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose()) {
		window.SwapBuffers()
		glfw.PollEvents()
	}

}
