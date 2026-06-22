package store

import (
	"fmt"

	"eino-researcher/internal/config"
)

func PostgresDSN(cfg config.PostgresConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)
}
