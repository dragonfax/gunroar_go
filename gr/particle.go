package main

import "github.com/dragonfax/gunroar/gr/actor"

/**
 * Sparks.
 */

var sparkRand = r.New(r.NewSource(time.Now().Unix()))

var _ LuminousActor = &Spark{}

type Spark struct{
  actor.ExistsImpl

  pos, ppos, vel vector.Vector;
  r, g, b float64;
  cnt int;
}

  func setSparkRandSeed(seed int64) {
    sparkRand = r.New(r.NewSource(seed))
  }

  func NewSpark() *Spark {
    return &Spark{}
  }

  func (*Spark) Init(args []interface{}) {
  }

  func (this *Spark) set(p Vector, vx, vy, r, g, b float64, c int) {
    this.ppos.X = p.X
    this.pos.X = p.X;
    this.ppos.Y = p.Y
    this.pos.Y = p.Y;
    this.vel.X = vx;
    this.vel.Y = vy;
    this.r = r;
    this.g = g;
    this.b = b;
    this.cnt = c;
    this.SetExists(true)
  }

  func (this *Spark) move() {
    this.cnt--;
    if this.cnt <= 0 || this.vel.dist(0,0) < 0.005 {
      this.SetExists(false)
      return;
    }
    this.ppos.X = this.pos.X;
    this.ppos.Y = this.pos.Y;
    this.pos += this.vel;
    this.vel *= 0.96;
  }

  func (this *Spark) draw() {
    ox := this.vel.X;
    oy := this.vel.Y;
    sdl.SetColor(this.r, this.g, this.b, 1);
    ox *= 2;
    oy *= 2;
    glVertex3d(this.pos.X - ox, this.pos.Y - oy, 0);
    ox *= 0.5;
    oy *= 0.5;
    sdl.SetColor(this.r * 0.5, this.g * 0.5, this.b * 0.5, 0);
    glVertex3d(this.pos.X - oy, this.pos.Y + ox, 0);
    glVertex3d(this.pos.X + oy, this.pos.Y - ox, 0);
  }

  func (this *Spark) drawLuminous() {
    ox := this.vel.X;
    oy := this.vel.Y;
    sdl.SetColor(this.r, this.g, this.b, 1);
    ox *= 2;
    oy *= 2;
    glVertex3d(this.pos.X - ox, this.pos.Y - oy, 0);
    ox *= 0.5;
    oy *= 0.5;
    sdl.SetColor(this.r * 0.5, this.g * 0.5, this.b * 0.5, 0);
    glVertex3d(this.pos.X - oy, this.pos.Y + ox, 0);
    glVertex3d(this.pos.X + oy, this.pos.Y - ox, 0);
  }

type SparkPool struct {
  actor.ActorPool
}
func NewSparkPool(n int, args []interface{}) *SparkPool {
  f := func () actor.Actor { return NewSpark() }
  this := &SparkPool{ActorPool: NewActorPool(f, n, args)}
  return this
}

/**
 * Smokes.
 */

 type SmokeType  int

 const(
    FIRE SmokeType = iota
    EXPLOSION
    SAND
    SPARK
    WAKE
    SMOKE
    LANCE_SPARK
 )

var smokeRand = r.New(r.NewSource(time.Now().Unix()))
var windVel = vector.Vector3{0.04, 0.04, 0.02};
var wakePos vector.Vector;

var _ LuminousActor = &Smoke{}

type Smoke struct {
  actor.ExistsImpl

  field Field;
  wakes WakePool
  pos, vel vector.Vector3;
  typ SmokeType
  cnt, startCnt int;
  size, r, g, b, a float64;
}

  func setSmokeRandSeed(seed int) {
    smokeRand = r.New(r.NewSource(seed))
  }

  func NewSmoke() *Smoke{
    this := &Smoke{
      startCnt:1,
      size: 1,
    }
    return this
  }

  func (this *Smoke) Init(args []interface{}) {
    this.field = args[0].(*Field);
    this.wakes = args[1].(*WakePool);
  }

 func (this *Smoke) setVector(p vector.Vector, mx, my, mz float64 , t int, c int /* = 60 */, sz float64 /* = 2 */ ) {
    this.set(p.X, p.Y, mx, my, mz, t, c, sz);
  }

  func (this *Smoke) setVector3(p vector.Vector3, mx, my, mz float64, t int, c int /* = 60 */, sz float64 /* = 2 */) {
    this.set(p.X, p.Y, mx, my, mz, t, c, sz);
    this.pos.Z = p.Z;
  }

  func (this *Smoke) set(x, y, mx, my, mz float64, t int, c int /* = 60 */, sz float64 /* = 2 */) {
    if !this.field.checkInOuterField(x, y) {
      return;
    }
    this.pos.X = x;
    this.pos.Y = y;
    this.pos.Z = 0;
    this.vel.X = mx;
    this.vel.Y = my;
    this.vel.Z = mz;
    this.typ = t;
    this.startCnt = c
    this.cnt = c;
    this.size = sz;
    switch (this.typ) {
    case FIRE:
      this.r = nextFloat(rand,0.1) + 0.9;
      this.g = nextFloat(rand,0.2) + 0.2;
      this.b = 0;
      this.a = 1;
    case EXPLOSION:
      this.r = nextFloat(rand,0.3) + 0.7;
      this.g = nextFloat(rand,0.3) + 0.3;
      this.b = 0;
      this.a = 1;
    case SAND:
      this.r = 0.8;
      this.g = 0.8;
      this.b = 0.6;
      this.a = 0.6;
    case SPARK:
      this.r = nextFloat(rand,0.3) + 0.7;
      this.g = nextFloat(rand,0.5) + 0.5;
      this.b = 0;
      this.a = 1;
    case WAKE:
      this.r = 0.6;
      this.g = 0.6;
      this.b = 0.8;
      this.a = 0.6;
    case SMOKE:
      this.r = nextFloat(rand,0.1) + 0.1;
      this.g = nextFloat(rand,0.1) + 0.1;
      this.b = 0.1;
      this.a = 0.5;
    case LANCE_SPARK:
      this.r = 0.4;
      this.g = nextFloat(rand,0.2) + 0.7;
      this.b = nextFloat(rand,0.2) + 0.7;
      this.a = 1;
    }
    this.SetExists(true);
  }

  func (this *Smoke) move() {
    this.cnt--;
    if this.cnt <= 0 || !this.field.checkInOuterField(this.pos.X, this.pos.Y) {
      this.SetExists(false);
      return;
    }
    if this.typ != WAKE {
      this.vel.X += (windVel.X - this.vel.X) * 0.01;
      this.vel.Y += (windVel.Y - this.vel.Y) * 0.01;
      this.vel.Z += (windVel.Z - this.vel.Z) * 0.01;
    }
    this.pos += this.vel;
    this.pos.Y -= this.field.lastScrollY;
    switch (this.typ) {
    case FIRE, EXPLOSION, SMOKE:
      if (cnt < startCnt / 2) {
        this.r *= 0.95;
        this.g *= 0.95;
        this.b *= 0.95;
      } else {
        this.a *= 0.97;
      }
      this.size *= 1.01;
    case SAND:
      this.r *= 0.98;
      this.g *= 0.98;
      this.b *= 0.98;
      this.a *= 0.98;
    case SPARK:
      this.r *= 0.92;
      this.g *= 0.92;
      this.a *= 0.95;
      this.vel *= 0.9;
    case WAKE:
      this.a *= 0.98;
      this.size *= 1.005;
    case LANCE_SPARK:
      this.a *= 0.95;
      this.size *= 0.97;
    }
    if this.size > 5 {
      this.size = 5;
    }
    if this.typ == EXPLOSION && this.pos.Z < 0.01 {
      bl := this.field.getBlock(this.pos.X, this.pos.Y);
      if bl >= 1 {
        this.vel *= 0.8;
      }
      if this.cnt % 3 == 0 && bl < -1 {
        sp := math.Sqrt(this.vel.X * this.vel.X + this.vel.Y * this.vel.Y);
        if sp > 0.3 {
          d := math.Atan2(this.vel.X, this.vel.Y);
          wakePos.X = this.pos.X + math.Sin(d + math.Pi / 2) * this.size * 0.25;
          wakePos.Y = this.pos.Y + math.Cos(d + math.Pi / 2) * this.size * 0.25;
          w := wakes.getInstanceForced();
          w.set(wakePos, d + math.Pi - 0.2 + nextSignedFloat(rand,0.1), sp * 0.33,
                20 + rand.Intn(12), this.size * (7.0 + nextFloat(rand,3)));
          wakePos.X = this.pos.X + math.Sin(d - math.Pi / 2) * this.size * 0.25;
          wakePos.Y = this.pos.Y + math.Cos(d - math.Pi / 2) * this.size * 0.25;
          w = wakes.getInstanceForced();
          w.set(wakePos, d + math.Pi + 0.2 + nextSignedFloat(rand,0.1), sp * 0.33,
                20 + rand.Intn(12), this.size * (7.0 + nextFloat(rand,3)));
        }
      }
    }
  }

  func (this *Smoke) draw() {
    quadSize := this.size / 2;
    sdl.SetColor(this.r, this.g, this.b, this.a);
    glVertex3d(this.pos.X - quadSize, this.pos.Y - quadSize, this.pos.Z);
    glVertex3d(this.pos.X + quadSize, this.pos.Y - quadSize, this.pos.Z);
    glVertex3d(this.pos.X + quadSize, this.pos.Y + quadSize, this.pos.Z);
    glVertex3d(this.pos.X - quadSize, this.pos.Y + quadSize, this.pos.Z);
  }

  func (this *Smoke) drawLuminous() {
    if this.r + this.g > 0.8 && this.b < 0.5 {
      quadSize := this.size / 2;
      sdl.SetColor(this.r, this.g, this.b, this.a);
      glVertex3d(this.pos.X - quadSize, this.pos.Y - quadSize, this.pos.Z);
      glVertex3d(this.pos.X + quadSize, this.pos.Y - quadSize, this.pos.Z);
      glVertex3d(this.pos.X + quadSize, this.pos.Y + quadSize, this.pos.Z);
      glVertex3d(this.pos.X - quadSize, this.pos.Y + quadSize, this.pos.Z);
    }
  }

var _ LuminousActorPool = &SmokePool{}

type SmokePool struct {
  actor.ActorPool
}

func NewSmokePool(n int, args []interface{}) *SmokePool {
  this := &SmokePool{ ActorPool: NewActorPool(NewSmoke,n, args) }
  return this
}

/**
 * Fragments of destroyed enemies.
 */
public class Fragment: Actor {
 private:
  static DisplayList displayList;
  static Rand rand;
  Field field;
  SmokePool smokes;
  Vector3 pos;
  Vector3 vel;
  float size;
  float d2, md2;

  invariant {
    assert(pos.x < 15 && pos.x > -15);
    assert(pos.y < 20 && pos.y > -20);
    assert(pos.z < 20 && pos.z > -10);
    assert(vel.x < 10 && vel.x > -10);
    assert(vel.y < 10 && vel.y > -10);
    assert(vel.z < 10 && vel.z > -10);
    assert(size >= 0 && size < 10);
    assert(d2 <>= 0);
    assert(md2 <>= 0);
  }

  public static void init() {
    rand = new Rand;
    displayList = new DisplayList(1);
    displayList.beginNewList();
    Screen.setColor(0.7f, 0.5f, 0.5f, 0.5f);
    glBegin(GL_TRIANGLE_FAN);
    glVertex2f(-0.5f, -0.25f);
    glVertex2f(0.5f, -0.25f);
    glVertex2f(0.5f, 0.25f);
    glVertex2f(-0.5f, 0.25f);
    glEnd();
    Screen.setColor(0.7f, 0.5f, 0.5f, 0.9f);
    glBegin(GL_LINE_LOOP);
    glVertex2f(-0.5f, -0.25f);
    glVertex2f(0.5f, -0.25f);
    glVertex2f(0.5f, 0.25f);
    glVertex2f(-0.5f, 0.25f);
    glEnd();
    displayList.endNewList();
  }

  public static void setRandSeed(long seed) {
    rand.setSeed(seed);
  }

  public this() {
    pos = new Vector3;
    vel = new Vector3;
    size = 1;
    d2 = md2 = 0;
  }

  public override void init(Object[] args) {
    field = cast(Field) args[0];
    smokes = cast(SmokePool) args[1];
  }

  public void set(Vector p, float mx, float my, float mz, float sz = 1) {
    if (!field.checkInOuterField(p.x, p.y))
      return;
    pos.x = p.x;
    pos.y = p.y;
    pos.z = 0;
    vel.x = mx;
    vel.y = my;
    vel.z = mz;
    size = sz;
    if (size > 5)
      size = 5;
    d2 = rand.nextFloat(360);
    md2 = rand.nextSignedFloat(20);
    exists = true;
  }

  public override void move() {
    if (!field.checkInOuterField(pos.x, pos.y)) {
      exists = false;
      return;
    }
    vel.x *= 0.96f;
    vel.y *= 0.96f;
    vel.z += (-0.04f - vel.z) * 0.01f;
    pos += vel;
    if (pos.z < 0) {
      Smoke s = smokes.getInstanceForced();
      if (field.getBlock(pos.x, pos.y) < 0)
        s.set(pos.x, pos.y, 0, 0, 0, Smoke.SmokeType.WAKE, 60, size * 0.66f);
      else
        s.set(pos.x, pos.y, 0, 0, 0, Smoke.SmokeType.SAND, 60, size * 0.75f);
      exists = false;
      return;
    }
    pos.y -= field.lastScrollY;
    d2 += md2;
  }

  public override void draw() {
    glPushMatrix();
    Screen.glTranslate(pos);
    glRotatef(d2, 1, 0, 0);
    glScalef(size, size, 1);
    displayList.call(0);
    glPopMatrix();
  }
}

public class FragmentPool: ActorPool!(Fragment) {
  public this(int n, Object[] args) {
    super(n, args);
  }
}

/**
 * Luminous fragments.
 */
public class SparkFragment: LuminousActor {
 private:
  static DisplayList displayList;
  static Rand rand;
  Field field;
  SmokePool smokes;
  Vector3 pos;
  Vector3 vel;
  float size;
  float d2, md2;
  int cnt;
  bool hasSmoke;

  invariant {
    assert(pos.x < 15 && pos.x > -15);
    assert(pos.y < 20 && pos.y > -20);
    assert(pos.z < 20 && pos.z > -10);
    assert(vel.x < 10 && vel.x > -10);
    assert(vel.y < 10 && vel.y > -10);
    assert(vel.z < 10 && vel.z > -10);
    assert(size >= 0 && size < 10);
    assert(d2 <>= 0);
    assert(md2 <>= 0);
    assert(cnt >= 0);
  }

  public static void init() {
    rand = new Rand;
    displayList = new DisplayList(1);
    displayList.beginNewList();
    glBegin(GL_TRIANGLE_FAN);
    glVertex2f(-0.25f, -0.25f);
    glVertex2f(0.25f, -0.25f);
    glVertex2f(0.25f, 0.25f);
    glVertex2f(-0.25f, 0.25f);
    glEnd();
    displayList.endNewList();
  }

  public static void setRandSeed(long seed) {
    rand.setSeed(seed);
  }


  public this() {
    pos = new Vector3;
    vel = new Vector3;
    size = 1;
    d2 = md2 = 0;
    cnt = 0;
  }

  public override void init(Object[] args) {
    field = cast(Field) args[0];
    smokes = cast(SmokePool) args[1];
  }

  public void set(Vector p, float mx, float my, float mz, float sz = 1) {
    if (!field.checkInOuterField(p.x, p.y))
      return;
    pos.x = p.x;
    pos.y = p.y;
    pos.z = 0;
    vel.x = mx;
    vel.y = my;
    vel.z = mz;
    size = sz;
    if (size > 5)
      size = 5;
    d2 = rand.nextFloat(360);
    md2 = rand.nextSignedFloat(15);
    if (rand.nextInt(4) == 0)
      hasSmoke = true;
    else
      hasSmoke = false;
    cnt = 0;
    exists = true;
  }

  public override void move() {
    if (!field.checkInOuterField(pos.x, pos.y)) {
      exists = false;
      return;
    }
    vel.x *= 0.99f;
    vel.y *= 0.99f;
    vel.z += (-0.08f - vel.z) * 0.01f;
    pos += vel;
    if (pos.z < 0) {
      Smoke s = smokes.getInstanceForced();
      if (field.getBlock(pos.x, pos.y) < 0)
        s.set(pos.x, pos.y, 0, 0, 0, Smoke.SmokeType.WAKE, 60, size * 0.66f);
      else
        s.set(pos.x, pos.y, 0, 0, 0, Smoke.SmokeType.SAND, 60, size * 0.75f);
      exists = false;
      return;
    }
    pos.y -= field.lastScrollY;
    d2 += md2;
    cnt++;
    if (hasSmoke && cnt % 5 == 0) {
      Smoke s = smokes.getInstance();
      if (s)
        s.set(pos, 0, 0, 0, Smoke.SmokeType.SMOKE, 90 + rand.nextInt(60), size * 0.5f);
    }
  }

  public override void draw() {
    glPushMatrix();
    Screen.setColor(1, rand.nextFloat(1), 0, 0.8f);
    Screen.glTranslate(pos);
    glRotatef(d2, 1, 0, 0);
    glScalef(size, size, 1);
    displayList.call(0);
    glPopMatrix();
  }

  public override void drawLuminous() {
    glPushMatrix();
    Screen.setColor(1, rand.nextFloat(1), 0, 0.8f);
    Screen.glTranslate(pos);
    glRotatef(d2, 1, 0, 0);
    glScalef(size, size, 1);
    displayList.call(0);
    glPopMatrix();
  }
}

public class SparkFragmentPool: LuminousActorPool!(SparkFragment) {
  public this(int n, Object[] args) {
    super(n, args);
  }
}

/**
 * Wakes of ships and smokes.
 */
public class Wake: Actor {
 private:
  Field field;
  Vector pos;
  Vector vel;
  float deg;
  float speed;
  float size;
  int cnt;
  bool revShape;

  invariant {
    assert(pos.x < 15 && pos.x > -15);
    assert(pos.y < 20 && pos.y > -20);
    assert(vel.x < 10 && vel.x > -10);
    assert(vel.y < 10 && vel.y > -10);
    assert(size > 0 && size < 1000);
    assert(deg <>= 0);
    assert(speed >= 0 && speed < 10);
    assert(cnt >= 0);
  }

  public this() {
    pos = new Vector;
    vel = new Vector;
    size = 1;
    deg = 0;
    speed = 0;
    cnt = 0;
  }

  public override void init(Object[] args) {
    field = cast(Field) args[0];
  }

  public void set(Vector p, float deg, float speed, int c = 60, float sz = 1, bool rs = false) {
    if (!field.checkInOuterField(p.x, p.y))
      return;
    pos.x = p.x;
    pos.y = p.y;
    this.deg = deg;
    this.speed = speed;
    vel.x = sin(deg) * speed;
    vel.y = cos(deg) * speed;
    cnt = c;
    size = sz;
    revShape = rs;
    exists = true;
  }

  public override void move() {
    cnt--;
    if (cnt <= 0 || vel.dist(0,0) < 0.005f || !field.checkInOuterField(pos.x, pos.y)) {
      exists = false;
      return;
    }
    pos += vel;
    pos.y -= field.lastScrollY;
    vel *= 0.96f;
    size *= 1.02f;
  }

  public override void draw() {
    float ox = vel.x;
    float oy = vel.y;
    Screen.setColor(0.33f, 0.33f, 1);
    ox *= size;
    oy *= size;
    if (revShape)
      glVertex3f(pos.x + ox, pos.y + oy, 0);
    else
      glVertex3f(pos.x - ox, pos.y - oy, 0);
    ox *= 0.2f;
    oy *= 0.2f;
    Screen.setColor(0.2f, 0.2f, 0.6f, 0.5f);
    glVertex3f(pos.x - oy, pos.y + ox, 0);
    glVertex3f(pos.x + oy, pos.y - ox, 0);
  }
}

public class WakePool: ActorPool!(Wake) {
  public this(int n, Object[] args) {
    super(n, args);
  }
}
