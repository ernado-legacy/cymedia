package query

import (
	"errors"
	"fmt"
	"github.com/ernado/cymedia/mediad/conventer"
	"github.com/ernado/cymedia/mediad/models"
	"github.com/ginuerzh/weedo"
	"io"
	"net/http"
)

var (
	ErrorBadType = errors.New("bad type")
)

type Server struct {
	weed  *weedo.Client
	video conventer.Conventer
	query Query
}

func NewTestServer() (QueryServer, Query) {
	s := &Server{}
	s.weed = weedo.NewClient("http://localhost:9333")
	s.query = NewMemoryQuery()
	s.video = &conventer.VideoConventer{}
	return s, s.query
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
