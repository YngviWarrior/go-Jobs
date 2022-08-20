package main

import (
	"database/sql"
	"fmt"
	"strings"

	"processgame/entities"
	finishgame "processgame/finishGame"
	"time"
)

func finishGames(db *sql.DB) {
	now := time.Now().Add(time.Second * 2).Add(time.Hour * 3).Format("2006-01-02 15:04:05")

	query := `SELECT GROUP_CONCAT(id) as id_games, id_moedas_pares, game_id_type_time, 4 as game_status
	FROM binary_option_game
	WHERE game_id_status = ?
	AND game_date_finish <= ?
	GROUP BY id_moedas_pares, game_id_type_time;`

	rows, err := db.Query(query, 3, now)

	if err != nil {
		fmt.Println("FG 1: " + err.Error())
	}

	var gList []*entities.GamesUpdate
	for rows.Next() {
		var g entities.GamesUpdate
		err := rows.Scan(&g.IdGames, &g.IdMoedasPares, &g.GameIdTypeTime, &g.GameIdStatus)

		if err != nil {
			fmt.Println("FG 0: " + err.Error())
			return
		}

		gList = append(gList, &g)
	}

	var totalGames int
	if len(gList) > 0 {
		for _, g := range gList {
			go func(g *entities.GamesUpdate) {
				list := strings.Split(g.IdGames, ",")
				totalGames += len(list)

				if len(list) > 0 {
					for _, game := range list {
						// if i%50 == 0 {
						// 	time.Sleep(time.Second * 1)
						// }

						finishgame.SetStatusFinishGame(db, game)
					}
				}
			}(g)
		}
	}

	push(gList)

	fmt.Printf("%v finish changed \n", totalGames)
}
