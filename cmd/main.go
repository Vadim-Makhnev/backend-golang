package main

import (
	"fmt"
	"log"
	"project/config"
	"project/handler"
	"project/internal/model"
	"project/internal/repository"
	"project/internal/service"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title Project API
// @version 1.0
// @description This is a sample server for a project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@project.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"
func main() {

	db := config.ConnectDB()

	err := db.AutoMigrate(&model.User{}, &model.Token{})
	if err != nil {
		fmt.Printf("migration failed %s", err)
	}

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	repo := repository.NewUserRepository(db)
	service := service.NewAuthService(repo)
	handle := handler.NewAuthHandler(service)

	handler.RegisterNewRoutes(app, handle)

	log.Fatal(app.Listen(":8080"))
}
