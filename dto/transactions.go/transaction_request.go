package transactionsdto

type CreateTransactionRequest struct {
	CounterQty int    `json:"qty" form:"qty"`
	Total      int    `json:"total" form:"total"`
	Status     string `json:"status" form:"status"`
	// Image      string `json:"image" form:"image"`
	TripId int `json:"trip_id" form:"trip_id"`
	UserId int `json:"user_id" form:"user_id"`
}

type UpdateTransactionRequest struct {
	CounterQty int    `json:"qty" form:"qty"`
	Total      int    `json:"total" form:"total"`
	Status     string `json:"status" form:"status"`
	// Image      string `json:"image" form:"image"`
	TripId int `json:"trip_id" form:"trip_id"`
	UserId int `json:"user_id" form:"user_id"`
}
