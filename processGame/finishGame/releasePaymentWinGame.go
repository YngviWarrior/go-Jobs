package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
)

func releasePaymentWinGame(db *sql.DB, id uint64) {
	rows, err := db.Query(`
	SELECT id, hash_id, id_game, id_usuario, id_choice, id_balance, bet_amount_dolar, amount_win_dolar, price_amount_selected, 
		status_received_win_payment, id_trader_follower, bot_use_status, date_register, bonus_trader_percent_from_tax_bet_win, 
		bonus_indication_percent_from_tax_bet_win, status_received_refund_payment, refund, deleted
	FROM binary_option_game_bet
	WHERE id_game = ? AND status_received_win_payment = 1 AND amount_win_dolar > 0
	LIMIT 0,1`, id)

	if err != nil {
		fmt.Println(err)
	}

	var b entities.BinaryOptionGameBet

	for rows.Next() {
		err := rows.Scan(&b.Id, &b.HashId, &b.IdGame, &b.IdUsuario, &b.IdChoice, &b.IdBalance, &b.BetAmountDolar, &b.AmountWinDolar,
			&b.PriceAmountSelected, &b.StatusReceivedWinPayment, &b.IdTraderFollower, &b.BotUseStatus, &b.DateRegister, &b.BonusTraderPercentFromTaxBetWin,
			&b.BonusIndicationPercentFromTaxBetWin, &b.StatusReceivedRefundPayment, &b.Refund, &b.Deleted)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if b.Id == 0 {
		return
	}

	_, err = db.Exec(`
		UPDATE usuarios u
		JOIN (
			SELECT b.id_usuario, SUM(b.bet_amount_dolar) as bet_amount_dolar
			FROM binary_option_game_bet b
			WHERE b.id_game = ?
				AND b.status_received_win_payment = 0 
				AND b.amount_win_dolar = 0 
				AND b.id_balance = ?
			
			GROUP BY b.id_usuario
		) as t ON t.id_usuario = u.id
		SET u.total_lose_balance_play = IF(
			(u.total_lose_balance_play + t.bet_amount_dolar) > u.total_deposit_balance_play,
			u.total_deposit_balance_play,
			(u.total_lose_balance_play + t.bet_amount_dolar)
			)
	`, id, 3)

	if err != nil {
		fmt.Println(err)
	}

	rows, err = db.Query(`
		SELECT b.id, b.id_usuario, (b.amount_win_dolar + b.bet_amount_dolar) as amount_win_dolar, b.id_balance
		FROM binary_option_game_bet b
		WHERE b.id_game = ? AND b.status_received_win_payment = 0 AND b.amount_win_dolar > 0
		GROUP BY b.id_usuario
	`, id)

	if err != nil {
		fmt.Println(err)
	}

	var bet entities.BinaryOptionGameBet
	var listBet []*entities.BinaryOptionGameBet
	for rows.Next() {
		err := rows.Scan(&bet.Id, &bet.IdUsuario, &bet.AmountWinDolar, &bet.IdBalance)

		if err != nil {
			fmt.Println(err)
		}

		listBet = append(listBet, &bet)
	}

	if len(listBet) > 0 {
		for _, v := range listBet {
			modifyBalance(db, v.IdUsuario, v.IdBalance, 7, v.AmountWinDolar, v.Id, false)
		}
	}

	_, err = db.Exec(`
		UPDATE binary_option_game_bet b
		SET b.status_received_win_payment = 1
		WHERE b.id_game = ?
	`, b.IdGame)

	if err != nil {
		fmt.Println(err)
	}
}
