package database

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	_ "github.com/lib/pq"
	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/domain"
)

type Database interface {
	IpInfo(ip net.IP) (*domain.IpInfo, error)
	UpdateIpDatabase() (nextUpdate time.Duration, err error)

	Close() error
}

type db struct {
	*sql.DB
	ctx context.Context

	l   *log.Logger
	cfg *domain.DatabaseConfig
}

func Connect(ctx context.Context, cfg *config.App, l *log.Logger) (d Database, err error) {
	sqlDb := &sql.DB{}
	if sqlDb, err = sql.Open("postgres", cfg.DatabaseUrl); err != nil {
		return
	}

	return &db{
		DB:  sqlDb,
		ctx: ctx,

		l: l,
		cfg: &domain.DatabaseConfig{
			LastUpdate:  time.Now().Add(-31 * 24 * time.Hour),
			ActiveTable: "ip_to_city_one",
			BackupTable: "ip_to_city_two",
		},
	}, nil
}

func (d *db) Close() (err error) {
	if d.DB == nil {
		return
	}
	return d.DB.Close()
}
