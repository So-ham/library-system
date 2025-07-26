package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type Book struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	Title       string    `json:"title" gorm:"not null" validate:"required"`
	Author      string    `json:"author" gorm:"not null" validate:"required"`
	ISBN        string    `json:"isbn" gorm:"unique;not null" validate:"required"`
	Publisher   string    `json:"publisher" gorm:"not null" validate:"required"`
	PublishDate time.Time `json:"publish_date" gorm:"not null" validate:"required"`
	Description string    `json:"description" gorm:"type:text"`
	Copies      int       `json:"copies" gorm:"not null;default:1" validate:"required,min=0"`
}

type BookRequest struct {
	Title       string    `json:"title" validate:"required"`
	Author      string    `json:"author" validate:"required"`
	ISBN        string    `json:"isbn" validate:"required"`
	Publisher   string    `json:"publisher" validate:"required"`
	PublishDate time.Time `json:"publish_date" validate:"required"`
	Description string    `json:"description"`
	Copies      int       `json:"copies" validate:"required,min=0"`
}

type BookResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	Publisher   string    `json:"publisher"`
	PublishDate time.Time `json:"publish_date"`
	Description string    `json:"description"`
	Copies      int       `json:"copies"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
