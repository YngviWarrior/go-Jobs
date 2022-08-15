package entities

import "database/sql"

type BonusTrader struct {
	Id                    uint64  `json:"id"`
	IdUser                uint64  `json:"id_user"`
	IdUserOrigin          uint64  `json:"id_user_origin"`
	IdGame                uint64  `json:"id_game"`
	IdGameBet             uint64  `json:"id_game_bet"`
	IdBalance             uint64  `json:"id_balance"`
	Valor                 float64 `json:"valor"`
	DateRegister          string  `json:"date_register"`
	StatusReceivedPayment uint64  `json:"status_received_payment"`
}

type Saldos struct {
	Id                       uint64        `json:"id"`
	IdUsuario                uint64        `json:"id_usuario"`
	IdTipo                   uint64        `json:"id_tipo"`
	IdOrigem                 uint64        `json:"id_origem"`
	Valor                    float64       `json:"valor"`
	TotalAntes               float64       `json:"total_antes"`
	TotalDepois              float64       `json:"total_depois"`
	IdPedido                 sql.NullInt64 `json:"id_pedido"`
	IdBinaryOptionGameBet    sql.NullInt64 `json:"id_binary_option_game_bet"`
	IdBinaryOptionGameBetWin sql.NullInt64 `json:"id_binary_option_game_bet_win"`
	IdBonusIndicacao         sql.NullInt64 `json:"id_bonus_indicacao"`
	IdBonusTrader            sql.NullInt64 `json:"id_bonus_trader"`
	IdBtcPayment             sql.NullInt64 `json:"id_btc_payment"`
	IdUserBalanceConvert     sql.NullInt64 `json:"id_user_balance_convert"`
	IdDepositAdmin           sql.NullInt64 `json:"id_deposit_admin"`
	DataRegistro             string        `json:"data_registro"`
}

type BonusIndicacao struct {
	Id                     uint64  `json:"id"`
	IdUser                 uint64  `json:"id_user"`
	IdUserOrigin           uint64  `json:"id_user_origin"`
	IdGame                 uint64  `json:"id_game"`
	IdGameBet              uint64  `json:"id_game_bet"`
	IdBalance              uint64  `json:"id_balance"`
	IdType                 uint64  `json:"id_type"`
	Valor                  float64 `json:"valor"`
	DateRegister           string  `json:"date_register"`
	StatusReceivedOPayment uint64  `json:"status_received_payment"`
}

type Candle struct {
	IdMoedasPares int64
	Mts           int64
	Open          float64
	Close         float64
	High          float64
	Low           float64
	Volume        float64
	IdGameCustom  sql.NullInt64
	OpenCustom    sql.NullFloat64
	CloseCustom   sql.NullFloat64
	HighCustom    sql.NullFloat64
	LowCustom     sql.NullFloat64
	VolumeCustom  sql.NullFloat64
	MakeHighLow   sql.NullBool
}

type BestResultGame struct {
	Id                      uint64                 `json:"id"`
	IdGame                  uint64                 `json:"id_game"`
	Price                   float64                `json:"price"`
	TotalWinDolar           float64                `json:"total_win_dolar"`
	TotalLoseDolar          float64                `json:"total_lose_dolar"`
	TotalDifWinLose         float64                `json:"total_dif_win_lose"`
	TotalEqualDolar         float64                `json:"total_equal_dolar"`
	ListPlayersWin          []*BinaryOptionGameBet `json:"list_players_win"`
	ListPlayersLose         []*BinaryOptionGameBet `json:"list_players_lose"`
	ListPlayersEqual        []*BinaryOptionGameBet `json:"list_players_equal"`
	TotalLoseDolarTraderBot int64                  `json:"total_lose_dolar_trader_bot"`
	TotalWinDolarTraderBot  int64                  `json:"total_win_dolar_trader_bot"`
}

type MoedasPares struct {
	Id                   uint64         `json:"id"`
	Nome                 sql.NullString `json:"nome"`
	Symbol               sql.NullString `json:"symbol"`
	SymbolSync           sql.NullString `json:"symbol_sync"`
	IdCoin1              sql.NullInt64  `json:"id_coin1"`
	IdCoin2              sql.NullInt64  `json:"id_coin2"`
	UtilizarGiro         int64          `json:"utilizar_giro"`
	SyncOnlyExchangeRate int64          `json:"sync_only_exchange_rate"`
	NotList              uint64         `json:"not_list"`
}

type BinaryOptionGameBet struct {
	Id                                  uint64  `json:"id"`
	HashId                              string  `json:"hash_id"`
	IdGame                              uint64  `json:"id_game"`
	IdUsuario                           uint64  `json:"id_usuario"`
	IdChoice                            uint64  `json:"id_choice"`
	IdBalance                           uint64  `json:"id_balance"`
	BetAmountDolar                      float64 `json:"bet_amount_dolar"`
	AmountWinDolar                      float64 `json:"amount_win_dolar"`
	PriceAmountSelected                 float64 `json:"price_amount_selected"`
	StatusReceivedWinPayment            uint64  `json:"status_received_win_payment"`
	IdTraderFollower                    uint64  `json:"id_trader_follower"`
	BotUseStatus                        uint64  `json:"bot_use_status"`
	DateRegister                        string  `json:"date_register"`
	BonusTraderPercentFromTaxBetWin     float64 `json:"bonus_trader_percent_from_tax_bet_win"`
	BonusIndicationPercentFromTaxBetWin float64 `json:"bonus_indication_percent_from_tax_bet_win"`
	StatusReceivedRefundPayment         uint64  `json:"status_received_refund_payment"`
	Refund                              uint64  `json:"refund"`
	Deleted                             uint64  `json:"deleted"`
}

type GamesUpdate struct {
	Id             uint64 `json:"id"`
	IdMoedasPares  uint64 `json:"id_moedas_pares"`
	GameIdTypeTime int64  `json:"game_id_type_time"`
}

type BinaryOptionGame struct {
	Id                               uint64          `json:"id"`
	IdMoedasPares                    uint64          `json:"id_moedas_pares"`
	GameIdStatus                     uint64          `json:"game_id_status"`
	GameIdTypeTime                   int64           `json:"game_id_type_time"`
	GameDateStart                    string          `json:"game_date_start"`
	GameDateProcess                  string          `json:"game_date_process"`
	GameDateFinish                   string          `json:"game_date_finish"`
	GameProfitPercent                float64         `json:"game_profit_percent"`
	GameBetAmountDolarTotal          float64         `json:"game_bet_amount_dolar_total"`
	BonusTraderTotalAmountDolar      float64         `json:"bonus_trader_total_amount_dolar"`
	BonusIndicationTotalAmountDolar  float64         `json:"bonus_indication_total_amount_dolar"`
	GameWinAmountPercent             float64         `json:"game_win_amount_percent"`
	GameWinAmountDolar               float64         `json:"game_win_amount_dolar"`
	GameLoseAmountPercent            float64         `json:"game_lose_amount_percent"`
	GameLoseAmountDolar              float64         `json:"game_lose_amount_dolar"`
	GameEqualAmountPercent           float64         `json:"game_equal_amount_percent"`
	GameEqualAmountDolar             float64         `json:"game_equal_amount_dolar"`
	GamePriceAmountSelectedFinish    float64         `json:"game_price_amount_selected_finish"`
	CompanyTaxAmountDolarFromGameTax float64         `json:"company_tax_amount_dolar_from_game_tax"`
	GameCalculated                   int64           `json:"game_calculated"`
	PartitionFilter                  sql.NullFloat64 `json:"partition_filter"`
	TotalWinDolarTraderBot           sql.NullFloat64 `json:"total_win_dolar_trader_bot"`
	IdTraderBotAction                sql.NullInt64   `json:"id_trader_bot_action"`
}
