package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
)

var statusBonusIndicacao bool
var statusBonusIndicacaofromlosePlayer bool = true

func SetStatusFinishGame(db *sql.DB, id string) {
	tx, _ := db.Begin()

	var g entities.BinaryOptionGame
	query := `SELECT id, id_moedas_pares, game_id_type_time, game_profit_percent, 4 as game_status
				FROM binary_option_game
				WHERE id = ?`

	err := tx.QueryRow(query, id).Scan(&g.Id, &g.IdMoedasPares, &g.GameIdTypeTime, &g.GameProfitPercent, &g.GameIdStatus)

	if err != nil {
		fmt.Println("STFG 1: " + err.Error())
		tx.Rollback()
		return
	}

	query = `UPDATE binary_option_game
		SET game_id_status = 4
		WHERE id = ?`

	_, err = tx.Exec(query, g.Id)

	if err != nil {
		fmt.Println("STFG 2: " + err.Error())
		tx.Rollback()
		return
	}

	bestResultGame, ok := searchBestPriceForEndGame(tx, &g, 100)

	if ok == 0 {
		tx.Rollback()
		return
	} else if ok == 1 {
		tx.Commit()
		return
	}

	if len(bestResultGame.ListPlayersWin) > 0 || len(bestResultGame.ListPlayersEqual) > 0 {
		if !saveBestResultGame(tx, &g, bestResultGame) {
			fmt.Println("saveBestResultGame")
			tx.Rollback()
			return
		}

		// //Não testado
		// if len(bestResultGame.ListPlayersWin) > 0 {
		// 	if statusBonusIndicacao {
		// 		if !generateIndicationBonus(tx, g.Id) {
		// 			fmt.Println("generateIndicationBonus")
		// 			tx.Rollback()
		// 			return
		// 		}
		// 	}
		// }

		// if statusBonusIndicacaofromlosePlayer {
		// 	// Query estáva errada, isso tá sendo usado ??? Não encontra nada nas tabelas, mt menos insere? é pra funcionar ?
		// 	if !generateIndicationBonusLosePlayer(tx, g.Id) {
		// 		fmt.Println("generateIndicationBonusLosePlayer")
		// 		tx.Rollback()
		// 		return
		// 	}
		// }

		if !releasePaymentWinGame(tx, g.Id) { // Será que aquele group by em RPRG 3 é necessário ?
			fmt.Println("releasePaymentWinGame")
			tx.Rollback()
			return
		}

		if !releasePaymentRefundGame(tx, g.Id) { // Será que aquele group by em RPRG 2 é necessário ?
			fmt.Println("releasePaymentRefundGame")
			tx.Rollback()
			return
		}

		// //Não testado
		// if statusBonusIndicacao {
		// 	if !releasePaymentIndicationBonus(tx, g.Id) {
		// 		fmt.Println("releasePaymentIndicationBonus")
		// 		tx.Rollback()
		// 		return
		// 	}
		// }

		// if statusBonusIndicacaofromlosePlayer {
		// 	if !releasePaymentIndicationBonusLosePlayer(tx, g.Id) {
		// 		fmt.Println("releasePaymentIndicationBonusLosePlayer")
		// 		tx.Rollback()
		// 		return
		// 	}
		// }
	}

	tx.Commit()
}
