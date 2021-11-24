package sdl

import (
	. "github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/sdl/file"
	"github.com/dragonfax/gunroar/gr/sdl/record"
	"github.com/veandco/go-sdl2/sdl"
)

const MOUSE_SCREEN_MAPPING_RATIO_X = 26.0
const MOUSE_SCREEN_MAPPING_RATIO_Y = 19.5

var _ Input = &Mouse{}

/**
 * Mouse input.
 */
type Mouse struct {
	screen SizableScreen
	state  MouseState
}

func New(screen SizableScreen) *Mouse {
	this := &Mouse{}
	this.state = MouseState{}
	this.screen = screen
	return this
}

func (this *Mouse) Init(screen SizableScreen) {
	this.screen = screen
}

func (this *Mouse) HandleEvent(event *sdl.Event) {
}

func (this *Mouse) GetState(doRecord bool /* = true */) MouseState {
	mx, my, btn := sdl.GetMouseState()
	this.state.X = float64(mx)
	this.state.Y = float64(my)
	this.state.Button = 0
	if btn&sdl.Button(1) != 0 {
		this.state.Button |= MouseButtonLEFT
	}
	if btn&sdl.Button(3) != 0 {
		this.state.Button |= MouseButtonRIGHT
	}
	this.adjustPos(&this.state)
	if doRecord {
		record(this.state)
	}
	return this.state
}

func (this *Mouse) adjustPos(ms *MouseState) {
	ms.X = (ms.X - this.screen.Width/2) * MOUSE_SCREEN_MAPPING_RATIO_X / this.screen.Width
	ms.Y = -(ms.Y - this.screen.Height/2) * MOUSE_SCREEN_MAPPING_RATIO_Y / this.screen.Height
}

func (this *Mouse) GetNullState() MouseState {
	this.state.Clear()
	return this.state
}

type MouseButton int

const MouseButtonLEFT MouseButton = 1
const MouseButtonRIGHT MouseButton = 2

type MouseState struct {
	X, Y   float64
	Button MouseButton
}

func (this *MouseState) Clear() {
	this.Button = 0
}

func (this *MouseState) Read(fd file.File) {
	this.X = fd.ReadFloat64()
	this.Y = fd.ReadFloat64()
	this.Button = MouseButton(fd.ReadInt())
}

func (this *MouseState) Write(fd file.File) {
	fd.WriteFloat64(this.X)
	fd.WriteFloat64(this.Y)
	fd.WriteInt(int(this.Button))
}

func (this MouseState) Equals(s MouseState) bool {
	return this.X == s.X && this.Y == s.Y && this.Button == s.Button
}

type RecordableMouse struct {
	Mouse
	record.RecordableInput
}

func (this *RecordableMouse) GetState(doRecord bool /* = true */) MouseState {
	var s = this.GetState()
	if doRecord {
		this.record(s)
	}
	return s
}
