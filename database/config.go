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

	d.dbIpCfg.LastUpdate = dto.LastUpdate
	d.dbIpCfg.ActiveTable = dto.ActiveTable
	d.dbIpCfg.BackupTable = dto.BackupTable

	return nil
}

func (d *db) updateDatabaseConfig(ctx context.Context) error {
	d.l.Println("updating database config")

	_, err := d.ExecContext(ctx, fmt.Sprintf(
		"update %s set last_update=now() at time zone 'utc', active_table='%s', backup_table='%s';",
		configTableName,
		d.dbIpCfg.ActiveTable,
		d.dbIpCfg.BackupTable,
	))
	if err == nil {
		d.dbIpCfg.LastUpdate = time.Now().UTC()
	} else {
		return fmt.Errorf("error updating config: %w", err)
	}

	return nil
}
