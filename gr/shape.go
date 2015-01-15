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
	shapeType                   ShapeType
	r, g, b                     float32
	pillarPos                   []Vector
	pointPos                    []Vector
	pointDeg                    []float32

	collidable bool
}

func (this *ComplexShape) InitComplexShape(size float32, distRatio float32, spinyRatio float32, shapeType int, r float32, g float32, b float32, collidable bool /* = false */) {

	this.size = size
	this.distRatio = distRatio
	this.spinyRatio = spinyRatio
	this.shapeType = shapeType
	this.r = r
	this.g = g
	this.b = b
	this.collidable = collidable
	if collidable {
		this.collision = Vector{size / 2, size / 2}
	} else {
		this.collision = nil
	}
	this.InitSimpleShape()
	this.createDisplayList()
}

func (this *ComplexShape) createDisplayList() {
	height := this.size * 0.5
	var z float32 = 0
	var sz float32 = 1
	if this.shapeType == BRIDGE {
		z += height
	}
	if this.shapeType != SHIP_DESTROYED {
		setScreenColor(r, g, b, 1)
	}
	gl.Begin(gl.LINE_LOOP)
	if this.shapeType != BRIDGE {
		this.createLoop(sz, z, false, true)
	} else {
		this.createSquareLoop(sz, z, false, true)
	}
	gl.End()
	if this.shapeType != SHIP_SHADOW && this.shapeType != SHIP_DESTROYED &&
		this.shapeType != PLATFORM_DESTROYED && this.shapeType != TURRET_DESTROYED {
		setScreenColor(r*0.4, g*0.4, b*0.4, 1)
		gl.Begin(gl.TRIANGLE_FAN)
		this.createLoop(sz, z, true)
		gl.End()
	}
	switch this.shapeType {
	case SHIP, SHIP_ROUNDTAIL, SHIP_SHADOW, SHIP_DAMAGED, SHIP_DESTROYED:
		if this.shapeType != SHIP_DESTROYED {
			setScreenColor(r*0.4, g*0.4, b*0.4)
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
		setScreenColor(r*0.4, g*0.4, b*0.4)
		for i := 0; i < 3; i++ {
			z -= height / 3
			for pp := range pillarPos {
				gl.Begin(gl.LINE_LOOP)
				this.createPillar(pp, size*0.2, z)
				gl.End()
			}
		}
		break
	case BRIDGE, TURRET, TURRET_DAMAGED:
		setScreenColor(r*0.6, g*0.6, b*0.6)
		z += height
		sz -= 0.33
		gl.Begin(gl.LINE_LOOP)
		if this.shapeType == BRIDGE {
			this.createSquareLoop(sz, z)
		} else {
			this.createSquareLoop(sz, z/2, false, 3)
		}
		gl.End()
		setScreenColor(r*0.25, g*0.25, b*0.25)
		gl.Begin(gl.TRIANGLE_FAN)
		if this.shapeType == BRIDGE {
			this.createSquareLoop(sz, z, true)
		} else {
			this.createSquareLoop(sz, z/2, true, 3)
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
		if this.shapeType != SHIP && this.shapeType != SHIP_DESTROYED && this.shapeType != SHIP_DAMAGED &&
			i > POINT_NUM*2/5 && i <= POINT_NUM*3/5 {
			continue
		}
		if (this.shapeType == TURRET || this.shapeType == TURRET_DAMAGED || this.shapeType == TURRET_DESTROYED) &&
			(i <= POINT_NUM/5 || i > POINT_NUM*4/5) {
			continue
		}
		d = Pi32 * 2 * i / POINT_NUM
		cx = Sin32(d) * this.size * s * (1 - this.distRatio)
		cy = Cos32(d) * this.size * s
		var sx, sy float32
		if i == POINT_NUM/4 || i == POINT_NUM/4*3 {
			sy = 0
		} else {
			sy = 1 / (1 + fabs32(tan32(d)))
		}
		sx = 1 - sy
		if i >= POINT_NUM/2 {
			sx *= -1
		}
		if i >= POINT_NUM/4 && i <= POINT_NUM/4*3 {
			sy *= -1
		}
		sx *= this.size * s * (1 - this.distRatio)
		sy *= this.size * s
		px := cx*(1-this.spinyRatio) + sx*this.spinyRatio
		py := cy*(1-this.spinyRatio) + sy*this.spinyRatio
		gl.Vertex3(px, py, z)
		if backToFirst && firstPoint {
			fpx = px
			fpy = py
			firstPoint = false
		}
		if record {
			if i == POINT_NUM/8 || i == POINT_NUM/8*3 ||
				i == POINT_NUM/8*5 || i == POINT_NUM/8*7 {
				this.pillarPos = append(this.pillarPos, Vector{px * 0.8, py * 0.8})
			}
			this.pointPos = append(this.pointPos, Vector{px, py})
			this.pointDeg = append(this.pointDeg, d)
		}
	}
	if backToFirst {
		gl.Vertex3(fpx, fpy, z)
	}
}

func (this *ComplexShape) createSquareLoop(s float32, z float32, backToFirst bool /*= false*/, yRatio float32 /*= 1*/) {
	var d float32
	var pn int
	if backToFirst {
		pn = 4
	} else {
		pn = 3
	}
	for i := 0; i <= pn; i++ {
		d := Pi32*2*i/4 + Pi32/4
		px := Sin32(d) * this.size * s
		py := Cos32(d) * this.size * s
		if py > 0 {
			py *= yRatio
		}
		gl.Vertex3(px, py, z)
	}
}

func (this *ComplexShape) createPillar(p Vector, s float32, z float32) {
	var d float32
	for i := 0; i < PILLAR_POINT_NUM; i++ {
		d := Pi32 * 2 * i / PILLAR_POINT_NUM
		gl.Vertex3(Sin32(d)*s+p.x, Cos32(d)*s+p.y, z)
	}
}

func (this *ComplexShape) addWake(pos Vector, deg float32, spd float32, sr float32 /*= 1*/) {
	sp := spd
	if sp > 0.1 {
		sp = 0.1
	}
	sz := size
	if sz > 10 {
		sz = 10
	}
	wakePos.x = pos.x + Sin32(deg+Pi32/2+0.7)*size*0.5*sr
	wakePos.y = pos.y + Cos32(deg+Pi32/2+0.7)*size*0.5*sr
	w := NewWake()
	w.set(wakePos, deg+Pi32-0.2+nextSignedFloat(0.1), sp, 40, sz*32*sr)
	wakePos.x = pos.x + Sin32(deg-Pi32/2-0.7)*size*0.5*sr
	wakePos.y = pos.y + Cos32(deg-Pi32/2-0.7)*size*0.5*sr
	w = NewWake()
	w.set(wakePos, deg+Pi32+0.2+nextSignedFloat(0.1), sp, 40, sz*32*sr)
}

func (this *ComplexShape) checkShipCollision(x float32, y float32, deg float32, sr float32 /*= 1*/) bool {
	cs := this.size * (1 - this.distRatio) * 1.1 * sr
	if this.dist(x, y, 0, 0) < cs {
		return true
	}
	var ofs float32 = 0
	for {
		ofs += cs
		cs *= this.distRatio
		if cs < 0.2 {
			return false
		}
		if this.dist(x, y, Sin32(deg)*ofs, Cos32(deg)*ofs) < cs ||
			this.dist(x, y, -Sin32(deg)*ofs, -Cos32(deg)*ofs) < cs {
			return true
		}
	}
}

func (this *ComplexShape) dist(x float32, y float32, px float32, py float32) float32 {
	ax := fabs32(x - px)
	ay := fabs32(y - py)
	if ax > ay {
		return ax + ay/2
	} else {
		return ay + ax/2
	}
}

type TurretShapeType int

const (
	TurretShapeTypeNORMAL TurretShapeType = iota
	TurretShapeTypeDAMAGED
	TurretShapeTypeDESTROYED
)

var turretShapes []ComplexShape

func InitTurretShapes() {
	turretShapes = append(turretShapes, NewCollidableComplexShape(1, 0, 0, ComplexShape.TURRET, 1, 0.8, 0.8))
	turretShapes = append(turretShapes, NewComplexShape(1, 0, 0, ComplexShape.TURRET_DAMAGED, 0.9, 0.9, 1))
	turretShapes = append(turretShapes, NewComplexShape(1, 0, 0, ComplexShape.TURRET_DESTROYED, 0.8, 0.33, 0.66))
}

type TurretShape struct {
	*ResizableShape
}

func closeTurretShapes() {
	for s := range turretShapes {
		s.close()
	}
}

func NewTurretShape(t TurretShapeType) TurretShape {
	turretShape = TurretShape{}
	turretShape.shape = turretShapes[t]
}

type EnemyShapeType int

const (
	EnemyShapeTypeSMALL EnemyShapeType = iota
	EnemyShapeTypeSMALL_DAMAGED
	EnemyShapeTypeSMALL_BRIDGE
	EnemyShapeTypeMIDDLE
	EnemyShapeTypeMIDDLE_DAMAGED
	EnemyShapeTypeMIDDLE_DESTROYED
	EnemyShapeTypeMIDDLE_BRIDGE
	EnemyShapeTypePLATFORM
	EnemyShapeTypePLATFORM_DAMAGED
	EnemyShapeTypePLATFORM_DESTROYED
	EnemyShapeTypePLATFORM_BRIDGE
)

const MIDDLE_COLOR_R = 1
const MIDDLE_COLOR_G = 0.6
const MIDDLE_COLOR_B = 0.5

var enemyShapes []ComplexShape

type EnemyShape struct {
	*ResizableShape
}

func InitEnemyShapes() {
	enemyShapes = append(enemyShapes, NewComplexShape(1, 0.5, 0.1, ComplexShape.SHIP, 0.9, 0.7, 0.5))
	enemyShapes = append(enemyShapes, NewComplexShape(1, 0.5, 0.1, ComplexShape.SHIP_DAMAGED, 0.5, 0.5, 0.9))
	enemyShapes = append(enemyShapes, NewCollidableComplexShape(0.66, 0, 0, ComplexShape.BRIDGE, 1, 0.2, 0.3))
	enemyShapes = append(enemyShapes, NewComplexShape(1, 0.7, 0.33, ComplexShape.SHIP, MIDDLE_COLOR_R, MIDDLE_COLOR_G, MIDDLE_COLOR_B))
	enemyShapes = append(enemyShapes, NewComplexShape(1, 0.7, 0.33, ComplexShape.SHIP_DAMAGED, 0.5, 0.5, 0.9))
	enemyShapes = append(enemyShapes, NewComplexShape(1, 0.7, 0.33, ComplexShape.SHIP_DESTROYED, 0, 0, 0))
	enemyShapes = append(enemyShapes, NewCollidableComplexShape(0.66, 0, 0, ComplexShape.BRIDGE, 1, 0.2, 0.3))
	enemyShapes = append(enemyShapes, NewComplexShape(1, 0, 0, ComplexShape.PLATFORM, 1, 0.6, 0.7))
	enemyShapes = append(enemyShapes, NewComplexShape(1, 0, 0, ComplexShape.PLATFORM_DAMAGED, 0.5, 0.5, 0.9))
	enemyShapes = append(enemyShapes, NewComplexShape(1, 0, 0, ComplexShape.PLATFORM_DESTROYED, 1, 0.6, 0.7))
	enemyShapes = append(enemyShapes, NewCollidableComplexShape(0.5, 0, 0, ComplexShape.BRIDGE, 1, 0.2, 0.3))
}

func closeEnemyShapes() {
	for s := range enemyShapes {
		s.close()
	}
}

func NewEnemyShape(t EnemyShapeType) EnemyShape {
	e = EnemyShape{}
	e.shape = enemyShapes[t]
	return e
}

func (this *EnemyShape) addWake(pos Vector, deg float32, sp float32) {
	cs, ok := this.shape.(ComplexShape)
	if ok {
		cs.addWake(wakes, pos, deg, sp, this.size)
	}
}

func (this *EnemyShape) checkShipCollision(x float32, y float32, deg float32) bool {
	cs, ok := this.shape.(ComplexShape)
	if ok {
		return cs.checkShipCollision(x, y, deg, this.size)
	} else {
		panic("enemy shape wasn't a complex shape")
	}
}

type BulletShapeType int

const (
	BulletShapeTypeNORMAL BulletShapeType = iota
	BulletShapeTypeSMALL
	BulletShapeTypeMOVING_TURRET
	BulletShapeTypeDESTRUCTIVE
)

var bulletShapes []SimpleShape

type BulletShape struct {
	*ResizableShape
}

func InitBulletShapes() {
	bulletShapes = append(bulletShapes, NewNormalBulletShape())
	bulletShapes = append(bulletShapes, NewSmallBulletShape())
	bulletShapes = append(bulletShapes, NewMovingTurretBulletShape())
	bulletShapes = append(bulletShapes, NewDestructiveBulletShape())
}

func closeBulletShapes() {
	for s := range bulletShapes {
		s.close()
	}
}

func NewBulletShape(t BulletShapeType) BulletShape {
	b = BulletShape{}
	b.shape = bulletShapes[t]
	return b
}

type NormalBulletShape struct {
	*SimpleShape
}

func NewNormalBulletShape() NormalBulletShape {
	nbs := NormalBulletShape{}
	nbs.startDisplayList()
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
	nbs.endDisplayList
	return nbs
}

type SmallBulletShape struct {
	*SimpleShape
}

func NewSmallBulletShape() SmallBulletShape {
	sbs := SmallBulletShape{}
	sbs.startDisplayList()
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
	sbs.endDisplayList()
	return sbs
}

type MovingTurretBulletShape struct {
	*SimpleShape
}

func NewMovingTurretBulletShape() MovingTurretBulletShape {
	mtbs := MovingTurretBulletShape{}
	mtbs.startDisplayList()
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
	mtbs.endDisplayList()
	return mtbs
}

type DestructiveBulletShape struct {
	*SimpleShape

	collision Vector
}

func NewDestructiveBulletShape() DestructiveBulletShape {
	dbs := DestructiveBulletShape{}
	dbs.startDisplayList()
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
	dbs.endDisplayList()
	dbs.collision = Vector{0.4, 0.4}
	return dbs
}

type CrystalShape struct {
	*SimpleShape
}

func NewCrystalShape() CrystalShape {
	cs := CrystalShape{}
	cs.startDisplayList()
	setScreenColor(0.6, 1, 0.7)
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex3(-0.2, 0.2, 0)
	gl.Vertex3(0.2, 0.2, 0)
	gl.Vertex3(0.2, -0.2, 0)
	gl.Vertex3(-0.2, -0.2, 0)
	gl.End()
	cs.endDisplayList()
	return cs
}

type ShieldShape struct {
	*SimpleShape
}

func NewShieldShape() ShieldShape {
	ss := ShieldShape{}
	ss.startDisplayList()
	setScreenColor(0.5, 0.5, 0.7)
	gl.Begin(gl.LINE_LOOP)
	var d float32 = 0
	for i := 0; i < 8; i++ {
		gl.Vertex3f(Sin32(d), Cos32(d), 0)
		d += Pi32 / 4
	}
	gl.End()
	gl.Begin(gl.TRIANGLE_FAN)
	setScreenColor(0, 0, 0)
	gl.Vertex3f(0, 0, 0)
	d = 0
	setScreenColor(0.3, 0.3, 0.5)
	for i := 0; i < 9; i++ {
		gl.Vertex3f(Sin32(d), Cos32(d), 0)
		d += Pi32 / 4
	}
	gl.End()
	ss.endDisplayList()
}

/**
 * Interface for drawing a shape.
 */
type Shape interface {
	draw()
	collision() Vector
	checkCollision(ax float32, ay float32, shape Shape /*= null */) bool
}

/* just a displaylist
   and a simple static collision, if collidable */
type SimpleShape struct {
	displayList DisplayList
	collision   Vector
}

func (ss *SimpleShape) checkCollision(ax float32, ay float32, shape Shape /* = null */) {
	return checkCollisionWithShapes(ax, ay, ss, shape)
}

func (ss *SimpleShape) startDisplayList() {
	ss.displayList = NewDisplayList(1)
	ss.displayList.beginNewList()
}

func (ss *SimpleShape) endDisplayList() {
	ss.displayList.endNewList()
}

func (ss *SimpleShape) collision() Vector {
	return ss.collision
}

func (ss *SimpleShape) close() {
	ss.displayList.close()
}

func (ss *SimpleShape) draw() {
	ss.displayList.call(0)
}

/*
 * a Shape that can change a size.
 *
 * proxies a Simple or Complex shape
 */
type ResizableShape struct {
	shape            Shape
	size             float32
	resizedCollision Vector
}

func (rd *ResizableShape) draw() {
	gl.Scalef(rs.size, rs.size, rs.size)
	rs.shape.Draw()
}

func (rd *ResizableShape) collision() *Vector {
	rs.resizedCollision = NewVector(cd.collision().x*rs.size, cd.collision().y*rs.size)
	return rs.resizedCollision
}

func checkCollisionWithShapes(ax float32, ay float32, shape1 Shape, shape2 Shape) bool {
	if shape1 == nil {
		// this shape doesn't collide
		return false
	}
	var cx, cy float32
	if shape2 != nil {
		cx = shape1.collision().x + shape2.collision().x
		cy = shape1.collision().y + shape2.collision().y
	} else {
		cx = shape1.collision().x
		cy = shape1.collision().y
	}
	return ax <= cx && ay <= cy
}

func (rd *ResizableShape) checkCollision(ax float32, ay float32, shape Shape /* = null */) {
	return checkCollisionWithShapes(ax, ay, rd, shape)
}
