package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Video struct {
	ID          int    `json:"id"`
	Filename    string `json:"filename"`
	UploadTime  string `json:"upload_time"`
	TotalPeople int    `json:"total_people"`
	Faces       []Face `json:"faces"`
}

type Face struct {
	ID         int      `json:"id"`
	VideoID    int      `json:"video_id"`
	ImagePath  string   `json:"image_path"`
	Timestamps []string `json:"timestamps"`
	Name       string   `json:"name"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create videos table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS videos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT NOT NULL,
			upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			total_people INTEGER DEFAULT 0
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create faces table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS faces (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			video_id INTEGER,
			image_path TEXT NOT NULL,
			timestamps TEXT,
			name TEXT DEFAULT 'Unknown',
			FOREIGN KEY (video_id) REFERENCES videos (id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func resetDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Resetting database...")

	// Close current database connection
	if db != nil {
		db.Close()
	}

	// Delete the database file
	err := os.Remove("database.db")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error removing database: %v", err)
		http.Error(w, "Failed to reset database", http.StatusInternalServerError)
		return
	}

	// Clear faces directory
	facesDir := "faces"
	if _, err := os.Stat(facesDir); err == nil {
		err = os.RemoveAll(facesDir)
		if err != nil {
			log.Printf("Error clearing faces directory: %v", err)
		}
	}

	// Clear videos directory
	videosDir := "videos"
	if _, err := os.Stat(videosDir); err == nil {
		err = os.RemoveAll(videosDir)
		if err != nil {
			log.Printf("Error clearing videos directory: %v", err)
		}
	}

	// Recreate directories
	os.MkdirAll(facesDir, 0755)
	os.MkdirAll(videosDir, 0755)

	// Reinitialize database
	initDB()

	log.Println("Database reset completed successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Database reset successfully",
		"status":  "success",
	})
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create videos directory if it doesn't exist
	os.MkdirAll("videos", 0755)

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
	filepath := filepath.Join("videos", filename)

	// Save file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	_, err = dst.ReadFrom(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save to database
	result, err := db.Exec("INSERT INTO videos (filename) VALUES (?)", filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videoID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Start analysis in background
	go analyzeVideo(int(videoID), filepath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Video uploaded successfully",
		"filename": filename,
		"video_id": videoID,
	})
}

func analyzeVideo(videoID int, videoPath string) {
	log.Printf("Starting analysis for video ID %d: %s", videoID, videoPath)

	// Get absolute paths
	absVideoPath, err := filepath.Abs(videoPath)
	if err != nil {
		log.Printf("Error getting absolute video path: %v", err)
		return
	}

	absPythonDir, err := filepath.Abs("../python")
	if err != nil {
		log.Printf("Error getting absolute python directory: %v", err)
		return
	}

	absFacesDir, err := filepath.Abs("faces")
	if err != nil {
		log.Printf("Error getting absolute faces directory: %v", err)
		return
	}

	// Create faces directory if it doesn't exist
	os.MkdirAll(absFacesDir, 0755)

	// Run Python analysis script
	cmd := exec.Command("python3", "analyze_video_clean.py", absVideoPath, absFacesDir)
	cmd.Dir = absPythonDir
	// Set the Python path to use the virtual environment
	cmd.Env = append(os.Environ(), "PATH="+absPythonDir+"/venv/bin:"+os.Getenv("PATH"))

	output, err := cmd.Output()
	if err != nil {
		log.Printf("Python analysis failed: %v", err)
		log.Printf("Command was: python3 analyze_video_clean.py %s %s", absVideoPath, absFacesDir)
		log.Printf("Working directory: %s", absPythonDir)
		return
	}

	// Parse JSON output from Python script
	outputStr := string(output)
	jsonStart := strings.Index(outputStr, "{")
	if jsonStart == -1 {
		log.Printf("No JSON found in Python output")
		return
	}

	jsonStr := outputStr[jsonStart:]
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		log.Printf("JSON string: %s", jsonStr)
		return
	}

	// Update video with total people count
	totalPeople := int(result["total_people"].(float64))
	_, err = db.Exec("UPDATE videos SET total_people = ? WHERE id = ?", totalPeople, videoID)
	if err != nil {
		log.Printf("Error updating video total_people: %v", err)
		return
	}

	// Save faces to database
	faces := result["faces"].([]interface{})
	for _, faceInterface := range faces {
		face := faceInterface.(map[string]interface{})
		imagePath := face["image"].(string)
		timestamps := face["timestamps"].([]interface{})

		// Convert timestamps to JSON string
		timestampsJSON, err := json.Marshal(timestamps)
		if err != nil {
			log.Printf("Error marshaling timestamps: %v", err)
			continue
		}

		_, err = db.Exec("INSERT INTO faces (video_id, image_path, timestamps) VALUES (?, ?, ?)",
			videoID, imagePath, string(timestampsJSON))
		if err != nil {
			log.Printf("Error inserting face: %v", err)
			continue
		}

		// Get the inserted face ID for logging
		var faceID int
		err = db.QueryRow("SELECT last_insert_rowid()").Scan(&faceID)
		if err != nil {
			log.Printf("Error getting face ID: %v", err)
		} else {
			log.Printf("Added face: ID=%d, Image=%s, Name=Unknown", faceID, imagePath)
		}
	}

	log.Printf("Analysis completed for video ID %d", videoID)
}

func getVideos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, filename, upload_time, total_people FROM videos ORDER BY upload_time DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var video Video
		err := rows.Scan(&video.ID, &video.Filename, &video.UploadTime, &video.TotalPeople)
		if err != nil {
			log.Printf("Error scanning video row: %v", err)
			continue
		}
		videos = append(videos, video)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func getVideoAnalysis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	// Get video info
	var video Video
	err = db.QueryRow("SELECT id, filename, upload_time, total_people FROM videos WHERE id = ?", videoID).Scan(
		&video.ID, &video.Filename, &video.UploadTime, &video.TotalPeople)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	// Get faces for this video
	rows, err := db.Query("SELECT id, image_path, timestamps, name FROM faces WHERE video_id = ?", videoID)
	if err != nil {
		log.Printf("Error querying faces: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var faces []Face
	for rows.Next() {
		var face Face
		var timestampsJSON string
		err := rows.Scan(&face.ID, &face.ImagePath, &timestampsJSON, &face.Name)
		if err != nil {
			log.Printf("Error scanning face row: %v", err)
			continue
		}

		// Parse timestamps JSON
		var timestamps []string
		err = json.Unmarshal([]byte(timestampsJSON), &timestamps)
		if err != nil {
			log.Printf("Error parsing timestamps JSON: %v", err)
			timestamps = []string{}
		}

		face.VideoID = videoID
		face.Timestamps = timestamps
		faces = append(faces, face)
	}

	video.Faces = faces

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

func serveVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	http.ServeFile(w, r, filepath.Join("videos", filename))
}

func serveFace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	http.ServeFile(w, r, filepath.Join("faces", filename))
}

func updateFaceName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	faceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid face ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE faces SET name = ? WHERE id = ?", request.Name, faceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Face name updated successfully",
		"face_id": faceID,
		"name":    request.Name,
	})
}

func main() {
	initDB()

	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/upload", uploadVideo).Methods("POST")
	r.HandleFunc("/api/videos", getVideos).Methods("GET")
	r.HandleFunc("/api/videos/{id}", getVideoAnalysis).Methods("GET")
	r.HandleFunc("/api/faces/{id}/name", updateFaceName).Methods("PUT")
	r.HandleFunc("/api/reset", resetDatabase).Methods("POST")

	// Static file serving
	r.HandleFunc("/videos/{filename}", serveVideo).Methods("GET")
	r.HandleFunc("/faces/{filename}", serveFace).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/index.html")
	})

	// CORS middleware
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}
