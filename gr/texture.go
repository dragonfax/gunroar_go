/*
 * $Id: texture.d,v 1.2 2005/07/03 07:05:23 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package sdl

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/go-gl/gl"
	"github.com/go-gl/glu"
	"github.com/veandco/go-sdl2/sdl"
)

const imagesDir = "images/"

var surface = make(map[string]*sdl.Surface)

type Texture struct {
	textures, maskTextures []gl.Texture
	pixels, maskPixels     [128 * 128]uint32
}

func LoadBmp(name string) *sdl.Surface {
	if val, ok := surface[name]; ok {
		return val
	} else {
		fileName := imagesDir + name
		s := sdl.LoadBMP(fileName)
		if s == nil {
			panic(errors.New("SDLInitFailedException: Unable to load: " + fileName))
		}
		var format sdl.PixelFormat
		format.Palette = nil
		format.BitsPerPixels = 32
		format.BytesPerPixel = 4
		format.Rmask = 0x000000ff
		format.Gmask = 0x0000ff00
		format.Bmask = 0x00ff0000
		format.Amask = 0xff000000
		format.Rshift = 0
		format.Gshift = 8
		format.Bshift = 16
		format.Ashift = 24
		format.Rloss = 0
		format.Gloss = 0
		format.Bloss = 0
		format.Aloss = 0
		/* format.Alpha = 0 */
		cs := s.Convert(&format, sdl.SWSURFACE)
		surface[name] = cs
		return cs
	}
}

func NewTextureFromBMP(name string) *Texture {
	this := new(Texture)
	s := LoadBmp(name)
	this.textures = make([]gl.Texture, 1)

	gl.GenTextures(this.textures)

	this.textures[0].Bind(gl.TEXTURE_2D)
	glu.Build2DMipmaps(gl.TEXTURE_2D, 4, int(s.W), int(s.H),
		gl.RGBA, gl.UNSIGNED_BYTE, s.Pixels())
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	return this
}

func NewTextureFromBMPOption(name string, sx int, sy int, xn int, yn int, panelWidth int, panelHeight int, maskColor uint32 /* = 0xffffffffu */) *Texture {
	s := LoadBmp(name)
	pi, err := ByteArrayToUint32Array(s.Pixels())
	if err != nil {
		panic(err)
	}
	return NewTexture(pi, int(s.W), sx, sy, xn, yn, panelWidth, panelHeight, maskColor)
}

func ByteArrayToUint32Array(ary []byte) ([]uint32, error) {
	var pi []uint32
	buf := bytes.NewReader(ary)
	err := binary.Read(buf, binary.BigEndian, &pi)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func NewTexture(surfacePixels []uint32, surfaceWidth int,
	sx int, sy int, xn int, yn int, panelWidth int, panelHeight int,
	maskColor uint32 /* = 0xffffffffu */) *Texture {
	this := new(Texture)

	textureNum := xn * yn
	this.textures = make([]gl.Texture, textureNum)
	gl.GenTextures(this.textures)
	if maskColor != 0xffffffff {
		maskTextureNum := textureNum
		this.maskTextures = make([]gl.Texture, maskTextureNum)
		gl.GenTextures(this.maskTextures)
	}
	ti := 0
	for oy := 0; oy < yn; oy++ {
		for ox := 0; ox < xn; ox++ {
			pi := 0
			for y := 0; y < panelHeight; y++ {
				for x := 0; x < panelWidth; x++ {
					var p uint32 = surfacePixels[ox*panelWidth+x+sx+(oy*panelHeight+y+sy)*surfaceWidth]
					var m uint32
					if p == maskColor {
						p = 0xff000000
						m = 0x00ffffff
					} else {
						m = 0x00000000
					}
					this.pixels[pi] = p
					if maskColor != 0xffffffff {
						this.maskPixels[pi] = m
					}
					pi++
				}
			}
			this.textures[ti].Bind(gl.TEXTURE_2D)
			glu.Build2DMipmaps(gl.TEXTURE_2D, 4, panelWidth, panelHeight,
				gl.RGBA, gl.UNSIGNED_BYTE, this.pixels)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
			if maskColor != 0xffffffff {
				this.maskTextures[ti].Bind(gl.TEXTURE_2D)
				glu.Build2DMipmaps(gl.TEXTURE_2D, 4, panelWidth, panelHeight,
					gl.RGBA, gl.UNSIGNED_BYTE, this.maskPixels)
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
			}
			ti++
		}
	}
	return this
}

func (t *Texture) Close() {
	gl.DeleteTextures(t.textures)
	if len(t.maskTextures) != 0 {
		gl.DeleteTextures(t.maskTextures)
	}
}

func (t *Texture) Bind(idx int /* = 0 */) {
	t.textures[idx].Bind(gl.TEXTURE_2D)
}

func (t *Texture) BindMask(idx int /* = 0 */) {
	t.maskTextures[idx].Bind(gl.TEXTURE_2D)
}
