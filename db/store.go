package db

import (
	"database/sql"
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
    SubmitCreateTask(task app.Task) (int, error)
    GetTaskById(id int) (app.Task, error)
    GetUserByEmail(email string) (app.User, error)
    GetAllTasks(user_id int) ([]app.Task, error)
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

    createTable := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        description TEXT,
        done INTEGER NOT NULL,
        due DATETIME NOT NULL
    )`
    if _, err := db.Exec(createTable); err != nil {
        log.Fatal(err)
    }

    return &SQLiteStore{db: db, Path: path}
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
        return 0, nil // todo: make this an error
    }
    var userId int
    err := s.db.QueryRow(`SELECT id FROM users WHERE session_token = ?`, sessionToken).
        Scan(&userId)
    return userId, err
}

func (s *SQLiteStore) SubmitCreateTask(t app.Task) (int, error) {
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
    rows, err := s.db.Query(`SELECT id, title, description, due FROM tasks WHERE user_id = ?`, user_id)
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

// Add UpdateTask, DeleteTask, GetAllTasks later
