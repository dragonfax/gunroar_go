/*
 * $Id: gamemanager.d,v 1.5 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Manage the game state
 */

var shipTurnSpeed float32 = 1
var shipReverseFire bool = false

type GameManager struct {
  pad Pad
  // twinStick TwinStick
  mouse Mouse
  mouseAndPad MouseAndPad
  screen Screen
  field Field
  ship Ship
  stageManager StageManager
  titleManager TitleManager
  scoreReel ScoreReel
  state GameState
  titleState TitleState
  inGameState InGameState
  escPressed bool
}

  func (this *GameManager) init() {
    InitLetter.init()
    InitShot.init()
    InitBulletShape.init()
    InitEnemyShape.init()
    InitTurret.init()
    InitTurretShape.init()
    InitFragment.init()
    InitSparkFragment.init()
    InitCrystal.init()
    this.pad = input.inputs[0]
    // twinStick = cast(TwinStick) (cast(MultipleInputDevice) input).inputs[1]
    // twinStick.openJoystick(pad.openJoystick())
    this.mouse = input.inputs[2]
    this.mouse.init(screen)
    this.mouseAndPad = NewMouseAndPad(mouse, pad)
    this.field = NewField()
    this.ship = NewShip(pad, twinStick, mouse, mouseAndPad, field, screen )
    this.scoreReel = NewScoreReel()
    this.stageManager = NewStageManager(field, enemies, ship)
    this.ship.setStageManager(stageManager)
    this.field.setStageManager(stageManager)
    this.field.setShip(ship)
    loadSounds()
    this.titleManager = NewTitleManager(pad, mouse, field, this)
    this.inGameState = NewInGameState(this, screen, pad, /*twinStick, */mouse, mouseAndPad,
                                  field, ship, stageManager, scoreReel)
    this.titleState = NewTitleState(this, screen, pad, /*twinStick, */mouse, mouseAndPad,
                                field, ship, stageManager, scoreReel,
                                titleManager, inGameState)
    this.ship.setGameState(this.inGameState)
  }

  func (this *GameManager) close() {
    this.ship.close()
    CloseBulletShape()
    CloseEnemyShape()
    CloseTurretShape()
    CloseFragment()
    CloseSparkFragment()
    CloseCrystal()
    this.titleState()
    CloseLetter()
  }

  func (this *GameManager)  start() {
    this.startTitle()
  }

  func (this *GameManager)  startTitle(fromGameover bool /*= false*/) {
    this.state = this.titleState
    this.startState()
  }

  func (this *GameManager) startInGame(gameMode GameMode) {
    this.state = this.inGameState
    this.inGameState.gameMode = gameMode
    this.startState()
  }

  func (this *GameManager)  startState() {
    this.state.start()
  }

  func (this *GameManager) initInterval() {
    mainLoop.initInterval()
  }

  func (this *GameManager)  addSlowdownRatio(sr float32) {
    mainLoop.addSlowdownRatio(sr)
  }

  func (this *GameManager)  move() {
    if (this.pad.keys[SDL.K_ESCAPE] == sdl.PRESSED) {
      if (!escPressed) {
        this.escPressed = true
        if (this.state == this.inGameState) {
          this.startTitle()
        } else {
          mainLoop.breakLoop()
        }
        return
      }
    } else {
      this.escPressed = false
    }
    this.state.move()
  }

  func (this *GameManager) draw() {
		e := mainLoop.event
    if (e.type == sdl.VIDEORESIZE) {
      sdl.ResizeEvent re = e.resize
      if (re.w > 150 && re.h > 100)
        screen.resized(re.w, re.h)
   }
   if (screen.startRenderToLuminousScreen()) {
      glPushMatrix()
      screen.setEyepos()
      state.drawLuminous()
      glPopMatrix()
      screen.endRenderToLuminousScreen()
    }
    screen.clear()
    glPushMatrix()
    screen.setEyepos()
    state.draw()
    glPopMatrix()
    screen.drawLuminous()
    glPushMatrix()
    screen.setEyepos()
    field.drawSideWalls()
    state.drawFront()
    glPopMatrix()
    screen.viewOrthoFixed()
    state.drawOrtho()
    screen.viewPerspective()
  }
}

/**
 * Manage the game state.
 * (e.g. title, in game, gameover, pause, ...)
 */
public class GameState {
 protected:
  GameManager gameManager
  Screen screen
  Pad pad
  TwinStick twinStick
  Mouse mouse
  MouseAndPad mouseAndPad
  Field field
  Ship ship
  StageManager stageManager
  ScoreReel scoreReel

  public this(GameManager gameManager, Screen screen,
              Pad pad, TwinStick twinStick, Mouse mouse, MouseAndPad mouseAndPad,
              Field field, Ship ship, StageManager stageManager, ScoreReel scoreReel) {
    this.gameManager = gameManager
    this.screen = screen
    this.pad = pad
    this.twinStick = twinStick
    this.mouse = mouse
    this.mouseAndPad = mouseAndPad
    this.field = field
    this.ship = ship
    this.stageManager = stageManager
    this.scoreReel = scoreReel
  }
}

type GameMode int

const (
  GameModeNORMAL GameMode = iota
  GameModeTWIN_STICK
  GameModeDOUBLE_PLAY
  GameModeMOUSE
)

const GAME_MODE_NUM := 4
var gameModeText []string = []string{"NORMAL", "TWIN STICK", "DOUBLE PLAY", "MOUSE"}
var isGameOver bool

const SCORE_REEL_SIZE_DEFAULT = 0.5
const SCORE_REEL_SIZE_SMALL = 0.01

public class InGameState: GameState {
  left, time, gameOverCnt int
  btnPressed bool
  pauseCnt int
  pausePressed bool
  scoreReelSize float32
  gameMode GameMode

  public this(GameManager gameManager, Screen screen,
              Pad pad, TwinStick twinStick, Mouse mouse, MouseAndPad mouseAndPad,
              Field field, Ship ship, StageManager stageManager, ScoreReel scoreReel) {
    super(gameManager, screen, pad, twinStick, mouse, mouseAndPad,
          field, ship, stageManager, scoreReel)
    left = 0
    gameOverCnt = pauseCnt = 0
    scoreReelSize = SCORE_REEL_SIZE_DEFAULT
  }

  public override void start() {
    switch (_gameMode) {
    case GameMode.NORMAL:
      break
    case GameMode.TWIN_STICK:
    case GameMode.DOUBLE_PLAY:
      break
    case GameMode.MOUSE:
      break
    }
    SoundManager.enableBgm()
    SoundManager.enableSe()
    startInGame()
  }

  public void startInGame() {
    clearAll()
    stageManager.start(1)
    field.start()
    ship.start(_gameMode)
    initGameState()
    screen.setScreenShake(0, 0)
    gameOverCnt = 0
    pauseCnt = 0
    scoreReelSize = SCORE_REEL_SIZE_DEFAULT
    isGameOver = false
    SoundManager.playBgm()
  }

  private void initGameState() {
    time = 0
    left = 2
    scoreReel.clear(9)
    NumIndicator.initTargetY()
  }

  public override void move() {
    if (pad.keys[SDL.K_p] == sdl.PRESSED) {
      if (!pausePressed) {
        if (pauseCnt <= 0 && !isGameOver)
          pauseCnt = 1
        else
          pauseCnt = 0
      }
      pausePressed = true
    } else {
      pausePressed = false
    }
    if (pauseCnt > 0) {
      pauseCnt++
      return
    }
    moveInGame()
    if (isGameOver) {
      gameOverCnt++
      PadState input =  pad.getState(false)
      MouseState mouseInput = mouse.getState(false)
      if ((input.button & PadState.Button.A) ||
          (gameMode == InGameState.GameMode.MOUSE &&
           (mouseInput.button & MouseState.Button.LEFT))) {
        if (gameOverCnt > 60 && !btnPressed)
          gameManager.startTitle(true)
        btnPressed = true
      } else {
        btnPressed = false
      }
      if (gameOverCnt == 120) {
        SoundManager.fadeBgm()
        SoundManager.disableBgm()
      }
      if (gameOverCnt > 1200)
        gameManager.startTitle(true)
    }
  }

  public void moveInGame() {
    field.move()
    ship.move()
    stageManager.move()
    enemiesMove()
    shotsMove()
    bulletsMove()
    crystalsMove()
    numIndicatorsMove()
    sparksMove()
    smokesMove()
    fragmentsMove()
    sparkFragmentsMove()
    wakesMove()
    screen.move()
    scoreReelSize += (SCORE_REEL_SIZE_DEFAULT - scoreReelSize) * 0.05
    scoreReel.move()
    if (!isGameOver)
      time += 17
    SoundManager.playMarkedSe()
  }

  public override void draw() {
    field.draw()
    glBegin(GL_TRIANGLES)
    wakesDraw()
    sparksDraw()
    glEnd()
    glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)
    glBegin(GL_QUADS)
    smokes.draw()
    glEnd()
    fragmentsDraw()
    sparkFragmentsDraw()
    crystalsDraw()
    glBlendFunc(GL_SRC_ALPHA, GL_ONE)
    enemiesDraw()
    shotsDraw()
    shipDraw()
    bulletsDraw()
  }

  public void drawFront() {
    ship.drawFront()
    scoreReel.draw(11.5 + (SCORE_REEL_SIZE_DEFAULT - scoreReelSize) * 3,
                   -8.2 - (SCORE_REEL_SIZE_DEFAULT - scoreReelSize) * 3,
                   scoreReelSize)
    float x = -12
    for (int i = 0; i < left; i++) {
      glPushMatrix()
      glTranslatef(x, -9, 0)
      glScalef(0.7, 0.7, 0.7)
      ship.drawShape()
      glPopMatrix()
      x += 0.7
    }
    numIndicators.draw()
  }

  public void drawGameParams() {
    stageManager.draw()
  }

  public void drawOrtho() {
    drawGameParams()
    if (isGameOver)
      Letter.drawString("GAME OVER", 190, 180, 15)
    if (pauseCnt > 0 && (pauseCnt % 64) < 32)
      Letter.drawString("PAUSE", 265, 210, 12)
  }

  public override void drawLuminous() {
    glBegin(GL_TRIANGLES)
    sparks.drawLuminous()
    glEnd()
    sparkFragments.drawLuminous()
    glBegin(GL_QUADS)
    smokes.drawLuminous()
    glEnd()
  }

  public void shipDestroyed() {
    clearBullets()
    stageManager.shipDestroyed()
    gameManager.initInterval()
    left--
    if (left < 0) {
      isGameOver = true
      btnPressed = true
      SoundManager.fadeBgm()
      scoreReel.accelerate()
    }
  }

  public void clearBullets() {
    bullets.clear()
  }

  public void shrinkScoreReel() {
    scoreReelSize += (SCORE_REEL_SIZE_SMALL - scoreReelSize) * 0.08
  }

  public int gameMode() {
    return _gameMode
  }

  public int gameMode(int v) {
    return _gameMode = v
  }
}

public class TitleState: GameState {
 private:
  TitleManager titleManager
  InGameState inGameState
  int gameOverCnt

  public this(GameManager gameManager, Screen screen,
              Pad pad, TwinStick twinStick, Mouse mouse, MouseAndPad mouseAndPad,
              Field field, Ship ship,  StageManager stageManager, ScoreReel scoreReel,
              TitleManager titleManager, InGameState inGameState) {
    super(gameManager, screen, pad, twinStick, mouse, mouseAndPad,
          field, ship,  stageManager, scoreReel)
    this.titleManager = titleManager
    this.inGameState = inGameState
    gameOverCnt = 0
  }

  public void close() {
    titleManager.close()
  }

  public override void start() {
    SoundManager.haltBgm()
    SoundManager.disableBgm()
    SoundManager.disableSe()
    titleManager.start()
  }

  public override void move() {
    titleManager.move()
  }

  public override void draw() {
    field.draw()
  }

  public void drawFront() {
  }

  public override void drawOrtho() {
    titleManager.draw()
  }

  public override void drawLuminous() {
    inGameState.drawLuminous()
  }
}
