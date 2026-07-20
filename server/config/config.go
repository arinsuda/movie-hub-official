package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port       string
	AppBaseURL string
	CORS       CORSConfig
	DB         DBConfig
	TMDB       TMDBConfig
	JWT        JWTConfig
	SMTP       SMTPConfig
	Cookie     CookieConfig
	MinIO      MinIOConfig
	Google     GoogleConfig
}

type GoogleConfig struct {
	Enabled            bool
	ClientID           string
	ClientSecret       string
	RedirectURL        string
	FrontendSuccessURL string
	FrontendErrorURL   string
}

type CORSConfig struct {
	AllowedOrigin string
}

type DBConfig struct {
	Host                        string
	User                        string
	Password                    string
	Name                        string
	Port                        string
	SSLMode                     string
	MigrationLockTimeoutMs      int
	MigrationStatementTimeoutMs int
}

type TMDBConfig struct {
	BaseURL string
	APIKey  string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type CookieConfig struct {
	Domain   string
	Secure   bool
	SameSite string
}

type MinIOConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

func Load() (*Config, error) {
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	migrationLockTimeoutMs, _ := strconv.Atoi(getEnv("MIGRATION_LOCK_TIMEOUT_MS", "10000"))
	migrationStatementTimeoutMs, _ := strconv.Atoi(getEnv("MIGRATION_STATEMENT_TIMEOUT_MS", "15000"))

	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),
		CORS: CORSConfig{
			AllowedOrigin: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),
		},
		DB: DBConfig{
			Host:                        requireEnv("POSTGRES_HOST"),
			User:                        requireEnv("POSTGRES_USER"),
			Password:                    requireEnv("POSTGRES_PASSWORD"),
			Name:                        requireEnv("POSTGRES_DB"),
			Port:                        getEnv("POSTGRES_PORT", "5432"),
			SSLMode:                     getEnv("POSTGRES_SSLMODE", ""),
			MigrationLockTimeoutMs:      migrationLockTimeoutMs,
			MigrationStatementTimeoutMs: migrationStatementTimeoutMs,
		},
		TMDB: TMDBConfig{
			BaseURL: requireEnv("THE_MOVIE_BASE_API"),
			APIKey:  requireEnv("THE_MOVIE_API_KEY"),
		},
		JWT: JWTConfig{
			AccessSecret:  requireEnv("JWT_ACCESS_SECRET"),
			RefreshSecret: requireEnv("JWT_REFRESH_SECRET"),
			AccessTTL:     15 * time.Minute,
			RefreshTTL:    7 * 24 * time.Hour,
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     smtpPort,
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@remov.app"),
		},
		Cookie: CookieConfig{
			Domain:   getEnv("COOKIE_DOMAIN", "localhost"),
			Secure:   getEnv("COOKIE_SECURE", "false") == "true",
			SameSite: getEnv("COOKIE_SAMESITE", "Strict"),
		},
		MinIO: MinIOConfig{
			Endpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:  getEnv("MINIO_ROOT_USER", "minioadmin"),
			SecretKey:  getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
			BucketName: getEnv("MINIO_BUCKET_NAME", "remov-private"),
			UseSSL:     getEnv("MINIO_USE_SSL", "false") == "true",
		},
		Google: GoogleConfig{
			Enabled:            getEnv("GOOGLE_OAUTH_ENABLED", "false") == "true",
			ClientID:           getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret:       getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:        getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
			FrontendSuccessURL: getEnv("GOOGLE_FRONTEND_SUCCESS_URL", "http://localhost:5173/"),
			FrontendErrorURL:   getEnv("GOOGLE_FRONTEND_ERROR_URL", "http://localhost:5173/login"),
		},
	}

	if cfg.Google.Enabled {
		if cfg.Google.ClientID == "" || cfg.Google.ClientSecret == "" || cfg.Google.RedirectURL == "" || cfg.Google.FrontendSuccessURL == "" || cfg.Google.FrontendErrorURL == "" {
			return nil, fmt.Errorf("GOOGLE_OAUTH_ENABLED is true but missing required Google OAuth configuration variables")
		}
	}

	return cfg, nil
}

func (d DBConfig) DSN() string {
	sslMode := d.SSLMode
	if sslMode == "" {
		if d.Host == "localhost" || d.Host == "127.0.0.1" {
			sslMode = "disable"
		} else {
			sslMode = "require"
		}
	}
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		d.Host, d.User, d.Password, d.Name, d.Port, sslMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}
