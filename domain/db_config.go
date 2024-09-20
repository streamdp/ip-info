package domain

import "time"

type DatabaseConfig struct {
	LastUpdate  time.Time `json:"last_update" db:"last_update"`
	ActiveTable string    `json:"active_table" db:"active_table"`
	BackupTable string    `json:"backup_table" db:"backup_table"`
}
