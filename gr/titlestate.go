package main

/* title screen state */

var titleState *TitleState

type TitleState struct {
	gameOverCnt int
}

func NewTitleState() *TitleState {
	this := new(TitleState)
	return this
}

func (this *TitleState) close() {
	titleManager.close()
}

func (this *TitleState) start() {
	haltBgm()
	disableBgm()
	disableSe()
	titleManager.start()
}

func (this *TitleState) move() {
	titleManager.move()
}

func (this *TitleState) draw() {
	field.draw()
}

func (this *TitleState) drawFront() {
}

func (this *TitleState) drawOrtho() {
	titleManager.draw()
}

func (this *TitleState) drawLuminous() {
	inGameState.drawLuminous()
}
