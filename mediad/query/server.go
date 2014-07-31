package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ernado/cymedia/mediad/conventer"
	"github.com/ernado/cymedia/mediad/models"
	"github.com/garyburd/redigo/redis"
	"github.com/ginuerzh/weedo"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

var (
	ErrorBadType  = errors.New("bad type")
	ErrorBadKey   = errors.New("bad key")
	ErrorCritical = errors.New("critical error")
)

type Server struct {
	conn  redis.Conn
	weed  *weedo.Client
	video conventer.Conventer
	query Query
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
		responce = models.Responce{Id: request.Id, Success: false, Error: err.Error()}
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
	var err error
	var sleep time.Duration
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
	s := &Server{}
	s.weed = weedo.NewClient("http://localhost:9333")
	s.query = NewMemoryQuery()
	s.video = &conventer.VideoConventer{}
	return s, s.query
}

func NewRedisServer(weedUrl, redisHost, redisKey string) (server QueryServer, err error) {
	s := new(Server)
	s.weed = weedo.NewClient(weedUrl)
	s.video = new(conventer.VideoConventer)
	s.query, err = NewRedisQuery(redisHost, redisKey)
	s.conn, err = redis.Dial("tcp", redisHost)
	return s, err
}

func (s *Server) Convert(req models.Request) (output io.ReadCloser, err error) {
	url, _, err := s.weed.GetUrl(req.File)
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
	if req.Type == "video" {
		return s.video.Convert(resp.Body, options)
	}
	if req.Type == "audio" {
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
	fid, _, err := s.weed.AssignUpload(fmt.Sprintf("file.%s", options.Extension()), options.Mime(), output)
	if err != nil {
		return
	}
	response.File = fid
	response.Success = true
	return
}
