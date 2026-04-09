package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBDSN      string
	RedisAddr  string
	JWTSecret  string
	UploadDir  string
	MaxQuotaMB int
}

func Load() (*Config, error) {
	// в проде .env может не быть, это ок
	_ = godotenv.Load()

	quota, _ := strconv.Atoi(env("MAX_QUOTA_MB", "100"))

	host := env("DB_HOST", "localhost")
	port := env("DB_PORT", "5432")
	user := env("DB_USER", "filestorage")
	pass := env("DB_PASSWORD", "filestorage")
	name := env("DB_NAME", "filestorage")

	dsn := "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + name + "?sslmode=disable"

	return &Config{
		Port:       env("SERVER_PORT", "8080"),
		DBDSN:      dsn,
		RedisAddr:  env("REDIS_ADDR", "localhost:6379"),
		JWTSecret:  env("JWT_SECRET", "dev-secret-change-me"),
		UploadDir:  env("UPLOAD_DIR", "./uploads"),
		MaxQuotaMB: quota,
	}, nil
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
