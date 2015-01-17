/*
 * $Id: vector.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2004 Kenta Cho. Some rights reserved.
 */
package main

type Vector struct {
	x, y float32
}

/* dot product */
func (this Vector) Mul(v2 Vector) float32 {
	return this.x*v2.x + this.y*v2.y
}

func (this Vector) getElement(v Vector) Vector {
	var rsl Vector
	ll := v.Mul(v)
	if ll != 0 {
		mag := this.Mul(v)
		rsl.x = mag * v.x / ll
		rsl.y = mag * v.y / ll
	} else {
		rsl.x = 0
		rsl.y = 0
	}
	return rsl
}

func (this *Vector) AddAssign(v Vector) {
	this.x += v.x
	this.y += v.y
}

func (this *Vector) SubAssign(v Vector) {
	this.x -= v.x
	this.y -= v.y
}

func (this *Vector) MulAssign(a float32) {
	this.x *= a
	this.y *= a
}

func (this *Vector) DivAssign(a float32) {
	this.x /= a
	this.y /= a
}

func (this Vector) checkSide(pos1 Vector, pos2 Vector) {
	xo := pos2.x - pos1.x
	yo := pos2.y - pos1.y
	if xo == 0 {
		if yo == 0 {
			return 0
		}
		if yo > 0 {
			return this.x - pos1.x
		} else {
			return pos1.x - this.x
		}
	} else if yo == 0 {
		if xo > 0 {
			return pos1.y - this.y
		} else {
			return this.y - pos1.y
		}
	} else {
		if xo*yo > 0 {
			return (this.x-pos1.x)/xo - (this.y-pos1.y)/yo
		} else {
			return -(this.x-pos1.x)/xo + (this.y-pos1.y)/yo
		}
	}
}

func (this Vector) checkSide3(pos1 Vector, pos2 Vector, ofs Vector) float32 {
	xo := pos2.x - pos1.x
	yo := pos2.y - pos1.y
	mx := this.x + ofs.x
	my := this.y + ofs.y
	if xo == 0 {
		if yo == 0 {
			return 0
		}
		if yo > 0 {
			return mx - pos1.x
		} else {
			return pos1.x - mx
		}
	} else if yo == 0 {
		if xo > 0 {
			return pos1.y - my
		} else {
			return my - pos1.y
		}
	} else {
		if xo*yo > 0 {
			return (mx-pos1.x)/xo - (my-pos1.y)/yo
		} else {
			return -(mx-pos1.x)/xo + (my-pos1.y)/yo
		}
	}
}

func (this Vector) checkCross(p Vector, p1 Vector, p2 Vector, width float32) bool {
	var a1x, a1y, a2x, a2y float32
	if this.x < p.x {
		a1x = this.x - width
		a2x = p.x + width
	} else {
		a1x = p.x - width
		a2x = this.x + width
	}
	if this.y < p.y {
		a1y = this.y - width
		a2y = p.y + width
	} else {
		a1y = p.y - width
		a2y = this.y + width
	}
	var b1x, b1y, b2x, b2y float32
	if p2.y < p1.y {
		b1y = p2.y - width
		b2y = p1.y + width
	} else {
		b1y = p1.y - width
		b2y = p2.y + width
	}
	if a2y >= b1y && b2y >= a1y {
		if p2.x < p1.x {
			b1x = p2.x - width
			b2x = p1.x + width
		} else {
			b1x = p1.x - width
			b2x = p2.x + width
		}
		if a2x >= b1x && b2x >= a1x {
			a := this.y - p.y
			b := p.x - this.x
			c := p.x*this.y - p.y*this.x
			d := p2.y - p1.y
			e := p1.x - p2.x
			f := p1.x*p2.y - p1.y*p2.x
			dnm := b*d - a*e
			if dnm != 0 {
				x := (b*f - c*e) / dnm
				y := (c*d - a*f) / dnm
				if a1x <= this.x && this.x <= a2x && a1y <= this.y && this.y <= a2y &&
					b1x <= this.x && this.x <= b2x && b1y <= this.y && this.y <= b2y {
					return true
				}
			}
		}
	}
	return false
}

func (this Vector) checkHitDist(p Vector, pp Vector, dist float32) bool {
	var bmvx, bmvy, inaa float32
	bmvx = pp.x
	bmvy = pp.y
	bmvx -= p.x
	bmvy -= p.y
	inaa = bmvx*bmvx + bmvy*bmvy
	if inaa > 0.00001 {
		var sofsx, sofsy, inab, hd float32
		sofsx = this.x
		sofsy = this.y
		sofsx -= p.x
		sofsy -= p.y
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
	return sqrt(x*x + y*y)
}

func (this Vector) distVector(v Vector) float32 {
	return dist(v.x, v.y)
}

func (this Vector) dist(px float32 /* = 0 */, py float32 /* = 0 */) float32 {
	ax := fabs32(this.x - px)
	ay := fabs32(this.y - py)
	if ax > ay {
		return ax + ay/2
	} else {
		return ay + ax/2
	}
}

func (this Vector) containsVector(p Vector, r float32 /* = 1 */) bool {
	return contains(p.x, p.y, r)
}

func (this Vector) contains(px float32, py float32, r float32 /*= 1*/) bool {
	return px >= -this.x*r && px <= this.x*r && py >= -this.y*r && py <= this.y*r
}

type Vector3 struct {
	x, y, z float32
}

func (this *Vector3) rollX(d float32) {
	ty := this.y*Cos32(d) - this.z*Sin32(d)
	this.z = this.y*Sin32(d) + this.z*Cos32(d)
	this.y = ty
}

func (this *Vector3) rollY(d float32) {
	tx := this.x*Cos32(d) - this.z*Sin32(d)
	this.z = this.x*Sin32(d) + this.z*Cos32(d)
	this.x = tx
}

func (this *Vector3) rollZ(d float32) {
	tx := this.x*Cos32(d) - this.y*Sin32(d)
	this.y = this.x*Sin32(d) + this.y*Cos32(d)
	this.x = tx
}

func (this *Vector3) blend(v1 Vector3, v2 Vector3, ratio float32) {
	this.x = v1.x*ratio + v2.x*(1-ratio)
	this.y = v1.y*ratio + v2.y*(1-ratio)
	this.z = v1.z*ratio + v2.z*(1-ratio)
}

func (this *Vector3) AddAssign(v Vector3) {
	this.x += v.x
	this.y += v.y
	this.z += v.z
}

func (this *Vector3) SubAssign(v Vector3) {
	this.x -= v.x
	this.y -= v.y
	this.z -= v.z
}

func (this *Vector3) MulAssign(a float32) {
	this.x *= a
	this.y *= a
	this.z *= a
}

func (this *Vector3) DivAssign(a float32) {
	x /= a
	y /= a
	z /= a
}
