package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid"`
	FullName       string    `gorm:"type:varchar(255);not null"`
	Username       string    `gorm:"type:varchar(255);not null;unique"`
	Email          string    `gorm:"type:varchar(255);unique;not null"`
	Password       string    `gorm:"type:varchar(500);not null"`
	Role           Role      `gorm:"type:varchar(50);not null;default:'user'"`
	EducationLevel EduLevel  `gorm:"type:varchar(50)"`
	Major          string    `gorm:"type:varchar(255)"`
	Institution    string    `gorm:"type:varchar(255)"`
	GraduationYear int
	IsPremium      bool      `gorm:"default:false"`
	RefreshToken   string    `gorm:"type:varchar(500)"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
