package query

import (
	"errors"
	"fmt"
	"github.com/ernado/cymedia/mediad/models"
	"github.com/ginuerzh/weedo"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func isFileOK(res models.Responce) error {
	client := weedo.NewClient("localhost:9333")
	url, _, err := client.GetUrl(res.File)
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.ContentLength != 0 {
		return nil
	}
	return errors.New(fmt.Sprintf("bad content length %d", resp.ContentLength))
}

func TestServer(t *testing.T) {
	Convey("Upload", t, func() {
		server, query := NewTestServer()
		client := NewTestClient(query)
		filename := "../conventer/samples/sample.webm"
		o := new(models.VideoOptions)
		o.Video.Format = "libvpx"
		o.Audio.Format = "libvorbis"
		o.Audio.Bitrate = 128 * 1024
		o.Video.Bitrate = 500 * 1024
		err := client.TestPush(filename, "video", o)
		So(err, ShouldBeNil)
		Convey("Convert", func() {
			req, err := query.Pull()
			So(err, ShouldBeNil)
			res, err := server.Process(req)
			So(err, ShouldBeNil)
			So(isFileOK(res), ShouldBeNil)
		})
	})
}
