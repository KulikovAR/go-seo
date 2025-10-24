package main

import (
	"log"
	"net/http"
	"time"

	_ "go-seo/docs"

	httpDelivery "go-seo/internal/delivery/http"
	"go-seo/internal/infrastructure/config"
	"go-seo/internal/infrastructure/database/postgres"
	"go-seo/internal/infrastructure/services"
	"go-seo/internal/repositories"
	"go-seo/internal/usecases"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := postgres.NewDatabase(postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	repos := repositories.NewContainer(db.DB)

	xmlRiverService, err := services.NewXMLRiverService(
		cfg.XMLRiver.BaseURL,
		cfg.XMLRiver.UserID,
		cfg.XMLRiver.APIKey,
	)
	if err != nil {
		log.Fatal("Failed to create XMLRiver service:", err)
	}
	defer xmlRiverService.Close()

	xmlStockService, err := services.NewXMLRiverService(
		cfg.XMLStock.BaseURL,
		cfg.XMLStock.UserID,
		cfg.XMLStock.APIKey,
	)
	if err != nil {
		log.Fatal("Failed to create XMLStock service:", err)
	}
	defer xmlStockService.Close()

	wordstatService, err := services.NewWordstatService(
		cfg.XMLRiver.BaseURL,
		cfg.XMLRiver.UserID,
		cfg.XMLRiver.APIKey,
	)
	if err != nil {
		log.Fatal("Failed to create Wordstat service:", err)
	}
	defer wordstatService.Close()

	kafkaService, err := services.NewKafkaService(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal("Failed to create Kafka service:", err)
	}
	defer kafkaService.Close()

	idGenerator := services.NewIDGeneratorService()
	retryService := services.NewRetryService(5, 10*time.Second)

	useCases := usecases.NewContainer(repos, xmlRiverService, xmlStockService, wordstatService, kafkaService, idGenerator, retryService, cfg.Async.WorkerCount, cfg.Async.BatchSize)

	r := gin.Default()

	if len(cfg.Server.TrustedProxies) > 0 {
		r.SetTrustedProxies(cfg.Server.TrustedProxies)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	httpDelivery.SetupRoutes(r, useCases)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}
