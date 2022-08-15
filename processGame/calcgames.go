package main

import (
	"database/sql"
	"fmt"
	"log"
	"processgame/entities"
	"time"
)

func calcGames(db *sql.DB) {
	tx, _ := db.Begin()

	_, err := tx.Exec("SET SESSION group_concat_max_len = 18446744073709551615;")

	if err != nil {
		fmt.Println("CG 1: " + err.Error())
		tx.Rollback()
		return
	}

	now := time.Now().Add(time.Second + 1).Add(time.Hour + 3).Format("2006-01-02 15:04:05")

	query := `SELECT id, id_moedas_pares, game_id_type_time
	FROM binary_option_game
	WHERE game_id_status = ?
	AND game_date_process <= ?
	FOR UPDATE;`

	rows, err := tx.Query(query, 2, now)

	if err != nil {
		fmt.Println("CG 1: " + err.Error())
		tx.Rollback()
		return
	}

	var list []*entities.GamesUpdate

	for rows.Next() {
		var g entities.GamesUpdate
		err = rows.Scan(&g.Id, &g.IdMoedasPares, &g.GameIdTypeTime)

		if err != nil {
			log.Println("fetch err: ", err)
		}

		list = append(list, &g)
	}

	if len(list) > 0 {
		var toCache []entities.GamesUpdate
		var gameList string

		for i, v := range list {
			if i+1 == len(list) {
				gameList += fmt.Sprintf("%v", v.Id)
			} else {
				gameList += fmt.Sprintf("%v,", v.Id)
			}

			var game entities.GamesUpdate
			game.Id = v.Id
			game.IdMoedasPares = v.IdMoedasPares
			game.GameIdTypeTime = v.GameIdTypeTime

			toCache = append(toCache, game)
		}

		query = `UPDATE binary_option_game
		SET game_id_status = ?
		WHERE id IN (` + gameList + `)`

		_, err = tx.Exec(query, 3)

		if err != nil {
			fmt.Println("CG 2: " + err.Error())
			tx.Rollback()
			return
		}

		tx.Commit()

		push(toCache)
	}

	fmt.Printf("%v calc changed \n", len(list))
}
