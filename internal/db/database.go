// ==============================================================================
// simply todo - Checklist Backend Database Logic (database.go)
// ==============================================================================
// Repository: checklist-backend
// Path: internal/db/database.go
// Purpose: Manages the PostgreSQL connection pool and health handshake.
// License: Apache 2.0
// ==============================================================================

package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool" // Requires: go get github.com/jackc/pgx/v5
)

// ConnectVault establishes a high-performance connection pool to the Postgres-db.
// It includes a retry loop to account for the database container startup time.
func ConnectVault(databaseURL string) (*pgxpool.Pool, error) {
	log.Printf("[INFO] checklist-backend: Initiating handshake with vault at %s", "todo-db:5432")

	// 1. Connection Pool Configuration
	// pgxpool handles multiple concurrent requests automatically, matching 
	// the high-concurrency architecture of Go.
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Printf("[ERROR] Unable to parse DATABASE_URL: %v", err)
		return nil, err
	}

	// 2. The Handshake Loop (Retry Logic)
	// Because the DB might be cold-starting (especially on Render or first Docker boot),
	// we attempt to connect multiple times before failing.
	var pool *pgxpool.Pool
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		log.Printf("[INFO] Connection attempt %d/%d...", i, maxRetries)
		
		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			// Validate the connection is actually alive
			err = pool.Ping(context.Background())
			if err == nil {
				log.Println("[SUCCESS] checklist-backend: Handshake complete. Vault is open.")
				return pool, nil
			}
		}

		log.Printf("[WARNING] Vault not ready. Retrying in 3 seconds... (%v)", err)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to vault after %d attempts", maxRetries)
}

// CloseVault safely drains the connection pool. 
// Should be called via 'defer' in main.go.
func CloseVault(pool *pgxpool.Pool) {
	if pool != nil {
		log.Println("[INFO] checklist-backend: Draining connection pool and closing vault.")
		pool.Close()
	}
}

// ------------------------------------------------------------------------------
// Implementation Notes:
// ------------------------------------------------------------------------------
// 1. Connection Pooling: Using pgxpool instead of a single connection ensures 
//    the backend can handle multiple simultaneous CRUD operations.
//
// 2. Healthcheck: The .Ping() call ensures we don't just have a network 
//    socket, but an authenticated, active PostgreSQL session.
//
// 3. Resilience: The 15-second total retry window (5 attempts * 3s) matches 
//    the expected warm-up time for our postgres:15-alpine container.
// ------------------------------------------------------------------------------
