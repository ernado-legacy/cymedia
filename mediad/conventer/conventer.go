package conventer

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ernado/cymedia/mediad/models"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Conventer interface {
	Convert(input io.Reader, options models.Options) (output io.ReadCloser, err error)
}

type VideoConventer struct{}

var (
	ErrBadFormat = errors.New("Bad media format")
)

const fileNameLength = 12

func randStr(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (c *VideoConventer) Probe(filename string) (*models.Probe, error) {
	probe := new(models.Probe)
	params := "-v quiet -print_format json -show_format -show_streams"
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("ffprobe %s %s", params, filename))
	buffer := new(bytes.Buffer)
	errBuffer := new(bytes.Buffer)
	cmd.Stdout = buffer
	cmd.Stderr = errBuffer
	if err := cmd.Run(); err != nil {
		log.Println(err)
		log.Println(cmd.Args)
		log.Println(buffer.String())
		return nil, err
	}
	decoder := json.NewDecoder(buffer)
	if err := decoder.Decode(probe); err != nil {
		return nil, err
	}
	return probe, nil
}

func (c *VideoConventer) Convert(input io.Reader, options models.Options) (output io.ReadCloser, err error) {
	extension := options.Extension()
	if extension == "" {
		return output, ErrBadFormat
	}
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", randStr(fileNameLength), extension))
	tempfile := filepath.Join(os.TempDir(), randStr(fileNameLength))
	f, err := os.Create(tempfile)
	if err != nil {
		log.Println(err)
		return
	}
	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := input.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			return output, err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := f.Write(buf[:n]); err != nil {
			log.Println(err)
			return output, err
		}
	}
	probe, err := c.Probe(tempfile)
	if err != nil {
		return output, err
	}
	if err := options.Process(probe); err != nil {
		return output, err
	}
	command := fmt.Sprintf("ffmpeg -i %s %s %s", tempfile, options, path)
	log.Println(command)
	cmd := exec.Command("/bin/bash", "-c", command)
	buffer := new(bytes.Buffer)
	cmd.Stdin = input
	cmd.Stderr = buffer
	if err = cmd.Run(); err != nil {
		if err.Error() == "write |1: broken pipe" {
			log.Println("ignoring broken pipe")
			return os.Open(path)
		}
		log.Println(err)
		log.Println(cmd.Args)
		log.Println(buffer.String())
		return
	}
	return os.Open(path)
}
