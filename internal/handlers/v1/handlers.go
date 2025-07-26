package v1

import (
	"net/http"

	"library-system/internal/services"

	"github.com/go-playground/validator/v10"
)

type handlerV1 struct {
	Service  services.Service
	Validate *validator.Validate
}

type HandlerV1 interface {
	GetBookByID(w http.ResponseWriter, r *http.Request)
	GetAllBooks(w http.ResponseWriter, r *http.Request)
	CreateBook(w http.ResponseWriter, r *http.Request)
	UpdateBook(w http.ResponseWriter, r *http.Request)
	DeleteBook(w http.ResponseWriter, r *http.Request)
}

func New(s services.Service, v *validator.Validate) HandlerV1 {
	return &handlerV1{Service: s, Validate: v}
}
