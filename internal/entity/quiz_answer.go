package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizAnswer struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid"`
	QuizSessionID uuid.UUID `gorm:"type:uuid;not null"`
	QuizSession   QuizSession
	QuestionID    uuid.UUID `gorm:"type:uuid;not null"`
	Question      Question
	UserAnswer    string `gorm:"type:char(1)"`
	IsCorrect     bool   `gorm:"default:false"`
}

func (q *QuizAnswer) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}
