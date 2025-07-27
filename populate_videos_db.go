package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "video_analysis.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	videoDir := "videos"
	files, err := ioutil.ReadDir(videoDir)
	if err != nil {
		log.Fatal("Failed to read videos directory:", err)
	}

	ins := `INSERT INTO videos (id, filename, original_filename, file_size, status, created_at, updated_at) VALUES (?, ?, ?, ?, 'uploaded', ?, ?)`
	count := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !isVideoFile(name) {
			continue
		}
		id := generateUUID()
		filename := name
		originalFilename := name
		fileSize := file.Size()
		now := time.Now().Format(time.RFC3339)
		_, err := db.Exec(ins, id, filename, originalFilename, fileSize, now, now)
		if err != nil {
			log.Printf("Failed to insert %s: %v", name, err)
			continue
		}
		fmt.Printf("Inserted: %s\n", name)
		count++
	}
	fmt.Printf("Inserted %d videos into the database.\n", count)
}

func isVideoFile(name string) bool {
	exts := []string{".mp4", ".avi", ".mov", ".mkv", ".wmv", ".flv", ".webm"}
	for _, ext := range exts {
		if filepath.Ext(name) == ext {
			return true
		}
	}
	return false
}

// generateUUID generates a random UUID (v4)
func generateUUID() string {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	b := make([]byte, 16)
	_, err = f.Read(b)
	f.Close()
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	// Set version and variant bits
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
