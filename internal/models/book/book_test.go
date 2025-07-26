package book

import (
	"context"
	"database/sql"
	"log"
	"regexp"
	"testing"
	"time"

	"library-system/internal/entities"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMock() (*gorm.DB, *sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening gorm stub database connection", err)
	}
	return gormDB, db, mock
}

func Test_book_Create(t *testing.T) {
	gdb, db, mock := NewMock()
	defer db.Close()

	testTime := time.Now()

	validBook := entities.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "1234567890",
		Publisher:   "Test Publisher",
		PublishDate: testTime,
		Description: "Test Description",
		Copies:      5,
	}

	invalidBook := validBook
	invalidBook.ISBN = "" // Simulate DB failure

	insertStmt := regexp.QuoteMeta(`INSERT INTO "books" ("id","created_at","updated_at","deleted_at","title","author","isbn","publisher","publish_date","description","copies") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`)

	// Setup valid expectation
	mock.ExpectBegin()
	mock.ExpectExec(insertStmt).
		WithArgs(
			sqlmock.AnyArg(), // id
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			validBook.Title,
			validBook.Author,
			validBook.ISBN,
			validBook.Publisher,
			validBook.PublishDate,
			validBook.Description,
			validBook.Copies,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Setup invalid expectation
	mock.ExpectBegin()
	mock.ExpectExec(insertStmt).
		WithArgs(
			sqlmock.AnyArg(), // id
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			invalidBook.Title,
			invalidBook.Author,
			invalidBook.ISBN,
			invalidBook.Publisher,
			invalidBook.PublishDate,
			invalidBook.Description,
			invalidBook.Copies,
		).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	type args struct {
		ctx  context.Context
		book *entities.Book
	}
	tests := []struct {
		name    string
		b       *book
		args    args
		wantErr bool
	}{
		{
			name:    "valid case",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background(), book: &validBook},
			wantErr: false,
		},
		{
			name:    "insert error",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background(), book: &invalidBook},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Create(tt.args.ctx, tt.args.book); (err != nil) != tt.wantErr {
				t.Errorf("book.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_book_GetByID(t *testing.T) {
	gdb, db, mock := NewMock()
	defer db.Close()

	// Create a valid book ID for testing
	validID, _ := uuid.NewV4()
	invalidID, _ := uuid.NewV4()
	testTime := time.Now()

	// Create a book for the result
	validBook := entities.Book{
		ID:          validID,
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

	// Mock rows result
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "title", "author", "isbn", "publisher", "publish_date", "description", "copies"}).
		AddRow(validBook.ID, validBook.CreatedAt, validBook.UpdatedAt, validBook.Title, validBook.Author, validBook.ISBN, validBook.Publisher, validBook.PublishDate, validBook.Description, validBook.Copies)

	// Generate SQL statements for mocking
	validStmt := gdb.Session(&gorm.Session{DryRun: true}).Model(&entities.Book{}).Where("id = ?", validID).First(&entities.Book{}).Statement.SQL.String()
	invalidStmt := gdb.Session(&gorm.Session{DryRun: true}).Model(&entities.Book{}).Where("id = ?", invalidID).First(&entities.Book{}).Statement.SQL.String()

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta(validStmt)).WithArgs().WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(invalidStmt)).WithArgs().WillReturnError(gorm.ErrRecordNotFound)

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		b       *book
		args    args
		want    *entities.Book
		wantErr bool
	}{
		{
			name:    "valid case",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background(), id: validID},
			want:    &validBook,
			wantErr: false,
		},
		{
			name:    "not found error",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background(), id: invalidID},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("book.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if got.ID != tt.want.ID || got.Title != tt.want.Title || got.ISBN != tt.want.ISBN {
					t.Errorf("book.GetByID() got = %v, want %v", got, tt.want)
				}
			} else if (got == nil) != (tt.want == nil) {
				t.Errorf("book.GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_book_GetAll(t *testing.T) {
	gdb, db, mock := NewMock()
	defer db.Close()

	// Create test data
	id1, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	testTime := time.Now()

	books := []entities.Book{
		{
			ID:          id1,
			CreatedAt:   testTime,
			UpdatedAt:   testTime,
			Title:       "Test Book 1",
			Author:      "Test Author 1",
			ISBN:        "1234567890",
			Publisher:   "Test Publisher 1",
			PublishDate: testTime,
			Description: "Test Description 1",
			Copies:      5,
		},
		{
			ID:          id2,
			CreatedAt:   testTime,
			UpdatedAt:   testTime,
			Title:       "Test Book 2",
			Author:      "Test Author 2",
			ISBN:        "0987654321",
			Publisher:   "Test Publisher 2",
			PublishDate: testTime,
			Description: "Test Description 2",
			Copies:      3,
		},
	}

	// Mock rows result
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "title", "author", "isbn", "publisher", "publish_date", "description", "copies"})
	for _, book := range books {
		rows.AddRow(book.ID, book.CreatedAt, book.UpdatedAt, book.Title, book.Author, book.ISBN, book.Publisher, book.PublishDate, book.Description, book.Copies)
	}

	// Generate SQL statement for mocking
	validStmt := gdb.Session(&gorm.Session{DryRun: true}).Model(&entities.Book{}).Find(&[]entities.Book{}).Statement.SQL.String()
	errorStmt := validStmt

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta(validStmt)).WithArgs().WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(errorStmt)).WithArgs().WillReturnError(sql.ErrConnDone)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		b       *book
		args    args
		want    []entities.Book
		wantErr bool
	}{
		{
			name:    "valid case",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background()},
			want:    books,
			wantErr: false,
		},
		{
			name:    "database error",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background()},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("book.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("book.GetAll() got %d books, want %d books", len(got), len(tt.want))
			}
		})
	}
}

func Test_book_Update(t *testing.T) {
	gdb, db, mock := NewMock()
	defer db.Close()

	testTime := time.Now()

	validBook := entities.Book{
		ID:          uuid.Must(uuid.NewV4()),
		CreatedAt:   testTime,
		UpdatedAt:   testTime,
		Title:       "Updated Book",
		Author:      "Updated Author",
		ISBN:        "1234567890",
		Publisher:   "Updated Publisher",
		PublishDate: testTime,
		Description: "Updated Description",
		Copies:      10,
	}

	nonExistentBook := validBook
	nonExistentBook.ID = uuid.Must(uuid.NewV4())

	// Adjusted to include "deleted_at"
	updateStmt := regexp.QuoteMeta(`UPDATE "books" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"title"=$4,"author"=$5,"isbn"=$6,"publisher"=$7,"publish_date"=$8,"description"=$9,"copies"=$10 WHERE "id" = $11`)

	// --- VALID BOOK EXPECTATION ---
	mock.ExpectBegin()
	mock.ExpectExec(updateStmt).WithArgs(
		sqlmock.AnyArg(), // created_at
		sqlmock.AnyArg(), // updated_at
		sqlmock.AnyArg(), // deleted_at (zero/null)
		validBook.Title,
		validBook.Author,
		validBook.ISBN,
		validBook.Publisher,
		validBook.PublishDate,
		validBook.Description,
		validBook.Copies,
		sqlmock.AnyArg(), // ID
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// --- BOOK NOT FOUND EXPECTATION ---
	mock.ExpectBegin()
	mock.ExpectExec(updateStmt).WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		nonExistentBook.Title,
		nonExistentBook.Author,
		nonExistentBook.ISBN,
		nonExistentBook.Publisher,
		nonExistentBook.PublishDate,
		nonExistentBook.Description,
		nonExistentBook.Copies,
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectCommit()

	// --- TEST CASES ---
	type args struct {
		ctx  context.Context
		book *entities.Book
	}
	tests := []struct {
		name    string
		b       *book
		args    args
		wantErr bool
	}{
		{
			name:    "valid case",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background(), book: &validBook},
			wantErr: false,
		},
		{
			name:    "book not found",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background(), book: &nonExistentBook},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Update(tt.args.ctx, tt.args.book); (err != nil) != tt.wantErr {
				t.Errorf("book.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_book_Delete(t *testing.T) {
	gdb, db, mock := NewMock()
	defer db.Close()

	// Create valid and invalid IDs for testing
	validID, _ := uuid.NewV4()
	invalidID, _ := uuid.NewV4()

	// Generate SQL statements for mocking
	validStmt := gdb.Session(&gorm.Session{DryRun: true}).Delete(&entities.Book{}, "id = ?", validID).Statement.SQL.String()
	invalidStmt := gdb.Session(&gorm.Session{DryRun: true}).Delete(&entities.Book{}, "id = ?", invalidID).Statement.SQL.String()

	// Set up expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(validStmt)).WithArgs(validID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(invalidStmt)).WithArgs(invalidID).WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectCommit()

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		b       *book
		args    args
		wantErr bool
	}{
		{
			name:    "valid case",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background(), id: validID},
			wantErr: false,
		},
		{
			name:    "book not found",
			b:       &book{db: gdb},
			args:    args{ctx: context.Background(), id: invalidID},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("book.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
