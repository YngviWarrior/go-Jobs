package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
)

func releasePaymentRefundGame(tx *sql.Tx, id uint64) bool {
	var b entities.BinaryOptionGameBet

	err := tx.QueryRow(`
	SELECT id, hash_id, id_game, id_usuario, id_choice, id_balance, bet_amount_dolar, amount_win_dolar, price_amount_selected, 
		status_received_win_payment, id_trader_follower, bot_use_status, date_register, bonus_trader_percent_from_tax_bet_win, 
		bonus_indication_percent_from_tax_bet_win, status_received_refund_payment, refund, deleted
	FROM binary_option_game_bet
	WHERE id_game = ? AND status_received_win_payment = 1 AND refund = 1
	LIMIT 0,1`, id).Scan(&b.Id, &b.HashId, &b.IdGame, &b.IdUsuario, &b.IdChoice, &b.IdBalance, &b.BetAmountDolar, &b.AmountWinDolar,
		&b.PriceAmountSelected, &b.StatusReceivedWinPayment, &b.IdTraderFollower, &b.BotUseStatus, &b.DateRegister, &b.BonusTraderPercentFromTaxBetWin,
		&b.BonusIndicationPercentFromTaxBetWin, &b.StatusReceivedRefundPayment, &b.Refund, &b.Deleted)

	if err != nil {
		fmt.Println("RPRG 1: " + err.Error())
	}

	rows, err := tx.Query(`
		SELECT b.id, b.id_usuario, (b.amount_win_dolar + b.bet_amount_dolar) as amount_win_dolar, b.id_balance
		FROM binary_option_game_bet b
		WHERE b.id_game = ? AND b.status_received_win_payment = 0 AND b.amount_win_dolar > 0
		-- GROUP BY b.id_usuario
	`, id)

	if err != nil {
		fmt.Println("RPRG 2: " + err.Error())
		return false
	}

	var bet entities.BinaryOptionGameBet
	var listBet []*entities.BinaryOptionGameBet
	for rows.Next() {
		err := rows.Scan(&bet.Id, &bet.IdUsuario, &bet.AmountWinDolar, &bet.IdBalance)

		if err != nil {
			fmt.Println("RPRG 3: " + err.Error())
			return false
		}

		listBet = append(listBet, &bet)
	}

	if len(listBet) > 0 {
		for _, v := range listBet {
			modifyBalance(tx, v.IdUsuario, v.IdBalance, 10, v.AmountWinDolar, v.Id, false)
		}
	} else {
		fmt.Println("RPRG 4: No Bets.")
		return false
	}

	res, err := tx.Exec(`
		UPDATE binary_option_game_bet b
		SET b.status_received_refund_payment = 1
		WHERE b.id_game = ?
	`, b.IdGame)

	affcRows, _ := res.RowsAffected()

	if err != nil || affcRows == 0 {
		fmt.Println("RPRG 5: ")
		return false
	}

	return true
}
