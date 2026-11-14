package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"stream-rec-core/config"

	storage "stream-rec-core/internal/adapter/mongo"

	httptransport "stream-rec-core/internal/controller/http"
	"stream-rec-core/internal/usecase"

	_ "stream-rec-core/docs"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title           StreamRec Core API
// @version         1.0
// @description     NoSQL Data Layer & Backend API for a Streaming Recommendation Platform (Assignment #6).
// @host            localhost:8080
// @BasePath                   /api
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type "Bearer " followed by your JWT token.
func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongodb")
	}

	db := client.Database(cfg.MongoDB)
	repo, err := storage.NewMongoRepository(ctx, db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize repository")
	}

	svc := usecase.NewPlatformService(repo, cfg.JWTSecret)
	handler := httptransport.NewHandler(svc, log)

	router := handler.InitRouter()

	log.Info().Str("env", cfg.AppEnv).Str("port", cfg.HTTPPort).Msg("starting stream-rec-core server")
	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}
