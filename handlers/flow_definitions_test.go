package handlers

import (
	"flowspell/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_TEST_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create the enum type before migrating the models
	db.Exec("CREATE TYPE flow_definitions_status AS ENUM ('active', 'inactive');")

	// Auto migrate the models
	db.AutoMigrate(&models.FlowDefinition{})
	return db, nil
}

// func TestGetFlowDefinitions(t *testing.T) {
// 	db, err := setupTestDB()
// 	assert.NoError(t, err)

// 	handler := &FlowDefinitionHandler{DB: db}

// 	router := gin.Default()
// 	router.GET("/flow/definitions", handler.GetFlowDefinitions)

// 	req, _ := http.NewRequest("GET", "/flow/definitions", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "[]") // Expecting an empty array response
// }
