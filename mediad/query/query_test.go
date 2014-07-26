package query

import (
	"github.com/ernado/cymedia/mediad/models"
	. "github.com/smartystreets/goconvey/convey"
	// "os"
	"testing"
)

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
}
