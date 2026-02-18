package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/domain"
)

type db struct {
	*sql.DB

	cfg     *config.Database
	l       *log.Logger
	dbIpCfg *domain.DatabaseConfig
}

func Connect(l *log.Logger, cfg *config.Database) (*db, error) {
	sqlDb, err := sql.Open("postgres", cfg.Url())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &db{
		DB:  sqlDb,
		cfg: cfg,

		l: l,
		dbIpCfg: &domain.DatabaseConfig{
			LastUpdate:  time.Now().Add(-31 * 24 * time.Hour),
			ActiveTable: "ip_to_city_one",
			BackupTable: "ip_to_city_two",
		},
	}, nil
}

func (d *db) Close() error {
	if d.DB == nil {
		return nil
	}

	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
