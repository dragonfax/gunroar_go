package main

import (
	"./gr"
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glu"
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
		fmt.Printf("%v\n", err)
		return
	}

	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle("test")
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(onKey)

	initGL()

	letter = new(gr.Letter)
	letter.Init()

	running = true
	for running && glfw.WindowParam(glfw.Opened) == 1 {
		drawScene()
	}
}

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	glu.Perspective(45.0, float64(w)/float64(h), 0.1, 100.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func onKey(key, state int) {
	switch key {
	case glfw.KeyEsc:
		running = false
	}
}

func initGL() {
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	gl.Translatef(-1.5, 0, -6)

	gl.Begin(gl.TRIANGLES)
	gl.Color3f(1, 0, 0)
	gl.Vertex3f(0, 1, 0)
	gl.Color3f(0, 1, 0)
	gl.Vertex3f(-1, -1, 0)
	gl.Color3f(0, 0, 1)
	gl.Vertex3f(1, -1, 0)
	gl.End()

	gl.Translatef(3, 0, 0)
	gl.Color3f(0.5, 0.5, 1.0)

	gl.Begin(gl.QUADS)
	gl.Vertex3f(-1, 1, 0)
	gl.Vertex3f(1, 1, 0)
	gl.Vertex3f(1, -1, 0)
	gl.Vertex3f(-1, -1, 0)
	gl.End()

	drawString()

	glfw.SwapBuffers()
}

var letter *gr.Letter

func drawString() {

	for i := float32(0); i < 5; i++ {
		letter.DrawStringOption("lets do this", 0, 0, i/10, gr.TO_RIGHT, int(i), false, i)
	}
}
