package database

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const lockTimeout = 30 * time.Minute

var (
	errLockTimeout  = errors.New("lock timeout")
	errLockAcquired = errors.New("lock already acquired")
)

func (d *db) isLocked(ctx context.Context) bool {
	var isLocked bool

	_ = d.QueryRowContext(ctx,
		"select true from pg_tables where schemaname='public' and tablename='_lock';",
	).Scan(&isLocked)

	return isLocked
}

func (d *db) lockStatus(ctx context.Context) error {
	var createdAt time.Time
	if err := d.QueryRowContext(ctx, "select created_at from _lock;").Scan(&createdAt); err != nil {
		return fmt.Errorf("failed to get lock status, %w", err)
	}

	if t := time.Since(createdAt); t < lockTimeout {
		return errLockAcquired
	}

	return errLockTimeout
}

func (d *db) acquireLock(ctx context.Context) error {
	if !d.isLocked(ctx) {
		d.l.Println("acquiring lock")

		if _, err := d.ExecContext(ctx, "select * into _lock from (values(now())) as a(created_at);"); err != nil {
			return fmt.Errorf("failed to acquire lock, %w", err)
		}

		return nil
	}

	if err := d.lockStatus(ctx); errors.Is(err, errLockTimeout) {
		d.l.Println("previous lock has expired")

		return d.resetLock(ctx)
	}

	return nil
}

func (d *db) resetLock(ctx context.Context) error {
	d.l.Println("resetting lock")

	if _, err := d.ExecContext(ctx, "update _lock set created_at=now();"); err != nil {
		return fmt.Errorf("failed to reset lock, %w", err)
	}

	return nil
}

func (d *db) releaseLock(ctx context.Context) error {
	d.l.Println("releasing lock")

	if _, err := d.ExecContext(ctx, "drop table if exists _lock;"); err != nil {
		return fmt.Errorf("failed to release lock, %w", err)
	}

	return nil
}
