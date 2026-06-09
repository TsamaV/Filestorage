package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db      DbConfig
	Auth    AuthConfig
	S3      S3Config
	Server  ServerConfig
	Storage StorageConfig
}

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
	TTL    string
}

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

type ServerConfig struct {
	Port string
}

type StorageConfig struct {
	MaxFileSize string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DATABASE_URL"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("JWT_SECRET"),
			TTL:    os.Getenv("JWT_TTL"),
		},
		S3: S3Config{
			Endpoint:  os.Getenv("S3_ENDPOINT"),
			AccessKey: os.Getenv("S3_ACCESS_KEY"),
			SecretKey: os.Getenv("S3_SECRET_KEY"),
			Bucket:    os.Getenv("S3_BUCKET"),
		},
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Storage: StorageConfig{
			MaxFileSize: os.Getenv("MAX_FILE_SIZE"),
		},
	}
}
