package main

import (
	"flowspell/db"
	"flowspell/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"flowspell/docs"
)

func main() {
	docs.SwaggerInfo.Title = "FlowSpell"
	docs.SwaggerInfo.Description = "FlowSpell API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "http://localhost:8266"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http"}

	godotenv.Load()

	dbConnection, err := db.GetDBConnection()
	if err != nil {
		log.Fatalf("Error connecting to database")
	}

	flowDefinitionHandler := handlers.FlowDefinitionHandler{DB: dbConnection}
	flowInstanceHandler := handlers.FlowInstanceHandler{DB: dbConnection}

	router := gin.Default()

	flowsGroup := router.Group("/flows")
	{
		flowDefinitionsGroup := flowsGroup.Group("/definitions")
		{
			flowDefinitionsGroup.GET("/", flowDefinitionHandler.GetFlowDefinitions)
			flowDefinitionsGroup.POST("/", flowDefinitionHandler.CreateFlowDefinition)
			flowDefinitionsGroup.GET("/:referenceId", flowDefinitionHandler.GetFlowDefinition)
			flowDefinitionsGroup.PUT("/:referenceId", flowDefinitionHandler.UpdateFlowDefinition)
			flowDefinitionsGroup.PATCH("/:referenceId", flowDefinitionHandler.UpdateFlowDefinition)
			flowDefinitionsGroup.DELETE("/:referenceId", flowDefinitionHandler.DeleteFlowDefinition)
		}

		flowInstancesGroup := flowsGroup.Group("/instances")
		{
			flowInstancesGroup.GET("/", flowInstanceHandler.GetFlowInstances)
			// flowInstancesGroup.POST("/", flowInstanceHandler.CreateFlowInstance)
			flowInstancesGroup.POST("/:referenceId/start", flowInstanceHandler.StartFlowInstance)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8266")
}
