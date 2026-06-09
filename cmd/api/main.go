package main

import (
	"context"
	config "filestorage/config"
	"filestorage/internal/auth"
	"filestorage/internal/db"
	"filestorage/internal/file"
	"filestorage/internal/user"
	pkg "filestorage/pkg/logger"
	"filestorage/pkg/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	config := config.LoadConfig()
	ctx := context.Background()
	s3Client, err := minio.New(config.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3.AccessKey, config.S3.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	db, err := db.CreateConnection(ctx, config.Db.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	logger, cleanup, err := pkg.New("debug")
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	validate := validator.New()

	// Repositories
	userRepo := user.NewUserRepository(db.DB)
	fileRepo := file.NewFileRepository(db.DB)

	// Services
	authService := auth.NewAuthService(userRepo, config.Auth.Secret)
	fileService := file.NewFileService(fileRepo, s3Client, config.S3.Bucket)

	// Handlers
	authHandler := auth.NewHandler(authService, logger, validate)
	fileHandler := file.NewHandler(fileService, logger)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(config.Auth.Secret, logger)

	// Server
	mux := http.NewServeMux()
	authHandler.Routes(mux)
	fileHandler.Routes(mux, authMiddleware)

	server := &http.Server{
		Addr: ":" + config.Server.Port,
		Handler: mux,
	}

	go func() {
		log.Println("Server listenin on :" + config.Server.Port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<- quit

	log.Println("shutting down...")
	server.Shutdown(context.Background())
	log.Println("server stopped")
}
