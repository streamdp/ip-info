package database

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const configTableName = "config"

type configDto struct {
	LastUpdate  time.Time `db:"last_update"`
	ActiveTable string    `db:"active_table"`
	BackupTable string    `db:"backup_table"`
}

var errLoadConfig = errors.New("couldn't load config from database")

func (d *db) loadDatabaseConfig(ctx context.Context) error {
	dto := &configDto{}

	if err := d.QueryRowContext(ctx, fmt.Sprintf("select * from %s;", configTableName)).Scan(
		&dto.LastUpdate,
		&dto.ActiveTable,
		&dto.BackupTable,
	); err != nil {
		return errLoadConfig
	}

	d.cfg.LastUpdate = dto.LastUpdate
	d.cfg.ActiveTable = dto.ActiveTable
	d.cfg.BackupTable = dto.BackupTable

	return nil
}

func (d *db) updateDatabaseConfig(ctx context.Context) error {
	d.l.Println("updating database config")

	_, err := d.ExecContext(ctx, fmt.Sprintf(
		"update %s set last_update=now() at time zone 'utc', active_table='%s', backup_table='%s';",
		configTableName,
		d.cfg.ActiveTable,
		d.cfg.BackupTable,
	))
	if err == nil {
		d.cfg.LastUpdate = time.Now().UTC()
	} else {
		return fmt.Errorf("error updating config: %w", err)
	}

	return nil
}
