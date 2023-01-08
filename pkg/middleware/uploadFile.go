package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	dto "project/dto/result"
)

// function Upload file untuk upload file
func UploadFile(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		file, _, err := r.FormFile("image")

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("Error Retrieving the File")
			return
		}
		defer file.Close()

		const MAX_UPLOAD_SIZE = 10 << 20 // masksimal file upload 10mb

		// var MAX_UPLOAD_SIZE akan diparse
		r.ParseMultipartForm(MAX_UPLOAD_SIZE)

		// if contentLength lebih besar dari file yang diupload maka panggil ErrorResult
		if r.ContentLength > MAX_UPLOAD_SIZE {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Max size in 1mb"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// jika ukuran file sudah dibawah maksimal upload file maka file masuk ke folder upload
		tempFile, err := ioutil.TempFile("uploads", "image-*.png")
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close()

		// baca semua isi file yang kita upload, jika ada error maka tampilkan err
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		// write this byte array to our temporary file
		tempFile.Write(fileBytes)

		data := tempFile.Name()
		filename := data[8:] // split uploads(huruf paling 8 depan akan diambil)

		// filename akan ditambahkan kedalam variable ctx. dan r.Context akan di panggil jika ingin upload file
		ctx := context.WithValue(r.Context(), "dataFile", filename)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
