package services

import (
	"context"

	"library-system/internal/entities"
	"library-system/internal/models"

	"github.com/gofrs/uuid"
)

// Service represents the service layer having
// all the services from all service packages
type service struct {
	model models.Model
}

// New creates a new instance of Service
func New(model *models.Model) Service {
	m := &service{model: *model}
	return m
}

type Service interface {

	// Book services
	CreateBook(ctx context.Context, req *entities.BookRequest) error
	GetBookByID(ctx context.Context, id uuid.UUID) (*entities.BookResponse, error)
	GetAllBooks(ctx context.Context) (resp []*entities.BookResponse, err error)
	UpdateBook(ctx context.Context, id uuid.UUID, req *entities.BookRequest) error
	DeleteBook(ctx context.Context, id uuid.UUID) error
}
