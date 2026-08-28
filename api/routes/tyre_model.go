package routes

import (
	"tyre-match-backend/api"
	"tyre-match-backend/api/handlers"
)

func tyreModelRoutes(server *api.Server) {
	tyreModelHandler := handlers.NewTyreModelHandler(server)

	tyreModel := server.Echo.Group("/tyre-models")

	// CRUD
	tyreModel.GET("", tyreModelHandler.List)
	tyreModel.GET("/:id", tyreModelHandler.Get)
	tyreModel.POST("", tyreModelHandler.Create)
	tyreModel.PUT("/:id", tyreModelHandler.Update)
	tyreModel.DELETE("/:id", tyreModelHandler.Delete)

	// Images
	tyreModel.POST("/:id", tyreModelHandler.Upload)
}
