package finishgame

import (
	"processgame/entities"
)

func processGameWinLose(betList []*entities.BinaryOptionGameBet, closePrice float64, gameProfitPercent float64) (b *entities.BestResultGame) {

	b.Price = closePrice
	for _, v := range betList {

		if closePrice == v.PriceAmountSelected {
			b.ListPlayersEqual = append(b.ListPlayersEqual, v)
			b.TotalEqualDolar = b.TotalEqualDolar + v.BetAmountDolar
		} else if (v.IdChoice == 1 && closePrice > v.PriceAmountSelected) || (v.IdChoice == 2 && closePrice < v.PriceAmountSelected) {
			b.ListPlayersWin = append(b.ListPlayersWin, v)

			if v.IdBalance == 3 {
				b.TotalWinDolar = b.TotalWinDolar + v.BetAmountDolar
			}
		} else {
			if v.IdBalance == 3 {
				b.TotalLoseDolar = b.TotalLoseDolar + v.BetAmountDolar
			}
		}
	}

	return
}
