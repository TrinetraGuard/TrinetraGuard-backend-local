package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"video-analysis-service/internal/config"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database connection
type Database struct {
	*sql.DB
}

// Initialize creates a new database connection and runs migrations
func Initialize(cfg config.DatabaseConfig) (*Database, error) {
	var db *sql.DB
	var err error

	switch cfg.Driver {
	case "sqlite3":
		// Ensure directory exists
		dir := filepath.Dir(cfg.Name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
		db, err = sql.Open("sqlite3", cfg.Name)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
		db, err = sql.Open("postgres", dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{DB: db}

	// Run migrations
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database initialized successfully")
	return database, nil
}

// migrate runs database migrations
func (db *Database) migrate() error {
	queries := []string{
		// Videos table
		`CREATE TABLE IF NOT EXISTS videos (
			id TEXT PRIMARY KEY,
			filename TEXT NOT NULL,
			original_filename TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			duration REAL,
			frame_count INTEGER,
			width INTEGER,
			height INTEGER,
			status TEXT DEFAULT 'uploaded',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Analysis jobs table
		`CREATE TABLE IF NOT EXISTS analysis_jobs (
			id TEXT PRIMARY KEY,
			video_id TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			progress INTEGER DEFAULT 0,
			error_message TEXT,
			started_at TIMESTAMP,
			completed_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
		)`,

		// Analysis results table
		`CREATE TABLE IF NOT EXISTS analysis_results (
			id TEXT PRIMARY KEY,
			job_id TEXT NOT NULL,
			video_id TEXT NOT NULL,
			total_frames INTEGER,
			total_people INTEGER,
			unique_people INTEGER,
			people_per_frame JSON,
			tracking_data JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (job_id) REFERENCES analysis_jobs (id) ON DELETE CASCADE,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
		)`,

		// Reference images table
		`CREATE TABLE IF NOT EXISTS reference_images (
			id TEXT PRIMARY KEY,
			filename TEXT NOT NULL,
			original_filename TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Search jobs table
		`CREATE TABLE IF NOT EXISTS search_jobs (
			id TEXT PRIMARY KEY,
			reference_image_id TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			progress INTEGER DEFAULT 0,
			error_message TEXT,
			started_at TIMESTAMP,
			completed_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (reference_image_id) REFERENCES reference_images (id) ON DELETE CASCADE
		)`,

		// Search results table
		`CREATE TABLE IF NOT EXISTS search_results (
			id TEXT PRIMARY KEY,
			search_job_id TEXT NOT NULL,
			video_id TEXT NOT NULL,
			reference_image_id TEXT NOT NULL,
			matches JSON,
			first_appearance REAL,
			last_appearance REAL,
			total_appearances INTEGER,
			confidence REAL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (search_job_id) REFERENCES search_jobs (id) ON DELETE CASCADE,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE,
			FOREIGN KEY (reference_image_id) REFERENCES reference_images (id) ON DELETE CASCADE
		)`,

		// Persons table for unique people detected
		`CREATE TABLE IF NOT EXISTS persons (
			id TEXT PRIMARY KEY,
			video_id TEXT NOT NULL,
			person_number INTEGER NOT NULL,
			first_frame INTEGER,
			last_frame INTEGER,
			first_time REAL,
			last_time REAL,
			total_frames INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
		)`,

		// Person faces table for storing face images
		`CREATE TABLE IF NOT EXISTS person_faces (
			id TEXT PRIMARY KEY,
			person_id TEXT NOT NULL,
			video_id TEXT NOT NULL,
			frame_number INTEGER,
			timestamp REAL,
			bounding_box_x REAL,
			bounding_box_y REAL,
			bounding_box_width REAL,
			bounding_box_height REAL,
			confidence REAL,
			face_image TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE,
			FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
		)`,

		// Create indexes for better performance
		`CREATE INDEX IF NOT EXISTS idx_videos_status ON videos (status)`,
		`CREATE INDEX IF NOT EXISTS idx_analysis_jobs_video_id ON analysis_jobs (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_analysis_jobs_status ON analysis_jobs (status)`,
		`CREATE INDEX IF NOT EXISTS idx_analysis_results_video_id ON analysis_results (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_search_jobs_status ON search_jobs (status)`,
		`CREATE INDEX IF NOT EXISTS idx_search_results_video_id ON search_results (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_search_results_reference_image_id ON search_results (reference_image_id)`,
		`CREATE INDEX IF NOT EXISTS idx_persons_video_id ON persons (video_id)`,
		`CREATE INDEX IF NOT EXISTS idx_person_faces_person_id ON person_faces (person_id)`,
		`CREATE INDEX IF NOT EXISTS idx_person_faces_video_id ON person_faces (video_id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query: %w", err)
		}
	}

	return nil
}
