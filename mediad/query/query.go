package query

import (
	"github.com/ernado/cymedia/mediad/models"
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

func (m *MemoryQuery) Pull() (models.Request, error) {
	return <-m.requests, nil
}

func (m *MemoryQuery) Push(request models.Request) error {
	m.requests <- request
	return nil
}

func NewMemoryQuery() Query {
	q := new(MemoryQuery)
	q.requests = make(chan models.Request, 10)
	return q
}
