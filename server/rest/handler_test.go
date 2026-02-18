package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/pkg/iplocator"
	"github.com/streamdp/ip-info/server"
)

var errCommon = errors.New("some_error")

func Test_httpClientIp(t *testing.T) {
	createRequestWithHeader := func(header string, value string) *http.Request {
		r := &http.Request{
			Header: http.Header{},
		}
		r.Header.Add(header, value)

		return r
	}

	tests := []struct {
		name    string
		request *http.Request
		want    string
	}{
		{
			name:    "get ip from cf-connecting-ip header",
			request: createRequestWithHeader(iplocator.CfConnectingIp, "127.0.0.1"),
			want:    "127.0.0.1",
		},
		{
			name:    "get ip from x-forwarded-for header",
			request: createRequestWithHeader(iplocator.XForwardedFor, "82.12.32.1"),
			want:    "82.12.32.1",
		},
		{
			name:    "get ip from x-real-ip header",
			request: createRequestWithHeader(iplocator.XRealIp, "8.8.8.8"),
			want:    "8.8.8.8",
		},
		{
			name: "get ip from request remoteAddr",
			request: &http.Request{
				RemoteAddr: "123.12.21.3:8080",
			},
			want: "123.12.21.3",
		},
		{
			name: "ipv6 placed in brackets with port",
			request: &http.Request{
				RemoteAddr: "[::1]:8080",
			},
			want: "::1",
		},
		{
			name: "ipv6 in brackets without port",
			request: &http.Request{
				RemoteAddr: "[::1]",
			},
			want: "::1",
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
			err:  server.ErrRateLimitExceeded,
			want: http.StatusTooManyRequests,
		},
		{
			name: "get http.StatusBadRequest",
			err:  server.ErrWrongIpAddress,
			want: http.StatusBadRequest,
		},
		{
			name: "get http.StatusNotFound",
			err:  database.ErrNoIpAddress,
			want: http.StatusNotFound,
		},
		{
			name: "get http.StatusInternalServerError",
			err:  errCommon,
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

func TestServer_healthz(t *testing.T) {
	handler := (&Server{}).healthz()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	handler(w, r)

	res := w.Result()
	t.Cleanup(func() { _ = res.Body.Close() })

	if res.StatusCode != http.StatusOK {
		t.Errorf("healthz() = %d, want %d", res.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read body: expected no error, got: %v", err)
	}

	if string(body) != "ok" {
		t.Fatalf("expected \"ok\", got: %v", string(body))
	}
}

func TestServer_ipInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		ip             string
		locator        server.Locator
		useClientIp    bool
		wantStatusCode int
		wantError      error
	}{
		{
			name:           "wrong ip address",
			ip:             "8.8.8.A",
			locator:        &mockLocator{err: server.ErrWrongIpAddress},
			wantStatusCode: http.StatusBadRequest,
			wantError:      server.ErrWrongIpAddress,
		},
		{
			name:           "get ip info",
			ip:             "8.8.8.8",
			locator:        &mockLocator{ipInfo: &domain.IpInfo{Ip: net.ParseIP("8.8.8.8")}},
			wantStatusCode: http.StatusOK,
			wantError:      nil,
		},
		{
			name:           "get client ip info",
			ip:             "127.0.0.1",
			locator:        &mockLocator{ipInfo: &domain.IpInfo{Ip: net.ParseIP("127.0.0.1")}},
			useClientIp:    true,
			wantStatusCode: http.StatusOK,
			wantError:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &Server{
				locator: tt.locator,
				l:       log.New(io.Discard, "", log.LstdFlags),
			}

			handler := s.ipInfo(tt.useClientIp)

			w := httptest.NewRecorder()

			var r *http.Request
			if tt.useClientIp {
				r = httptest.NewRequest(http.MethodGet, "/client-ip", nil)
				r.Header.Set(iplocator.XRealIp, tt.ip)
			} else {
				r = httptest.NewRequest(http.MethodGet, "/ip-info?ip="+tt.ip, nil)
			}

			handler(w, r)

			res := w.Result()
			t.Cleanup(func() { _ = res.Body.Close() })

			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("ipInfo() = %d, want %d", res.StatusCode, tt.wantStatusCode)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("read body: expected no error, got: %v", err)
			}

			resp := domain.Response{}
			_ = json.Unmarshal(body, &resp)

			if tt.wantError != nil {
				if !strings.Contains(resp.Err, tt.wantError.Error()) {
					t.Fatalf("unexcpected error: want %v, got: %v", tt.wantError, resp.Err)
				}
			} else {
				if resp.Err != "" {
					t.Fatalf("response contain error: expected no error, got: %v", resp.Err)
				}

				content, ok := resp.Content.(map[string]any)
				if !ok {
					t.Fatalf("failed to get response content")

					return
				}

				if content["ip"] != tt.ip {
					t.Fatalf("expected \"%s\", got: %v", tt.ip, content["ip"])
				}
			}
		})
	}
}

func TestServer_version(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		appVersion string
	}{
		{
			name:       "version v0.0.1",
			appVersion: "v0.0.1",
		},
		{
			name:       "version v0.0.2",
			appVersion: "v0.0.2",
		},
		{
			name:       "empty version",
			appVersion: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := (&Server{appVersion: tt.appVersion}).version()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/app/version", nil)

			handler(w, r)

			res := w.Result()
			t.Cleanup(func() { _ = res.Body.Close() })

			if res.StatusCode != http.StatusOK {
				t.Errorf("version() = %d, want %d", res.StatusCode, http.StatusOK)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("read body: expected no error, got: %v", err)
			}

			out := domain.Response{}
			if err = json.Unmarshal(body, &out); err != nil {
				t.Fatalf("unmarshal body: expected no error, got: %v", err)
			}

			if content, ok := out.Content.(map[string]any); !ok || content["version"] != tt.appVersion {
				t.Fatalf("expected \"version\":\"%s\", got: %v", tt.appVersion, out.Content)
			}
		})
	}
}

type mockLocator struct {
	ipInfo *domain.IpInfo
	err    error
}

func (ml *mockLocator) GetIpInfo(_ context.Context, _ string) (*domain.IpInfo, error) {
	if ml.err != nil {
		return nil, ml.err
	}

	return ml.ipInfo, nil
}
