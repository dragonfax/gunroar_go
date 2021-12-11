package main

import (
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

func Build2DMipmaps(target int, internalFormat int, width int32, height int32, format uint32, typ uint32, data unsafe.Pointer) {
	// gl.TEXTURE_2D, 4, int(s.W), int(s.H), gl.RGBA, gl.UNSIGNED_BYTE, s.Pixels())

	// num_mipmaps := int32(4)
	// gl.TexStorage2D(gl.TEXTURE_2D, num_mipmaps, typ, width, height)
	// checkGLError()
	checkGLError()
	gl.TexImage2D(gl.TEXTURE_2D, 0, int32(format), width, height, 0, format, typ, data)
	checkGLError()
	gl.GenerateMipmap(gl.TEXTURE_2D) //Generate num_mipmaps number of mipmaps here.
	checkGLError()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	checkGLError()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	checkGLError()

}
