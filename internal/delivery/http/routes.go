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
	groupHandler := handlers.NewGroupHandler(useCases.Group)
	positionHandler := handlers.NewPositionHandler(useCases.PositionTracking, useCases.AsyncPositionTracking)
	trackingJobHandler := handlers.NewTrackingJobHandler(useCases.TrackingJob)
	debugHandler := handlers.NewDebugHandler(useCases.Debug)

	api := r.Group("/api")
	{
		sites := api.Group("/sites")
		{
			sites.POST("", siteHandler.CreateSite)
			sites.GET("", siteHandler.GetSites)
			sites.DELETE("/:id", siteHandler.DeleteSite)
		}

		groups := api.Group("/groups")
		{
			groups.POST("", groupHandler.CreateGroup)
			groups.GET("", groupHandler.GetGroups)
			groups.PUT("/:id", groupHandler.UpdateGroup)
			groups.DELETE("/:id", groupHandler.DeleteGroup)
		}

		keywords := api.Group("/keywords")
		{
			keywords.POST("", keywordHandler.CreateKeyword)
			keywords.POST("/batch", keywordHandler.CreateKeywordsBatch)
			keywords.GET("", keywordHandler.GetKeywords)
			keywords.PUT("/:id", keywordHandler.UpdateKeyword)
			keywords.DELETE("/:id", keywordHandler.DeleteKeyword)
		}

		positions := api.Group("/positions")
		{
			positions.POST("/track-google", positionHandler.TrackGooglePositions)
			positions.POST("/track-yandex", positionHandler.TrackYandexPositions)
			positions.POST("/track-wordstat", positionHandler.TrackWordstatPositions)
			positions.GET("/history", positionHandler.GetPositionsHistory)
			positions.GET("/latest", positionHandler.GetLatestPositions)
			positions.POST("/statistics", positionHandler.GetPositionStatistics)
			positions.GET("/combined", positionHandler.GetCombinedPositions)
		}

		trackingJobs := api.Group("/tracking-jobs")
		{
			trackingJobs.GET("", trackingJobHandler.GetTrackingJobs)
		}

		debug := api.Group("/debug")
		{
			debug.POST("/kafka/job-status", debugHandler.SendKafkaJobStatus)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
