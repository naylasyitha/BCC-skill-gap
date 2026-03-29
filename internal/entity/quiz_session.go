package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizSession struct {
	ID                  uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID              uuid.UUID `gorm:"type:uuid;not null"`
	User                User
	UserCareerSessionID uuid.UUID `gorm:"type:uuid;not null"`
	UserCareerSession   UserCareerSession
	Status              StatusEnum `gorm:"type:varchar(50);default:'on_process'"`
	Score               float64    `gorm:"type:float8"`
	StartedAt           time.Time  `gorm:"autoCreateTime"`
	CompletedAt         *time.Time
}

func (q *QuizSession) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}
