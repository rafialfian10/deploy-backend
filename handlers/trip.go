package handlers

import (
	"encoding/json"
	"net/http"
	dto "project/dto"
	"project/models"
	"project/repositories"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

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

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertMultipleTripResponse(trips, r)}
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

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOneTripResponse(trip, r)}
	json.NewEncoder(w).Encode(response) // response akan diEncode dan akan dikirim sebagai respon
}

// membuat struct function CreateTrip . parameter adalah struct handlerTrip
func (h *handlerTrip) CreateTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//parse data
	CountryId, _ := strconv.Atoi(r.FormValue("country_id"))
	day, _ := strconv.Atoi(r.FormValue("day"))
	night, _ := strconv.Atoi(r.FormValue("night"))
	price, _ := strconv.Atoi(r.FormValue("price"))
	quota, _ := strconv.Atoi(r.FormValue("quota"))

	// struct createTripRequest (dto) untuk menampung data
	request := dto.CreateTripRequest{
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
	// middleware image
	dataContex := r.Context().Value("arrImages").([]string)
	request.Images = append(request.Images, dataContex...)

	// validasi request jika ada error maka panggil ErrorResult(jika ada request kosong maka error)
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
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
	}

	// mengisikan array image file ke array image milik object trip, image yang ditambahkan disini nantinya akan otomatis ditambahkan pula ke tabel Image yang berelasi (sesuai dengan association method yang ada di doc gorm)
	for _, image := range request.Images {
		imgData := models.ImageResponse{
			Name: image,
		}
		trip.Image = append(trip.Image, imgData)
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
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOneTripResponse(tripResponse, r)}
	json.NewEncoder(w).Encode(response)
}

// membuat struct function UpdateTrip . parameter adalah struct handlerTrip
func (h *handlerTrip) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// mengambil data dari request body
	var updateTrip dto.UpdateTripRequest
	updateTrip.Title = r.FormValue("title")
	updateTrip.CountryId, _ = strconv.Atoi(r.FormValue("country_id"))
	updateTrip.Accomodation = r.FormValue("accomodation")
	updateTrip.Transportation = r.FormValue("transportation")
	updateTrip.Eat = r.FormValue("eat")
	updateTrip.Day, _ = strconv.Atoi(r.FormValue("day"))
	updateTrip.Night, _ = strconv.Atoi(r.FormValue("night"))
	updateTrip.DateTrip = r.FormValue("datetrip")
	updateTrip.Price, _ = strconv.Atoi(r.FormValue("price"))
	updateTrip.Quota, _ = strconv.Atoi(r.FormValue("quota"))
	updateTrip.Description = r.FormValue("description")

	// middleware
	dataContex := r.Context().Value("arrImages").([]string)
	updateTrip.Images = append(updateTrip.Images, dataContex...)

	// panggil function GetTrip didalam handlerTrip dengan index tertentu
	trip, err := h.TripRepository.GetTrip(int(id))

	// jika ada error maka panggil ErrorResult
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
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
	// mereplace gambar jika ada gambar yang diuplad
	arrReqImagesLength := len(updateTrip.Images)
	if arrReqImagesLength > 0 {
		// mengambil panjang array image dari data yang akan diupdate
		arrPrevImagesLength := len(trip.Image)
		// fmt.Println(arrPrevImagesLength)

		// mengupdate array image file ke array image milik object tripWantToUpdate
		for i, image := range updateTrip.Images {
			// mengganti gambar sebelumnya (sesuai jumlah gambar)
			if i < arrPrevImagesLength {
				trip.Image[i].Name = string(image)
			} else { // jika gambar yang diupload lebih banyak dari gambar sebelumnya, maka tambahkan gambar baru
				imgData := models.ImageResponse{
					Name: string(image),
				}
				trip.Image = append(trip.Image, imgData)
			}
		}

		// menghapus gambar lama jika gambar yang baru yang direquest lebih sedikit
		for i := range trip.Image {
			if i >= arrReqImagesLength {
				trip.Image[i].Name = "unused"
			}
		}
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
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOneTripResponse(newtripResponse, r)}
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
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertOneTripResponse(data, r)}
	json.NewEncoder(w).Encode(response)
}

// function convert response trip
func convertOneTripResponse(t models.Trip, r *http.Request) dto.TripResponse {
	var result dto.TripResponse
	result.Id = t.Id
	result.Title = t.Title
	result.Country = t.Country
	result.Accomodation = t.Accomodation
	result.Transportation = t.Transportation
	result.Eat = t.Eat
	result.Day = t.Day
	result.Night = t.Night
	result.DateTrip = t.DateTrip.Format("02 January 2006")
	result.Price = t.Price
	result.Quota = t.Quota
	result.Description = t.Description

	for _, img := range t.Image {
		result.Images = append(result.Images, img.Name)
	}

	// fmt.Println(result.Images)
	return result
}

// convert multiple trip
func convertMultipleTripResponse(t []models.Trip, r *http.Request) []dto.TripResponse {
	var results []dto.TripResponse

	for _, trip := range t {
		var t dto.TripResponse

		t.Id = trip.Id
		t.Title = trip.Title
		t.Country = trip.Country
		t.Accomodation = trip.Accomodation
		t.Transportation = trip.Transportation
		t.Eat = trip.Eat
		t.Day = trip.Day
		t.Night = trip.Night
		t.DateTrip = trip.DateTrip.Format("02 January 2006")
		t.Price = trip.Price
		t.Quota = trip.Quota
		t.Description = trip.Description

		for _, img := range trip.Image {
			t.Images = append(t.Images, img.Name)
		}

		results = append(results, t)
	}

	return results
}
