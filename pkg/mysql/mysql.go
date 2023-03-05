package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// variable DB akan dipanggil di untuk migrasi database (database/migration)
var DB *gorm.DB

func DatabaseInit() {
	var err error

	// username:root, password:kosong, route:localhost:3306, database name:projects
	dsn := "root:SJICET9zDR0xrEjJic89@tcp(containers-us-west-134.railway.app:5847)/railway?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// dsn := "root:@tcp(localhost:3306)/dewetour?charset=utf8mb4&parseTime=True&loc=Local"
	// DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Database")
}
