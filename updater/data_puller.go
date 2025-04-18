package updater

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/streamdp/ip-info/database"
)

const repeatIntervalOnError = 1 * time.Minute

type DataPuller interface {
	PullUpdates()
}

type DatabaseUpdater interface {
	UpdateIpDatabase(ctx context.Context) (duration time.Duration, err error)
}

type puller struct {
	ctx context.Context

	d DatabaseUpdater
	l *log.Logger
}

func New(ctx context.Context, d DatabaseUpdater, l *log.Logger) DataPuller {
	return &puller{
		ctx: ctx,

		d: d,
		l: l,
	}
}

func (p *puller) PullUpdates() {
	t := time.NewTimer(time.Second)

	for {
		select {
		case <-p.ctx.Done():
			t.Stop()
			return
		case <-t.C:
			t.Reset(repeatIntervalOnError)

			p.l.Println("ip database update started")
			nextUpdate, err := p.d.UpdateIpDatabase(p.ctx)
			if err != nil {
				p.l.Println(err)

				if !errors.Is(err, database.ErrNoUpdateRequired) {
					p.l.Printf("ip database update interrupted, retry after %0.1fs",
						repeatIntervalOnError.Seconds())
					continue
				}
			}

			t.Reset(nextUpdate)
			p.l.Printf("ip database update completed, next update through %0.1fh",
				nextUpdate.Hours())
		}
	}
}
