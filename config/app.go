package config

import (
	"errors"
	"os"
)

var (
	errEmptyDatabaseUrl    = errors.New("database url cannot be blank")
	errEmptyDatabaseUrlEnv = errors.New("IP_INFO_DATABASE_URL environment variable not set")
)

type App struct {
	Http    *Http
	Grpc    *Grpc
	Limiter *Limiter
	Cache   *Cache
	Redis   *Redis

	DatabaseUrl string
	Version     string
}

func newAppConfig() *App {
	return &App{
		Http:    newHttpConfig(),
		Grpc:    newGrpcConfig(),
		Limiter: newLimiterConfig(),
		Cache:   newCacheConfig(),
		Redis:   newRedisConfig(),

		DatabaseUrl: "",
		Version:     "",
	}
}

func (a *App) loadEnvs() error {
	if a.DatabaseUrl = os.Getenv("IP_INFO_DATABASE_URL"); a.DatabaseUrl == "" {
		return errEmptyDatabaseUrlEnv
	}

	a.Limiter.loadEnvs()
	a.Cache.loadEnvs()

	return nil
}

func (a *App) validate() error {
	if a.DatabaseUrl == "" {
		return errEmptyDatabaseUrl
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
