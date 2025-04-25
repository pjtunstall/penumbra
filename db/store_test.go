package db

import (
	"testing"
	"time"

	"dts/app"
)

func TestSQLiteStore_CreateAndGetTask(t *testing.T) {
	store := NewSQLiteStore(":memory:")
	defer store.Close()

	due := time.Now().Add(24 * time.Hour)
	task := app.Task{
		Title:       "Test Task",
		Description: "Testing",
		Status:      "open",
		Due:         due,
	}

	id, err := store.CreateTask(task)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	got, err := store.GetTaskById(id)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}

	if got.Title != task.Title || got.Description != task.Description || got.Status != task.Status || !got.Due.Equal(task.Due) {
		t.Errorf("Got %+v, want %+v", got, task)
	}
}

func TestSQLiteStore_Close(t *testing.T) {
	store := NewSQLiteStore(":memory:")
	if err := store.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
