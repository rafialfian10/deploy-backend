package repositories

import "gorm.io/gorm"

// struct gorm akan di panggil ke semua repositories untuk koneksi ke gorm
type repository struct {
	db *gorm.DB
}
