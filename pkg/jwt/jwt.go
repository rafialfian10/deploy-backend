package jwtToken

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

var SecretKey = "SECRET_KEY"

// function GenerateToken untuk membuat token
func GenerateToken(claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	webtoken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return webtoken, nil
}

// function verify token untuk verifikasi apakah token yang kita buat sama dengan token yang dimasukkan
func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

// function DecodeToken
func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, isOk := token.Claims.(jwt.MapClaims)
	if isOk && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// function DecodeToken berfungsi ketika request masuk, middleware akan mengecek apakah ada auth?, jika ada maka token akan diambil lalu dikirim ke fungsi decodeToken didalam decode token. lalu token akan diperiksa menggunakan function verifyToken, apabila token valid maka function decodeToken akan mengambil data yang disisipkan kedalam token
