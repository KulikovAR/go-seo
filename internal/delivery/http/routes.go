package http

import (
	"go-seo/internal/delivery/http/handlers"
	"go-seo/internal/usecases"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, useCases *usecases.Container) {
	siteHandler := handlers.NewSiteHandler(useCases.Site)
	keywordHandler := handlers.NewKeywordHandler(useCases.Keyword)
	positionHandler := handlers.NewPositionHandler(useCases.PositionTracking)

	api := r.Group("/api")
	{
		sites := api.Group("/sites")
		{
			sites.POST("", siteHandler.CreateSite)
			sites.GET("", siteHandler.GetSites)
			sites.DELETE("/:id", siteHandler.DeleteSite)
		}

		keywords := api.Group("/keywords")
		{
			keywords.POST("", keywordHandler.CreateKeyword)
			keywords.GET("", keywordHandler.GetKeywords)
			keywords.DELETE("/:id", keywordHandler.DeleteKeyword)
		}

		positions := api.Group("/positions")
		{
			positions.POST("/track-site", positionHandler.TrackSitePositions)
			positions.GET("/history", positionHandler.GetPositionsHistory)
			positions.GET("/latest", positionHandler.GetLatestPositions)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
