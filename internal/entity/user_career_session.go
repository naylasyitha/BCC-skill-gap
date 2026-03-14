package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserCareerSession struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	User        User
	CareerID    uuid.UUID `gorm:"type:uuid;not null"`
	Career      Career
	Status      StatusEnum `gorm:"type:varchar(50);default:'on_process'"`
	StartedAt   time.Time  `gorm:"autoCreateTime"`
	CompletedAt *time.Time // pakai pointer agar dapat null
}

func (us *UserCareerSession) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}
