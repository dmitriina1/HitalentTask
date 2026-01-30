package models

import "time"

type Message struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	ChatID    int       `json:"chat_id" gorm:"not null;index"`
	Text      string    `json:"text" gorm:"size:5000;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
