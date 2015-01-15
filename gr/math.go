/*
 * $Id: math.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

import (
	"github.com/dragonfax/gunroar_go/rand"
	"math"
)

const Pi32 = float32(math.Pi)

func normalizeDeg(d float32) float32 {
	if d < -Pi32 {
		d = Pi32*2 - (Mod32(-d, (Pi32 * 2)))
	}
	d = Mod32((d+Pi32), (Pi32*2)) - Pi32
	return d
}

func normalizeDeg360(d float32) float32 {
	if d < -180 {
		d = 360 - Mod32(-d, 360)
	}
	d = Mod32((d+180), 360) - 180
	return d
}
