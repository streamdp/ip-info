package config

import (
	"errors"
	"testing"
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
				Cacher: "microcache",
				TTL:    100,
			},
			wantErr: nil,
		},
		{
			name: "wrong ttl",
			cfg: &Cache{
				Cacher: "redis",
				TTL:    -1,
			},
			wantErr: errCacheTTL,
		},
		{
			name: "wrong cacher",
			cfg: &Cache{
				Cacher: "redis111",
				TTL:    2,
			},
			wantErr: errWrongCache,
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
