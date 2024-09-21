package puller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/streamdp/ip-info/database"
)

const (
	repeatIntervalOnError = 1 * time.Minute
)

type DataPuller interface {
	PullUpdates()
}

type DatabaseUpdater interface {
	UpdateIpDatabase() (duration time.Duration, err error)
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
			nextUpdate, err := p.d.UpdateIpDatabase()
			if err != nil {
				p.l.Println(err)

				if !errors.Is(err, database.ErrNoUpdateRequired) {
					p.l.Println(fmt.Sprintf("ip database update interrupted, retry after %0.1fs",
						repeatIntervalOnError.Seconds()))
					continue
				}
			}

			t.Reset(nextUpdate)
			p.l.Println(fmt.Sprintf("ip database update completed, next update through %0.1fh",
				nextUpdate.Hours()))
		}
	}
}
