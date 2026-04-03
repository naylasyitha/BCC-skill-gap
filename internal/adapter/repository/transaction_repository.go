package repository

import (
	"context"
	"project-bcc/internal/entity"
	"project-bcc/internal/usecase"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) usecase.PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) CreateTransaction(ctx context.Context, tx *entity.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *paymentRepository) GetTransactionByID(ctx context.Context, orderID string) (*entity.Transaction, error) {
	var tx entity.Transaction
	err := r.db.WithContext(ctx).Where("id = ?", orderID).First(&tx).Error
	return &tx, err
}

func (r *paymentRepository) UpdateTransactionStatus(ctx context.Context, orderID string, status entity.TransactionStatus) error {
	return r.db.WithContext(ctx).Model(&entity.Transaction{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *paymentRepository) UpgradeUserAccount(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Table("users").Where("id = ?", userID).Update("is_premium", true).Error
}
