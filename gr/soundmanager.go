package main

import (
	"fmt"
	"math/rand"
	"path/filepath"

	"github.com/dragonfax/gunroar/gr/sdl"
)

/**
 * Manage BGMs and SEs.
 */

const RANDOM_BGM_START_INDEX = 1

var seFileName = []string{
	"shot.wav", "lance.wav", "hit.wav",
	"turret_destroyed.wav", "destroyed.wav", "small_destroyed.wav", "explode.wav",
	"ship_destroyed.wav", "ship_shield_lost.wav", "score_up.wav"}
var seChannel = []int{0, 1, 2, 3, 4, 5, 6, 7, 7, 6}
var bgm map[string]*sdl.Music
var se map[string]*sdl.Chunk
var seMark map[string]bool
var bgmDisabled = false
var seDisabled = false
var bgmFileName []string
var currentBgm string
var prevBgmIdx int
var nextIdxMv int

func loadSounds() {
	loadMusics()
	loadChunks()
}

func loadMusics() {
	files, err := filepath.Glob(sdl.MusicDir + "/*")
	if err != nil {
		panic(err)
	}
	for _, fileName := range files {
		ext := filepath.Ext(fileName)
		if ext != "ogg" && ext != "wav" {
			continue
		}
		music := &sdl.Music{}
		music.Load(fileName)
		bgm[fileName] = music
		bgmFileName = append(bgmFileName, fileName)
		fmt.Println("Load bgm: " + fileName)
	}
}

func loadChunks() {
	i := 0
	for _, fileName := range seFileName {
		chunk := &sdl.Chunk{}
		chunk.LoadWithChannel(fileName, seChannel[i])
		se[fileName] = chunk
		seMark[fileName] = false
		fmt.Println("Load SE: " + fileName)
		i++
	}
}

func playBgmWithName(name string) {
	currentBgm = name
	if bgmDisabled {
		return
	}
	sdl.HaltMusic()
	bgm[name].Play()
}

func playBgm() {
	bgmIdx := rand.Intn(len(bgm)-RANDOM_BGM_START_INDEX) + RANDOM_BGM_START_INDEX
	nextIdxMv = rand.Intn(2)*2 - 1
	prevBgmIdx = bgmIdx
	playBgmWithName(bgmFileName[bgmIdx])
}

func nextBgm() {
	bgmIdx := prevBgmIdx + nextIdxMv
	if bgmIdx < RANDOM_BGM_START_INDEX {
		bgmIdx = len(bgm) - 1
	} else if bgmIdx >= len(bgm) {
		bgmIdx = RANDOM_BGM_START_INDEX
	}
	prevBgmIdx = bgmIdx
	playBgmWithName(bgmFileName[bgmIdx])
}

func playCurrentBgm() {
	playBgmWithName(currentBgm)
}

func fadeBgm() {
	sdl.FadeMusic()
}

func haltBgm() {
	sdl.HaltMusic()
}

func playSe(name string) {
	if seDisabled {
		return
	}
	seMark[name] = true
}

func playMarkedSe() {
	for key := range seMark {
		if seMark[key] {
			se[key].Play()
			seMark[key] = false
		}
	}
}

func disableSe() {
	seDisabled = true
}

func enableSe() {
	seDisabled = false
}

func disableBgm() {
	bgmDisabled = true
}

func enableBgm() {
	bgmDisabled = false
}
