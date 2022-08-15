package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
	"time"
)

func generateIndicationBonusLosePlayer(db *sql.DB, id uint64) bool {
	var b entities.BonusIndicacao
	err := db.QueryRow(`
		SELECT id, id_user, id_user_origin, id_game, id_game_bet, id_balance, valor, date_register, status_received_payment
		FROM bonus_indicacao
		WHERE id_game = ?
		LIMIT 0,1
	`, id).Scan(&b.Id, &b.IdUser, &b.IdUserOrigin, &b.IdGame, &b.IdGameBet, &b.IdBalance, &b.Valor, &b.DateRegister, &b.StatusReceivedOPayment)

	if err != nil {
		fmt.Println("GIBLP 1: " + err.Error())
		// return false
	}

	activePeriod := 1
	now := time.Now().Format("2006-01-02 15:04:05")

	res, err := db.Exec(`
		INSERT INTO bonus_indicacao (id_user, id_user_origin, id_game, id_game_bet, id_balance, date_register, valor, status_received_payment) 
		SELECT u.id_indicador, u.id, b.id_game, b.id, ?, ?
			,IF(b.amount_win_dolar = 0,
				IF(
				TRUNCATE(b.bet_amount_dolar / 100 * COALESCE(indicador.bonus_indication_percent,0) ,2) > (u.total_deposit_balance_play - u.total_lose_balance_play), 
				(u.total_deposit_balance_play - u.total_lose_balance_play),
				TRUNCATE(b.bet_amount_dolar / 100 * COALESCE(indicador.bonus_indication_percent,0) ,2)
				), (b.amount_win_dolar * -1)  
			), ?			
		FROM binary_option_game g
		JOIN binary_option_game_bet b ON g.id = b.id_game
		JOIN usuarios u ON u.id = b.id_usuario AND u.id_indicador > 0
		JOIN usuarios indicador 
			ON indicador.id = u.id_indicador 
			AND (u.total_deposit_balance_play - u.total_lose_balance_play) > 0
			
		WHERE b.id_game = ?
		AND b.id_balance = ?
		`, 24, now, activePeriod, id, 3)

	inserId, _ := res.LastInsertId()

	if err != nil || inserId == 0 {
		fmt.Println("GIBLP 2: ")
	}

	return true
}
