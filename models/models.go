package models

import "time"

type Question struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Answer struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID uint      `gorm:"not null;index" json:"question_id"`
	UserID     string    `gorm:"type:uuid;not null" json:"user_id"`
	Text       string    `gorm:"type:text;not null" json:"text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
