package finishgame

import math "math/big"

func bcSimplesComp(value1 float64, operator string, value2 float64, decimals int64) bool {
	v1 := math.NewFloat(value1)
	v2 := math.NewFloat(value2)

	result := v1.Cmp(v2)

	if result == -1 && operator == "<" {
		return true
	}

	return false
}
