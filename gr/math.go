package main

import "math"

func normalizeDeg(d float64) float64 {
	if d < -math.Pi {
		d = math.Pi*2 - math.Mod(-d, (math.Pi*2))
	}
	d = math.Mod((d+math.Pi), (math.Pi*2)) - math.Pi
	return d
}

func normalizeDeg360(d float64) float64 {
	if d < -180 {
		d = 360 - math.Mod(-d, 360)
	}
	d = math.Mod((d+180), 360) - 180
	return d
}
