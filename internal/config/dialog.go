package config

import (
	"os"
	"strings"
)

const (
	dialogStorageTypeEnvName = "DIALOG_STORAGE_TYPE"
)

type DialogConfig interface {
	StorageType() string
}

type dialogConfig struct {
	storageType string
}

func NewDialogConfig() (DialogConfig, error) {
	storageType := os.Getenv(dialogStorageTypeEnvName)
	if len(storageType) == 0 {
		// По умолчанию используем PostgreSQL
		storageType = "postgres"
	}

	// Нормализуем значение
	storageType = strings.ToLower(storageType)
	if storageType != "postgres" && storageType != "redis" {
		storageType = "postgres"
	}

	return &dialogConfig{
		storageType: storageType,
	}, nil
}

func (cfg *dialogConfig) StorageType() string {
	return cfg.storageType
}
