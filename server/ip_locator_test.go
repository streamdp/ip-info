package ip_locator

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
)

func TestLocateIp(t *testing.T) {
	type args struct {
		d        database.Database
		ipString string
	}
	tests := []struct {
		name       string
		args       args
		wantIpInfo *domain.IpInfo
		wantErr    bool
	}{
		{
			name: "get ip location info",
			args: args{
				&databaseMock{
					err: nil,
					ipInfo: &domain.IpInfo{
						Ip:        net.ParseIP("82.28.25.43"),
						Continent: "NA",
						Country:   "US",
						StateProv: "California",
						City:      "Mountain View",
						Latitude:  -122.085,
						Longitude: 37.4223,
					},
				},
				"82.28.25.43",
			},
			wantIpInfo: &domain.IpInfo{
				Ip:        net.ParseIP("82.28.25.43"),
				Continent: "NA",
				Country:   "US",
				StateProv: "California",
				City:      "Mountain View",
				Latitude:  -122.085,
				Longitude: 37.4223,
			},
			wantErr: false,
		},
		{
			name: "get database error",
			args: args{
				&databaseMock{
					err:    database.ErrNoIpAddress,
					ipInfo: nil,
				},
				"82.28.25.43",
			},
			wantIpInfo: nil,
			wantErr:    true,
		},
		{
			name: "get ip parsing error",
			args: args{
				&databaseMock{},
				"256.28.25.43",
			},
			wantIpInfo: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIpInfo, err := LocateIp(tt.args.d, tt.args.ipString)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocateIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIpInfo, tt.wantIpInfo) {
				t.Errorf("LocateIp() gotIpInfo = %v, want %v", gotIpInfo, tt.wantIpInfo)
			}
		})
	}
}

type databaseMock struct {
	err    error
	ipInfo *domain.IpInfo
}

func (d *databaseMock) IpInfo(_ net.IP) (*domain.IpInfo, error) {
	return d.ipInfo, d.err
}

func (d *databaseMock) UpdateIpDatabase() (nextUpdate time.Duration, err error) {
	return 0, nil
}

func (d *databaseMock) Close() error {
	return nil
}
