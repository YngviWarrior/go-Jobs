package finishgame

import (
	"database/sql"
	"fmt"
	math "math/big"
	"processgame/entities"
)

func saveBestResultGame(tx *sql.Tx, g *entities.BinaryOptionGame, bestResultGame *entities.BestResultGame) bool {
	if len(bestResultGame.ListPlayersWin) > 0 {
		query := `UPDATE binary_option_game_bet SET amount_win_dolar = ? WHERE id = ?;`

		var amountWinDolar float64

		for _, v := range bestResultGame.ListPlayersWin {
			x := math.NewFloat(g.GameProfitPercent)
			percent := math.NewFloat(100)

			x.Quo(x, percent)

			betAmountDollar := math.NewFloat(v.BetAmountDolar)

			x.Mul(x, betAmountDollar)

			amountWinDolar, _ = x.Float64()
			fmt.Printf("ID: %v // User: %v %v\n", g.Id, v.IdUsuario, amountWinDolar)

			res, err := tx.Exec(query, amountWinDolar, v.Id)

			if err != nil {
				fmt.Println("SBRG 1: " + err.Error())
				return false
			}

			affcRows, _ := res.RowsAffected()
			if affcRows == 0 {
				fmt.Println("SBRG 2: ")
				return false
			}
		}
	}

	if len(bestResultGame.ListPlayersEqual) > 0 {
		var ids string
		for i, v := range bestResultGame.ListPlayersEqual {
			if i+1 == len(bestResultGame.ListPlayersEqual) {
				ids += fmt.Sprintf("%v", v.Id)
			} else {
				ids += fmt.Sprintf("%v,", v.Id)
			}
		}
		fmt.Printf("IDS: %v \n", ids)
		query := `UPDATE binary_option_game_bet SET refund = 1 WHERE id IN (` + ids + `);`

		res, err := tx.Exec(query)

		if err != nil {
			fmt.Println("SBRG 4: " + err.Error())
			return false
		}

		affcRows, _ := res.RowsAffected()

		if affcRows == 0 {
			fmt.Println("SBRG 4: ")
			return false
		}
	}

	gameTotalBet := math.NewFloat(0)
	TotalWinDolar := math.NewFloat(bestResultGame.TotalWinDolar)
	TotalLoseDolar := math.NewFloat(bestResultGame.TotalLoseDolar)
	TotalEqualDolar := math.NewFloat(bestResultGame.TotalEqualDolar)

	x := math.NewFloat(bestResultGame.TotalLoseDolar)
	gameTotalBet.Add(TotalWinDolar, x)

	y := math.NewFloat(bestResultGame.TotalEqualDolar)
	gameTotalBet.Add(TotalWinDolar, y)

	gameWinAmountPercent := resultToPercent(TotalWinDolar, gameTotalBet, 2)
	gameLoseAmountPercent := resultToPercent(TotalLoseDolar, gameTotalBet, 2)
	gameEqualAmountPercent := resultToPercent(TotalEqualDolar, gameTotalBet, 2)

	// OBS: Necessita estar sincronizado com fechamento de candles !
	// O jogo termina no mesmo valor para todos.
	// TotalDolar & gameTotalBet serão sempre iguais, gerando % de 100 sempre...

	query := `
	UPDATE binary_option_game g
	LEFT JOIN binary_option_game_bet b ON b.id_game = g.id
	SET  g.game_win_amount_dolar = ?
		,g.game_win_amount_percent = ?
		,g.game_lose_amount_dolar = ?
		,g.game_lose_amount_percent = ?
		,g.game_equal_amount_dolar = ?
		,g.game_equal_amount_percent = ?
		,g.game_price_amount_selected_finish = ?
		,g.game_calculated = 1
		,g.total_lose_dolar_trader_bot = ?
		,g.total_win_dolar_trader_bot = ?
	WHERE g.id = ?;
	`

	_, err := tx.Exec(query, bestResultGame.TotalWinDolar, gameWinAmountPercent, bestResultGame.TotalLoseDolar, gameLoseAmountPercent,
		bestResultGame.TotalEqualDolar, gameEqualAmountPercent, bestResultGame.Price, bestResultGame.TotalLoseDolarTraderBot,
		bestResultGame.TotalWinDolarTraderBot, g.Id)

	if err != nil {
		fmt.Println("SBRG 5: " + err.Error())
		return false
	}

	query = `
		UPDATE binary_option_game_bet b
		SET b.status_received_win_payment = 1
		WHERE b.id_game = ? AND b.amount_win_dolar = 0 AND b.refund = 0;`

	_, err = tx.Exec(query, g.Id)

	if err != nil {
		fmt.Println("SBRG 6: " + err.Error())
		return false
	}

	return true
}
