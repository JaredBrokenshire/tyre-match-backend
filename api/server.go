package api

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"time"
	"tyre-match-backend/db"
	"tyre-match-backend/db/repositories"
	"tyre-match-backend/pkg/dependencies"
	"tyre-match-backend/services"
)

type Server struct {
	Echo         *echo.Echo
	Db           *gorm.DB
	Repos        *repositories.Repos
	Dependencies *dependencies.DependencyService
	Services     *services.Services
}

func NewServer() *Server {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = utc

	s := &Server{
		Echo: echo.New(),
		Db:   db.Init(),
	}

	s.Echo.HideBanner = true

	s.Repos = repositories.NewRepos(s.Db)
	s.Dependencies = dependencies.NewDependencyService(s.Db)

	s.Services = services.NewServices(
		s.Repos,
		s.Dependencies,
		// CV Processor
	)

	return s
}

// Start runs anything needed for the application before starting the server
func (s *Server) Start(addr string) error {
	return s.Echo.Start(":" + addr)
}
