package model

import "time"

type Post struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Body      string    `json:"body"`
    UserID    int       `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}