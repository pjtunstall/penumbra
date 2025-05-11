package db

import (
	"bytes"
	"database/sql"
	"testing"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"penumbra/app"
)

func TestCreateUser(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE users (
        id INTEGER PRIMARY KEY,
        name TEXT,
        password_hash BLOB,
        email TEXT,
        phone TEXT,
        session_token_hash TEXT,
        session_expires_at TIMESTAMP
    )`)
    if err != nil {
        t.Fatalf("failed to create table: %v", err)
    }

    store := &SQLiteStore{db: db}
    user := app.User{
        Name:         "Alice",
        PasswordHash: []byte("hashed-password"),
        Email:        "alice@example.com",
        Phone:        "1234567890",
    }

    err = store.CreateUser(user)
    if err != nil {
        t.Errorf("CreateUser returned error: %v", err)
    }

    var name, email, phone string
    var passwordHash []byte
    var sessionToken string
    var sessionExpires time.Time
    row := db.QueryRow(`
        SELECT name, password_hash, email, phone, session_token_hash, session_expires_at 
        FROM users WHERE email = ?`, user.Email)
    err = row.Scan(&name, &passwordHash, &email, &phone, &sessionToken, &sessionExpires)
    if err != nil {
        t.Errorf("failed to retrieve user: %v", err)
    }

    if name != user.Name || email != user.Email || phone != user.Phone {
        t.Errorf("retrieved user doesn't match input")
    }

    if !bytes.Equal(passwordHash, user.PasswordHash) {
        t.Errorf("expected password hash %v, got %v", user.PasswordHash, passwordHash)
    }

    if sessionToken != "" {
        t.Errorf("expected empty session token hash, got %q", sessionToken)
    }

    if !sessionExpires.Equal(time.Unix(0, 0)) {
        t.Errorf("expected session expires at Unix(0, 0), got %v", sessionExpires)
    }
}

func TestGetUserByEmail(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        password_hash BLOB,
        email TEXT UNIQUE,
        phone TEXT,
        session_token_hash TEXT,
        session_expires_at TIMESTAMP
    )`)
    if err != nil {
        t.Fatalf("failed to create table: %v", err)
    }

    store := &SQLiteStore{db: db}
    expectedUser := app.User{
        Name:         "Bob",
        PasswordHash: []byte("secure-hash"),
        Email:        "bob@example.com",
        Phone:        "9876543210",
    }

    err = store.CreateUser(expectedUser)
    if err != nil {
        t.Fatalf("CreateUser failed: %v", err)
    }

    user, err := store.GetUserByEmail(expectedUser.Email)
    if err != nil {
        t.Fatalf("GetUserByEmail returned error: %v", err)
    }

    if user.Name != expectedUser.Name || user.Email != expectedUser.Email || user.Phone != expectedUser.Phone {
        t.Errorf("retrieved user doesn't match inserted user")
    }

    if !bytes.Equal(user.PasswordHash, expectedUser.PasswordHash) {
        t.Errorf("expected password hash %v, got %v", expectedUser.PasswordHash, user.PasswordHash)
    }

    if user.Id == 0 {
        t.Errorf("expected non-zero user ID")
    }
}

func TestAddSessionToken(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        password_hash BLOB,
        email TEXT,
        phone TEXT,
        session_token_hash BLOB,
        session_expires_at TIMESTAMP
    )`)
    if err != nil {
        t.Fatalf("failed to create table: %v", err)
    }

    store := &SQLiteStore{db: db}

    user := app.User{
        Name:         "Charlie",
        PasswordHash: []byte("passhash"),
        Email:        "charlie@example.com",
        Phone:        "1112223333",
    }

    err = store.CreateUser(user)
    if err != nil {
        t.Fatalf("CreateUser failed: %v", err)
    }

    var userID int
    err = db.QueryRow(`SELECT id FROM users WHERE email = ?`, user.Email).Scan(&userID)
    if err != nil {
        t.Fatalf("failed to get user ID: %v", err)
    }

    token, expiresAt, err := store.AddSessionToken(userID)
    if err != nil {
        t.Fatalf("AddSessionToken failed: %v", err)
    }

    var hashFromDB []byte
    var expiresFromDB time.Time
    err = db.QueryRow(`
        SELECT session_token_hash, session_expires_at FROM users WHERE id = ?`, userID).
        Scan(&hashFromDB, &expiresFromDB)
    if err != nil {
        t.Fatalf("failed to fetch session info: %v", err)
    }

    if err := bcrypt.CompareHashAndPassword(hashFromDB, token[:]); err != nil {
        t.Errorf("session token hash doesn't match: %v", err)
    }

    if !expiresAt.Equal(expiresFromDB) {
        t.Errorf("expiresAt mismatch: expected %v, got %v", expiresAt, expiresFromDB)
    }
}

func TestGetTaskById(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE tasks (
        id TEXT PRIMARY KEY,
        title TEXT,
        description TEXT,
        done BOOLEAN,
        due TIMESTAMP
    )`)
    if err != nil {
        t.Fatalf("failed to create tasks table: %v", err)
    }

    store := &SQLiteStore{db: db}

    id := uuid.New()
    due := time.Now().Add(48 * time.Hour)

    _, err = db.Exec(`
        INSERT INTO tasks (id, title, description, done, due)
        VALUES (?, ?, ?, ?, ?)`,
        id.String(), "Test Task", "Do the thing", false, due)
    if err != nil {
        t.Fatalf("failed to insert test task: %v", err)
    }

    task, err := store.GetTaskById(id)
    if err != nil {
        t.Fatalf("GetTaskById failed: %v", err)
    }

    if task.Id != id {
        t.Errorf("expected ID %v, got %v", id, task.Id)
    }
    if task.Title != "Test Task" || task.Description != "Do the thing" || task.Done != 0 {
        t.Errorf("unexpected task fields: %+v", task)
    }
    if !task.Due.Equal(due) {
        t.Errorf("expected due %v, got %v", due, task.Due)
    }
}

func TestGetAllTasks(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    store := &SQLiteStore{db: db}

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id BLOB PRIMARY KEY,
            user_id INTEGER NOT NULL,
            title TEXT NOT NULL,
            description TEXT,
            done INTEGER NULL,
            due DATETIME NOT NULL
        )`)
    if err != nil {
        t.Fatalf("failed to create table: %v", err)
    }

    userID := 1
    taskID1 := uuid.New()
    taskID2 := uuid.New()

    _, err = db.Exec(`
        INSERT INTO tasks (id, user_id, title, description, done, due)
        VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)`,
        taskID1[:], userID, "Test Task 1", "Description 1", 0, "2025-05-13T00:00:00Z",
        taskID2[:], userID, "Test Task 2", "Description 2", 1, "2025-05-14T00:00:00Z")
    if err != nil {
        t.Fatalf("failed to insert tasks: %v", err)
    }

    rows, err := store.GetAllTasks(userID)
    if err != nil {
        t.Fatalf("failed to get all tasks: %v", err)
    }

    if len(rows) != 2 {
        t.Fatalf("expected 2 tasks, got %d", len(rows))
    }

    if rows[0].Title != "Test Task 1" || rows[1].Title != "Test Task 2" {
        t.Fatalf("unexpected task titles: got %v", rows)
    }

    if rows[0].Due.IsZero() || rows[1].Due.IsZero() {
        t.Fatalf("expected non-zero due dates, got %v", rows)
    }
}

func TestCreateTask(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE tasks (
        id TEXT PRIMARY KEY,
        user_id INTEGER,
        title TEXT,
        description TEXT,
        done BOOLEAN,
        due TIMESTAMP
    )`)
    if err != nil {
        t.Fatalf("failed to create tasks table: %v", err)
    }

    store := &SQLiteStore{db: db}

    id := uuid.New()
    due := time.Now().Add(48 * time.Hour).UTC() // Convert to UTC to avoid mismatch between how Go stores the time (including timezone information) and SQLite, which doesn't.

    task := app.Task{
        Id:          id,
        UserId:      1,
        Title:       "Test Task",
        Description: "Do the thing",
        Done:        0,
        Due:         due,
    }

    err = store.CreateTask(task)
    if err != nil {
        t.Fatalf("CreateTask failed: %v", err)
    }

    var dbTask app.Task
    err = db.QueryRow(`
        SELECT id, user_id, title, description, done, due
        FROM tasks WHERE id = ?`, id.String()).Scan(&dbTask.Id, &dbTask.UserId, &dbTask.Title, &dbTask.Description, &dbTask.Done, &dbTask.Due)
    if err != nil {
        t.Fatalf("failed to fetch task: %v", err)
    }

    if !dbTask.Due.Equal(task.Due) {
        t.Errorf("expected due %v, got %v", task.Due, dbTask.Due)
    }

    if dbTask != task {
        t.Errorf("expected task %+v, got %+v", task, dbTask)
    }
}

func TestDeleteTask(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE tasks (
        id TEXT PRIMARY KEY,
        user_id INTEGER,
        title TEXT,
        description TEXT,
        done BOOLEAN,
        due TIMESTAMP
    )`)
    if err != nil {
        t.Fatalf("failed to create tasks table: %v", err)
    }

    store := &SQLiteStore{db: db}

    id := uuid.New()
    due := time.Now().Add(48 * time.Hour)

    task := app.Task{
        Id:          id,
        UserId:      1,
        Title:       "Test Task",
        Description: "Do the thing",
        Done:        0,
        Due:         due,
    }

    err = store.CreateTask(task)
    if err != nil {
        t.Fatalf("CreateTask failed: %v", err)
    }

    err = store.DeleteTask(id)
    if err != nil {
        t.Fatalf("DeleteTask failed: %v", err)
    }

    var count int
    err = db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE id = ?`, id.String()).Scan(&count)
    if err != nil {
        t.Fatalf("failed to query tasks: %v", err)
    }

    if count != 0 {
        t.Errorf("expected task to be deleted, but found %d task(s)", count)
    }
}
