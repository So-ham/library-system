package services

import (
	"context"

	"library-system/internal/entities"

	"github.com/gofrs/uuid"
)

// CreateBook adds a new book to the library
func (s *service) CreateBook(ctx context.Context, req *entities.BookRequest) error {
	// Create book entity from request
	book := &entities.Book{
		Title:       req.Title,
		Author:      req.Author,
		ISBN:        req.ISBN,
		Publisher:   req.Publisher,
		PublishDate: req.PublishDate,
		Description: req.Description,
		Copies:      req.Copies,
	}

	// Save to database
	err := s.model.Book.Create(ctx, book)
	if err != nil {
		return err
	}

	return nil
}

// GetBookByID retrieves a book by its ID
func (s *service) GetBookByID(ctx context.Context, id uuid.UUID) (*entities.BookResponse, error) {
	// Get book from database
	book, err := s.model.Book.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entities.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		ISBN:        book.ISBN,
		Publisher:   book.Publisher,
		PublishDate: book.PublishDate,
		Description: book.Description,
		Copies:      book.Copies,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}, nil
}

// GetAllBooks retrieves all books from the library
func (s *service) GetAllBooks(ctx context.Context) ([]*entities.BookResponse, error) {
	// Get all books from database
	books, err := s.model.Book.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]*entities.BookResponse, len(books))
	for i, book := range books {
		resp[i] = &entities.BookResponse{
			ID:          book.ID,
			Title:       book.Title,
			Author:      book.Author,
			ISBN:        book.ISBN,
			Publisher:   book.Publisher,
			PublishDate: book.PublishDate,
			Description: book.Description,
			Copies:      book.Copies,
			CreatedAt:   book.CreatedAt,
			UpdatedAt:   book.UpdatedAt,
		}
	}

	return resp, nil
}

// UpdateBook modifies an existing book in the library
func (s *service) UpdateBook(ctx context.Context, id uuid.UUID, req *entities.BookRequest) error {
	// Check if book exists
	existingBook, err := s.model.Book.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existingBook.Title = req.Title
	existingBook.Author = req.Author
	existingBook.ISBN = req.ISBN
	existingBook.Publisher = req.Publisher
	existingBook.PublishDate = req.PublishDate
	existingBook.Description = req.Description
	existingBook.Copies = req.Copies

	err = s.model.Book.Update(ctx, existingBook)
	if err != nil {
		return err
	}

	// Return response
	return nil
}

// DeleteBook removes a book from the library
func (s *service) DeleteBook(ctx context.Context, id uuid.UUID) error {

	err := s.model.Book.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
