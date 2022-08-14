package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
)

var statusBonusIndicacao bool
var statusBonusIndicacaofromlosePlayer bool = true

func SetStatusFinishGame(db *sql.DB, tx *sql.Tx, id uint64) (uint64, uint64, int64) {
	var g entities.BinaryOptionGame
	query := `SELECT id, id_moedas_pares, game_id_type_time
	FROM binary_option_game
	WHERE id = ?
	AND game_id_status <= ?
	FOR UPDATE;`

	err := db.QueryRow(query, id, 3).Scan(&g.Id, &g.IdMoedasPares, &g.GameIdTypeTime)

	if err != nil {
		fmt.Println(err)
	}

	query = `UPDATE binary_option_game
		SET game_id_status = 4
		WHERE id = ?`

	_, err = db.Exec(query, id)

	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	}

	bestResultGame := searchBestPriceForEndGame(db, &g, 0)

	saveBestResultGame(db, &g, bestResultGame)

	if len(bestResultGame.ListPlayersWin) > 0 {
		if statusBonusIndicacao {
			generateIndicationBonus(db, g.Id)

		}
	}

	if statusBonusIndicacaofromlosePlayer {
		generateIndicationBonusLosePlayer(db, g.Id)
	}

	releasePaymentWinGame(db, g.Id)
	releasePaymentRefundGame(db, g.Id)

	if statusBonusIndicacao {
		releasePaymentIndicationBonus(db, g.Id)
	}

	if statusBonusIndicacaofromlosePlayer {
		releasePaymentIndicationBonusLosePlayer(db, g.Id)
	}

	tx.Commit()

	return g.Id, g.IdMoedasPares, g.GameIdTypeTime
}
