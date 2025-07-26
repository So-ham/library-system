package postgres

import (
	"os"
	"time"

	"library-system/internal/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		panic("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		panic("failed to ping database: " + err.Error())
	}

	if err := db.AutoMigrate(
		&entities.Book{},
	); err != nil {
		panic("failed to auto-migrate database: " + err.Error())
	}

	if err := SeedData(db); err != nil {
		panic("failed to seed database: " + err.Error())
	}

	return db
}
