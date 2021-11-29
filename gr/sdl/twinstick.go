package sdl

import (
	"fmt"
	"math"

	"github.com/dragonfax/gunroar/gr/sdl/file"
	"github.com/dragonfax/gunroar/gr/sdl/record"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/veandco/go-sdl2/sdl"
)

var _ Input = &TwinStick{}

const JOYSTICK_AXIS_MAX = 32768

/**
 * Twinstick input.
 */
type TwinStick struct {
	Rotate      float64
	Reverse     float64
	keys        []uint8
	EnableAxis5 bool
	stick       *sdl.Joystick
	state       TwinStickState
}

func NewTwinStick() TwinStick {
	this := TwinStick{Reverse: 1}
	return this
}

func (this *TwinStick) openJoystick(st *sdl.Joystick) *sdl.Joystick {
	var stick *sdl.Joystick
	if st == nil {

		err := sdl.InitSubSystem(sdl.INIT_JOYSTICK)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		stick = sdl.JoystickOpen(0)
	} else {
		stick = st
	}
	return stick
}

func (this *TwinStick) HandleEvent(event sdl.Event) {
	this.keys = sdl.GetKeyboardState()
}

func (this *TwinStick) GetState() TwinStickState {
	if this.stick != nil {
		this.state.left.X = this.adjustAxis(this.stick.Axis(0))
		this.state.left.Y = -this.adjustAxis(this.stick.Axis(1))
		var rx int16 = 0
		if this.EnableAxis5 {
			rx = this.stick.Axis(4)
		} else {
			rx = this.stick.Axis(2)
		}
		ry := this.stick.Axis(3)
		if rx == 0 && ry == 0 {
			this.state.right.X = 0
			this.state.right.Y = 0
		} else {
			ry = -ry
			rd := math.Atan2(float64(rx), float64(ry))*this.Reverse + this.Rotate
			rl := math.Sqrt(float64(rx*rx + ry*ry))
			this.state.right.X = this.adjustAxis(int16(math.Sin(rd) * rl))
			this.state.right.Y = this.adjustAxis(int16(math.Cos(rd) * rl))
		}
	} else {
		this.state.left.X = 0
		this.state.left.Y = 0
		this.state.right.X = 0
		this.state.right.Y = 0
	}
	if this.keys[sdl.K_d] == sdl.PRESSED {
		this.state.left.X = 1
	}
	if this.keys[sdl.K_l] == sdl.PRESSED {
		this.state.right.X = 1
	}
	if this.keys[sdl.K_a] == sdl.PRESSED {
		this.state.left.X = -1
	}
	if this.keys[sdl.K_j] == sdl.PRESSED {
		this.state.right.X = -1
	}
	if this.keys[sdl.K_s] == sdl.PRESSED {
		this.state.left.Y = -1
	}
	if this.keys[sdl.K_k] == sdl.PRESSED {
		this.state.right.Y = -1
	}
	if this.keys[sdl.K_w] == sdl.PRESSED {
		this.state.left.Y = 1
	}
	if this.keys[sdl.K_i] == sdl.PRESSED {
		this.state.right.Y = 1
	}
	return this.state
}

func (this *TwinStick) adjustAxis(v int16) float64 {
	var a int16
	if v > JOYSTICK_AXIS_MAX/3 {
		a = (v - JOYSTICK_AXIS_MAX/3) /
			(JOYSTICK_AXIS_MAX - JOYSTICK_AXIS_MAX/3)
		if a > 1 {
			a = 1
		}
	} else if v < -(JOYSTICK_AXIS_MAX / 3) {
		a = (v + JOYSTICK_AXIS_MAX/3) /
			(JOYSTICK_AXIS_MAX - JOYSTICK_AXIS_MAX/3)
		if a < -1 {
			a = -1
		}
	}
	return float64(a) // TODO its possible float cast shoudl be deeper in this funcdtion,to avoid precision loss.
}

func (this *TwinStick) GetNullState() TwinStickState {
	this.state.clear()
	return this.state
}

type TwinStickState struct {
	left, right vector.Vector
}

func NewTwinStickState(i record.InputState) record.InputState {

	if i == nil {
		return &TwinStickState{}
	} else {
		s, ok := i.(*TwinStickState)
		if !ok {
			panic("wrong state given to NewTwinStickStat")
		}
		// copy it.
		s2 := *s
		return &s2
	}
}

func (this *TwinStickState) Set(i record.InputState) {
	s, ok := i.(*TwinStickState)
	if !ok {
		panic("wrong state type given to TwinStickState.Set")
	}
	this.left.X = s.left.X
	this.left.Y = s.left.Y
	this.right.X = s.right.X
	this.right.Y = s.right.Y
}

func (this *TwinStickState) clear() {
	this.left.X = 0
	this.left.Y = 0
	this.right.X = 0
	this.right.Y = 0
}

func (this *TwinStickState) Read(fd file.File) {
	this.left.X = fd.ReadFloat64()
	this.left.Y = fd.ReadFloat64()
	this.right.X = fd.ReadFloat64()
	this.right.Y = fd.ReadFloat64()
}

func (this *TwinStickState) Write(fd file.File) {
	fd.WriteFloat64(this.left.X)
	fd.WriteFloat64(this.left.Y)
	fd.WriteFloat64(this.right.X)
	fd.WriteFloat64(this.right.Y)
}

func (this *TwinStickState) Equals(i record.InputState) bool {
	s, ok := i.(*TwinStickState)
	if !ok {
		panic("wrong state given to TwinStickState")
	}
	return this.left.X == s.left.X && this.left.Y == s.left.Y &&
		this.right.X == s.right.X && this.right.Y == s.right.Y
}

type RecordableTwinStick struct {
	TwinStick
	record.RecordableInput
}

func NewRecordableTwinStick() *RecordableTwinStick {
	this := &RecordableTwinStick{
		TwinStick:       NewTwinStick(),
		RecordableInput: record.NewRecordableInput(),
	}
	return this
}

func (this RecordableTwinStick) GetState(doRecord bool /*= true */) TwinStickState {
	s := this.TwinStick.GetState()
	if doRecord {
		this.Record(&s)
	}
	return s
}
