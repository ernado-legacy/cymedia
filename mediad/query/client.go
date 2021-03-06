package query

import (
	"github.com/ernado/cymedia/mediad/models"
	"github.com/ginuerzh/weedo"
	"os"
)

type QueryClient interface {
	Push(id, fid, requestType string, options models.Options) error
}

type Client struct {
	resultKey string
	query     Query
	weed      *weedo.Client
}

func NewTestClient(query Query) *Client {
	c := new(Client)
	c.weed = weedo.NewClient("http://localhost:9333")
	c.query = query
	return c
}

func NewRedisClient(weedUrl, redisHost, redisKey, resultKey string) (client QueryClient, err error) {
	c := new(Client)
	c.weed = weedo.NewClient(weedUrl)
	c.query, err = NewRedisQuery(redisHost, redisKey)
	c.resultKey = resultKey
	return c, err
}

func (t *Client) Push(id, fid, requestType string, options models.Options) error {
	req := models.Request{}
	req.Id = id
	req.Type = requestType
	req.Options = options
	req.File = fid
	req.ResultKey = t.resultKey
	return t.query.Push(req)
}

func (t *Client) FilePush(filename, requestType string, options models.Options) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	fid, _, err := t.weed.AssignUpload(filename, "file", f)
	if err != nil {
		return err
	}
	return t.Push(fid, fid, requestType, options)
}
