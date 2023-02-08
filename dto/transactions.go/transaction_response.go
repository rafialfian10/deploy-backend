package transactionsdto

import "project/models"

type TransactionResponse struct {
	Id          int                 `json:"id"`
	CounterQty  int                 `json:"qty"`
	Token       string              `json:"token"`
	Total       int                 `json:"total"`
	Status      string              `json:"status"`
	BookingDate string              `json:"booking_date"`
	Trip        models.TripResponse `json:"trip"`
	User        models.UserResponse `json:"user"`
	// Image      string `json:"image" form:"image"`
}
