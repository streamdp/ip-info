package database

import (
	"errors"
	"fmt"
	"time"
)

const lockTimeout = 30 * time.Minute

var errLockTimeout = errors.New("lock timeout")

func (d *db) isLocked() (l bool) {
	_ = d.DB.QueryRowContext(d.ctx,
		"select true from pg_tables where schemaname='public' and tablename='_lock';").Scan(&l)
	return
}

func (d *db) lockStatus() (err error) {
	var createdAt time.Time
	if err = d.DB.QueryRowContext(d.ctx, "select created_at from _lock;").Scan(&createdAt); err != nil {
		return
	}

	if t := time.Now().Sub(createdAt); t < lockTimeout {
		return fmt.Errorf("lock already acquired %0.1fm ago (timeout=%0.1fm)",
			t.Minutes(), lockTimeout.Minutes())
	}

	return errLockTimeout
}

func (d *db) acquireLock() (err error) {
	if !d.isLocked() {
		d.l.Println("acquiring lock")
		_, err = d.DB.ExecContext(d.ctx, "select * into _lock from (values(now())) as a(created_at)")
		return
	}

	if err = d.lockStatus(); errors.Is(err, errLockTimeout) {
		d.l.Println("previous lock has expired")
		return d.resetLock()
	}
	return
}

func (d *db) resetLock() (err error) {
	d.l.Println("resetting lock")
	_, err = d.DB.ExecContext(d.ctx, "update _lock set created_at=now();")
	return
}

func (d *db) releaseLock() (err error) {
	d.l.Println("releasing lock")
	_, err = d.DB.ExecContext(d.ctx, "drop table if exists _lock;")
	return
}
