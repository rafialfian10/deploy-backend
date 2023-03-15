package repositories

import (
	"project/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactions() ([]models.Transaction, error)
	FindTransactionsByUser(UserId int) ([]models.Transaction, error)
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

func (r *repository) FindTransactionsByUser(UserId int) ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").Where("user_id = ?", UserId).Order("booking_date desc").Find(&transaction).Error

	return transaction, err
}

func (r *repository) GetTransaction(Id int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Trip.Country").Preload("Trip").Preload("User").First(&transaction, "id = ?", Id).Error

	return transaction, err
}

func (r *repository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("Trip.Country").Preload("Trip").Preload("User").Create(&transaction).Error

	return transaction, err
}

func (r *repository) UpdateTransaction(status string, Id int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Trip.Country").Preload("Trip").Preload("User").First(&transaction, "id = ?", Id)

	// jika status dan transaksi status berbeda & Status adalah "reject" maka quota trip akan dikurangi
	if status != transaction.Status && status == "success" {
		var trip models.Trip
		r.db.First(&trip, transaction.TripID)
		trip.Quota = trip.Quota - transaction.CounterQty
		r.db.Model(&trip).Updates(trip)
	}

	// jika status dan transaksi status berbeda & Status adalah "reject" maka quota trip akan tetap
	if status != transaction.Status && status == "reject" {
		var trip models.Trip
		r.db.First(&trip, transaction.TripID)
		trip.Quota = trip.Quota + transaction.CounterQty
		r.db.Model(&trip).Updates(trip)
	}

	// mengubah transaction status
	transaction.Status = status
	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) UpdateTokenTransaction(token string, Id int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Trip.Country").Preload("Trip").Preload("User").First(&transaction, "id = ?", Id)

	// mengubah transaction token
	transaction.Token = token

	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) DeleteTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("Trip").Delete(&transaction).Error

	return transaction, err
}
