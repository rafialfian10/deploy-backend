package dto

import "project/models"

type CreateTransactionRequest struct {
	CounterQty int    `json:"counter_qty" form:"counter_qty"`
	Total      int    `json:"total" form:"total"`
	Status     string `json:"status" form:"status"`
	TripId     int    `json:"trip_id" form:"trip_id"`
	UserId     int    `json:"user_id" form:"user_id"`
	// Image      string `json:"image" form:"image"`
}

type UpdateTransactionRequest struct {
	Status string `json:"status" form:"status"`
}

type TransactionResponse struct {
	Id          int                 `json:"id"`
	CounterQty  int                 `json:"counter_qty"`
	Token       string              `json:"token" gorm:"type: varchar(255)"`
	Total       int                 `json:"total"`
	Status      string              `json:"status"`
	BookingDate string              `json:"booking_date"`
	Trip        TripResponse        `json:"trip"`
	User        models.UserResponse `json:"user"`
	// Image      string `json:"image" form:"image"`
}
