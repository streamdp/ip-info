package ip_locator

import (
	"context"
	"errors"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/server"
)

func TestLocateIp(t *testing.T) {
	tests := []struct {
		name       string
		locator    server.Locator
		ipString   string
		wantIpInfo *domain.IpInfo
		wantErr    bool
	}{
		{
			name: "get ip location info from db",
			locator: New(
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
				nil,
			),
			ipString: "82.28.25.43",
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
			name: "get ip location info from cache",
			locator: New(
				&databaseMock{},
				&cacheMock{
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
			),
			ipString: "82.28.25.43",
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
			locator: New(
				&databaseMock{
					err:    database.ErrNoIpAddress,
					ipInfo: nil,
				},
				&cacheMock{
					getErr: errors.New("redis: nil"),
				},
			),
			ipString:   "82.28.25.43",
			wantIpInfo: nil,
			wantErr:    true,
		},
		{
			name: "get database error, when cache uninitialized",
			locator: New(
				&databaseMock{
					err:    database.ErrNoIpAddress,
					ipInfo: nil,
				},
				nil,
			),
			ipString:   "82.28.25.43",
			wantIpInfo: nil,
			wantErr:    true,
		},
		{
			name: "get set cache error",
			locator: New(
				&databaseMock{},
				&cacheMock{
					getErr: errors.New("redis: nil"),
					setErr: errors.New("cache: redis: nil"),
				},
			),
			ipString:   "82.28.25.43",
			wantIpInfo: nil,
			wantErr:    true,
		},
		{
			name:       "get ip parsing error",
			locator:    New(&databaseMock{}, &cacheMock{}),
			ipString:   "256.28.25.43",
			wantIpInfo: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIpInfo, err := tt.locator.GetIpInfo(context.Background(), tt.ipString)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIpInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIpInfo, tt.wantIpInfo) {
				t.Errorf("GetIpInfo() gotIpInfo = %v, want %v", gotIpInfo, tt.wantIpInfo)
			}
		})
	}
}

type databaseMock struct {
	err    error
	ipInfo *domain.IpInfo
}

func (d *databaseMock) IpInfo(_ context.Context, _ net.IP) (*domain.IpInfo, error) {
	return d.ipInfo, d.err
}

func (d *databaseMock) UpdateIpDatabase(_ context.Context) (nextUpdate time.Duration, err error) {
	return 0, nil
}

func (d *databaseMock) Close() error {
	return nil
}

type cacheMock struct {
	getErr, setErr error

	ipInfo *domain.IpInfo
}

func (c *cacheMock) Set(_ *domain.IpInfo) error {
	return c.setErr
}

func (c *cacheMock) Get(string) (*domain.IpInfo, error) {
	return c.ipInfo, c.getErr
}
