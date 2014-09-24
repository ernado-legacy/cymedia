package models

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProbe(t *testing.T) {
	Convey("Probe", t, func() {
		data := `
{
    "streams": [
        {
            "index": 0,
            "codec_name": "aac",
            "codec_long_name": "AAC (Advanced Audio Coding)",
            "codec_type": "audio",
            "codec_time_base": "1/44100",
            "codec_tag_string": "mp4a",
            "duration_ts": 311296,
            "duration": "7.058866",
            "bit_rate": "63456",
            "tags": {
                "creation_time": "2014-08-26 08:54:56",
                "language": "und",
                "handler_name": "Core Media Data Handler"
            }
        },
        {
            "index": 1,
            "codec_name": "h264",
            "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
            "profile": "Main",
            "codec_type": "video",
            "codec_time_base": "1/1200",
            "codec_tag_string": "avc1",
            "codec_tag": "0x31637661",
            "width": 1280,
            "height": 720,
            "duration_ts": 4224,
            "duration": "7.040000",
            "tags": {
                "rotate": "90",
                "creation_time": "2014-08-26 08:54:56",
                "language": "und",
                "handler_name": "Core Media Data Handler",
                "encoder": "H.264"
            }
        }
    ],
    "format": {
        "filename": "test.mov",
        "nb_streams": 2,
        "nb_programs": 0,
        "format_name": "mov,mp4,m4a,3gp,3g2,mj2",
        "format_long_name": "QuickTime / MOV",
        "start_time": "0.000000",
        "duration": "7.023333",
        "size": "4167866",
        "bit_rate": "4747450",
        "probe_score": 100,
        "tags": {
            "major_brand": "qt  ",
            "minor_version": "0",
            "compatible_brands": "qt  ",
            "creation_time": "2014-08-26 08:54:56",
            "make": "Apple",
            "make-heb": "Apple",
            "encoder": "7.1",
            "encoder-heb": "7.1",
            "date": "2014-08-10T19:34:28+0300",
            "date-heb": "2014-08-10T19:34:28+0300",
            "location": "+32.8176+034.9998+019.000/",
            "location-heb": "+32.8176+034.9998+019.000/",
            "model": "iPhone 4S",
            "model-heb": "iPhone 4S"
        }
    }
}

		`
		probe := new(Probe)
		So(json.Unmarshal([]byte(data), probe), ShouldBeNil)
		video := probe.Stream("video")
		audio := probe.Stream("audio")
		So(video, ShouldNotBeNil)
		So(audio, ShouldNotBeNil)
		So(video.Tag("rotate"), ShouldEqual, "90")
		So(video.CodecName, ShouldEqual, "h264")
	})
}
