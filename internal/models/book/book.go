package book

import (
	"context"
	"errors"
	"time"

	"library-system/internal/entities"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Book interface {
	Create(ctx context.Context, book *entities.Book) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Book, error)
	GetAll(ctx context.Context) ([]entities.Book, error)
	Update(ctx context.Context, book *entities.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type book struct {
	db *gorm.DB
}

func New(db *gorm.DB) Book {
	return &book{db: db}
}

func (b *book) Create(ctx context.Context, book *entities.Book) error {
	book.ID, _ = uuid.NewV4()
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	result := b.db.Create(book)
	return result.Error
}

func (b *book) GetByID(ctx context.Context, id uuid.UUID) (*entities.Book, error) {
	var book entities.Book
	result := b.db.Where("id = ?", id).First(&book)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, result.Error
	}

	return &book, nil
}

func (b *book) GetAll(ctx context.Context) ([]entities.Book, error) {
	var books []entities.Book
	result := b.db.Find(&books)

	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil
}

func (b *book) Update(ctx context.Context, book *entities.Book) error {
	book.UpdatedAt = time.Now()

	result := b.db.Save(book)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}

func (b *book) Delete(ctx context.Context, id uuid.UUID) error {
	result := b.db.Delete(&entities.Book{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}
