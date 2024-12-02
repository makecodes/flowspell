package server

import (
	"flowspell/db"
	"flowspell/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

const referenceIdPath = "/:referenceId"

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
			flowDefinitionsGroup.GET(referenceIdPath, flowDefinitionHandler.GetFlowDefinition)
			flowDefinitionsGroup.PUT(referenceIdPath, flowDefinitionHandler.UpdateFlowDefinition)
			flowDefinitionsGroup.DELETE(referenceIdPath, flowDefinitionHandler.DeleteFlowDefinition)
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
			taskDefinitionsGroup.GET(referenceIdPath, taskDefinitionHandler.GetTaskDefinition)
			taskDefinitionsGroup.DELETE(referenceIdPath, taskDefinitionHandler.DeleteTaskDefinition)
		}

		tasksGroup.POST("/queue", taskInstanceHandler.GetTaskQueue)
	}

	// JSONSchema
	router.GET("/schemas/flow_definitions/:referenceId/:type", flowDefinitionHandler.GetFlowDefinitionSchema)

	err = router.Run(":8266")
	if err != nil {
		return
	}
}
