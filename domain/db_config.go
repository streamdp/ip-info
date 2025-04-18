package domain

import "time"

type DatabaseConfig struct {
	LastUpdate  time.Time `db:"last_update"  json:"last_update"`
	ActiveTable string    `db:"active_table" json:"active_table"`
	BackupTable string    `db:"backup_table" json:"backup_table"`
}
