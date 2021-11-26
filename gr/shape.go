package main

type TurretShapeType int 

const (
    NORMAL TurretShapeType = iota + 1
    DAMAGED
    DESTROYED
)

var turretShapes []*BaseShape

type TurretShape struct {
  ResizableDrawable
}

func TurretInit() {
  turretShapes = append(turretShapes, NewCollidableBaseShape(1, 0, 0, BaseShape.ShapeType.TURRET, 1, 0.8, 0.8));
  turretShapes = append(turretShapes, NewBaseShape(1, 0, 0, BaseShape.ShapeType.TURRET_DAMAGED, 0.9, 0.9, 1));
  turretShapes = append(turretShapes, NewBaseShape(1, 0, 0, BaseShape.ShapeType.TURRET_DESTROYED, 0.8, 0.33, 0.66));
}

  public static void close() {
    foreach (BaseShape s; shapes)
      s.close();
  }

  public this(int t) {
    shape = shapes[t];
  }
}

public class EnemyShape: ResizableDrawable {
 public:
  static enum EnemyShapeType {
    SMALL, SMALL_DAMAGED, SMALL_BRIDGE,
    MIDDLE, MIDDLE_DAMAGED, MIDDLE_DESTROYED, MIDDLE_BRIDGE,
    PLATFORM, PLATFORM_DAMAGED, PLATFORM_DESTROYED, PLATFORM_BRIDGE,
  };
  static const float MIDDLE_COLOR_R = 1, MIDDLE_COLOR_G = 0.6, MIDDLE_COLOR_B = 0.5;
 private:
  static BaseShape[] shapes;

  public static void init() {
    shapes ~= new BaseShape
      (1, 0.5, 0.1, BaseShape.ShapeType.SHIP, 0.9, 0.7, 0.5);
    /hapes ~= new BaseShape
      (1, 0.5, 0.1, BaseShape.ShapeType.SHIP_DAMAGED, 0.5, 0.5, 0.9);
    shapes ~= new CollidableBaseShape
      (0.66, 0, 0, BaseShape.ShapeType.BRIDGE, 1, 0.2, 0.3);
    shapes ~= new BaseShape
      (1, 0.7, 0.33, BaseShape.ShapeType.SHIP, MIDDLE_COLOR_R, MIDDLE_COLOR_G, MIDDLE_COLOR_B);
    shapes ~= new BaseShape
      (1, 0.7, 0.33, BaseShape.ShapeType.SHIP_DAMAGED, 0.5, 0.5, 0.9);
    shapes ~= new BaseShape
      (1, 0.7, 0.33, BaseShape.ShapeType.SHIP_DESTROYED, 0, 0, 0);
    shapes ~= new CollidableBaseShape
      (0.66, 0, 0, BaseShape.ShapeType.BRIDGE, 1, 0.2, 0.3);
    shapes ~= new BaseShape
      (1, 0, 0, BaseShape.ShapeType.PLATFORM, 1, 0.6, 0.7);
    shapes ~= new BaseShape
      (1, 0, 0, BaseShape.ShapeType.PLATFORM_DAMAGED, 0.5, 0.5, 0.9);
    shapes ~= new BaseShape
      (1, 0, 0, BaseShape.ShapeType.PLATFORM_DESTROYED, 1, 0.6, 0.7);
    shapes ~= new CollidableBaseShape
      (0.5, 0, 0, BaseShape.ShapeType.BRIDGE, 1, 0.2, 0.3);
  }

  public static void close() {
    foreach (BaseShape s; shapes)
      s.close();
  }

  public this(int t) {
    shape = shapes[t];
  }

  public void addWake(WakePool wakes, Vector pos, float deg, float sp) {
    (cast(BaseShape) shape).addWake(wakes, pos, deg, sp, size);
  }

  public bool checkShipCollision(float x, float y, float deg) {
    return (cast(BaseShape) shape).checkShipCollision(x, y, deg, size);
  }
}

public class BulletShape: ResizableDrawable {
 public:
  static enum BulletShapeType {
    NORMAL, SMALL, MOVING_TURRET, DESTRUCTIVE,
  };
 private:
  static DrawableShape[] shapes;

  public static void init() {
    shapes ~= new NormalBulletShape;
    shapes ~= new SmallBulletShape;
    shapes ~= new MovingTurretBulletShape;
    shapes ~= new DestructiveBulletShape;
  }

  public static void close() {
    foreach (DrawableShape s; shapes)
      s.close();
  }

  public void set(int t) {
    shape = shapes[t];
  }
}

public class NormalBulletShape: DrawableShape {
  public override void createDisplayList() {
    glDisable(GL_BLEND);
    Screen.setColor(1, 1, 0.3);
    glBegin(GL_LINE_STRIP);
    glVertex3f(0.2, -0.25, 0.2);
    glVertex3f(0, 0.33, 0);
    glVertex3f(-0.2, -0.25, -0.2);
    glEnd();
    glBegin(GL_LINE_STRIP);
    glVertex3f(-0.2, -0.25, 0.2);
    glVertex3f(0, 0.33, 0);
    glVertex3f(0.2, -0.25, -0.2);
    glEnd();
    glEnable(GL_BLEND);
    Screen.setColor(0.5, 0.2, 0.1);
    glBegin(GL_TRIANGLE_FAN);
    glVertex3f(0, 0.33, 0);
    glVertex3f(0.2, -0.25, 0.2);
    glVertex3f(-0.2, -0.25, 0.2);
    glVertex3f(-0.2, -0.25, -0.2);
    glVertex3f(0.2, -0.25, -0.2);
    glVertex3f(0.2, -0.25, 0.2);
    glEnd();
  }
}

public class SmallBulletShape: DrawableShape {
  public override void createDisplayList() {
    glDisable(GL_BLEND);
    Screen.setColor(0.6, 0.9, 0.3);
    glBegin(GL_LINE_STRIP);
    glVertex3f(0.25, -0.25, 0.25);
    glVertex3f(0, 0.33, 0);
    glVertex3f(-0.25, -0.25, -0.25);
    glEnd();
    glBegin(GL_LINE_STRIP);
    glVertex3f(-0.25, -0.25, 0.25);
    glVertex3f(0, 0.33, 0);
    glVertex3f(0.25, -0.25, -0.25);
    glEnd();
    glEnable(GL_BLEND);
    Screen.setColor(0.2, 0.4, 0.1);
    glBegin(GL_TRIANGLE_FAN);
    glVertex3f(0, 0.33, 0);
    glVertex3f(0.25, -0.25, 0.25);
    glVertex3f(-0.25, -0.25, 0.25);
    glVertex3f(-0.25, -0.25, -0.25);
    glVertex3f(0.25, -0.25, -0.25);
    glVertex3f(0.25, -0.25, 0.25);
    glEnd();
  }
}

public class MovingTurretBulletShape: DrawableShape {
  public override void createDisplayList() {
    glDisable(GL_BLEND);
    Screen.setColor(0.7, 0.5, 0.9);
    glBegin(GL_LINE_STRIP);
    glVertex3f(0.25, -0.25, 0.25);
    glVertex3f(0, 0.33, 0);
    glVertex3f(-0.25, -0.25, -0.25);
    glEnd();
    glBegin(GL_LINE_STRIP);
    glVertex3f(-0.25, -0.25, 0.25);
    glVertex3f(0, 0.33, 0);
    glVertex3f(0.25, -0.25, -0.25);
    glEnd();
    glEnable(GL_BLEND);
    Screen.setColor(0.2, 0.2, 0.3);
    glBegin(GL_TRIANGLE_FAN);
    glVertex3f(0, 0.33, 0);
    glVertex3f(0.25, -0.25, 0.25);
    glVertex3f(-0.25, -0.25, 0.25);
    glVertex3f(-0.25, -0.25, -0.25);
    glVertex3f(0.25, -0.25, -0.25);
    glVertex3f(0.25, -0.25, 0.25);
    glEnd();
  }
}

public class DestructiveBulletShape: DrawableShape, Collidable {
  mixin CollidableImpl;
 private:
  Vector _collision;

  public override void createDisplayList() {
    glDisable(GL_BLEND);
    Screen.setColor(0.9, 0.9, 0.6);
    glBegin(GL_LINE_LOOP);
    glVertex3f(0.2, 0, 0);
    glVertex3f(0, 0.4, 0);
    glVertex3f(-0.2, 0, 0);
    glVertex3f(0, -0.4, 0);
    glEnd();
    glEnable(GL_BLEND);
    Screen.setColor(0.7, 0.5, 0.4);
    glBegin(GL_TRIANGLE_FAN);
    glVertex3f(0.2, 0, 0);
    glVertex3f(0, 0.4, 0);
    glVertex3f(-0.2, 0, 0);
    glVertex3f(0, -0.4, 0);
    glEnd();
    _collision = new Vector(0.4, 0.4);
  }

  public Vector collision() {
    return _collision;
  }
}

public class CrystalShape: DrawableShape {
  public override void createDisplayList() {
    Screen.setColor(0.6, 1, 0.7);
    glBegin(GL_LINE_LOOP);
    glVertex3f(-0.2, 0.2, 0);
    glVertex3f(0.2, 0.2, 0);
    glVertex3f(0.2, -0.2, 0);
    glVertex3f(-0.2, -0.2, 0);
    glEnd();
  }
}

public class ShieldShape: DrawableShape {
  public override void createDisplayList() {
    Screen.setColor(0.5, 0.5, 0.7);
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
    Screen.setColor(0.3, 0.3, 0.5);
    for (int i = 0; i < 9; i++) {
      glVertex3f(sin(d), cos(d), 0);
      d += PI / 4;
    }
    glEnd();
  }
}
