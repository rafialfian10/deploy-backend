package models

import "time"

type Transaction struct {
	Id          int          `json:"id" gorm:"primary_key:auto_increment"`
	CounterQty  int          `json:"counter_qty" gorm:"type: int"`
	Total       int          `json:"total" gorm:"type: int"`
	BookingDate time.Time    `json:"booking_date"`
	Status      string       `json:"status" form:"status" gorm:"type: varchar(255)"`
	Token       string       `json:"token" gorm:"type: varchar(255)"`
	Image       string       `json:"image" form:"image" gorm:"type: varchar(255)"`
	TripId      int          `json:"-"`
	UserId      int          `json:"-"`
	Trip        TripResponse `json:"trip"`
	User        UserResponse `json:"user"`
}

// type TransactionResponse struct {
// 	Id          int          `json:"id"`
// 	CounterQty  int          `json:"counter_qty" gorm:"type: int"`
// 	Total       int          `json:"total" gorm:"type: int"`
// 	BookingDate time.Time    `json:"booking_date"`
// 	Status      string       `json:"status" gorm:"type: varchar(255)"`
// 	Token       string       `json:"token" gorm:"type: varchar(255)"`
// 	TripId      int          `json:"-"`
// 	UserId      int          `json:"-"`
// 	Trip        TripResponse `json:"trip"`
// 	User        UserResponse `json:"user"`
// }

// func (TransactionResponse) TableName() string {
// 	return "transactions"
// }
