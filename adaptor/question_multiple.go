// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

const (
	TypeText  = `text`
	TypeImage = `image_url`
	TypeAudio = `input_audio`
	TypeVideo = `video_url`
)

type ImageUrl struct {
	Url    string `json:"url"`
	Detail string `json:"detail,omitzero"`
}

type InputAudio struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

type VedioUrl struct {
	Url string `json:"url"`
}

type QuestionMultiple []struct {
	Type       string     `json:"type"`
	Text       string     `json:"text,omitzero"`
	ImageUrl   ImageUrl   `json:"image_url,omitzero"`
	InputAudio InputAudio `json:"input_audio,omitzero"`
	VedioUrl   VedioUrl   `json:"vedio_url,omitzero"`
}
