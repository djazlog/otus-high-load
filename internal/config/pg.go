package config

import (
	"errors"
	"os"
)

const (
	dsnEnvName        = "PG_DSN"
	dsnReplicaEnvName = "PG_REPLICA_DSN"
)

type PGConfig interface {
	DSN() string
	DSNReplica() string
}

type pgConfig struct {
	dsn        string
	dsnReplica string
}

func NewPGConfig() (PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	dsnReplica := os.Getenv(dsnReplicaEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn:        dsn,
		dsnReplica: dsnReplica,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
func (cfg *pgConfig) DSNReplica() string {
	return cfg.dsnReplica
}
