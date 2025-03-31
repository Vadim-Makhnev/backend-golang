package handler

import (
	"project/internal/dto"
	"project/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s * service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func RegisterNewRoutes(app *fiber.App, a * AuthHandler){
	api := app.Group("api")
	
	api.Post("/register", a.CreateUser)
	api.Post("/login", a.LoginUser)
}



func (a *AuthHandler) CreateUser(c *fiber.Ctx)error{
	var user dto.UserRegister

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse a body",
		})
	}

	validate := validator.New()

	err = validate.Struct(user)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid credentials: email, password, username",
		})
	}


	err = a.service.CreateUser(user.Username, user.Email, user.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot create a user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user registered successfully",
	})
}

func (a *AuthHandler) LoginUser(c *fiber.Ctx) error{
	var user dto.UserLogin
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":"cannot parse a body",
		})
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":"invalid credentials",
		})
	}
	err = a.service.LoginUser(user.Email, user.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":"invalid credentials",
		})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message":"You login successfully",
	})
}