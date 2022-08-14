package main

import (
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	for range time.Tick(time.Second * 5) {
		memcacheConnect()
		conn := dbConnect()

		go playGames(conn)

		go calcGames(conn)

		go finishGames(conn)
	}
}
