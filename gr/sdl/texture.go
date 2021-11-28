package sdl

import (
	gl "github.com/go-gl/gl/v4.1-compatibility/gl"
	sdl2 "github.com/veandco/go-sdl2/sdl"
)

const imagesDir = "images/"

var surface = make(map[string]*sdl2.Surface)

/**
 * Manage OpenGL textures.
 */
type Texture struct {
	num, maskNum               uint32
	textureNum, maskTextureNum uint32
	pixels                     [128 * 128]uint32
	maskPixels                 [128 * 128]uint32
}

func loadBmp(name string) *sdl2.Surface {
	if _, ok := surface[name]; ok {
		return surface[name]
	} else {
		fileName := imagesDir + name
		s, err := sdl2.LoadBMP(fileName)
		if err != nil {
			panic(err)
		}
		if s == nil {
			panic("Unable to load: " + fileName)
		}
		var format sdl2.PixelFormat
		format.Palette = nil
		format.BitsPerPixel = 32
		format.BytesPerPixel = 4
		format.Rmask = 0x000000ff
		format.Gmask = 0x0000ff00
		format.Bmask = 0x00ff0000
		format.Amask = 0xff000000
		/* TODO
		format.Rshift = 0
		format.Gshift = 8
		format.Bshift = 16
		format.Ashift = 24
		*/
		cs, err := s.Convert(&format, sdl2.SWSURFACE)
		if err != nil {
			panic(err)
		}
		surface[name] = cs
		return cs
	}
}

func NewTexture(name string) *Texture {
	this := &Texture{}
	s := loadBmp(name)
	gl.GenTextures(1, &this.num)
	gl.BindTexture(gl.TEXTURE_2D, this.num)

	var mode int32 = gl.RGB
	if s.Format.BytesPerPixel == 4 {
		mode = gl.RGBA
	}

	// code copied from the internet
	// data := s.Pixels()
	//sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&data))

	gl.TexImage2D(gl.TEXTURE_2D, 0, mode, s.w, s.h, 0, uint32(mode), gl.UNSIGNED_BYTE, s.Data())
	//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
	//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

	//gluBuild2DMipmaps(gl.TEXTURE_2D, 4, s.W, s.H, gl.RGBA, gl.UNSIGNED_BYTE, s.Pixels) # not found
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	return this
}

/*
I think we don't use this.
func NewTextureWithScale(name string, sx, sy, xn, yn, panelWidth, panelHeight int, maskColor uint32 = 0xffffffffu ) *Texture {
	s := loadBmp(name)
	surfacePixels := s.Pixels()
	return NewTextureWithPixels(surfacePixels, s.W, sx, sy, xn, yn, panelWidth, panelHeight, maskColor)
}

func NewTextureWithPixels(surfacePixels []uint32, surfaceWidth, sx, sy, xn, yn, panelWidth, panelHeight int, maskColor uint32 = 0xffffffffu ) *Texture {
	this := &Texture{}
	this.textureNum = xn * yn
	gl.GenTextures(int32(this.textureNum), &this.num)
	if maskColor != 0xffffffff {
		this.maskTextureNum = this.textureNum
		gl.GenTextures(int32(this.maskTextureNum), &this.maskNum)
	}
	ti := uint32(0)
	for oy := 0; oy < yn; oy++ {
		for ox := 0; ox < xn; ox++ {
			pi := 0
			for y := 0; y < panelHeight; y++ {
				for x := 0; x < panelWidth; x++ {
					p := surfacePixels[ox*panelWidth+x+sx+(oy*panelHeight+y+sy)*surfaceWidth]
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
			gl.BindTexture(gl.TEXTURE_2D, this.num+ti)
			// TODO gluBuild2DMipmaps(gl.TEXTURE_2D, 4, panelWidth, panelHeight, gl.RGBA, gl.UNSIGNED_BYTE, this.pixels)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
			if maskColor != 0xffffffff {
				gl.BindTexture(gl.TEXTURE_2D, this.maskNum+ti)
				// TODO gluBuild2DMipmaps(gl.TEXTURE_2D, 4, panelWidth, panelHeight, gl.RGBA, gl.UNSIGNED_BYTE, this.maskPixels)
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
			}
			ti++
		}
	}
	return this
}
*/

func (this *Texture) Close() {
	gl.DeleteTextures(int32(this.textureNum), &this.num)
	if this.maskTextureNum > 0 {
		gl.DeleteTextures(int32(this.maskTextureNum), &this.maskNum)
	}
}

func (this *Texture) bind(idx int /* = 0 */) {
	gl.BindTexture(gl.TEXTURE_2D, this.num+uint32(idx))
}

func (this *Texture) bindMask(idx int /* = 0 */) {
	gl.BindTexture(gl.TEXTURE_2D, this.maskNum+uint32(idx))
}
