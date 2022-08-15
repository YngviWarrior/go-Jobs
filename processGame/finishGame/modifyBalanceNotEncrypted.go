package finishgame

import (
	"database/sql"
	"fmt"
	math "math/big"
	"processgame/entities"
)

func modifyBalanceNotEncrypted(tx *sql.Tx, idUser uint64, idBalance uint64, value float64, acceptNegative bool) (beforeValue float64, afterValue float64) {
	query := `
		SELECT valor
		FROM saldo_valor
		WHERE id_usuario = ?
		AND id = ?
		FOR UPDATE;
	`

	var s entities.Saldos
	err := tx.QueryRow(query, idUser, idBalance).Scan(&s.Valor)

	if err != nil {
		fmt.Println("MCNE 1: " + err.Error())
	}

	beforeValue = s.Valor

	val := math.NewFloat(value)
	bVal := math.NewFloat(beforeValue)
	bVal.Add(bVal, val)
	v, _ := bVal.Float64()

	afterValue = v

	if !acceptNegative && bcSimplesComp(afterValue, "<", 0, 8) {
		fmt.Println("After Value Neg Value")
		return
	}

	query = `
		UPDATE saldo_valor
		SET valor = ?
		WHERE id_usuario = ?
		AND id = ?
	`

	res, err := tx.Exec(query, afterValue, idUser, idBalance)

	if err != nil {
		fmt.Println("MCNE 2: " + err.Error())
	}

	affcRows, _ := res.RowsAffected()
	if affcRows == 0 {
		fmt.Println("MCNE 3: ")
		return
	}

	return
}
