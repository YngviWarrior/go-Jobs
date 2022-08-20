package main

import (
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	for range time.Tick(time.Second * 3) {
		memcacheConnect()
		conn := dbConnect()

		// createBet(conn)
		playGames(conn)
		// time.Sleep(time.Second * 2)

		calcGames(conn)
		// time.Sleep(time.Second * 2)

		finishGames(conn)
		// conn.Close()
	}
}
