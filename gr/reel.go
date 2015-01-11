/*
 * $Id: reel.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Rolling reel that displays the score.
 */

static const int MAX_DIGIT = 16;

type ScoreReel struct {
  score, targetScore int
  actualScore int
  digit int
  numReel [MAX_DIGIT]NumReel
}

func (sr *ScoreReel) Init() {
	for i,_ := range sr.numReel {
		sr.numReel[i].Init()
	}
	sr.digit = 1
}

func (sr *ScoreReel) clear(digit int /*= 9 */) {
	sr.score = sr.targetScore = sr.actualScore = 0
	sr.digit = digit
	for i := 0; i < digit; i++ {
		sr.numReel[i].clear()
	}
}

func (sr *ScoreReel)  move() {
	for i := 0; i < digit; i++ {
		sr.numReel[i].move()
	}
}

func (sr *ScoreReel)  draw(x float32, y float32, s float32) {
	lx = x, ly = y float32
	for i := 0; i < digit; i++ {
		sr.numReel[i].draw(lx, ly, s);
		lx -= s * 2
	}
}

func (sr *ScoreReel) addReelScore(as int) {
	sr.targetScore += as
	ts := sr.targetScore
	for i := 0; i < digit; i++ {
		sr.numReel[i].targetDeg = float32(ts * 360 / 10)
		ts /= 10
		if (ts < 0) {
			break
		}
	}
}

func (sr *ScoreReel)  accelerate() {
	for i := 0; i < digit; i++ {
		sr.numReel[i].accelerate()
	}
}

func (sr *ScoreReel)  addActualScore(as int) {
	sr.actualScore += as
}


const VEL_MIN float32 = 5

type NumReel struct {
  deg float32
  _targetDeg float32
  ofs float32
  velRatio float32
}

func (nr *NumReel) Init() {
  nr.deg = 0
	nr.ofs = 0
	nr.velRatio = 1;
}

func (nr *NumReel) clear() {
	nr.Init()
}

func (nr *NumReel)  move() {
	vd := nr.targetDeg - nr.deg
	vd *= 0.05 * nr.velRatio
	if (vd < VEL_MIN * nr.velRatio) {
		vd = VEL_MIN * nr.velRatio
	}
	nr.deg += vd
	if (nr.deg > nr.targetDeg) {
		nr.deg = nr.targetDeg
	}
}

func (nr *NumReel) draw(x float, y float, s float) {
	n := int((deg * 10 / 360 + 0.99f) + 1) % 10;
	d := deg % 360
	od := d - n * 360 / 10
	od -= 15;
	normalizeDeg360(od)
	od *= 1.5
	for i := 0; i < 3; i++ {
		gl.PushMatrix()
		if (nr.ofs > 0.005) {
			gl.Translatef(x + rand.Float32() * nr.ofs, y + rand.Float32() * nr.ofs, 0)
		} else {
			gl.Translatef(x, y, 0)
		}
		gl.Rotatef(od, 1, 0, 0)
		gl.Translatef(0, 0, s * 2.4)
		gl.Scalef(s, -s, s)
		a := float32(1 - fabs((od + 15) / (360 / 10 * 1.5)) / 2)
		if (a < 0) {
			a = 0
		}
		Screen.setColor(a, a, a)
		Letter.drawLetter(n, 2)
		Screen.setColor(a / 2, a / 2, a / 2)
		Letter.drawLetter(n, 3)
		gl.PopMatrix()
		n--
		if (n < 0) {
			n = 9
		}
		od += 360 / 10 * 1.5
		Math.normalizeDeg360(od)
	}
	ofs *= 0.95f;
}

fund (nr *NumReel) targetDeg(td float) {
	if ((td - nr.targetDeg) > 1) {
		nr.ofs += 0.1
	}
	nr.targetDeg = td
	return nr.targetDeg
}

func (nr *NumReel) accelerate() {
	nr.velRatio = 4
}

/**
 * Flying indicator that shows the score and the multiplier.
 */

static enum IndicatorType {
	SCORE, MULTIPLIER,
};
static enum FlyingToType {
	RIGHT, BOTTOM,
};
static const float TARGET_Y_MIN = -7;
static const float TARGET_Y_MAX = 7;
static const float TARGET_Y_INTERVAL = 1;
static float targetY;
struct Target {
	Vector pos;
	int flyingTo;
	float initialVelRatio;
	float size;
	int n;
	int cnt;
};

type NumIndicator struct {
	ActorImpl

  ScoreReel scoreReel;
  Vector pos, vel;
  int n, type;
  float size;
  int cnt;
  float alpha;
  Target[4] target;
  int targetIdx;
  int targetNum;

  invariant {
    assert(targetY <= TARGET_Y_MAX && targetY >= TARGET_Y_MIN);
    assert(pos.x < 15 && pos.x > -15);
    assert(pos.y < 20 && pos.y > -20);
    assert(vel.x < 10 && vel.x > -10);
    assert(vel.y < 10 && vel.y > -10);
    assert(alpha >= 0 && alpha <= 1);
    foreach (Target t; target) {
      assert(t.pos.x < 15 && t.pos.x > -15);
      assert(t.pos.y < 20 && t.pos.y > -20);
      assert(t.initialVelRatio >= 0);
      assert(t.size >= 0);
    }
    assert(targetIdx >= -1 && targetIdx <= 4);
    assert(targetNum >= 0 && targetNum <= 4);
  }

  public static this() {
    rand = new Rand;
    targetY = TARGET_Y_MIN;
  }

  public static void setRandSeed(long seed) {
    rand.setSeed(seed);
  }

  public static void initTargetY() {
    targetY = TARGET_Y_MIN;
  }

  public static float getTargetY() {
    float ty = targetY;
    targetY += TARGET_Y_INTERVAL;
    if (targetY > TARGET_Y_MAX)
      targetY = TARGET_Y_MIN;
    return ty;
  }

  public static void decTargetY() {
    targetY -= TARGET_Y_INTERVAL;
    if (targetY < TARGET_Y_MIN)
      targetY = TARGET_Y_MAX;
  }

  public this() {
    pos = new Vector;
    vel = new Vector;
    foreach (inout Target t; target) {
      t.pos = new Vector;
      t.initialVelRatio = 0;
      t.size = 0;
    }
    targetIdx = targetNum = 0;
    alpha = 1;
  }

  public void init(Object[] args) {
    scoreReel = cast(ScoreReel) args[0];
  }

  public void set(int n, IndicatorType type, float size, Vector p) {
    set(n, type, size, p.x, p.y);
  }

  public void set(int n, IndicatorType type, float size, float x, float y) {
    if (exists && this.type == IndicatorType.SCORE) {
      if (this.target[targetIdx].flyingTo == FlyingToType.RIGHT)
        decTargetY();
      scoreReel.addReelScore(target[targetNum - 1].n);
    }
    this.n = n;
    this.type = type;
    this.size = size;
    pos.x = x;
    pos.y = y;
    targetIdx = -1;
    targetNum = 0;
    alpha = 0.1f;
    exists = true;
  }

  public void addTarget(float x, float y, FlyingToType flyingTo, float initialVelRatio,
                        float size, int n, int cnt) {
    target[targetNum].pos.x = x;
    target[targetNum].pos.y = y;
    target[targetNum].flyingTo = flyingTo;
    target[targetNum].initialVelRatio = initialVelRatio;
    target[targetNum].size = size;
    target[targetNum].n = n;
    target[targetNum].cnt = cnt;
    targetNum++;
  }

  public void gotoNextTarget() {
    targetIdx++;
    if (targetIdx > 0)
      SoundManager.playSe("score_up.wav");
    if (targetIdx >= targetNum) {
      if (target[targetIdx - 1].flyingTo == FlyingToType.BOTTOM)
        scoreReel.addReelScore(target[targetIdx - 1].n);
      exists = false;
      return;
    }
    switch (target[targetIdx].flyingTo) {
    case FlyingToType.RIGHT:
      vel.x = -0.3f + rand.nextSignedFloat(0.05f);
      vel.y = rand.nextSignedFloat(0.1f);
      break;
    case FlyingToType.BOTTOM:
      vel.x = rand.nextSignedFloat(0.1f);
      vel.y = 0.3f + rand.nextSignedFloat(0.05f);
      decTargetY();
      break;
    }
    vel *= target[targetIdx].initialVelRatio;
    cnt = target[targetIdx].cnt;
  }

  public void move() {
    if (targetIdx < 0)
      return;
    Vector tp = target[targetIdx].pos;
    switch (target[targetIdx].flyingTo) {
    case FlyingToType.RIGHT:
      vel.x += (tp.x - pos.x) * 0.0036f;
      pos.y += (tp.y - pos.y) * 0.1f;
      if (fabs(pos.y - tp.y) < 0.5f)
        pos.y += (tp.y - pos.y) * 0.33f;
      alpha += (1 - alpha) * 0.03f;
      break;
    case FlyingToType.BOTTOM:
      pos.x += (tp.x - pos.x) * 0.1f;
      vel.y += (tp.y - pos.y) * 0.0036f;
      alpha *= 0.97f;
      break;
    }
    vel *= 0.98f;
    size += (target[targetIdx].size - size) * 0.025f;
    pos += vel;
    int vn = cast(int) ((target[targetIdx].n - n) * 0.2f);
    if (vn < 10 && vn > -10)
      n = target[targetIdx].n;
    else
      n += vn;
    switch (target[targetIdx].flyingTo) {
    case FlyingToType.RIGHT:
      if (pos.x > tp.x) {
        pos.x = tp.x;
        vel.x *= -0.05f;
      }
      break;
    case FlyingToType.BOTTOM:
      if (pos.y < tp.y) {
        pos.y = tp.y;
        vel.y *= -0.05f;
      }
      break;
    }
    cnt--;
    if (cnt < 0)
      gotoNextTarget();
  }

  public void draw() {
    Screen.setColor(alpha, alpha, alpha);
    switch (type) {
    case IndicatorType.SCORE:
      Letter.drawNumSign(n, pos.x, pos.y, size, Letter.LINE_COLOR);
      break;
    case IndicatorType.MULTIPLIER:
      Screen.setColor(alpha, alpha, alpha);
      Letter.drawNumSign(n, pos.x, pos.y, size, Letter.LINE_COLOR, 33, 3);
      break;
    }
  }
}

