package routes

import (
	"tyre-match-backend/api"
	"tyre-match-backend/api/handlers"
)

func tyreImpressionRoutes(server *api.Server) {
	tyreImpressionHandler := handlers.NewTyreImpressionHandler(server)

	tyreImpression := server.Echo.Group("/tyre-impressions")

	// CRUD
	tyreImpression.GET("", tyreImpressionHandler.List)
	tyreImpression.GET("/:id", tyreImpressionHandler.Get)
	tyreImpression.POST("", tyreImpressionHandler.Create)
	tyreImpression.DELETE("/:id", tyreImpressionHandler.Delete)

	// Images
	tyreImpression.POST("/:id", tyreImpressionHandler.Upload)
}
