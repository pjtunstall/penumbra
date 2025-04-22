package main

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    exePath, err := os.Executable()
    if err != nil {
        log.Fatalf("finding executable: %v", err)
    }

    dbFile := getDBPath(exePath)

    if err := resetDB(dbFile); err != nil {
        log.Fatalf("resetting database: %v", err)
    }

    migrations := "file://" + filepath.Join(filepath.Dir(exePath), "migrations")

    m, err := migrate.New(
        migrations,
        "sqlite3://"+dbFile,
    )
    if err != nil {
        log.Fatal("Error creating migration instance:", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal("Error applying migrations:", err)
    }

    log.Println("Migrations applied successfully!")
}

func getDBPath(exePath string) string {
    root := filepath.Dir(exePath)

    dataDir := filepath.Join(root, "data")
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        log.Fatalf("creating data directory: %v", err)
    }

    return filepath.Join(dataDir, "dev.db")
}

func resetDB(path string) error {
    if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
        return err
    }
    return nil
}