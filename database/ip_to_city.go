package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/streamdp/ip-info/domain"
)

const downloadUrl = "https://download.db-ip.com/free/dbip-city-lite-%d-%s.csv.gz"

var (
	ErrNoUpdateRequired = errors.New("no update required")
	ErrNoIpAddress      = errors.New("no ip address in the database")
	errDatabaseError    = errors.New("database error")
)

type ipToCityDto struct {
	ipStart   string
	ipEnd     string
	Continent string  `db:"continent"`
	Country   string  `db:"country"`
	StateProv string  `db:"state_prov"`
	City      string  `db:"city"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
}

func (d *db) importCsv(ctx context.Context, url string) error {
	d.l.Println("import ip database updates")

	_, err := d.ExecContext(ctx,
		fmt.Sprintf("copy %s from program 'wget -qO- %s|gzip -d' csv null 'null' delimiter ',';",
			d.dbIpCfg.BackupTable,
			url,
		))
	if err != nil {
		return fmt.Errorf("failed to import csv: %w", err)
	}

	return nil
}

func (d *db) truncate(ctx context.Context) error {
	d.l.Printf("truncate %s table before importing update", d.dbIpCfg.BackupTable)

	if _, err := d.ExecContext(ctx, fmt.Sprintf("truncate table %s;", d.dbIpCfg.BackupTable)); err != nil {
		return fmt.Errorf("failed to truncate table: %w", err)
	}

	return nil
}

func (d *db) createIndex(ctx context.Context) error {
	d.l.Printf("creating %s_ip_start_gist_idx index on %s table", d.dbIpCfg.BackupTable, d.dbIpCfg.BackupTable)

	_, err := d.ExecContext(ctx,
		fmt.Sprintf("create index if not exists %s_ip_start_gist_idx on %s using gist (ip_start inet_ops);",
			d.dbIpCfg.BackupTable,
			d.dbIpCfg.BackupTable,
		))
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

func (d *db) dropIndex(ctx context.Context) error {
	d.l.Printf("droping %s_ip_start_gist_idx index", d.dbIpCfg.BackupTable)

	_, err := d.ExecContext(ctx, fmt.Sprintf("drop index if exists %s_ip_start_gist_idx;", d.dbIpCfg.BackupTable))
	if err != nil {
		return fmt.Errorf("failed to drop index: %w", err)
	}

	return nil
}

func (d *db) switchTables() {
	d.l.Println("switching backup and working tables")
	d.dbIpCfg.ActiveTable, d.dbIpCfg.BackupTable = d.dbIpCfg.BackupTable, d.dbIpCfg.ActiveTable
}

func (d *db) IpInfo(ctx context.Context, ip net.IP) (*domain.IpInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cfg.RequestTimeout())
	defer cancel()

	dto := &ipToCityDto{}

	if err := d.QueryRowContext(ctx,
		fmt.Sprintf("select * from %s where '%s' >= ip_start order by ip_start desc limit 1;",
			d.dbIpCfg.ActiveTable,
			ip.String(),
		),
	).Scan(
		&dto.ipStart,
		&dto.ipEnd,
		&dto.Continent,
		&dto.Country,
		&dto.StateProv,
		&dto.City,
		&dto.Longitude,
		&dto.Latitude,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoIpAddress
		}

		return nil, errDatabaseError
	}

	return &domain.IpInfo{
		Ip:        ip,
		Continent: dto.Continent,
		Country:   dto.Country,
		StateProv: dto.StateProv,
		City:      dto.City,
		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,
	}, nil
}

func (d *db) UpdateIpDatabase(ctx context.Context) (time.Duration, error) {
	if err := d.loadDatabaseConfig(ctx); err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	if d.dbIpCfg.LastUpdate.Year() == now.Year() && d.dbIpCfg.LastUpdate.Month() == now.Month() {
		return nextUpdateInterval(d.dbIpCfg.LastUpdate), ErrNoUpdateRequired
	}

	if err := d.acquireLock(ctx); err != nil {
		return 0, err
	}
	defer func() {
		if errReleaseLock := d.releaseLock(ctx); errReleaseLock != nil {
			d.l.Printf("update ip database: %v", errReleaseLock)
		}
	}()

	if err := d.truncate(ctx); err != nil {
		return 0, err
	}
	if err := d.dropIndex(ctx); err != nil {
		return 0, err
	}
	if err := d.importCsv(ctx, buildDownloadUrl(time.Now())); err != nil {
		return 0, err
	}
	if err := d.createIndex(ctx); err != nil {
		return 0, err
	}
	d.switchTables()

	if err := d.updateDatabaseConfig(ctx); err != nil {
		d.l.Printf("update ip database: %v", err)
	}

	return nextUpdateInterval(d.dbIpCfg.LastUpdate), nil
}

func buildDownloadUrl(t time.Time) string {
	year, month, _ := t.Date()

	monthStr := strconv.Itoa(int(month))
	if month < 10 {
		monthStr = "0" + monthStr
	}

	return fmt.Sprintf(downloadUrl, year, monthStr)
}

func nextUpdateInterval(t time.Time) time.Duration {
	year, month, _ := t.Date()

	return time.Date(year, month+1, 2, 0, 0, -1, 0, time.UTC).Sub(time.Now().UTC())
}
