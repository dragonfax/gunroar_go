package vector

import (
	"fmt"
	"math"
)

type Vector struct {
	X, Y float64
}

func New(x, y float64) Vector {
	v := Vector{X: x, Y: y}
	return v
}

func (this Vector) OpMul(v Vector) float64 {
	return this.X*v.X + this.Y*v.Y
}

func (this Vector) getElement(v Vector) Vector {
	var rsl Vector
	ll := v.OpMul(v)
	if ll != 0 {
		mag := this.OpMul(v)
		rsl.X = mag * v.X / ll
		rsl.Y = mag * v.Y / ll
	} else {
		rsl.X = 0
		rsl.Y = 0
	}
	return rsl
}

func (this *Vector) OpAddAssign(v Vector) {
	this.X += v.X
	this.Y += v.Y
}

func (this *Vector) OpSubAssign(v Vector) {
	this.X -= v.X
	this.Y -= v.Y
}

func (this *Vector) OpMulAssign(a float64) {
	this.X *= a
	this.Y *= a
}

func (this *Vector) OpDivAssign(a float64) {
	this.X /= a
	this.Y /= a
}

func (this Vector) CheckSide(pos1, pos2 Vector) float64 {
	xo := pos2.X - pos1.X
	yo := pos2.Y - pos1.Y
	if xo == 0 {
		if yo == 0 {
			return 0
		}
		if yo > 0 {
			return this.X - pos1.X
		} else {
			return pos1.X - this.X
		}
	} else if yo == 0 {
		if xo > 0 {
			return pos1.Y - this.Y
		} else {
			return this.Y - pos1.Y
		}
	} else {
		if xo*yo > 0 {
			return (this.X-pos1.X)/xo - (this.Y-pos1.Y)/yo
		} else {
			return -(this.X-pos1.X)/xo + (this.Y-pos1.Y)/yo
		}
	}
}

func (this Vector) CheckSide3(pos1, pos2, ofs Vector) float64 {
	xo := pos2.X - pos1.X
	yo := pos2.Y - pos1.Y
	mx := this.X + ofs.X
	my := this.Y + ofs.Y
	if xo == 0 {
		if yo == 0 {
			return 0
		}
		if yo > 0 {
			return mx - pos1.X
		} else {
			return pos1.X - mx
		}
	} else if yo == 0 {
		if xo > 0 {
			return pos1.Y - my
		} else {
			return my - pos1.Y
		}
	} else {
		if xo*yo > 0 {
			return (mx-pos1.X)/xo - (my-pos1.Y)/yo
		} else {
			return -(mx-pos1.X)/xo + (my-pos1.Y)/yo
		}
	}
}

func (this Vector) CheckCross(p, p1, p2 Vector, width float64) bool {
	var a1x, a1y, a2x, a2y float64
	if this.X < p.X {
		a1x = this.X - width
		a2x = p.X + width
	} else {
		a1x = p.X - width
		a2x = this.X + width
	}
	if this.Y < p.Y {
		a1y = this.Y - width
		a2y = p.Y + width
	} else {
		a1y = p.Y - width
		a2y = this.Y + width
	}
	var b1x, b1y, b2x, b2y float64
	if p2.Y < p1.Y {
		b1y = p2.Y - width
		b2y = p1.Y + width
	} else {
		b1y = p1.Y - width
		b2y = p2.Y + width
	}
	if a2y >= b1y && b2y >= a1y {
		if p2.X < p1.X {
			b1x = p2.X - width
			b2x = p1.X + width
		} else {
			b1x = p1.X - width
			b2x = p2.X + width
		}
		if a2x >= b1x && b2x >= a1x {
			a := this.Y - p.Y
			b := p.X - this.X
			c := p.X*this.Y - p.Y*this.X
			d := p2.Y - p1.Y
			e := p1.X - p2.X
			f := p1.X*p2.Y - p1.Y*p2.X
			dnm := b*d - a*e
			if dnm != 0 {
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

func (this Vector) CheckHitDist(p, pp Vector, dist float64) bool {
	var bmvx, bmvy, inaa float64
	bmvx = pp.X
	bmvy = pp.Y
	bmvx -= p.X
	bmvy -= p.Y
	inaa = bmvx*bmvx + bmvy*bmvy
	if inaa > 0.00001 {
		var sofsx, sofsy, inab, hd float64
		sofsx = this.X
		sofsy = this.Y
		sofsx -= p.X
		sofsy -= p.Y
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

func (this Vector) VctSize() float64 {
	return math.Sqrt(this.X*this.X + this.Y*this.Y)
}

func (this *Vector) DistVector(v Vector) float64 {
	return this.Dist(v.X, v.Y)
}

func (this Vector) Dist(px, py float64) float64 {
	ax := math.Abs(this.X - px)
	ay := math.Abs(this.Y - py)
	if ax > ay {
		return ax + ay/2
	}
	return ay + ax/2
}

func (this Vector) ContainsVector(p Vector, r float64) bool {
	return this.Contains(p.X, p.Y, r)
}

func (this Vector) Contains(px, py, r float64) bool {
	if px >= -this.X*r && px <= this.X*r && py >= -this.Y*r && py <= this.Y*r {
		return true
	}
	return false
}

func (this Vector) ToString() string {
	return fmt.Sprintf("(%f, %f)", this.X, this.Y)
}

type Vector3 struct {
	X, Y, Z float64
}

func New3(x, y, z float64) Vector3 {
	this := Vector3{}
	this.X = x
	this.Y = y
	this.Z = z
	return this
}

func (this *Vector3) RollX(d float64) {
	ty := this.Y*math.Cos(d) - this.Z*math.Sin(d)
	this.Z = this.Y*math.Sin(d) + this.Z*math.Cos(d)
	this.Y = ty
}

func (this *Vector3) RollY(d float64) {
	tx := this.X*math.Cos(d) - this.Z*math.Sin(d)
	this.Z = this.X*math.Sin(d) + this.Z*math.Cos(d)
	this.X = tx
}

func (this *Vector3) RollZ(d float64) {
	tx := this.X*math.Cos(d) - this.Y*math.Sin(d)
	this.Y = this.X*math.Sin(d) + this.Y*math.Cos(d)
	this.X = tx
}

func (this *Vector3) Blend(v1, v2 Vector3, ratio float64) {
	this.X = v1.X*ratio + v2.X*(1-ratio)
	this.Y = v1.Y*ratio + v2.Y*(1-ratio)
	this.Z = v1.Z*ratio + v2.Z*(1-ratio)
}

func (this *Vector3) OpAddAssign(v Vector3) {
	this.X += v.X
	this.Y += v.Y
	this.Z += v.Z
}

func (this *Vector3) OpSubAssign(v Vector3) {
	this.X -= v.X
	this.Y -= v.Y
	this.Z -= v.Z
}

func (this *Vector3) OpMulAssign(a float64) {
	this.X *= a
	this.Y *= a
	this.Z *= a
}

func (this *Vector3) OpDivAssign(a float64) {
	this.X /= a
	this.Y /= a
	this.Z /= a
}
