package dto

import "project/models"

type CreateTransactionRequest struct {
	CounterQty int    `json:"counter_qty" form:"counter_qty"`
	Total      int    `json:"total" form:"total"`
	Status     string `json:"status" form:"status"`
	TripID     int    `json:"tripId" form:"tripId"`
	UserId     int    `json:"userId" form:"userId"`
	// Image      string `json:"image" form:"image"`
}

type UpdateTransactionRequest struct {
	Status string `json:"status" form:"status"`
}

type TransactionResponse struct {
	Id          int                 `json:"id" gorm:"primary_key:auto_increment"`
	CounterQty  int                 `json:"counter_qty"`
	Total       int                 `json:"total"`
	Status      string              `json:"status"`
	BookingDate string              `json:"booking_date"`
	Token       string              `json:"token"`
	Trip        TripResponse        `json:"trip" gorm:"foreignKey:TripID"`
	User        models.UserResponse `json:"user"`
	// Image      string `json:"image" form:"image"`
}
