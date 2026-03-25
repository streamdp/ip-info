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

func (d *db) loadConfig(ctx context.Context) (*configDto, error) {
	dto := &configDto{}

	if err := d.QueryRowContext(ctx, "select last_update, active_table, backup_table from config;").Scan(
		&dto.LastUpdate,
		&dto.ActiveTable,
		&dto.BackupTable,
	); err != nil {
		return nil, errLoadConfig
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.dbIpCfg.LastUpdate = dto.LastUpdate
	d.dbIpCfg.ActiveTable = dto.ActiveTable
	d.dbIpCfg.BackupTable = dto.BackupTable

	return dto, nil
}

func (d *db) updateConfig(ctx context.Context, activeTable, backupTable string) error {
	d.l.Println("updating database config")

	_, err := d.ExecContext(ctx,
		`update config set last_update=now() at time zone 'utc', active_table=$1, backup_table=$2;`,
		activeTable,
		backupTable,
	)
	if err != nil {
		return fmt.Errorf("error updating config: %w", err)
	}

	d.mu.Lock()
	d.dbIpCfg.LastUpdate = time.Now().UTC()
	d.mu.Unlock()

	return nil
}

func (d *db) activeTable() string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.dbIpCfg.ActiveTable
}

func (d *db) backupTable() string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.dbIpCfg.BackupTable
}

func (d *db) lastUpdate() time.Time {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.dbIpCfg.LastUpdate
}

func (d *db) swapTables() (string, string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.l.Println("swapping working and backup tables")
	d.dbIpCfg.ActiveTable, d.dbIpCfg.BackupTable = d.dbIpCfg.BackupTable, d.dbIpCfg.ActiveTable

	return d.dbIpCfg.ActiveTable, d.dbIpCfg.BackupTable
}
