package services

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"library-system/internal/entities"
	"library-system/internal/models"
	bookMock "library-system/internal/models/book/mocks"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

func Test_service_CreateBook(t *testing.T) {
	bookID, _ := uuid.NewV4()
	testTime := time.Now()

	req := &entities.BookRequest{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "1234567890",
		Publisher:   "Test Publisher",
		PublishDate: testTime,
		Description: "Test Description",
		Copies:      5,
	}

	successMock := bookMock.Book{}
	successMock.On("Create", mock.Anything, mock.MatchedBy(func(b *entities.Book) bool {
		return b.Title == req.Title &&
			b.Author == req.Author &&
			b.ISBN == req.ISBN &&
			b.Publisher == req.Publisher &&
			b.PublishDate.Equal(req.PublishDate) &&
			b.Description == req.Description &&
			b.Copies == req.Copies
	})).Return(nil).Run(func(args mock.Arguments) {
		book := args.Get(1).(*entities.Book)
		book.ID = bookID
		book.CreatedAt = testTime
		book.UpdatedAt = testTime
	})

	errorMock := bookMock.Book{}
	errorMock.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))

	tests := []struct {
		name    string
		s       *service
		req     *entities.BookRequest
		wantErr bool
	}{
		{
			name: "successful creation",
			s:    &service{model: models.Model{Book: &successMock}},
			req:  req, wantErr: false,
		},
		{
			name: "database error",
			s:    &service{model: models.Model{Book: &errorMock}},
			req:  req, wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.CreateBook(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBook() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.s.model.Book.(*bookMock.Book).AssertExpectations(t)
		})
	}
}

func Test_service_GetBookByID(t *testing.T) {
	bookID, _ := uuid.NewV4()
	invalidID, _ := uuid.NewV4()
	testTime := time.Now()

	book := &entities.Book{
		ID:          bookID,
		CreatedAt:   testTime,
		UpdatedAt:   testTime,
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "1234567890",
		Publisher:   "Test Publisher",
		PublishDate: testTime,
		Description: "Test Description",
		Copies:      5,
	}

	expected := &entities.BookResponse{
		ID:          bookID,
		Title:       book.Title,
		Author:      book.Author,
		ISBN:        book.ISBN,
		Publisher:   book.Publisher,
		PublishDate: book.PublishDate,
		Description: book.Description,
		Copies:      book.Copies,
		CreatedAt:   testTime,
		UpdatedAt:   testTime,
	}

	successMock := bookMock.Book{}
	successMock.On("GetByID", mock.Anything, bookID).Return(book, nil)

	notFoundMock := bookMock.Book{}
	notFoundMock.On("GetByID", mock.Anything, invalidID).Return(nil, errors.New("book not found"))

	tests := []struct {
		name    string
		s       *service
		id      uuid.UUID
		want    *entities.BookResponse
		wantErr bool
	}{
		{
			name: "found",
			s:    &service{model: models.Model{Book: &successMock}},
			id:   bookID,
			want: expected,
		},
		{
			name: "not found",
			s:    &service{model: models.Model{Book: &notFoundMock}},
			id:   invalidID,
			want: nil, wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetBookByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBookByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBookByID() = %v, want %v", got, tt.want)
			}
			tt.s.model.Book.(*bookMock.Book).AssertExpectations(t)
		})
	}
}

func Test_service_GetAllBooks(t *testing.T) {
	id1, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	now := time.Now()

	books := []entities.Book{
		{ID: id1, Title: "Book 1", Author: "Author 1", ISBN: "111", Publisher: "Pub 1", PublishDate: now, Description: "Desc 1", Copies: 2, CreatedAt: now, UpdatedAt: now},
		{ID: id2, Title: "Book 2", Author: "Author 2", ISBN: "222", Publisher: "Pub 2", PublishDate: now, Description: "Desc 2", Copies: 3, CreatedAt: now, UpdatedAt: now},
	}

	expected := []*entities.BookResponse{
		{ID: id1, Title: "Book 1", Author: "Author 1", ISBN: "111", Publisher: "Pub 1", PublishDate: now, Description: "Desc 1", Copies: 2, CreatedAt: now, UpdatedAt: now},
		{ID: id2, Title: "Book 2", Author: "Author 2", ISBN: "222", Publisher: "Pub 2", PublishDate: now, Description: "Desc 2", Copies: 3, CreatedAt: now, UpdatedAt: now},
	}

	successMock := bookMock.Book{}
	successMock.On("GetAll", mock.Anything).Return(books, nil)

	errorMock := bookMock.Book{}
	errorMock.On("GetAll", mock.Anything).Return(nil, errors.New("db error"))

	emptyMock := bookMock.Book{}
	emptyMock.On("GetAll", mock.Anything).Return([]entities.Book{}, nil)

	tests := []struct {
		name    string
		s       *service
		want    []*entities.BookResponse
		wantErr bool
	}{
		{
			name:    "successful retrieval",
			s:       &service{model: models.Model{Book: &successMock}},
			want:    expected,
			wantErr: false,
		},
		{
			name:    "database error",
			s:       &service{model: models.Model{Book: &errorMock}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty result",
			s:       &service{model: models.Model{Book: &emptyMock}},
			want:    []*entities.BookResponse{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetAllBooks(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllBooks() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllBooks() = %v, want %v", got, tt.want)
			}
			tt.s.model.Book.(*bookMock.Book).AssertExpectations(t)
		})
	}
}

func Test_service_UpdateBook(t *testing.T) {
	bookID, _ := uuid.NewV4()
	invalidID, _ := uuid.NewV4()
	now := time.Now()

	req := &entities.BookRequest{
		Title:       "Updated",
		Author:      "Updated Author",
		ISBN:        "999",
		Publisher:   "Updated Publisher",
		PublishDate: now,
		Description: "Updated Desc",
		Copies:      7,
	}

	existing := &entities.Book{
		ID: bookID, Title: "Old", Author: "Old Author", ISBN: "111", Publisher: "Old Pub",
		PublishDate: now, Description: "Old Desc", Copies: 1, CreatedAt: now, UpdatedAt: now,
	}

	successMock := bookMock.Book{}
	successMock.On("GetByID", mock.Anything, bookID).Return(existing, nil)
	successMock.On("Update", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		book := args.Get(1).(*entities.Book)
		book.Title = req.Title
		book.Author = req.Author
		book.ISBN = req.ISBN
		book.Publisher = req.Publisher
		book.PublishDate = req.PublishDate
		book.Description = req.Description
		book.Copies = req.Copies
	})

	notFoundMock := bookMock.Book{}
	notFoundMock.On("GetByID", mock.Anything, invalidID).Return(nil, errors.New("not found"))

	updateErrMock := bookMock.Book{}
	updateErrMock.On("GetByID", mock.Anything, bookID).Return(existing, nil)
	updateErrMock.On("Update", mock.Anything, mock.Anything).Return(errors.New("update error"))

	tests := []struct {
		name    string
		s       *service
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "success",
			s:       &service{model: models.Model{Book: &successMock}},
			id:      bookID,
			wantErr: false,
		},
		{
			name: "not found",
			s:    &service{model: models.Model{Book: &notFoundMock}},
			id:   invalidID, wantErr: true,
		},
		{
			name: "update error",
			s:    &service{model: models.Model{Book: &updateErrMock}},
			id:   bookID, wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.UpdateBook(context.Background(), tt.id, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBook() error = %v, wantErr %v", err, tt.wantErr)
			}

			tt.s.model.Book.(*bookMock.Book).AssertExpectations(t)
		})
	}
}

func Test_service_DeleteBook(t *testing.T) {
	bookID, _ := uuid.NewV4()
	invalidID, _ := uuid.NewV4()

	successMock := bookMock.Book{}
	successMock.On("Delete", mock.Anything, bookID).Return(nil)

	notFoundMock := bookMock.Book{}
	notFoundMock.On("Delete", mock.Anything, invalidID).Return(errors.New("not found"))

	tests := []struct {
		name    string
		s       *service
		id      uuid.UUID
		wantErr bool
	}{
		{
			name: "delete success",
			s:    &service{model: models.Model{Book: &successMock}},
			id:   bookID,
		},
		{
			name:    "not found",
			s:       &service{model: models.Model{Book: &notFoundMock}},
			id:      invalidID,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.DeleteBook(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBook() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.s.model.Book.(*bookMock.Book).AssertExpectations(t)
		})
	}
}
