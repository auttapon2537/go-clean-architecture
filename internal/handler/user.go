package handler

import (
	"fmt"
	"strconv"

	"github.com/example/go-clean-architecture/internal/entity"
	"github.com/example/go-clean-architecture/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

// UserHandler represents the HTTP handler for user
type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// CreateHandler handles the creation of a new user
func (h *UserHandler) CreateHandler(c *fiber.Ctx) error {
	var req entity.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := h.userUsecase.CreateUser(req)
	if err != nil {
		// Check if it's a specific error type
		switch err.(type) {
		case *usecase.EmailAlreadyExistsError:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetByIDHandler handles retrieving a user by ID
func (h *UserHandler) GetByIDHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	response, err := h.userUsecase.GetUserByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetByEmailHandler handles retrieving a user by email
func (h *UserHandler) GetByEmailHandler(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email parameter is required"})
	}

	response, err := h.userUsecase.GetUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetAllHandler handles retrieving all users
func (h *UserHandler) GetAllHandler(c *fiber.Ctx) error {
	responses, err := h.userUsecase.GetAllUsers()
	fmt.Println(responses)
	fmt.Println("debug")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// UpdateHandler handles updating a user
func (h *UserHandler) UpdateHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req entity.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := h.userUsecase.UpdateUser(uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteHandler handles deleting a user
func (h *UserHandler) DeleteHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.userUsecase.DeleteUser(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}
