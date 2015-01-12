/*
 * $Id: vector.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2004 Kenta Cho. Some rights reserved.
 */
package gr

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Vector struct {
	*mgl32.Vec2
}

func NewVector(x float32, y float32) Vector {
	return Vector{&mgl32.Vec2{x, y}}
}

func (this Vector) DotV(v Vector) float32 {
	return this.Dot(mgl32.Vec2{v.X(), v.Y()})
}

func (this Vector) getElement(v Vector) Vector {
	var rsl Vector
	ll := v.DotV(v)
	if ll != 0 {
		mag := this.DotV(v)
		x := mag * v.X() / ll
		y := mag * v.Y() / ll
		rsl = NewVector(x, y)
	}
	return rsl
}

func (this Vector) checkSide(pos1 Vector, pos2 Vector) float32 {
	xo := pos2.X() - pos1.X()
	yo := pos2.Y() - pos1.Y()
	if xo == 0 {
		if yo == 0 {
			return 0
		}
		if yo > 0 {
			return this.X() - pos1.X()
		} else {
			return pos1.X() - this.X()
		}
	} else if yo == 0 {
		if xo > 0 {
			return pos1.Y() - this.Y()
		} else {
			return this.Y() - pos1.Y()
		}
	} else {
		if xo*yo > 0 {
			return (this.X()-pos1.X())/xo - (this.Y()-pos1.Y())/yo
		} else {
			return -(this.X()-pos1.X())/xo + (this.Y()-pos1.Y())/yo
		}
	}
}

func (this Vector) checkSide3(pos1 Vector, pos2 Vector, ofs Vector) float32 {
	xo := pos2.X() - pos1.X()
	yo := pos2.Y() - pos1.Y()
	mx := this.X() + ofs.X()
	my := this.Y() + ofs.Y()
	if xo == 0 {
		if yo == 0 {
			return 0
		}
		if yo > 0 {
			return mx - pos1.X()
		} else {
			return pos1.X() - mx
		}
	} else if yo == 0 {
		if xo > 0 {
			return pos1.Y() - my
		} else {
			return my - pos1.Y()
		}
	} else {
		if xo*yo > 0 {
			return (mx-pos1.X())/xo - (my-pos1.Y())/yo
		} else {
			return -(mx-pos1.X())/xo + (my-pos1.Y())/yo
		}
	}
}

func (this Vector) checkCross(p Vector, p1 Vector, p2 Vector, width float32) bool {
	var a1x, a1y, a2x, a2y float32
	if this.X() < p.X() {
		a1x = this.X() - width
		a2x = p.X() + width
	} else {
		a1x = p.X() - width
		a2x = this.X() + width
	}
	if this.Y() < p.Y() {
		a1y = this.Y() - width
		a2y = p.Y() + width
	} else {
		a1y = p.Y() - width
		a2y = this.Y() + width
	}
	var b1x, b1y, b2x, b2y float32
	if p2.Y() < p1.Y() {
		b1y = p2.Y() - width
		b2y = p1.Y() + width
	} else {
		b1y = p1.Y() - width
		b2y = p2.Y() + width
	}
	if a2y >= b1y && b2y >= a1y {
		if p2.X() < p1.X() {
			b1x = p2.X() - width
			b2x = p1.X() + width
		} else {
			b1x = p1.X() - width
			b2x = p2.X() + width
		}
		if a2x >= b1x && b2x >= a1x {
			a := this.Y() - p.Y()
			b := p.X() - this.X()
			c := p.X()*this.Y() - p.Y()*this.X()
			d := p2.Y() - p1.Y()
			e := p1.X() - p2.X()
			f := p1.X()*p2.Y() - p1.Y()*p2.X()
			dnm := b*d - a*e
			if dnm != 0 {
				/* TODO should these (x & ) have modified "this"? */
				x := (b*f - c*e) / dnm
				y := (c*d - a*f) / dnm
				if a1x <= x && x <= a2x && a1y <= y && y <= a2y &&
					b1x <= x && x <= b2x && b1y <= y && y <= b2y {
					return true
				}
			}
		}
	}
	return false
}

func (this Vector) checkHitDist(p Vector, pp Vector, dist float32) bool {
	var bmvx, bmvy, inaa float32
	bmvx = pp.X()
	bmvy = pp.Y()
	bmvx -= p.X()
	bmvy -= p.Y()
	inaa = bmvx*bmvx + bmvy*bmvy
	if inaa > 0.00001 {
		var sofsx, sofsy, inab, hd float32
		sofsx = this.X()
		sofsy = this.Y()
		sofsx -= p.X()
		sofsy -= p.Y()
		inab = bmvx*sofsx + bmvy*sofsy
		if inab >= 0 && inab <= inaa {
			hd = sofsx*sofsx + sofsy*sofsy - inab*inab/inaa
			if hd >= 0 && hd <= dist {
				return true
			}
		}
	}
	return false
}

func (this Vector) vctSize() float32 {
	return float32(math.Sqrt(float64(this.X()*this.X() + this.Y()*this.Y())))
}

func (this Vector) dist(v Vector) float32 {
	return this.distFloat(v.X(), v.Y())
}

func fabs(f float32) float32 {
	return float32(math.Abs(float64(f)))
}

func (this Vector) distFloat(px float32 /* = 0 */, py float32 /* = 0 */) float32 {
	ax := fabs(this.X() - px)
	ay := fabs(this.Y() - py)
	if ax > ay {
		return ax + ay/2
	} else {
		return ay + ax/2
	}
}

func (this Vector) containsVector(p Vector, r float32 /*= 1*/) bool {
	return this.containsFloat(p.X(), p.Y(), r)
}

func (this Vector) containsFloat(px float32, py float32, r float32 /* = 1 */) bool {
	return px >= -this.X()*r && px <= this.X()*r && py >= -this.Y()*r && py <= this.Y()*r
}

func (this Vector) toString() string {
	return fmt.Sprintf("(%v, %v)", this.X(), this.Y())
}

type Vector3 struct {
	*mgl32.Vec3
}

func NewVector3(x float32, y float32, z float32) Vector3 {
	return Vector3{&mgl32.Vec3{x, y, z}}
}

func (this Vector) MulV(f float32) Vector {
	v := this.Mul(f)
	return Vector{&v}
}

func (this Vector) AddV(v2 Vector) Vector {
	v := this.Add(mgl32.Vec2{v2.X(), v2.Y()})
	return Vector{&v}
}

func (this Vector) SetX(x float32) Vector {
	return NewVector(x, this.Y())
}

func (this Vector) SetY(y float32) Vector {
	return NewVector(this.X(), y)
}

func Cos32(d float32) float32 {
	return float32(math.Cos(float64(d)))
}

func Sin32(d float32) float32 {
	return float32(math.Sin(float64(d)))
}

func (this Vector3) rollX(d float32) Vector3 {
	ty := this.Y()*Cos32(d) - this.Z()*Sin32(d)
	z := this.Y()*Sin32(d) + this.Z()*Cos32(d)
	y := ty
	return NewVector3(this.X(), y, z)
}

func (this Vector3) rollY(d float32) Vector3 {
	tx := this.X()*Cos32(d) - this.Z()*Sin32(d)
	z := this.X()*Sin32(d) + this.Z()*Cos32(d)
	x := tx
	return NewVector3(x, this.Y(), z)
}

func (this Vector3) rollZ(d float32) Vector3 {
	tx := this.X()*Cos32(d) - this.Y()*Sin32(d)
	y := this.X()*Sin32(d) + this.Y()*Cos32(d)
	x := tx
	return NewVector3(x, y, this.Z())
}

func (this Vector3) blend(v1 Vector3, v2 Vector3, ratio float32) Vector3 {
	x := v1.X()*ratio + v2.X()*(1-ratio)
	y := v1.Y()*ratio + v2.Y()*(1-ratio)
	z := v1.Z()*ratio + v2.Z()*(1-ratio)
	return NewVector3(x, y, z)
}
