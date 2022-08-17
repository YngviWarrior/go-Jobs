package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
)

func releasePaymentIndicationBonus(tx *sql.Tx, id uint64) bool {
	var b entities.BonusIndicacao
	err := tx.QueryRow(`
		SELECT id, id_user, id_user_origin, id_game, id_game_bet, id_balance, valor, date_register, status_received_payment
		FROM bonus_indicacao
		WHERE id_game = ? AND status_received_payment = 1
		LIMIT 0,1
	`, id).Scan(&b.Id, &b.IdUser, &b.IdUserOrigin, &b.IdGame, &b.IdGameBet, &b.IdBalance, &b.Valor, &b.DateRegister, &b.StatusReceivedOPayment)

	if err != nil {
		fmt.Println("RPIB 1: " + err.Error())
	}

	_, err = tx.Exec(`CREATE TEMPORARY TABLE bonus_indication_process(
		id bigint(20) NOT NULL,
		id_type tinyint(1) NOT NULL,
		id_user bigint(20) NOT NULL,
		valor decimal(16,8) NOT NULL,
		KEY ix_bonus_indication_process_id_user (id_user)
	) ENGINE=InnoDB;`)

	if err != nil {
		fmt.Println("RPIB 2: " + err.Error())
	}

	res, err := tx.Exec(`
		INSERT INTO bonus_indication_process VALUES(1000, 1, 6, 1);
	`)

	if err != nil {
		fmt.Println("RPIB TEST RPIBLP 2.1: " + err.Error())
		return false
	}

	affcRows, _ := res.RowsAffected()

	if affcRows == 0 {
		fmt.Println("RPIB TEST RPIBLP 2.2: " + err.Error())
		return false
	}

	var bonus_indicacao_process float64
	err = tx.QueryRow(`SELECT SUM(valor) as valor
		FROM bonus_indication_process
		WHERE id_type = 1`).Scan(&bonus_indicacao_process)

	if err != nil {
		fmt.Println("RPIB TEST RPIBLP 2.3: " + err.Error())
		return false
	}

	if bonus_indicacao_process == 0 {
		fmt.Println("RPIB TEST RPIBLP 2.4: Table nos created.")
		return false
	}

	rows, err := tx.Query(`
		SELECT *
		FROM (
			(
				SELECT b.id, ? as id_type, b.id_user, b.valor
				FROM bonus_indicacao b
				WHERE b.id_game = ?
				AND b.status_received_payment = 0
			)
			UNION ALL
			(
				SELECT b.id, ? as id_type, b.id_user, b.valor
				FROM bonus_trader b
				WHERE b.id_game = ?
				AND b.status_received_payment = 0
			)
		) as t;
	`, 1, id, 1, id)

	if err != nil {
		fmt.Println("RPIB 3: " + err.Error())
	}

	var bonusList []*entities.BonusIndicacao
	for rows.Next() {
		var b entities.BonusIndicacao

		err := rows.Scan(&b.Id, &b.IdType, &b.IdUser, &b.Valor)

		if err != nil {
			fmt.Println("RPIB 4: " + err.Error())
			return false
		}

		if b.Id != 0 {
			bonusList = append(bonusList, &b)
		}
	}

	if len(bonusList) == 0 {
		fmt.Println("RPIB 5: No Bonus Indication.")
	}

	res, err = tx.Exec(`
		UPDATE usuarios
		JOIN (
			SELECT b.id_user, SUM(b.valor) as amount
			FROM bonus_indication_process b
			GROUP BY b.id_user
		) as t ON t.id_user = usuarios.id
		SET usuarios.total_bonus = usuarios.total_bonus + t.amount;
	`)

	if err != nil {
		fmt.Println("RPIB 6: " + err.Error())
		return false
	}

	affcRows, _ = res.RowsAffected()

	if affcRows == 0 {
		fmt.Println("RPIB 7: Total Bonus not updated.")
		// return false
	}

	res, err = tx.Exec(`
		INSERT INTO usuarios_total_bonus (id_usuario, id_bonus, pontos)
		SELECT t.id_user, t.id_bonus, t.amount
		FROM (
			SELECT b.id_user, b.id_type as id_bonus, SUM(b.valor) as amount
			FROM bonus_indication_process b
			GROUP BY b.id_user, b.id_type
		) as t ON DUPLICATE KEY UPDATE pontos = pontos + t.amount;
	`)

	if err != nil {
		fmt.Println("RPIB 8: " + err.Error())
		return false
	}

	lastInsert, _ := res.LastInsertId()

	if lastInsert == 0 {
		fmt.Println("RPIB 9: UserTotalBonus not inserted.")
		// return false
	}

	_, err = tx.Exec(`
		SET @totalBonusIndication := COALESCE(
			(SELECT SUM(b.valor)
			FROM bonus_indication_process b
			WHERE b.id_type = ?
		),0)
	`, 1)

	if err != nil {
		fmt.Println("RPIB 10: " + err.Error())
		return false
	}

	var totalBonusIndication float64
	err = tx.QueryRow(`
		SELECT @totalBonusIndication;
	`).Scan(&totalBonusIndication)

	if err != nil {
		fmt.Println("RPIB Select 1: " + err.Error())
		return false
	}

	_, err = tx.Exec(`
		SET @totalBonusTrader := COALESCE(
				(SELECT SUM(b.valor)
				FROM bonus_indication_process b
				WHERE b.id_type = ?
			)
		,0)
	`, 1)

	if err != nil {
		fmt.Println("RPIB 11: " + err.Error())
		return false
	}

	var totalBonusTrader float64
	err = tx.QueryRow(`
		SELECT @totalBonusTrader;
	`).Scan(&totalBonusTrader)

	if err != nil {
		fmt.Println("RPIB Select 2 :" + err.Error())
		return false
	}

	_, err = tx.Exec(`
		SET @campanyLeft := ABS(COALESCE(
			(SELECT (SUM(b.bet_amount_dolar) - SUM(IF(b.amount_win_dolar > 0,b.bet_amount_dolar + b.amount_win_dolar,0)))
			FROM binary_option_game_bet b
			WHERE b.id_game = ? AND b.id_balance = ?
		),0) - ? - ?);
	`, id, 20, totalBonusTrader, totalBonusIndication) // 3

	if err != nil {
		fmt.Println("RPIB 13: " + err.Error())
		return false
	}

	var campanyLeft float64
	err = tx.QueryRow(`
		SELECT @campanyLeft;
	`).Scan(&campanyLeft)

	if err != nil {
		fmt.Println("RPIB Select 3: " + err.Error())
		return false
	}

	res, err = tx.Exec(`UPDATE binary_option_game g
	SET g.bonus_trader_total_amount_dolar = TRUNCATE(?,8)
		,g.bonus_indication_total_amount_dolar = TRUNCATE(?,8)
		,g.company_tax_amount_dolar_from_game_tax = TRUNCATE(?,8)
	WHERE  g.id = ?;`, totalBonusTrader, totalBonusIndication, campanyLeft, id)

	if err != nil {
		fmt.Println("RPIB 14: " + err.Error())
		return false
	}

	affcRows, _ = res.RowsAffected()
	if affcRows == 0 {
		fmt.Println("RPIB 15: ")
	}

	if err != nil {
		fmt.Println("RPIB 16: " + err.Error())
		return false
	}

	rows, err = tx.Query(`
		SELECT b.id, b.id_user, b.valor
		FROM bonus_indicacao b
		WHERE b.id_game = ? AND b.status_received_payment = 0
	`, id)

	if err != nil {
		fmt.Println("RPIB 17: " + err.Error())
	}

	var bonusIndicationList []*entities.BonusIndicacao
	for rows.Next() {
		var bonus entities.BonusIndicacao

		err := rows.Scan(&bonus.Id, &bonus.IdUser, &bonus.Valor)

		if err != nil {
			fmt.Println("RPIB 18: " + err.Error())
			return false
		}
		if bonus.Id != 0 {
			bonusIndicationList = append(bonusIndicationList, &bonus)
		}
	}

	if len(bonusIndicationList) > 0 {
		for _, v := range bonusIndicationList {
			modifyBalance(tx, v.IdUser, 16, 3, v.Valor, v.Id, false)
		}
	} else {
		fmt.Println("RPIB 19: No Bets.")
	}

	rows, err = tx.Query(`
		SELECT b.id, b.id_user, b.valor
		FROM bonus_trader b
		WHERE b.id_game = ? AND b.status_received_payment = 0
	`, id)

	if err != nil {
		fmt.Println("RPIB 20: " + err.Error())
	}

	var bonusTraderList []*entities.BonusTrader
	for rows.Next() {
		var bonus entities.BonusTrader

		err := rows.Scan(&bonus.Id, &bonus.IdUser, &bonus.Valor)

		if err != nil {
			fmt.Println("RPIB 21: " + err.Error())
			return false
		}

		if bonus.Id != 0 {
			bonusTraderList = append(bonusTraderList, &bonus)
		}
	}

	if len(bonusTraderList) > 0 {
		for _, v := range bonusTraderList {
			modifyBalance(tx, v.IdUser, 16, 9, v.Valor, v.Id, false)
		}
	} else {
		fmt.Println("RPIB 22: No Bonus Trader.")
	}

	_, err = tx.Exec(`
		UPDATE bonus_indicacao b
		JOIN bonus_indication_process bp ON bp.id = b.id AND bp.id_type = ?
		SET b.status_received_payment = 1;
	`, 1)

	if err != nil {
		fmt.Println("RPIB 23: " + err.Error())
		return false
	}

	_, err = tx.Exec(`
		UPDATE bonus_trader b
		JOIN bonus_indication_process bp ON bp.id = b.id AND bp.id_type = ?
		SET b.status_received_payment = 1;
	`, 2)

	if err != nil {
		fmt.Println("RPIB 24: " + err.Error())
		return false
	}

	_, err = tx.Exec(`
		DROP TEMPORARY TABLE IF EXISTS bonus_indication_process;
	`)

	if err != nil {
		fmt.Println("RPIB 25: " + err.Error())
		return false
	}

	return true
}
