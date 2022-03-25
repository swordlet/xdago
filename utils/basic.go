package utils

import (
	"math"
)

func Amount2xdag(amount uint64) float64 {
	integer := amount >> 32
	temp := amount - (integer << 32)
	decimal := float64(temp) / math.Pow(2, 32)

	return math.Round((float64(integer)+decimal)*100) / 100
}
