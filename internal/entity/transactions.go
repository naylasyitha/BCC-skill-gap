package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    uuid.UUID `gorm:"column:user_id"`
	User      User
	Amount    int64             `gorm:"column:amount"`
	Status    TransactionStatus `gorm:"column:status"`
	SnapToken string            `gorm:"column:snap_token"`
	SnapURL   string            `gorm:"column:snap_url"`
	CreatedAt time.Time         `gorm:"column:created_at"`
	UpdatedAt time.Time         `gorm:"column:updated_at"`
}
