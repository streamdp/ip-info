package server

import (
	"testing"
)

func TestExtractIpAddress(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want string
	}{
		{
			name: "ip4 address without port",
			ip:   "123.12.21.3",
			want: "123.12.21.3",
		},
		{
			name: "ip4 address with port",
			ip:   "123.12.21.3:8080",
			want: "123.12.21.3",
		},
		{
			name: "ipv6 placed in brackets with port",
			ip:   "[::1]:8080",
			want: "::1",
		},
		{
			name: "ipv6 in brackets without port",
			ip:   "[::1]",
			want: "::1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractIpAddress(tt.ip); got != tt.want {
				t.Errorf("ExtractIpAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
