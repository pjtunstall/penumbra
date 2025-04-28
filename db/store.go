package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"dts/app"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
)

type Store interface {
    CreateUser(user app.User) error
    AddSessionToken(user_id int) (string, time.Time, error)
    GetUserIdFromSessionToken(sessionToken string) (int, error)
    GetTaskById(id int) (app.Task, error)
    SetTaskDone(id int) error
    GetUserByEmail(email string) (app.User, error)
    GetAllTasks(user_id int) ([]app.Task, error)
    SubmitCreate(task app.Task) (int, error)
    UpdateTask(task app.Task) error
    DeleteTask(id int) error
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
    _, err := s.db.Exec(`INSERT INTO users (name, password_hash, email, phone) VALUES (?, ?, ?, ?)`,
        user.Name, user.PasswordHash, user.Email, user.Phone)
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

	_, err := s.db.Exec(`
        UPDATE users SET session_token = ?, session_expires_at = ? WHERE id = ?
    `, sessionToken, expiresAt, user_id)
    if err != nil {
        return "", time.Time{}, err
    }

    return sessionToken, expiresAt, err
}

func (s *SQLiteStore) GetUserIdFromSessionToken(sessionToken string) (int, error) {
    if sessionToken == "" {
        return 0, errors.New("session token is empty")
    }

    var userId int
    var sessionExpiresAt time.Time

    err := s.db.QueryRow(`SELECT id, session_expires_at FROM users WHERE session_token = ?`, sessionToken).
        Scan(&userId, &sessionExpiresAt)

    if time.Now().After(sessionExpiresAt) {
        return 0, errors.New("session token has expired")
    }

    return userId, err
}

func (s *SQLiteStore) SubmitCreate(t app.Task) (int, error) {
    res, err := s.db.Exec(`INSERT INTO tasks (user_id, title, description, done, due) VALUES (?, ?, ?, ?, ?)`,
        t.UserId, t.Title, t.Description, t.Done, t.Due)
    if err != nil {
        return 0, err
    }
    id, _ := res.LastInsertId()
    return int(id), nil
}

func (s *SQLiteStore) GetTaskById(id int) (app.Task, error) {
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

func (s *SQLiteStore) DeleteTask(id int) error {
    _, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
    return err
}

func (s *SQLiteStore) UpdateTask(t app.Task) error {
    _, err := s.db.Exec(`UPDATE tasks SET title = ?, description = ?, done = ?, due = ? WHERE id = ?`,
        t.Title, t.Description, t.Done, t.Due, t.Id)
    return err
}

func (s *SQLiteStore) SetTaskDone(id int) error {
    _, err := s.db.Exec(`UPDATE tasks SET done = 1 WHERE id = ?`, id)
    return err
}