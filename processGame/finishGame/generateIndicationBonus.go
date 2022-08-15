package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
	"time"
)

func generateIndicationBonus(tx *sql.Tx, id uint64) bool {
	var b entities.BonusIndicacao
	err := tx.QueryRow(`
		SELECT id, id_user, id_user_origin, id_game, id_game_bet, id_balance, valor, date_register, status_received_payment
		FROM bonus_indicacao
		WHERE id_game = ?
		LIMIT 0,1
	`, id).Scan(&b.Id, &b.IdUser, &b.IdUserOrigin, &b.IdGame, &b.IdGameBet, &b.IdBalance, &b.Valor, &b.DateRegister, &b.StatusReceivedOPayment)

	if err != nil {
		fmt.Println("GIB 1:" + err.Error())
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	res, err := tx.Exec(`
		INSERT INTO bonus_indicacao (id_user, id_user_origin, id_game, id_game_bet, id_balance, date_register, valor) 
		SELECT u.id_indicador, u.id, b.id_game, b.id, ?, ? ,TRUNCATE(b.bonus_indication_percent_from_tax_bet_win / 100 * ((100 - g.game_profit_percent) / 100 * b.bet_amount_dolar) ,2)
		FROM binary_option_game_bet b
		JOIN usuarios u ON u.id = b.id_usuario AND u.id_indicador > 0
		JOIN binary_option_game g ON g.id = b.id_game
		WHERE b.id_game = ?
		AND b.amount_win_dolar > 0
		AND b.id_balance = ?
		-- AND b.id_trader_follower = 0;

		INSERT INTO bonus_trader (id_user, id_user_origin, id_game, id_game_bet, id_balance, date_register, valor)  
		SELECT b.id_trader_follower, u.id, b.id_game, b.id, ?, ? ,TRUNCATE(b.bonus_trader_percent_from_tax_bet_win / 100 * ((100 - g.game_profit_percent) / 100 * b.bet_amount_dolar) ,2)
		FROM binary_option_game_bet b
		JOIN usuarios u ON u.id = b.id_usuario
		JOIN binary_option_game g ON g.id = b.id_game
		WHERE b.id_game = ?
		AND b.amount_win_dolar > 0
		AND b.id_trader_follower > 0
		AND b.id_balance = ? ;
	`, 16, now, id, 3, 16, now, id, 3)

	lastInsert, _ := res.LastInsertId()

	if err != nil || lastInsert == 0 {
		fmt.Println("GIB 2: ")
		return false
	}

	return true
}
