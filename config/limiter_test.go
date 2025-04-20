package config

import (
	"errors"
	"testing"
	"time"
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
				limiter:   "golimiter",
				rateLimit: 100,
				ttl:       60,
			},
			wantErr: nil,
		},
		{
			name: "wrong rate limit",
			cfg: &Limiter{
				limiter:   "redis_rate",
				rateLimit: -1,
				ttl:       60,
			},
			wantErr: errWrongRateLimit,
		},
		{
			name: "wrong limiter",
			cfg: &Limiter{
				limiter:   "redis_rate111",
				rateLimit: 2,
				ttl:       60,
			},
			wantErr: errWrongLimiter,
		},
		{
			name: "wrong ttl",
			cfg: &Limiter{
				limiter:   "redis_rate",
				rateLimit: 2,
				ttl:       0,
			},
			wantErr: errRateLimitTTL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.validate(); err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLimiter_Enabled(t *testing.T) {
	tests := []struct {
		name string
		l    *Limiter
		want bool
	}{
		{
			name: "not enabled",
			l: &Limiter{
				enabled: false,
			},
			want: false,
		},
		{
			name: "enabled",
			l: &Limiter{
				enabled: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Enabled(); got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimiter_Limiter(t *testing.T) {
	tests := []struct {
		name string
		l    *Limiter
		want string
	}{
		{
			name: "redis limiter",
			l: &Limiter{
				limiter: "redis",
			},
			want: "redis",
		},
		{
			name: "golimiter limiter",
			l: &Limiter{
				limiter: "golimiter",
			},
			want: "golimiter",
		},
		{
			name: "empty",
			l: &Limiter{
				limiter: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Limiter(); got != tt.want {
				t.Errorf("Limiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimiter_RateLimit(t *testing.T) {
	tests := []struct {
		name string
		l    *Limiter
		want int
	}{
		{
			name: "rate limit is 100",
			l: &Limiter{
				rateLimit: 100,
			},
			want: 100,
		},
		{
			name: "rate limit is 0",
			l: &Limiter{
				rateLimit: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.RateLimit(); got != tt.want {
				t.Errorf("RateLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimiter_Ttl(t *testing.T) {
	tests := []struct {
		name string
		l    *Limiter
		want time.Duration
	}{
		{
			name: "ttl is one hour",
			l: &Limiter{
				ttl: 3600,
			},
			want: time.Hour,
		},
		{
			name: "ttl is 36 seconds",
			l: &Limiter{
				ttl: 36,
			},
			want: 36 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Ttl(); got != tt.want {
				t.Errorf("Ttl() = %v, want %v", got, tt.want)
			}
		})
	}
}
