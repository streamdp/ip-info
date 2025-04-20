package config

import (
	"testing"
	"time"
)

func TestGrpc_ReadTimeout(t *testing.T) {
	tests := []struct {
		name string
		g    *Grpc
		want time.Duration
	}{
		{
			name: "5 second",
			g: &Grpc{
				readTimeout: 5000,
			},
			want: 5 * time.Second,
		},
		{
			name: "one hour",
			g: &Grpc{
				readTimeout: 60000,
			},
			want: 1 * time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.ReadTimeout(); got != tt.want {
				t.Errorf("ReadTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpc_UseReflection(t *testing.T) {
	tests := []struct {
		name string
		g    *Grpc
		want bool
	}{
		{
			name: "reflection is not used",
			g: &Grpc{
				useReflection: false,
			},
			want: false,
		},
		{
			name: "use reflection",
			g: &Grpc{
				useReflection: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.UseReflection(); got != tt.want {
				t.Errorf("UseReflection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpc_Port(t *testing.T) {
	tests := []struct {
		name string
		g    *Grpc
		want int
	}{
		{
			name: "port is 100",
			g: &Grpc{
				port: 100,
			},
			want: 100,
		},
		{
			name: "port is 1000",
			g: &Grpc{
				port: 1000,
			},
			want: 1000,
		},
		{
			name: "port is 0",
			g: &Grpc{
				port: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Port(); got != tt.want {
				t.Errorf("Port() = %v, want %v", got, tt.want)
			}
		})
	}
}
