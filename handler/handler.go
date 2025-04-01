package handler

import (
	"errors"
	"fmt"
	"project/internal/dto"
	"project/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandlerI interface {
    CreateUser(c *fiber.Ctx)error
    LoginUser(c *fiber.Ctx)error
}

type AuthHandler struct {
	service service.AuthServiceI
}

func NewAuthHandler(s service.AuthServiceI) AuthHandlerI {
	return &AuthHandler{
		service: s,
	}
}

func RegisterNewRoutes(app *fiber.App, a AuthHandlerI){
	api := app.Group("api")
	
	api.Post("/register", a.CreateUser)
	api.Post("/login", a.LoginUser)
}



func (a *AuthHandler) CreateUser(c *fiber.Ctx) error {
    var user dto.UserRegister

    // Парсинг тела запроса
    err := c.BodyParser(&user)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
            "error":   err.Error(),
        })
    }

    // Валидация данных
    validate := validator.New()
    err = validate.Struct(user)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Validation failed",
            "error":   err.Error(),
        })
    }

    // Вызов сервиса
    err = a.service.CreateUser(user.Username, user.Email, user.Password)
    if err != nil {
        // Проверяем тип ошибки
        if errors.Is(err, fmt.Errorf("user already exists")) {
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "message": "User with this email already exists",
            })
        }

        // Общая ошибка создания пользователя
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to create user",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "User registered successfully",
    })
}

func (a *AuthHandler) LoginUser(c *fiber.Ctx) error {
    var user dto.UserLogin

    // Парсинг тела запроса
    err := c.BodyParser(&user)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
            "error":   err.Error(),
        })
    }

    // Валидация данных
    validate := validator.New()
    err = validate.Struct(user)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Validation failed",
            "error":   err.Error(),
        })
    }

    // Вызов сервиса
    err = a.service.LoginUser(c, user.Email, user.Password)
    if err != nil {
        if errors.Is(err, fmt.Errorf("user is not found")) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid credentials",
            })
        }

        // Общая ошибка входа
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to login",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Logged in successfully",
    })
}