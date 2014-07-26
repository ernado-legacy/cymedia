package query

import (
	"github.com/ernado/cymedia/mediad/models"
	"github.com/ginuerzh/weedo"
	"os"
)

type QueryClient interface {
	Push(fid, requestType string, options models.Options) error
}

type TestClient struct {
	query Query
	weed  *weedo.Client
}

func NewTestClient(query Query) *TestClient {
	c := &TestClient{}
	c.weed = weedo.NewClient("http://localhost:9333")
	c.query = query
	return c
}

func (t *TestClient) Push(fid, requestType string, options models.Options) error {
	req := models.Request{}
	req.Id = fid + requestType
	req.Type = requestType
	req.Options = options
	req.File = fid
	return t.query.Push(req)
}

func (t *TestClient) TestPush(filename, requestType string, options models.Options) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	fid, _, err := t.weed.AssignUpload(filename, "file", f)
	if err != nil {
		return err
	}

	return t.Push(fid, requestType, options)
}
