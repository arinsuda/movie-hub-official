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
}

type CORSConfig struct {
	AllowedOrigin string
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
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

	return &Config{
		Port:       getEnv("PORT", "8080"),
		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),
		CORS: CORSConfig{
			AllowedOrigin: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),
		},
		DB: DBConfig{
			Host:     requireEnv("POSTGRES_HOST"),
			User:     requireEnv("POSTGRES_USER"),
			Password: requireEnv("POSTGRES_PASSWORD"),
			Name:     requireEnv("POSTGRES_DB"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
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
			Host:     requireEnv("SMTP_HOST"),
			Port:     smtpPort,
			Username: requireEnv("SMTP_USERNAME"),
			Password: requireEnv("SMTP_PASSWORD"),
			From:     requireEnv("SMTP_FROM"),
		},
		Cookie: CookieConfig{
			Domain:   getEnv("COOKIE_DOMAIN", "localhost"),
			Secure:   getEnv("COOKIE_SECURE", "false") == "true",
			SameSite: getEnv("COOKIE_SAMESITE", "Strict"),
		},
		MinIO: MinIOConfig{
			Endpoint:   requireEnv("MINIO_ENDPOINT"),
			AccessKey:  requireEnv("MINIO_ROOT_USER"),
			SecretKey:  requireEnv("MINIO_ROOT_PASSWORD"),
			BucketName: getEnv("MINIO_BUCKET_NAME", "remov-private"),
			UseSSL:     getEnv("MINIO_USE_SSL", "false") == "true",
		},
	}, nil
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		d.Host, d.User, d.Password, d.Name, d.Port,
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
