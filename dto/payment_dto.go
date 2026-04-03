package dto

type PaymentResponse struct {
	OrderID     string `json:"order_id"`
	SnapToken   string `json:"snap_token"`
	RedirectURL string `json:"redirect_url"`
}

type MidtransNotificationRequest struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	FraudStatus       string `json:"fraud_status"`
}
