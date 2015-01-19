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

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
)

var noSound bool

const RANDOM_BGM_START_INDEX = 1

var seFileName = []string{"shot.wav", "lance.wav", "hit.wav",
	"turret_destroyed.wav", "destroyed.wav", "small_destroyed.wav", "explode.wav",
	"ship_destroyed.wav", "ship_shield_lost.wav", "score_up.wav"}

var seChannel = []int{0, 1, 2, 3, 4, 5, 6, 7, 7, 6}
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
	for i, fileName := range seFileName {
		chunk := &Chunk{}
		chunk.LoadToChannel(fileName, seChannel[i])
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
	if sdl.InitSubSystem(sdl.INIT_AUDIO) < 0 {
		noSound = true
		panic(errors.New(fmt.Sprintf("SDLInitFailedException Unable to initialize SDL_AUDIO: %v", sdl.GetError())))
	}
	audio_rate := 44100
	var audio_format uint16 = sdl.AUDIO_S16
	audio_channels := 1
	audio_buffers := 4096
	if !mix.OpenAudio(audio_rate, audio_format, audio_channels, audio_buffers) {
		noSound = true
		panic(errors.New(fmt.Sprintf("SDLInitFailedException Couldn't open audio: %v", sdl.GetError())))
	}
	mix.QuerySpec(&audio_rate, &audio_format, &audio_channels)
}

func CloseSoundManager() {
	if noSound {
		return
	}
	if mix.MusicPlaying() {
		mix.HaltMusic()
	}
	mix.CloseAudio()
}

type Sound interface {
	Load(name string)
	LoadToChannel(name string, ch int)
	Free()
	Play()
	Fade()
	Halt()
}

const fadeOutSpeed = 1280
const musicDir = "sounds/musics"

type Music struct {
	music *mix.Music
}

func (m *Music) Load(name string) {
	if noSound {
		return
	}
	fileName := musicDir + "/" + name
	m.music = mix.LoadMUS(fileName)
	if m.music == nil {
		noSound = true
		panic("Couldn't load: " + fileName)
	}
}

func (m *Music) LoadToChannel(name string, ch int) {
	m.Load(name)
}

func (m *Music) Free() {
	if m.music != nil {
		m.Halt()
		m.music.Free()
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

func (m *Music) Fade() {
	fadeMusic()
}

func (m *Music) Halt() {
	haltMusic()
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
	if mix.MusicPlaying() {
		mix.HaltMusic()
	}
}

const soundDir = "sounds/chunks"

type Chunk struct {
	chunk        *mix.Chunk
	chunkChannel int
}

func (c *Chunk) Load(name string) {
	c.LoadToChannel(name, 0)
}

func (c *Chunk) LoadToChannel(name string, ch int) {
	if noSound {
		return
	}
	fileName := soundDir + "/" + name
	c.chunk = mix.LoadWAV(fileName)
	if c.chunk == nil {
		noSound = true
		panic("SDLException Couldn't load: " + fileName)
	}
	c.chunkChannel = ch
}

func (c *Chunk) Free() {
	if c.chunk != nil {
		c.Halt()
		c.chunk.Free()
	}
}

func (c *Chunk) Play() {
	if noSound {
		return
	}
	c.chunk.PlayChannel(c.chunkChannel, 0)
}

func (c *Chunk) Halt() {
	if noSound {
		return
	}
	mix.HaltChannel(c.chunkChannel)
}

func (c *Chunk) Fade() {
	c.Halt()
}
