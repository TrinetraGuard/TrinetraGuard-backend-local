package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database connection
type Database struct {
	DB *sql.DB
}

// Initialize sets up the database connection and runs migrations
func Initialize(dbType, connectionString string) (*Database, error) {
	var db *sql.DB
	var err error

	switch dbType {
	case "sqlite":
		db, err = sql.Open("sqlite3", connectionString)
	case "postgres":
		db, err = sql.Open("postgres", connectionString)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	database := &Database{DB: db}

	// Run migrations
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully")
	return database, nil
}

// migrate creates all necessary tables and indexes
func (d *Database) migrate() error {
	queries := []string{
		// Videos table
		`CREATE TABLE IF NOT EXISTS videos (
			id TEXT PRIMARY KEY,
			filename TEXT NOT NULL,
			original_name TEXT NOT NULL,
			file_path TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			duration REAL,
			frame_count INTEGER,
			width INTEGER,
			height INTEGER,
			status TEXT NOT NULL DEFAULT 'uploaded',
			uploaded_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,

		// Analysis jobs table
		`CREATE TABLE IF NOT EXISTS analysis_jobs (
			id TEXT PRIMARY KEY,
			video_id TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			progress INTEGER NOT NULL DEFAULT 0,
			started_at DATETIME NOT NULL,
			completed_at DATETIME,
			error TEXT,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
		)`,

		// Analysis results table
		`CREATE TABLE IF NOT EXISTS analysis_results (
			id TEXT PRIMARY KEY,
			video_id TEXT NOT NULL,
			analysis_job_id TEXT NOT NULL,
			total_frames INTEGER NOT NULL DEFAULT 0,
			total_people INTEGER NOT NULL DEFAULT 0,
			unique_people INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE,
			FOREIGN KEY (analysis_job_id) REFERENCES analysis_jobs (id) ON DELETE CASCADE
		)`,

		// People per frame table
		`CREATE TABLE IF NOT EXISTS people_per_frame (
			id TEXT PRIMARY KEY,
			analysis_result_id TEXT NOT NULL,
			frame_number INTEGER NOT NULL,
			timestamp REAL NOT NULL,
			count INTEGER NOT NULL,
			FOREIGN KEY (analysis_result_id) REFERENCES analysis_results (id) ON DELETE CASCADE
		)`,

		// Tracking data table
		`CREATE TABLE IF NOT EXISTS tracking_data (
			id TEXT PRIMARY KEY,
			analysis_result_id TEXT NOT NULL,
			person_id TEXT NOT NULL,
			frame_number INTEGER NOT NULL,
			timestamp REAL NOT NULL,
			x REAL NOT NULL,
			y REAL NOT NULL,
			width REAL NOT NULL,
			height REAL NOT NULL,
			confidence REAL NOT NULL,
			FOREIGN KEY (analysis_result_id) REFERENCES analysis_results (id) ON DELETE CASCADE
		)`,

		// Persons table (unique people detected)
		`CREATE TABLE IF NOT EXISTS persons (
			id TEXT PRIMARY KEY,
			video_id TEXT NOT NULL,
			person_id TEXT NOT NULL,
			first_seen REAL NOT NULL,
			last_seen REAL NOT NULL,
			total_frames INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE,
			UNIQUE(video_id, person_id)
		)`,

		// Person faces table (face images for each person)
		`CREATE TABLE IF NOT EXISTS person_faces (
			id TEXT PRIMARY KEY,
			person_id TEXT NOT NULL,
			video_id TEXT NOT NULL,
			frame_number INTEGER NOT NULL,
			timestamp REAL NOT NULL,
			face_image TEXT NOT NULL,
			x REAL NOT NULL,
			y REAL NOT NULL,
			width REAL NOT NULL,
			height REAL NOT NULL,
			confidence REAL NOT NULL,
			is_best_face BOOLEAN NOT NULL DEFAULT FALSE,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
		)`,

		// Person activity timeline table
		`CREATE TABLE IF NOT EXISTS person_activity (
			id TEXT PRIMARY KEY,
			person_id TEXT NOT NULL,
			video_id TEXT NOT NULL,
			start_time REAL NOT NULL,
			end_time REAL NOT NULL,
			start_frame INTEGER NOT NULL,
			end_frame INTEGER NOT NULL,
			duration REAL NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
		)`,

		// Reference images table
		`CREATE TABLE IF NOT EXISTS reference_images (
			id TEXT PRIMARY KEY,
			filename TEXT NOT NULL,
			original_name TEXT NOT NULL,
			file_path TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			uploaded_at DATETIME NOT NULL
		)`,

		// Search jobs table
		`CREATE TABLE IF NOT EXISTS search_jobs (
			id TEXT PRIMARY KEY,
			reference_image_id TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			progress INTEGER NOT NULL DEFAULT 0,
			started_at DATETIME NOT NULL,
			completed_at DATETIME,
			error TEXT,
			FOREIGN KEY (reference_image_id) REFERENCES reference_images (id) ON DELETE CASCADE
		)`,

		// Search results table
		`CREATE TABLE IF NOT EXISTS search_results (
			id TEXT PRIMARY KEY,
			search_job_id TEXT NOT NULL,
			video_id TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (search_job_id) REFERENCES search_jobs (id) ON DELETE CASCADE,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
		)`,

		// Search matches table
		`CREATE TABLE IF NOT EXISTS search_matches (
			id TEXT PRIMARY KEY,
			search_result_id TEXT NOT NULL,
			person_id TEXT NOT NULL,
			first_frame INTEGER NOT NULL,
			last_frame INTEGER NOT NULL,
			first_time REAL NOT NULL,
			last_time REAL NOT NULL,
			confidence REAL NOT NULL,
			x REAL NOT NULL,
			y REAL NOT NULL,
			width REAL NOT NULL,
			height REAL NOT NULL,
			FOREIGN KEY (search_result_id) REFERENCES search_results (id) ON DELETE CASCADE
		)`,

		// Indexes for better performance
		`CREATE INDEX IF NOT EXISTS idx_videos_status ON videos (status)`,
		`CREATE INDEX IF NOT EXISTS idx_videos_uploaded_at ON videos (uploaded_at)`,
		`CREATE INDEX IF NOT EXISTS idx_analysis_jobs_video_id ON analysis_jobs (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_analysis_jobs_status ON analysis_jobs (status)`,
		`CREATE INDEX IF NOT EXISTS idx_analysis_results_video_id ON analysis_results (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_people_per_frame_analysis_result_id ON people_per_frame (analysis_result_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tracking_data_analysis_result_id ON tracking_data (analysis_result_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tracking_data_person_id ON tracking_data (person_id)`,
		`CREATE INDEX IF NOT EXISTS idx_persons_video_id ON persons (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_persons_person_id ON persons (person_id)`,
		`CREATE INDEX IF NOT EXISTS idx_person_faces_person_id ON person_faces (person_id)`,
		`CREATE INDEX IF NOT EXISTS idx_person_faces_video_id ON person_faces (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_person_faces_is_best_face ON person_faces (is_best_face)`,
		`CREATE INDEX IF NOT EXISTS idx_person_activity_person_id ON person_activity (person_id)`,
		`CREATE INDEX IF NOT EXISTS idx_person_activity_video_id ON person_activity (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_reference_images_uploaded_at ON reference_images (uploaded_at)`,
		`CREATE INDEX IF NOT EXISTS idx_search_jobs_reference_image_id ON search_jobs (reference_image_id)`,
		`CREATE INDEX IF NOT EXISTS idx_search_jobs_status ON search_jobs (status)`,
		`CREATE INDEX IF NOT EXISTS idx_search_results_search_job_id ON search_results (search_job_id)`,
		`CREATE INDEX IF NOT EXISTS idx_search_results_video_id ON search_results (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_search_matches_search_result_id ON search_matches (search_result_id)`,
		`CREATE INDEX IF NOT EXISTS idx_search_matches_person_id ON search_matches (person_id)`,
	}

	for _, query := range queries {
		if _, err := d.DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query: %v", err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}
