package database

import (
	"context"
	"database/sql"
	"net"
	"time"

	_ "github.com/lib/pq"
	"github.com/streamdp/ip-info/domain"
)

type Database interface {
	IpInfo(ip net.IP) (*domain.IpInfo, error)
	LastUpdate() time.Time

	UpdateDatabaseConfig() error
	LoadDatabaseConfig() error

	Import(url string) error
	Truncate() error
	CreateIndex() error
	DropIndex() error
	SwitchTables()

	Close() error
}

type db struct {
	*sql.DB
	ctx context.Context

	cfg *domain.DatabaseConfig
}

func Connect(ctx context.Context, cfg *domain.Config) (d Database, err error) {
	sqlDb := &sql.DB{}
	if sqlDb, err = sql.Open("postgres", cfg.DatabaseUrl); err != nil {
		return
	}

	return &db{
		DB:  sqlDb,
		ctx: ctx,

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
