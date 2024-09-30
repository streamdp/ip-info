package database

import (
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

func (d *db) loadDatabaseConfig() (err error) {
	dto := &configDto{}

	if err = d.DB.QueryRowContext(d.ctx, fmt.Sprintf("select * from %s;", configTableName)).Scan(
		&dto.LastUpdate,
		&dto.ActiveTable,
		&dto.BackupTable,
	); err != nil {
		return errors.New("couldn't load config from database")
	}

	d.cfg.LastUpdate = dto.LastUpdate
	d.cfg.ActiveTable = dto.ActiveTable
	d.cfg.BackupTable = dto.BackupTable

	return
}

func (d *db) updateDatabaseConfig() (err error) {
	d.l.Println("updating database config")
	if _, err = d.DB.ExecContext(
		d.ctx, fmt.Sprintf(
			"update %s set last_update=now() at time zone 'utc', active_table='%s', backup_table='%s';",
			configTableName,
			d.cfg.ActiveTable,
			d.cfg.BackupTable,
		),
	); err == nil {
		d.cfg.LastUpdate = time.Now().UTC()
	}

	return
}
