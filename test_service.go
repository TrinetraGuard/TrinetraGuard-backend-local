package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "./video_analysis.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("✅ Database connection successful")

	// Test the exact query from the service
	limit, offset := 10, 0

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM videos`
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		log.Fatal("Failed to count videos:", err)
	}
	fmt.Printf("✅ Total videos: %d\n", total)

	// Get videos with the exact query from service
	query := `SELECT id, filename, original_filename, file_size, duration, frame_count, 
			  width, height, status, created_at, updated_at FROM videos 
			  ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		log.Fatal("Failed to query videos:", err)
	}
	defer rows.Close()

	var videos []map[string]interface{}
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
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		video := map[string]interface{}{
			"id":                id,
			"filename":          filename,
			"original_filename": originalFilename,
			"file_size":         fileSize,
			"duration":          duration,
			"frame_count":       frameCount,
			"width":             width,
			"height":            height,
			"status":            status,
			"created_at":        createdAt,
			"updated_at":        updatedAt,
		}
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error iterating rows:", err)
	}

	fmt.Printf("✅ Successfully retrieved %d videos\n", len(videos))
	for i, video := range videos {
		fmt.Printf("Video %d: ID=%s, Name=%s, Size=%d, Status=%s\n",
			i+1, video["id"], video["original_filename"], video["file_size"], video["status"])
	}
}
