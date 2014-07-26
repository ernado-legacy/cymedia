package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ernado/cymedia/mediad/models"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Converter struct {
}

var (
	ErrBadFormat = errors.New("Bad media format")
	extensions   = map[string]string{"h264": "mp4", "libvpx": "webm", "libvorbis": "ogg", "aac": "aac"}
)

const fileNameLength = 12

func randStr(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (c *Converter) Convert(input io.Reader, options models.Options) (output io.ReadCloser, err error) {
	extension, ok := extensions[options.GetFormat()]
	if !ok {
		return output, ErrBadFormat
	}
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", randStr(fileNameLength), extension))
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("ffmpeg -i - %s %s", options, path))
	cmd.Stdin = input
	if err = cmd.Run(); err != nil {
		return
	}
	return os.Open(path)
}
