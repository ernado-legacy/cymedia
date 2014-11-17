package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/ernado/cymedia/mediad/conventer"
	"github.com/ernado/cymedia/mediad/models"
	"github.com/ernado/weed"
	"github.com/garyburd/redigo/redis"
)

var (
	ErrorBadType  = errors.New("bad type")
	ErrorBadKey   = errors.New("bad key")
	ErrorCritical = errors.New("critical error")
)

type StorageAdapter interface {
	GetUrl(fid string) (url string, err error)
	Upload(reader io.Reader, t, format string) (fid string, purl string, size int64, err error)
}

type Server struct {
	conn    redis.Conn
	storage StorageAdapter
	video   conventer.Conventer
	query   Query
}

func (s *Server) MakeResponce(request models.Request) (err error) {
	defer func() {
		rec := recover()
		if rec != nil {
			log.Println(rec)
			err = ErrorCritical
			debug.PrintStack()
		}
	}()
	log.Println("making responce")
	responce, err := s.Process(request)
	log.Println("responce generated", responce, err)
	if err != nil {
		responce = models.Responce{Id: request.Id, Success: false, Error: err.Error(), Type: request.Type}
	}

	log.Println("marshaling")
	data, err := json.Marshal(responce)
	if err != nil {
		return err
	}
	log.Println("pushing", string(data))
	if request.ResultKey == "" {
		return ErrorBadKey
	}
	_, err = s.conn.Do("LPUSH", request.ResultKey, data)
	return err
}

func (s *Server) Iteration() error {
	request, err := s.query.Pull()
	log.Println("got request", request.Id)
	log.Printf("%+v\n", request)
	if err != nil {
		log.Println("error processing request", err)
		return err
	}
	return s.MakeResponce(request)
}

func (s *Server) Main() {
	var (
		sleep time.Duration
		err   error
	)
	log.Println("started")
	for {
		if err = s.Iteration(); err != nil {
			if sleep == time.Second*0 {
				sleep = time.Second * 1
			}
			sleep = sleep * 2
			log.Println(err, "sleeping for", sleep)
		} else {
			sleep = time.Second * 0
		}
		time.Sleep(sleep)
	}
}

func NewTestServer() (QueryServer, Query) {
	s := new(Server)
	s.storage = weed.NewAdapter("http://localhost:9333")
	s.query = NewMemoryQuery()
	s.video = new(conventer.VideoConventer)
	return s, s.query
}

func NewRedisServer(weedUrl, redisHost, redisKey string) (server QueryServer, err error) {
	s := new(Server)
	s.storage = weed.NewAdapter(weedUrl)
	s.video = new(conventer.VideoConventer)
	s.query, err = NewRedisQuery(redisHost, redisKey)
	s.conn, err = redis.Dial("tcp", redisHost)
	return s, err
}

func NewRedisSelectelServer(storage StorageAdapter, redisHost, redisKey string) (server QueryServer, err error) {
	s := new(Server)
	s.storage = storage
	s.video = new(conventer.VideoConventer)
	s.query, err = NewRedisQuery(redisHost, redisKey)
	s.conn, err = redis.Dial("tcp", redisHost)
	return s, err
}

func (s *Server) Convert(req models.Request) (output io.ReadCloser, err error) {
	url, err := s.storage.GetUrl(req.File)
	if err != nil {
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	log.Println("getting options")
	options := req.GetOptions()
	if options == nil {
		err = ErrorBadType
		return
	}
	log.Println("converting")
	if req.Type == models.VideoType {
		return s.video.Convert(resp.Body, options)
	}
	if req.Type == models.AudioType {
		return s.video.Convert(resp.Body, options)
	}
	if req.Type == models.ThumbnailType {
		return s.video.Convert(resp.Body, options)
	}
	return output, ErrorBadType
}

func (s *Server) Process(request models.Request) (response models.Responce, err error) {
	log.Println("processing")
	options := request.GetOptions()
	response.Id = request.Id
	response.Format = options.Extension()
	response.Type = request.Type
	output, err := s.Convert(request)
	if err != nil {
		return
	}
	fid, _, _, err := s.storage.Upload(output, fmt.Sprintf("file.%s", options.Extension()), options.Mime())
	if err != nil {
		return
	}
	response.File = fid
	response.Success = true
	return
}
