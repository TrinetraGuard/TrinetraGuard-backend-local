package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database connection (same as in database.go)
type Database struct {
	*sql.DB
}

func main() {
	// Same initialization as in database.go
	dbName := "video_analysis.db"

	// Ensure directory exists
	dir := filepath.Dir(dbName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal("Failed to create database directory:", err)
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	database := &Database{DB: db}
	fmt.Println("✅ Database connection successful")

	// Test the exact query from the service
	limit, offset := 10, 0

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM videos`
	err = database.DB.QueryRow(countQuery).Scan(&total)
	if err != nil {
		log.Fatal("Failed to count videos:", err)
	}
	fmt.Printf("✅ Total videos: %d\n", total)

	// Get videos with the exact query from service
	query := `SELECT id, filename, original_filename, file_size, duration, frame_count, 
			  width, height, status, created_at, updated_at FROM videos 
			  ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := database.DB.Query(query, limit, offset)
	if err != nil {
		log.Fatal("Failed to query videos:", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var id, filename, originalFilename, status, createdAt, updatedAt string
		var fileSize int
		var duration, frameCount, width, height interface{}

		err := rows.Scan(
			&id, &filename, &originalFilename, &fileSize,
			&duration, &frameCount, &width, &height,
			&status, &createdAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan row %d: %v", count, err)
			continue
		}
		count++
		fmt.Printf("Video %d: ID=%s, Name=%s, Size=%d, Status=%s\n",
			count, id, originalFilename, fileSize, status)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error iterating rows:", err)
	}

	fmt.Printf("✅ Successfully retrieved %d videos\n", count)

	// Don't close the database connection to simulate the service behavior
	fmt.Println("✅ Database connection kept open (simulating service behavior)")
}
