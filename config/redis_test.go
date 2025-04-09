package config

import (
	"errors"
	"reflect"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestRedis_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Redis
		wantErr error
	}{
		{
			name: "redis config is valid",
			cfg: &Redis{
				Host:     "127.0.0.1",
				Port:     6349,
				Password: "qwerty",
				Db:       1,
			},
			wantErr: nil,
		},
		{
			name: "wrong host",
			cfg: &Redis{
				Host:     "",
				Port:     6379,
				Password: "qwerty",
				Db:       1,
			},
			wantErr: errRedisHost,
		},
		{
			name: "wrong port",
			cfg: &Redis{
				Host:     "127.0.0.1",
				Port:     -1,
				Password: "qwerty",
				Db:       1,
			},
			wantErr: errWrongNetworkPort,
		},
		{
			name: "wrong db",
			cfg: &Redis{
				Host:     "127.0.0.1",
				Port:     6379,
				Password: "qwerty",
				Db:       150,
			},
			wantErr: errRedisDb,
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

func TestRedis_Options(t *testing.T) {
	tests := []struct {
		name    string
		l       *Redis
		envs    map[string]string
		want    *redis.Options
		wantErr error
	}{
		{
			name: "regular get options from config",
			l: &Redis{
				Host:     "127.0.0.1",
				Port:     6379,
				Password: "qwerty",
				Db:       1,
			},
			want:    &redis.Options{Addr: "127.0.0.1:6379", Password: "qwerty", DB: 1},
			wantErr: nil,
		},
		{
			name: "get options from env REDIS_URL",
			envs: map[string]string{
				"REDIS_URL": "redis://:qwerty@redis:6379/0",
			},
			want:    &redis.Options{Network: "tcp", Addr: "redis:6379", Password: "qwerty", DB: 0},
			wantErr: nil,
		},
		{
			name: "get options from separated envs",
			l:    &Redis{},
			envs: map[string]string{
				"REDIS_HOSTNAME": "127.0.0.1",
				"REDIS_PORT":     "6379",
				"REDIS_PASSWORD": "qwerty",
				"REDIS_DB":       "0",
			},
			want:    &redis.Options{Addr: "127.0.0.1:6379", Password: "qwerty", DB: 0},
			wantErr: nil,
		},
		{
			name: "err config not initialized",
			envs: map[string]string{
				"REDIS_HOSTNAME": "127.0.0.1",
			},
			want:    nil,
			wantErr: errConfigNotInitialized,
		},
		{
			name: "err parse redis port",
			l:    &Redis{},
			envs: map[string]string{
				"REDIS_PORT": "six",
			},
			want:    nil,
			wantErr: errWrongNetworkPort,
		},
		{
			name: "err parse db index",
			l:    &Redis{},
			envs: map[string]string{
				"REDIS_DB": "seven",
			},
			want:    nil,
			wantErr: errRedisDb,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}
			got, err := tt.l.Options()
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Options() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() got = %v, want %v", got, tt.want)
			}
		})
	}
}
