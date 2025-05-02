package config

import (
	"fmt"
)

type App struct {
	Http     *Http
	Grpc     *Grpc
	Limiter  *Limiter
	Cache    *Cache
	Redis    *Redis
	Database *Database

	version string
}

func newAppConfig() *App {
	return &App{
		Http:     newHttpConfig(),
		Grpc:     newGrpcConfig(),
		Limiter:  newLimiterConfig(),
		Cache:    newCacheConfig(),
		Redis:    newRedisConfig(),
		Database: newDatabaseConfig(),

		version: version,
	}
}

func (a *App) Version() string {
	return a.version
}

func (a *App) loadEnvs() error {
	if err := a.Database.loadEnvs(); err != nil {
		return fmt.Errorf("failed to load 'IP_INFO_DATABASE_URL' env: %w", err)
	}

	a.Limiter.loadEnvs()
	a.Cache.loadEnvs()

	return nil
}

func (a *App) validate() error {
	if err := a.Database.validate(); err != nil {
		return err
	}
	if err := a.Http.validate(); err != nil {
		return err
	}
	if err := a.Grpc.validate(); err != nil {
		return err
	}
	if err := a.Limiter.validate(); err != nil {
		return err
	}
	if err := a.Cache.validate(); err != nil {
		return err
	}
	if err := a.Redis.validate(); err != nil {
		return err
	}

	return nil
}
