/*
 * $Id: shape.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package sdl

type Drawable interface {
	Draw()
}

type Collidable interface {
	Collision() *Vector
	CheckCollision(ax float32, ay float32, shape *Collidable /* = null */) bool
}

type CollidableImpl struct {
}

func (c *CollidableImpl) CheckCollision(ax float32, ay float32, shape *Collidable /* = nil */) bool {
	var cx, cy float32
	if shape != nil {
		cx = collision.x + shape.collision.x
		cy = collision.y + shape.collision.y
	} else {
		cx = collision.x
		cy = collision.y
	}
	if ax <= cx && ay <= cy {
		return true
	} else {
		return false
	}
}

/**
 * Drawable that has a single displaylist.
 */
type DrawableShape struct {
	displayList DisplayList
}

func NewDrawableShape() *DrawableShape {
	this := New(DrawableShape)
	this.displayList = NewDisplayList(1)
	this.displayList.beginNewList()
	this.createDisplayList()
	this.displayList.endNewList()
	return this
}

func (ds *DrawableShape) Close() {
	ds.displayList.Close()
}

func (ds *DrawableShape) Draw() {
	ds.displayList.Call(0)
}

/**
 * DrawableShape that has a collision.
 */
type CollidableDrawable struct {
	DrawableShape
	CollidableImpl

	collision Vector
}

func NewCollidableDrawable() CollidableDrawable {
	this := new(CollidableDrawable)
	this.setCollision()
}

/**
 * Drawable that can change a size.
 */
type ResizableDrawable struct {
	CollidableImpl

	Shape     *Drawable
	Size      float32
	collision *Vector
}

func (rd *ResizableDrawable) Draw() {
	gl.Scalef(rd.size, rd.size, rd.size)
	shape.Draw()
}

func (rd *ResizableDrawable) SetShape(v Drawable) *Drawable {
	rd.collision = new(Vector)
	rd.shape = v
	return rd.shape
}

func (rd *ResizableDrawable) Collision() *Vector {
	cd := Collidable(rd.shape)
	if cd != nil {
		rd.collision.x = cd.collision.x * rd.Size
		rd.collision.y = cd.collision.y * rd.Size
		return rd.collision
	} else {
		return nil
	}
}
