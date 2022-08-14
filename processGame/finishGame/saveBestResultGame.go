package finishgame

import (
	"database/sql"
	"fmt"
	math "math/big"
	"processgame/entities"
)

func saveBestResultGame(db *sql.DB, g *entities.BinaryOptionGame, bestResultGame *entities.BestResultGame) {
	if len(bestResultGame.ListPlayersWin) > 0 {
		query := `UPDATE binary_option_game_bet SET amount_win_dolar = ? WHERE id IN (?);`
		var ids string
		x := math.NewFloat(g.GameProfitPercent)

		for _, v := range bestResultGame.ListPlayersWin {
			percent := math.NewFloat(100)
			// precision := math.NewFloat(8)
			x.Quo(x, percent)

			betAmountDollar := math.NewFloat(v.BetAmountDolar)
			x.Mul(x, betAmountDollar)
		}

		_, err := db.Exec(query, x, ids)

		if err != nil {
			fmt.Println(err)
		}
	}

	if len(bestResultGame.ListPlayersEqual) > 0 {
		query := `UPDATE binary_option_game_bet SET refund = 1 WHERE id IN (?);`
		var ids string

		_, err := db.Exec(query, ids)

		if err != nil {
			fmt.Println(err)
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

	_, err := db.Exec(query, bestResultGame.TotalWinDolar, gameWinAmountPercent, bestResultGame.TotalLoseDolar, gameLoseAmountPercent,
		bestResultGame.TotalEqualDolar, gameEqualAmountPercent, bestResultGame.Price, bestResultGame.TotalLoseDolarTraderBot,
		bestResultGame.TotalWinDolarTraderBot, g.Id)

	if err != nil {
		fmt.Println(err)
	}

	query = `
	UPDATE binary_option_game_bet b
	SET b.status_received_win_payment = 1
	WHERE b.id_game = ? AND b.amount_win_dolar = 0 AND b.refund = 0;`

	_, err = db.Exec(query, g.Id)

	if err != nil {
		fmt.Println(err)
	}
}
