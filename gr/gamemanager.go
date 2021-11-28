package main

import (
	"fmt"
	"math"
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/sdl"
)

/**
 * Manage the game state and actor pools.
 */
var _ sdl.GameManager = &GameManager{}

var shipTurnSpeed = 1.0
var shipReverseFire = false

type GameManager struct {
	sdl.GameManagerBase

	twinStick      TwinStick
	prefManager    PrefManager
	screen         Screen
	field          Field
	ship           *Ship
	shots          ShotPool
	bullets        BulletPool
	enemies        EnemyPool
	sparks         SparkPool
	smokes         SmokePool
	fragments      FragmentPool
	sparkFragments SparkFragmentPool
	wakes          WakePool
	crystals       CrystalPool
	numIndicators  NumIndicatorPool
	stageManager   StageManager
	titleManager   TitleManager
	scoreReel      ScoreReel
	state          GameState
	titleState     TitleState
	inGameState    InGameState
	escPressed     bool
}

func NewGameManager() *GameManager {
	this := &GameManager{GameManagerBase: sdl.NewGameManagerBaseInternal()}
	return this
}

func (this *GameManager) init() {
	letter.Init()
	Shot.init()
	BulletShape.init()
	EnemyShape.init()
	Turret.init()
	TurretShape.init()
	Fragment.init()
	SparkFragment.init()
	Crystal.init()
	this.prefManager = abstPrefManager
	this.screen = abstScreen
	this.twinStick = input.inputs[1]
	this.twinStick.openJoystick(pad.openJoystick())
	this.field = NewField()
	pargs := make([]interface{}, 0)
	this.sparks = NewSparkPool(120, pargs)
	pargs = append(pargs, field)
	this.wakes = NewWakePool(100, pargs)
	pargs = append(pargs, wakes)
	this.smokes = NewSmokePool(200, pargs)
	fargs := []interface{}{field, smokes}
	this.fragments = NewFragmentPool(60, fargs)
	this.sparkFragments = NewSparkFragmentPool(40, fargs)
	this.ship = NewShip(twinStick,
		field, screen, sparks, smokes, fragments, wakes)
	cargs := []interface{}{ship}
	crystals := NewCrystalPool(80, cargs)
	this.scoreReel = NewScoreReel()
	nargs := []interface{}{scoreReel}
	numIndicators := NewNumIndicatorPool(50, nargs)
	bargs := []interface{}{this, field, ship, smokes, wakes, crystals}
	this.bullets = NewBulletPool(240, bargs)
	eargs := []interface{}{field, screen, bullets, ship, sparks, smokes, fragments, sparkFragments, numIndicators, scoreReel}
	this.enemies = NewEnemyPool(40, eargs)
	sargs := []interface{}{field, enemies, sparks, smokes, bullets}
	this.shots = NewShotPool(50, sargs)
	this.ship.setShots(shots)
	this.ship.setEnemies(enemies)
	this.stageManager = NewStageManager(field, enemies, ship, bullets, sparks, smokes, fragments, wakes)
	this.ship.setStageManager(stageManager)
	this.field.setStageManager(stageManager)
	this.field.setShip(ship)
	this.enemies.setStageManager(stageManager)
	loadSounds()
	this.titleManager = NewTitleManager(prefManager, field, this)
	this.inGameState = NewInGameState(this, screen, twinStick,
		field, ship, shots, bullets, enemies,
		sparks, smokes, fragments, sparkFragments, wakes,
		crystals, numIndicators, stageManager, scoreReel,
		prefManager)
	this.titleState = NewTitleState(this, screen, twinStick,
		field, ship, shots, bullets, enemies,
		sparks, smokes, fragments, sparkFragments, wakes,
		crystals, numIndicators, stageManager, scoreReel,
		titleManager, inGameState)
	this.ship.setGameState(this.inGameState)
}

func (this *GameManager) close() {
	this.ship.close()
	this.BulletShape.close()
	this.EnemyShape.close()
	this.TurretShape.close()
	this.Fragment.close()
	this.SparkFragment.close()
	this.Crystal.close()
	this.titleState.close()
	this.Letter.close()
}

func (this *GameManager) start() {
	this.loadLastReplay()
	this.startTitle()
}

func (this *GameManager) startTitle(fromGameover bool /* = false */) {
	if fromGameover {
		this.saveLastReplay()
	}
	this.titleState.replayData = this.inGameState.replayData
	this.state = this.titleState
	this.startState()
}

func (this *GameManager) startInGame(gameMode int) {
	this.state = this.inGameState
	this.inGameState.gameMode = gameMode
	this.startState()
}

func (this *GameManager) startState() {
	this.state.start()
}

func (this *GameManager) saveErrorReplay() {
	if this.state == this.inGameState {
		this.inGameState.saveReplay("error.rpl")
	}
}

func (this *GameManager) saveLastReplay() {
	err := this.inGameState.saveReplay("last.rpl")
	if err != nil {
		fmt.Printf("warn : %s \n", err.Error())
	}
}

func (this *GameManager) loadLastReplay() {
	err := this.inGameState.loadReplay("last.rpl")
	if err != nil {
		fmt.Printf("warn : %s \n", err.Error())
		this.inGameState.resetReplay()
	}
}

func (this *GameManager) loadErrorReplay() {
	err := this.inGameState.loadReplay("error.rpl")
	if err != nil {
		fmt.printf("warn : %s\n", err.Error())
		this.inGameState.resetReplay()
	}
}

func (this *GameManager) initInterval() {
	mainLoop.initInterval()
}

func (this *GameManager) addSlowdownRatio(sr float64) {
	mainLoop.addSlowdownRatio(sr)
}

func (this *GameManager) move() {
	if pad.keys[SDLK_ESCAPE] == SDL_PRESSED {
		if !escPressed {
			escPressed = true
			if state == inGameState {
				startTitle()
			} else {
				mainLoop.breakLoop()
			}
			return
		}
	} else {
		escPressed = false
	}
	state.move()
}

func (this *GameManager) draw() {
	e := mainLoop.event
	if e.GetType() == sdl.VIDEORESIZE {
		re := e.resize
		if re.w > 150 && re.h > 100 {
			this.screen.resized(re.w, re.h)
		}
	}
	if this.screen.startRenderToLuminousScreen() {
		glPushMatrix()
		this.screen.setEyepos()
		this.state.drawLuminous()
		glPopMatrix()
		this.screen.endRenderToLuminousScreen()
	}
	this.screen.clear()
	glPushMatrix()
	this.screen.setEyepos()
	this.state.draw()
	glPopMatrix()
	this.screen.drawLuminous()
	glPushMatrix()
	this.screen.setEyepos()
	this.field.drawSideWalls()
	this.state.drawFront()
	glPopMatrix()
	this.screen.viewOrthoFixed()
	this.state.drawOrtho()
	this.screen.viewPerspective()
}

/**
 * Manage the game state.
 * (e.g. title, in game, gameover, pause, ...)
 */
type GameState interface {
	start()
	move()
	draw()
	drawLuminous()
	drawFront()
	drawOrtho()
}

type GameStateBase struct {
	gameManager    GameManager
	screen         Screen
	twinStick      TwinStick
	field          Field
	ship           *Ship
	shots          ShotPool
	bullets        BulletPool
	enemies        EnemyPool
	sparks         SparkPool
	smokes         SmokePool
	fragments      FragmentPool
	sparkFragments SparkFragmentPool
	wakes          WakePool
	crystals       CrystalPool
	numIndicators  NumIndicatorPool
	stageManager   StageManager
	scoreReel      ScoreReel
	_replayData    *ReplayData
}

func NewGameStateBase(gameManager GameManager, screen Screen,
	twinStick TwinStick,
	field Field, ship *Ship, shots ShotPool, bullets BulletPool, enemies EnemyPool,
	sparks SparkPool, smokes SmokePool,
	fragments FragmentPool, sparkFragments SparkFragmentPool, wakes WakePool,
	crystals CrystalPool, numIndicators NumIndicatorPool,
	stageManager StageManager, scoreReel ScoreReel) GameStateBase {
	this := GameStateBase{}
	this.gameManager = gameManager
	this.screen = screen
	this.twinStick = twinStick
	this.field = field
	this.ship = ship
	this.shots = shots
	this.bullets = bullets
	this.enemies = enemies
	this.sparks = sparks
	this.smokes = smokes
	this.fragments = fragments
	this.sparkFragments = sparkFragments
	this.wakes = wakes
	this.crystals = crystals
	this.numIndicators = numIndicators
	this.stageManager = stageManager
	this.scoreReel = scoreReel
	return this
}

func (this *GameStateBase) clearAll() {
	this.shots.clear()
	this.bullets.clear()
	this.enemies.clear()
	this.sparks.clear()
	this.smokes.clear()
	this.fragments.clear()
	this.sparkFragments.clear()
	this.wakes.clear()
	this.crystals.clear()
	this.numIndicators.clear()
}

func (this *GameStateBase) setReplayData(ReplayData v) ReplayData {
	this._replayData = v
	return v
}

func (this *GameStateBase) replayData() ReplayData {
	return this._replayData
}

type GameMode int

const (
	TWIN_STICK GameMode = iota
	DOUBLE_PLAY
)
const GAME_MODE_NUM = 2

var gameModeText = []string{"TWIN STICK", "DOUBLE PLAY"}

const SCORE_REEL_SIZE_DEFAULT = 0.5
const SCORE_REEL_SIZE_SMALL = 0.01

var _ GameState = &InGameState{}

type InGameState struct {
	GameStateBase

	isGameOver    bool
	rand          *r.Rand
	prefManager   PrefManager
	left          int
	time          int
	gameOverCnt   int
	btnPressed    bool
	pauseCnt      int
	pausePressed  bool
	scoreReelSize float64
	_gameMode     int
}

func NewInGameState(gameManager GameManager, screen Screen,
	twinStick TwinStick,
	field Field, ship *Ship, shots ShotPool, bullets BulletPool, enemies EnemyPool,
	sparks SparkPool, smokes SmokePool,
	fragments FragmentPool, sparkFragments SparkFragmentPool, wakes WakePool,
	crystals CrystalPool, numIndicators NumIndicatorPool,
	stageManager StageManager, scoreReel ScoreReel,
	prefManager PrefManager) *InGameState {
	this := &InGameState{GameStateBase: NewGameStateBase(
		gameManager, screen, twinStick,
		field, ship, shots, bullets, enemies,
		sparks, smokes, fragments, sparkFragments, wakes, crystals, numIndicators,
		stageManager, scoreReel,
	)}
	this.prefManager = prefManager
	this.rand = r.New(r.NewSource(time.Now().Unix()))
	this.scoreReelSize = SCORE_REEL_SIZE_DEFAULT
	return this
}

func (this *InGameState) start() {
	this.ship.unsetReplayMode()
	this._replayData = NewReplayData()
	this.prefManager.prefData.recordGameMode(this._gameMode)
	switch this._gameMode {
	case GameMode.TWIN_STICK, GameMode.DOUBLE_PLAY:
		rts := twinStick
		rts.startRecord()
		this._replayData.twinStickInputRecord = rts.inputRecord
	}
	this._replayData.seed = this.rand.nextInt32()
	this._replayData.shipTurnSpeed = shipTurnSpeed
	this._replayData.shipReverseFire = shipReverseFire
	this._replayData.gameMode = this._gameMode
	enableBgm()
	enableSe()
	this.startInGame()
}

func (this *InGameState) startInGame() {
	thuis.clearAll()
	seed := this._replayData.seed
	this.field.setRandSeed(seed)
	setEnemyStateRandSeed(seed)
	setEnemySpecRandSeed(seed)
	setTurretRandSeed(seed)
	setSparkRandSeed(seed)
	setSmokeRandSeed(seed)
	setFragmentRandSeed(seed)
	setSparkFragmentRandSeed(seed)
	setScreenRandSeed(seed)
	setBaseShapeRandSeed(seed)
	this.ship.setRandSeed(seed)
	setShotRandSeed(seed)
	this.stageManager.setRandSeed(seed)
	setNumReelRandSeed(seed)
	setNumIndicatorRandSeed(seed)
	setSoundManagerRandSeed(seed)
	this.stageManager.start(1)
	this.field.start()
	this.ship.start(this._gameMode)
	this.initGameState()
	this.screen.setScreenShake(0, 0)
	this.gameOverCnt = 0
	this.pauseCnt = 0
	this.scoreReelSize = SCORE_REEL_SIZE_DEFAULT
	this.isGameOver = false
	playBgm()
}

func (this *InGameState) initGameState() {
	this.time = 0
	this.left = 2
	this.scoreReel.clear(9)
	initTargetY()
}

func (this *InGameState) move() {
	if pad.keys[SDLK_p] == SDL_PRESSED {
		if !pausePressed {
			if pauseCnt <= 0 && !isGameOver {
				this.pauseCnt = 1
			} else {
				this.pauseCnt = 0
			}
		}
		pausePressed = true
	} else {
		pausePressed = false
	}
	if this.pauseCnt > 0 {
		this.pauseCnt++
		return
	}
	this.moveInGame()
	if this.isGameOver {
		this.gameOverCnt++
		if (this.input.button & PadState.Button.A) || (this.gameMode == InGameState.GameMode.MOUSE &&
			(mouseInput.button & MouseState.Button.LEFT)) {
			if this.gameOverCnt > 60 && !btnPressed {
				this.gameManager.startTitle(true)
			}
			this.btnPressed = true
		} else {
			this.btnPressed = false
		}
		if this.gameOverCnt == 120 {
			fadeBgm()
			disableBgm()
		}
		if this.gameOverCnt > 1200 {
			this.gameManager.startTitle(true)
		}
	}
}

func (this *InGameState) moveInGame() {
	this.field.move()
	this.ship.move()
	this.stageManager.move()
	this.enemies.move()
	this.shots.move()
	this.bullets.move()
	this.crystals.move()
	this.numIndicators.move()
	this.sparks.move()
	this.smokes.move()
	this.fragments.move()
	this.sparkFragments.move()
	this.wakes.move()
	this.screen.move()
	this.scoreReelSize += (SCORE_REEL_SIZE_DEFAULT - this.scoreReelSize) * 0.05
	this.scoreReel.move()
	if !this.isGameOver {
		this.time += 17
	}
	playMarkedSe()
}

func (this *InGameState) draw() {
	this.field.draw()
	gl.Begin(gl.TRIANGLES)
	this.wakes.draw()
	this.sparks.draw()
	gl.End()
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Begin(gl.QUADS)
	this.smokes.draw()
	gl.End()
	this.fragments.draw()
	this.sparkFragments.draw()
	this.crystals.draw()
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	this.enemies.draw()
	this.shots.draw()
	this.ship.draw()
	this.bullets.draw()
}

func (this *InGameState) drawFront() {
	this.ship.drawFront()
	this.scoreReel.draw(11.5+(SCORE_REEL_SIZE_DEFAULT-this.scoreReelSize)*3,
		-8.2-(SCORE_REEL_SIZE_DEFAULT-this.scoreReelSize)*3,
		this.scoreReelSize)
	x := -12.0
	for i := 0; i < this.left; i++ {
		gl.PushMatrix()
		gl.Translatef(x, -9, 0)
		gl.Scalef(0.7, 0.7, 0.7)
		this.ship.drawShape()
		gl.PopMatrix()
		x += 0.7
	}
	this.numIndicators.draw()
}

func (this *InGameState) drawGameParams() {
	this.stageManager.draw()
}

func (this *InGameState) drawOrtho() {
	this.drawGameParams()
	if this.isGameOver {
		letter.drawString("GAME OVER", 190, 180, 15)
	}
	if this.pauseCnt > 0 && math.Mod(this.pauseCnt, 64) < 32 {
		letter.drawString("PAUSE", 265, 210, 12)
	}
}

func (this *InGameState) drawLuminous() {
	gl.Begin(gl.TRIANGLES)
	this.sparks.drawLuminous()
	gl.End()
	this.sparkFragments.drawLuminous()
	gl.Begin(gl.QUADS)
	this.smokes.drawLuminous()
	gl.End()
}

func (this *InGameState) shipDestroyed() {
	this.clearBullets()
	this.stageManager.shipDestroyed()
	this.gameManager.initInterval()
	this.left--
	if this.left < 0 {
		this.isGameOver = true
		this.btnPressed = true
		fadeBgm()
		this.scoreReel.accelerate()
		if !this.ship.replayMode {
			disableSe()
			this.prefManager.prefData.recordResult(this.scoreReel.actualScore, this._gameMode)
			this._replayData.score = this.scoreReel.actualScore
		}
	}
}

func (this *InGameState) clearBullets() {
	this.bullets.clear()
}

func (this *InGameState) shrinkScoreReel() {
	this.scoreReelSize += (SCORE_REEL_SIZE_SMALL - this.scoreReelSize) * 0.08
}

func (this *InGameState) saveReplay(fileName string) {
	this._replayData.save(fileName)
}

func (this *InGameState) loadReplay(fileName string) {
	this._replayData = NewReplayData()
	this._replayData.load(fileName)
}

func (this *InGameState) resetReplay() {
	this._replayData = null
}

func (this *InGameState) gameMode() int {
	return this._gameMode
}

func (this *InGameState) setGameMode(v int) int {
	this._gameMode = v
	return v
}

var _ GameState = &TitleState{}

type TitleState struct {
	GameStateBase

	titleManager TitleManager
	inGameState  InGameState
	gameOverCnt  int
}

func NewTitleState(GameManager gameManager, Screen screen,
	TwinStick twinStick,
	Field field, Ship ship, ShotPool shots, BulletPool bullets, EnemyPool enemies,
	SparkPool sparks, SmokePool smokes,
	FragmentPool fragments, SparkFragmentPool sparkFragments, WakePool wakes,
	CrystalPool crystals, NumIndicatorPool numIndicators,
	StageManager stageManager, ScoreReel scoreReel,
	TitleManager titleManager, InGameState inGameState) *TitleState {

	this := &TitleState{NewGameStateBase(gameManager, screen, twinStick,
		field, ship, shots, bullets, enemies,
		sparks, smokes, fragments, sparkFragments, wakes, crystals, numIndicators,
		stageManager, scoreReel)}
	this.titleManager = titleManager
	this.inGameState = inGameState
	this.gameOverCnt = 0
	return this
}

func (this *TitleState) close() {
	this.titleManager.close()
}

func (this *TitleState) start() {
	haltBgm()
	disableBgm()
	disableSe()
	this.titleManager.start()
	if this.replayData != nil {
		this.startReplay()
	} else {
		this.titleManager.replayData = nil
	}
}

func (this *TitleState) startReplay() {
	this.ship.setReplayMode(this._replayData.shipTurnSpeed, this._replayData.shipReverseFire)
	switch this._replayData.gameMode {
	case TWIN_STICK, DOUBLE_PLAY:
		rts := twinStick
		rts.startReplay(this._replayData.twinStickInputRecord)
	}
	this.titleManager.replayData = this._replayData
	this.inGameState.gameMode = this._replayData.gameMode
	this.inGameState.startInGame()
}

func (this *TitleState) move() {
	if this._replayData != nil {
		if this.inGameState.isGameOver {
			this.gameOverCnt++
			if this.gameOverCnt > 120 {
				this.startReplay()
			}
		}
		this.inGameState.moveInGame()
	}
	this.titleManager.move()
}

func (this *TitleState) draw() {
	if this._replayData != nil {
		this.inGameState.draw()
	} else {
		this.field.draw()
	}
}

func (this *TitleState) drawFront() {
	if this._replayData != nil {
		this.inGameState.drawFront()
	}
}

func (this *TitleState) drawOrtho() {
	if this._replayData != nil {
		this.inGameState.drawGameParams()
	}
	this.titleManager.draw()
}

func (this *TitleState) drawLuminous() {
	this.inGameState.drawLuminous()
}
