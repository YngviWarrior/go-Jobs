package finishgame

import math "math/big"

func resultToPercent(result *math.Float, total *math.Float, decimais int64) float64 {
	temp := math.NewFloat(0)
	x, _ := result.Float64()
	y, _ := total.Float64()

	if x > 0 && y > 0 {
		temp = temp.Quo(result, total)

		percent := math.NewFloat(100)
		temp.Mul(temp, percent)
	} else {
		temp = math.NewFloat(0)
	}

	resp, _ := temp.Float64()
	return resp
}
