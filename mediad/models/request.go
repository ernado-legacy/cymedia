package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	extensions = map[string]string{"h264": "mp4", "libvpx": "webm",
		"libvorbis": "ogg", "aac": "m4a", "jpg": "jpg", "png": "png", "mp3": "mp3"}
)

type Options interface {
	String() string
	Mime() string
	Extension() string
	Process(probe *Probe) error
}

const (
	VideoType     = "video"
	AudioType     = "audio"
	ThumbnailType = "thumbnail"
	transposeFlag = "transpose=1"
)

type Request struct {
	Id          string      `json:"id"`
	File        string      `json:"file"`
	Type        string      `json:"type"`
	ProgressKey string      `json:"progress_key"`
	ResultKey   string      `json:"result_key"`
	Options     interface{} `json:"options"`
}

func convert(input interface{}, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}
func (r *Request) GetOptions() (options Options) {
	if r.Type == VideoType {
		o := new(VideoOptions)
		convert(r.Options, o)
		options = o
	}
	if r.Type == AudioType {
		o := new(AudioOptions)
		convert(r.Options, o)
		options = o
	}
	if r.Type == ThumbnailType {
		o := new(ThumbnailOptions)
		convert(r.Options, o)
		options = o
	}
	return
}

type Responce struct {
	Id      string `json:"id"`
	File    string `json:"file,omitempty"`
	Type    string `json:"type"`
	Format  string `json:"format,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type VideoOptions struct {
	Start    int  `json:"start"`
	End      int  `json:"end"`
	Duration int  `json:"duration"`
	Rotated  bool `json:"rotated,omitempty"`
	Video    struct {
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		Square  bool   `json:"square"`
		Format  string `json:"format"`
		Bitrate int    `json:"birtate"`
	} `json:"video"`

	Audio struct {
		Format  string `json:"format"`
		Bitrate int    `json:"birtate"`
	} `json:"audio"`
}

type AudioOptions struct {
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Duration int    `json:"duration"`
	Format   string `json:"format"`
	Bitrate  int    `json:"birtate"`
}

type PictureOptions struct {
	Format    string `json:"format"`
	Thumbnail bool   `json:"thumbnail"`
	Width     int    `json:"width,omitempty"`
	Heigth    int    `json:"heigth,omitempty"`
	Quality   int    `json:"quality,omitempty"`
}

type ThumbnailOptions struct {
	Format  string `json:"format"`
	Rotated bool   `json:"rotated,omitempty"`
	Time    int    `json:"time"`
}

func fixAAC(params []string) []string {
	params = append(params, "-strict")
	return append(params, "-2")
}

func (v *VideoOptions) Process(probe *Probe) error {
	video := probe.Stream("video")
	if video != nil {
		rotation := video.Tag("rotate")
		if rotation == "90" {
			v.Rotated = true
		}
	}
	return nil
}

func (t *ThumbnailOptions) Process(probe *Probe) error {
	video := probe.Stream("video")
	if video != nil {
		rotation := video.Tag("rotate")
		if rotation == "90" {
			t.Rotated = true
		}
	}
	return nil
}

func (a *AudioOptions) Process(probe *Probe) error {
	return nil
}

func (v *VideoOptions) String() string {
	var params []string
	params = append(params, fmt.Sprintf("-c:v %s", v.Video.Format))
	params = append(params, fmt.Sprintf("-b:v %d", v.Video.Bitrate))

	if v.Audio.Format == "aac" {
		params = fixAAC(params)
	}

	params = append(params, fmt.Sprintf("-c:a %s", v.Audio.Format))
	if v.Audio.Bitrate != 0 {
		params = append(params, fmt.Sprintf("-b:a %d", v.Audio.Bitrate))
	}
	if v.Start != 0 {
		params = append(params, fmt.Sprintf("-ss %d", v.Start))
	}
	if v.End != 0 {
		params = append(params, fmt.Sprintf("-to %d", v.End))
	}
	if v.Duration != 0 {
		params = append(params, fmt.Sprintf("-t %d", v.Duration))
	}
	if v.Video.Square {
		param := fmt.Sprintf("-vf crop=ih:ih,scale=%d:%d", v.Video.Width, v.Video.Height)
		if v.Rotated {
			param = fmt.Sprintf("-vf crop=ih:ih,scale=%d:%d,%s", v.Video.Width, v.Video.Height, transposeFlag)
		}
		params = append(params, param)
	} else if v.Rotated {
		params = append(params, fmt.Sprintf("-vf %s", transposeFlag))
	}

	return strings.Join(params, " ")
}

func (v *VideoOptions) GetFormat() string {
	return v.Video.Format
}

func (v *VideoOptions) Extension() string {
	return extensions[v.Video.Format]
}

func (v *VideoOptions) Mime() string {
	return fmt.Sprintf("video/%s", v.Extension())
}

func (a *AudioOptions) Extension() string {
	return extensions[a.Format]
}

func (a *AudioOptions) Mime() string {
	return fmt.Sprintf("audio/%s", a.Extension())
}

func (a *AudioOptions) String() string {
	var params []string

	if a.Format == "aac" {
		params = fixAAC(params)
	}

	params = append(params, fmt.Sprintf("-c:a %s", a.Format))
	if a.Bitrate != 0 {
		params = append(params, fmt.Sprintf("-b:a %d", a.Bitrate))
	}
	if a.Start != 0 {
		params = append(params, fmt.Sprintf("-ss %d", a.Start))
	}
	if a.End != 0 {
		params = append(params, fmt.Sprintf("-to %d", a.End))
	}
	if a.Duration != 0 {
		params = append(params, fmt.Sprintf("-t %d", a.Duration))
	}

	params = append(params, "-vn")

	return strings.Join(params, " ")
}

func (p *PictureOptions) Mime() string {
	return fmt.Sprintf("image/%s", p.Extension())
}

func (p *PictureOptions) Extension() string {
	return p.Format
}

func (p *PictureOptions) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (t *ThumbnailOptions) String() string {
	var params []string
	if t.Time != 0 {
		params = append(params, fmt.Sprintf("-ss %d", t.Time))
	}
	params = append(params, "-vframes 1")
	if t.Rotated {
		params = append(params, fmt.Sprintf("-vf %s", transposeFlag))
	}
	return strings.Join(params, " ")
}

func (t *ThumbnailOptions) Mime() string {
	return fmt.Sprintf("image/%s", t.Extension())
}

func (t *ThumbnailOptions) Extension() string {
	return t.Format
}
