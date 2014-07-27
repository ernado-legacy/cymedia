package query

import (
	"encoding/json"
	"github.com/ernado/cymedia/mediad/models"
	"github.com/garyburd/redigo/redis"
)

type Query interface {
	Pull() (models.Request, error)
	Push(request models.Request) error
}

type QueryServer interface {
	Process(request models.Request) (models.Responce, error)
}

type MemoryQuery struct {
	requests chan models.Request
}

type RedisQuery struct {
	conn redis.Conn
	key  string
}

func (m *MemoryQuery) Pull() (models.Request, error) {
	return <-m.requests, nil
}

func (m *MemoryQuery) Push(request models.Request) error {
	m.requests <- request
	return nil
}

func (r *RedisQuery) Pull() (req models.Request, err error) {
	reply, err := redis.Strings(r.conn.Do("BLPOP", r.key, 0))
	if err != nil {
		return
	}
	return req, json.Unmarshal([]byte(reply[1]), &req)
}

func (r *RedisQuery) Push(request models.Request) error {
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	_, err = r.conn.Do("RPUSH", r.key, string(data))
	return err
}

func NewMemoryQuery() Query {
	q := new(MemoryQuery)
	q.requests = make(chan models.Request, 10)
	return q
}

func NewRedisQuery(host, key string) (Query, error) {
	conn, err := redis.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	q := new(RedisQuery)
	q.key = key
	q.conn = conn
	return q, nil
}
