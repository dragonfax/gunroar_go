package sdl

import (
	"fmt"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var NoSound = false

/**
 * Initialize and close SDL_mixer.
 */

func SoundManagerInit() {
	if NoSound {
		return
	}
	var audio_rate int
	var audio_format uint16
	var audio_channels int
	var audio_buffers int
	err := sdl.InitSubSystem(sdl.INIT_AUDIO)
	if err != nil {
		NoSound = true
		panic("Unable to initialize SDL_AUDIO: " + err.Error())
	}
	audio_rate = 44100
	audio_format = mix.DEFAULT_FORMAT // mix.AUDIO_S16
	audio_channels = 1
	audio_buffers = 4096
	err = mix.OpenAudio(audio_rate, audio_format, audio_channels, audio_buffers)
	if err != nil {
		NoSound = true
		panic("Couldn't open audio: " + err.Error())
	}
	// audio_rate, audio_format, audio_channels, audio_opened, err := mix.QuerySpec() # serves no purpose
	fmt.Println("audio opened")
}

/**
 * Music / Chunk.
 */
type Sound interface {
	Load(name string)
	LoadWithChannel(name string, ch int)
	Free()
	Play()
	Fade()
	Halt()
}

var _ Sound = &Music{}

var fadeOutSpeed = 1280

const MusicDir = "sounds/musics"

type Music struct {
	music *mix.Music
}

func (this *Music) Load(name string) {
	if NoSound {
		return
	}
	fileName := name
	m, err := mix.LoadMUS(fileName)
	this.music = m
	if err != nil {
		panic(err)
	}
	if this.music == nil {
		NoSound = true
		panic("music not loaded")
	}
}

func (this *Music) LoadWithChannel(name string, ch int) {
	this.Load(name)
}

func (this *Music) Free() {
	if this.music != nil {
		this.Halt()
		this.music.Free()
	}
}

func (this *Music) Play() {
	if NoSound {
		return
	}
	this.music.Play(-1)
}

func (this *Music) PlayOnce() {
	if NoSound {
		return
	}
	this.music.Play(1)
}

func (this *Music) Fade() {
	FadeMusic()
}

func (this *Music) Halt() {
	HaltMusic()
}

func FadeMusic() {
	if NoSound {
		return
	}
	mix.FadeOutMusic(fadeOutSpeed)
}

func HaltMusic() {
	if NoSound {
		return
	}
	if mix.PlayingMusic() {
		mix.HaltMusic()
	}
}

var _ Sound = &Chunk{}

const sound_dir = "sounds/chunks"

type Chunk struct {
	chunk        *mix.Chunk
	chunkChannel int
}

func (this *Chunk) Load(name string) {
	this.LoadWithChannel(name, 0)
}

func (this *Chunk) LoadWithChannel(name string, ch int) {
	if NoSound {
		return
	}
	fileName := sound_dir + "/" + name
	c, err := mix.LoadWAV(fileName)
	this.chunk = c
	if err != nil {
		panic(err)
	}
	if this.chunk == nil {
		NoSound = true
		panic("no chunk loaded")
	}
	this.chunkChannel = ch
}

func (this *Chunk) Free() {
	if this.chunk != nil {
		this.Halt()
		this.chunk.Free()
	}
}

func (this *Chunk) Play() {
	if NoSound {
		return
	}
	this.chunk.Play(this.chunkChannel, 0)
}

func (this *Chunk) Halt() {
	if NoSound {
		return
	}
	mix.HaltChannel(this.chunkChannel)
}

func (this *Chunk) Fade() {
	this.Halt()
}
