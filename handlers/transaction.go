package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	dto "project/dto"
	"project/models"
	"project/repositories"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// var c = coreapi.Client{
// 	ServerKey: os.Getenv("SERVER_KEY"),
// 	ClientKey: os.Getenv("CLIENT_KEY"),
// }

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

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertMultipleTransactionResponse(transaction)}

	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetAllTransactionByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := r.Context().Value("userInfo").(jwt.MapClaims)
	id := int(claims["id"].(float64))

	transaction, err := h.TransactionRepository.FindTransactionsByUser(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertMultipleTransactionResponse(transaction)}

	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOneTransactionResponse(transaction)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))
	counterqty, _ := strconv.Atoi(r.FormValue("counter_qty"))
	total, _ := strconv.Atoi(r.FormValue("total"))
	tripId, _ := strconv.Atoi(r.FormValue("tripId"))
	request := dto.CreateTransactionRequest{
		CounterQty: counterqty,
		Total:      total,
		TripID:     tripId,
		UserId:     userId,
	}

	json.NewDecoder(r.Body).Decode(&request)

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var TrxIdMatch = false
	var TrxId int
	for !TrxIdMatch {
		TrxId = int(time.Now().Unix())
		transactionData, _ := h.TransactionRepository.GetTransaction(TrxId)
		if transactionData.Id == 0 {
			TrxIdMatch = true
		}
	}

	newTransaction := models.Transaction{
		Id:          TrxId,
		CounterQty:  request.CounterQty,
		Total:       request.Total,
		BookingDate: time.Now().UTC(),
		Status:      "pending",
		TripID:      request.TripID,
		UserId:      request.UserId,
	}

	transaction, err := h.TransactionRepository.CreateTransaction(newTransaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	TransactionAdded, _ := h.TransactionRepository.GetTransaction(transaction.Id)

	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(TransactionAdded.Id),
			GrossAmt: int64(TransactionAdded.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: TransactionAdded.User.Name,
			Email: TransactionAdded.User.Email,
		},
	}

	snapResp, _ := s.CreateTransaction(req)
	updateTransaction, _ := h.TransactionRepository.UpdateTokenTransaction(snapResp.Token, TransactionAdded.Id)
	transactionUpdated, _ := h.TransactionRepository.GetTransaction(updateTransaction.Id)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOneTransactionResponse(transactionUpdated)}

	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.Id),
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
	transaction, _ = h.TransactionRepository.UpdateTokenTransaction(snapResp.Token, id)
	transactionUpdated, _ := h.TransactionRepository.GetTransaction(id)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOneTransactionResponse(transactionUpdated)}

	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) UpdateTransactionByAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	status := r.FormValue("status")
	request := dto.UpdateTransactionRequest{
		Status: status,
	}

	json.NewDecoder(r.Body).Decode(&request)

	// fmt.Println("ID before", id)

	// mengambil data yang ingin diupdate berdasarkan id yang didapatkan dari url
	_, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusNotFound, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// fmt.Println("ID after", id)

	// mengirim data transaction yang sudah diupdate ke database
	transactionUpdated, err := h.TransactionRepository.UpdateTransaction(request.Status, id)
	// fmt.Println("Transaction Updated", transactionUpdated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// mengambil detail transaction yang baru saja ditambahkan (perlu diambil ulang, karena hasil dari transactionAdded hanya ada country_id saja, tanpa ada detail country nya)
	getTransactionUpdated, err := h.TransactionRepository.GetTransaction(transactionUpdated.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// menyiapkan response
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOneTransactionResponse(getTransactionUpdated)}

	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func SendEmail(status string, transaction models.Transaction) {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = "dewetour <rafialfian770@gmail.com>"
	var CONFIG_AUTH_EMAIL = os.Getenv("SYSTEM_EMAIL")
	var CONFIG_AUTH_PASSWORD = os.Getenv("SYSTEM_PASSWORD")

	var tripName = transaction.Trip.Title
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

func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)
	order_id, _ := strconv.Atoi(orderId)

	transaction, err := h.TransactionRepository.GetTransaction(order_id)
	if err != nil {
		fmt.Println("Transaction not found")
		return
	}
	// fmt.Println(transactionStatus, fraudStatus, orderId, transaction)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			h.TransactionRepository.UpdateTransaction("pending", transaction.Id)
		} else if fraudStatus == "accept" {
			SendEmail("Transaction Success", transaction)
			transaction.Status = "success"
			h.TransactionRepository.UpdateTransaction("success", transaction.Id)
		}
	} else if transactionStatus == "settlement" {
		SendEmail("Transaction Success", transaction)
		transaction.Status = "success"
		h.TransactionRepository.UpdateTransaction("success", transaction.Id)
	} else if transactionStatus == "deny" {
		SendEmail("Transaction Failed", transaction)
		transaction.Status = "failed"
		h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		SendEmail("Transaction Failed", transaction)
		transaction.Status = "failed"
		h.TransactionRepository.UpdateTransaction("failed", transaction.Id)
	} else if transactionStatus == "pending" {
		SendEmail("Transaction Pending", transaction)
		transaction.Status = "pending"
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
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOneTransactionResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func convertOneTransactionResponse(t models.Transaction) dto.TransactionResponse {
	result := dto.TransactionResponse{
		Id:         t.Id,
		CounterQty: t.CounterQty,
		Total:      t.Total,
		Status:     t.Status,
		Token:      t.Token,
		User:       t.User,
		Trip: dto.TripResponse{
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
	result.Trip.DateTrip = t.Trip.DateTrip.Format("Monday, 2 January 2006")

	for _, img := range t.Trip.Image {
		result.Trip.Images = append(result.Trip.Images, img.Name)
	}

	return result
}

func convertMultipleTransactionResponse(t []models.Transaction) []dto.TransactionResponse {
	var result []dto.TransactionResponse

	for _, trans := range t {
		transaction := dto.TransactionResponse{
			Id:         trans.Id,
			CounterQty: trans.CounterQty,
			Total:      trans.Total,
			Status:     trans.Status,
			Token:      trans.Token,
			User:       trans.User,
			Trip: dto.TripResponse{
				Id:             trans.Trip.Id,
				Title:          trans.Trip.Title,
				Country:        trans.Trip.Country,
				Accomodation:   trans.Trip.Accomodation,
				Transportation: trans.Trip.Transportation,
				Eat:            trans.Trip.Eat,
				Day:            trans.Trip.Day,
				Night:          trans.Trip.Night,
				Price:          trans.Trip.Price,
				Quota:          trans.Trip.Quota,
				Description:    trans.Trip.Description,
			},
		}
		transaction.BookingDate = trans.BookingDate.Format("Monday, 2 January 2006")
		transaction.Trip.DateTrip = trans.Trip.DateTrip.Format("Monday, 2 January 2006")

		for _, img := range trans.Trip.Image {
			transaction.Trip.Images = append(transaction.Trip.Images, img.Name)
		}

		result = append(result, transaction)
	}

	return result
}

// fungsi untuk mendapatkan waktu sesuai zona indonesia
// func timeIn(name string) time.Time {
// 	loc, err := time.LoadLocation(name)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return time.Now().In(loc)
// }
