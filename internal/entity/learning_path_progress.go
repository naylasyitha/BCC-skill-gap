package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LearningPathProgress struct {
	ID                  uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID              uuid.UUID `gorm:"type:uuid;not null"`
	User                User
	UserCareerSessionID uuid.UUID `gorm:"type:uuid;not null"`
	UserCareerSession   UserCareerSession
	MaterialID          uuid.UUID `gorm:"type:uuid;not null"`
	Material            Material
	Status              StatusEnum `gorm:"type:varchar(50);default:'not_started'"`
	CompletedAt         *time.Time
}

func (l *LearningPathProgress) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
