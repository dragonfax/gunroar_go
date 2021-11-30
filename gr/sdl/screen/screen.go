package screen

/**
 * SDL screen handler interface.
 */
type Screen interface {
	InitSDL()
	CloseSDL()
	Flip()
	Clear()
	HandleError()
}

type SizableScreen interface {
	WindowMode() bool
	Width() int
	Height() int
}
