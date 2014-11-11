package photo

import (
	"os"
	"testing"

	"github.com/ernado/weed"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPhoto(t *testing.T) {
	Convey("Photo upload", t, func() {
		f, err := os.Open("test/image.jpg")
		So(err, ShouldBeNil)
		adapter := weed.NewAdapter("localhost:9333")
		uploader := Uploader{adapter, 1000, 100}
		_, err = uploader.Upload(f)
		So(err, ShouldBeNil)
	})
}
