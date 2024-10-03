package rest

import (
	"errors"
	"net/http"
	"testing"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/server"
)

func Test_httpClientIp(t *testing.T) {
	createRequestWithHeader := func(header string, value string) (r *http.Request) {
		r = &http.Request{
			Header: http.Header{},
		}
		r.Header.Add(header, value)

		return
	}

	tests := []struct {
		name    string
		request *http.Request
		want    string
	}{
		{
			name:    "get ip from cf-connecting-ip header",
			request: createRequestWithHeader(ip_locator.CfConnectingIp, "127.0.0.1"),
			want:    "127.0.0.1",
		},
		{
			name:    "get ip from x-forwarded-for header",
			request: createRequestWithHeader(ip_locator.XForwardedFor, "82.12.32.1"),
			want:    "82.12.32.1",
		},
		{
			name:    "get ip from x-real-ip header",
			request: createRequestWithHeader(ip_locator.XRealIp, "8.8.8.8"),
			want:    "8.8.8.8",
		},
		{
			name: "get ip from request remoteAddr",
			request: &http.Request{
				RemoteAddr: "123.12.21.3:8080",
			},
			want: "123.12.21.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := httpClientIp(tt.request); got != tt.want {
				t.Errorf("httpClientIp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getHttpStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "get http.StatusOk",
			err:  nil,
			want: http.StatusOK,
		},
		{
			name: "get http.StatusBadRequest",
			err:  ip_locator.ErrWrongIpAddress,
			want: http.StatusBadRequest,
		},
		{
			name: "get http.StatusNotFound",
			err:  database.ErrNoIpAddress,
			want: http.StatusNotFound,
		},
		{
			name: "get http.StatusInternalServerError",
			err:  errors.New("some_error"),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHttpStatus(tt.err); got != tt.want {
				t.Errorf("getHttpStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
