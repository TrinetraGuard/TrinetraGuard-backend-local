package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Storage     StorageConfig
	Analysis    AnalysisConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// StorageConfig holds file storage configuration
type StorageConfig struct {
	VideosDir    string
	FinderDir    string
	MaxFileSize  int64
	AllowedTypes []string
}

// AnalysisConfig holds analysis service configuration
type AnalysisConfig struct {
	MaxConcurrentJobs int
	JobTimeout        int
	FrameRate         int
	Confidence        float64
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite3"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "video_analysis"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Storage: StorageConfig{
			VideosDir:    getEnv("STORAGE_VIDEOS_DIR", "./videos"),
			FinderDir:    getEnv("STORAGE_FINDER_DIR", "./finder"),
			MaxFileSize:  getEnvAsInt64("STORAGE_MAX_FILE_SIZE", 100*1024*1024), // 100MB
			AllowedTypes: []string{".mp4", ".avi", ".mov", ".mkv", ".wmv", ".flv", ".webm"},
		},
		Analysis: AnalysisConfig{
			MaxConcurrentJobs: getEnvAsInt("ANALYSIS_MAX_CONCURRENT_JOBS", 3),
			JobTimeout:        getEnvAsInt("ANALYSIS_JOB_TIMEOUT", 3600), // 1 hour
			FrameRate:         getEnvAsInt("ANALYSIS_FRAME_RATE", 1),     // 1 frame per second
			Confidence:        getEnvAsFloat64("ANALYSIS_CONFIDENCE", 0.7),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsInt64 gets an environment variable as int64 or returns a default value
func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsFloat64 gets an environment variable as float64 or returns a default value
func getEnvAsFloat64(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
