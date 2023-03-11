package models

import "time"

type Transaction struct {
	Id          string    `json:"id" gorm:"type: varchar(255); PRIMARY_KEY"`
	CounterQty  int       `json:"qty" gorm:"type: int"`
	Total       int       `json:"total" gorm:"type: int"`
	BookingDate time.Time `json:"booking_date"`
	Status      string    `json:"status" form:"status" gorm:"type: varchar(255)"`
	Token       string    `json:"token" gorm:"type: varchar(255)"`
	Image       string    `json:"image" form:"image" gorm:"type: varchar(255)"`
	TripId      int       `json:"trip_id" gorm:"type: int"`
	UserId      int       `json:"user_id" gorm:"type: int"`
	Trip        TripResponse
	User        UserResponse
}

type TransactionResponse struct {
	Id          int       `json:"id"`
	CounterQty  int       `json:"qty" gorm:"type: int"`
	BookingDate time.Time `json:"booking_date"`
	Total       int       `json:"total" gorm:"type: int"`
	Status      string    `json:"status" gorm:"type: varchar(255)"`
	TripId      int       `json:"trip_id" gorm:"type: int"`
	UserId      int       `json:"user_id" gorm:"type: int"`
}

func (TransactionResponse) TableName() string {
	return "transactions"
}
