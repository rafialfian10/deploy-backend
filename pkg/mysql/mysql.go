package mysql

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// variable DB akan dipanggil di untuk migrasi database (database/migration)
var DB *gorm.DB

func DatabaseInit() {
	var err error

	// username:root, password:kosong, route:localhost:3306, database name:projects
	// dsn := "root:xxumSGtJulD69dkIHhZ9@tcp(containers-us-west-134.railway.app:5847)/railway?charset=utf8mb4&parseTime=True&loc=Local"
	// DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// dsn := "root:@tcp(localhost:3306)/dewetour?charset=utf8mb4&parseTime=True&loc=Local"
	// DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	var DB_HOST = os.Getenv("DB_HOST")
	var DB_USER = os.Getenv("DB_USER")
	var DB_PASSWORD = os.Getenv("DB_PASSWORD")
	var DB_NAME = os.Getenv("DB_NAME")
	var DB_PORT = os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Database")
}
