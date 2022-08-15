package finishgame

import (
	"processgame/entities"
)

// Separate win, lose and equal bets and sum values generating win, loss and equal amounts of all bets.
func processGameWinLose(betList []*entities.BinaryOptionGameBet, closePrice float64, gameProfitPercent float64) (b entities.BestResultGame) {
	b.Price = closePrice

	for _, bet := range betList {
		if closePrice == bet.PriceAmountSelected {
			b.ListPlayersEqual = append(b.ListPlayersEqual, bet)
			b.TotalEqualDolar = b.TotalEqualDolar + bet.BetAmountDolar
		} else if (bet.IdChoice == 1 && closePrice > bet.PriceAmountSelected) || (bet.IdChoice == 2 && closePrice < bet.PriceAmountSelected) {
			b.ListPlayersWin = append(b.ListPlayersWin, bet)

			if bet.IdBalance == 3 {
				b.TotalWinDolar = b.TotalWinDolar + bet.BetAmountDolar
			}
		} else {
			if bet.IdBalance == 3 {
				b.TotalLoseDolar = b.TotalLoseDolar + bet.BetAmountDolar
			}
		}
	}

	return
}
