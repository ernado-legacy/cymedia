package query

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ernado/cymedia/mediad/models"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func randStr(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func TestQuery(t *testing.T) {
	Convey("Memory query", t, func() {

		q := NewMemoryQuery()
		request := models.Request{}
		request.Id = "request"

		So(q.Push(request), ShouldBeNil)
		r, err := q.Pull()
		So(err, ShouldBeNil)
		So(r.Id, ShouldEqual, request.Id)
	})

	Convey("Redis query", t, func() {
		q, err := NewRedisQuery(":6379", fmt.Sprintf("test:cymedia:%s", randStr(20)))
		So(err, ShouldBeNil)
		request := models.Request{}
		request.Id = "request"
		So(q.Push(request), ShouldBeNil)
		r, err := q.Pull()
		So(err, ShouldBeNil)
		So(r.Id, ShouldEqual, request.Id)
	})
}
