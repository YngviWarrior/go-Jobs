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
	now := time.Now().Add(time.Second * 2).Add(time.Hour * 3).Format("2006-01-02 15:04:05")

	query := `SELECT id
	FROM binary_option_game
	WHERE game_id_status = ?
	AND game_date_finish <= ?`

	rows, err := db.Query(query, 3, now)

	if err != nil {
		fmt.Println("FG 1: " + err.Error())
		return
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

	if len(list) > 0 {
		var toCache []entities.GamesUpdate
		for i, v := range list {
			if i%50 == 0 {
				time.Sleep(time.Second * 1)
			}

			go func(v *entities.BinaryOptionGame) {
				// fmt.Println(v.Id)

				idGame, idSymbol, idTime := finishgame.SetStatusFinishGame(db, v.Id)
				if idGame >= 0 {
					var game entities.GamesUpdate
					game.Id = uint64(idGame)
					game.IdMoedasPares = idSymbol
					game.GameIdTypeTime = idTime

					toCache = append(toCache, game)
				}
			}(v)
		}

		push(toCache)
	}

	fmt.Printf("%v finish changed \n", len(list))
	defer db.Close()
}
