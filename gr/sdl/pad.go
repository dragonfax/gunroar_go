package sdl

import (
	"github.com/dragonfax/gunroar/gr/sdl/file"
	"github.com/dragonfax/gunroar/gr/sdl/record"
	sdl2 "github.com/veandco/go-sdl2/sdl"
)

/**
 * Joystick and keyboard input.
 */

const JOYSTICK_AXIS = 16384

var _ Input = &Pad{}

type Pad struct {
	keys           []uint8
	buttonReversed bool
	stick          *sdl2.Joystick
	state          PadState
}

func NewPad() Pad {
	return Pad{}
}

func (this *Pad) openJoystick(st *sdl2.Joystick /* = null*/) *sdl2.Joystick {
	if st == nil {
		err := sdl2.InitSubSystem(sdl2.INIT_JOYSTICK)
		if err != nil {
			return nil
		}
		this.stick = sdl2.JoystickOpen(0)
	} else {
		this.stick = st
	}
	return this.stick
}

func (this *Pad) HandleEvent(event sdl2.Event) {
	this.keys = sdl2.GetKeyboardState()
}

func (this *Pad) getState() PadState {
	var x, y int16
	this.state.dir = NONE
	if this.stick != nil {
		x = this.stick.Axis(0)
		y = this.stick.Axis(1)
	}
	if this.keys[sdl2.K_RIGHT] == sdl2.PRESSED || this.keys[sdl2.K_KP_6] == sdl2.PRESSED ||
		this.keys[sdl2.K_d] == sdl2.PRESSED || this.keys[sdl2.K_l] == sdl2.PRESSED ||
		x > JOYSTICK_AXIS {
		this.state.dir |= RIGHT
	}
	if this.keys[sdl2.K_LEFT] == sdl2.PRESSED || this.keys[sdl2.K_KP_4] == sdl2.PRESSED ||
		this.keys[sdl2.K_a] == sdl2.PRESSED || this.keys[sdl2.K_j] == sdl2.PRESSED ||
		x < -JOYSTICK_AXIS {
		this.state.dir |= LEFT
	}
	if this.keys[sdl2.K_DOWN] == sdl2.PRESSED || this.keys[sdl2.K_KP_2] == sdl2.PRESSED ||
		this.keys[sdl2.K_s] == sdl2.PRESSED || this.keys[sdl2.K_k] == sdl2.PRESSED ||
		y > JOYSTICK_AXIS {
		this.state.dir |= DOWN
	}
	if this.keys[sdl2.K_UP] == sdl2.PRESSED || this.keys[sdl2.K_KP_8] == sdl2.PRESSED ||
		this.keys[sdl2.K_w] == sdl2.PRESSED || this.keys[sdl2.K_i] == sdl2.PRESSED ||
		y < -JOYSTICK_AXIS {
		this.state.dir |= UP
	}
	this.state.button = 0
	var btn1, btn2 byte
	if this.stick != nil {
		btn1 = this.stick.Button(0) + this.stick.Button(3) +
			this.stick.Button(4) + this.stick.Button(7) +
			this.stick.Button(8) + this.stick.Button(11)
		btn2 = this.stick.Button(1) + this.stick.Button(2) +
			this.stick.Button(5) + this.stick.Button(6) +
			this.stick.Button(9) + this.stick.Button(10)
	}
	if this.keys[sdl2.K_z] == sdl2.PRESSED || this.keys[sdl2.K_PERIOD] == sdl2.PRESSED ||
		this.keys[sdl2.K_LCTRL] == sdl2.PRESSED || this.keys[sdl2.K_RCTRL] == sdl2.PRESSED ||
		btn1 > 0 {
		if !this.buttonReversed {
			this.state.button |= ButtonA
		} else {
			this.state.button |= ButtonB
		}
	}
	if this.keys[sdl2.K_x] == sdl2.PRESSED || this.keys[sdl2.K_SLASH] == sdl2.PRESSED ||
		this.keys[sdl2.K_LALT] == sdl2.PRESSED || this.keys[sdl2.K_RALT] == sdl2.PRESSED ||
		this.keys[sdl2.K_LSHIFT] == sdl2.PRESSED || this.keys[sdl2.K_RSHIFT] == sdl2.PRESSED ||
		this.keys[sdl2.K_RETURN] == sdl2.PRESSED ||
		btn2 > 0 {
		if !this.buttonReversed {
			this.state.button |= ButtonB
		} else {
			this.state.button |= ButtonA
		}
	}
	return this.state
}

func (this *Pad) getNullState() PadState {
	this.state.clear()
	return this.state
}

type Dir int

const (
	NONE  Dir = 0
	UP    Dir = 1
	DOWN  Dir = 2
	LEFT  Dir = 4
	RIGHT Dir = 8
)

type Button int

const (
	ButtonA   = 16
	ButtonB   = 32
	ButtonANY = 48
)

type PadState struct {
	dir    Dir
	button Button
}

func NewPadState(s PadState) PadState {
	this := PadState{}
	this.set(s)
	return this
}

func (this *PadState) set(s PadState) {
	this.dir = s.dir
	this.button = s.button
}

func (this *PadState) clear() {
	this.dir = 0
	this.button = 0
}

func (this *PadState) Read(fd file.File) {
	s := fd.ReadInt()
	this.dir = Dir(s & (int(UP) | int(DOWN) | int(LEFT) | int(RIGHT)))
	this.button = Button(s & int(ButtonANY))
}

func (this *PadState) Write(fd file.File) {
	s := int(this.dir) | int(this.button)
	fd.WriteInt(s)
}

func (this *PadState) Equals(s PadState) bool {
	return this.dir == s.dir && this.button == s.button
}

type RecordablePad struct {
	Pad
	record.RecordableInput
}

func NewRecordablePad() *RecordablePad {
	this := &RecordablePad{
		Pad:             NewPad(),
		RecordableInput: record.NewRecordableInput(),
	}
	return this
}

func (this RecordablePad) GetState(doRecord bool /*= true */) PadState {
	s := this.Pad.getState()
	if doRecord {
		this.Record(s)
	}
	return s
}
