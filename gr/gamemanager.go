package main

import (
	r "math/rand"
	"time"

	"github.com/dragonfax/gunroar/gr/letter"
	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
	sdl2 "github.com/veandco/go-sdl2/sdl"
)

/**
 * Manage the game state and actor pools.
 */
var _ sdl.GameManager = &GameManager{}

var shipTurnSpeed = 1.0
var shipReverseFire = false

type GameManager struct {
	sdl.GameManagerBase

	pad            *sdl.RecordablePad
	twinStick      *sdl.RecordableTwinStick
	prefManager    *PrefManager
	screen         *Screen
	field          *Field
	ship           *Ship
	shots          *ShotPool
	bullets        *BulletPool
	enemies        *EnemyPool
	sparks         *SparkPool
	smokes         *SmokePool
	fragments      *FragmentPool
	sparkFragments *SparkFragmentPool
	wakes          *WakePool
	crystals       *CrystalPool
	numIndicators  *NumIndicatorPool
	stageManager   *StageManager
	titleManager   *TitleManager
	scoreReel      *ScoreReel
	state          GameState
	titleState     *TitleState
	inGameState    *InGameState
	escPressed     bool
}

func NewGameManager() *GameManager {
	this := &GameManager{GameManagerBase: sdl.NewGameManagerBaseInternal()}
	return this
}

func (this *GameManager) Init() {
	letter.LetterInit()
	shotInit()
	fragmentInit()
	sparkFragmentInit()
	crystalInit()
	this.prefManager = this.GetPrefManager().(*PrefManager)
	this.screen = this.GetScreen().(*Screen)
	this.pad = input.Inputs[0].(*sdl.RecordablePad)
	this.twinStick = input.Inputs[1].(*sdl.RecordableTwinStick)
	this.twinStick.OpenJoystick(pad.OpenJoystick(nil))
	this.field = NewField()
	pargs := make([]interface{}, 0)
	this.sparks = NewSparkPool(120, pargs)
	pargs = append(pargs, this.field)
	this.wakes = NewWakePool(100, pargs)
	pargs = append(pargs, this.wakes)
	this.smokes = NewSmokePool(200, pargs)
	fargs := []interface{}{this.field, this.smokes}
	this.fragments = NewFragmentPool(60, fargs)
	this.sparkFragments = NewSparkFragmentPool(40, fargs)
	this.ship = NewShip(twinStick,
		this.field, screen, this.sparks, this.smokes, this.fragments, this.wakes)
	cargs := []interface{}{this.ship}
	crystals := NewCrystalPool(80, cargs)
	this.scoreReel = NewScoreReel()
	nargs := []interface{}{this.scoreReel}
	numIndicators := NewNumIndicatorPool(50, nargs)
	bargs := []interface{}{this, this.field, this.ship, this.smokes, this.wakes, crystals}
	this.bullets = NewBulletPool(240, bargs)
	eargs := []interface{}{this.field, screen, this.bullets, this.ship, this.sparks, this.smokes, this.fragments, this.sparkFragments, numIndicators, this.scoreReel}
	this.enemies = NewEnemyPool(40, eargs)
	sargs := []interface{}{this.field, this.enemies, this.sparks, this.smokes, this.bullets}
	this.shots = NewShotPool(50, sargs)
	this.ship.setShots(this.shots)
	this.ship.setEnemies(this.enemies)
	this.stageManager = NewStageManager(this.field, this.enemies, this.ship, this.bullets, this.sparks, this.smokes, this.fragments, this.wakes)
	this.ship.setStageManager(this.stageManager)
	this.field.setStageManager(this.stageManager)
	this.field.setShip(this.ship)
	this.enemies.setStageManager(this.stageManager)
	loadSounds()
	this.titleManager = NewTitleManager(prefManager, pad, this.field, this)
	this.inGameState = NewInGameState(this, this.screen, this.twinStick,
		this.field, this.ship, this.shots, this.bullets, this.enemies,
		this.sparks, this.smokes, this.fragments, this.sparkFragments, this.wakes,
		this.crystals, this.numIndicators, this.stageManager, this.scoreReel,
		this.prefManager)
	this.titleState = NewTitleState(this, screen, twinStick,
		this.field, this.ship, this.shots, this.bullets, this.enemies,
		this.sparks, this.smokes, this.fragments, this.sparkFragments, this.wakes,
		this.crystals, this.numIndicators, this.stageManager, this.scoreReel,
		this.titleManager, this.inGameState)
	this.ship.setGameState(this.inGameState)
}

func (this *GameManager) Start() {
	this.loadLastReplay()
	this.startTitle(false)
}

func (this *GameManager) startTitle(fromGameover bool /* = false */) {
	if fromGameover {
		this.saveLastReplay()
	}
	this.titleState.replayData = this.inGameState.replayData
	this.state = this.titleState
	this.startState()
}

func (this *GameManager) startInGame(gameMode GameMode) {
	this.state = this.inGameState
	this.inGameState.setGameMode(gameMode)
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
	this.inGameState.saveReplay("last.rpl")
}

func (this *GameManager) loadLastReplay() {
	this.inGameState.loadReplay("last.rpl")
}

func (this *GameManager) loadErrorReplay() {
	this.inGameState.loadReplay("error.rpl")
}

func (this *GameManager) initInterval() {
	mainLoop.InitInterval()
}

func (this *GameManager) addSlowdownRatio(sr float64) {
	mainLoop.AddSlowdownRatio(sr)
}

func (this *GameManager) Move() {
	if pad.Keys[sdl2.K_ESCAPE] == sdl2.PRESSED {
		if !this.escPressed {
			this.escPressed = true
			if this.state == this.inGameState {
				this.startTitle(false)
			} else {
				mainLoop.BreakLoop()
			}
			return
		}
	} else {
		this.escPressed = false
	}
	this.state.move()
}

func (this *GameManager) Draw() {
	e := mainLoop.Event
	if e.GetType() == sdl2.WINDOWEVENT_RESIZED {
		we := e.(*sdl2.WindowEvent)
		rew := we.Data1
		reh := we.Data2
		if rew > 150 && reh > 100 {
			this.screen.resized(int(rew), int(reh))
		}
	}
	if this.screen.startRenderToLuminousScreen() {
		gl.PushMatrix()
		this.screen.setEyepos()
		this.state.drawLuminous()
		gl.PopMatrix()
		this.screen.endRenderToLuminousScreen()
	}
	this.screen.clear()
	gl.PushMatrix()
	this.screen.setEyepos()
	this.state.draw()
	gl.PopMatrix()
	this.screen.drawLuminous()
	gl.PushMatrix()
	this.screen.setEyepos()
	this.field.drawSideWalls()
	this.state.drawFront()
	gl.PopMatrix()
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
	gameManager    *GameManager
	screen         *Screen
	twinStick      *sdl.RecordableTwinStick
	field          *Field
	ship           *Ship
	shots          *ShotPool
	bullets        *BulletPool
	enemies        *EnemyPool
	sparks         *SparkPool
	smokes         *SmokePool
	fragments      *FragmentPool
	sparkFragments *SparkFragmentPool
	wakes          *WakePool
	crystals       *CrystalPool
	numIndicators  *NumIndicatorPool
	stageManager   *StageManager
	scoreReel      *ScoreReel
	_replayData    *ReplayData
}

func NewGameStateBase(gameManager *GameManager, screen *Screen,
	twinStick *sdl.RecordableTwinStick,
	field *Field, ship *Ship, shots *ShotPool, bullets *BulletPool, enemies *EnemyPool,
	sparks *SparkPool, smokes *SmokePool,
	fragments *FragmentPool, sparkFragments *SparkFragmentPool, wakes *WakePool,
	crystals *CrystalPool, numIndicators *NumIndicatorPool,
	stageManager *StageManager, scoreReel *ScoreReel) GameStateBase {
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
	this.shots.Clear()
	this.bullets.Clear()
	this.enemies.Clear()
	this.sparks.Clear()
	this.smokes.Clear()
	this.fragments.Clear()
	this.sparkFragments.Clear()
	this.wakes.Clear()
	this.crystals.Clear()
	this.numIndicators.Clear()
}

func (this *GameStateBase) setReplayData(v ReplayData) ReplayData {
	this._replayData = v
	return v
}

func (this *GameStateBase) replayData() *ReplayData {
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
	prefManager   *PrefManager
	left          int
	time          int
	gameOverCnt   int
	btnPressed    bool
	pauseCnt      int
	pausePressed  bool
	scoreReelSize float64
	_gameMode     GameMode
}

func NewInGameState(gameManager *GameManager, screen *Screen,
	twinStick *sdl.RecordableTwinStick,
	field *Field, ship *Ship, shots *ShotPool, bullets *BulletPool, enemies *EnemyPool,
	sparks *SparkPool, smokes *SmokePool,
	fragments *FragmentPool, sparkFragments *SparkFragmentPool, wakes *WakePool,
	crystals *CrystalPool, numIndicators *NumIndicatorPool,
	stageManager *StageManager, scoreReel *ScoreReel,
	prefManager *PrefManager) *InGameState {
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
	this.prefManager.prefData().recordGameMode(this._gameMode)
	switch this._gameMode {
	case TWIN_STICK, DOUBLE_PLAY:
		rts := twinStick
		rts.StartRecord()
		this._replayData.twinStickInputRecord = rts.inputRecord
	}
	this._replayData.seed = this.rand.Int63()
	this._replayData.shipTurnSpeed = shipTurnSpeed
	this._replayData.shipReverseFire = shipReverseFire
	this._replayData.gameMode = this._gameMode
	enableBgm()
	enableSe()
	this.startInGame()
}

func (this *InGameState) startInGame() {
	this.clearAll()
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
	SetBaseShapeRandSeed(seed)
	this.ship.setRandSeed(seed)
	setShotRandSeed(seed)
	this.stageManager.setRandSeed(seed)
	setNumReelRandSeed(seed)
	setNumIndicatorRandSeed(seed)
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
	if pad.Keys[sdl2.K_p] == sdl2.PRESSED {
		if !this.pausePressed {
			if this.pauseCnt <= 0 && !this.isGameOver {
				this.pauseCnt = 1
			} else {
				this.pauseCnt = 0
			}
		}
		this.pausePressed = true
	} else {
		this.pausePressed = false
	}
	if this.pauseCnt > 0 {
		this.pauseCnt++
		return
	}
	this.moveInGame()
	if this.isGameOver {
		this.gameOverCnt++
		if this.input.button & sdl.ButtonA /* || (this.gameMode == InGameState.GameMode.MOUSE && (mouseInput.button & MouseState.Button.LEFT)) */ {
			if this.gameOverCnt > 60 && !this.btnPressed {
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
	this.enemies.Move()
	this.shots.Move()
	this.bullets.Move()
	this.crystals.Move()
	this.numIndicators.Move()
	this.sparks.Move()
	this.smokes.Move()
	this.fragments.Move()
	this.sparkFragments.Move()
	this.wakes.Move()
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
	this.wakes.Draw()
	this.sparks.Draw()
	gl.End()
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Begin(gl.QUADS)
	this.smokes.Draw()
	gl.End()
	this.fragments.Draw()
	this.sparkFragments.Draw()
	this.crystals.Draw()
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	this.enemies.Draw()
	this.shots.Draw()
	this.ship.draw()
	this.bullets.Draw()
}

func (this *InGameState) drawFront() {
	this.ship.drawFront()
	this.scoreReel.draw(11.5+(SCORE_REEL_SIZE_DEFAULT-this.scoreReelSize)*3,
		-8.2-(SCORE_REEL_SIZE_DEFAULT-this.scoreReelSize)*3,
		this.scoreReelSize)
	x := -12.0
	for i := 0; i < this.left; i++ {
		gl.PushMatrix()
		gl.Translated(x, -9, 0)
		gl.Scalef(0.7, 0.7, 0.7)
		this.ship.drawShape()
		gl.PopMatrix()
		x += 0.7
	}
	this.numIndicators.Draw()
}

func (this *InGameState) drawGameParams() {
	this.stageManager.draw()
}

func (this *InGameState) drawOrtho() {
	this.drawGameParams()
	if this.isGameOver {
		letter.DrawString("GAME OVER", 190, 180, 15, letter.TO_RIGHT, 0, false, 0)
	}
	if this.pauseCnt > 0 && this.pauseCnt%64 < 32 {
		letter.DrawString("PAUSE", 265, 210, 12, letter.TO_RIGHT, 0, false, 0)
	}
}

func (this *InGameState) drawLuminous() {
	gl.Begin(gl.TRIANGLES)
	this.sparks.DrawLuminous()
	gl.End()
	this.sparkFragments.DrawLuminous()
	gl.Begin(gl.QUADS)
	this.smokes.DrawLuminous()
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
			this._replayData.score = this.scoreReel.actualScore()
		}
	}
}

func (this *InGameState) clearBullets() {
	this.bullets.Clear()
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
	this._replayData = nil
}

func (this *InGameState) gameMode() GameMode {
	return this._gameMode
}

func (this *InGameState) setGameMode(v GameMode) GameMode {
	this._gameMode = v
	return v
}

var _ GameState = &TitleState{}

type TitleState struct {
	GameStateBase

	titleManager *TitleManager
	inGameState  *InGameState
	gameOverCnt  int
}

func NewTitleState(gameManager *GameManager, screen *Screen,
	twinStick *sdl.RecordableTwinStick,
	field *Field, ship *Ship, shots *ShotPool, bullets *BulletPool, enemies *EnemyPool,
	sparks *SparkPool, smokes *SmokePool,
	fragments *FragmentPool, sparkFragments *SparkFragmentPool, wakes *WakePool,
	crystals *CrystalPool, numIndicators *NumIndicatorPool,
	stageManager *StageManager, scoreReel *ScoreReel,
	titleManager *TitleManager, inGameState *InGameState) *TitleState {

	this := &TitleState{GameStateBase: NewGameStateBase(gameManager, screen, twinStick,
		field, ship, shots, bullets, enemies,
		sparks, smokes, fragments, sparkFragments, wakes, crystals, numIndicators,
		stageManager, scoreReel)}
	this.titleManager = titleManager
	this.inGameState = inGameState
	this.gameOverCnt = 0
	return this
}

func (this *TitleState) start() {
	haltBgm()
	disableBgm()
	disableSe()
	this.titleManager.start()
	if this.replayData != nil {
		this.startReplay()
	} else {
		this.titleManager.replayData(nil)
	}
}

func (this *TitleState) startReplay() {
	this.ship.setReplayMode(this._replayData.shipTurnSpeed, this._replayData.shipReverseFire)
	switch this._replayData.gameMode {
	case TWIN_STICK, DOUBLE_PLAY:
		rts := twinStick
		rts.startReplay(this._replayData.twinStickInputRecord)
	}
	this.titleManager.replayData(this._replayData)
	this.inGameState.setGameMode(this._replayData.gameMode)
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
