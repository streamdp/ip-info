package config

import (
	"errors"
	"testing"
)

func TestLimiter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Limiter
		wantErr error
	}{
		{
			name: "limiter config is valid",
			cfg: &Limiter{
				Limiter:   "golimiter",
				RateLimit: 100,
				TTL:       60,
			},
			wantErr: nil,
		},
		{
			name: "wrong rate limit",
			cfg: &Limiter{
				Limiter:   "redis_rate",
				RateLimit: -1,
				TTL:       60,
			},
			wantErr: errWrongRateLimit,
		},
		{
			name: "wrong limiter",
			cfg: &Limiter{
				Limiter:   "redis_rate111",
				RateLimit: 2,
				TTL:       60,
			},
			wantErr: errWrongLimiter,
		},
		{
			name: "wrong TTL",
			cfg: &Limiter{
				Limiter:   "redis_rate",
				RateLimit: 2,
				TTL:       0,
			},
			wantErr: errRateLimitTTL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
