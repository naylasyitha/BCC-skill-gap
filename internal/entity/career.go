package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Career struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name string    `gorm:"type:varchar(255);not null"`
	Desc string    `gorm:"type:text"`
}

func (c *Career) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil

}
