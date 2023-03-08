package main

import (
	"fmt"
	"net/http"
	"os"
	"project/database"
	"project/pkg/mysql"
	"project/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// env
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	// route
	route := mux.NewRouter()

	// initial DB
	mysql.DatabaseInit()

	// run migration
	database.RunMigration()

	// route untuk menginisialisasi folder dengan file, image css, js agar dapat diakses kedalam project
	route.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// pathPrefix untuk membuat route baru. Subrouter untuk menguji route pada pathPrefix. RouteInit dari (routes/routes)
	routes.RouteInit(route.PathPrefix("/api/v1").Subrouter())

	// Setup header, method dan CORS
	var AllowedHeaders = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	var AllowedOrigins = handlers.AllowedOrigins([]string{"https://dewetour-frontend.netlify.app"})
	var AllowedMethods = handlers.AllowedMethods([]string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"})

	var PORT = os.Getenv("PORT")

	fmt.Println("server running localhost:5000")
	http.ListenAndServe(":"+PORT, handlers.CORS(AllowedHeaders, AllowedOrigins, AllowedMethods)(route))
}

// lifecycle: models ---> koneksi mysql ---> database migration ---> repositories ---> dto ---> handlers ---> routers
// run xampp: sudo /opt/lampp/lampp start

// Ngrok adalah proxy server untuk membuat atau membuka jaringan private melalui NAT atau firewall, lalu Menghubungkan localhost ke internet dengan tunnel yang aman. Atau bahasa gampangnya ngrok adalah layanan gratis yang memberikan kemampuan kepada aplikasi kita agar bisa diakses online.

// Midtrans adalah payment gateway yang memfasilitasi kebutuhan bisnis online dengan menyediakan layanan dalam berbagai metode pembayaran. Layanan ini memungkinkan pelaku industri beroperasi lebih mudah dan meningkatkan penjualan. Metode pembayaran yang disediakan adalah pembayaran kartu, transfer bank, debit langsung, e-wallet, over the counter, dan lain-lain.

// SECRET_KEY=bolehapaaja
// PATH_FILE=http://localhost:5000/uploads/
// SERVER_KEY=your_midtrans_server_key...
// CLIENT_KEY=your_midtrans_client_key
// EMAIL_SYSTEM=email_here...
// PASSWORD_SYSTEM=password_app...

// SNAP adalah portal pembayaran yang memungkinkan merchant menampilkan halaman pembayaran Midtrans langsung di website. Permintaan API harus dilakukan dari backend merchant untuk mendapatkan token transaksi Snap dengan memberikan informasi pembayaran dan Server Key. Setidaknya ada tiga komponen yang diperlukan untuk mendapatkan token Snap

// gomail
// Dengan memanfaatkan package net/smtp.
// Menggunakan Library gomail.
// C.18.1. Kirim Email Menggunakan net/smtp
// Golang menyediakan package net/smtp, isinya banyak API untuk berkomunikasi via protokol SMTP. Lewat package ini kita bisa melakukan operasi kirim email.

// Sebuah akun email diperlukan dalam mengirim email, silakan gunakan provider email apa saja. Pada chapter ini kita gunakan Google Mail (gmail), jadi siapkan satu buah akun gmail untuk keperluan testing

// const CONFIG_SMTP_HOST = "smtp.gmail.com"
// const CONFIG_SMTP_PORT = 587
// const CONFIG_SENDER_NAME = "PT. Makmur Subur Jaya <emailanda@gmail.com>"
// const CONFIG_AUTH_EMAIL = "emailanda@gmail.com"
// const CONFIG_AUTH_PASSWORD = "passwordemailanda"

// cara menjalankan ngrok
// ngrok --authtoken 2JWW4Ppm0nhPLYeDcyyGIAUdd65_6pYH13uBZUv4DFSpryA2j
// ngrok http 5000
