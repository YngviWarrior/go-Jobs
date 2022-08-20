package finishgame

import (
	"database/sql"
	"fmt"
	"processgame/entities"
	"time"
)

func modifyBalance(tx *sql.Tx, idUser uint64, idBalance uint64, idOrigin uint64, value float64, idRef uint64, acceptNegativeBalance bool) (afterValue float64) {
	var balance float64

	err := tx.QueryRow(`
		SELECT valor
		FROM saldo_valor
		WHERE id_usuario = ?
		AND id = ?
		FOR UPDATE;
	`, idUser, idBalance).Scan(&balance)

	if err != nil {
		fmt.Println("MB 1: no balance.")
	}

	var s entities.Saldos

	switch idBalance {
	case 3, 20:
		s.IdBonusIndicacao = sql.NullInt64{Int64: int64(idRef), Valid: true}
	case 7, 10:
		s.IdBinaryOptionGameBet = sql.NullInt64{Int64: int64(idRef), Valid: true}
	default:
		return
	}

	beforeValue, afterValue := modifyBalanceNotEncrypted(tx, idUser, idBalance, value, acceptNegativeBalance)

	s.IdUsuario = idUser
	s.IdTipo = idBalance
	s.IdOrigem = idOrigin
	s.Valor = value
	s.TotalAntes = beforeValue
	s.TotalDepois = afterValue
	s.DataRegistro = time.Now().Format("2006-01-02 15:04:05")

	query := `
		UPDATE saldos
		SET valor = ?, total_antes = ?, total_depois = ?, data_registro = ?, id_origem = ?
		WHERE id_usuario = ? AND id_tipo = ?
	`

	res, err := tx.Exec(query, s.Valor, s.TotalAntes, s.TotalDepois, s.DataRegistro, s.IdOrigem, s.IdUsuario, s.IdTipo)

	if err != nil {
		fmt.Println("MB 2: " + err.Error())
	}

	affcRows, _ := res.RowsAffected()
	if affcRows == 0 {
		fmt.Println("MB 3: ")
		return
	}

	return
}
