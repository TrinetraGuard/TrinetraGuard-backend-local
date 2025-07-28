package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

type VideoAnalysis struct {
	ID          int       `json:"id"`
	Filename    string    `json:"filename"`
	UploadTime  time.Time `json:"upload_time"`
	TotalPeople int       `json:"total_people"`
	Faces       []Face    `json:"faces"`
}

type Face struct {
	ID         int      `json:"id"`
	VideoID    int      `json:"video_id"`
	ImagePath  string   `json:"image_path"`
	Timestamps []string `json:"timestamps"`
}

type PythonAnalysisResult struct {
	TotalPeople int `json:"total_people"`
	Faces       []struct {
		Image      string   `json:"image"`
		Timestamps []string `json:"timestamps"`
	} `json:"faces"`
}

var db *sql.DB

func main() {
	// Initialize database
	initDatabase()

	// Create directories if they don't exist
	os.MkdirAll("videos", 0755)
	os.MkdirAll("faces", 0755)

	// Setup router
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/upload", uploadVideo).Methods("POST")
	r.HandleFunc("/api/videos", getVideos).Methods("GET")
	r.HandleFunc("/api/videos/{id}", getVideoAnalysis).Methods("GET")
	r.HandleFunc("/videos/{filename}", serveVideo).Methods("GET")
	r.HandleFunc("/faces/{filename}", serveFace).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend")))

	// CORS middleware
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	createTables := `
	CREATE TABLE IF NOT EXISTS videos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT NOT NULL,
		upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		total_people INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS faces (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		video_id INTEGER,
		image_path TEXT NOT NULL,
		timestamps TEXT,
		FOREIGN KEY (video_id) REFERENCES videos (id)
	);
	`

	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	}
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "No video file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
	filepath := filepath.Join("videos", filename)

	// Save video file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert into database
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

	// Run Python analysis
	go func() {
		analyzeVideo(int(videoID), filepath)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Video uploaded successfully",
		"video_id": videoID,
		"filename": filename,
	})
}

func analyzeVideo(videoID int, videoPath string) {
	// Get absolute paths
	absVideoPath, _ := filepath.Abs(videoPath)
	absFacesDir, _ := filepath.Abs("faces")
	absPythonDir, _ := filepath.Abs("../python")

	// Run Python script with virtual environment
	cmd := exec.Command("python3", "analyze_video_final.py", absVideoPath, absFacesDir)
	cmd.Dir = absPythonDir
	// Set the Python path to use the virtual environment
	cmd.Env = append(os.Environ(), "PATH="+absPythonDir+"/venv/bin:"+os.Getenv("PATH"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		log.Printf("Python analysis failed: %v", err)
		return
	}

	// Parse Python output - find the JSON in the output
	outputStr := string(output)
	jsonStart := strings.Index(outputStr, "{")
	if jsonStart == -1 {
		log.Printf("No JSON found in Python output")
		return
	}

	jsonStr := outputStr[jsonStart:]
	var result PythonAnalysisResult
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		log.Printf("Failed to parse Python output: %v", err)
		log.Printf("Output was: %s", jsonStr)
		return
	}

	// Update database with results
	_, err = db.Exec("UPDATE videos SET total_people = ? WHERE id = ?", result.TotalPeople, videoID)
	if err != nil {
		log.Printf("Failed to update video: %v", err)
		return
	}

	// Insert faces
	for _, face := range result.Faces {
		timestampsJSON, _ := json.Marshal(face.Timestamps)
		_, err = db.Exec("INSERT INTO faces (video_id, image_path, timestamps) VALUES (?, ?, ?)",
			videoID, face.Image, string(timestampsJSON))
		if err != nil {
			log.Printf("Failed to insert face: %v", err)
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

	var videos []VideoAnalysis
	for rows.Next() {
		var video VideoAnalysis
		err := rows.Scan(&video.ID, &video.Filename, &video.UploadTime, &video.TotalPeople)
		if err != nil {
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
	var video VideoAnalysis
	err = db.QueryRow("SELECT id, filename, upload_time, total_people FROM videos WHERE id = ?", videoID).
		Scan(&video.ID, &video.Filename, &video.UploadTime, &video.TotalPeople)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	// Get faces
	rows, err := db.Query("SELECT id, image_path, timestamps FROM faces WHERE video_id = ?", videoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var face Face
		var timestampsJSON string
		err := rows.Scan(&face.ID, &face.ImagePath, &timestampsJSON)
		if err != nil {
			continue
		}
		json.Unmarshal([]byte(timestampsJSON), &face.Timestamps)
		face.VideoID = videoID
		video.Faces = append(video.Faces, face)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

func serveVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	filepath := filepath.Join("videos", filename)
	http.ServeFile(w, r, filepath)
}

func serveFace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	filepath := filepath.Join("faces", filename)
	http.ServeFile(w, r, filepath)
}
