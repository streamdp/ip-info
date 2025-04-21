package config

import (
	"testing"
	"time"
)

func TestHttp_HttpClientTimeout(t *testing.T) {
	tests := []struct {
		name string
		h    *Http
		want time.Duration
	}{
		{
			name: "get client timeout",
			h: &Http{
				clientTimeout: 100,
			},
			want: 100 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.ClientTimeout(); got != tt.want {
				t.Errorf("ClientTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttp_ServerWriteTimeout(t *testing.T) {
	tests := []struct {
		name string
		h    *Http
		want time.Duration
	}{
		{
			name: "get server write timeout",
			h: &Http{
				serverWriteTimeout: 100,
			},
			want: 100 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.ServerWriteTimeout(); got != tt.want {
				t.Errorf("ServerWriteTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttp_ServerReadHeaderTimeout(t *testing.T) {
	tests := []struct {
		name string
		h    *Http
		want time.Duration
	}{
		{
			name: "get server read header timeout",
			h: &Http{
				serverReadHeaderTimeout: 100,
			},
			want: 100 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.ServerReadHeaderTimeout(); got != tt.want {
				t.Errorf("ServerReadHeaderTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttp_ServerReadTimeout(t *testing.T) {
	tests := []struct {
		name string
		h    *Http
		want time.Duration
	}{
		{
			name: "get server read timeout",
			h: &Http{
				serverReadTimeout: 100,
			},
			want: 100 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.ServerReadTimeout(); got != tt.want {
				t.Errorf("ReadTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttp_Port(t *testing.T) {
	tests := []struct {
		name string
		h    *Http
		want int
	}{
		{
			name: "port is 100",
			h: &Http{
				port: 100,
			},
			want: 100,
		},
		{
			name: "port is 1000",
			h: &Http{
				port: 1000,
			},
			want: 1000,
		},
		{
			name: "port is 0",
			h: &Http{
				port: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Port(); got != tt.want {
				t.Errorf("Port() = %v, want %v", got, tt.want)
			}
		})
	}
}
