package main

import (
	"fmt"
	"os"

	"processgame/entities"

	"github.com/bradfitz/gomemcache/memcache"
)

var cache *memcache.Client

func memcacheConnect() {
	var memcachedAddr = os.Getenv("memcached")
	cache = memcache.New(memcachedAddr)

	err := cache.Ping()

	if err != nil {
		fmt.Println("MC 1: " + err.Error())
		os.Exit(2)
	}

	i := memcache.Item{}
	i.Key = "test"
	i.Value = []byte{}
	err = cache.Set(&i)

	if err != nil {
		fmt.Println("MC 2: " + err.Error())
		os.Exit(2)
	}
}

func push(toCache []entities.GamesUpdate) {
	gameList, err := cache.Get("GamesUpdate")
	i := memcache.Item{}
	i.Key = "GamesUpdate"

	if err != nil {
		i.Value = []byte(fmt.Sprintf("%v", toCache))
		cache.Set(&i)
	} else if gameList != nil {
		newList := append(gameList.Value, []byte(fmt.Sprintf("%v", toCache))...)

		i.Value = []byte(fmt.Sprintf("%v", newList))
		cache.Set(&i)
	}
}
