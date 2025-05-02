package rest

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/server"
)

func Test_rateLimiterMW(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		request        *http.Request
		limiter        server.Limiter
		wantStatusCode int
		wantError      bool
	}{
		{
			name:           "client has reached its limits",
			request:        httptest.NewRequest(http.MethodGet, "/client-ip", nil),
			limiter:        &mockLimiter{err: server.ErrRateLimitExceeded},
			wantStatusCode: http.StatusTooManyRequests,
			wantError:      true,
		},
		{
			name:           "client not limited",
			request:        httptest.NewRequest(http.MethodGet, "/ip-info", nil),
			limiter:        &mockLimiter{},
			wantStatusCode: http.StatusOK,
			wantError:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mw := rateLimiterMW(
				tt.limiter,
				log.New(io.Discard, "", log.LstdFlags),
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			)

			w := httptest.NewRecorder()

			mw.ServeHTTP(w, tt.request)

			res := w.Result()
			t.Cleanup(func() { _ = res.Body.Close() })

			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("rateLimiterMW() = %d, want %d", res.StatusCode, tt.wantStatusCode)
			}

			if tt.wantError {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("read body: expected no error, got: %v", err)
				}

				resp := domain.Response{}
				_ = json.Unmarshal(body, &resp)

				if resp.Err == "" {
					t.Fatalf("response doesn't contain error: expected error, got: \"\"")
				}
			}
		})
	}
}

func Test_contentTypeRestrictionMW(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		request            *http.Request
		allowedContentType string
		requestContentType string
		wantStatusCode     int
		wantError          bool
	}{
		{
			name:               "content type not implemented",
			request:            httptest.NewRequest(http.MethodGet, "/client-ip", nil),
			allowedContentType: jsonContentType,
			requestContentType: "application/xml",
			wantStatusCode:     http.StatusNotImplemented,
			wantError:          true,
		},
		{
			name:               "content type implemented",
			request:            httptest.NewRequest(http.MethodGet, "/ip-info", nil),
			allowedContentType: jsonContentType,
			requestContentType: jsonContentType,
			wantStatusCode:     http.StatusOK,
			wantError:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mw := contentTypeRestrictionMW(
				log.New(io.Discard, "", log.LstdFlags),
				func(w http.ResponseWriter, r *http.Request) {},
				tt.allowedContentType,
			)

			w := httptest.NewRecorder()
			tt.request.Header.Set(contentTypeHeader, tt.requestContentType)

			mw.ServeHTTP(w, tt.request)

			res := w.Result()
			t.Cleanup(func() {
				_ = res.Body.Close()
			})

			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("contentTypeRestrictionMW() = %d, want %d", res.StatusCode, tt.wantStatusCode)
			}

			if tt.wantError {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("read body: expected no error, got: %v", err)
				}

				resp := domain.Response{}
				_ = json.Unmarshal(body, &resp)

				if resp.Err == "" {
					t.Fatalf("response doesn't contain error: expected error, got: \"\"")
				}
			}
		})
	}
}

func Test_isAllowedContentType(t *testing.T) {
	type args struct {
		c            string
		allowedTypes []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "allowed type",
			args: args{
				jsonContentType,
				[]string{jsonContentType, textPlainContentType},
			},
			want: true,
		},
		{
			name: "content type is not allowed",
			args: args{
				"application/xml",
				[]string{jsonContentType},
			},
			want: false,
		},
		{
			name: "empty content type",
			args: args{
				"",
				[]string{jsonContentType},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAllowedContentType(tt.args.c, tt.args.allowedTypes); got != tt.want {
				t.Errorf("isAllowedContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockLimiter struct {
	err error
}

func (ml *mockLimiter) Limit(_ context.Context, _ string) error {
	return ml.err
}
