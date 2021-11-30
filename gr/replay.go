package main

import (
	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/sdl/file"
	"github.com/dragonfax/gunroar/gr/sdl/record"
)

/**
 * Save/Load a replay data.
 */

const dir = "replay"
const REPLAY_VERSION_NUM = 11

type ReplayData struct {
	padInputRecord       record.InputRecord
	twinStickInputRecord record.InputRecord
	// InputRecord!(MouseAndPadState) mouseAndPadInputRecord;
	seed            int64
	score           int
	shipTurnSpeed   float64
	shipReverseFire bool
	gameMode        GameMode
}

func NewReplayData() *ReplayData {
	this := &ReplayData{}
	return this
}

func (this *ReplayData) save(fileName string) {
	fd := file.New()
	fd.Create(dir + "/" + fileName)
	fd.WriteInt(REPLAY_VERSION_NUM)
	fd.WriteInt64(this.seed)
	fd.WriteInt(this.score)
	fd.WriteFloat64(this.shipTurnSpeed)
	if this.shipReverseFire {
		fd.WriteInt(1)
	} else {
		fd.WriteInt(0)
	}
	fd.WriteInt(int(this.gameMode))
	switch this.gameMode {
	/* case InGameState.GameMode.NORMAL:
	   padInputRecord.save(fd);
	*/
	case TWIN_STICK, DOUBLE_PLAY:
		this.twinStickInputRecord.Save(fd)
		/* case InGameState.GameMode.MOUSE:
		this.mouseAndPadInputRecord.save(fd);
		*/
	}
	fd.Close()
}

func (this *ReplayData) load(fileName string) {
	fd := file.New()
	fd.Open(dir + "/" + fileName)
	ver := fd.ReadInt()
	if ver != REPLAY_VERSION_NUM {
		panic("Wrong version num")
	}
	this.seed = fd.ReadInt64()
	this.score = fd.ReadInt()
	this.shipTurnSpeed = fd.ReadFloat64()
	srf := fd.ReadInt()
	if srf == 1 {
		shipReverseFire = true
	} else {
		shipReverseFire = false
	}
	this.gameMode = GameMode(fd.ReadInt())
	switch this.gameMode {
	/* case InGameState.GameMode.NORMAL:
	   padInputRecord = new InputRecord!(PadState);
	   padInputRecord.load(fd);
	*/
	case TWIN_STICK, DOUBLE_PLAY:
		this.twinStickInputRecord = record.New(sdl.NewTwinStickState)
		this.twinStickInputRecord.Load(fd)
		/* case InGameState.GameMode.MOUSE:
		   mouseAndPadInputRecord = new InputRecord!(MouseAndPadState);
		   mouseAndPadInputRecord.load(fd);
		*/
	}
	fd.Close()
}
