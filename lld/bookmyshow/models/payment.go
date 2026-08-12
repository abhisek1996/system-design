package models

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodPayPal     PaymentMethod = "PAYPAL"
	PaymentMethodNetBanking PaymentMethod = "NET_BANKING"
)

type Payment struct {
	Id        int
	BookingId int
	Amount    float64
	Status    PaymentStatus
	Date      time.Time
	Method    PaymentMethod // ENUM: Credit Card, PayPal, etc.
}
