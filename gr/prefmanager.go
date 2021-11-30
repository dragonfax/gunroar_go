package main

import "github.com/dragonfax/gunroar/gr/sdl/file"

/**
 * Save/Load the high score.
 */

const PREFS_VERSION_NUM = 14
const PREF_FILE = "gr.prf"

type PrefManager struct {
	_prefData *PrefData
}

func NewPrefManager() *PrefManager {
	this := &PrefManager{}
	return this
}

func (this *PrefManager) Load() {
	fd := file.New()
	err := fd.Open(PREF_FILE)
	if err != nil {
		ver := fd.ReadInt()
		if ver != PREFS_VERSION_NUM {
			panic("Wrong version num")
		} else {
			this._prefData.load(fd)
		}
	} else {
		this._prefData.init()
	}

	if fd.IsOpen() {
		fd.Close()
	}
}

func (this *PrefManager) Save() {
	fd := file.New()
	fd.Create(PREF_FILE)
	fd.WriteInt(PREFS_VERSION_NUM)
	this._prefData.save(fd)
	fd.Close()
}

func (this *PrefManager) prefData() *PrefData {
	return this._prefData
}

type PrefData struct {
	_highScore [4]int
	_gameMode  GameMode
}

func (this *PrefData) init() {
	for i := range this._highScore {
		this._highScore[i] = 0
	}
	this._gameMode = 0
}

func (this *PrefData) load(fd file.File) {
	for i := range this._highScore {
		this._highScore[i] = fd.ReadInt()
	}
	this._gameMode = GameMode(fd.ReadInt())
}

func (this *PrefData) save(fd file.File) {
	for i := range this._highScore {
		fd.WriteInt(this._highScore[i])
	}
	fd.WriteInt(int(this._gameMode))
}

func (this *PrefData) recordGameMode(gm GameMode) {
	this._gameMode = gm
}

func (this *PrefData) recordResult(score int, gm GameMode) {
	if score > this._highScore[gm] {
		this._highScore[gm] = score
	}
	this._gameMode = gm
}

func (this PrefData) highScore(gm GameMode) int {
	return this._highScore[int(gm)]
}

func (this PrefData) gameMode() GameMode {
	return this._gameMode
}
