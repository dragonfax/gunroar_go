package glu

// #cgo darwin LDFLAGS: -framework Carbon -framework OpenGL -framework GLUT -framework GLKit
//
// #ifdef __APPLE__
//   #include <OpenGL/glu.h>
//	 #include <GLKit/GLKMatrix4.h>
// #endif
import "C"

func LookAt(eyeX, eyeY, eyeZ, centerX, centerY, centerZ, upX, upY, upZ float64) {
	/* m4 := */ C.GLKMatrix4MakeLookAt(
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

	// TODO gl.LoadMatrixf((*float32)(unsafe.Pointer(&m4)))
}
