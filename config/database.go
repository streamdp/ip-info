package config

import (
	"errors"
	"os"
	"time"
)

const databaseRequestTimeout = 5000

var errEmptyDatabaseUrl = errors.New("database url cannot be blank")

type Database struct {
	url            string
	requestTimeout int
}

func newDatabaseConfig() *Database {
	return &Database{
		url:            "",
		requestTimeout: databaseRequestTimeout,
	}
}

func (d *Database) RequestTimeout() time.Duration {
	return time.Duration(d.requestTimeout) * time.Millisecond
}

func (d *Database) SetRequestTimeout(timeout int) *Database {
	d.requestTimeout = timeout

	return d
}

func (d *Database) Url() string {
	return d.url
}

func (d *Database) loadEnvs() error {
	if d.url = os.Getenv("IP_INFO_DATABASE_URL"); d.url == "" {
		return errEmptyDatabaseUrl
	}

	return nil
}

func (d *Database) validate() error {
	if d.url == "" {
		return errEmptyDatabaseUrl
	}

	return nil
}
