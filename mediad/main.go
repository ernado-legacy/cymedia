package main

import (
	"flag"
	"github.com/ernado/cymedia/mediad/query"
	"log"
)

var (
	weedUrl   = flag.String("weed", "http://localhost:9333", "Weed master url")
	redisHost = flag.String("redis.addr", ":6379", "Redis server address")
	redisKey  = flag.String("redis.key", "cymedia:query", "Redis query key")
)

func main() {
	flag.Parse()
	log.Println("connecting")
	server, err := query.NewRedisServer(*weedUrl, *redisHost, *redisKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("starting")
	server.Main()
}
