package domain

import (
	"net"
	"testing"
)

func TestIpInfo_String(t *testing.T) {
	tests := []struct {
		name string
		info *IpInfo
		want string
	}{
		{
			name: "marshal ip info struct",
			info: &IpInfo{
				Ip:        net.ParseIP("8.8.8.8"),
				Continent: "NA",
				Country:   "US",
				StateProv: "California",
				City:      "Mountain View",
				Latitude:  -122.085,
				Longitude: 37.4223,
			},
			want: "{\n  \"ip\": \"8.8.8.8\",\n  \"continent\": \"NA\",\n  \"country\": \"US\",\n  \"state_prov\": \"California\",\n  \"city\": \"Mountain View\",\n  \"latitude\": -122.085,\n  \"longitude\": 37.4223\n}",
		},
		{
			name: "marshal empty ip info struct",
			info: &IpInfo{},
			want: "{\n  \"ip\": \"\",\n  \"continent\": \"\",\n  \"country\": \"\",\n  \"state_prov\": \"\",\n  \"city\": \"\",\n  \"latitude\": 0,\n  \"longitude\": 0\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
