package conventer

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
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

func (c *VideoConventer) Convert(input io.Reader, options models.Options) (output io.ReadCloser, err error) {
	extension := options.Extension()
	if extension == "" {
		return output, ErrBadFormat
	}
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", randStr(fileNameLength), extension))
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("ffmpeg -i - %s %s", options, path))
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
