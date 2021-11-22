package vector

type Vector struct {
  X, Y float64
}

func New(x, y float64) Vector {
  v := Vector{X: x, Y: y}
  return v
}

func (this Vector) OpMul(v Vector) float64 {
  return this.X * v.X + this.Y * v.Y
}

func (this Vector) getElement(v Vector) Vector {
  var rsl Vector
  ll := v.OpMul(v);
  if (ll != 0) {
    mag := this.OpMul(v);
    rsl.X = mag * v.X / ll;
    rsl.Y = mag * v.Y / ll;
  } else {
    rsl.X = 0
    rsl.Y = 0
  }
  return rsl
}

func (this *Vector) OpAddAssign(v Vector) {
  this.X += v.X;
  this.Y += v.Y;
}

func (this *Vector) OpSubAssign(v Vector) {
  this.X -= v.X;
  this.Y -= v.Y;
}

func (this *Vector) OpMulAssign(a float64) {	
  this.X *= a;
  this.Y *= a;
}

func (this *Vector) OpDivAssign(a float64) {	
  this.X /= a
  this.Y /= a
}

func (this Vector) CheckSide(pos1, pos2 Vector) float64 {
  xo := pos2.X - pos1.X
  yo := pos2.Y - pos1.Y
  if (xo == 0) {
    if (yo == 0) {
      return 0
    }
    if (yo > 0) {
      return this.X - pos1.X
    } else {
      return pos1.X - this.X;
    }
  } else if (yo == 0) {
    if (xo > 0) {
      return pos1.Y - this.Y
  } else {
    return this.Y - pos1.Y
  }
  } else {
    if (xo * yo > 0) {
      return (this.X - pos1.X) / xo - (this.Y - pos1.Y) / yo
    } else {
      return -(this.X - pos1.X) / xo + (this.Y - pos1.Y) / yo
    }
  }
}

func (this *Vector) CheckSide(Vector pos1, Vector pos2, Vector ofs) float64 {
  xo := pos2.x - pos1.x
  yo := pos2.y - pos1.y
  mx := x + ofs.x
  my := y + ofs.y
  if (xo == 0) {
    if (yo == 0) {
return 0
    }
    if (yo > 0) {
return mx - pos1.x
    } else {
return pos1.x - mx
    }
  } else if (yo == 0) {
    if (xo > 0) {
return pos1.y - my
    } else {
return my - pos1.y
    }
  } else {
    if (xo * yo > 0) {
return (mx - pos1.x) / xo - (my - pos1.y) / yo
    } else {
return -(mx - pos1.x) / xo + (my - pos1.y) / yo
    }
  }
}

func (this *Vector) CheckCross(Vector p, Vector p1, Vector p2, float width) bool {
  var a1x, a1y, a2x, a2y float64
  if (x < p.x) {
    a1x = x - width 
    a2x = p.x + width
  } else {
    a1x = p.x - width 
    a2x = x + width
  }
  if (y < p.y) {
    a1y = y - width 
    a2y = p.y + width
  } else {
    a1y = p.y - width 
    a2y = y + width
  }
  var b1x, b1y, b2x, b2y float64
  if (p2.y < p1.y) {
    b1y = p2.y - width 
    b2y = p1.y + width
  } else {
    b1y = p1.y - width 
    b2y = p2.y + width
  }
  if (a2y >= b1y && b2y >= a1y) {
    if (p2.x < p1.x) {
      b1x = p2.x - width 
      b2x = p1.x + width
    } else {
      b1x = p1.x - width 
      b2x = p2.x + width
    }
    if (a2x >= b1x && b2x >= a1x) {
      a := y - p.y
      b := p.x - x
      c := p.x * y - p.y * x
      d := p2.y - p1.y
      e := p1.x - p2.x
      f := p1.x * p2.y - p1.y * p2.x
      dnm := b * d - a * e
      if (dnm != 0) {
        x := (b*f - c*e) / dnm
        y := (c*d - a*f) / dnm
        if (a1x <= x && x <= a2x && a1y <= y && y <= a2y &&
            b1x <= x && x <= b2x && b1y <= y && y <= b2y) {
          return true
            }
      }
    }
  }
  return false
}

func (this *Vector) CheckHitDist(Vector p, Vector pp, float dist) bool {
  var bmvx, bmvy, inaa float64
  bmvx = pp.x
  bmvy = pp.y
  bmvx -= p.x
  bmvy -= p.y
  inaa = bmvx * bmvx + bmvy * bmvy
  if (inaa > 0.00001) {
    var sofsx, sofsy, inab, hd float64
    sofsx = x
    sofsy = y
    sofsx -= p.x
    sofsy -= p.y
    inab = bmvx * sofsx + bmvy * sofsy
    if (inab >= 0 && inab <= inaa) {
hd = sofsx * sofsx + sofsy * sofsy - inab * inab / inaa
if (hd >= 0 && hd <= dist) {
  return true
}
    }
  }
  return false
}

func (this *Vector) VctSize() float64 {
  return sqrt(x * x + y * y)
}

func (this *Vector) DistVector(Vector v) float64 {
  return dist(v.x, v.y)
}

func (this *Vector) Dist(px, py float64) float64 {
  ax := fabs(x - px)
  ay := fabs(y - py)
  if (ax > ay) {
    return ax + ay / 2
  } else {
    return ay + ax / 2
  }
}

func (this *Vector) Contains(Vector p, float r = 1) bool {
  return contains(p.x, p.y, r)
}

func (this *Vector) Contains(float px, float py, float r = 1) bool {
if (px >= -x * r && px <= x * r && py >= -y * r && py <= y * r) {
    return true
} else {
    return false
}
}

func (this *Vector) ToString() string {
  return "(" ~ std.string.toString(x) ~ ", " ~ std.string.toString(y) ~ ")"
}

type Vector3 struct {
X, Y, Z float64
}

func New3(x, y, z float64) Vector3 {
  this := Vector3{}
  this.X= x
  this.Y = y
  this.Z = Z
  return this
}

func (this *Vector3) RollX(float d) {
  ty := y * cos(d) - z * sin(d)
  z = y * sin(d) + z * cos(d)
  y = ty
}

func (this *Vector3) RollY(float d) {
  tx := x * cos(d) - z * sin(d)
  z = x * sin(d) + z * cos(d)
  x = tx
}

func (this *Vector3) RollZ(float d) {
  tx := x * cos(d) - y * sin(d)
  y = x * sin(d) + y * cos(d)
  x = tx
}

func (this *Vector3) Blend(Vector3 v1, Vector3 v2, float ratio) {
  x = v1.x * ratio + v2.x * (1 - ratio)
  y = v1.y * ratio + v2.y * (1 - ratio)
  z = v1.z * ratio + v2.z * (1 - ratio)
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