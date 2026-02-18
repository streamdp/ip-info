package domain

import (
	"encoding/json"
)

type Response struct {
	Err     string `json:"error"`
	Content any    `json:"content"`
}

func NewResponse(err error, content any) *Response {
	return &Response{
		Err: func(err error) string {
			if err != nil {
				return err.Error()
			}

			return ""
		}(err),
		Content: content,
	}
}

func (r *Response) Bytes() []byte {
	b, _ := json.MarshalIndent(r, "", "  ")

	return b
}
