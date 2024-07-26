package db

import (
    "log"
    "os"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    _ "github.com/lib/pq"
    "gorm.io/gorm"
    "database/sql"
    "fmt"
)

var db *sql.DB

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found")
    }
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
