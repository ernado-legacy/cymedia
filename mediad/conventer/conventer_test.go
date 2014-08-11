package conventer

import (
	"github.com/ernado/cymedia/mediad/models"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestVideoConvertation(t *testing.T) {
	Convey("Error handling", t, func() {
		Convey("Bad format", func() {
			filename := "samples/sample.webm"
			f, err := os.Open(filename)
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			defer f.Close()

			c := VideoConventer{}
			o := new(models.VideoOptions)
			o.Video.Format = "badformat"
			o.Audio.Format = "aac"
			o.Audio.Bitrate = 128 * 1024
			o.Video.Bitrate = 500 * 1024
			_, err = c.Convert(f, o)
			So(err, ShouldNotBeNil)
		})
		Convey("ffmpeg error", func() {
			filename := "samples/sample.webm"
			f, err := os.Open(filename)
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			defer f.Close()

			c := VideoConventer{}
			o := new(models.VideoOptions)
			o.Video.Format = "libvpx"
			o.Audio.Format = "libvo"
			o.Audio.Bitrate = 128 * 1024
			o.Video.Bitrate = 500 * 1024
			_, err = c.Convert(f, o)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Video", t, func() {
		filename := "samples/sample.webm"
		Convey("mp4", func() {
			f, err := os.Open(filename)
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			defer f.Close()

			c := VideoConventer{}
			o := new(models.VideoOptions)
			o.Video.Format = "h264"
			o.Audio.Format = "aac"
			o.Audio.Bitrate = 128 * 1024
			o.Video.Bitrate = 500 * 1024
			_, err = c.Convert(f, o)
			So(err, ShouldBeNil)
		})
		Convey("webm", func() {
			f, err := os.Open(filename)
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			defer f.Close()

			c := VideoConventer{}
			o := new(models.VideoOptions)
			o.Video.Format = "libvpx"
			o.Audio.Format = "libvorbis"
			o.Audio.Bitrate = 128 * 1024
			o.Video.Bitrate = 500 * 1024
			_, err = c.Convert(f, o)
			So(err, ShouldBeNil)
		})
	})

	Convey("Thumbnail", t, func() {
		filename := "samples/sample.webm"
		Convey("start", func() {
			f, err := os.Open(filename)
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			defer f.Close()

			c := VideoConventer{}
			o := new(models.ThumbnailOptions)
			o.Format = "png"
			_, err = c.Convert(f, o)
			So(err, ShouldBeNil)
		})
		Convey("5s", func() {
			f, err := os.Open(filename)
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			defer f.Close()

			c := VideoConventer{}
			o := new(models.ThumbnailOptions)
			o.Format = "png"
			o.Time = 5
			_, err = c.Convert(f, o)
			So(err, ShouldBeNil)
		})
	})
	Convey("Audio", t, func() {
		filename := "samples/sample.ogg"
		Convey("aac", func() {
			f, err := os.Open(filename)
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			defer f.Close()

			c := VideoConventer{}
			o := new(models.AudioOptions)
			o.Format = "aac"
			o.Bitrate = 128 * 1024
			_, err = c.Convert(f, o)
			So(err, ShouldBeNil)
		})
		Convey("ogg", func() {
			f, err := os.Open(filename)
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			defer f.Close()

			c := VideoConventer{}
			o := new(models.AudioOptions)
			o.Format = "libvorbis"
			o.Bitrate = 128 * 1024
			_, err = c.Convert(f, o)
			So(err, ShouldBeNil)
		})
	})
}
