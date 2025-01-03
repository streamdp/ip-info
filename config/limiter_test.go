package config

import (
	"testing"
)

func TestLimiter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Limiter
		wantErr bool
	}{
		{
			name: "limiter config is valid",
			cfg: &Limiter{
				RateLimit: 100,
			},
			wantErr: false,
		},
		{
			name: "wrong rate limit",
			cfg: &Limiter{
				RateLimit: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
