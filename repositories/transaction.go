package repositories

import (
	"project/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactions() ([]models.Transaction, error)
	GetTransaction(Id int) (models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(status string, Id int) (models.Transaction, error)
	UpdateTokenTransaction(token string, Id int) (models.Transaction, error)
	DeleteTransaction(transaction models.Transaction) (models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").Find(&transaction).Error

	return transaction, err
}

func (r *repository) GetTransaction(Id int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").First(&transaction, Id).Error

	return transaction, err
}

func (r *repository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&transaction).Error

	return transaction, err
}

func (r *repository) UpdateTransaction(status string, Id int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Trip.Country").Preload("Trip.Image").Preload("Trip").Preload("User").First(&transaction, "id = ?", Id)

	// If is different & Status is "success" decrement available quota on data trip
	if status != transaction.Status && status == "success" {
		var trip models.Trip
		r.db.First(&trip, transaction.TripId)
		trip.Quota = trip.Quota - transaction.CounterQty
		r.db.Model(&trip).Updates(trip)
	}

	// If is different & Status is "reject" decrement available quota on data trip
	if status != transaction.Status && status == "reject" {
		var trip models.Trip
		r.db.First(&trip, transaction.TripId)
		trip.Quota = trip.Quota + transaction.CounterQty
		r.db.Model(&trip).Updates(trip)
	}

	// change transaction status
	transaction.Status = status

	// fmt.Println(status)
	// fmt.Println(transaction.Status)
	// fmt.Println(transaction.ID)

	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) UpdateTokenTransaction(token string, Id int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Trip").Preload("Trip.Country").Preload("User").First(&transaction, "id = ?", Id)

	// change transaction token
	transaction.Token = token

	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) DeleteTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("Trip").Delete(&transaction).Error

	return transaction, err
}
