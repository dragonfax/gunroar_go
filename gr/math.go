/*
 * $Id: math.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2005 Kenta Cho. Some rights reserved.
 */
package gr

import (
	"math"
)

func NormalizeDeg(d *float32) {
	if d < -math.Pi {
		d = math.Pi*2 - (-d % (math.Pi * 2))
	}
	d = (d+math.Pi)%(math.Pi*2) - math.Pi
}

func NormalizeDeg360(d *float32) {
	if d < -180 {
		d = 360 - (-d % 360)
	}
	d = (d+180)%360 - 180
}
