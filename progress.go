package cymedia

import (
	"io"
)

type Progressor struct {
	Length   int64
	Rate     int64
	Reader   *io.PipeReader
	Writer   *io.PipeWriter
	Progress chan float32
}

func (p *Progressor) Start() {
	defer p.Writer.Close()
	var total int64
	bufLen := p.Length * 1. / p.Rate
	for {
		buffer := make([]byte, bufLen)
		read, err := p.Reader.Read(buffer)
		if err == io.EOF {
			break
		}
		total += int64(read)
		p.Progress <- float32(total) / float32(p.Length) * 100
	}
	p.Progress <- 100.0
	close(p.Progress)
}

func Progress(length int64, rate int64, progress chan float32) *io.PipeWriter {
	progressReader, progressWriter := io.Pipe()
	p := Progressor{length, rate, progressReader, progressWriter, progress}
	go p.Start()
	return progressWriter
}
