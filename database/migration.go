package database

import (
	"fmt"
	"project/models"
	"project/pkg/mysql"
)

// Jika aplikasi berjalan maka auto migration akan berjalan
func RunMigration() {
	// koneksi database akan melakukan auto migrasi struct/models ke dalam database mysql
	err := mysql.DB.AutoMigrate( // panggil mysql lalu DB(pkg/mysql) lalu panggil function AutoMigrate()
		&models.User{},
		&models.Trip{},
		&models.Country{},
		&models.Transaction{},
	)
	// jika ada error maka panggil panic
	if err != nil {
		fmt.Println(err)
		panic("Migration failed")
	}

	fmt.Println("Migration success")
}
