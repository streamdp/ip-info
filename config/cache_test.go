package config

import (
	"testing"
)

func TestCache_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Cache
		wantErr bool
	}{
		{
			name: "cache config is valid",
			cfg: &Cache{
				Provider: "microcache",
				TTL:      100,
			},
			wantErr: false,
		},
		{
			name: "wrong ttl",
			cfg: &Cache{
				Provider: "redis",
				TTL:      -1,
			},
			wantErr: true,
		},
		{
			name: "wrong provider",
			cfg: &Cache{
				Provider: "redis111",
				TTL:      2,
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
