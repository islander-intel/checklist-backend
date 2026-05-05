// ==============================================================================
// simply todo - Checklist Backend Entry Point (main.go)
// ==============================================================================
// Repository: checklist-backend
// Path: main.go
// Purpose: Main execution thread for the Go Worker tier.
// License: Apache 2.0
// ==============================================================================

package main

import (
	"log"
	"os"

	"checklist-backend/internal/db"
	"checklist-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Requires: go get github.com/joho/godotenv
)

func main() {
	// 1. Logging Initialization
	log.Println("[STARTUP] simply todo: checklist-backend is waking up...")

	// 2. Load Environment Variables
	// Looks for the .env file we defined in our project manifest.
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARNING] No .env file found, falling back to system environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("[CRITICAL] DATABASE_URL is not set. checklist-backend cannot function without its vault.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001" // Default port for the Worker tier in our 4-tier architecture
	}

	// 3. Database Handshake
	// ConnectVault includes our 5-attempt retry loop to wait for the DB container.
	vaultPool, err := db.ConnectVault(dbURL)
	if err != nil {
		log.Fatalf("[CRITICAL] Failed to open the vault: %v", err)
	}
	defer db.CloseVault(vaultPool)

	// 4. Initialize Handlers
	taskHandler := handlers.NewTaskHandler(vaultPool)

	// 5. Router Setup (Using Gin for high performance)
	// We use gin.Default() which includes basic Logging and Recovery middleware.
	router := gin.Default()

	// 6. Route Registration
	// These endpoints match the API Contract defined in our DATASHEET.md.
	v1 := router.Group("/api/v1")
	{
		v1.POST("/tasks", taskHandler.CreateTask)      // CREATE
		v1.GET("/tasks", taskHandler.GetTasks)        // READ (ALL)
		v1.PUT("/tasks/:id", taskHandler.UpdateTask)   // UPDATE
		v1.DELETE("/tasks/:id", taskHandler.DeleteTask) // DELETE
	}

	// 7. Health Check (For Router/Orchestrator monitoring)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "worker": "alive", "vault": "connected"})
	})

	// 8. Execution
	log.Printf("[READY] checklist-backend: Listening on port %s", port)
	log.Printf("[INFO] Mapping: POST/GET/PUT/DELETE /api/v1/tasks")
	
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("[CRITICAL] Server failed to start: %v", err)
	}
}

/*
LOGGING TRACE:
[STARTUP] checklist-backend: Environment loaded.
[INFO] Attempting to reach postgres-db:5432...
[SUCCESS] Handshake complete. Vault is open.
[READY] Server listening on :8001
*/
