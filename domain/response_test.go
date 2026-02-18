package domain

import (
	"errors"
	"reflect"
	"testing"
)

var errCommon = errors.New("error")

func TestNewResponse(t *testing.T) {
	type args struct {
		err     error
		content any
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "get response",
			args: args{
				err:     nil,
				content: "content",
			},
			want: &Response{
				Err:     "",
				Content: "content",
			},
		},
		{
			name: "get response with err",
			args: args{
				err:     errCommon,
				content: "content",
			},
			want: &Response{
				Err:     "error",
				Content: "content",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResponse(tt.args.err, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_Bytes(t *testing.T) {
	tests := []struct {
		name string
		resp *Response
		want []byte
	}{
		{
			name: "marshal response struct",
			resp: &Response{
				Err:     "error",
				Content: "content",
			},
			want: []byte("{\n  \"error\": \"error\",\n  \"content\": \"content\"\n}"),
		},
		{
			name: "marshal empty response struct",
			resp: &Response{},
			want: []byte("{\n  \"error\": \"\",\n  \"content\": null\n}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.resp.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
