package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	dto "project/dto/result"
	transactionsdto "project/dto/transactions.go"
	"project/models"
	"project/repositories"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

var path_file_trans = "http://localhost:5000/uploads/"

var c = coreapi.Client{
	ServerKey: os.Getenv("SERVER_KEY"),
	ClientKey: os.Getenv("CLIENT_KEY"),
}

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) FindTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transaction, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	for i, p := range transaction {
		transaction[i].Image = path_file_trans + p.Image
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: transaction}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetAllTransactionByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := r.Context().Value("userInfo").(jwt.MapClaims)
	id := int(claims["id"].(float64))

	// mengambil seluruh data transaction
	transaction, err := h.TransactionRepository.FindTransactionsByUser(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// menyiapkan response
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertMultipleTransactionResponse(transaction),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	trans, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	trans.Image = path_file_trans + trans.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trans}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get data user token
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	// // middleware upload file
	// dataContex := r.Context().Value("dataFile")
	// filename := dataContex.(string)

	//parse data
	counterQty, _ := strconv.Atoi(r.FormValue("qty"))
	total, _ := strconv.Atoi(r.FormValue("total"))
	tripId, _ := strconv.Atoi(r.FormValue("trip_id"))

	request := transactionsdto.CreateTransactionRequest{
		CounterQty: counterQty,
		Total:      total,
		Status:     r.FormValue("status"),
		// Image:      filename,
		TripId: tripId,
		UserId: userId,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// buat transaction id unik
	var transIdIsMatch = false
	var transactionId int

	// buat transaksi id dengan time unix
	for !transIdIsMatch {
		transactionId = int(time.Now().Unix())
		transaction, _ := h.TransactionRepository.GetTransaction(transactionId)
		if transaction.Id == 0 {
			transIdIsMatch = true
		}
	}

	// request ke model transaction
	transaction := models.Transaction{
		Id:          transactionId,
		CounterQty:  request.CounterQty,
		Total:       request.Total,
		Status:      request.Status,
		BookingDate: timeIn("Asia/Jakarta"),
		// Image:      request.Image,
		TripId: request.TripId,
		UserId: request.UserId,
	}

	// panggil Transaction repository dan masukkan transaction ke dalam kedalam function CreateTransaction
	data, err := h.TransactionRepository.CreateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	// panggil function getTrip agar setelah data di create data id akan keluar response
	transactionResponse, err := h.TransactionRepository.GetTransaction(data.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// snap
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transactionResponse.Id),
			GrossAmt: int64(transactionResponse.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: transactionResponse.User.Name,
			Email: transactionResponse.User.Email,
			BillAddr: &midtrans.CustomerAddress{
				FName:   transactionResponse.User.Name,
				Phone:   transactionResponse.User.Phone,
				Address: transactionResponse.User.Address,
				// Postcode: transaction.User.PostCode,
			},
			ShipAddr: &midtrans.CustomerAddress{
				FName:   transactionResponse.User.Name,
				Phone:   transactionResponse.User.Phone,
				Address: transactionResponse.User.Address,
				// Postcode: transaction.User.PostCode,
			},
		},
	}

	snapResp, _ := s.CreateTransaction(req)
	fmt.Println(snapResp)

	// mengupdate token di database
	transaction, _ = h.TransactionRepository.UpdateTokenTransaction(snapResp.Token, transactionResponse.Id)

	// mengambil data transaction yang baru diupdate
	transactionUpdated, _ := h.TransactionRepository.GetTransaction(transaction.Id)

	// menyiapkan response
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOneTransactionResponse(transactionUpdated),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id_transaction"])

	// mengambil data transaction yang baru ditambahkan
	transaction, _ := h.TransactionRepository.GetTransaction(id)
	transactionId := strconv.Itoa(transaction.Id)

	var s = snap.Client{}
	s.New(os.Getenv("SB-Mid-server-CBYg0a0CWSxQrUrIYbcaHJvM"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionId,
			GrossAmt: int64(transaction.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: transaction.User.Name,
			Email: transaction.User.Email,
		},
	}

	snapResp, _ := s.CreateTransaction(req)

	// mengupdate token di database
	transaction, _ = h.TransactionRepository.UpdateTokenTransaction(snapResp.Token, id)

	// mengambil data transaction yang baru diupdate
	transactionUpdated, _ := h.TransactionRepository.GetTransaction(id)

	// menyiapkan response
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertOneTransactionResponse(transactionUpdated),
	}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

// function send email
func SendEmail(status string, transaction models.Transaction) {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = "dewetour <rafialfian770@gmail.com>"
	var CONFIG_AUTH_EMAIL = os.Getenv("SYSTEM_EMAIL")
	var CONFIG_AUTH_PASSWORD = os.Getenv("SYSTEM_PASSWORD")

	var tripName = transaction.User.Name
	var price = strconv.Itoa(transaction.Total)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", transaction.User.Email)
	mailer.SetHeader("Subject", "Status Transaction")
	mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
    <html lang="en">
      <head>
      <meta charset="UTF-8" />
      <meta http-equiv="X-UA-Compatible" content="IE=edge" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <title>Document</title>
      <style>
        h1 {
        color: brown;
        }
      </style>
      </head>
      <body>
      <h2>Product payment :</h2>
      <ul style="list-style-type:none;">
        <li>Name : %s</li>
        <li>Total payment: Rp.%s</li>
        <li>Status : %s</li>
		<li>Iklan : %s</li>
      </ul>
      </body>
    </html>`, tripName, price, status, "Terima kasih"))

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// function notification (mengixinkan mitrans untuk mengupdate status transaksi)
func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// transaksi status
	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)

	orderIdConvert, _ := strconv.Atoi(orderId)

	// panggil function get transaction
	transaction, _ := h.TransactionRepository.GetTransaction(orderIdConvert)
	fmt.Println(transactionStatus, fraudStatus, orderId, transaction)

	// kondisi transaksi
	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			SendEmail("Transaction Failed", transaction)
			h.TransactionRepository.UpdateTransaction("pending", transaction.Id)
		} else if fraudStatus == "accept" {
			SendEmail("Transaction Success", transaction)
			h.TransactionRepository.UpdateTransaction("success", transaction.Id)
		}
	} else if transactionStatus == "settlement" {
		SendEmail("Transaction Success", transaction)
		h.TransactionRepository.UpdateTransaction("success", transaction.Id)
	} else if transactionStatus == "deny" {
		SendEmail("Transaction Failed", transaction)
		h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		SendEmail("Transaction Failed", transaction)
		h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
	} else if transactionStatus == "pending" {
		SendEmail("Transaction Pending", transaction)
		h.TransactionRepository.UpdateTransaction("pending", transaction.Id)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handlerTransaction) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.TransactionRepository.DeleteTransaction(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseTransaction(data)}
	json.NewEncoder(w).Encode(response)
}

// membuat fungsi konversi data yang akan disajikan sebagai response sesuai requirement
func convertResponseTransaction(u models.Transaction) transactionsdto.TransactionResponse {
	return transactionsdto.TransactionResponse{
		Id:         u.Id,
		CounterQty: u.CounterQty,
		Total:      u.Total,
		Status:     u.Status,
		Token:      u.Token,
		// Image:      u.Image,
		Trip: u.Trip,
		User: u.User,
	}
}

// membuat fungsi konversi data yang akan disajikan sebagai response sesuai requirement
func convertOneTransactionResponse(t models.Transaction) transactionsdto.TransactionResponse {
	result := transactionsdto.TransactionResponse{
		Id:         t.Id,
		CounterQty: t.CounterQty,
		Total:      t.Total,
		Status:     t.Status,
		Token:      t.Token,
		User:       t.User,
		Trip: models.TripResponse{
			Id:             t.Trip.Id,
			Title:          t.Trip.Title,
			Country:        t.Trip.Country,
			Accomodation:   t.Trip.Accomodation,
			Transportation: t.Trip.Transportation,
			Eat:            t.Trip.Eat,
			Day:            t.Trip.Day,
			Night:          t.Trip.Night,
			Price:          t.Trip.Price,
			Quota:          t.Trip.Quota,
			Description:    t.Trip.Description,
		},
	}
	result.BookingDate = t.BookingDate.Format("Monday, 2 January 2006")
	result.Trip.DateTrip = t.Trip.DateTrip
	result.Trip.Image = t.Trip.Image
	// for _, img := range t.Trip.Image {
	// 	result.Trip.Images = append(result.Trip.Images, img.FileName)
	// }

	return result
}

// membuat fungsi konversi data yang akan disajikan sebagai response sesuai requirement
func convertMultipleTransactionResponse(t []models.Transaction) []transactionsdto.TransactionResponse {
	var result []transactionsdto.TransactionResponse

	for _, t := range t {
		transaction := transactionsdto.TransactionResponse{
			Id:         t.Id,
			CounterQty: t.CounterQty,
			Total:      t.Total,
			Status:     t.Status,
			Token:      t.Token,
			User:       t.User,
			Trip: models.TripResponse{
				Id:             t.Trip.Id,
				Title:          t.Trip.Title,
				Country:        t.Trip.Country,
				Accomodation:   t.Trip.Accomodation,
				Transportation: t.Trip.Transportation,
				Eat:            t.Trip.Eat,
				Day:            t.Trip.Day,
				Night:          t.Trip.Night,
				Price:          t.Trip.Price,
				Quota:          t.Trip.Quota,
				Description:    t.Trip.Description,
			},
		}
		transaction.BookingDate = t.BookingDate.Format("Monday, 2 January 2006")
		transaction.Trip.DateTrip = t.Trip.DateTrip
		// for _, img := range t.Trip.Image {
		// 	transaction.Trip.Image = append(transaction.Trip.image, img.FileName)
		// }
		result = append(result, transaction)
	}
	return result
}

// fungsi untuk mendapatkan waktu sesuai zona indonesia
func timeIn(name string) time.Time {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}
