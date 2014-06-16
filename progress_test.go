package cymedia

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProgress(t *testing.T) {
	Convey("Progress", t, func() {
		var length int64
		length = 1024 * 1024
		buffer := make([]byte, length)
		bufferReader := bytes.NewReader(buffer)
		progress := make(chan float32, 4)
		rate := int64(10)
		step := float64(100.0 / rate)
		reader := Progress(bufferReader, length, rate, progress)
		go func() {
			newBuffer := make([]byte, length)
			reader.Read(newBuffer)
		}()
		progressShould := float64(0.0)
		precision := 0.01
		for p := range progress {
			So(p, ShouldBeBetween, progressShould-precision*step, progressShould+precision*step)
			progressShould += step
		}
		So(true, ShouldBeTrue)
	})
}
