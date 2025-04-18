package config

import (
	"errors"
	"flag"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		envs        map[string]string
		wantApp     *App
		wantRedis   *Redis
		wantLimiter *Limiter
		wantCache   *Cache
		wantErr     error
	}{
		{
			name: "regular loading config",
			envs: map[string]string{
				"GRPC_USE_REFLECTION":    "false",
				"IP_INFO_DATABASE_URL":   "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_RATE_LIMIT":     "10",
				"IP_INFO_RATE_LIMIT_TTL": "65",
				"IP_INFO_CACHE_PROVIDER": "redis",
				"IP_INFO_CACHE_TTL":      "3600",
				"IP_INFO_ENABLE_LIMITER": "true",
			},
			wantApp: &App{
				HttpPort:              httpServerDefaultPort,
				GrpcPort:              gRPCServerDefaultPort,
				GrpcUseReflection:     false,
				DatabaseUrl:           "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				GrpcReadTimeout:       serverDefaultTimeout,
				HttpReadTimeout:       serverDefaultTimeout,
				HttpReadHeaderTimeout: serverDefaultTimeout,
				HttpWriteTimeout:      serverDefaultTimeout,
				Version:               "",
				EnableLimiter:         true,
				DisableCache:          false,
			},
			wantRedis: &Redis{
				Host:     redisDefaultHost,
				Port:     redisDefaultPort,
				Password: "",
				Db:       redisDefaultDb,
			},
			wantLimiter: &Limiter{
				Provider:  "golimiter",
				RateLimit: 10,
				TTL:       65,
			},
			wantCache: &Cache{
				Provider: "redis",
				TTL:      3600,
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
			gotApp, gotRedis, gotLimiter, gotCache, err := LoadConfig()
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(gotApp, tt.wantApp) {
				t.Errorf("LoadConfig() gotApp = %v, want %v", gotApp, tt.wantApp)
			}
			if !reflect.DeepEqual(gotRedis, tt.wantRedis) {
				t.Errorf("LoadConfig() gotRedis = %v, want %v", gotRedis, tt.wantRedis)
			}
			if !reflect.DeepEqual(gotLimiter, tt.wantLimiter) {
				t.Errorf("LoadConfig() gotLimiter = %v, want %v", gotLimiter, tt.wantLimiter)
			}
			if !reflect.DeepEqual(gotCache, tt.wantCache) {
				t.Errorf("LoadConfig() gotCache = %v, want %v", gotCache, tt.wantCache)
			}
		})
	}
}
