package photo

import (
	"github.com/ernado/weed"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestPhoto(t *testing.T) {
	Convey("Photo upload", t, func() {
		f, err := os.Open("test/image.jpg")
		So(err, ShouldBeNil)
		adapter := weed.NewAdapter("localhost:9333")
		uploader := Uploader{adapter, 1000, 100}
		progress := make(chan float32)
		stat, err := f.Stat()
		So(err, ShouldBeNil)
		var p float32
		go func() {
			for p = range progress {
				continue
			}
		}()
		_, err = uploader.Upload(stat.Size(), f, progress)
		So(err, ShouldBeNil)
		So(p, ShouldAlmostEqual, 100.0)
	})
}
