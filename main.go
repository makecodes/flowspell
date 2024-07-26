package main

import (
    "log"
    "flowspell/handlers"
    "flowspell/db"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    dbConnection, err := db.GetDBConnection()
    if err != nil {
        log.Fatalf("Error connecting to database")
    }

    flowDefinitionHandler := handlers.FlowDefinitionHandler{DB: dbConnection}

    router := gin.Default()

    florDefinitionsGroup := router.Group("/flows/definitions")
    {
        florDefinitionsGroup.GET("/", flowDefinitionHandler.GetFlowDefinitions)
        florDefinitionsGroup.POST("/", flowDefinitionHandler.CreateFlowDefinition)
        florDefinitionsGroup.GET("/:id", flowDefinitionHandler.GetFlowDefinition)
        florDefinitionsGroup.PUT("/:id", flowDefinitionHandler.UpdateFlowDefinition)
        florDefinitionsGroup.DELETE("/:id", flowDefinitionHandler.DeleteFlowDefinition)
    }

    router.Run(":8266")
}

