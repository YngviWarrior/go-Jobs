package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
)

func releasePaymentIndicationBonusLosePlayer(tx *sql.Tx, id uint64) bool {
	var b entities.BonusIndicacao
	err := tx.QueryRow(`
		SELECT id, id_user, id_user_origin, id_game, id_game_bet, id_balance, valor, date_register, status_received_payment
		FROM bonus_indicacao
		WHERE id_game = ? AND status_received_payment = 1
		LIMIT 0,1
	`, id).Scan(&b.Id, &b.IdUser, &b.IdUserOrigin, &b.IdGame, &b.IdGameBet, &b.IdBalance, &b.Valor, &b.DateRegister, &b.StatusReceivedOPayment)

	if err != nil {
		fmt.Println("RPIBLP 1: That's OK." + err.Error())
	}

	_, err = tx.Exec(`
		CREATE TEMPORARY TABLE bonus_indication_process(
			id bigint(20) NOT NULL,
			id_type tinyint(1) NOT NULL,
			id_user bigint(20) NOT NULL,
			valor decimal(60,8) NOT NULL,
			KEY ix_bonus_indication_process_id_user (id_user)
		) ENGINE=InnoDB;
	`)

	if err != nil {
		fmt.Println("RPIBLP 2:" + err.Error())
	}

	res, err := tx.Exec(`
		INSERT INTO bonus_indication_process VALUES(1000, 1, 6, 1);
	`)

	if err != nil {
		fmt.Println("RPIBLP TEST RPIBLP 2.1: " + err.Error())
		return false
	}

	affcRows, _ := res.RowsAffected()

	if affcRows == 0 {
		fmt.Println("RPIBLP TEST RPIBLP 2.2: " + err.Error())
		return false
	}

	var bonus_indicacao_process float64
	err = tx.QueryRow(`SELECT SUM(valor) as valor
		FROM bonus_indication_process
		WHERE id_type = 1`).Scan(&bonus_indicacao_process)

	if err != nil {
		fmt.Println("RPIBLP TEST RPIBLP 2.3: " + err.Error())
		return false
	}

	if bonus_indicacao_process == 0 {
		fmt.Println("RPIBLP TEST RPIBLP 2.4: Table nos created.")
		return false
	}

	rows, err := tx.Query(`
		SELECT *
			FROM (
				SELECT b.id, ? as id_type, b.id_user, b.valor
				FROM bonus_indicacao b
				WHERE b.id_game = ?
				AND b.status_received_payment = 0
			)
		as t;
	`, 1, id)

	if err != nil {
		fmt.Println("RPIBLP 3: " + err.Error())
		// return false
	}

	var bi []*entities.BonusIndicacao
	for rows.Next() {
		var b entities.BonusIndicacao

		err := rows.Scan(&b.Id, &b.IdType, &b.IdUser, &b.Valor)

		if err != nil {
			fmt.Println("RPIBLP 4: " + err.Error())
		}

		if b.Id != 0 {
			bi = append(bi, &b)
		}
	}

	if len(bi) == 0 {
		fmt.Println("RPIBLP 5: No Bonus Indication.")
	}

	_, err = tx.Exec(`
		SET @totalBonusIndication := COALESCE(
			(SELECT SUM(b.valor)
			FROM bonus_indication_process b
			WHERE b.id_type = ?
		),0)
	`, 1)

	if err != nil {
		fmt.Println("RPIBLP 6:" + err.Error())
		return false
	}

	var totalBonusIndication float64
	err = tx.QueryRow(`
		SELECT @totalBonusIndication;
	`).Scan(&totalBonusIndication)

	if err != nil {
		fmt.Println("RPIBLP Select 1: " + err.Error())
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
		fmt.Println("RPIBLP 8:" + err.Error())
		return false
	}

	var totalBonusTrader float64
	err = tx.QueryRow(`
		SELECT @totalBonusTrader;
	`).Scan(&totalBonusTrader)

	if err != nil {
		fmt.Println("RPIBLP Select 2 :" + err.Error())
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
		fmt.Println("RPIBLP 9:" + err.Error())
		return false
	}

	var campanyLeft float64
	err = tx.QueryRow(`
		SELECT @campanyLeft;
	`).Scan(&campanyLeft)

	if err != nil {
		fmt.Println("RPIBLP Select 3: " + err.Error())
		return false
	}

	res, err = tx.Exec(`UPDATE binary_option_game g
	SET g.bonus_trader_total_amount_dolar = TRUNCATE(?,8)
		,g.bonus_indication_total_amount_dolar = TRUNCATE(?,8)
		,g.company_tax_amount_dolar_from_game_tax = TRUNCATE(?,8)
	WHERE  g.id = ?;`, totalBonusTrader, totalBonusIndication, campanyLeft, id)

	if err != nil {
		fmt.Println("RPIBLP 11:" + err.Error())
		return false
	}

	affcRows, _ = res.RowsAffected()
	if affcRows == 0 {
		fmt.Println("RPIBLP 12: ")
	}

	rows, err = tx.Query(`
		SELECT b.id, b.id_user, b.valor
		FROM bonus_indicacao b
		WHERE b.id_game = ? AND b.status_received_payment = 0
	`, id)

	if err != nil {
		fmt.Println("RPIBLP 13:" + err.Error())
		// return false
	}

	var bonusIndicationList []*entities.BonusIndicacao
	for rows.Next() {
		var bonus entities.BonusIndicacao
		err := rows.Scan(&bonus.Id, &bonus.IdUser, &bonus.Valor)

		if err != nil {
			fmt.Println("RPIBLP 14:" + err.Error())
			return false
		}

		if b.Id != 0 {
			bonusIndicationList = append(bonusIndicationList, &b)
		}
	}

	if len(bonusIndicationList) > 0 {
		for _, v := range bonusIndicationList {
			modifyBalance(tx, v.IdUser, 24, 3, v.Valor, v.Id, true)
		}
	} else {
		fmt.Println("RPIBLP 15: No Bonus Indication")
		// return false
	}

	rows, err = tx.Query(`
		SELECT b.id, b.id_user, b.valor
		FROM bonus_trader b
		WHERE b.id_game = ? AND b.status_received_payment = 0
	`, id)

	if err != nil {
		fmt.Println("RPIBLP 16:" + err.Error())
		return false
	}

	var bonusTraderList []*entities.BonusTrader
	for rows.Next() {
		var bonus entities.BonusTrader

		err := rows.Scan(&bonus.Id, &bonus.IdUser, &bonus.Valor)

		if err != nil {
			fmt.Println("RPIBLP 17:" + err.Error())
			return false
		}

		if b.Id != 0 {
			bonusTraderList = append(bonusTraderList, &bonus)
		}
	}

	if len(bonusTraderList) > 0 && bonusTraderList[0].Id != 0 {
		for _, v := range bonusTraderList {
			modifyBalance(tx, v.IdUser, 16, 9, v.Valor, v.Id, false)
		}
	} else {
		fmt.Println("RPIBLP 18: No Bonus Trader")
	}

	res, err = tx.Exec(`
		UPDATE bonus_indicacao b
		JOIN bonus_indication_process bp ON bp.id = b.id AND bp.id_type = ?
		SET b.status_received_payment = 1;
	`, 1)

	if err != nil {
		fmt.Println("RPIBLP 19:" + err.Error())
		return false
	}

	affcRows, _ = res.RowsAffected()
	if affcRows == 0 {
		fmt.Println("RPIBLP 20: ")
		// return false
	}

	res, err = tx.Exec(`
		UPDATE bonus_trader b
		JOIN bonus_indication_process bp ON bp.id = b.id AND bp.id_type = ?
		SET b.status_received_payment = 1;
	`, 1)

	if err != nil {
		fmt.Println("RPIBLP 21:" + err.Error())
		return false
	}

	affcRows, _ = res.RowsAffected()
	if affcRows == 0 {
		fmt.Println("RPIBLP 22: ")
		// return false
	}

	_, err = tx.Exec(`
		DROP TEMPORARY TABLE IF EXISTS bonus_indication_process;
	`)

	if err != nil {
		fmt.Println("RPIBLP 23:" + err.Error())
		return false
	}

	return true
}
