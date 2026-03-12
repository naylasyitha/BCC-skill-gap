package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid"`
	FullName       string    `gorm:"not null"`
	Username       string    `gorm:"not null"`
	Email          string    `gorm:"unique;not null"`
	Password       string    `gorm:"not null"`
	Role           string    `gorm:"not null;default:'user'"`
	EducationLevel string
	Major          string
	Institution    string
	GraduationYear int
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
