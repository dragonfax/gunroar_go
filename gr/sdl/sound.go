package sdl

import (
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var noSound = false

/**
 * Initialize and close SDL_mixer.
 */

func SoundManagerInit() {
	if noSound {
		return
	}
	var audio_rate int
	var audio_format uint16
	var audio_channels int
	var audio_buffers int
	err := sdl.InitSubSystem(sdl.INIT_AUDIO)
	if err != nil {
		noSound = true
		panic("Unable to initialize SDL_AUDIO: " + err.Error())
	}
	audio_rate = 44100
	audio_format = mix.AUDIO_S16
	audio_channels = 1
	audio_buffers = 4096
	err = mix.OpenAudio(audio_rate, audio_format, audio_channels, audio_buffers)
	if err != nil {
		noSound = true
		panic("Couldn't open audio: " + err.Error())
	}
	mix.QuerySpec(&audio_rate, &audio_format, &audio_channels)
}

func close() {
	if noSound {
		return
	}
	if mix.PlayingMusic() {
		mix.HaltMusic()
	}
	mix.CloseAudio()
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
	if noSound {
		return
	}
	fileName := MusicDir + "/" + name
	this.music = mix.LoadMUS(fileName)
	if this.music == nil {
		noSound = true
		panic("Couldn't load: " + fileName +
			" (" + mix.GetError().Error() + ")")
	}
}

func (this *Music) LoadWithChannel(name string, ch int) {
	this.Load(name)
}

func (this *Music) Free() {
	if this.music != nil {
		this.Halt()
		mix.FreeMusic(this.music)
	}
}

func (this *Music) Play() {
	if noSound {
		return
	}
	mix.PlayMusic(this.music, -1)
}

func (this *Music) PlayOnce() {
	if noSound {
		return
	}
	mix.PlayMusic(this.music, 1)
}

func (this *Music) Fade() {
	fadeMusic()
}

func (this *Music) Halt() {
	HaltMusic()
}

func fadeMusic() {
	if noSound {
		return
	}
	mix.FadeOutMusic(fadeOutSpeed)
}

func HaltMusic() {
	if noSound {
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
	if noSound {
		return
	}
	fileName := sound_dir + "/" + name
	this.chunk = mix.LoadWAV(fileName)
	if this.chunk == nil {
		noSound = true
		panic("Couldn't load: " + fileName +
			" (" + mix.GetError().Error() + ")")
	}
	this.chunkChannel = ch
}

func (this *Chunk) Free() {
	if this.chunk != nil {
		this.Halt()
		mix.FreeChunk(this.chunk)
	}
}

func (this *Chunk) Play() {
	if noSound {
		return
	}
	mix.PlayChannel(this.chunkChannel, this.chunk, 0)
}

func (this *Chunk) Halt() {
	if noSound {
		return
	}
	mix.HaltChannel(this.chunkChannel)
}

func (this *Chunk) Fade() {
	this.Halt()
}
