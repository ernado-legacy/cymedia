package main

import (
	"flag"
	"io"
	"log"

	"github.com/ernado/cymedia/mediad/query"
	"github.com/ernado/selectel/storage"
	"gopkg.in/mgo.v2/bson"
)

var (
	weedUrl           = flag.String("weed", "http://localhost:9333", "Weed master url")
	redisHost         = flag.String("redis.addr", ":6379", "Redis server address")
	redisKey          = flag.String("redis.key", "cymedia:query", "Redis query key")
	selectel          bool
	selectelContainer = flag.String("selectel.container", "", "Selectel container")
)

func init() {
	flag.BoolVar(&selectel, "selectel", false, "Use selectel storage")
}

type selectelAdapter struct {
	api storage.ContainerAPI
}

func (s selectelAdapter) GetUrl(name string) (string, error) {
	return s.api.URL(name), nil
}

func (s selectelAdapter) URL(name string) (string, error) {
	return s.GetUrl(name)
}

func (s selectelAdapter) Upload(reader io.Reader, t, format string) (fid string, purl string, size int64, err error) {
	fid = bson.NewObjectId().Hex()
	err = s.api.Upload(reader, fid, t+"/"+format)
	purl = s.api.URL(fid)
	return
}

func main() {
	flag.Parse()
	log.Println("connecting")
	var (
		server     query.QueryServer
		err        error
		selStorage storage.API
	)
	if selectel {
		selStorage, err = storage.NewEnv()
		if err != nil {
			log.Fatal(err)
		}
		container := selStorage.Container(*selectelContainer)
		server, err = query.NewRedisSelectelServer(selectelAdapter{container}, *redisHost, *redisKey)
	} else {
		server, err = query.NewRedisServer(*weedUrl, *redisHost, *redisKey)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("starting")
	server.Main()
}
