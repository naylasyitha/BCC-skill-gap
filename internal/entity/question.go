package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid"`
	SkillID         uuid.UUID `gorm:"type:uuid;not null"`
	Skill           Skill
	Level           LevelEnum `gorm:"type:varchar(50);not null"`
	QuestionContent string    `gorm:"type:text;not null"`
	OptionA         string    `gorm:"type:varchar(255);not null"`
	OptionB         string    `gorm:"type:varchar(255);not null"`
	OptionC         string    `gorm:"type:varchar(255);not null"`
	OptionD         string    `gorm:"type:varchar(255);not null"`
	Answer          string    `gorm:"type:char(1);not null"`
	Explanation     string    `gorm:"type:text"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}
