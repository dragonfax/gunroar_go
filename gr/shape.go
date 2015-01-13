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

const POINT_NUM = 16
const PILLAR_POINT_NUM = 8
var wakePos = Vector{}

type ComplexShape struct {
	*SimpleShape

  size, distRatio, spinyRatio float32
  shapeType ShapeType
  r, g, b float32
  pillarPos []Vector
  pointPos []Vector
  pointDeg []float32
}

func (this *ComplexShape) InitComplexShape( size float32, distRatio float32, spinyRatio float32, shapeType int, r float32, g float32, b float32, collidable bool /* = false */) {

	this.size = size
	this.distRatio = distRatio
	this.spinyRatio = spinyRatio
	this.shapeType = shapeType
	this.r = r
	this.g = g
	this.b = b
	if (collidable ) {
		this.collision = Vector{size / 2, size / 2}
	} else{
		this.collision = nil
	}
	this.InitSimpleShape()
}

func (this *ComplexShape)  createDisplayList() {
	height := this.size * 0.5
	var z float32 = 0
	var sz float32 = 1
	if (this.shapeType == BRIDGE) {
		z += height
	}
	if (this.shapeType != SHIP_DESTROYED) {
		setScreenColor(r, g, b, 1)
	}
	gl.Begin(gl.LINE_LOOP)
	if (this.shapeType != BRIDGE) {
		this.createLoop(sz, z, false, true)
	} else {
		this.createSquareLoop(sz, z, false, true)
	}
	gl.End()
	if (this.shapeType != SHIP_SHADOW && this.shapeType != SHIP_DESTROYED &&
			this.shapeType != PLATFORM_DESTROYED && this.shapeType != TURRET_DESTROYED) {
		setScreenColor(r * 0.4, g * 0.4, b * 0.4, 1)
		gl.Begin(gl.TRIANGLE_FAN)
		this.createLoop(sz, z, true)
		gl.End()
	}
	switch (this.shapeType) {
	case SHIP, SHIP_ROUNDTAIL, SHIP_SHADOW, SHIP_DAMAGED, SHIP_DESTROYED:
		if (this.shapeType != SHIP_DESTROYED) {
			setScreenColor(r * 0.4, g * 0.4, b * 0.4)
		}
		for i := 0; i < 3; i++ {
			z -= height / 4
			sz -= 0.2
			gl.Begin(gl.LINE_LOOP)
			this.createLoop(sz, z)
			gl.End()
		}
		break
	case PLATFORM, PLATFORM_DAMAGED, PLATFORM_DESTROYED:
		setScreenColor(r * 0.4, g * 0.4, b * 0.4)
		for i := 0; i < 3; i++ {
			z -= height / 3
			for pp := range pillarPos {
				gl.Begin(gl.LINE_LOOP)
				this.createPillar(pp, size * 0.2, z)
				gl.End()
			}
		}
		break
	case BRIDGE, TURRET, TURRET_DAMAGED:
		setScreenColor(r * 0.6, g * 0.6, b * 0.6)
		z += height
		sz -= 0.33
		gl.Begin(gl.LINE_LOOP)
		if (this.shapeType == BRIDGE) {
			this.createSquareLoop(sz, z)
		} else {
			this.createSquareLoop(sz, z / 2, false, 3)
		}
		gl.End()
		setScreenColor(r * 0.25, g * 0.25, b * 0.25)
		gl.Begin(gl.TRIANGLE_FAN)
		if (this.shapeType == BRIDGE) {
			this.createSquareLoop(sz, z, true)
		} else {
			this.createSquareLoop(sz, z / 2, true, 3)
		}
		gl.End()
		break
	case TURRET_DESTROYED:
		break
	}
}

func (this *ComplexShape) createLoop(s float32, z float32, backToFirst bool /*= false*/, record bool /*= false*/) {
	var d float32 = 0
	var pn int
	firstPoint := true
	var fpx, fpy float32
	for i := 0; i < POINT_NUM; i++ {
		if (this.shapeType != SHIP && this.shapeType != SHIP_DESTROYED && this.shapeType != SHIP_DAMAGED &&
				i > POINT_NUM * 2 / 5 && i <= POINT_NUM * 3 / 5) {
			continue
		}
		if ((this.shapeType == TURRET || this.shapeType == TURRET_DAMAGED || this.shapeType == TURRET_DESTROYED) &&
				(i <= POINT_NUM / 5 || i > POINT_NUM * 4 / 5)) {
			continue
		}
		d = Pi32 * 2 * i / POINT_NUM
		cx = Sin32(d) * this.size * s * (1 - this.distRatio)
		cy = Cos32(d) * this.size * s
		var sx, sy float32
		if (i == POINT_NUM / 4 || i == POINT_NUM / 4 * 3) {
			sy = 0
		} else {
			sy = 1 / (1 + fabs32(tan32(d)))
		}
		sx = 1 - sy
		if (i >= POINT_NUM / 2) {
			sx *= -1
		}
		if (i >= POINT_NUM / 4 && i <= POINT_NUM / 4 * 3) {
			sy *= -1
		}
		sx *= this.size * s * (1 - this.distRatio)
		sy *= this.size * s
		px := cx * (1 - this.spinyRatio) + sx * this.spinyRatio
		py := cy * (1 - this.spinyRatio) + sy * this.spinyRatio
		gl.Vertex3(px, py, z)
		if (backToFirst && firstPoint) {
			fpx = px
			fpy = py
			firstPoint = false
		}
		if record {
			if (i == POINT_NUM / 8 || i == POINT_NUM / 8 * 3 ||
					i == POINT_NUM / 8 * 5 || i == POINT_NUM / 8 * 7) {
				this.pillarPos = append(this.pillarPos,Vector{px * 0.8, py * 0.8})
			}
			this.pointPos = append(this.pointPos,Vector{px, py})
			this.pointDeg = append(this.pointDeg,d)
		}
	}
	if (backToFirst) {
		gl.Vertex3(fpx, fpy, z)
	}
}

private void createSquareLoop(float32 s, float32 z, bool backToFirst = false, float32 yRatio = 1) {
	float32 d
	int pn
	if (backToFirst) {
		pn = 4
	} else {
		pn = 3
	}
	for i := 0; i <= pn; i++ {
		d = PI * 2 * i / 4 + PI / 4
		float32 px = sin(d) * size * s
		float32 py = cos(d) * size * s
		if (py > 0) {
			py *= yRatio
		}
		gl.Vertex3(px, py, z)
	}
}

private void createPillar(Vector p, float32 s, float32 z) {
	float32 d
	for i := 0; i < PILLAR_POINT_NUM; i++ {
		d = PI * 2 * i / PILLAR_POINT_NUM
		gl.Vertex3(sin(d) * s + p.x, cos(d) * s + p.y, z)
	}
}

public void addWake(WakePool wakes, Vector pos, float32 deg, float32 spd, float32 sr = 1) {
	float32 sp = spd
	if (sp > 0.1) {
		sp = 0.1
	}
	float32 sz = size
	if (sz > 10) {
		sz = 10
	}
	wakePos.x = pos.x + sin(deg + PI / 2 + 0.7) * size * 0.5 * sr
	wakePos.y = pos.y + cos(deg + PI / 2 + 0.7) * size * 0.5 * sr
	Wake w = wakes.getInstanceForced()
	w.set(wakePos, deg + PI - 0.2 + rand.nextSignedfloat32(0.1), sp, 40, sz * 32 * sr)
	wakePos.x = pos.x + sin(deg - PI / 2 - 0.7) * size * 0.5 * sr
	wakePos.y = pos.y + cos(deg - PI / 2 - 0.7) * size * 0.5 * sr
	w = wakes.getInstanceForced()
	w.set(wakePos, deg + PI + 0.2 + rand.nextSignedfloat32(0.1), sp, 40, sz * 32 * sr)
}

public Vector[] pointPos() {
	return _pointPos
}

public float32[] pointDeg() {
	return _pointDeg
}

public bool checkShipCollision(float32 x, float32 y, float32 deg, float32 sr = 1) {
	float32 cs = size * (1 - distRatio) * 1.1 * sr
	if (dist(x, y, 0, 0) < cs) {
		return true
	}
	float32 ofs = 0
	for {
		ofs += cs
		cs *= distRatio
		if (cs < 0.2) {
			return false
		}
		if (dist(x, y, sin(deg) * ofs, cos(deg) * ofs) < cs ||
				dist(x, y, -sin(deg) * ofs, -cos(deg) * ofs) < cs) {
			return true
		}
	}
}

private float32 dist(float32 x, float32 y, float32 px, float32 py) {
	float32 ax = fabs(x - px)
	float32 ay = fabs(y - py)
	if (ax > ay) {
		return ax + ay / 2
	} else {
		return ay + ax / 2
	}
}
}

public class TurretShape: ResizableShape {
 public:
  static enum TurretShapeType {
    NORMAL, DAMAGED, DESTROYED,
  }
 private:
  static ComplexShape[] shapes

  public static void init() {
    shapes = append(shapes,new CollidableComplexShape(1, 0, 0, ComplexShape.TURRET, 1, 0.8, 0.8))
    shapes = append(shapes,new ComplexShape(1, 0, 0, ComplexShape.TURRET_DAMAGED, 0.9, 0.9, 1))
    shapes = append(shapes,new ComplexShape(1, 0, 0, ComplexShape.TURRET_DESTROYED, 0.8, 0.33, 0.66))
  }

  public static void close() {
    foreach (ComplexShape s; shapes) {
      s.close()
		}
  }

  public this(int t) {
    shape = shapes[t]
  }
}

public class EnemyShape: ResizableShape {
 public:
  static enum EnemyShapeType {
    SMALL, SMALL_DAMAGED, SMALL_BRIDGE,
    MIDDLE, MIDDLE_DAMAGED, MIDDLE_DESTROYED, MIDDLE_BRIDGE,
    PLATFORM, PLATFORM_DAMAGED, PLATFORM_DESTROYED, PLATFORM_BRIDGE,
  }
  static const float32 MIDDLE_COLOR_R = 1, MIDDLE_COLOR_G = 0.6, MIDDLE_COLOR_B = 0.5
 private:
  static ComplexShape[] shapes

  public static void init() {
    shapes = append(shapes,new ComplexShape
      (1, 0.5, 0.1, ComplexShape.SHIP, 0.9, 0.7, 0.5))
    shapes = append(shapes,new ComplexShape
      (1, 0.5, 0.1, ComplexShape.SHIP_DAMAGED, 0.5, 0.5, 0.9))
    shapes = append(shapes,new CollidableComplexShape
      (0.66, 0, 0, ComplexShape.BRIDGE, 1, 0.2, 0.3))
    shapes = append(shapes,new ComplexShape
      (1, 0.7, 0.33, ComplexShape.SHIP, MIDDLE_COLOR_R, MIDDLE_COLOR_G, MIDDLE_COLOR_B))
    shapes = append(shapes, new ComplexShape
      (1, 0.7, 0.33, ComplexShape.SHIP_DAMAGED, 0.5, 0.5, 0.9))
    shapes = append(shapes, new ComplexShape
      (1, 0.7, 0.33, ComplexShape.SHIP_DESTROYED, 0, 0, 0))
    shapes = append(shapes, new CollidableComplexShape
      (0.66, 0, 0, ComplexShape.BRIDGE, 1, 0.2, 0.3))
    shapes = append(shapes, new ComplexShape
      (1, 0, 0, ComplexShape.PLATFORM, 1, 0.6, 0.7))
    shapes = append(shapes, new ComplexShape
      (1, 0, 0, ComplexShape.PLATFORM_DAMAGED, 0.5, 0.5, 0.9))
    shapes = append(shapes, new ComplexShape
      (1, 0, 0, ComplexShape.PLATFORM_DESTROYED, 1, 0.6, 0.7))
    shapes = append(shapes, new CollidableComplexShape
      (0.5, 0, 0, ComplexShape.BRIDGE, 1, 0.2, 0.3))
  }

  public static void close() {
    foreach (ComplexShape s; shapes) {
      s.close()
		}
  }

  public this(int t) {
    shape = shapes[t]
  }

  public void addWake(WakePool wakes, Vector pos, float32 deg, float32 sp) {
    (cast(ComplexShape) shape).addWake(wakes, pos, deg, sp, size)
  }

  public bool checkShipCollision(float32 x, float32 y, float32 deg) {
    return (cast(ComplexShape) shape).checkShipCollision(x, y, deg, size)
  }
}

public class BulletShape: ResizableShape {
 public:
  static enum BulletShapeType {
    NORMAL, SMALL, MOVING_TURRET, DESTRUCTIVE,
  }
 private:
  static SimpleShape[] shapes

  public static void init() {
    shapes = append(shapes, new NormalBulletShape)
    shapes = append(shapes, new SmallBulletShape)
    shapes = append(shapes, new MovingTurretBulletShape)
    shapes = append(shapes, new DestructiveBulletShape)
  }

  public static void close() {
    foreach (SimpleShape s; shapes) {
      s.close()
		}
  }

  public void set(int t) {
    shape = shapes[t]
  }
}

public class NormalBulletShape: SimpleShape {
  public override void createDisplayList() {
    gl.Disable(gl.BLEND)
    setScreenColor(1, 1, 0.3)
    gl.Begin(gl.LINE_STRIP)
    gl.Vertex3(0.2, -0.25, 0.2)
    gl.Vertex3(0, 0.33, 0)
    gl.Vertex3(-0.2, -0.25, -0.2)
    gl.End()
    gl.Begin(gl.LINE_STRIP)
    gl.Vertex3(-0.2, -0.25, 0.2)
    gl.Vertex3(0, 0.33, 0)
    gl.Vertex3(0.2, -0.25, -0.2)
    gl.End()
    gl.Enable(gl.BLEND)
    setScreenColor(0.5, 0.2, 0.1)
    gl.Begin(gl.TRIANGLE_FAN)
    gl.Vertex3(0, 0.33, 0)
    gl.Vertex3(0.2, -0.25, 0.2)
    gl.Vertex3(-0.2, -0.25, 0.2)
    gl.Vertex3(-0.2, -0.25, -0.2)
    gl.Vertex3(0.2, -0.25, -0.2)
    gl.Vertex3(0.2, -0.25, 0.2)
    gl.End()
  }
}

public class SmallBulletShape: SimpleShape {
  public override void createDisplayList() {
    gl.Disable(gl.BLEND)
    setScreenColor(0.6, 0.9, 0.3)
    gl.Begin(gl.LINE_STRIP)
    gl.Vertex3(0.25, -0.25, 0.25)
    gl.Vertex3(0, 0.33, 0)
    gl.Vertex3(-0.25, -0.25, -0.25)
    gl.End()
    gl.Begin(gl.LINE_STRIP)
    gl.Vertex3(-0.25, -0.25, 0.25)
    gl.Vertex3(0, 0.33, 0)
    gl.Vertex3(0.25, -0.25, -0.25)
    gl.End()
    gl.Enable(gl.BLEND)
    setScreenColor(0.2, 0.4, 0.1)
    gl.Begin(gl.TRIANGLE_FAN)
    gl.Vertex3(0, 0.33, 0)
    gl.Vertex3(0.25, -0.25, 0.25)
    gl.Vertex3(-0.25, -0.25, 0.25)
    gl.Vertex3(-0.25, -0.25, -0.25)
    gl.Vertex3(0.25, -0.25, -0.25)
    gl.Vertex3(0.25, -0.25, 0.25)
    gl.End()
  }
}

public class MovingTurretBulletShape: SimpleShape {
  public override void createDisplayList() {
    gl.Disable(gl.BLEND)
    setScreenColor(0.7, 0.5, 0.9)
    gl.Begin(gl.LINE_STRIP)
    gl.Vertex3(0.25, -0.25, 0.25)
    gl.Vertex3(0, 0.33, 0)
    gl.Vertex3(-0.25, -0.25, -0.25)
    gl.End()
    gl.Begin(gl.LINE_STRIP)
    gl.Vertex3(-0.25, -0.25, 0.25)
    gl.Vertex3(0, 0.33, 0)
    gl.Vertex3(0.25, -0.25, -0.25)
    gl.End()
    gl.Enable(gl.BLEND)
    setScreenColor(0.2, 0.2, 0.3)
    gl.Begin(gl.TRIANGLE_FAN)
    gl.Vertex3(0, 0.33, 0)
    gl.Vertex3(0.25, -0.25, 0.25)
    gl.Vertex3(-0.25, -0.25, 0.25)
    gl.Vertex3(-0.25, -0.25, -0.25)
    gl.Vertex3(0.25, -0.25, -0.25)
    gl.Vertex3(0.25, -0.25, 0.25)
    gl.End()
  }
}

public class DestructiveBulletShape: SimpleShape {
 private:
  Vector _collision

  public override void createDisplayList() {
    gl.Disable(gl.BLEND)
    setScreenColor(0.9, 0.9, 0.6)
    gl.Begin(gl.LINE_LOOP)
    gl.Vertex3(0.2, 0, 0)
    gl.Vertex3(0, 0.4, 0)
    gl.Vertex3(-0.2, 0, 0)
    gl.Vertex3(0, -0.4, 0)
    gl.End()
    gl.Enable(gl.BLEND)
    setScreenColor(0.7, 0.5, 0.4)
    gl.Begin(gl.TRIANGLE_FAN)
    gl.Vertex3(0.2, 0, 0)
    gl.Vertex3(0, 0.4, 0)
    gl.Vertex3(-0.2, 0, 0)
    gl.Vertex3(0, -0.4, 0)
    gl.End()
    _collision = new Vector(0.4, 0.4)
  }

  public Vector collision() {
    return _collision
  }
}

public class CrystalShape: SimpleShape {
  public override void createDisplayList() {
    setScreenColor(0.6, 1, 0.7)
    gl.Begin(gl.LINE_LOOP)
    gl.Vertex3(-0.2, 0.2, 0)
    gl.Vertex3(0.2, 0.2, 0)
    gl.Vertex3(0.2, -0.2, 0)
    gl.Vertex3(-0.2, -0.2, 0)
    gl.End()
  }
}

public class ShieldShape: SimpleShape {
  public override void createDisplayList() {
    setScreenColor(0.5, 0.5, 0.7)
    gl.Begin(gl.LINE_LOOP)
    float32 d = 0
		for i := 0; i < 8; i++ {
      gl.Vertex3(sin(d), cos(d), 0)
      d += PI / 4
    }
    gl.End()
    gl.Begin(gl.TRIANGLE_FAN)
    setScreenColor(0, 0, 0)
    gl.Vertex3(0, 0, 0)
    d = 0
    setScreenColor(0.3, 0.3, 0.5)
		for i := 0; i < 9; i++ {
      gl.Vertex3(sin(d), cos(d), 0)
      d += PI / 4
    }
    gl.End()
  }
}


/**
 * Interface for drawing a shape.
 */
public interface Shape {
  public void draw()
  public Vector collision()
  public bool checkCollision(float32 ax, float32 ay, Shape shape = null)
}

/* just a displaylist
   and a simple static collision, if collidable */
public template SimpleShape() {
  protected DisplayList displayList
  collision Vector

	func (ss *SimpleShape) checkCollision(ax float32, ay float32, shape Shape/* = null */) {
		return checkCollisionWithShapes(ax,ay,ss,shape)
	}

  public this() {
    displayList = new DisplayList(1)
    displayList.beginNewList()
    createDisplayList()
    displayList.endNewList()
  }

  protected abstract void createDisplayList()

  public Vector collision() {
    return collision
  }

  public void close() {
    displayList.close()
  }

  public void draw() {
    displayList.call(0)
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

func checkCollisionWithShapes(ax float32, ay float32, shape1 Shape, shape2 Shape) bool {
	if shape1 == nil {
		// this shape doesn't collide
		return false
	}
  float32 cx, cy
  if shape2 != nil {
    cx = shape1.collision().X() + shape2.collision().X()
    cy = shape1.collision().Y() + shape2.collision().Y()
  } else {
    cx = shape1.collision().X()
    cy = shape1.collision().Y()
  }
  return ax <= cx && ay <= cy
}

func (rd *ResizableShape) checkCollision(ax float32, ay float32, shape Shape/* = null */) {
	return checkCollisionWithShapes(ax,ay,rd,shape)
}

