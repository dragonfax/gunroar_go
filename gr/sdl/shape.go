package sdl

import (
	"github.com/dragonfax/gunroar/gr/vector"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

/**
 * Interface for drawing a shape.
 */
type Drawable interface {
	Draw()
}

/**
 * Interface and implmentation for a shape that has a collision.
 */
type Collidable interface {
	HasCollision
	CheckCollision(ax, ay float64, shape Collidable /* = null */) bool
}

type HasCollision interface {
	Collision() *vector.Vector // could be nil
}

// CollidableImpl will wrap the collidable "sub classes" to give them their CheckCollision() method.
type CollidableImpl struct {
	Child HasCollision
}

func NewCollidable(child HasCollision) *CollidableImpl {
	this := NewCollidableInternal(child)
	return &this
}

func NewCollidableInternal(child HasCollision) CollidableImpl {
	return CollidableImpl{Child: child}
}

func (this *CollidableImpl) CheckCollision(ax, ay float64, shape Collidable /* = null */) bool {
	if this.Child == nil {
		panic("collidable impl never got a chid with collision")
	}
	var cx, cy float64
	if shape != nil {
		cx = this.Child.Collision().X + shape.Collision().X
		cy = this.Child.Collision().Y + shape.Collision().Y
	} else {
		cx = this.Child.Collision().X
		cy = this.Child.Collision().Y
	}
	return ax <= cx && ay <= cy
}

type HasCreateDisplayList interface {
	CreateDisplayList()
}

/**
 * Drawable that has a single displaylist.
 */
type DrawableShape struct {
	displayList *DisplayList
}

func NewDrawableShape(child HasCreateDisplayList) *DrawableShape {
	this := NewDrawableShapeInternal(child)
	return &this
}
func NewDrawableShapeInternal(child HasCreateDisplayList) DrawableShape {
	this := DrawableShape{}
	this.displayList = NewDisplayList(1)
	this.displayList.BeginNewList()
	child.CreateDisplayList()
	this.displayList.EndNewList()
	return this
}

func (this *DrawableShape) Close() {
	this.displayList.Close()
}

func (this *DrawableShape) Draw() {
	this.displayList.Call(0)
}

/**
 * DrawableShape that has a collision.
 */

var _ Collidable = &CollidableDrawable{}

type CollidableDrawable struct {
	*CollidableImpl
	*DrawableShape

	_collision vector.Vector
}

type HasSetCollision interface {
	SetCollision() vector.Vector
}

func NewCollidableDrawable(collidable HasCollision, shape HasCreateDisplayList, setCollision HasSetCollision) CollidableDrawable {
	this := CollidableDrawable{
		CollidableImpl: NewCollidable(collidable),
		DrawableShape:  NewDrawableShape(shape),
	}
	this._collision = setCollision.SetCollision()
	return this
}

func (this *CollidableDrawable) Collision() *vector.Vector {
	v := this._collision
	return &v
}

/**
 * Drawable that can change a size.
 */

var _ Drawable = &ResizableDrawable{}
var _ Collidable = &ResizableDrawable{}

type ResizableDrawable struct {
	CollidableImpl

	_shape     Drawable
	_size      float64
	_collision vector.Vector
}

func NewResizableDrawableInternal() ResizableDrawable {
	this := ResizableDrawable{}
	this.CollidableImpl = NewCollidableInternal(&this)
	return this
}

func (this *ResizableDrawable) Draw() {
	gl.Scaled(this._size, this._size, this._size)
	this._shape.Draw()
}

func (this *ResizableDrawable) SetShape(v Drawable) Drawable {
	this._collision = vector.Vector{}
	this._shape = v
	return v
}

func (this *ResizableDrawable) Shape() Drawable {
	return this._shape
}

func (this *ResizableDrawable) SetSize(v float64) float64 {
	this._size = v
	return v
}

func (this *ResizableDrawable) Size() float64 {
	return this._size
}

func (this *ResizableDrawable) Collision() *vector.Vector {
	cd := this._shape.(Collidable)
	if cd != nil {
		this._collision.X = cd.Collision().X * this._size
		this._collision.Y = cd.Collision().Y * this._size
		v := this._collision
		return &v
	}
	return nil
}
