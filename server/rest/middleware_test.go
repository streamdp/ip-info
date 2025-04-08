package rest

import (
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
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)

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
	tests := []struct {
		name           string
		r              *http.Request
		contentType    string
		wantStatusCode int
		wantError      bool
	}{
		{
			name:           "content type not implemented",
			r:              httptest.NewRequest(http.MethodGet, "/client-ip", nil),
			contentType:    "application/xml",
			wantStatusCode: http.StatusNotImplemented,
			wantError:      true,
		},
		{
			name:           "content type implemented",
			r:              httptest.NewRequest(http.MethodGet, "/ip-info", nil),
			contentType:    jsonContentType,
			wantStatusCode: http.StatusOK,
			wantError:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mw := contentTypeRestrictionMW(
				log.New(io.Discard, "", log.LstdFlags),
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			)

			w := httptest.NewRecorder()
			tt.r.Header.Set(contentTypeHeader, tt.contentType)

			mw.ServeHTTP(w, tt.r)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)

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

type mockLimiter struct {
	err error
}

func (ml *mockLimiter) Limit(_ string) error {
	return ml.err
}
