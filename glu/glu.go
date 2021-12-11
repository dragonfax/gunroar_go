package glu

// #cgo darwin LDFLAGS: -framework Carbon -framework OpenGL -framework GLUT -framework GLKit
//
// #ifdef __APPLE__
//   #include <OpenGL/glu.h>
//	 #include <GLKit/GLKMatrix4.h>
// #endif
import "C"
import (
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

func LookAt(eyeX, eyeY, eyeZ, centerX, centerY, centerZ, upX, upY, upZ float64) {
	m4 := C.GLKMatrix4MakeLookAt(
		C.GLfloat(eyeX),
		C.GLfloat(eyeY),
		C.GLfloat(eyeZ),
		C.GLfloat(centerX),
		C.GLfloat(centerY),
		C.GLfloat(centerZ),
		C.GLfloat(upX),
		C.GLfloat(upY),
		C.GLfloat(upZ),
	)

	gl.LoadMatrixf((*float32)(unsafe.Pointer(&m4)))
}
