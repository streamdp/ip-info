package config

import (
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
		wantErr     bool
	}{
		{
			name: "regular loading config",
			envs: map[string]string{
				"GRPC_USE_REFLECTION":    "false",
				"IP_INFO_DATABASE_URL":   "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_RATE_LIMIT":     "10",
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
				CacheProvider:         "redis",
			},
			wantRedis: &Redis{
				Host:     redisDefaultHost,
				Port:     redisDefaultPort,
				Password: "",
				Db:       redisDefaultDb,
			},
			wantLimiter: &Limiter{
				RateLimit: 10,
			},
			wantCache: &Cache{
				TTL: 3600,
			},
			wantErr: false,
		},
		{
			name: "wrong rate limit",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL": "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_RATE_LIMIT":   "wrong",
			},
			wantErr: true,
		},
		{
			name: "wrong rate limit",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL": "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_RATE_LIMIT":   "wrong",
			},
			wantErr: true,
		},
		{
			name: "negative rate limit",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL":   "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_RATE_LIMIT":     "-1",
				"IP_INFO_ENABLE_LIMITER": "true",
			},
			wantErr: true,
		},
		{
			name: "wrong cache ttl",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL": "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_CACHE_TTL":    "wrong",
			},
			wantErr: true,
		},
		{
			name: "negative cache ttl",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL": "postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable",
				"IP_INFO_CACHE_TTL":    "-1",
			},
			wantErr: true,
		},
		{
			name: "empty db url",
			envs: map[string]string{
				"IP_INFO_DATABASE_URL": "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}
			gotApp, gotRedis, gotLimiter, gotCache, err := LoadConfig()
			if (err != nil) != tt.wantErr {
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
