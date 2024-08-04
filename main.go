package main

import (
	"flowspell/db"
	"flowspell/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"flowspell/docs"
)

// @title FlowSpell API
// @version 1.0
// @description FlowSpell will handle the flows of your applications
// @termsOfService https://flowspell.org/terms/

// @host localhost:8266
// @BasePath /
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
	taskDefinitionHandler := handlers.TaskDefinitionHandler{DB: dbConnection}

	router := gin.Default()

	flowsGroup := router.Group("/flows")
	{
		flowDefinitionsGroup := flowsGroup.Group("/definitions")
		{
			flowDefinitionsGroup.GET("/", flowDefinitionHandler.GetFlowDefinitions)
			flowDefinitionsGroup.POST("/", flowDefinitionHandler.CreateFlowDefinition)
			flowDefinitionsGroup.GET("/:referenceId", flowDefinitionHandler.GetFlowDefinition)
			flowDefinitionsGroup.PUT("/:referenceId", flowDefinitionHandler.UpdateFlowDefinition)
			flowDefinitionsGroup.DELETE("/:referenceId", flowDefinitionHandler.DeleteFlowDefinition)
		}

		flowInstancesGroup := flowsGroup.Group("/instances")
		{
			flowInstancesGroup.GET("/", flowInstanceHandler.GetFlowInstances)
			flowInstancesGroup.POST("/:referenceId/start", flowInstanceHandler.StartFlow)
		}
	}

	tasksGroup := router.Group("/tasks")
	{
		taskDefinitionsGroup := tasksGroup.Group("/definitions")
		{
			taskDefinitionsGroup.GET("/", taskDefinitionHandler.GetTaskDefinitions)
			taskDefinitionsGroup.POST("/", taskDefinitionHandler.CreateTaskDefinition)
            taskDefinitionsGroup.DELETE("/:referenceId", taskDefinitionHandler.DeleteTaskDefinition)
		}
	}

	// JSONSchema
	router.GET("/schemas/flow_definitions/:referenceId/:type", flowDefinitionHandler.GetFlowDefinitionSchema)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // flowSpellURL := os.Getenv("FLOWSPELL_HOST")
    // url := ginSwagger.URL(flowSpellURL + "/swagger/openapi.json")
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Run(":8266")
}
