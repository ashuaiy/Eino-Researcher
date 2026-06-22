package store

import (
	"net/url"
	"testing"

	"eino-researcher/internal/config"
)

func TestPostgresDSNEncodesCredentials(t *testing.T) {
	dsn := PostgresDSN(config.PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user/name",
		Password: "p@ss/word?#",
		Database: "eino_researcher",
		SSLMode:  "disable",
	})
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse DSN: %v", err)
	}
	if parsed.User.Username() != "user/name" {
		t.Fatalf("unexpected user: %q", parsed.User.Username())
	}
	password, ok := parsed.User.Password()
	if !ok || password != "p@ss/word?#" {
		t.Fatalf("unexpected password: %q", password)
	}
}
