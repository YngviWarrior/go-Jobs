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

		playGames(conn)
		// time.Sleep(time.Second * 1)

		calcGames(conn)
		// time.Sleep(time.Second * 1)

		finishGames(conn)
	}
}
