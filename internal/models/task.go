// ==============================================================================
// simply todo - Checklist Backend Models (task.go)
// ==============================================================================
// Repository: checklist-backend
// Path: internal/models/task.go
// Purpose: Defines the data structures (structs) that mirror the todo-db schema.
// License: Apache 2.0
// ==============================================================================

package models

import (
	"time"

	"github.com/google/uuid" // Requires: go get github.com/google/uuid
)

// Task represents the core "simply todo" data entity.
// This struct is used for scanning rows from the database and 
// for serializing responses to the Router/Frontend.
type Task struct {
	// ID maps to the UUID PRIMARY KEY in todo-db.
	// We use the google/uuid type to ensure native compatibility.
	ID uuid.UUID `json:"id"`

	// Title maps to VARCHAR(255) NOT NULL.
	Title string `json:"title"`

	// Description maps to TEXT (NULLABLE). 
	// Using a pointer (*string) allows Go to represent a database NULL as nil.
	Description *string `json:"description"`

	// IsCompleted maps to BOOLEAN NOT NULL.
	IsCompleted bool `json:"is_completed"`

	// CreatedAt maps to TIMESTAMPTZ.
	// We use time.Time which handles UTC offsets automatically.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt maps to TIMESTAMPTZ.
	// This field is managed by the PostgreSQL trigger in the database vault.
	UpdatedAt time.Time `json:"updated_at"`
}

// ------------------------------------------------------------------------------
// Implementation Notes & Constraints:
// ------------------------------------------------------------------------------
// 1. UUIDs: The 'id' field uses the google/uuid library. During DB scans, 
//    ensure the pgx driver is configured to handle UUID -> uuid.UUID mapping.
//
// 2. Timestamps: We use 'time.Time' for both created_at and updated_at. 
//    Postgres 'TIMESTAMPTZ' ensures that these are always stored in UTC, 
//    preventing "time-drift" bugs between the backend and the database.
//
// 3. Nullability: By using *string for Description, the JSON output will 
//    render as "null" if the database field is empty, matching the expected 
//    behavior for modern REST APIs.
// ------------------------------------------------------------------------------

/*
LOGGING TRACE:
[INFO] Model "Task" initialized.
[INFO] Mapping: ID (uuid) -> Title (string) -> Description (*string) -> IsCompleted (bool)
[INFO] Handshake: Models synchronized with DATASHEET.md v1.0
*/
