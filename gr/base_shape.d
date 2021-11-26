/*
 * $Id: shape.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
module abagames.gr.shape;

private import std.math;
private import opengl;
private import abagames.util.vector;
private import abagames.util.rand;
private import abagames.util.sdl.shape;
private import abagames.gr.screen;
private import abagames.gr.particle;

/**
 * Shape of a ship/platform/turret/bridge.
 */
public class BaseShape: DrawableShape {
 public:
  static enum ShapeType {
    SHIP, SHIP_ROUNDTAIL, SHIP_SHADOW, PLATFORM, TURRET, BRIDGE,
    SHIP_DAMAGED, SHIP_DESTROYED,
    PLATFORM_DAMAGED, PLATFORM_DESTROYED,
    TURRET_DAMAGED, TURRET_DESTROYED,
  };
 private:
  static const int POINT_NUM = 16;
  static Rand rand;
  static Vector wakePos;
  float size, distRatio, spinyRatio;
  int type;
  float r, g, b;
  static const int PILLAR_POINT_NUM = 8;
  Vector[] pillarPos;
  Vector[] _pointPos;
  float[] _pointDeg;

  invariant {
    assert(wakePos.x < 15 && wakePos.x > -15);
    assert(wakePos.y < 60 && wakePos.y > -40);
    assert(size > 0 && size < 20);
    assert(distRatio >= 0 && distRatio <= 1);
    assert(spinyRatio >= 0 && spinyRatio <= 1);
    assert(type >= 0);
    assert(r >= 0 && r <= 1);
    assert(g >= 0 && g <= 1);
    assert(b >= 0 && b <= 1);
    foreach (Vector p; pillarPos) {
      assert(p.x < 20 && p.x > -20);
      assert(p.y < 20 && p.x > -20);
    }
    foreach (Vector p; _pointPos) {
      assert(p.x < 20 && p.x > -20);
      assert(p.y < 20 && p.x > -20);
    }
    foreach (float d; _pointDeg)
      assert(d <>= 0);
  }

  public static this() {
    rand = new Rand;
    wakePos = new Vector;
  }

  public static void setRandSeed(long seed) {
    rand.setSeed(seed);
  }

  public this(float size, float distRatio, float spinyRatio,
              int type, float r, float g, float b) {
    this.size = size;
    this.distRatio = distRatio;
    this.spinyRatio = spinyRatio;
    this.type = type;
    this.r = r;
    this.g = g;
    this.b = b;
    super();
  }

  public override void createDisplayList() {
    float height = size * 0.5f;
    float z = 0;
    float sz = 1;
    if (type == ShapeType.BRIDGE)
      z += height;
    if (type != ShapeType.SHIP_DESTROYED)
      Screen.setColor(r, g, b);
    glBegin(GL_LINE_LOOP);
    if (type != ShapeType.BRIDGE)
      createLoop(sz, z, false, true);
    else
      createSquareLoop(sz, z, false, true);
    glEnd();
    if (type != ShapeType.SHIP_SHADOW && type != ShapeType.SHIP_DESTROYED &&
        type != ShapeType.PLATFORM_DESTROYED && type != ShapeType.TURRET_DESTROYED) {
      Screen.setColor(r * 0.4f, g * 0.4f, b * 0.4f);
      glBegin(GL_TRIANGLE_FAN);
      createLoop(sz, z, true);
      glEnd();
    }
    switch (type) {
    case ShapeType.SHIP:
    case ShapeType.SHIP_ROUNDTAIL:
    case ShapeType.SHIP_SHADOW:
    case ShapeType.SHIP_DAMAGED:
    case ShapeType.SHIP_DESTROYED:
      if (type != ShapeType.SHIP_DESTROYED)
        Screen.setColor(r * 0.4f, g * 0.4f, b * 0.4f);
      for (int i = 0; i < 3; i++) {
        z -= height / 4;
        sz -= 0.2f;
        glBegin(GL_LINE_LOOP);
        createLoop(sz, z);
        glEnd();
      }
      break;
    case ShapeType.PLATFORM:
    case ShapeType.PLATFORM_DAMAGED:
    case ShapeType.PLATFORM_DESTROYED:
      Screen.setColor(r * 0.4f, g * 0.4f, b * 0.4f);
      for (int i = 0; i < 3; i++) {
        z -= height / 3;
        foreach (Vector pp; pillarPos) {
          glBegin(GL_LINE_LOOP);
          createPillar(pp, size * 0.2f, z);
          glEnd();
        }
      }
      break;
    case ShapeType.BRIDGE:
    case ShapeType.TURRET:
    case ShapeType.TURRET_DAMAGED:
      Screen.setColor(r * 0.6f, g * 0.6f, b * 0.6f);
      z += height;
      sz -= 0.33f;
      glBegin(GL_LINE_LOOP);
      if (type == ShapeType.BRIDGE)
        createSquareLoop(sz, z);
      else
        createSquareLoop(sz, z / 2, false, 3);
      glEnd();
      Screen.setColor(r * 0.25f, g * 0.25f, b * 0.25f);
      glBegin(GL_TRIANGLE_FAN);
      if (type == ShapeType.BRIDGE)
        createSquareLoop(sz, z, true);
      else
        createSquareLoop(sz, z / 2, true, 3);
      glEnd();
      break;
    case ShapeType.TURRET_DESTROYED:
      break;
    }
  }

  private void createLoop(float s, float z, bool backToFirst = false, bool record = false) {
    float d = 0;
    int pn;
    bool firstPoint = true;
    float fpx, fpy;
    for (int i = 0; i < POINT_NUM; i++) {
      if (type != ShapeType.SHIP && type != ShapeType.SHIP_DESTROYED && type != ShapeType.SHIP_DAMAGED &&
          i > POINT_NUM * 2 / 5 && i <= POINT_NUM * 3 / 5)
        continue;
      if ((type == ShapeType.TURRET || type == ShapeType.TURRET_DAMAGED || type == ShapeType.TURRET_DESTROYED) &&
          (i <= POINT_NUM / 5 || i > POINT_NUM * 4 / 5))
        continue;
      d = PI * 2 * i / POINT_NUM;
      float cx = sin(d) * size * s * (1 - distRatio);
      float cy = cos(d) * size * s;
      float sx, sy;
      if (i == POINT_NUM / 4 || i == POINT_NUM / 4 * 3)
        sy = 0;
      else
        sy = 1 / (1 + fabs(tan(d)));
      assert(sy <>= 0);
      sx = 1 - sy;
      if (i >= POINT_NUM / 2)
        sx *= -1;
      if (i >= POINT_NUM / 4 && i <= POINT_NUM / 4 * 3)
        sy *= -1;
      sx *= size * s * (1 - distRatio);
      sy *= size * s;
      float px = cx * (1 - spinyRatio) + sx * spinyRatio;
      float py = cy * (1 - spinyRatio) + sy * spinyRatio;
      glVertex3f(px, py, z);
      if (backToFirst && firstPoint) {
        fpx = px;
        fpy = py;
        firstPoint = false;
      }
      if (record) {
        if (i == POINT_NUM / 8 || i == POINT_NUM / 8 * 3 ||
            i == POINT_NUM / 8 * 5 || i == POINT_NUM / 8 * 7)
          pillarPos ~= new Vector(px * 0.8f, py * 0.8f);
        _pointPos ~= new Vector(px, py);
        _pointDeg ~= d;
      }
    }
    if (backToFirst)
      glVertex3f(fpx, fpy, z);
  }

  private void createSquareLoop(float s, float z, bool backToFirst = false, float yRatio = 1) {
    float d;
    int pn;
    if (backToFirst)
      pn = 4;
    else
      pn = 3;
    for (int i = 0; i <= pn; i++) {
      d = PI * 2 * i / 4 + PI / 4;
      float px = sin(d) * size * s;
      float py = cos(d) * size * s;
      if (py > 0)
        py *= yRatio;
      glVertex3f(px, py, z);
    }
  }

  private void createPillar(Vector p, float s, float z) {
    float d;
    for (int i = 0; i < PILLAR_POINT_NUM; i++) {
      d = PI * 2 * i / PILLAR_POINT_NUM;
      glVertex3f(sin(d) * s + p.x, cos(d) * s + p.y, z);
    }
  }

  public void addWake(WakePool wakes, Vector pos, float deg, float spd, float sr = 1) {
    float sp = spd;
    if (sp > 0.1f)
      sp = 0.1f;
    float sz = size;
    if (sz > 10)
      sz = 10;
    wakePos.x = pos.x + sin(deg + PI / 2 + 0.7f) * size * 0.5f * sr;
    wakePos.y = pos.y + cos(deg + PI / 2 + 0.7f) * size * 0.5f * sr;
    Wake w = wakes.getInstanceForced();
    w.set(wakePos, deg + PI - 0.2f + rand.nextSignedFloat(0.1f), sp, 40, sz * 32 * sr);
    wakePos.x = pos.x + sin(deg - PI / 2 - 0.7f) * size * 0.5f * sr;
    wakePos.y = pos.y + cos(deg - PI / 2 - 0.7f) * size * 0.5f * sr;
    w = wakes.getInstanceForced();
    w.set(wakePos, deg + PI + 0.2f + rand.nextSignedFloat(0.1f), sp, 40, sz * 32 * sr);
  }

  public Vector[] pointPos() {
    return _pointPos;
  }

  public float[] pointDeg() {
    return _pointDeg;
  }

  public bool checkShipCollision(float x, float y, float deg, float sr = 1) {
    float cs = size * (1 - distRatio) * 1.1f * sr;
    if (dist(x, y, 0, 0) < cs)
      return true;
    float ofs = 0;
    for (;;) {
      ofs += cs;
      cs *= distRatio;
      if (cs < 0.2f)
        return false;
      if (dist(x, y, sin(deg) * ofs, cos(deg) * ofs) < cs ||
          dist(x, y, -sin(deg) * ofs, -cos(deg) * ofs) < cs)
        return true;
    }
  }

  private float dist(float x, float y, float px, float py) {
    float ax = fabs(x - px);
    float ay = fabs(y - py);
    if (ax > ay)
      return ax + ay / 2;
    else
      return ay + ax / 2;
  }
}

public class CollidableBaseShape: BaseShape, Collidable {
  mixin CollidableImpl;
 private:
  Vector _collision;

  public this(float size, float distRatio, float spinyRatio,
              int type,
              float r, float g, float b) {
    super(size, distRatio, spinyRatio, type, r, g, b);
    _collision = new Vector(size / 2, size / 2);
  }

  public Vector collision() {
    return _collision;
  }
}

