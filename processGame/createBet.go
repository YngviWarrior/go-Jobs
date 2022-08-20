package main

import (
	"database/sql"
	"fmt"
	"time"
)

type usuarios struct {
	Id int64 `json:"id"`
}

func createBet(db *sql.DB) {
	rows, err := db.Query(`SELECT id FROM usuarios`)

	if err != nil {
		fmt.Println(err)
	}

	var uL []usuarios
	for rows.Next() {
		var u usuarios

		err := rows.Scan(&u.Id)

		if err != nil {
			fmt.Println(err)
		}

		uL = append(uL, u)
	}

	for _, u := range uL {
		_, err := db.Exec(`
		INSERT INTO binary_option_game_bet(hash_id, id_game, id_usuario, id_choice, id_balance, bet_amount_dolar, amount_win_dolar,
		price_amount_selected, status_received_win_payment, id_trader_follower, bot_use_status, date_register, 
		bonus_trader_percent_from_tax_bet_win, bonus_indication_percent_from_tax_bet_win, status_received_refund_payment, 
		refund, deleted) VALUES 
		(1044148127, ?, ?, 2, 20, 2.00000000, 0, 24123.79000000, 0, 0, 0, ?, 0.00000000, 0.00000000, 0, 0, 0);
		`, 36619, u.Id, time.Now().Format("2006-01-02 15:04:05"))

		if err != nil {
			fmt.Println(err)
		}

		// li, _ := res.LastInsertId()

		// if li == 0 {
		// 	fmt.Println("Não inseriu")
		// }

	}
}
