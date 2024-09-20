package database

import (
	"fmt"
	"net"

	"github.com/streamdp/ip-info/domain"
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

func (d *db) Import(url string) (err error) {
	_, err = d.DB.ExecContext(d.ctx,
		fmt.Sprintf("copy %s from program 'wget -qO- %s|gzip -d' delimiter ',' csv;", d.cfg.BackupTable, url),
	)
	return err
}

func (d *db) Truncate() (err error) {
	_, err = d.DB.ExecContext(d.ctx, fmt.Sprintf("truncate table %s;", d.cfg.BackupTable))
	return err
}

func (d *db) CreateIndex() (err error) {
	_, err = d.DB.ExecContext(d.ctx,
		fmt.Sprintf("create index if not exists %s_ip_start_gist_idx on %s using gist (ip_start inet_ops);",
			d.cfg.BackupTable,
			d.cfg.BackupTable,
		),
	)
	return err
}

func (d *db) DropIndex() (err error) {
	_, err = d.DB.ExecContext(d.ctx,
		fmt.Sprintf("drop index if exists %s_ip_start_gist_idx;", d.cfg.BackupTable),
	)
	return err
}

func (d *db) SwitchTables() {
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
