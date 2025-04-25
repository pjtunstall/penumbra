package db

import (
	"database/sql"
	"log"
	"time"

	"dts/app"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
    CreateUser(user app.User) error
    AddSessionToken(user_id int) (string, time.Time, error)
    GetUserIDFromSessionToken(sessionToken string) (int, error)
    CreateTask(task app.Task) (int, error)
    GetTask(id int) (app.Task, error)
    GetUserByEmail(email string) (app.User, error)
    ListTasks(user_id int) ([]app.Task, error)
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
        title TEXT NOT NULL,
        description TEXT,
        status TEXT NOT NULL,
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
        Scan(&user.ID, &user.Name, &user.PasswordHash, &user.Email, &user.Phone)
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

func (s *SQLiteStore) GetUserIDFromSessionToken(sessionToken string) (int, error) {
    if sessionToken == "" {
        return 0, nil
    }

    sessionTokenHash, err := bcrypt.GenerateFromPassword([]byte(sessionToken), 10)
    if err != nil {
        return 0, err
    }

    var userID int
    err = s.db.QueryRow(`SELECT user_id FROM users WHERE session_token_hash = ?`, sessionTokenHash).
        Scan(&userID)
    return userID, err
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

func (s *SQLiteStore) ListTasks(user_id int) ([]app.Task, error) {
    rows, err := s.db.Query(`SELECT id, title, description, status, due FROM tasks WHERE user_id = ?`, user_id)
    if err != nil {
        return nil, err
    }
    tasks := []app.Task{}
    for rows.Next() {
        var t app.Task
        err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Due)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }
    return tasks, nil
}

// Add UpdateTask, DeleteTask, ListTasks later
