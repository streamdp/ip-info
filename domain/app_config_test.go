package domain

import (
	"errors"
	"testing"
)

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *AppConfig
		wantErr error
	}{
		{
			name: "app config is valid",
			cfg: &AppConfig{
				HttpPort:    8080,
				GrpcPort:    8081,
				DatabaseUrl: "postgres://postgres:postgres@localhost:5432/postgres",
			},
			wantErr: nil,
		},
		{
			name: "wrong http port",
			cfg: &AppConfig{
				HttpPort:    -1,
				GrpcPort:    8081,
				DatabaseUrl: "postgres://postgres:postgres@localhost:5432/postgres",
			},
			wantErr: errWrongNetworkPort,
		},
		{
			name: "wrong grpc port",
			cfg: &AppConfig{
				HttpPort:    8080,
				GrpcPort:    -1,
				DatabaseUrl: "postgres://postgres:postgres@localhost:5432/postgres",
			},
			wantErr: errWrongNetworkPort,
		},
		{
			name: "wrong database url",
			cfg: &AppConfig{
				HttpPort:    8080,
				GrpcPort:    8081,
				DatabaseUrl: "",
			},
			wantErr: errEmptyDatabaseUrl,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
