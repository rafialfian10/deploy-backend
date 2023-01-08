package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	dto "project/dto/result"
	tripsdto "project/dto/trips"
	"project/models"
	"project/repositories"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var path_file_trip = "http://localhost:5000/uploads/"

// membuat struct handlerTrip untuk menghandle TripRepository. handlerTrip akan dipanggil ke setiap function
type handlerTrip struct {
	TripRepository repositories.TripRepository
}

func HandlerTrip(TripRepository repositories.TripRepository) *handlerTrip {
	return &handlerTrip{TripRepository}
}

// membuat struct function findTrips (all trip). parameter adalah struct handlerTrip
func (h *handlerTrip) FindTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Header berfungsi untuk menampilkan data.(text-html /json)

	// panggil function FindTrip didalam handlerTrip
	trips, err := h.TripRepository.FindTrips()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error()) // Error akan diEncode dan akan dikirim sebagai respon
	}

	// looping image pada trip, lalu trips akan di isi dengan data image dari struct
	for i, p := range trips {
		imagePath := os.Getenv("PATH_FILE") + p.Image
		trips[i].Image = imagePath
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trips}
	json.NewEncoder(w).Encode(response)
}

// membuat struct function GetTrip . parameter adalah struct handlerTrip
func (h *handlerTrip) GetTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// panggil function GetTrip didalam handlerTrip dengan index tertentu
	trip, err := h.TripRepository.GetTrip(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// jika tidak ada error maka image akan di isi dengan path image
	trip.Image = os.Getenv("PATH_FILE") + trip.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trip}
	json.NewEncoder(w).Encode(response) // response akan diEncode dan akan dikirim sebagai respon
}

// membuat struct function CreateTrip . parameter adalah struct handlerTrip
func (h *handlerTrip) CreateTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// middleware image
	dataContex := r.Context().Value("dataFile")
	filepath := dataContex.(string) // filename akan dipanggil di request

	// cloudinary
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	//parse data
	CountryId, _ := strconv.Atoi(r.FormValue("country_id"))
	day, _ := strconv.Atoi(r.FormValue("day"))
	night, _ := strconv.Atoi(r.FormValue("night"))
	price, _ := strconv.Atoi(r.FormValue("price"))
	quota, _ := strconv.Atoi(r.FormValue("quota"))

	// struct createTripRequest (dto) untuk menampung data
	request := tripsdto.CreateTripRequest{
		Title:          r.FormValue("title"),
		CountryId:      CountryId,
		Accomodation:   r.FormValue("accomodation"),
		Transportation: r.FormValue("transportation"),
		Eat:            r.FormValue("eat"),
		Day:            day,
		Night:          night,
		DateTrip:       r.FormValue("datetrip"),
		Price:          price,
		Quota:          quota,
		Description:    r.FormValue("description"),
	}

	// validasi request jika ada error maka panggil ErrorResult(jika ada request kosong maka error)
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "dewetour"})

	if err != nil {
		fmt.Println(err.Error())
	}

	// parse DateTrip menjadi string
	dateTrip, _ := time.Parse("2006-01-02", r.FormValue("datetrip"))

	// struct trip di isi dengan request
	trip := models.Trip{
		Title:          request.Title,
		CountryId:      request.CountryId,
		Accomodation:   request.Accomodation,
		Transportation: request.Transportation,
		Eat:            request.Eat,
		Day:            request.Day,
		Night:          request.Night,
		DateTrip:       dateTrip,
		Price:          request.Price,
		Quota:          request.Quota,
		Description:    request.Description,
		Image:          resp.SecureURL,
	}

	// panggil function CreateTrip didalam handlerTrip
	data, err := h.TripRepository.CreateTrip(trip)

	// jika tidak ada error maka panggil ErrorResult
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// panggil function getTrip agar setelah data di create data id akan keluar response
	tripResponse, err := h.TripRepository.GetTrip(data.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// jika  tidak ada error maka panggil SuccessResult
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseTrip(tripResponse)}
	json.NewEncoder(w).Encode(response)
}

// membuat struct function UpdateTrip . parameter adalah struct handlerTrip
func (h *handlerTrip) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// panggil function GetTrip didalam handlerTrip dengan index tertentu
	trip, err := h.TripRepository.GetTrip(int(id))

	// jika ada error maka panggil ErrorResult
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// middleware
	dataContex := r.Context().Value("dataFile")
	filename := dataContex.(string)

	// request image agar nantinya image dapat diupdate
	request := tripsdto.UpdateTripRequest{
		Image: filename,
	}

	// title
	if r.FormValue("title") != "" {
		trip.Title = r.FormValue("title")
	}

	// country id
	countryId, _ := strconv.Atoi(r.FormValue("country_id"))
	if countryId != 0 {
		trip.CountryId = countryId
	}

	// accomodation
	if r.FormValue("accomodation") != "" {
		trip.Accomodation = r.FormValue("accomodation")
	}

	// transportation
	if r.FormValue("transportation") != "" {
		trip.Transportation = r.FormValue("transportation")
	}

	// eat
	if r.FormValue("eat") != "" {
		trip.Eat = r.FormValue("eat")
	}

	// parse day
	day, _ := strconv.Atoi(r.FormValue("day"))
	if day != 0 {
		trip.Day = day
	}

	// parse night
	night, _ := strconv.Atoi(r.FormValue("night"))
	if night != 0 {
		trip.Night = night
	}

	// parse time
	date, _ := time.Parse("2006-01-02", r.FormValue("datetrip"))
	time := time.Now()
	if date != time {
		trip.DateTrip = date
	}

	// parse price
	price, _ := strconv.Atoi(r.FormValue("price"))
	if price != 0 {
		trip.Price = price
	}

	// parse quota
	quota, _ := strconv.Atoi(r.FormValue("quota"))
	if quota != 0 {
		trip.Quota = quota
	}

	// description
	if r.FormValue("description") != "" {
		trip.Description = r.FormValue("description")
	}

	// image
	if request.Image != "" {
		trip.Image = request.Image
	}

	// panggil function UpdateTrip didalam handlerTrip untuk update semua data trip lalu tampung ke var new trip
	newTrip, err := h.TripRepository.UpdateTrip(trip)

	// jika ada error maka tampilkan ErrorResult
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// panggil function getTrip agar setelah data di create data id akan keluar response
	newtripResponse, err := h.TripRepository.GetTrip(newTrip.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// jika tidak ada error maka SuccessResult
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseTrip(newtripResponse)}
	json.NewEncoder(w).Encode(response)
}

// function delete trip
func (h *handlerTrip) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// panggil function GetTrip didalam handlerTrip dengan index tertentu
	trip, err := h.TripRepository.GetTrip(id)

	// jika ada error panggil Errorresult
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// panggil function DeleteTrip berdasarkan id
	data, err := h.TripRepository.DeleteTrip(trip)

	// jika ada error maka tampilkan errorResult
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// jika tidak ada error maka
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseTrip(data)}
	json.NewEncoder(w).Encode(response)
}

// function convert response trip
func convertResponseTrip(u models.Trip) tripsdto.TripResponse {
	return tripsdto.TripResponse{
		Id:             u.Id,
		Title:          u.Title,
		CountryId:      u.CountryId,
		Accomodation:   u.Accomodation,
		Transportation: u.Transportation,
		Eat:            u.Eat,
		Day:            u.Day,
		Night:          u.Night,
		DateTrip:       u.DateTrip.Format("2 January 2006"),
		Price:          u.Price,
		Quota:          u.Quota,
		Description:    u.Description,
		Image:          u.Image,
	}
}
