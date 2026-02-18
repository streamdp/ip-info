package config

import (
	"errors"
	"testing"
)

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *App
		wantErr error
	}{
		{
			name: "app config is valid",
			cfg: &App{
				Http:    newHttpConfig(),
				Grpc:    newGrpcConfig(),
				Limiter: newLimiterConfig(),
				Cache:   newCacheConfig(),
				Redis:   newRedisConfig(),
				Database: &Database{
					url:            "postgres://postgres:postgres@localhost:5432/postgres",
					requestTimeout: 0,
				},
				version: "",
			},
			wantErr: nil,
		},
		{
			name: "wrong http port",
			cfg: &App{
				Http: &Http{
					port: -8080,
				},
				Grpc:    newGrpcConfig(),
				Limiter: newLimiterConfig(),
				Cache:   newCacheConfig(),
				Redis:   newRedisConfig(),
				Database: &Database{
					url:            "postgres://postgres:postgres@localhost:5432/postgres",
					requestTimeout: 0,
				},
				version: "",
			},
			wantErr: errWrongNetworkPort,
		},
		{
			name: "wrong grpc port",
			cfg: &App{
				Http: newHttpConfig(),
				Grpc: &Grpc{
					port: -50051,
				},
				Limiter: newLimiterConfig(),
				Cache:   newCacheConfig(),
				Redis:   newRedisConfig(),
				Database: &Database{
					url:            "postgres://postgres:postgres@localhost:5432/postgres",
					requestTimeout: 0,
				},
				version: "",
			},
			wantErr: errWrongNetworkPort,
		},
		{
			name: "wrong database url",
			cfg: &App{
				Http:    newHttpConfig(),
				Grpc:    newGrpcConfig(),
				Limiter: newLimiterConfig(),
				Cache:   newCacheConfig(),
				Redis:   newRedisConfig(),
				Database: &Database{
					url:            "",
					requestTimeout: 0,
				},
				version: "",
			},
			wantErr: errEmptyDatabaseUrl,
		},
		{
			name: "wrong redis config",
			cfg: &App{
				Http:    newHttpConfig(),
				Grpc:    newGrpcConfig(),
				Limiter: newLimiterConfig(),
				Cache:   newCacheConfig(),
				Redis: &Redis{
					host:     "",
					port:     6379,
					password: "",
					db:       0,
				},
				Database: &Database{
					url:            "postgres://postgres:postgres@localhost:5432/postgres",
					requestTimeout: 0,
				},
				version: "",
			},
			wantErr: errRedisHost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_Version(t *testing.T) {
	tests := []struct {
		name   string
		appCfg *App
		want   string
	}{
		{
			name:   "version v0.0.1",
			appCfg: &App{version: "v0.0.1"},
			want:   "v0.0.1",
		},
		{
			name:   "version v0.0.2",
			appCfg: &App{version: "v0.0.2"},
			want:   "v0.0.2",
		},
		{
			name:   "empty version",
			appCfg: &App{},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.appCfg.Version(); got != tt.want {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
		})
	}
}
