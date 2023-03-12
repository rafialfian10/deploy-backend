package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	dto "project/dto"
	jwtToken "project/pkg/jwt"
	"strings"
)

type ResultAdmin struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func AuthAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// mengambil token dari request header
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "unauthorized"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// memvalidasi token dan mengambil nilai claim jika token tersebut valid
		token = strings.Split(token, " ")[1]
		claims, err := jwtToken.DecodeToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "unauthorized"}
			json.NewEncoder(w).Encode(response)
			return
		}

		role := claims["role"].(string)
		if role != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "unauthorized, you're not admin !"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// menyiapkan context dengan key "userInfo" yang berisi jwt claim
		ctx := context.WithValue(r.Context(), "userInfo", claims)

		// mengirim nilai context ke object http.HandlerFunc yang menjadi parameter saat fungsi middleware ini dipanggil
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
