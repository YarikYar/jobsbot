package models

import "time"

type Post struct {
	Id        int       `json:"id,omitempty" db:"id"`
	Type      string    `json:"type,omitempty" db:"post_type"`
	Title     string    `json:"title,omitempty" db:"title"`
	Link      string    `json:"link,omitempty" db:"post_link"`
	MessageId int       `json:"message_id,omitempty" db:"post_id"`
	UserId    int64     `json:"user_id,omitempty" db:"user_id"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
