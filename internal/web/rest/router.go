package rest

import (
	"library-system/internal/handlers"

	"github.com/gorilla/mux"
)

// NewRouter returns a new router instance with configured routes
func NewRouter(h *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	// Book endpoints
	router.HandleFunc("/api/books", h.V1.GetAllBooks).Methods("GET")
	router.HandleFunc("/api/books", h.V1.CreateBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", h.V1.GetBookByID).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.V1.UpdateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", h.V1.DeleteBook).Methods("DELETE")

	return router
}
