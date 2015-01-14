/*
 * $Id: gamemanager.d,v 1.5 2005/09/11 00:47:40 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Manage the game state
 */
public class GameManager: abagames.util.sdl.gamemanager.GameManager {
 public:
  static float shipTurnSpeed = 1;
  static bool shipReverseFire = false;
 private:
  Pad pad;
  TwinStick twinStick;
  Mouse mouse;
  MouseAndPad mouseAndPad;
  Screen screen;
  Field field;
  Ship ship;
  StageManager stageManager;
  TitleManager titleManager;
  ScoreReel scoreReel;
  GameState state;
  TitleState titleState;
  InGameState inGameState;
  bool escPressed;

  public override void init() {
    Letter.init();
    Shot.init();
    BulletShape.init();
    EnemyShape.init();
    Turret.init();
    TurretShape.init();
    Fragment.init();
    SparkFragment.init();
    Crystal.init();
    screen = cast(Screen) abstScreen;
    pad = cast(Pad) (cast(MultipleInputDevice) input).inputs[0];
    twinStick = cast(TwinStick) (cast(MultipleInputDevice) input).inputs[1];
    twinStick.openJoystick(pad.openJoystick());
    mouse = cast(Mouse) (cast(MultipleInputDevice) input).inputs[2];
    mouse.init(screen);
    mouseAndPad = new MouseAndPad(mouse, pad);
    field = new Field;
    ship = new Ship(pad, twinStick, mouse, mouseAndPad,
                    field, screen, sparks, smokes, fragments, wakes);
    scoreReel = new ScoreReel;
    stageManager = new StageManager(field, enemies, ship, bullets,
                                    sparks, smokes, fragments, wakes);
    ship.setStageManager(stageManager);
    field.setStageManager(stageManager);
    field.setShip(ship);
    SoundManager.loadSounds();
    titleManager = new TitleManager(pad, mouse, field, this);
    inGameState = new InGameState(this, screen, pad, twinStick, mouse, mouseAndPad,
                                  field, ship, shots, bullets, enemies,
                                  sparks, smokes, fragments, sparkFragments, wakes,
                                  crystals, numIndicators, stageManager, scoreReel);
    titleState = new TitleState(this, screen, pad, twinStick, mouse, mouseAndPad,
                                field, ship, shots, bullets, enemies,
                                sparks, smokes, fragments, sparkFragments, wakes,
                                crystals, numIndicators, stageManager, scoreReel,
                                titleManager, inGameState);
    ship.setGameState(inGameState);
  }

  public override void close() {
    ship.close();
    BulletShape.close();
    EnemyShape.close();
    TurretShape.close();
    Fragment.close();
    SparkFragment.close();
    Crystal.close();
    titleState.close();
    Letter.close();
  }

  public override void start() {
    startTitle();
  }

  public void startTitle(bool fromGameover = false) {
    state = titleState;
    startState();
  }

  public void startInGame(int gameMode) {
    state = inGameState;
    inGameState.gameMode = gameMode;
    startState();
  }

  private void startState() {
    state.start();
  }

  public void initInterval() {
    mainLoop.initInterval();
  }

  public void addSlowdownRatio(float sr) {
    mainLoop.addSlowdownRatio(sr);
  }

  public override void move() {
    if (pad.keys[SDLK_ESCAPE] == SDL_PRESSED) {
      if (!escPressed) {
        escPressed = true;
        if (state == inGameState) {
          startTitle();
        } else {
          mainLoop.breakLoop();
        }
        return;
      }
    } else {
      escPressed = false;
    }
    state.move();
  }

  public override void draw() {
    SDL_Event e = mainLoop.event;
    if (e.type == SDL_VIDEORESIZE) {
      SDL_ResizeEvent re = e.resize;
      if (re.w > 150 && re.h > 100)
        screen.resized(re.w, re.h);
   }
   if (screen.startRenderToLuminousScreen()) {
      glPushMatrix();
      screen.setEyepos();
      state.drawLuminous();
      glPopMatrix();
      screen.endRenderToLuminousScreen();
    }
    screen.clear();
    glPushMatrix();
    screen.setEyepos();
    state.draw();
    glPopMatrix();
    screen.drawLuminous();
    glPushMatrix();
    screen.setEyepos();
    field.drawSideWalls();
    state.drawFront();
    glPopMatrix();
    screen.viewOrthoFixed();
    state.drawOrtho();
    screen.viewPerspective();
  }
}

/**
 * Manage the game state.
 * (e.g. title, in game, gameover, pause, ...)
 */
public class GameState {
 protected:
  GameManager gameManager;
  Screen screen;
  Pad pad;
  TwinStick twinStick;
  Mouse mouse;
  MouseAndPad mouseAndPad;
  Field field;
  Ship ship;
  StageManager stageManager;
  ScoreReel scoreReel;

  public this(GameManager gameManager, Screen screen,
              Pad pad, TwinStick twinStick, Mouse mouse, MouseAndPad mouseAndPad,
              Field field, Ship ship, StageManager stageManager, ScoreReel scoreReel) {
    this.gameManager = gameManager;
    this.screen = screen;
    this.pad = pad;
    this.twinStick = twinStick;
    this.mouse = mouse;
    this.mouseAndPad = mouseAndPad;
    this.field = field;
    this.ship = ship;
    this.stageManager = stageManager;
    this.scoreReel = scoreReel;
  }
}

public class InGameState: GameState {
 public:
  static enum GameMode {
    NORMAL, TWIN_STICK, DOUBLE_PLAY, MOUSE,
  };
  static int GAME_MODE_NUM = 4;
  static char[][] gameModeText = ["NORMAL", "TWIN STICK", "DOUBLE PLAY", "MOUSE"];
  bool isGameOver;
 private:
  static const float SCORE_REEL_SIZE_DEFAULT = 0.5f;
  static const float SCORE_REEL_SIZE_SMALL = 0.01f;
  int left;
  int time;
  int gameOverCnt;
  bool btnPressed;
  int pauseCnt;
  bool pausePressed;
  float scoreReelSize;
  int _gameMode;

  invariant {
    assert(left >= -1 && left < 10);
    assert(gameOverCnt >= 0);
    assert(pauseCnt >= 0);
    assert(scoreReelSize >= SCORE_REEL_SIZE_SMALL && scoreReelSize <= SCORE_REEL_SIZE_DEFAULT);
  }

  public this(GameManager gameManager, Screen screen,
              Pad pad, TwinStick twinStick, Mouse mouse, MouseAndPad mouseAndPad,
              Field field, Ship ship, StageManager stageManager, ScoreReel scoreReel) {
    super(gameManager, screen, pad, twinStick, mouse, mouseAndPad,
          field, ship, stageManager, scoreReel);
    left = 0;
    gameOverCnt = pauseCnt = 0;
    scoreReelSize = SCORE_REEL_SIZE_DEFAULT;
  }

  public override void start() {
    switch (_gameMode) {
    case GameMode.NORMAL:
      break;
    case GameMode.TWIN_STICK:
    case GameMode.DOUBLE_PLAY:
      break;
    case GameMode.MOUSE:
      break;
    }
    SoundManager.enableBgm();
    SoundManager.enableSe();
    startInGame();
  }

  public void startInGame() {
    clearAll();
    stageManager.start(1);
    field.start();
    ship.start(_gameMode);
    initGameState();
    screen.setScreenShake(0, 0);
    gameOverCnt = 0;
    pauseCnt = 0;
    scoreReelSize = SCORE_REEL_SIZE_DEFAULT;
    isGameOver = false;
    SoundManager.playBgm();
  }

  private void initGameState() {
    time = 0;
    left = 2;
    scoreReel.clear(9);
    NumIndicator.initTargetY();
  }

  public override void move() {
    if (pad.keys[SDLK_p] == SDL_PRESSED) {
      if (!pausePressed) {
        if (pauseCnt <= 0 && !isGameOver)
          pauseCnt = 1;
        else
          pauseCnt = 0;
      }
      pausePressed = true;
    } else {
      pausePressed = false;
    }
    if (pauseCnt > 0) {
      pauseCnt++;
      return;
    }
    moveInGame();
    if (isGameOver) {
      gameOverCnt++;
      PadState input =  pad.getState(false);
      MouseState mouseInput = mouse.getState(false);
      if ((input.button & PadState.Button.A) ||
          (gameMode == InGameState.GameMode.MOUSE &&
           (mouseInput.button & MouseState.Button.LEFT))) {
        if (gameOverCnt > 60 && !btnPressed)
          gameManager.startTitle(true);
        btnPressed = true;
      } else {
        btnPressed = false;
      }
      if (gameOverCnt == 120) {
        SoundManager.fadeBgm();
        SoundManager.disableBgm();
      }
      if (gameOverCnt > 1200)
        gameManager.startTitle(true);
    }
  }

  public void moveInGame() {
    field.move();
    ship.move();
    stageManager.move();
    enemiesMove();
    shotsMove();
    bulletsMove();
    crystalsMove();
    numIndicatorsMove();
    sparksMove();
    smokesMove();
    fragmentsMove();
    sparkFragmentsMove();
    wakesMove();
    screen.move();
    scoreReelSize += (SCORE_REEL_SIZE_DEFAULT - scoreReelSize) * 0.05f;
    scoreReel.move();
    if (!isGameOver)
      time += 17;
    SoundManager.playMarkedSe();
  }

  public override void draw() {
    field.draw();
    glBegin(GL_TRIANGLES);
    wakesDraw();
    sparksDraw();
    glEnd();
    glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);
    glBegin(GL_QUADS);
    smokes.draw();
    glEnd();
    fragmentsDraw();
    sparkFragmentsDraw();
    crystalsDraw();
    glBlendFunc(GL_SRC_ALPHA, GL_ONE);
    enemiesDraw();
    shotsDraw();
    shipDraw();
    bulletsDraw();
  }

  public void drawFront() {
    ship.drawFront();
    scoreReel.draw(11.5f + (SCORE_REEL_SIZE_DEFAULT - scoreReelSize) * 3,
                   -8.2f - (SCORE_REEL_SIZE_DEFAULT - scoreReelSize) * 3,
                   scoreReelSize);
    float x = -12;
    for (int i = 0; i < left; i++) {
      glPushMatrix();
      glTranslatef(x, -9, 0);
      glScalef(0.7f, 0.7f, 0.7f);
      ship.drawShape();
      glPopMatrix();
      x += 0.7f;
    }
    numIndicators.draw();
  }

  public void drawGameParams() {
    stageManager.draw();
  }

  public void drawOrtho() {
    drawGameParams();
    if (isGameOver)
      Letter.drawString("GAME OVER", 190, 180, 15);
    if (pauseCnt > 0 && (pauseCnt % 64) < 32)
      Letter.drawString("PAUSE", 265, 210, 12);
  }

  public override void drawLuminous() {
    glBegin(GL_TRIANGLES);
    sparks.drawLuminous();
    glEnd();
    sparkFragments.drawLuminous();
    glBegin(GL_QUADS);
    smokes.drawLuminous();
    glEnd();
  }

  public void shipDestroyed() {
    clearBullets();
    stageManager.shipDestroyed();
    gameManager.initInterval();
    left--;
    if (left < 0) {
      isGameOver = true;
      btnPressed = true;
      SoundManager.fadeBgm();
      scoreReel.accelerate();
    }
  }

  public void clearBullets() {
    bullets.clear();
  }

  public void shrinkScoreReel() {
    scoreReelSize += (SCORE_REEL_SIZE_SMALL - scoreReelSize) * 0.08f;
  }

  public int gameMode() {
    return _gameMode;
  }

  public int gameMode(int v) {
    return _gameMode = v;
  }
}

public class TitleState: GameState {
 private:
  TitleManager titleManager;
  InGameState inGameState;
  int gameOverCnt;

  invariant {
    assert(gameOverCnt >= 0);
  }

  public this(GameManager gameManager, Screen screen,
              Pad pad, TwinStick twinStick, Mouse mouse, MouseAndPad mouseAndPad,
              Field field, Ship ship,  StageManager stageManager, ScoreReel scoreReel,
              TitleManager titleManager, InGameState inGameState) {
    super(gameManager, screen, pad, twinStick, mouse, mouseAndPad,
          field, ship,  stageManager, scoreReel);
    this.titleManager = titleManager;
    this.inGameState = inGameState;
    gameOverCnt = 0;
  }

  public void close() {
    titleManager.close();
  }

  public override void start() {
    SoundManager.haltBgm();
    SoundManager.disableBgm();
    SoundManager.disableSe();
    titleManager.start();
  }

  public override void move() {
    titleManager.move();
  }

  public override void draw() {
    field.draw();
  }

  public void drawFront() {
  }

  public override void drawOrtho() {
    titleManager.draw();
  }

  public override void drawLuminous() {
    inGameState.drawLuminous();
  }
}
