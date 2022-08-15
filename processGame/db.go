package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func dbConnect() (db *sql.DB) {
	if db != nil {
		err := db.Ping()

		if err != nil {
			fmt.Println("DC 1: " + err.Error())
		}

		return
	}

	db, err := sql.Open("mysql", os.Getenv("db"))

	if err != nil {
		fmt.Println("DC 2: ", err.Error())
	}

	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return
}
