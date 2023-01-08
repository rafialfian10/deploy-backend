package routes

import (
	"project/handlers"
	"project/pkg/middleware"
	"project/pkg/mysql"
	"project/repositories"

	"github.com/gorilla/mux"
)

func CountryRoutes(r *mux.Router) {
	CountryRepository := repositories.RepositoryCountry(mysql.DB)
	h := handlers.HandlerCountry(CountryRepository)

	r.HandleFunc("/countries", h.FindCountries).Methods("GET")
	r.HandleFunc("/country/{id}", h.GetCountry).Methods("GET")
	r.HandleFunc("/country", middleware.AuthAdmin(h.CreateCountry)).Methods("POST")
	r.HandleFunc("/country/{id}", middleware.AuthAdmin(h.UpdateCountry)).Methods("PATCH")
	r.HandleFunc("/country/{id}", middleware.AuthAdmin(h.DeleteCountry)).Methods("DELETE")
}
