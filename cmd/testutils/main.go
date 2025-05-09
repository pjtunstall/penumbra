package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"penumbra/app"
	"penumbra/db"

	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("finding executable: %v", err)
	}

	migrations := "file://" + filepath.Join(filepath.Dir(exePath), "migrations")

	m, err := migrate.New(migrations, "sqlite3://:memory:")
	if err != nil {
		log.Fatalf("Error setting up migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error applying migrations: %v", err)
	}

	if err := runTests(); err != nil {
		log.Fatalf("Tests failed: %v", err)
	}

	fmt.Println("Migrations applied and tests passed successfully!")
}

func runTests() error {
	t := testing.T{}
	TestSQLiteStore_CreateAndGetTask(&t)
	TestSQLiteStore_Close(&t)

	return nil
}

func TestSQLiteStore_CreateAndGetTask(t *testing.T) {
	store, err := db.NewSQLiteStore(":memory:")
	defer store.Close()

	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}

	due := time.Now().Add(24 * time.Hour)
	task := app.Task{
		Id:          uuid.New(),
		Title:       "Test Task",
		UserId:      1,
		Description: "Testing",
		Due:         due,
		Done:        0,
	}

	err = store.UpsertTask(task)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	got, err := store.GetTaskById(task.Id)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}

	if got.Title != task.Title || got.Description != task.Description || !got.Due.Equal(task.Due) || got.Done != task.Done {
		t.Errorf("Got %+v, want %+v", got, task)
	}
}

func TestSQLiteStore_Close(t *testing.T) {
	store, err := db.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("NewSQLiteStore failed: %v", err)
	}

	if err := store.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
