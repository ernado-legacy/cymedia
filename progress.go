package cymedia

import "io"

type progressor struct {
	Length   int64
	Rate     int64
	Reader   *io.PipeReader
	Progress chan float32
}

func (p *progressor) Start() {
	var total int64
	buffer := make([]byte, p.Length*1./p.Rate)
	p.Progress <- float32(0)
	for {
		read, err := p.Reader.Read(buffer)
		total += int64(read)
		p.Progress <- float32(total) / float32(p.Length) * 100
		if err != nil || total >= p.Length {
			break
		}
	}
	close(p.Progress)
}

func Progress(f io.Reader, length int64, rate int64, progress chan float32) io.Reader {
	progressReader, progressWriter := io.Pipe()
	p := progressor{length, rate, progressReader, progress}
	go p.Start()
	return io.TeeReader(f, progressWriter)
}
