/*
 * $Id: soundmanager.d,v 1.5 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var noSound bool

const RANDOM_BGM_START_INDEX = 1
const soundDir = "sounds/chunks"
const fadeOutSpeed = 1280
const musicDir = "sounds/musics"

var bgm map[string]*Music = make(map[string]*Music)
var se map[string]*Chunk = make(map[string]*Chunk)
var seMark map[string]bool = make(map[string]bool)
var bgmDisabled bool
var seDisabled bool
var bgmFileName []string
var currentBgm string
var prevBgmIdx int
var nextIdxMv int

func loadSounds() {
	loadMusics()
	loadChunks()
}

func loadMusics() {
	files, err := ioutil.ReadDir(musicDir)
	if err != nil {
		panic(err.Error())
	}
	if len(files) == 0 {
		panic("no bgms found")
	}
	for _, fileInfo := range files {
		fileName := fileInfo.Name()
		ext := filepath.Ext(fileName)
		if ext != ".ogg" && ext != ".wav" {
			fmt.Printf("skipping extension %s\n", ext)
			continue
		}
		music := &Music{}
		music.Load(fileName)
		bgm[fileName] = music
		bgmFileName = append(bgmFileName, fileName)
		fmt.Println("Load bgm: " + fileName)
	}
}

func loadChunks() {
	files, err := ioutil.ReadDir(soundDir)
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		panic("no sound effects found")
	}
	for _, fileInfo := range files {
		fileName := fileInfo.Name()
		chunk := &Chunk{}
		chunk.LoadToChannel(fileName)
		se[fileName] = chunk
		seMark[fileName] = false
		fmt.Println("Load SE: " + fileName)
	}
}

func playBgmByName(name string) {
	currentBgm = name
	if bgmDisabled {
		return
	}
	haltMusic()
	bgm[name].Play()
}

func playBgm() {
	bgmIdx := nextInt(len(bgm)-RANDOM_BGM_START_INDEX) + RANDOM_BGM_START_INDEX
	nextIdxMv = nextInt(2)*2 - 1
	prevBgmIdx = bgmIdx
	playBgmByName(bgmFileName[bgmIdx])
}

func nextBgm() {
	bgmIdx := prevBgmIdx + nextIdxMv
	if bgmIdx < RANDOM_BGM_START_INDEX {
		bgmIdx = len(bgm) - 1
	} else if bgmIdx >= len(bgm) {
		bgmIdx = RANDOM_BGM_START_INDEX
	}
	prevBgmIdx = bgmIdx
	playBgmByName(bgmFileName[bgmIdx])
}

func playCurrentBgm() {
	playBgmByName(currentBgm)
}

func fadeBgm() {
	fadeMusic()
}

func haltBgm() {
	haltMusic()
}

func playSe(name string) {
	if seDisabled {
		return
	}
	seMark[name] = true
}

func playMarkedSe() {
	for key, _ := range seMark {
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

func InitSoundManager() {
	if noSound {
		return
	}
	if err := sdl.InitSubSystem(sdl.INIT_AUDIO); err != nil {
		noSound = true
		panic(errors.New(fmt.Sprintf("SDLInitFailedException Unable to initialize SDL_AUDIO: %s", err)))
	}
	audio_rate := 44100
	var audio_format uint16 = sdl.AUDIO_S16
	audio_channels := 1
	audio_buffers := 4096
	if err := mix.OpenAudio(audio_rate, audio_format, audio_channels, audio_buffers); err != nil {
		noSound = true
		panic(errors.New(fmt.Sprintf("SDLInitFailedException Couldn't open audio: %s", err)))
	}
	var err error
	audio_rate, audio_format, audio_channels, _, err = mix.QuerySpec()
	if err != nil {
		panic(err)
	}
}

func CloseSoundManager() {
	if noSound {
		return
	}
	if mix.Playing(-1) > 0 {
		mix.HaltMusic()
	}
	mix.CloseAudio()
}

type Music struct {
	music *mix.Music
}

func (m *Music) Load(name string) {
	if noSound {
		return
	}
	fileName := musicDir + "/" + name
	mus, err := mix.LoadMUS(fileName)
	if err != nil {
		panic(err)
	}
	m.music = mus
	if m.music == nil {
		noSound = true
		panic("Couldn't load: " + fileName)
	}
}

func (m *Music) Play() {
	if noSound {
		return
	}
	m.music.Play(-1)
}

func (m *Music) PlayOnce() {
	if noSound {
		return
	}
	m.music.Play(1)
}

func fadeMusic() {
	if noSound {
		return
	}
	mix.FadeOutMusic(fadeOutSpeed)
}

func haltMusic() {
	if noSound {
		return
	}
	if mix.Playing(-1) > 0 {
		mix.HaltMusic()
	}
}

type Chunk struct {
	chunk *mix.Chunk
}

func (c *Chunk) LoadToChannel(name string) {
	if noSound {
		return
	}
	fileName := soundDir + "/" + name
	chunk, err := mix.LoadWAV(fileName)
	if err != nil {
		panic(err)
	}
	c.chunk = chunk
	if c.chunk == nil {
		noSound = true
		panic("SDLException Couldn't load: " + fileName)
	}
}

func (c *Chunk) Play() {
	if noSound {
		return
	}
	c.chunk.Play(-1, 0)
}
