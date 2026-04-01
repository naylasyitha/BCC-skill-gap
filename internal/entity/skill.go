package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Skill struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name string    `gorm:"type:varchar(255);not null"`
	Desc string    `gorm:"column:description;type:text" json:"desc"`
}

func (s *Skill) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
