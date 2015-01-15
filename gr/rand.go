package gr

import (
	builtin "math/rand"
)

func nextSignedFloat(n float32) float32 {
	return builtin.Float32()*(n*2) - n
}

func nextFloat(n float32) float32 {
	return builtin.Float32() * n
}

func nextInt(n int) int {
	return builtin.Int() % n
}

func nextSignedInt(n int) int {
	return builtin.Int()%(n*2+1) - n
}
