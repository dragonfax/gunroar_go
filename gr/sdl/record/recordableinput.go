package record

import (
	"fmt"

	"github.com/dragonfax/gunroar/gr/sdl/file"
)

type InputState interface {
	Read(file.File)
	Write(file.File)
	Equals(InputState) bool
	Set(InputState)
}

/**
 * Record an input for a replay.
 * T represents a data structure of specific device input.
 */
type RecordableInput struct {
	inputRecord *InputRecord
}

func NewRecordableInput() RecordableInput {
	this := RecordableInput{}
	return this
}

func (this *RecordableInput) startRecord() {
	this.inputRecord = new(InputRecord)
	this.inputRecord.clear()
}

func (this *RecordableInput) Record(d InputState) {
	this.inputRecord.add(d)
}

func (this *RecordableInput) startReplay(pr *InputRecord) {
	this.inputRecord = pr
	this.inputRecord.reset()
}

var EndRecordingErr = fmt.Errorf("end of recording")

func (this *RecordableInput) Replay() (InputState, error) {
	if !this.inputRecord.hasNext() {
		return nil, EndRecordingErr
	}
	return this.inputRecord.next(), nil
}

type InputStateConstructor func(InputState) InputState

type Record struct {
	series int
	data   InputState
}

type InputRecord struct {
	record           []Record
	idx, series      int
	replayData       InputState
	stateConstructor InputStateConstructor
}

func New(constructor InputStateConstructor) InputRecord {
	this := InputRecord{
		stateConstructor: constructor,
		record:           make([]Record, 0),
		replayData:       constructor(nil),
	}
	return this
}

func (this *InputRecord) clear() { //lint:ignore ST1006
	this.record = make([]Record, 0)
}

func (this *InputRecord) add(d InputState) {
	if len(this.record) > 0 && this.record[len(this.record)-1].data.Equals(d) {
		this.record[len(this.record)-1].series++
	} else {
		var r Record
		r.series = 1
		r.data = this.stateConstructor(d)
		this.record = append(this.record, r)
	}
}

func (this *InputRecord) reset() {
	this.idx = 0
	this.series = 0
}

func (this *InputRecord) hasNext() bool {
	return this.idx < len(this.record)
}

func (this *InputRecord) next() InputState {
	if this.idx >= len(this.record) {
		panic("No more items")
	}
	if this.series <= 0 {
		this.series = this.record[this.idx].series
	}
	this.replayData.Set(this.record[this.idx].data)
	this.series--
	if this.series <= 0 {
		this.idx++
	}
	return this.replayData
}

func (this *InputRecord) Save(fd file.File) {
	fd.WriteInt(len(this.record))
	for _, r := range this.record {
		fd.WriteInt(r.series)
		r.data.Write(fd)
	}
}

func (this *InputRecord) Load(fd file.File) {
	this.clear()
	l := fd.ReadInt()
	for i := 0; i < l; i++ {
		s := fd.ReadInt()
		d := this.stateConstructor(nil)
		d.Read(fd)
		var r Record
		r.series = s
		r.data = d
		this.record = append(this.record, r)
	}
}
