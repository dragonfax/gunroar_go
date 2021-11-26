public class TurretShape: ResizableDrawable {
 public:
  static enum TurretShapeType {
    NORMAL, DAMAGED, DESTROYED,
  };
 private:
  static BaseShape[] shapes;

  public static void init() {
    shapes ~= new CollidableBaseShape(1, 0, 0, BaseShape.ShapeType.TURRET, 1, 0.8f, 0.8f);
    shapes ~= new BaseShape(1, 0, 0, BaseShape.ShapeType.TURRET_DAMAGED, 0.9f, 0.9f, 1);
    shapes ~= new BaseShape(1, 0, 0, BaseShape.ShapeType.TURRET_DESTROYED, 0.8f, 0.33f, 0.66f);
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
  static const float MIDDLE_COLOR_R = 1, MIDDLE_COLOR_G = 0.6f, MIDDLE_COLOR_B = 0.5f;
 private:
  static BaseShape[] shapes;

  public static void init() {
    shapes ~= new BaseShape
      (1, 0.5f, 0.1f, BaseShape.ShapeType.SHIP, 0.9f, 0.7f, 0.5f);
    shapes ~= new BaseShape
      (1, 0.5f, 0.1f, BaseShape.ShapeType.SHIP_DAMAGED, 0.5f, 0.5f, 0.9f);
    shapes ~= new CollidableBaseShape
      (0.66f, 0, 0, BaseShape.ShapeType.BRIDGE, 1, 0.2f, 0.3f);
    shapes ~= new BaseShape
      (1, 0.7f, 0.33f, BaseShape.ShapeType.SHIP, MIDDLE_COLOR_R, MIDDLE_COLOR_G, MIDDLE_COLOR_B);
    shapes ~= new BaseShape
      (1, 0.7f, 0.33f, BaseShape.ShapeType.SHIP_DAMAGED, 0.5f, 0.5f, 0.9f);
    shapes ~= new BaseShape
      (1, 0.7f, 0.33f, BaseShape.ShapeType.SHIP_DESTROYED, 0, 0, 0);
    shapes ~= new CollidableBaseShape
      (0.66f, 0, 0, BaseShape.ShapeType.BRIDGE, 1, 0.2f, 0.3f);
    shapes ~= new BaseShape
      (1, 0, 0, BaseShape.ShapeType.PLATFORM, 1, 0.6f, 0.7f);
    shapes ~= new BaseShape
      (1, 0, 0, BaseShape.ShapeType.PLATFORM_DAMAGED, 0.5f, 0.5f, 0.9f);
    shapes ~= new BaseShape
      (1, 0, 0, BaseShape.ShapeType.PLATFORM_DESTROYED, 1, 0.6f, 0.7f);
    shapes ~= new CollidableBaseShape
      (0.5f, 0, 0, BaseShape.ShapeType.BRIDGE, 1, 0.2f, 0.3f);
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

public class SmallBulletShape: DrawableShape {
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

public class MovingTurretBulletShape: DrawableShape {
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

public class DestructiveBulletShape: DrawableShape, Collidable {
  mixin CollidableImpl;
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

public class CrystalShape: DrawableShape {
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

public class ShieldShape: DrawableShape {
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
