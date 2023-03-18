package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"project/dto"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Middleware untuk handle upload file user
func UpdateUserImage(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		file, _, err := r.FormFile("image")

		if err != nil && r.Method == "PATCH" {
			ctx := context.WithValue(r.Context(), "userImage", "false")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

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
		// filepath := data[8:] // split uploads(huruf paling 8 depan akan diambil)

		// filename akan ditambahkan kedalam variable ctx. dan r.Context akan di panggil jika ingin upload file
		ctx := context.WithValue(r.Context(), "userImage", data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// membuat middleware untuk menghandle upload file
func UploadTripImage(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handling dan parsing data dari form data yang ada data file nya. Argumen 1024 pada method tersebut adalah maxMemory sebesar 1024byte, apabila file yang diupload lebih besar maka akan disimpan di file sementara
		if err := r.ParseMultipartForm(1024); err != nil {
			panic(err.Error())
		}

		var arrImages []string

		files := r.MultipartForm.File["images"]
		for _, f := range files {
			// mengambil file dari form
			file, err := f.Open()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				response := dto.ErrorResult{
					Code:    http.StatusBadRequest,
					Message: "Please upload a JPG, JPEG or PNG image",
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			defer file.Close()

			// Apabila format file bukan .jpg, .jpeg atau .png, maka tampilkan error
			if filepath.Ext(f.Filename) != ".jpg" && filepath.Ext(f.Filename) != ".jpeg" && filepath.Ext(f.Filename) != ".png" {
				w.WriteHeader(http.StatusBadRequest)
				response := dto.ErrorResult{
					Code:    http.StatusBadRequest,
					Message: "The provided file format is not allowed. Please upload a JPG, JPEG or PNG image",
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			// create empty context
			var ctx = context.Background()

			// setup cloudinary credentials
			var CLOUD_NAME = os.Getenv("CLOUD_NAME")
			var API_KEY = os.Getenv("API_KEY")
			var API_SECRET = os.Getenv("API_SECRET")

			// create new instance of cloudinary object using cloudinary credentials
			cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

			// Upload file to Cloudinary
			resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "dewetour"})
			if err != nil {
				fmt.Println(err.Error())
			}
			// cek respon dari cloudinary
			// fmt.Println("respon from cloudinary", resp)

			// arrImages = append(arrImages, fileLocation)
			arrImages = append(arrImages, resp.SecureURL)
		}

		// membuat sebuah context baru dengan menyisipkan value di dalamnya, valuenya adalah array sring yang berisikan url img yang didapat dari cloudinary
		ctx := context.WithValue(r.Context(), "arrImages", arrImages)

		// mengirim nilai context ke object http.HandlerFunc yang menjadi parameter saat fungsi middleware ini dipanggil
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
