package sdl

/**
 * SDL screen handler interface.
 */
type Screen interface {
	InitSDL()
	CloseSDL()
	Flip()
	Clear()
}

type SizableScreen interface {
	WindowMode() bool
	Width() int
	Height() int
}
