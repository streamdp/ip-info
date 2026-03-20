package database

import (
	"context"
	"errors"
	"fmt"
)

var errLockAcquired = errors.New("lock already acquired")

const advisoryLockID = 4788125797291836 // Unique ID for the ip-info update process lock

func (d *db) acquireLock(ctx context.Context) error {
	d.l.Println("acquiring advisory lock")

	var acquired bool
	err := d.QueryRowContext(ctx, "select pg_try_advisory_lock($1);", advisoryLockID).Scan(&acquired)
	if err != nil {
		return fmt.Errorf("failed to acquire advisory lock: %w", err)
	}

	if !acquired {
		return errLockAcquired
	}

	return nil
}

func (d *db) releaseLock(ctx context.Context) error {
	d.l.Println("releasing advisory lock")

	_, err := d.ExecContext(ctx, "select pg_advisory_unlock($1);", advisoryLockID)
	if err != nil {
		return fmt.Errorf("failed to release advisory lock: %w", err)
	}

	return nil
}
