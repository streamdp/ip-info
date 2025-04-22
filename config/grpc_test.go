package config

import (
	"testing"
)

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
