package database

import (
	"context"
	"os"
	"testing"
)

func TestConnectDB(t *testing.T) {
	// Set env variables for test (or load from a .env.test file)
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PORT", "5432")

	db := ConnectDB()
	if db.Pool == nil {
		t.Fatal("Expected database pool, got nil")
	}

	// Try a simple ping
	err := db.Pool.Ping(context.Background())
	if err != nil {
		t.Fatalf("Database ping failed: %v", err)
	}
}
