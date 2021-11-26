package main

import (
	"math"

	"github.com/dragonfax/gunroar/gr/sdl"
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

type TurretShapeType int

const (
	TurretNORMAL TurretShapeType = iota
	TurretDAMAGED
	TurretDESTROYED
)

var turretShapes = []ShapeI{
	NewCollidableBaseShape(1, 0, 0, TURRET, 1, 0.8, 0.8),
	NewBaseShape(1, 0, 0, TURRET_DAMAGED, 0.9, 0.9, 1),
	NewBaseShape(1, 0, 0, TURRET_DESTROYED, 0.8, 0.33, 0.66),
}

type TurretShape struct {
	sdl.ResizableDrawable
}

func NewTurretShape(t TurretShapeType) *TurretShape {
	this := &TurretShape{sdl.NewResizableDrawableInternal()}
	this.SetShape(turretShapes[int(t)])
	return this
}

const (
	MIDDLE_COLOR_R = 1
	MIDDLE_COLOR_G = 0.6
	MIDDLE_COLOR_B = 0.5
)

type EnemyShapeType int

const (
	EnemySMALL EnemyShapeType = iota
	EnemySMALL_DAMAGED
	EnemySMALL_BRIDGE
	EnemyMIDDLE
	EnemyMIDDLE_DAMAGED
	EnemyMIDDLE_DESTROYED
	EnemyMIDDLE_BRIDGE
	EnemyPLATFORM
	EnemyPLATFORM_DAMAGED
	EnemyPLATFORM_DESTROYED
	EnemyPLATFORM_BRIDGE
)

var enemyShapes = []ShapeI{
	NewBaseShape(1, 0.5, 0.1, SHIP, 0.9, 0.7, 0.5),
	NewBaseShape(1, 0.5, 0.1, SHIP_DAMAGED, 0.5, 0.5, 0.9),
	NewCollidableBaseShape(0.66, 0, 0, BRIDGE, 1, 0.2, 0.3),
	NewBaseShape(1, 0.7, 0.33, SHIP, MIDDLE_COLOR_R, MIDDLE_COLOR_G, MIDDLE_COLOR_B),
	NewBaseShape(1, 0.7, 0.33, SHIP_DAMAGED, 0.5, 0.5, 0.9),
	NewBaseShape(1, 0.7, 0.33, SHIP_DESTROYED, 0, 0, 0),
	NewCollidableBaseShape(0.66, 0, 0, BRIDGE, 1, 0.2, 0.3),
	NewBaseShape(1, 0, 0, PLATFORM, 1, 0.6, 0.7),
	NewBaseShape(1, 0, 0, PLATFORM_DAMAGED, 0.5, 0.5, 0.9),
	NewBaseShape(1, 0, 0, PLATFORM_DESTROYED, 1, 0.6, 0.7),
	NewCollidableBaseShape(0.5, 0, 0, BRIDGE, 1, 0.2, 0.3),
}

type EnemyShape struct {
	sdl.ResizableDrawable
}

func NewEnemyShape(t int) *EnemyShape {
	this := &EnemyShape{sdl.NewResizableDrawableInternal()}
	this.SetShape(enemyShapes[t])
	return this
}

func (this *EnemyShape) addWake(wakes WakePool, pos vector.Vector, deg, sp float64) {
	this.Shape().(*BaseShape).addWake(wakes, pos, deg, sp, this.Size())
}

func (this *EnemyShape) checkShipCollision(x, y, deg float64) bool {
	return this.Shape().(*BaseShape).checkShipCollision(x, y, deg, this.Size())
}

type BulletShapeType int

const (
	BulletNORMAL BulletShapeType = iota
	BulletSMALL
	BulletMOVING_TURRET
	BulletDESTRUCTIVE
)

var bulletShapes = []ShapeI{
	NewNormalBulletShape(),
	NewSmallBulletShape(),
	NewMovingTurretBulletShape(),
	NewDestructiveBulletShape(),
}

type BulletShape struct {
	sdl.ResizableDrawable
}

func (this *BulletShape) Set(t int) {
	this.SetShape(bulletShapes[t])
}

type NormalBulletShape struct {
	sdl.DrawableShape
}

func NewNormalBulletShape() *NormalBulletShape {
	this := &NormalBulletShape{}
	this.DrawableShape = sdl.NewDrawableShapeInternal(this)
	return this
}

func (this *NormalBulletShape) CreateDisplayList() {
	gl.Disable(gl.BLEND)
	sdl.SetColor(1, 1, 0.3, 1)
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex3f(0.2, -0.25, 0.2)
	gl.Vertex3f(0, 0.33, 0)
	gl.Vertex3f(-0.2, -0.25, -0.2)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex3f(-0.2, -0.25, 0.2)
	gl.Vertex3f(0, 0.33, 0)
	gl.Vertex3f(0.2, -0.25, -0.2)
	gl.End()
	gl.Enable(gl.BLEND)
	sdl.SetColor(0.5, 0.2, 0.1, 1)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex3f(0, 0.33, 0)
	gl.Vertex3f(0.2, -0.25, 0.2)
	gl.Vertex3f(-0.2, -0.25, 0.2)
	gl.Vertex3f(-0.2, -0.25, -0.2)
	gl.Vertex3f(0.2, -0.25, -0.2)
	gl.Vertex3f(0.2, -0.25, 0.2)
	gl.End()
}

type SmallBulletShape struct {
	sdl.DrawableShape
}

func NewSmallBulletShape() *SmallBulletShape {
	this := &SmallBulletShape{}
	this.DrawableShape = sdl.NewDrawableShapeInternal(this)
	return this
}

func (this *SmallBulletShape) CreateDisplayList() {
	gl.Disable(gl.BLEND)
	sdl.SetColor(0.6, 0.9, 0.3, 1)
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex3f(0.25, -0.25, 0.25)
	gl.Vertex3f(0, 0.33, 0)
	gl.Vertex3f(-0.25, -0.25, -0.25)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex3f(-0.25, -0.25, 0.25)
	gl.Vertex3f(0, 0.33, 0)
	gl.Vertex3f(0.25, -0.25, -0.25)
	gl.End()
	gl.Enable(gl.BLEND)
	sdl.SetColor(0.2, 0.4, 0.1, 1)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex3f(0, 0.33, 0)
	gl.Vertex3f(0.25, -0.25, 0.25)
	gl.Vertex3f(-0.25, -0.25, 0.25)
	gl.Vertex3f(-0.25, -0.25, -0.25)
	gl.Vertex3f(0.25, -0.25, -0.25)
	gl.Vertex3f(0.25, -0.25, 0.25)
	gl.End()
}

type MovingTurretBulletShape struct {
	sdl.DrawableShape
}

func NewMovingTurretBulletShape() *MovingTurretBulletShape {
	this := &MovingTurretBulletShape{}
	this.DrawableShape = sdl.NewDrawableShapeInternal(this)
	return this
}

func (this *MovingTurretBulletShape) CreateDisplayList() {
	gl.Disable(gl.BLEND)
	sdl.SetColor(0.7, 0.5, 0.9, 1)
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex3f(0.25, -0.25, 0.25)
	gl.Vertex3f(0, 0.33, 0)
	gl.Vertex3f(-0.25, -0.25, -0.25)
	gl.End()
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex3f(-0.25, -0.25, 0.25)
	gl.Vertex3f(0, 0.33, 0)
	gl.Vertex3f(0.25, -0.25, -0.25)
	gl.End()
	gl.Enable(gl.BLEND)
	sdl.SetColor(0.2, 0.2, 0.3, 1)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex3f(0, 0.33, 0)
	gl.Vertex3f(0.25, -0.25, 0.25)
	gl.Vertex3f(-0.25, -0.25, 0.25)
	gl.Vertex3f(-0.25, -0.25, -0.25)
	gl.Vertex3f(0.25, -0.25, -0.25)
	gl.Vertex3f(0.25, -0.25, 0.25)
	gl.End()
}

var _ sdl.Collidable = &DestructiveBulletShape{}

type DestructiveBulletShape struct {
	sdl.DrawableShape
	sdl.CollidableImpl

	_collision vector.Vector
}

func NewDestructiveBulletShape() *DestructiveBulletShape {
	this := &DestructiveBulletShape{}
	this.DrawableShape = sdl.NewDrawableShapeInternal(this)
	this.CollidableImpl = sdl.NewCollidableInternal(this)
	return this
}

func (this *DestructiveBulletShape) CreateDisplayList() {
	gl.Disable(gl.BLEND)
	sdl.SetColor(0.9, 0.9, 0.6, 1)
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex3f(0.2, 0, 0)
	gl.Vertex3f(0, 0.4, 0)
	gl.Vertex3f(-0.2, 0, 0)
	gl.Vertex3f(0, -0.4, 0)
	gl.End()
	gl.Enable(gl.BLEND)
	sdl.SetColor(0.7, 0.5, 0.4, 1)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex3f(0.2, 0, 0)
	gl.Vertex3f(0, 0.4, 0)
	gl.Vertex3f(-0.2, 0, 0)
	gl.Vertex3f(0, -0.4, 0)
	gl.End()
	this._collision = vector.Vector{0.4, 0.4}
}

func (this *DestructiveBulletShape) Collision() *vector.Vector {
	return &this._collision
}

type CrystalShape struct {
	sdl.DrawableShape
}

func NewDrawableShape() *CrystalShape {
	this := &CrystalShape{}
	this.DrawableShape = sdl.NewDrawableShapeInternal(this)
	return this
}

func (this *CrystalShape) CreateDisplayList() {
	sdl.SetColor(0.6, 1, 0.7, 1)
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex3f(-0.2, 0.2, 0)
	gl.Vertex3f(0.2, 0.2, 0)
	gl.Vertex3f(0.2, -0.2, 0)
	gl.Vertex3f(-0.2, -0.2, 0)
	gl.End()
}

type ShieldShape struct {
	sdl.DrawableShape
}

func NewShieldShape() *ShieldShape {
	this := &ShieldShape{}
	this.DrawableShape = sdl.NewDrawableShapeInternal(this)
	return this
}

func (this *ShieldShape) CreateDisplayList() {
	sdl.SetColor(0.5, 0.5, 0.7, 1)
	gl.Begin(gl.LINE_LOOP)
	var d float64
	for i := 0; i < 8; i++ {
		gl.Vertex3f(math.Sin(d), math.Cos(d), 0)
		d += math.Pi / 4
	}
	gl.End()
	gl.Begin(gl.TRIANGLE_FAN)
	sdl.SetColor(0, 0, 0, 1)
	gl.Vertex3f(0, 0, 0)
	d = 0
	sdl.SetColor(0.3, 0.3, 0.5, 1)
	for i := 0; i < 9; i++ {
		gl.Vertex3f(math.Sin(d), math.Cos(d), 0)
		d += math.Pi / 4
	}
	gl.End()
}
