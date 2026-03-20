package database

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type configDto struct {
	LastUpdate  time.Time `db:"last_update"`
	ActiveTable string    `db:"active_table"`
	BackupTable string    `db:"backup_table"`
}

var errLoadConfig = errors.New("couldn't load config from database")

func (d *db) loadDatabaseConfig(ctx context.Context) error {
	dto := &configDto{}

	if err := d.QueryRowContext(ctx, "select last_update, active_table, backup_table from config;").Scan(
		&dto.LastUpdate,
		&dto.ActiveTable,
		&dto.BackupTable,
	); err != nil {
		return errLoadConfig
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.dbIpCfg.LastUpdate = dto.LastUpdate
	d.dbIpCfg.ActiveTable = dto.ActiveTable
	d.dbIpCfg.BackupTable = dto.BackupTable

	return nil
}

func (d *db) updateDatabaseConfig(ctx context.Context) error {
	d.l.Println("updating database config")

	d.mu.RLock()
	active := d.dbIpCfg.ActiveTable
	backup := d.dbIpCfg.BackupTable
	d.mu.RUnlock()

	_, err := d.ExecContext(ctx,
		`update config set last_update=now() at time zone 'utc', active_table=$1, backup_table=$2;`,
		active,
		backup,
	)
	if err != nil {
		return fmt.Errorf("error updating config: %w", err)
	}

	d.mu.Lock()
	d.dbIpCfg.LastUpdate = time.Now().UTC()
	d.mu.Unlock()

	return nil
}
