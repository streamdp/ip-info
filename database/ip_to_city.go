package database

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/streamdp/ip-info/domain"
)

const downloadUrl = "https://download.db-ip.com/free/dbip-city-lite-%d-%s.csv.gz"

var ErrNoUpdateRequired = fmt.Errorf("no update required")

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

func (d *db) importCsv(url string) (err error) {
	d.l.Println("import ip database updates")
	_, err = d.DB.ExecContext(d.ctx,
		fmt.Sprintf("copy %s from program 'wget -qO- %s|gzip -d' delimiter ',' csv;", d.cfg.BackupTable, url),
	)
	return err
}

func (d *db) truncate() (err error) {
	d.l.Println(fmt.Sprintf("truncate %s table before importing update", d.cfg.BackupTable))
	_, err = d.DB.ExecContext(d.ctx, fmt.Sprintf("truncate table %s;", d.cfg.BackupTable))
	return err
}

func (d *db) createIndex() (err error) {
	d.l.Println(
		fmt.Sprintf("creating  %s_ip_start_gist_idx index on %s table", d.cfg.BackupTable, d.cfg.BackupTable),
	)
	_, err = d.DB.ExecContext(d.ctx,
		fmt.Sprintf("create index if not exists %s_ip_start_gist_idx on %s using gist (ip_start inet_ops);",
			d.cfg.BackupTable,
			d.cfg.BackupTable,
		),
	)
	return err
}

func (d *db) dropIndex() (err error) {
	d.l.Println(fmt.Sprintf("droping %s_ip_start_gist_idx index", d.cfg.BackupTable))
	_, err = d.DB.ExecContext(d.ctx,
		fmt.Sprintf("drop index if exists %s_ip_start_gist_idx;", d.cfg.BackupTable),
	)
	return err
}

func (d *db) switchTables() {
	d.l.Println("switching backup and working tables")
	d.cfg.ActiveTable, d.cfg.BackupTable = d.cfg.BackupTable, d.cfg.ActiveTable
}

func (d *db) IpInfo(ip net.IP) (*domain.IpInfo, error) {
	dto := &ipToCityDto{}

	if err := d.DB.QueryRowContext(
		d.ctx,
		fmt.Sprintf("select * from %s where '%s' >= ip_start order by ip_start desc limit 1;",
			d.cfg.ActiveTable,
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
		return nil, fmt.Errorf("couldn't determine ip location")
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

func (d *db) UpdateIpDatabase() (nextUpdate time.Duration, err error) {
	if err = d.loadDatabaseConfig(); err != nil {
		return 0, err
	}
	if d.cfg.LastUpdate.Month() == time.Now().Month() {
		return nextUpdateInterval(d.cfg.LastUpdate), ErrNoUpdateRequired
	}

	if err = d.truncate(); err != nil {
		return 0, fmt.Errorf("failed to truncate table: %w", err)
	}
	if err = d.dropIndex(); err != nil {
		return 0, fmt.Errorf("failed to drop index: %w", err)
	}
	if err = d.importCsv(buildDownloadUrl()); err != nil {
		return 0, fmt.Errorf("failed to import database: %w", err)
	}
	if err = d.createIndex(); err != nil {
		return 0, fmt.Errorf("failed to create index: %w", err)
	}
	d.switchTables()

	if err = d.updateDatabaseConfig(); err != nil {
		d.l.Println(fmt.Errorf("error updating config: %w", err))
	}

	return nextUpdateInterval(d.cfg.LastUpdate), nil
}

func buildDownloadUrl() string {
	year, month, _ := time.Now().Date()

	monthStr := strconv.Itoa(int(month))
	if month < 10 {
		monthStr = fmt.Sprintf("0%s", monthStr)
	}

	return fmt.Sprintf(downloadUrl, year, monthStr)
}

func nextUpdateInterval(t time.Time) time.Duration {
	year, month, _ := t.Date()
	return time.Date(year, month+1, 2, 0, 0, -1, 0, time.UTC).Sub(time.Now().UTC())
}
