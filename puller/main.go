package puller

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/streamdp/ip-info/database"
)

const (
	downloadUrl = "https://download.db-ip.com/free/dbip-city-lite-%d-%s.csv.gz"

	repeatIntervalOnError = 1 * time.Minute
)

type DataPuller interface {
	PullUpdates()
}

type puller struct {
	ctx context.Context

	d database.Database
	l *log.Logger
}

func New(ctx context.Context, d database.Database, l *log.Logger) DataPuller {
	return &puller{
		ctx: ctx,

		d: d,
		l: l,
	}
}

func buildDownloadUrl() string {
	year, month, _ := time.Now().Date()

	monthStr := strconv.Itoa(int(month))
	if month < 10 {
		monthStr = fmt.Sprintf("0%s", monthStr)
	}

	return fmt.Sprintf(downloadUrl, year, monthStr)
}

func (p *puller) nextUpdateInterval(t time.Time) (nextUpdate time.Duration) {
	year, month, _ := t.Date()
	nextUpdate = time.Date(year, month+1, 2, 0, 0, -1, 0, time.UTC).Sub(time.Now().UTC())
	p.l.Println(fmt.Sprintf("next database update through %0.1f hours", nextUpdate.Hours()))
	return
}

func (p *puller) PullUpdates() {
	var (
		err error
		t   = time.NewTimer(time.Second)
	)

	for {
		select {
		case <-p.ctx.Done():
			t.Stop()
			return
		case <-t.C:
			t.Reset(repeatIntervalOnError)

			if err = p.d.LoadDatabaseConfig(); err != nil {
				p.l.Println(err)
				continue
			}

			if p.d.LastUpdate().Month() == time.Now().Month() {
				t.Reset(p.nextUpdateInterval(p.d.LastUpdate()))
				continue
			}

			p.l.Println("ip database update started")

			p.l.Println(fmt.Sprintf("truncate ip database table before importing update"))
			if err = p.d.Truncate(); err != nil {
				p.l.Println(fmt.Errorf("failed to truncate table: %w", err))
				continue
			}

			p.l.Println(fmt.Sprintf("droping index on ip database table"))
			if err = p.d.DropIndex(); err != nil {
				p.l.Println(fmt.Errorf("failed to drop index: %w", err))
				continue
			}

			p.l.Println("import ip database updates")
			if err = p.d.Import(buildDownloadUrl()); err != nil {
				p.l.Println(fmt.Errorf("failed to import database: %w", err))
				continue
			}

			p.l.Println(fmt.Sprintf("creating new gist index on ip database table"))
			if err = p.d.CreateIndex(); err != nil {
				p.l.Println(fmt.Errorf("failed to create index: %w", err))
				continue
			}

			p.l.Println(fmt.Sprintf("switch backup and working tables"))
			p.d.SwitchTables()

			p.l.Println(fmt.Sprintf("update database config"))
			if err = p.d.UpdateDatabaseConfig(); err != nil {
				p.l.Println(fmt.Errorf("error updating config: %w", err))
			}

			p.l.Println("ip database updated successfully")
			t.Reset(p.nextUpdateInterval(p.d.LastUpdate()))
		}
	}
}
