package models

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"strings"
)

type Request struct {
	Id          bson.ObjectId `json:"id"`
	File        string        `json:"file"`
	Type        string        `json:"type"`
	ProgressKey string        `json:"progress_key"`
	ResultKey   string        `json:"result_key"`
	Options     interface{}   `json:"options"`
}

type VideoOptions struct {
	Start    int `json:"start"`
	End      int `json:"end"`
	Duration int `json:"duration"`
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

func (v *VideoOptions) String() string {
	var params []string
	params = append(params, fmt.Sprintf("-c:v %s", v.Video.Format))
	params = append(params, fmt.Sprintf("-b:v %d", v.Video.Bitrate))

	if v.Audio.Format == "aac" {
		params = append(params, "-strict")
		params = append(params, "-2")
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
		params = append(params, fmt.Sprintf("-vf crop=ih:ih,scale=%d:%d", v.Video.Width, v.Video.Height))
	}

	return strings.Join(params, " ")
}

func (a *AudioOptions) String() string {
	var params []string

	if a.Format == "aac" {
		params = append(params, "-strict")
		params = append(params, "-2")
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
