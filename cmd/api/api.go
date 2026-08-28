package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"log"
	"os"
	version "tyre-match-backend"
	"tyre-match-backend/api"
	"tyre-match-backend/api/routes"
	"tyre-match-backend/docs"
	"tyre-match-backend/pkg/validation"
)

// @title TyreMatch Backend
// @version 0.0.1
// @description MSc Forensic Investigation dissertation project reviewing computer vision approaches for matching tyre impressions from crime scenes to a known dataset of tyre makes/models

// @contact.name Jared Brokenshire
// @contact.email jbrokenshire0306@gmail.com

// @BasePath /
func main() {
	version.Get()
	log.Printf("Starting TyreMatch API. v%v\n", version.Get())
	
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	if os.Getenv("ENVIRONMENT") == "development" {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("EXPOSE_PORT"))
	}

	app := api.NewServer()
	app.Echo.Validator = validation.NewCustomValidator(validator.New())

	routes.ConfigureRoutes(app)
	err = app.Start(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Port already used")
	}
}
