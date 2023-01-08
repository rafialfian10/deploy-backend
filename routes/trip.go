package routes

import (
	"project/handlers"
	"project/pkg/middleware"
	"project/pkg/mysql"
	"project/repositories"

	"github.com/gorilla/mux"
)

func TripRoutes(r *mux.Router) {
	// panggil repositoryTrip isikan parameter mywal.DB dan simpan ke dalam variabel
	TripRepository := repositories.RepositoriyTrip(mysql.DB)

	// panggil HandlerTrip dari (handlers/trip)
	h := handlers.HandlerTrip(TripRepository)

	r.HandleFunc("/trips", h.FindTrips).Methods("GET")
	r.HandleFunc("/trip/{id}", h.GetTrip).Methods("GET")
	r.HandleFunc("/trip", middleware.Auth(middleware.UploadFile(h.CreateTrip))).Methods("POST")
	r.HandleFunc("/trip/{id}", middleware.AuthAdmin(middleware.UploadFile(h.UpdateTrip))).Methods("PATCH")
	r.HandleFunc("/trip/{id}", middleware.AuthAdmin(h.DeleteTrip)).Methods("DELETE")
}
