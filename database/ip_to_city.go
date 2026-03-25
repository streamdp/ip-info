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
	ipRange   string  `db:"ip_range"`
}

func (d *db) importCsv(ctx context.Context, url string) error {
	d.l.Println("import ip database updates")

	_, err := d.ExecContext(ctx,
		fmt.Sprintf("copy %s from program 'wget -qO- %s|gzip -d' csv null 'null' delimiter ',';",
			d.backupTable(),
			url,
		))
	if err != nil {
		return fmt.Errorf("failed to import csv: %w", err)
	}

	return nil
}

func (d *db) truncate(ctx context.Context) error {
	backupTable := d.backupTable()

	d.l.Printf("truncate %s table before importing update", backupTable)

	_, err := d.ExecContext(ctx, fmt.Sprintf("truncate table %s;", backupTable))
	if err != nil {
		return fmt.Errorf("failed to truncate table: %w", err)
	}

	return nil
}

func (d *db) createIndex(ctx context.Context) error {
	backupTable := d.backupTable()

	indexName := fmt.Sprintf("%s_ip_range_spgist_idx", backupTable)

	d.l.Printf("creating %s index on %s table", indexName, backupTable)

	_, err := d.ExecContext(ctx,
		fmt.Sprintf(`create index if not exists %s on %s using spgist(ip_range);`,
			indexName,
			backupTable,
		))
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

func (d *db) dropIndex(ctx context.Context) error {
	indexName := fmt.Sprintf("%s_ip_range_spgist_idx", d.backupTable())

	d.l.Println("droping ", indexName)

	_, err := d.ExecContext(ctx, fmt.Sprintf("drop index if exists %s;", indexName))
	if err != nil {
		return fmt.Errorf("failed to drop index: %w", err)
	}

	return nil
}

func (d *db) IpInfo(ctx context.Context, ip net.IP) (*domain.IpInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cfg.RequestTimeout())
	defer cancel()

	dto := &ipToCityDto{}
	if err := d.QueryRowContext(ctx, fmt.Sprintf(
		`select * from %s where ip_range::inet>>='%s';`,
		d.activeTable(),
		ip.String(),
	)).Scan(
		&dto.ipStart,
		&dto.ipEnd,
		&dto.Continent,
		&dto.Country,
		&dto.StateProv,
		&dto.City,
		&dto.Longitude,
		&dto.Latitude,
		&dto.ipRange,
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
	cfg, err := d.loadConfig(ctx)
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	if cfg.LastUpdate.Year() == now.Year() && cfg.LastUpdate.Month() == now.Month() {
		return nextUpdateInterval(d.dbIpCfg.LastUpdate), ErrNoUpdateRequired
	}

	if err = d.acquireLock(ctx); err != nil {
		return 0, err
	}
	defer func() {
		if errReleaseLock := d.releaseLock(ctx); errReleaseLock != nil {
			d.l.Printf("update ip database: %v", errReleaseLock)
		}
	}()

	if err = d.truncate(ctx); err != nil {
		return 0, err
	}
	if err = d.dropIndex(ctx); err != nil {
		return 0, err
	}
	if err = d.importCsv(ctx, buildDownloadUrl(time.Now().UTC())); err != nil {
		return 0, err
	}
	if err = d.createIndex(ctx); err != nil {
		return 0, err
	}

	activeTable, backupTable := d.swapTables()
	if err = d.updateConfig(ctx, activeTable, backupTable); err != nil {
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
