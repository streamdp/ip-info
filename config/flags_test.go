package config

import (
	"errors"
	"flag"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		wantApp *App
		wantErr error
	}{
		{
			name: "regular loading config",
			envs: map[string]string{
				"GRPC_USE_REFLECTION":    "false",
				"IP_INFO_DATABASE_URL":   "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_RATE_LIMIT":     "10",
				"IP_INFO_RATE_LIMIT_TTL": "65",
				"IP_INFO_CACHER":         "redis",
				"IP_INFO_CACHE_TTL":      "3600",
				"IP_INFO_ENABLE_LIMITER": "true",
			},
			wantApp: &App{
				Http: &Http{
					port:                    httpServerDefaultPort,
					serverReadTimeout:       httpServerDefaultTimeout,
					serverReadHeaderTimeout: httpServerDefaultTimeout,
					serverWriteTimeout:      httpServerDefaultTimeout,
					clientTimeout:           httpClientDefaultTimeout,
				},
				Grpc: &Grpc{
					port:          gRPCServerDefaultPort,
					readTimeout:   gRPCServerDefaultTimeout,
					useReflection: false,
				},
				Limiter: &Limiter{
					limiter:   limiterDefaultLimiter,
					rateLimit: limiterDefaultRateLimit,
					ttl:       65,
					enabled:   true,
				},
				Cache: &Cache{
					cacher:   "redis",
					ttl:      cacheDefaultTTL,
					disabled: false,
				},
				Redis: &Redis{
					host:     redisDefaultHost,
					port:     redisDefaultPort,
					Password: "",
					db:       redisDefaultDb,
				},

				DatabaseUrl: "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				Version:     "",
			},
			wantErr: nil,
		},
		{
			name: "wrong rate limit",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL":   "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_ENABLE_LIMITER": "true",
				"IP_INFO_RATE_LIMIT":     "wrong",
			},
			wantErr: errWrongRateLimit,
		},
		{
			name: "set wrong limiter evn",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL":   "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_ENABLE_LIMITER": "true",
				"IP_INFO_LIMITER":        "wrong",
			},
			wantErr: errWrongLimiter,
		},
		{
			name: "negative rate limit",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL":   "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_RATE_LIMIT":     "-1",
				"IP_INFO_ENABLE_LIMITER": "true",
			},
			wantErr: errWrongRateLimit,
		},
		{
			name: "wrong cache ttl",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL": "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_CACHE_TTL":    "wrong",
			},
			wantErr: errCacheTTL,
		},
		{
			name: "negative cache ttl",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL": "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_CACHE_TTL":    "-1",
			},
			wantErr: errCacheTTL,
		},
		{
			name: "empty db url environment variable",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL": "",
			},
			wantErr: errEmptyDatabaseUrlEnv,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}
			gotApp, err := LoadConfig()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if tt.wantErr == nil {
				if !reflect.DeepEqual(gotApp, tt.wantApp) {
					t.Errorf("LoadConfig() gotApp = %v, want %v", gotApp, tt.wantApp)
				}
			}
		})
	}
}
