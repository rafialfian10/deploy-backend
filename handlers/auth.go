package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	authdto "project/dto/auth"
	dto "project/dto/result"
	"project/models"
	"project/pkg/bcrypt"
	jwtToken "project/pkg/jwt"
	"project/repositories"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

// Handle struct
type handlerAuth struct {
	AuthRepository repositories.AuthRepository
}

// Handle Authentication
func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

// membuat struct function Register
func (h *handlerAuth) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// panggil method new, dan dto RegisterRequest akan digunakan sebagai parameter
	request := new(authdto.RegisterRequest)

	// err akan decode menjadi data aslinya dan akan di request di body, dan jika ada error maka panggil ErrorResult lalu encode response
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// lakukan validasi
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Hashing password request.Password(registerRequest) dengan method HashingPassword
	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// jika tidak ada error struct user akan diisi dengan request
	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: password,
		Gender:   request.Gender,
		Phone:    request.Phone,
		Address:  request.Address,
		// Image:    request.Image,
		Role: "user",
	}

	// panggil Register lalu user akan digunakan sebagai parameter
	userData, err := h.AuthRepository.Register(user)

	// jika ada error maka panggil ErrorResult
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	// jika tidak ada error maka panggil SuccesResult dan data akan di isi dengan func convertresponseRegister
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseRegister(userData)}
	json.NewEncoder(w).Encode(response)
}

// function convertResponseRegister
func convertResponseRegister(u models.User) authdto.RegisterResponse {
	return authdto.RegisterResponse{
		Email:    u.Email,
		Password: u.Password,
	}
}

// membuat struct function Login
func (h *handlerAuth) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// panggil method new, dan dto LoginRequest akan digunakan sebagai parameter
	request := new(authdto.LoginRequest)

	// err akan decode menjadi data aslinya dan akan di request di body, dan jika ada error maka panggil ErrorResult lalu encode response
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// jika tidak ada error struct user akan diisi dengan request
	user := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	// panggil Register lalu user.Email akan digunakan sebagai parameter
	user, err := h.AuthRepository.Login(user.Email)

	// jika ada error maka panggil ErrorResult
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// check password dengan method CheckPasswordHash. par request.Password dan user.Password akan di cek
	isValid := bcrypt.CheckPasswordHash(request.Password, user.Password)

	// jika tidak valid maka panggil ErrorResult
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "salah email atau password"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// membuat data yang akan disimpan di jwt dan claim akan digunakan untuk generate token
	claims := jwt.MapClaims{}

	claims["id"] = user.Id // buat key id valuenya user.Id
	claims["role"] = user.Role
	claims["email"] = user.Email
	claims["password"] = user.Password
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // mak token 2 jam

	// panggil method GenerateToken(agar dibuatkan token) dan claim akan dijadikan parameter
	token, err := jwtToken.GenerateToken(&claims)

	// jika ada error / tidak ada token maka err
	if err != nil {
		log.Println(err)
		fmt.Println("Unauthorize")
		return
	}

	// jika tidak ada error struct LoginResponse akan di isi data request user
	loginResponse := authdto.LoginResponse{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Token:    token,
		Role:     user.Role,
	}

	// dan login loginResponse akan dijadikan value dari data
	w.Header().Set("Content-Type", "application/json")
	response := dto.SuccessResult{Code: http.StatusOK, Data: loginResponse}
	json.NewEncoder(w).Encode(response)
}

// membuat fungsi yang digunakan untuk mendaftarkan user baru
func (h *handlerAuth) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// mengambil data user dari request body
	var request authdto.RegisterRequest
	json.NewDecoder(r.Body).Decode(&request)

	// memvalidasi inputan dari request body berdasarkan struct dto.CountryRequest
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// menghashing password
	hashedPassword, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// membuat object user baru dengan cetakan models.User
	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
		Gender:   request.Gender,
		Phone:    request.Phone,
		Address:  request.Address,
		// Image:    request.Image,
		Role: "admin",
	}

	// mengirim data user baru ke database
	adminData, err := h.AuthRepository.Register(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// menyiapkan response
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{
		Code: http.StatusOK,
		Data: convertResponseRegister(adminData),
	}
	// mengirim response
	json.NewEncoder(w).Encode(response)
}

func (h *handlerAuth) CheckAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	// Check User by Id
	user, err := h.AuthRepository.Getuser(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	CheckAuthResponse := authdto.CheckAuth{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	response := dto.SuccessResult{Code: http.StatusOK, Data: CheckAuthResponse}
	json.NewEncoder(w).Encode(response)
}
