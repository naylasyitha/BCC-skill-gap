package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CareerSkill struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid"`
	CareerID      uuid.UUID `gorm:"type:uuid;not null"`
	Career        Career
	SkillID       uuid.UUID `gorm:"type:uuid;not null"`
	Skill         Skill
	Priority      int       `gorm:"not null"`
	RequiredLevel LevelEnum `gorm:"type:varchar(50);not null"`
}

func (c *CareerSkill) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
