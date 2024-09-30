package grpc

import (
	"net"
	"reflect"
	"testing"

	"github.com/streamdp/ip-info/domain"
)

func Test_convertIpInfoDto(t *testing.T) {
	tests := []struct {
		name string
		dto  *domain.IpInfo
		want *Response
	}{
		{
			name: "regular conversion",
			dto: &domain.IpInfo{
				Ip:        net.ParseIP("8.8.8.8"),
				Continent: "NA",
				Country:   "US",
				StateProv: "California",
				City:      "Mountain View",
				Latitude:  -122.085,
				Longitude: 37.4223,
			},
			want: &Response{
				Ip:        "8.8.8.8",
				Continent: "NA",
				Country:   "US",
				StateProv: "California",
				City:      "Mountain View",
				Latitude:  -122.085,
				Longitude: 37.4223,
			},
		},
		{
			name: "empty conversion",
			dto:  &domain.IpInfo{},
			want: &Response{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertIpInfoDto(tt.dto); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertIpInfoDto() = %v, want %v", got, tt.want)
			}
		})
	}
}
