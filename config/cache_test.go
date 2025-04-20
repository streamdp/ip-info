package config

import (
	"errors"
	"testing"
	"time"
)

func TestCache_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Cache
		wantErr error
	}{
		{
			name: "cache config is valid",
			cfg: &Cache{
				cacher: "microcache",
				ttl:    100,
			},
			wantErr: nil,
		},
		{
			name: "wrong ttl",
			cfg: &Cache{
				cacher: "redis",
				ttl:    -1,
			},
			wantErr: errCacheTTL,
		},
		{
			name: "wrong cacher",
			cfg: &Cache{
				cacher: "redis111",
				ttl:    2,
			},
			wantErr: errWrongCache,
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

func TestCache_Cacher(t *testing.T) {
	tests := []struct {
		name string
		c    *Cache
		want string
	}{
		{
			name: "redis cacher",
			c: &Cache{
				cacher: "redis",
			},
			want: "redis",
		},
		{
			name: "microcache cacher",
			c: &Cache{
				cacher: "microcache",
			},
			want: "microcache",
		},
		{
			name: "empty cacher",
			c: &Cache{
				cacher: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Cacher(); got != tt.want {
				t.Errorf("Cacher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Enabled(t *testing.T) {
	tests := []struct {
		name string
		c    *Cache
		want bool
	}{
		{
			name: "enabled",
			c: &Cache{
				disabled: false,
			},
			want: true,
		},
		{
			name: "disabled",
			c: &Cache{
				disabled: true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Enabled(); got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Ttl(t *testing.T) {
	tests := []struct {
		name string
		c    *Cache
		want time.Duration
	}{
		{
			name: "5 seconds",
			c: &Cache{
				ttl: 5,
			},
			want: 5 * time.Second,
		},
		{
			name: "1 hour",
			c: &Cache{
				ttl: 3600,
			},
			want: time.Hour,
		},
		{
			name: "0",
			c: &Cache{
				ttl: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Ttl(); got != tt.want {
				t.Errorf("Cacher() = %v, want %v", got, tt.want)
			}
		})
	}
}
