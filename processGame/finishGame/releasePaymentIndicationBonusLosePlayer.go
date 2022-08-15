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
		fmt.Println("RPIBLP 1:" + err.Error())
	}

	query := `
	CREATE TEMPORARY TABLE bonus_indication_process
	(
		id bigint(20) NOT NULL,
		id_type tinyint(1) NOT NULL,
		id_user bigint(20) NOT NULL,
		valor decimal(16,8) NOT NULL,
		KEY ix_bonus_indication_process_id_user (id_user)
	) ENGINE=InnoDB;
		
	SELECT *
	FROM (
		SELECT b.id, ? as id_type, b.id_user, b.valor
		FROM bonus_indicacao b
		WHERE b.id_game = ?
		AND b.status_received_payment = 0
	) as t;


	SET @totalBonusIndication := COALESCE(
		(SELECT SUM(b.valor)
		FROM bonus_indication_process b
		WHERE b.id_type = ?
	),0);

	SET @totalBonusTrader := COALESCE(
		(SELECT SUM(b.valor)
		FROM bonus_indication_process b
		WHERE b.id_type = ?
	),0);

	SET @campanyLeft := (COALESCE(
		(SELECT (SUM(b.bet_amount_dolar) - SUM(IF(b.amount_win_dolar > 0,b.bet_amount_dolar + b.amount_win_dolar,0)))
		FROM binary_option_game_bet b
		WHERE b.id_game = ? AND b.id_balance = ?
	),0) - @totalBonusTrader - @totalBonusIndication);

	UPDATE binary_option_game g
	SET
		g.bonus_trader_total_amount_dolar = TRUNCATE(@totalBonusTrader,8)
		,g.bonus_indication_total_amount_dolar = TRUNCATE(@totalBonusIndication,8)
		,g.company_tax_amount_dolar_from_game_tax = TRUNCATE(@campanyLeft,8)
	WHERE  g.id = ?; `

	_, err = tx.Query(query, 1, id, 1, 2, id, 3, id)

	if err != nil {
		fmt.Println("RPIBLP 2:" + err.Error())
		return false
	}

	rows, err := tx.Query(`
		SELECT b.id, b.id_user, b.valor
		FROM bonus_indicacao b
		WHERE b.id_game = ? AND b.status_received_payment = 0
	`, id)

	if err != nil {
		fmt.Println("RPIBLP 3:" + err.Error())
		return false
	}

	var bonusIndicationList []*entities.BonusIndicacao
	for rows.Next() {
		var bonus entities.BonusIndicacao

		err := rows.Scan(&bonus.Id, &bonus.IdUser, &bonus.Valor)

		if err != nil {
			fmt.Println("RPIBLP 4:" + err.Error())
			return false
		}

		bonusIndicationList = append(bonusIndicationList, &b)
	}

	if len(bonusIndicationList) > 0 {
		for _, v := range bonusIndicationList {
			modifyBalance(tx, v.IdUser, 24, 3, v.Valor, v.Id, true)
		}
	} else {
		fmt.Println("RPIBLP 5: No Bonus" + err.Error())
		return false
	}

	rows, err = tx.Query(`
		SELECT b.id, b.id_user, b.valor
		FROM bonus_trader b
		WHERE b.id_game = ? AND b.status_received_payment = 0
	`, id)

	if err != nil {
		fmt.Println("RPIBLP 6:" + err.Error())
		return false
	}

	var bonusTraderList []*entities.BonusTrader
	for rows.Next() {
		var bonus entities.BonusTrader

		err := rows.Scan(&bonus.Id, &bonus.IdUser, &bonus.Valor)

		if err != nil {
			fmt.Println("RPIBLP 7:" + err.Error())
			return false
		}

		bonusTraderList = append(bonusTraderList, &bonus)
	}

	if len(bonusTraderList) > 0 {
		for _, v := range bonusTraderList {
			modifyBalance(tx, v.IdUser, 16, 9, v.Valor, v.Id, false)
		}
	} else {
		fmt.Println("RPIBLP 8: No Bonus Trader")
		return false
	}

	res, err := tx.Exec(`
		UPDATE bonus_indicacao b
		JOIN bonus_indication_process bp ON bp.id = b.id AND bp.id_type = ?
		SET b.status_received_payment = 1;
		UPDATE bonus_trader b
		JOIN bonus_indication_process bp ON bp.id = b.id AND bp.id_type = ?
		SET b.status_received_payment = 1;

		DROP TEMPORARY TABLE IF EXISTS bonus_indication_process;
	`, 1, 2)

	affcRows, _ := res.RowsAffected()

	if err != nil || affcRows == 0 {
		fmt.Println("RPIBLP 9:" + err.Error())
		return false
	}

	return true
}
