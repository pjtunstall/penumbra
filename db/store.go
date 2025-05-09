package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"penumbra/app"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
    CreateUser(user app.User) error
    AddSessionToken(user_id int) (string, time.Time, error)
    GetUserIdFromSessionToken(sessionToken string) (int, error)
    GetTaskById(id string) (app.Task, error)
    SetTaskDone(id string) error
    GetUserByEmail(email string) (app.User, error)
    GetAllTasks(user_id int) ([]app.Task, error)
    UpsertTask(task app.Task) error
    DeleteTask(id string) error
}

type SQLiteStore struct {
    db *sql.DB
    Path string
}

// Ensure SQLiteStore implements the Store interface
var _ Store = &SQLiteStore{}

func NewSQLiteStore(path string) *SQLiteStore {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        log.Fatal(err)
    }

    panicIfTablesDoNotExist(db)

    return &SQLiteStore{db: db, Path: path}
}

func panicIfTablesDoNotExist(db *sql.DB) {
    var tableName string
    err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tasks';").Scan(&tableName)
    if err != nil {
        if err == sql.ErrNoRows {
            panic("tasks table does not exist in the database")
        }
        log.Fatal(err)
    }

    err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users';").Scan(&tableName)
    if err != nil {
        if err == sql.ErrNoRows {
            panic("users table does not exist in the database")
        }
        log.Fatal(err)
    }
}

func (s *SQLiteStore) CreateUser(user app.User) error {
    _, err := s.db.Exec(`
    INSERT INTO users (name, password_hash, email, phone, session_token_hash, session_expires_at)
    VALUES (?, ?, ?, ?, '', ?)`,
    user.Name, user.PasswordHash, user.Email, user.Phone, time.Unix(0, 0))
    return err
}

func (s *SQLiteStore) GetUserByEmail(email string) (app.User, error) {
    var user app.User
    err := s.db.QueryRow(`SELECT id, name, password_hash, email, phone FROM users WHERE email = ?`, email).
        Scan(&user.Id, &user.Name, &user.PasswordHash, &user.Email, &user.Phone)
    return user, err
}

func (s *SQLiteStore) AddSessionToken(user_id int) (string, time.Time, error) {
    sessionToken := uuid.NewString()
    expiresAt := time.Now().Add(24 * time.Hour)

    sessionTokenHash, err := bcrypt.GenerateFromPassword([]byte(sessionToken), 10)
    if err != nil {
        return "", time.Time{}, err
    }

	_, err = s.db.Exec(`
        UPDATE users SET session_token_hash = ?, session_expires_at = ? WHERE id = ?
    `, sessionTokenHash, expiresAt, user_id)
    if err != nil {
        return "", time.Time{}, err
    }

    return sessionToken, expiresAt, err
}

func (s *SQLiteStore) GetUserIdFromSessionToken(sessionToken string) (int, error) {
    if sessionToken == "" {
        return 0, errors.New("session token is empty")
    }

    rows, err := s.db.Query(`SELECT id, session_token_hash, session_expires_at FROM users`)
    if err != nil {
        return 0, err
    }
    defer rows.Close()

    for rows.Next() {
        var userId int
        var hash sql.NullString
        var expiresAt time.Time

        if err := rows.Scan(&userId, &hash, &expiresAt); err != nil {
            return 0, err
        }

        if !hash.Valid {
            continue // Skip users with no session token.
        }

        if err := bcrypt.CompareHashAndPassword([]byte(hash.String), []byte(sessionToken)); err == nil {
            if time.Now().After(expiresAt) {
                return 0, errors.New("session token has expired")
            }
            return userId, nil
        }
    }

    return 0, errors.New("session token not found")
}

func (s *SQLiteStore) GetTaskById(id string) (app.Task, error) {
    var t app.Task
    err := s.db.QueryRow(`SELECT id, title, description, done, due FROM tasks WHERE id = ?`, id).
        Scan(&t.Id, &t.Title, &t.Description, &t.Done, &t.Due)
    t.SetStatus()
    return t, err
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) GetAllTasks(user_id int) ([]app.Task, error) {
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

func (s *SQLiteStore) DeleteTask(id string) error {
    _, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
    return err
}

func (s *SQLiteStore) UpsertTask(t app.Task) error {
    _, err := s.db.Exec(`
        INSERT INTO tasks (id, user_id, title, description, done, due)
        VALUES (?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            title = excluded.title,
            description = excluded.description,
            done = excluded.done,
            due = excluded.due
    `, t.Id, t.UserId, t.Title, t.Description, t.Done, t.Due)
    return err
}

func (s *SQLiteStore) SetTaskDone(id string) error {
    _, err := s.db.Exec(`UPDATE tasks SET done = 1 WHERE id = ?`, id)
    return err
}