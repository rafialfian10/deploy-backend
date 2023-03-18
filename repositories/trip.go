package repositories

import (
	"project/models"

	"gorm.io/gorm"
)

// membuat interface TripRepository
type TripRepository interface {
	FindTrips() ([]models.Trip, error)
	GetTrip(Id int) (models.Trip, error)
	CreateTrip(trip models.Trip) (models.Trip, error)
	UpdateTrip(trip models.Trip) (models.Trip, error)
	DeleteTrip(trip models.Trip) (models.Trip, error)
}

// membuat function RepositoryTrip. parameter pointer ke gorm, return repository{db}. ini akan dipanggil di routes
func RepositoriyTrip(db *gorm.DB) *repository {
	return &repository{db}
}

// membuat struct method FindTrips(memanggil struct dengan struct function)
func (r *repository) FindTrips() ([]models.Trip, error) {
	// panggil struct Trip lalu preload(berfungsi agar data dapat auto load saat create/update data)
	var trips []models.Trip
	err := r.db.Debug().Preload("Country").Preload("Image").Find(&trips).Error

	return trips, err
}

// membuat struct method GetTrip(memanggil struct dengan struct function)
func (r *repository) GetTrip(Id int) (models.Trip, error) {
	var trip models.Trip
	err := r.db.Debug().Preload("Country").Preload("Image").First(&trip, Id).Error

	return trip, err
}

// membuat struct method CreateTrip(memanggil struct dengan struct function)
func (r *repository) CreateTrip(trip models.Trip) (models.Trip, error) {
	err := r.db.Debug().Preload("Country").Create(&trip).Error

	return trip, err
}

// membuat struct method UpdateTrip(memanggil struct dengan struct function)
func (r *repository) UpdateTrip(trip models.Trip) (models.Trip, error) {
	// err := r.db.Save(&trip).Error // jika bertemu duplikat key, tidak akan mengupdate

	// db.Session digunakan agar jika bertemu duplikat key, value dari key tersebut akan diupdate sesuai dengan yang terbaru
	// Model dan Updates digunakan karena tabel Trip memiliki relasi belongsto dengan country
	err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&trip).Updates(trip).Error

	// menghapus gambar yang tidak lagi terpakai
	r.db.Exec("DELETE from images where file_name = ?", "deleted")

	return trip, err
}

// membuat struct method Deletetrip(memanggil struct dengan struct function)
func (r *repository) DeleteTrip(trip models.Trip) (models.Trip, error) {
	// err := r.db.Debug().Preload("Country").Delete(&trip).Error
	err := r.db.Debug().Preload("Country").Select("Image").Delete(&trip).Error

	return trip, err
}
