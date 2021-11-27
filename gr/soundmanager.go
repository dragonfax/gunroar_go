package main

import (
	r "math/rand"
	"time"

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
var soundRand *r.Rand
var bgmFileName []string
var currentBgm string
var prevBgmIdx int
var nextIdxMv int

func setSoundRandSeed(seed int64) {
	rand = r.New(r.NewSource(seed))
}

func loadSounds() {
	loadMusics()
	loadChunks()
	rand = r.New(r.NewSource(time.Now().Unix()))
}

func loadMusics() {
	files := listdir(sdl.MusicDir)
	for _, fileName := range files {
		ext := getExt(fileName)
		if ext != "ogg" && ext != "wav" {
			continue
		}
		music := &sdl.Music{}
		music.Load(fileName)
		bgm[fileName] = music
		bgmFileName = append(bgmFileName, fileName)
		Logger.info("Load bgm: " + fileName)
	}
}

func loadChunks() {
	i := 0
	for _, fileName := range seFileName {
		chunk := &sdl.Chunk{}
		chunk.LoadWithChannel(fileName, seChannel[i])
		se[fileName] = chunk
		seMark[fileName] = false
		Logger.info("Load SE: " + fileName)
		i++
	}
}

func playBgmWithName(name string) {
	currentBgm = name
	if bgmDisabled {
		return
	}
	sdl.Music.HaltMusic()
	bgm[name].Play()
}

func playBgm() {
	bgmIdx := rand.nextInt(bgm.length-RANDOM_BGM_START_INDEX) + RANDOM_BGM_START_INDEX
	nextIdxMv = rand.nextInt(2)*2 - 1
	prevBgmIdx = bgmIdx
	playBgm(bgmFileName[bgmIdx])
}

func nextBgm() {
	bgmIdx := prevBgmIdx + nextIdxMv
	if bgmIdx < RANDOM_BGM_START_INDEX {
		bgmIdx = bgm.length - 1
	} else if bgmIdx >= bgm.length {
		bgmIdx = RANDOM_BGM_START_INDEX
	}
	prevBgmIdx = bgmIdx
	playBgm(bgmFileName[bgmIdx])
}

func playCurrentBgm() {
	playBgm(currentBgm)
}

func fadeBgm() {
	sdl.Music.fadeMusic()
}

func haltBgm() {
	sdl.Music.HaltMusic()
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
			se[key].play()
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
