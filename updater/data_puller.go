package updater

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/streamdp/ip-info/database"
)

const repeatIntervalOnError = 1 * time.Minute

type DatabaseUpdater interface {
	UpdateIpDatabase(ctx context.Context) (duration time.Duration, err error)
}

type puller struct {
	d DatabaseUpdater
	l *log.Logger
}

func New(d DatabaseUpdater, l *log.Logger) *puller {
	return &puller{
		d: d,
		l: l,
	}
}

func (p *puller) PullUpdates(ctx context.Context) {
	t := time.NewTimer(time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			p.l.Println("ip database update started")
			nextUpdate, err := p.d.UpdateIpDatabase(ctx)
			if err != nil {
				p.l.Printf("failed to update ip database: %v", err)

				if !errors.Is(err, database.ErrNoUpdateRequired) {
					p.l.Printf("ip database update interrupted, retry after %0.1fs",
						repeatIntervalOnError.Seconds())
					t.Reset(repeatIntervalOnError)

					continue
				}
			}

			p.l.Printf("ip database update completed, next update through %0.1fh", nextUpdate.Hours())
			t.Reset(nextUpdate)
		}
	}
}
