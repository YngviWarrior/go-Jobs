package main

import (
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	for range time.Tick(time.Second * 1) {
		memcacheConnect()
		conn := dbConnect()

		// go playGames(conn)

		// go calcGames(conn)

		go finishGames(conn)
	}
}
