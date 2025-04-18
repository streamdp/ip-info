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

	l   *log.Logger
	cfg *domain.DatabaseConfig
}

func Connect(cfg *config.App, l *log.Logger) (*db, error) {
	sqlDb, err := sql.Open("postgres", cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &db{
		DB: sqlDb,

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
