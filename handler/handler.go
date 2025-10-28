package handler

import (
	"errors"
	"fmt"
	"project/internal/dto"
	"project/internal/service"
	"project/middleware"

	_ "project/cmd/docs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type AuthHandlerI interface {
	CreateUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
	Protected(c *fiber.Ctx) error
}

type AuthHandler struct {
	service service.AuthServiceI
}

func NewAuthHandler(s service.AuthServiceI) AuthHandlerI {
	return &AuthHandler{
		service: s,
	}
}

func RegisterNewRoutes(app *fiber.App, a AuthHandlerI) {
	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("api")
	pages := app.Group("pages")

	api.Post("/register", a.CreateUser)
	api.Post("/login", a.LoginUser)

	pages.Use(middleware.JWTMiddleware())
	pages.Get("main", a.Protected)

}

// CreateUser creates a new user
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserRegister true "User registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/register [post]
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

// LoginUser authenticates a user
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserLogin true "User login data"
// @Success 200 {object} map[string]interface{} "Logged in successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/login [post]
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

// Protected is a protected route that requires JWT authentication
// @Summary Protected route
// @Description Example of a protected route that requires JWT token
// @Tags protected
// @Produce plain
// @Security BearerAuth
// @Success 200 {string} string "Protected route"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /pages/main [get]
func (a *AuthHandler) Protected(c *fiber.Ctx) error {
	return c.SendString("Protected route")
}
