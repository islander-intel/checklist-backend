// ==============================================================================
// simply todo - Checklist Backend Handlers (tasks.go)
// ==============================================================================
// Repository: checklist-backend
// Path: internal/handlers/tasks.go
// Purpose: Implements CRUD logic for the tasks domain.
// License: Apache 2.0
// ==============================================================================

package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid" // Requires: go get github.com/google/uuid
)

// TaskHandler encapsulates the database pool to keep handlers stateless
type TaskHandler struct {
	db *pgxpool.Pool
}

// NewTaskHandler initializes the handler with a validated DB connection
func NewTaskHandler(db *pgxpool.Pool) *TaskHandler {
	return &TaskHandler{db: db}
}

// TaskRequest defines the incoming JSON structure for Create/Update
// Matches the DATASHEET.md contract.
type TaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

// ------------------------------------------------------------------------------
// CREATE: POST /tasks
// ------------------------------------------------------------------------------
func (h *TaskHandler) CreateTask(c *gin.Context) {
	log.Println("[INFO] checklist-backend: Received request to CREATE task")

	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: title is required"})
		return
	}

	query := `
		INSERT INTO tasks (title, description, is_completed)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	var id uuid.UUID
	var createdAt, updatedAt string
	err := h.db.QueryRow(context.Background(), query, req.Title, req.Description, req.IsCompleted).
		Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		log.Printf("[ERROR] DB Insertion failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task to vault"})
		return
	}

	log.Printf("[SUCCESS] Task created with UUID: %s", id)
	c.JSON(http.StatusCreated, gin.H{
		"id":           id,
		"title":        req.Title,
		"description":  req.Description,
		"is_completed": req.IsCompleted,
		"created_at":   createdAt,
		"updated_at":   updatedAt,
	})
}

// ------------------------------------------------------------------------------
// READ: GET /tasks
// ------------------------------------------------------------------------------
func (h *TaskHandler) GetTasks(c *gin.Context) {
	log.Println("[INFO] checklist-backend: Fetching all tasks")

	query := `SELECT id, title, description, is_completed, created_at, updated_at FROM tasks ORDER BY created_at DESC;`
	rows, err := h.db.Query(context.Background(), query)
	if err != nil {
		log.Printf("[ERROR] Query failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	defer rows.Close()

	var tasks []gin.H
	for rows.Next() {
		var id uuid.UUID
		var title, description, createdAt, updatedAt string
		var isCompleted bool

		err := rows.Scan(&id, &title, &description, &isCompleted, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("[ERROR] Row scan failed: %v", err)
			continue
		}
		tasks = append(tasks, gin.H{
			"id":           id,
			"title":        title,
			"description":  description,
			"is_completed": isCompleted,
			"created_at":   createdAt,
			"updated_at":   updatedAt,
		})
	}

	c.JSON(http.StatusOK, tasks)
}

// ------------------------------------------------------------------------------
// UPDATE: PUT /tasks/:id
// ------------------------------------------------------------------------------
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	log.Printf("[INFO] checklist-backend: Updating task %s", idStr)

	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Note: updated_at is handled by the Postgres trigger we wrote in todo-db/init.sql
	query := `
		UPDATE tasks 
		SET title = $1, description = $2, is_completed = $3 
		WHERE id = $4
		RETURNING updated_at;
	`

	var updatedAt string
	err := h.db.QueryRow(context.Background(), query, req.Title, req.Description, req.IsCompleted, idStr).Scan(&updatedAt)

	if err != nil {
		log.Printf("[ERROR] Update failed for %s: %v", idStr, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated", "updated_at": updatedAt})
}

// ------------------------------------------------------------------------------
// DELETE: DELETE /tasks/:id
// ------------------------------------------------------------------------------
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	log.Printf("[INFO] checklist-backend: Deleting task %s", idStr)

	query := `DELETE FROM tasks WHERE id = $1;`
	tag, err := h.db.Exec(context.Background(), query, idStr)

	if err != nil || tag.RowsAffected() == 0 {
		log.Printf("[ERROR] Delete failed for %s", idStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
