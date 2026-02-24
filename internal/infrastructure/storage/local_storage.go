package storage

import (
	"fmt"
	"os"

	"github.com/galihaleanda/event-invitation/internal/config"
)

func EnsureStorageDirs(cfg *config.Config) error {
	dirs := []string{
		cfg.Storage.BasePath,
		fmt.Sprintf("%s/events", cfg.Storage.BasePath),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
