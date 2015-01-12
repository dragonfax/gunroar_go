package main

import (
	"fmt"
	"github.com/dragonfax/gunroar_go/gr"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	// "github.com/go-gl/glu"
)

const Width = 600
const Height = 600

var running bool = false

func main() {
	var err error
	if err = glfw.Init(); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	defer glfw.Terminate()

	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 0, 8, glfw.Windowed); err != nil {
		panic(err)
	}

	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle("test")
	glfw.SetKeyCallback(onKey)

	gr.InitLetter()
	defer gr.CloseLetter()

	running = true
	for running && glfw.WindowParam(glfw.Opened) == 1 {
		drawScene()
	}

}

func onKey(key, state int) {
	switch key {
	case glfw.KeyEsc:
		running = false
	}
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gr.DrawStringOption("lets do this", 0, 0, 0.1, gr.TO_RIGHT, 0, false, 0)

	glfw.SwapBuffers()
}
