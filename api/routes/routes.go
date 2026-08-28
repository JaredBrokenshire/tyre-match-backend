package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	s "tyre-match-backend/api"
)

func ConfigureRoutes(server *s.Server) {
	// Log requests
	server.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [${status}] ${method} ${uri}\n",
	}))

	// Configure CORS
	server.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localtyrematch.com:8080",
		},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Routes
	tyreModelRoutes(server)
	tyreImpressionRoutes(server)
	fileRoutes(server)

}
