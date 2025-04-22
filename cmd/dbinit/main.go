package main

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"dts/db"
)

func main() {
    exePath, err := os.Executable()
    if err != nil {
        log.Fatalf("finding executable: %v", err)
    }

    dbFile := getDBPath(exePath)

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

func ApplyMigrations(store *db.SQLiteStore) error {
    migrationsPath := "file://" + filepath.Join("migrations")

    m, err := migrate.New(migrationsPath, "sqlite3://"+store.Path)
    if err != nil {
        return err
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    log.Println("Migrations applied successfully.")
    return nil
}