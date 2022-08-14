package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
)

func searchBestPriceForEndGame(db *sql.DB, game *entities.BinaryOptionGame, profit float64) (info *entities.BestResultGame) {
	query := `
	SELECT id, hash_id, id_game, id_usuario, id_choice, id_balance, bet_amount_dolar, amount_win_dolar, price_amount_selected, 
		status_received_win_payment, id_trader_follower, bot_use_status, date_register, bonus_trader_percent_from_tax_bet_win, 
		bonus_indication_percent_from_tax_bet_win, status_received_refund_payment, refund, deleted
	FROM binary_option_game_bet
	WHERE id_game = ?
	-- GROUP BY price_amount_selected`

	rows, err := db.Query(query, game.Id)

	if err != nil {
		fmt.Println(err)
	}

	var betList []*entities.BinaryOptionGameBet

	for rows.Next() {
		var b entities.BinaryOptionGameBet
		err = rows.Scan(
			&b.Id, &b.HashId, &b.IdGame, &b.IdUsuario, &b.IdChoice, &b.IdBalance, &b.BetAmountDolar, &b.AmountWinDolar,
			&b.PriceAmountSelected, &b.StatusReceivedWinPayment, &b.IdTraderFollower, &b.BotUseStatus, &b.DateRegister,
			&b.BonusTraderPercentFromTaxBetWin, &b.BonusIndicationPercentFromTaxBetWin, &b.StatusReceivedRefundPayment, &b.Refund,
			&b.Deleted)

		if err != nil {
			fmt.Println(err)
		}

		betList = append(betList, &b)
	}

	temp := lastCandleInfo(db, game)
	info = &temp

	if len(betList) == 0 {
		return
	}

	info = processGameWinLose(betList, info.Price, game.GameProfitPercent)

	return
}
