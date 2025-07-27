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

	// Test count query
	var total int
	countQuery := `SELECT COUNT(*) FROM videos`
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		log.Fatal("Failed to count videos:", err)
	}
	fmt.Printf("✅ Total videos in database: %d\n", total)

	// Test select query
	query := `SELECT id, filename, original_filename, file_size, duration, frame_count, 
			  width, height, status, created_at, updated_at FROM videos 
			  ORDER BY created_at DESC LIMIT 10 OFFSET 0`

	rows, err := db.Query(query)
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
}
