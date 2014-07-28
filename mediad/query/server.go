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
	"time"
)

var (
	ErrorBadType = errors.New("bad type")
)

type Server struct {
	conn  redis.Conn
	weed  *weedo.Client
	video conventer.Conventer
	query Query
}

func (s *Server) MakeResponce(request models.Request) error {
	responce, err := s.Process(request)
	if err != nil {
		responce = models.Responce{Id: request.Id, Success: false, Error: err.Error()}
	}

	data, err := json.Marshal(responce)
	if err != nil {
		return err
	}
	_, err = s.conn.Do("LPUSH", request.ResultKey, data)
	return err
}

func (s *Server) Iteration() error {
	request, err := s.query.Pull()
	log.Println("got request", request.Id)
	if err != nil {
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
	if req.Type == "video" {
		return s.video.Convert(resp.Body, req.Options)
	}
	return output, ErrorBadType
}

func (s *Server) Process(request models.Request) (response models.Responce, err error) {
	output, err := s.Convert(request)
	if err != nil {
		return
	}
	options := request.Options
	fid, _, err := s.weed.AssignUpload(fmt.Sprintf("file.%s", options.Extension()), options.Mime(), output)
	if err != nil {
		return
	}
	response.Id = request.Id
	response.File = fid
	response.Format = options.Extension()
	return
}
