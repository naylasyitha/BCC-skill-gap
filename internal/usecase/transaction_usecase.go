package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"project-bcc/dto"
	"project-bcc/internal/entity"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentUsecase struct {
	paymentRepo PaymentRepository
	serverKey   string
}

func NewPaymentUsecase(repo PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: repo,
		serverKey:   os.Getenv("MIDTRANS_SERVER_KEY"),
	}
}

func (u *PaymentUsecase) CreatePayment(ctx context.Context, userID string) (*dto.PaymentResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("User ID tidak valid")
	}

	var s snap.Client
	s.New(u.serverKey, midtrans.Sandbox)

	orderID := fmt.Sprintf("ORD-%s-%d", userUUID.String()[:8], time.Now().Unix())
	amount := int64(50000)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
	}

	snapResp, midErr := s.CreateTransaction(req)
	if midErr != nil {
		return nil, errors.New("Gagal terhubung ke payment gateway")
	}

	tx := &entity.Transaction{
		ID:        orderID,
		UserID:    userUUID,
		Amount:    amount,
		Status:    entity.StatusPending,
		SnapToken: snapResp.Token,
		SnapURL:   snapResp.RedirectURL,
	}

	if err := u.paymentRepo.CreateTransaction(ctx, tx); err != nil {
		return nil, errors.New("Gagal menyimpan transaksi")
	}

	return &dto.PaymentResponse{
		OrderID:     orderID,
		SnapToken:   snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

func (u *PaymentUsecase) HandleNotification(ctx context.Context, req dto.MidtransNotificationRequest) error {
	var c coreapi.Client
	c.New(u.serverKey, midtrans.Sandbox)

	transStatus, midErr := c.CheckTransaction(req.OrderID)
	if midErr != nil {
		return errors.New("Gagal memverifikasi transaksi ke Midtrans")
	}

	tx, err := u.paymentRepo.GetTransactionByID(ctx, req.OrderID)
	if err != nil {
		return errors.New("Transaksi tidak ditemukan")
	}

	if tx.Status == entity.StatusSettlement {
		return nil
	}

	var newStatus entity.TransactionStatus
	switch transStatus.TransactionStatus {
	case "capture", "settlement":
		newStatus = entity.StatusSettlement
	case "deny", "cancel", "expire", "failure":
		newStatus = entity.StatusExpire
	case "pending":
		newStatus = entity.StatusPending
	default:
		return nil
	}

	if err := u.paymentRepo.UpdateTransactionStatus(ctx, req.OrderID, newStatus); err != nil {
		return err
	}

	if newStatus == entity.StatusSettlement {
		if err := u.paymentRepo.UpgradeUserAccount(ctx, tx.UserID.String()); err != nil {
			return errors.New("Gagal mengupgrade akun user")
		}
	}

	return nil
}

func (u *PaymentUsecase) GetPaymentStatus(ctx context.Context, orderID string) (*entity.Transaction, error) {
	tx, err := u.paymentRepo.GetTransactionByID(ctx, orderID)
	if err != nil {
		return nil, errors.New("Transaksi tidak ditemukan")
	}
	return tx, nil
}
