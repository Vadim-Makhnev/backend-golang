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


func main() {

	db := config.ConnectDB()

	err := db.AutoMigrate(&model.User{}, &model.Token{})
	if err != nil{
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
