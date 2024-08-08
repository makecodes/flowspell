package server

import (
	"flowspell/db"
	"flowspell/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func FlowSpellServer() {
    dbConnection, err := db.GetDBConnection()
	if err != nil {
		log.Fatalf("Error connecting to database")
	}

    flowDefinitionHandler := handlers.FlowDefinitionHandler{DB: dbConnection}
	flowInstanceHandler := handlers.FlowInstanceHandler{DB: dbConnection}
	taskDefinitionHandler := handlers.TaskDefinitionHandler{DB: dbConnection}
	taskInstanceHandler := handlers.TaskInstanceHandler{DB: dbConnection}

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
			taskDefinitionsGroup.GET("/:referenceId", taskDefinitionHandler.GetTaskDefinition)
			taskDefinitionsGroup.DELETE("/:referenceId", taskDefinitionHandler.DeleteTaskDefinition)
		}

		tasksGroup.POST("/queue", taskInstanceHandler.GetTaskQueue)
	}

	// JSONSchema
	router.GET("/schemas/flow_definitions/:referenceId/:type", flowDefinitionHandler.GetFlowDefinitionSchema)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// flowSpellURL := os.Getenv("FLOWSPELL_HOST")
	// url := ginSwagger.URL(flowSpellURL + "/swagger/openapi.json")
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Run(":8266")
}