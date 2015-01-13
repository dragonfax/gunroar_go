/*
 * $Id: shot.d,v 1.2 2005/07/03 07:05:22 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */

/**
 * Player's shot.
 */

const SPEED = 0.6
const LANCE_SPEED = 0.5

var shotShape ShotShape
var lanceShape LanceShape

type Shot struct {
  field Field
  pos Vector
  cnt int
  hitCnt int
  deg float32
  damage int
  lance bool
}

func initShots() {
	shotShape = NewShotShape()
	lanceShape = NewLanceShape()
}

func closeShots() {
	shotShape.close()
}

func NewShot(f Field, p Vector, d float32, lance bool /*= false*/, dmg int /*= -1*/) *Shot {
	s = new(Shot)
	s.damage = 1
	s.field = f
	s.pos.x = p.x
	s.pos.y = p.y
	s.deg = d
	s.lance = lance
	if (lance) {
		s.damage = 10
	} else {
		s.damage = 1
	}
	if (dmg >= 0) {
		s.damage = dmg
	}
	actors[s] = true
	return s
}

  public override void move() {
    cnt++
    if (hitCnt > 0) {
      hitCnt++
      if (hitCnt > 30)
        remove()
      return
    }
    float sp
    if (!lance) {
      sp = SPEED
    } else {
      if (cnt < 10)
        sp = LANCE_SPEED * cnt / 10
      else
        sp = LANCE_SPEED
    }
    pos.x += sin(_deg) * sp
    pos.y += cos(_deg) * sp
    pos.y -= field.lastScrollY
    if (field.getBlock(pos) >= Field.ON_BLOCK_THRESHOLD ||
        !field.checkInOuterField(pos) || pos.y > field.size.y)
      remove()
    if (lance) {
      enemies.checkShotHit(pos, lanceShape, this)
    } else {
      bullets.checkShotHit(pos, shape, this)
      enemies.checkShotHit(pos, shape, this)
    }
  }

func (s *Shot) remove() {
	if (s.lance && s.hitCnt <= 0) {
		s.hitCnt = 1
		return
	}
	delete(actors,s)
}

func (s *Shot)  removeHitToBullet() {
	s.removeHit()
}

  public void removeHitToEnemy(bool isSmallEnemy = false) {
    if (isSmallEnemy && lance)
      return
    SoundManager.playSe("hit.wav")
    removeHit()
  }

  private void removeHit() {
    remove()
    int sn
    if (lance) {
      for (int i = 0; i < 10; i++) {
        Smoke s = smokes.getInstanceForced()
        float d = _deg + rand.nextSignedFloat(0.1)
        float sp = rand.nextFloat(LANCE_SPEED)
        s.set(pos, sin(d) * sp, cos(d) * sp, 0,
              Smoke.SmokeType.LANCE_SPARK, 30 + rand.nextInt(30), 1)
        s = smokes.getInstanceForced()
        d = _deg + rand.nextSignedFloat(0.1)
        sp = rand.nextFloat(LANCE_SPEED)
        s.set(pos, -sin(d) * sp, -cos(d) * sp, 0,
              Smoke.SmokeType.LANCE_SPARK, 30 + rand.nextInt(30), 1)
      }
    } else {
      Spark s = sparks.getInstanceForced()
      float d = _deg + rand.nextSignedFloat(0.5)
      s.set(pos, sin(d) * SPEED, cos(d) * SPEED,
            0.6 + rand.nextSignedFloat(0.4), 0.6 + rand.nextSignedFloat(0.4), 0.1, 20)
      s = sparks.getInstanceForced()
      d = _deg + rand.nextSignedFloat(0.5)
      s.set(pos, -sin(d) * SPEED, -cos(d) * SPEED,
            0.6 + rand.nextSignedFloat(0.4), 0.6 + rand.nextSignedFloat(0.4), 0.1, 20)
    }
  }

  public override void draw() {
    if (lance) {
      float x = pos.x, y = pos.y
      float size = 0.25, a = 0.6
      int hc = hitCnt
      for (int i = 0; i < cnt / 4 + 1; i++) {
        size *= 0.9
        a *= 0.8
        if (hc > 0) {
          hc--
          continue
        }
        float d = i * 13 + cnt * 3
        for (int j = 0; j < 6; j++) {
          gl.PushMatrix()
          gl.Translatef(x, y, 0)
          gl.Rotatef(-_deg * 180 / PI, 0, 0, 1)
          gl.Rotatef(d, 0, 1, 0)
          Screen.setColor(0.4, 0.8, 0.8, a)
          gl.Begin(gl.LINE_LOOP)
          gl.Vertex3f(-size, LANCE_SPEED, size / 2)
          gl.Vertex3f(size, LANCE_SPEED, size / 2)
          gl.Vertex3f(size, -LANCE_SPEED, size / 2)
          gl.Vertex3f(-size, -LANCE_SPEED, size / 2)
          gl.End()
          Screen.setColor(0.2, 0.5, 0.5, a / 2)
          gl.Begin(gl.TRIANGLE_FAN)
          gl.Vertex3f(-size, LANCE_SPEED, size / 2)
          gl.Vertex3f(size, LANCE_SPEED, size / 2)
          gl.Vertex3f(size, -LANCE_SPEED, size / 2)
          gl.Vertex3f(-size, -LANCE_SPEED, size / 2)
          gl.End()
          gl.PopMatrix()
          d += 60
        }
        x -= sin(deg) * LANCE_SPEED * 2
        y -= cos(deg) * LANCE_SPEED * 2
      }
    } else {
      gl.PushMatrix()
      Screen.gl.Translate(pos)
      gl.Rotatef(-_deg * 180 / PI, 0, 0, 1)
      gl.Rotatef(cnt * 31, 0, 1, 0)
      shape.draw()
      gl.PopMatrix()
    }
  }

  public float deg() {
    return _deg
  }

  public int damage() {
    return _damage
  }

  public bool removed() {
    if (hitCnt > 0)
      return true
    else
      return false
  }
}

public class ShotPool: ActorPool!(Shot) {
  public this(int n, Object[] args) {
    super(n, args)
  }

  public bool existsLance() {
    foreach (Shot s; actor)
      if (s.exists)
        if (s.lance && !s.removed)
          return true
    return false
  }
}

public class ShotShape: CollidableDrawable {
  protected override void createDisplayList() {
    Screen.setColor(0.1, 0.33, 0.1)
    gl.Begin(gl.QUADS)
    gl.Vertex3f(0, 0.3, 0.1)
    gl.Vertex3f(0.066, 0.3, -0.033)
    gl.Vertex3f(0.1, -0.3, -0.05)
    gl.Vertex3f(0, -0.3, 0.15)
    gl.Vertex3f(0.066, 0.3, -0.033)
    gl.Vertex3f(-0.066, 0.3, -0.033)
    gl.Vertex3f(-0.1, -0.3, -0.05)
    gl.Vertex3f(0.1, -0.3, -0.05)
    gl.Vertex3f(-0.066, 0.3, -0.033)
    gl.Vertex3f(0, 0.3, 0.1)
    gl.Vertex3f(0, -0.3, 0.15)
    gl.Vertex3f(-0.1, -0.3, -0.05)
    gl.End()
  }

  protected override void setCollision() {
    _collision = new Vector(0.33, 0.33)
  }
}

public class LanceShape: Collidable {
  mixin CollidableImpl
 private:
  Vector _collision

  public this() {
    _collision = new Vector(0.66, 0.66)
  }

  public Vector collision() {
    return _collision
  }
}
