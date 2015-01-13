/*
 * $Id: shape.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Shape of a ship/platform/turret/bridge.
 */

type ShapeType int

const (
  SHIP ShapeType = iota
  SHIP_ROUNDTAIL
  SHIP_SHADOW
  PLATFORM
  TURRET
  BRIDGE
  SHIP_DAMAGED
  SHIP_DESTROYED
  PLATFORM_DAMAGED
  PLATFORM_DESTROYED
  TURRET_DAMAGED
  TURRET_DESTROYED
)

public class ComplexShape: SimpleShape {
 public:
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

  public static this() {
    rand = new Rand;
    wakePos = new Vector;
  }

  public static void setRandSeed(long seed) {
    rand.setSeed(seed);
  }

  public this(float size, float distRatio, float spinyRatio,
              int type, float r, float g, float b, bool collidable = false) {
    this.size = size;
    this.distRatio = distRatio;
    this.spinyRatio = spinyRatio;
    this.type = type;
    this.r = r;
    this.g = g;
    this.b = b;
    if (collidable ) {
      collision = new Vector(size / 2, size / 2);
    }
    else{
      collision = nil
    }
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

public class TurretShape: ResizableShape {
 public:
  static enum TurretShapeType {
    NORMAL, DAMAGED, DESTROYED,
  };
 private:
  static ComplexShape[] shapes;

  public static void init() {
    shapes ~= new CollidableComplexShape(1, 0, 0, ComplexShape.ShapeType.TURRET, 1, 0.8f, 0.8f);
    shapes ~= new ComplexShape(1, 0, 0, ComplexShape.ShapeType.TURRET_DAMAGED, 0.9f, 0.9f, 1);
    shapes ~= new ComplexShape(1, 0, 0, ComplexShape.ShapeType.TURRET_DESTROYED, 0.8f, 0.33f, 0.66f);
  }

  public static void close() {
    foreach (ComplexShape s; shapes)
      s.close();
  }

  public this(int t) {
    shape = shapes[t];
  }
}

public class EnemyShape: ResizableShape {
 public:
  static enum EnemyShapeType {
    SMALL, SMALL_DAMAGED, SMALL_BRIDGE,
    MIDDLE, MIDDLE_DAMAGED, MIDDLE_DESTROYED, MIDDLE_BRIDGE,
    PLATFORM, PLATFORM_DAMAGED, PLATFORM_DESTROYED, PLATFORM_BRIDGE,
  };
  static const float MIDDLE_COLOR_R = 1, MIDDLE_COLOR_G = 0.6f, MIDDLE_COLOR_B = 0.5f;
 private:
  static ComplexShape[] shapes;

  public static void init() {
    shapes ~= new ComplexShape
      (1, 0.5f, 0.1f, ComplexShape.ShapeType.SHIP, 0.9f, 0.7f, 0.5f);
    shapes ~= new ComplexShape
      (1, 0.5f, 0.1f, ComplexShape.ShapeType.SHIP_DAMAGED, 0.5f, 0.5f, 0.9f);
    shapes ~= new CollidableComplexShape
      (0.66f, 0, 0, ComplexShape.ShapeType.BRIDGE, 1, 0.2f, 0.3f);
    shapes ~= new ComplexShape
      (1, 0.7f, 0.33f, ComplexShape.ShapeType.SHIP, MIDDLE_COLOR_R, MIDDLE_COLOR_G, MIDDLE_COLOR_B);
    shapes ~= new ComplexShape
      (1, 0.7f, 0.33f, ComplexShape.ShapeType.SHIP_DAMAGED, 0.5f, 0.5f, 0.9f);
    shapes ~= new ComplexShape
      (1, 0.7f, 0.33f, ComplexShape.ShapeType.SHIP_DESTROYED, 0, 0, 0);
    shapes ~= new CollidableComplexShape
      (0.66f, 0, 0, ComplexShape.ShapeType.BRIDGE, 1, 0.2f, 0.3f);
    shapes ~= new ComplexShape
      (1, 0, 0, ComplexShape.ShapeType.PLATFORM, 1, 0.6f, 0.7f);
    shapes ~= new ComplexShape
      (1, 0, 0, ComplexShape.ShapeType.PLATFORM_DAMAGED, 0.5f, 0.5f, 0.9f);
    shapes ~= new ComplexShape
      (1, 0, 0, ComplexShape.ShapeType.PLATFORM_DESTROYED, 1, 0.6f, 0.7f);
    shapes ~= new CollidableComplexShape
      (0.5f, 0, 0, ComplexShape.ShapeType.BRIDGE, 1, 0.2f, 0.3f);
  }

  public static void close() {
    foreach (ComplexShape s; shapes)
      s.close();
  }

  public this(int t) {
    shape = shapes[t];
  }

  public void addWake(WakePool wakes, Vector pos, float deg, float sp) {
    (cast(ComplexShape) shape).addWake(wakes, pos, deg, sp, size);
  }

  public bool checkShipCollision(float x, float y, float deg) {
    return (cast(ComplexShape) shape).checkShipCollision(x, y, deg, size);
  }
}

public class BulletShape: ResizableShape {
 public:
  static enum BulletShapeType {
    NORMAL, SMALL, MOVING_TURRET, DESTRUCTIVE,
  };
 private:
  static SimpleShape[] shapes;

  public static void init() {
    shapes ~= new NormalBulletShape;
    shapes ~= new SmallBulletShape;
    shapes ~= new MovingTurretBulletShape;
    shapes ~= new DestructiveBulletShape;
  }

  public static void close() {
    foreach (SimpleShape s; shapes)
      s.close();
  }

  public void set(int t) {
    shape = shapes[t];
  }
}

public class NormalBulletShape: SimpleShape {
  public override void createDisplayList() {
    glDisable(GL_BLEND);
    Screen.setColor(1, 1, 0.3f);
    glBegin(GL_LINE_STRIP);
    glVertex3f(0.2f, -0.25f, 0.2f);
    glVertex3f(0, 0.33f, 0);
    glVertex3f(-0.2f, -0.25f, -0.2f);
    glEnd();
    glBegin(GL_LINE_STRIP);
    glVertex3f(-0.2f, -0.25f, 0.2f);
    glVertex3f(0, 0.33f, 0);
    glVertex3f(0.2f, -0.25f, -0.2f);
    glEnd();
    glEnable(GL_BLEND);
    Screen.setColor(0.5f, 0.2f, 0.1f);
    glBegin(GL_TRIANGLE_FAN);
    glVertex3f(0, 0.33f, 0);
    glVertex3f(0.2f, -0.25f, 0.2f);
    glVertex3f(-0.2f, -0.25f, 0.2f);
    glVertex3f(-0.2f, -0.25f, -0.2f);
    glVertex3f(0.2f, -0.25f, -0.2f);
    glVertex3f(0.2f, -0.25f, 0.2f);
    glEnd();
  }
}

public class SmallBulletShape: SimpleShape {
  public override void createDisplayList() {
    glDisable(GL_BLEND);
    Screen.setColor(0.6f, 0.9f, 0.3f);
    glBegin(GL_LINE_STRIP);
    glVertex3f(0.25f, -0.25f, 0.25f);
    glVertex3f(0, 0.33f, 0);
    glVertex3f(-0.25f, -0.25f, -0.25f);
    glEnd();
    glBegin(GL_LINE_STRIP);
    glVertex3f(-0.25f, -0.25f, 0.25f);
    glVertex3f(0, 0.33f, 0);
    glVertex3f(0.25f, -0.25f, -0.25f);
    glEnd();
    glEnable(GL_BLEND);
    Screen.setColor(0.2f, 0.4f, 0.1f);
    glBegin(GL_TRIANGLE_FAN);
    glVertex3f(0, 0.33f, 0);
    glVertex3f(0.25f, -0.25f, 0.25f);
    glVertex3f(-0.25f, -0.25f, 0.25f);
    glVertex3f(-0.25f, -0.25f, -0.25f);
    glVertex3f(0.25f, -0.25f, -0.25f);
    glVertex3f(0.25f, -0.25f, 0.25f);
    glEnd();
  }
}

public class MovingTurretBulletShape: SimpleShape {
  public override void createDisplayList() {
    glDisable(GL_BLEND);
    Screen.setColor(0.7f, 0.5f, 0.9f);
    glBegin(GL_LINE_STRIP);
    glVertex3f(0.25f, -0.25f, 0.25f);
    glVertex3f(0, 0.33f, 0);
    glVertex3f(-0.25f, -0.25f, -0.25f);
    glEnd();
    glBegin(GL_LINE_STRIP);
    glVertex3f(-0.25f, -0.25f, 0.25f);
    glVertex3f(0, 0.33f, 0);
    glVertex3f(0.25f, -0.25f, -0.25f);
    glEnd();
    glEnable(GL_BLEND);
    Screen.setColor(0.2f, 0.2f, 0.3f);
    glBegin(GL_TRIANGLE_FAN);
    glVertex3f(0, 0.33f, 0);
    glVertex3f(0.25f, -0.25f, 0.25f);
    glVertex3f(-0.25f, -0.25f, 0.25f);
    glVertex3f(-0.25f, -0.25f, -0.25f);
    glVertex3f(0.25f, -0.25f, -0.25f);
    glVertex3f(0.25f, -0.25f, 0.25f);
    glEnd();
  }
}

public class DestructiveBulletShape: SimpleShape {
 private:
  Vector _collision;

  public override void createDisplayList() {
    glDisable(GL_BLEND);
    Screen.setColor(0.9f, 0.9f, 0.6f);
    glBegin(GL_LINE_LOOP);
    glVertex3f(0.2f, 0, 0);
    glVertex3f(0, 0.4f, 0);
    glVertex3f(-0.2f, 0, 0);
    glVertex3f(0, -0.4f, 0);
    glEnd();
    glEnable(GL_BLEND);
    Screen.setColor(0.7f, 0.5f, 0.4f);
    glBegin(GL_TRIANGLE_FAN);
    glVertex3f(0.2f, 0, 0);
    glVertex3f(0, 0.4f, 0);
    glVertex3f(-0.2f, 0, 0);
    glVertex3f(0, -0.4f, 0);
    glEnd();
    _collision = new Vector(0.4f, 0.4f);
  }

  public Vector collision() {
    return _collision;
  }
}

public class CrystalShape: SimpleShape {
  public override void createDisplayList() {
    Screen.setColor(0.6f, 1, 0.7f);
    glBegin(GL_LINE_LOOP);
    glVertex3f(-0.2f, 0.2f, 0);
    glVertex3f(0.2f, 0.2f, 0);
    glVertex3f(0.2f, -0.2f, 0);
    glVertex3f(-0.2f, -0.2f, 0);
    glEnd();
  }
}

public class ShieldShape: SimpleShape {
  public override void createDisplayList() {
    Screen.setColor(0.5f, 0.5f, 0.7f);
    glBegin(GL_LINE_LOOP);
    float d = 0;
    for (int i = 0; i < 8; i++) {
      glVertex3f(sin(d), cos(d), 0);
      d += PI / 4;
    }
    glEnd();
    glBegin(GL_TRIANGLE_FAN);
    Screen.setColor(0, 0, 0);
    glVertex3f(0, 0, 0);
    d = 0;
    Screen.setColor(0.3f, 0.3f, 0.5f);
    for (int i = 0; i < 9; i++) {
      glVertex3f(sin(d), cos(d), 0);
      d += PI / 4;
    }
    glEnd();
  }
}


/**
 * Interface for drawing a shape.
 */
public interface Shape {
  public void draw();
  public Vector collision();
  public bool checkCollision(float ax, float ay, Shape shape = null);
}

/* just a displaylist
   and a simple static collision, if collidable */
public template SimpleShape() {
  protected DisplayList displayList;
  collision Vector

  public bool checkCollision(float ax, float ay, Shape shape = null) {
    float cx, cy;
    if (shape) {
      cx = collision.x + shape.collision.x;
      cy = collision.y + shape.collision.y;
    } else {
      cx = collision.x;
      cy = collision.y;
    }
    if (ax <= cx && ay <= cy)
      return true;
    else
      return false;
  }

  public this() {
    displayList = new DisplayList(1);
    displayList.beginNewList();
    createDisplayList();
    displayList.endNewList();
  }

  protected abstract void createDisplayList();

  public Vector collision() {
    return collision;
  }

  public void close() {
    displayList.close();
  }

  public void draw() {
    displayList.call(0);
  }
}

/*
 * a Shape that can change a size.
 *
 * proxies a Simple or Complex shape
 */
type ResizableShape struct {
  shape Shape
  size float32
  resizedCollision Vector
}

func (rd *ResizableShape) draw() {
  gl.Scalef(rs.size, rs.size, rs.size)
  rs.shape.Draw()
}

func (rd *ResizableShape) collision() *Vector {
  rs.resizedCollision = NewVector(cd.collision().X() * rs.size, cd.collision().Y() * rs.size)
  return rs.collision
}

public bool checkCollision(float ax, float ay, Shape shape = null) {
  float cx, cy;
  if (shape) {
    cx = collision.x + shape.collision.x;
    cy = collision.y + shape.collision.y;
  } else {
    cx = collision.x;
    cy = collision.y;
  }
  if (ax <= cx && ay <= cy)
    return true;
  else
    return false;
}

