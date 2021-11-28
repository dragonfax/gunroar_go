package sdl

import "github.com/go-gl/gl/v4.1-compatibility/gl"

const TEXTURE_SIZE_MIN = 0.02
const TEXTURE_SIZE_MAX = 0.98
const LUMINOUS_TEXTURE_WIDTH_MAX = 64
const LUMINOUS_TEXTURE_HEIGHT_MAX = 64

const textureLen = LUMINOUS_TEXTURE_WIDTH_MAX * LUMINOUS_TEXTURE_HEIGHT_MAX * 4 /* * uint.sizeof */
const lmOfsBs = 3.0

/**
 * Luminous effect texture.
 */
type LuminousScreen struct {
	luminousTexture           uint32
	td                        [textureLen]uint32
	luminousTextureWidth      int
	luminousTextureHeight     int
	screenWidth, screenHeight int
	luminosity                float64
	lmOfs                     [2][2]float64
}

func NewLuminousScreen() *LuminousScreen {
	this := &LuminousScreen{
		luminousTextureWidth:  64,
		luminousTextureHeight: 64,
		lmOfs:                 [2][2]float64{{-2, -1}, {2, 1}},
	}
	return this
}

func (this *LuminousScreen) Init(luminosity float64, width, height int) {
	this.makeLuminousTexture()
	this.luminosity = luminosity
	this.Resized(width, height)
}

func (this *LuminousScreen) makeLuminousTexture() {
	// uint *data = td;
	// TODO data := this.td
	// TODO I don't really know what this is doing? generating a new one? clearing it? what?
	// TODO memset(data, 0, this.luminousTextureWidth*this.luminousTextureHeight*4*uint.sizeof)
	gl.GenTextures(1, &this.luminousTexture)
	gl.BindTexture(gl.TEXTURE_2D, this.luminousTexture)
	// TODO gl.TexImage2D(gl.TEXTURE_2D, 0, 4, this.luminousTextureWidth, this.luminousTextureHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
}

func (this *LuminousScreen) Resized(width, height int) {
	this.screenWidth = width
	this.screenHeight = height
}

func (this *LuminousScreen) StartRender() {
	gl.Viewport(0, 0, int32(this.luminousTextureWidth), int32(this.luminousTextureHeight))
}

func (this *LuminousScreen) EndRender() {
	gl.BindTexture(gl.TEXTURE_2D, this.luminousTexture)
	gl.CopyTexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		0, 0, int32(this.luminousTextureWidth), int32(this.luminousTextureHeight), 0)
	gl.Viewport(0, 0, int32(this.screenWidth), int32(this.screenHeight))
}

func (this *LuminousScreen) viewOrtho() {
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(0, float64(this.screenWidth), float64(this.screenHeight), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()
}

func (this *LuminousScreen) viewPerspective() {
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.MatrixMode(gl.MODELVIEW)
	gl.PopMatrix()
}

func (this *LuminousScreen) Draw() {
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, this.luminousTexture)
	this.viewOrtho()
	gl.Color4d(1, 0.8, 0.9, this.luminosity)
	gl.Begin(gl.QUADS)
	for i := 0; i < 2; i++ {
		gl.TexCoord2f(TEXTURE_SIZE_MIN, TEXTURE_SIZE_MAX)
		gl.Vertex2d(0+this.lmOfs[i][0]*lmOfsBs, 0+this.lmOfs[i][1]*lmOfsBs)
		gl.TexCoord2d(TEXTURE_SIZE_MIN, TEXTURE_SIZE_MIN)
		gl.Vertex2d(0+this.lmOfs[i][0]*lmOfsBs, float64(this.screenHeight)+this.lmOfs[i][1]*lmOfsBs)
		gl.TexCoord2d(TEXTURE_SIZE_MAX, TEXTURE_SIZE_MIN)
		gl.Vertex2d(float64(this.screenWidth)+this.lmOfs[i][0]*lmOfsBs, float64(this.screenHeight)+this.lmOfs[i][0]*lmOfsBs)
		gl.TexCoord2d(TEXTURE_SIZE_MAX, TEXTURE_SIZE_MAX)
		gl.Vertex2d(float64(this.screenWidth)+this.lmOfs[i][0]*lmOfsBs, 0+this.lmOfs[i][0]*lmOfsBs)
	}
	gl.End()
	this.viewPerspective()
	gl.Disable(gl.TEXTURE_2D)
}
