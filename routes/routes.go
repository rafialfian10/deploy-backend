package routes

import (
	"github.com/gorilla/mux"
)

// membuat function RouteInit untuk membuat route ke masing-masing route
func RouteInit(r *mux.Router) {
	AuthRoutes(r)
	UserRoutes(r)
	CountryRoutes(r)
	TripRoutes(r)
	TransactionRoutes(r)
}
