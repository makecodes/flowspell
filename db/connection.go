package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *sql.DB

func init() {
	godotenv.Load()
}

func GetDatabaseConnectionString() string {
	environment := os.Getenv("ENV")
	var dbURL string
	if environment == "test" {
		dbURL = os.Getenv("DATABASE_TEST_URL")
	} else {
		dbURL = os.Getenv("DATABASE_URL")
	}

	if dbURL == "" {
		log.Println("Database URL not found in environment variables")
	}
	return dbURL
}

func GetDBConnection() (*gorm.DB, error) {
	dsn := GetDatabaseConnectionString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

var sqlOpenFunc = func(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

func NewConnection() (*sql.DB, error) {
	dsn := GetDatabaseConnectionString()
	db, err := sqlOpenFunc("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	return db, nil
}
