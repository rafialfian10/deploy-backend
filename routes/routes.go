package routes

import (
	"github.com/gorilla/mux"
)

// membuat function RouteInit untuk membuat route ke masing-masing route
func RouteInit(r *mux.Router) {
	UserRoutes(r)
	TripRoutes(r)
	AuthRoutes(r)
	CountryRoutes(r)
	TransactionRoutes(r)
}
