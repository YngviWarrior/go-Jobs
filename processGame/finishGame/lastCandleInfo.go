package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
)

func lastCandleInfo(db *sql.DB, game *entities.BinaryOptionGame) (b entities.BestResultGame) {
	// query = `SELECT id, nome, symbol, symbol_sync, id_coin1, id_coin2, utilizar_giro, sync_only_exchange_rate, not_list
	// FROM moedas_pares m
	// WHERE m.id = ?`

	// var m MoedasPares
	// err = db.QueryRow(query, game.IdMoedasPares).Scan(&m.Id, &m.Nome, &m.Symbol, &m.SymbolSync, &m.IdCoin1, &m.IdCoin2,
	// 	&m.UtilizarGiro, &m.SyncOnlyExchangeRate, &m.NotList)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	query :=
		`SELECT mts,
			COALESCE(open_custom,open) as open,
			COALESCE(close_custom,close) as close,
			COALESCE(high_custom,high) as high,
			COALESCE(low_custom,low) as low,
			COALESCE(volume_custom,volume) as volume
	
		FROM candle_1m
		WHERE id_moedas_pares = ?
				
		ORDER BY mts DESC
		LIMIT 0,1`

	var c entities.Candle

	err := db.QueryRow(query, game.IdMoedasPares).Scan(&c.Mts, &c.Open, &c.Close, &c.High, &c.Low, &c.Volume)

	if err != nil {
		fmt.Println(err)
		return
	}

	b.Price = c.Close

	return
}
