package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	dto "project/dto/result"
	jwtToken "project/pkg/jwt"
	"strings"
)

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// function Auth berfungsi untuk validasi token(user baru dapat melakukan CRUD setelah memasukkan token)
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("Authorization")

		// jika token kosong maka panggil ErrorResult
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "unauthorized"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// token akan displit dan diambil index ke 1 dan token akan dipanggil di DecodeToken
		token = strings.Split(token, " ")[1]
		claims, err := jwtToken.DecodeToken(token)

		// jika ada error maka panggil Result dan tampilkan pesan
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			response := Result{Code: http.StatusUnauthorized, Message: "unauthorized"}
			json.NewEncoder(w).Encode(response)
			return
		}

		//
		ctx := context.WithValue(r.Context(), "userInfo", claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
