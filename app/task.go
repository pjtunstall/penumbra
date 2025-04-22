package app

import "time"

type User struct {
    ID   int
    Name string
    Phone string
    Email string
    PasswordHash string
}

type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    Status      string    `json:"status"`
    Due         time.Time `json:"due"`
}
