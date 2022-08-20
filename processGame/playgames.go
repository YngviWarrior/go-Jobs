package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"processgame/entities"

	_ "github.com/go-sql-driver/mysql"
)

func playGames(db *sql.DB) {
	tx, _ := db.Begin()

	now := time.Now().Add(time.Second + 1).Add(time.Hour * 3).Format("2006-01-02 15:04:05")

	query := `SELECT GROUP_CONCAT(id) as id_games, id_moedas_pares, game_id_type_time, 2 as game_status
	FROM binary_option_game
	WHERE game_id_status = ?
	AND game_date_start <= ?
	GROUP BY id_moedas_pares, game_id_type_time;`

	rows, err := tx.Query(query, 1, now)

	if err != nil {
		fmt.Println("PG 1: " + err.Error())
	}

	var gList []*entities.GamesUpdate
	for rows.Next() {
		var g entities.GamesUpdate
		err := rows.Scan(&g.IdGames, &g.IdMoedasPares, &g.GameIdTypeTime, &g.GameIdStatus)

		if err != nil {
			fmt.Println("PG 0: " + err.Error())
			return
		}

		gList = append(gList, &g)
	}

	_, err = tx.Exec("SET SESSION group_concat_max_len = 18446744073709551615;")

	if err != nil {
		fmt.Println("PG 2: " + err.Error())
	}

	var totalGames int
	if len(gList) > 0 {
		for _, g := range gList {
			list := strings.Split(g.IdGames, ",")
			totalGames += len(list)

			if len(list) > 0 {
				var gameList string
				for i, game := range list {
					if i+1 == len(list) {
						gameList += fmt.Sprintf("%v", game)
					} else {
						gameList += fmt.Sprintf("%v,", game)
					}
				}

				query = `UPDATE binary_option_game
					SET game_id_status = ?
					WHERE id IN (` + gameList + `)`

				_, err = tx.Exec(query, 2)

				if err != nil {
					fmt.Println("PG 3: " + err.Error())
					tx.Rollback()

					return
				}
			} else {
				tx.Rollback()
				return
			}
		}
	}

	err = tx.Commit()

	if err != nil {
		fmt.Println("PG 4: " + err.Error())
		return
	}

	push(gList)

	fmt.Printf("%v play changed \n", totalGames)
}
