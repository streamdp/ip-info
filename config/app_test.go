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
				Http:        newHttpConfig(),
				Grpc:        newGrpcConfig(),
				Limiter:     newLimiterConfig(),
				Cache:       newCacheConfig(),
				Redis:       newRedisConfig(),
				DatabaseUrl: "postgres://postgres:postgres@localhost:5432/postgres",
				Version:     "",
			},
			wantErr: nil,
		},
		{
			name: "wrong http port",
			cfg: &App{
				Http: &Http{
					port: -8080,
				},
				Grpc:        newGrpcConfig(),
				Limiter:     newLimiterConfig(),
				Cache:       newCacheConfig(),
				Redis:       newRedisConfig(),
				DatabaseUrl: "postgres://postgres:postgres@localhost:5432/postgres",
				Version:     "",
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
				Limiter:     newLimiterConfig(),
				Cache:       newCacheConfig(),
				Redis:       newRedisConfig(),
				DatabaseUrl: "postgres://postgres:postgres@localhost:5432/postgres",
				Version:     "",
			},
			wantErr: errWrongNetworkPort,
		},
		{
			name: "wrong database url",
			cfg: &App{
				Http:        newHttpConfig(),
				Grpc:        newGrpcConfig(),
				Limiter:     newLimiterConfig(),
				Cache:       newCacheConfig(),
				Redis:       newRedisConfig(),
				DatabaseUrl: "",
				Version:     "",
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
					Password: "",
					db:       0,
				},
				DatabaseUrl: "postgres://postgres:postgres@localhost:5432/postgres",
				Version:     "",
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
