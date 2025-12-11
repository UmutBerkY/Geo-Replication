package model

import "time"

type Article struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	ContentLong string    `json:"content_long"`
	Author      string    `json:"author"`
	Region      string    `json:"region"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateArticleInput struct {
	Title       string `json:"title" binding:"required"`
	Summary     string `json:"summary" binding:"required"`
	ContentLong string `json:"content_long" binding:"required"`
	Author      string `json:"author" binding:"required"`
}

type ReplicationStatus struct {
	Replica string    `json:"replica"`
	Status  string    `json:"status"`
	LastAt  time.Time `json:"last_at"`
}
