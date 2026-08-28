package routes

import (
	"tyre-match-backend/api"
	"tyre-match-backend/api/handlers"
)

func fileRoutes(server *api.Server) {
	fileHandler := handlers.NewFileHandler(server)

	file := server.Echo.Group("/files")

	file.GET("/*filepath", fileHandler.Get)
}
