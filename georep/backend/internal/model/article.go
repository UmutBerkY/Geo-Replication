package model

import "time"

type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Region    string    `json:"region"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateArticleInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author" binding:"required"`
}

type ReplicationStatus struct {
	Replica string    `json:"replica"`
	Status  string    `json:"status"`
	LastAt  time.Time `json:"last_at"`
}
