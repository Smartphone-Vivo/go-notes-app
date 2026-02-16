package domain

import "time"

type Note struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
