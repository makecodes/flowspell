package main

import (
    "log"
    "os"
    "flowspell/db"
    "flowspell/handlers"
    "flowspell/models"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/joho/godotenv"
    "fmt"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    if len(os.Args) > 1 {
        command := os.Args[1]
        switch command {
        case "migrate":
            db.MigrateUp()
            return
        case "rollback":
            db.MigrateDown()
            return
        }
    }

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    db.AutoMigrate(&models.User{})

    userHandler := handlers.UserHandler{DB: db}

    router := gin.Default()

    router.GET("/users", userHandler.GetUsers)
    router.GET("/users/:id", userHandler.GetUser)
    router.POST("/users", userHandler.CreateUser)
    router.PUT("/users/:id", userHandler.UpdateUser)
    router.DELETE("/users/:id", userHandler.DeleteUser)

    router.Run(":8266")
}

