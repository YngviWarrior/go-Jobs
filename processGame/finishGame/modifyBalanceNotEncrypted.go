package finishgame

import (
	"database/sql"
	"fmt"
	math "math/big"
	"processgame/entities"
)

func modifyBalanceNotEncrypted(db *sql.DB, idUser uint64, idBalance uint64, value float64, acceptNegative bool) (beforeValue float64, afterValue float64) {
	query := `
		SELECT valor
		FROM saldo_valor
		WHERE id_usuario = ?
		AND id = ?
		FOR UPDATE;
	`

	var s entities.Saldos
	err := db.QueryRow(query, idUser, idBalance).Scan(&s.Valor)

	if err != nil {
		fmt.Println(err)
	}

	beforeValue = s.Valor

	val := math.NewFloat(value)
	bVal := math.NewFloat(beforeValue)
	bVal.Add(bVal, val)
	v, _ := bVal.Float64()

	beforeValue = v

	if !acceptNegative && bcSimplesComp(afterValue, "<", 0, 8) {
		fmt.Println("After Value Neg Value")
		return
	}

	query = `
		UPDATE saldo_valor
		SET valor = valor + :valor
		WHERE id_usuario = ?
		AND id = ?
	`

	_, err = db.Exec(query, idUser, idBalance)

	if err != nil {
		fmt.Println(err)
	}

	return
}
