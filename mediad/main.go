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
	selectelKey       = flag.String("selectel.key", "", "Selectel key")
	selectelUser      = flag.String("selectel.user", "", "Selectel user")
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
		if len(*selectelKey) != 0 && len(*selectelUser) != 0 {
			log.Println(*selectelUser, *selectelKey)
			selStorage = storage.NewAsync(*selectelUser, *selectelKey)
			selStorage.Debug(true)
			err = selStorage.Auth(*selectelUser, *selectelKey)
		} else {
			selStorage, err = storage.NewEnv()
		}
		if err != nil {
			log.Fatal(err)
		}
		selStorage.Debug(true)
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
