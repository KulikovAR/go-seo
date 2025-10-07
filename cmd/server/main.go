package main

import (
	"log"

	_ "go-seo/docs"

	"go-seo/internal/delivery/http"
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

	useCases := usecases.NewContainer(repos, xmlRiverService)

	r := gin.Default()

	if len(cfg.Server.TrustedProxies) > 0 {
		r.SetTrustedProxies(cfg.Server.TrustedProxies)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	http.SetupRoutes(r, useCases)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
