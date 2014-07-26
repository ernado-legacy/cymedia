package models

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVideoParams(t *testing.T) {
	Convey("Audio options", t, func() {
		Convey("Simple", func() {
			p := new(AudioOptions)
			p.Format = "aac"
			p.Bitrate = 1024 * 128
			expected := "-strict -2 -c:a aac -b:a 131072 -vn"
			actual := p.String()
			So(actual, ShouldEqual, expected)
		})
		Convey("Duration", func() {
			p := new(AudioOptions)
			p.Format = "aac"
			p.Bitrate = 1024 * 128
			p.Start = 5
			p.Duration = 10
			expected := "-strict -2 -c:a aac -b:a 131072 -ss 5 -t 10 -vn"
			actual := p.String()
			So(actual, ShouldEqual, expected)
		})
		Convey("End", func() {
			p := new(AudioOptions)
			p.Format = "aac"
			p.Bitrate = 1024 * 128
			p.Start = 5
			p.End = 10
			expected := "-strict -2 -c:a aac -b:a 131072 -ss 5 -to 10 -vn"
			actual := p.String()
			So(actual, ShouldEqual, expected)
		})
	})
	Convey("Video options", t, func() {
		Convey("Simple", func() {
			p := new(VideoOptions)
			p.Audio.Format = "aac"
			p.Audio.Bitrate = 1024 * 128
			p.Video.Format = "h264"
			p.Video.Bitrate = 1024 * 500
			expected := "-c:v h264 -b:v 512000 -strict -2 -c:a aac -b:a 131072"
			actual := p.String()
			So(actual, ShouldEqual, expected)
		})
		Convey("Crop", func() {
			p := new(VideoOptions)
			p.Audio.Format = "libvorbis"
			p.Video.Format = "libvpx"
			p.Video.Bitrate = 1024 * 500
			p.Video.Square = true
			p.Video.Height = 200
			p.Video.Width = 200
			expected := "-c:v libvpx -b:v 512000 -c:a libvorbis -vf crop=ih:ih,scale=200:200"
			actual := p.String()
			So(actual, ShouldEqual, expected)
		})
		Convey("Duration", func() {
			p := new(VideoOptions)
			p.Audio.Format = "aac"
			p.Audio.Bitrate = 1024 * 128
			p.Video.Format = "h264"
			p.Video.Bitrate = 1024 * 500
			p.Start = 5
			p.Duration = 10
			expected := "-c:v h264 -b:v 512000 -strict -2 -c:a aac -b:a 131072 -ss 5 -t 10"
			actual := p.String()
			So(actual, ShouldEqual, expected)
		})
		Convey("End", func() {
			p := new(VideoOptions)
			p.Audio.Format = "aac"
			p.Audio.Bitrate = 1024 * 128
			p.Video.Format = "h264"
			p.Video.Bitrate = 1024 * 500
			p.Start = 5
			p.End = 15
			expected := "-c:v h264 -b:v 512000 -strict -2 -c:a aac -b:a 131072 -ss 5 -to 15"
			actual := p.String()
			So(actual, ShouldEqual, expected)
		})
	})
}
