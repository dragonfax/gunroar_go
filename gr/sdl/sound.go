package sdl

import (
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var noSound = false

/**
 * Initialize and close SDL_mixer.
 */
type SoundManager struct {
}

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
	load(name string)
	loadWithChannel(name string, ch int)
	free()
	play()
	fade()
	halt()
}

var _ Sound = &Music{}

var fadeOutSpeed = 1280

const music_dir = "sounds/musics"

type Music struct {
	music *mix.Music
}

func (this *Music) load(name string) {
	if noSound {
		return
	}
	fileName := music_dir + "/" + name
	this.music = mix.LoadMUS(fileName)
	if this.music == nil {
		noSound = true
		panic("Couldn't load: " + fileName +
			" (" + mix.GetError().Error() + ")")
	}
}

func (this *Music) loadWithChannel(name string, ch int) {
	this.load(name)
}

func (this *Music) free() {
	if this.music != nil {
		this.halt()
		mix.FreeMusic(this.music)
	}
}

func (this *Music) play() {
	if noSound {
		return
	}
	mix.PlayMusic(this.music, -1)
}

func (this *Music) playOnce() {
	if noSound {
		return
	}
	mix.PlayMusic(this.music, 1)
}

func (this *Music) fade() {
	fadeMusic()
}

func (this *Music) halt() {
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

var _ Sound = &Chunk{}

const sound_dir = "sounds/chunks"

type Chunk struct {
	chunk        *mix.Chunk
	chunkChannel int
}

func (this *Chunk) load(name string) {
	this.loadWithChannel(name, 0)
}

func (this *Chunk) loadWithChannel(name string, ch int) {
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

func (this *Chunk) free() {
	if this.chunk != nil {
		this.halt()
		mix.FreeChunk(this.chunk)
	}
}

func (this *Chunk) play() {
	if noSound {
		return
	}
	mix.PlayChannel(this.chunkChannel, this.chunk, 0)
}

func (this *Chunk) halt() {
	if noSound {
		return
	}
	mix.HaltChannel(this.chunkChannel)
}

func (this *Chunk) fade() {
	this.halt()
}
