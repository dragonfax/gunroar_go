package main

/**
 * Save/Load a replay data.
 */

const dir = "replay"
const REPLAY_VERSION_NUM = 11

type ReplayData struct {
	// jInputRecord!(PadState) padInputRecord;
	twinStickInputRecord sdl.InputRecord
	// InputRecord!(MouseAndPadState) mouseAndPadInputRecord;
	seed            int64
	score           int
	shipTurnSpeed   float64
	shipReverseFire bool
	gameMode        int
}

func NewReplayData() *ReplayData {
	this := &ReplayData{}
	return this
}

func (this *ReplayData) save(fileName string) {
	fd := NewFile()
	fd.create(dir + "/" + fileName)
	fd.writeInt(REPLAY_VERSION_NUM)
	fd.writeInt(this.seed)
	fd.writeInt(this.score)
	fd.writeInt(this.shipTurnSpeed)
	if this.shipReverseFire {
		fd.writeInt(1)
	} else {
		fd.writeInt(0)
	}
	fd.write(this.gameMode)
	switch this.gameMode {
	/* case InGameState.GameMode.NORMAL:
	   padInputRecord.save(fd);
	*/
	case InGameState.GameMode.TWIN_STICK, InGameState.GameMode.DOUBLE_PLAY:
		this.twinStickInputRecord.save(fd)
		/* case InGameState.GameMode.MOUSE:
		   this.mouseAndPadInputRecord.save(fd);
		*/
	}
	fd.close()
}

func (this *ReplayData) load(fileName string) {
	fd := file.NewFile()
	fd.open(dir + "/" + fileName)
	ver := fd.readInt()
	if ver != REPLAY_VERSION_NUM {
		panic("Wrong version num")
	}
	this.seed = fd.readInt()
	this.score = fd.readInt()
	this.shipTurnSpeed = fd.readInt()
	srf := fd.readInt()
	if srf == 1 {
		shipReverseFire = true
	} else {
		shipReverseFire = false
	}
	this.gameMode = fd.readInt()
	switch gameMode {
	/* case InGameState.GameMode.NORMAL:
	   padInputRecord = new InputRecord!(PadState);
	   padInputRecord.load(fd);
	*/
	case InGameState.GameMode.TWIN_STICK, InGameState.GameMode.DOUBLE_PLAY:
		this.twinStickInputRecord = NewInputRecord(TwinStickState)
		this.twinStickInputRecord.load(fd)
		/* case InGameState.GameMode.MOUSE:
		   mouseAndPadInputRecord = new InputRecord!(MouseAndPadState);
		   mouseAndPadInputRecord.load(fd);
		*/
	}
	fd.close()
}
