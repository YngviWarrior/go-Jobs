package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
)

var statusBonusIndicacao bool
var statusBonusIndicacaofromlosePlayer bool = true

func SetStatusFinishGame(db *sql.DB, id uint64) (uint64, uint64, int64) {
	tx, _ := db.Begin()

	var g entities.BinaryOptionGame
	query := `SELECT id, id_moedas_pares, game_id_type_time, game_profit_percent
	FROM binary_option_game
	WHERE id = ?
	AND game_id_status <= ?
	FOR UPDATE;`

	err := tx.QueryRow(query, id, 3).Scan(&g.Id, &g.IdMoedasPares, &g.GameIdTypeTime, &g.GameProfitPercent)

	if err != nil {
		fmt.Println("STFG 1: " + err.Error())
		return 0, 0, 0
	}

	query = `UPDATE binary_option_game
		SET game_id_status = 4
		WHERE id = ?`

	_, err = tx.Exec(query, id)

	if err != nil {
		fmt.Println("STFG 2: " + err.Error())
		tx.Rollback()
		return 0, 0, 0
	}

	bestResultGame, ok := searchBestPriceForEndGame(tx, &g, 100)

	if !ok {
		tx.Rollback()
		return 0, 0, 0
	}

	if !saveBestResultGame(tx, &g, bestResultGame) {
		fmt.Println("saveBestResultGame")
		tx.Rollback()
		return 0, 0, 0
	}

	//Não testado
	if len(bestResultGame.ListPlayersWin) > 0 {
		if statusBonusIndicacao {
			if !generateIndicationBonus(tx, g.Id) {
				fmt.Println("generateIndicationBonus")
				tx.Rollback()
				return 0, 0, 0
			}

		}
	}

	if statusBonusIndicacaofromlosePlayer {
		// Query estáva errada, isso tá sendo usado ??? Não encontra nada nas tabelas, mt menos insere? é pra funcionar ?
		if !generateIndicationBonusLosePlayer(tx, g.Id) {
			fmt.Println("generateIndicationBonusLosePlayer")
			tx.Rollback()
			return 0, 0, 0
		}
	}

	if len(bestResultGame.ListPlayersWin) > 0 {
		// Será que aquele group by em RPWG 3 é necessário ?
		if !releasePaymentWinGame(tx, g.Id) {
			fmt.Println("releasePaymentWinGame")
			tx.Rollback()
			return 0, 0, 0
		}
	} else if !releasePaymentRefundGame(tx, g.Id) { // Será que aquele group by em RPRG 2 é necessário ?
		fmt.Println("releasePaymentRefundGame")
		tx.Rollback()
		return 0, 0, 0
	}

	//Não testado
	if statusBonusIndicacao {
		if !releasePaymentIndicationBonus(tx, g.Id) {
			fmt.Println("releasePaymentIndicationBonus")
			tx.Rollback()
			return 0, 0, 0
		}
	}

	if statusBonusIndicacaofromlosePlayer {
		if !releasePaymentIndicationBonusLosePlayer(tx, g.Id) {
			fmt.Println("releasePaymentIndicationBonusLosePlayer")
			tx.Rollback()
			return 0, 0, 0
		}
	}

	tx.Commit()

	return g.Id, g.IdMoedasPares, g.GameIdTypeTime
}
