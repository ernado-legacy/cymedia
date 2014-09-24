package models

type Stream struct {
	Index     int               `json:"index"`
	CodecName string            `json:"codec_name"`
	CodecType string            `json:"codec_type"`
	Tags      map[string]string `json:"tags,omitempty"`
	Duration  int64             `json:"duration_ts"`
	Width     int               `json:"width,omitempty"`
	Height    int               `json:"height,omitempty"`
}

type Probe struct {
	Streams []Stream `json:"streams,omitempty"`
}

func (p Probe) Stream(title string) *Stream {
	if len(p.Streams) == 0 {
		return nil
	}

	for _, stream := range p.Streams {
		if stream.CodecType == title {
			return &stream
		}
	}
	return nil
}

func (s Stream) Tag(title string) (value string) {
	if s.Tags == nil {
		return
	}
	return s.Tags[title]
}
