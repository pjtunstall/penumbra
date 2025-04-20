package db

import (
	"database/sql"
	"log"

	"dts/backend/internal/app"

	_ "modernc.org/sqlite"
)

type Store interface {
    CreateTask(task app.Task) (int, error)
    GetTask(id int) (app.Task, error)
}

type SQLiteStore struct {
    db *sql.DB
}

// Ensure SQLiteStore implements the Store interface
var _ Store = &SQLiteStore{}

func NewSQLiteStore(path string) *SQLiteStore {
    db, err := sql.Open("sqlite", path+"?_busy_timeout=5000&_journal_mode=WAL")
    if err != nil {
        log.Fatal(err)
    }

    createTable := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT,
        status TEXT NOT NULL,
        due DATETIME NOT NULL
    )`
    if _, err := db.Exec(createTable); err != nil {
        log.Fatal(err)
    }

    return &SQLiteStore{db: db}
}

func (s *SQLiteStore) CreateTask(t app.Task) (int, error) {
    res, err := s.db.Exec(`INSERT INTO tasks (title, description, status, due) VALUES (?, ?, ?, ?)`,
        t.Title, t.Description, t.Status, t.Due)
    if err != nil {
        return 0, err
    }
    id, _ := res.LastInsertId()
    return int(id), nil
}

func (s *SQLiteStore) GetTask(id int) (app.Task, error) {
    var t app.Task
    err := s.db.QueryRow(`SELECT id, title, description, status, due FROM tasks WHERE id = ?`, id).
        Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Due)
    return t, err
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}


// Add UpdateTask, DeleteTask, ListTasks later
