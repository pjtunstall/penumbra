package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"penumbra/app"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
    CreateUser(user app.User) error
    AddSessionToken(user_id int) (uuid.UUID, time.Time, error)
    GetUserIdFromSessionToken(sessionToken uuid.UUID) (int, error)
    GetTaskById(id uuid.UUID) (app.Task, error)
    SetTaskDone(id uuid.UUID) error
    GetUserByEmail(email string) (app.User, error)
    GetAllTasks(user_id int) ([]app.Task, error)
    CreateTask(task app.Task) error
    UpdateTask(task app.Task) error
    DeleteTask(id uuid.UUID) error
}

type SQLiteStore struct {
    db *sql.DB
    path string
    mu   sync.RWMutex
}

// Ensure SQLiteStore implements the Store interface.
var _ Store = &SQLiteStore{}

func (s *SQLiteStore) Close() error {
    s.mu.Lock()
    defer s.mu.Unlock()

	return s.db.Close()
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    if err := checkAllTablesExist(db); err != nil {
        return nil, err
    }

    return &SQLiteStore{db: db, path: path}, nil
}

func checkAllTablesExist(db *sql.DB) error {
    tables := []string{"tasks", "users"}
    for _, table := range tables {
        if err := checkTableExists(db, table); err != nil {
            return err
        }
    }
    return nil
}

func checkTableExists(db *sql.DB, tableName string) error {
    var name string
    err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?;", tableName).Scan(&name)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("table %s does not exist in the database", tableName)
        }
        return fmt.Errorf("error checking table %s: %w", tableName, err)
    }
    return nil
}

func (s *SQLiteStore) CreateUser(user app.User) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    _, err := s.db.Exec(`
    INSERT INTO users (name, password_hash, email, phone, session_token_hash, session_expires_at)
    VALUES (?, ?, ?, ?, '', ?)`,
    user.Name, user.PasswordHash, user.Email, user.Phone, time.Unix(0, 0))

    return err
}

func (s *SQLiteStore) GetUserByEmail(email string) (app.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var user app.User
    err := s.db.QueryRow(`SELECT id, name, password_hash, email, phone FROM users WHERE email = ?`, email).
        Scan(&user.Id, &user.Name, &user.PasswordHash, &user.Email, &user.Phone)

    return user, err
}

func (s *SQLiteStore) AddSessionToken(user_id int) (uuid.UUID, time.Time, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    sessionToken := uuid.New()
    expiresAt := time.Now().Add(24 * time.Hour)

    sessionTokenHash, err := bcrypt.GenerateFromPassword(sessionToken[:], 10)
    if err != nil {
        return uuid.Nil, time.Time{}, err
    }

	_, err = s.db.Exec(`
        UPDATE users SET session_token_hash = ?, session_expires_at = ? WHERE id = ?
    `, sessionTokenHash, expiresAt, user_id)
    if err != nil {
        return uuid.Nil, time.Time{}, err
    }

    return sessionToken, expiresAt, err
}

func (s *SQLiteStore) GetUserIdFromSessionToken(sessionToken uuid.UUID) (int, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if sessionToken == uuid.Nil {
        return 0, errors.New("session token is empty")
    }

    rows, err := s.db.Query(`SELECT id, session_token_hash, session_expires_at FROM users`)
    if err != nil {
        return 0, err
    }
    defer rows.Close()

    for rows.Next() {
        var userId int
        var hash []byte
        var expiresAt time.Time

        if err := rows.Scan(&userId, &hash, &expiresAt); err != nil {
            return 0, err
        }

        if len(hash) == 0 {
            continue // Skip users with no session token.
        }

        if err := bcrypt.CompareHashAndPassword(hash, sessionToken[:]); err == nil {
            if time.Now().After(expiresAt) {
                return 0, errors.New("session token has expired")
            }
            return userId, nil
        }
    }

    return 0, errors.New("session token not found")
}

func (s *SQLiteStore) GetTaskById(id uuid.UUID) (app.Task, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var t app.Task
    err := s.db.QueryRow(`SELECT id, title, description, done, due FROM tasks WHERE id = ?`, id).
        Scan(&t.Id, &t.Title, &t.Description, &t.Done, &t.Due)
    t.SetStatus()

    return t, err
}

func (s *SQLiteStore) GetAllTasks(user_id int) ([]app.Task, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    rows, err := s.db.Query(`SELECT id, title, description, done, due FROM tasks WHERE user_id = ?`, user_id)
    if err != nil {
        return nil, err
    }
    tasks := []app.Task{}
    for rows.Next() {
        var t app.Task
        err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Done, &t.Due)
        if err != nil {
            return nil, err
        }
        t.SetStatus()
        tasks = append(tasks, t)
    }

    return tasks, nil
}

func (s *SQLiteStore) DeleteTask(id uuid.UUID) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    _, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
    
    return err
}

func (s *SQLiteStore) CreateTask(t app.Task) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    _, err := s.db.Exec(`
        INSERT INTO tasks (id, user_id, title, description, done, due)
        VALUES (?, ?, ?, ?, ?, ?)
    `, t.Id, t.UserId, t.Title, t.Description, t.Done, t.Due)

    return err
}

func (s *SQLiteStore) UpdateTask(t app.Task) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    _, err := s.db.Exec(`
        UPDATE tasks
        SET title = ?, description = ?, done = ?, due = ?
        WHERE id = ?
    `, t.Title, t.Description, t.Done, t.Due, t.Id)

    return err
}

func (s *SQLiteStore) SetTaskDone(id uuid.UUID) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    _, err := s.db.Exec(`UPDATE tasks SET done = 1 WHERE id = ?`, id)

    return err
}

func TestSetTaskDone(t *testing.T) {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    store := &SQLiteStore{db: db}

    _, err = db.Exec(`
        CREATE TABLE tasks (
            id TEXT PRIMARY KEY,
            title TEXT,
            description TEXT,
            done INTEGER,
            due TEXT
        )`)
    if err != nil {
        t.Fatalf("failed to create table: %v", err)
    }

    taskID := uuid.New()
    _, err = db.Exec(`
        INSERT INTO tasks (id, title, description, done, due)
        VALUES (?, ?, ?, ?, ?)`, taskID, "Test Task", "Test Description", 0, "2025-05-13")
    if err != nil {
        t.Fatalf("failed to insert task: %v", err)
    }

    err = store.SetTaskDone(taskID)
    if err != nil {
        t.Fatalf("failed to set task done: %v", err)
    }

    var done int
    err = db.QueryRow(`SELECT done FROM tasks WHERE id = ?`, taskID).Scan(&done)
    if err != nil {
        t.Fatalf("failed to query task done status: %v", err)
    }

    if done != 1 {
        t.Fatalf("expected task done to be 1, got %d", done)
    }
}
