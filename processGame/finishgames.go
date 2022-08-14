package main

import (
	"database/sql"
	"fmt"
	"log"

	"processgame/entities"
	finishgame "processgame/finishGame"
	"time"
)

func finishGames(db *sql.DB) {
	now := time.Now().Add(time.Second + 2).Format("2006-01-02 15:04:05")

	query := `SELECT id
	FROM binary_option_game
	WHERE game_id_status = ?
	AND game_date_finish <= ?`

	rows, err := db.Query(query, 3, now)

	if err != nil {
		fmt.Println(err)
	}

	var list []*entities.BinaryOptionGame

	for rows.Next() {
		var g entities.BinaryOptionGame
		err = rows.Scan(&g.Id)

		if err != nil {
			log.Println("fetch err: ", err)
		}

		list = append(list, &g)
	}

	fmt.Printf("%v finish changed \n", len(list))

	if len(list) > 0 {
		// var toCache []entities.GamesUpdate
		for _, v := range list {
			tx, _ := db.Begin()
			fmt.Println(tx)

			idGame, idSymbol, idTime := finishgame.SetStatusFinishGame(db, tx, v.Id)
			fmt.Println(idGame, idSymbol, idTime)
			// var game entities.GamesUpdate
			// game.Id = idGame
			// game.IdMoedasPares = idSymbol
			// game.GameIdTypeTime = idTime

			// toCache = append(toCache, game)
		}

		// push(toCache)
	}

	// defer db.Close()
}
