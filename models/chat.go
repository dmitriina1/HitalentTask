package models

import "time"

type Chat struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:200;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Messages  []Message `json:"messages" gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}
