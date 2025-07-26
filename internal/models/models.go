package models

import (
	"library-system/internal/models/book"

	"gorm.io/gorm"
)

type Model struct {
	Book book.Book
}

// New creates a new instance of Model
func New(gdb *gorm.DB) *Model {
	return &Model{
		Book: book.New(gdb),
	}
}
