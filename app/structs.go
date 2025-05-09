package app

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
    Id   int
    Name string
    Phone string
    Email string
    PasswordHash []byte
}

type Task struct {
    Id          uuid.UUID   `json:"id"`
    UserId      int `json:"userId"`
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    Status      string    `json:"status"`
    Done        int      `json:"done"`
    Due         time.Time `json:"due"`
    Uuid        string  `json:"uuid"`
}

func (t *Task) SetStatus() {
    if t.Done == 1 {
        t.Status = "done"
    } else {
        if t.Due.Before(time.Now()) {
            t.Status = "overdue"
        } else {
            t.Status = "pending"
        }
    }
}
