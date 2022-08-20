package main

import (
	"encoding/json"
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

func push(toCache []*entities.GamesUpdate) {
	_, err := cache.Get("GamesUpdate")
	i := memcache.Item{}
	i.Key = "GamesUpdate"

	jsonToCache, err2 := json.Marshal(toCache)

	if err2 != nil {
		fmt.Println("Marshal Json: " + err.Error())
	}

	i.Value = jsonToCache
	cache.Set(&i)
}
