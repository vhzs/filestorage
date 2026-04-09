package model

import "time"

type File struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Path      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
