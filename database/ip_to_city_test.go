package database

import (
	"testing"
	"time"
)

func Test_nextUpdateInterval(t *testing.T) {
	tests := []struct {
		name         string
		t            time.Time
		futureResult bool
	}{
		{
			name: "next update interval in the future",
			t: time.Date(
				time.Now().Year(), time.Now().Month()+2, 1, 0, 0, 0, 0, time.UTC,
			),
			futureResult: true,
		},
		{
			name: "next update interval in the past",
			t: time.Date(
				time.Now().Year(), time.Now().Month()-2, 1, 0, 0, 0, 0, time.UTC,
			),
			futureResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (nextUpdateInterval(tt.t) > 0) != tt.futureResult {
				t.Errorf("nextUpdateInterval() should return %s duration", func(b bool) string {
					if b {
						return "positive"
					}

					return "negative"
				}(tt.futureResult))
			}
		})
	}
}

func Test_buildDownloadUrl(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{
			name: "get download url for September",
			t:    time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			want: "https://download.db-ip.com/free/dbip-city-lite-2024-09.csv.gz",
		},
		{
			name: "get download url for December",
			t:    time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			want: "https://download.db-ip.com/free/dbip-city-lite-2024-12.csv.gz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildDownloadUrl(tt.t); got != tt.want {
				t.Errorf("buildDownloadUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
