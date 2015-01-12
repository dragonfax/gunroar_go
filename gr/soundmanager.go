/*
 * $Id: soundmanager.d,v 1.5 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

import (
	"github.com/veandco/go-sdl2/sdl_mixer"
)

var noSound bool

const RANDOM_BGM_START_INDEX = 1

var seFileName = []string{"shot.wav", "lance.wav", "hit.wav",
	"turret_destroyed.wav", "destroyed.wav", "small_destroyed.wav", "explode.wav",
	"ship_destroyed.wav", "ship_shield_lost.wav", "score_up.wav"}

var seChannel = []int{0, 1, 2, 3, 4, 5, 6, 7, 7, 6}
var bgm map[string]Music
var se map[string]Chunk
var seMark map[string]bool
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
	musics := make(map[string]Music)
	files := listdir(Music.dir)
	for _, fileName := range files {
		ext = getExt(fileName)
		if ext != "ogg" && ext != "wav" {
			continue
		}
		music := &Music{}
		music.Load(fileName)
		bgm[fileName] = music
		bgmFileName += fileName
		fmt.Println("Load bgm: " + fileName)
	}
}

func loadChunks() {
	for i, fileName := range seFileName {
		chunk := &Chunk{}
		chunk.Load(fileName, seChannel[i])
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
	Music.haltMusic()
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
	Music.fadeMusic()
}

func haltBgm() {
	Music.haltMusic()
}

func playSe(name string) {
	if seDisabled {
		return
	}
	seMark[name] = true
}

func playMarkedSe() {
	keys := seMark.keys
	for key, _ := range keys {
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
	audio_format := AUDIO_S16
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
	if music == nil {
		noSound = true
		panic("Couldn't load: " + fileName + " (" + Mix_GetError() + ")")
	}
}

func (m *Music) LoadToChannel(name string, ch int) {
	m.load(name)
}

func (m *Music) Free() {
	if m.music != nil {
		m.Halt()
		mix.FreeMusic(m.music)
	}
}

func (m *Music) Play() {
	if noSound {
		return
	}
	mix.PlayMusic(m.music, -1)
}

func (m *Music) PlayOnce() {
	if noSound {
		return
	}
	mix.PlayMusic(m.music, 1)
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
	if mix.PlayingMusic() {
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
	if chunk == nil {
		noSound = true
		panic("SDLException Couldn't load: " + fileName + " (" + mix.GetError() + ")")
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
