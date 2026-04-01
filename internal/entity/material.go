package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Material struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	SkillID     uuid.UUID `gorm:"type:uuid;not null"`
	Skill       Skill
	Level       LevelEnum `gorm:"type:varchar(50)"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Content     string    `gorm:"type:text"`
	VideoUrl    string    `gorm:"type:varchar(500)"`
	OrderNumber int       `gorm:"not null"`
}

func (m *Material) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
