package config

import (
	"reflect"
	"testing"
	"time"
)

func TestDatabase_Url(t *testing.T) {
	tests := []struct {
		name string
		d    *Database
		want string
	}{
		{
			name: "get url",
			d: &Database{
				url:            "postgres://postgres:postgres@localhost:5432/postgres",
				requestTimeout: databaseRequestTimeout,
			},
			want: "postgres://postgres:postgres@localhost:5432/postgres",
		},
		{
			name: "empty url",
			d: &Database{
				url:            "",
				requestTimeout: databaseRequestTimeout,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Url(); got != tt.want {
				t.Errorf("Url() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_RequestTimeout(t *testing.T) {
	tests := []struct {
		name string
		d    *Database
		want time.Duration
	}{
		{
			name: "get 5 second",
			d: &Database{
				url:            "",
				requestTimeout: databaseRequestTimeout,
			},
			want: 5 * time.Second,
		},
		{
			name: "get one minute",
			d: &Database{
				url:            "",
				requestTimeout: 60000,
			},
			want: time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.RequestTimeout(); got != tt.want {
				t.Errorf("RequestTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_SetRequestTimeout(t *testing.T) {
	tests := []struct {
		name    string
		d       *Database
		timeout int
		want    time.Duration
	}{
		{
			name:    "set timeout to 5 seconds",
			timeout: 5000,
			want:    5 * time.Second,
		},
		{
			name:    "set timeout to 1 minute",
			timeout: 60000,
			want:    time.Minute,
		},
	}
	for _, tt := range tests {
		d := &Database{}
		t.Run(tt.name, func(t *testing.T) {
			if got := d.SetRequestTimeout(tt.timeout).RequestTimeout(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRequestTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}
