package finishgame

import (
	"database/sql"
	"fmt"
	math "math/big"
	"processgame/entities"
)

func saveBestResultGame(tx *sql.Tx, g *entities.BinaryOptionGame, bestResultGame *entities.BestResultGame) bool {
	if len(bestResultGame.ListPlayersWin) > 0 {
		query := `UPDATE binary_option_game_bet SET amount_win_dolar = ? WHERE id IN (?);`

		x := math.NewFloat(g.GameProfitPercent)

		for _, v := range bestResultGame.ListPlayersWin {
			percent := math.NewFloat(100)

			x.Quo(x, percent)

			betAmountDollar := math.NewFloat(v.BetAmountDolar)
			x.Mul(x, betAmountDollar)
		}

		_, err := tx.Exec(query, x, g.Id)

		if err != nil {
			fmt.Println("SBRG 1: " + err.Error())
			return false
		}
	}

	if len(bestResultGame.ListPlayersEqual) > 0 {
		query := `UPDATE binary_option_game_bet SET refund = 1 WHERE id IN (?);`

		_, err := tx.Exec(query, g.Id)

		if err != nil {
			fmt.Println("SBRG 2: " + err.Error())
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
		fmt.Println(err)
		return false
	}

	query = `
	UPDATE binary_option_game_bet b
	SET b.status_received_win_payment = 1
	WHERE b.id_game = ? AND b.amount_win_dolar = 0 AND b.refund = 0;`

	_, err = tx.Exec(query, g.Id)

	if err != nil {
		fmt.Println("SBRG 3: " + err.Error())
		return false
	}

	return true
}
