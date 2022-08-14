package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func dbConnect() (db *sql.DB) {
	if db != nil {
		err := db.Ping()

		if err != nil {
			fmt.Println(err)
		}
		return
	}

	db, err := sql.Open("mysql", os.Getenv("db"))

	if err != nil {
		log.Println("mysql conn err: ", err)
	}

	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return
}
