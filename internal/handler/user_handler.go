package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"user-api/internal/models"
	"user-api/internal/service"
)

type UserHandler struct {
	service  *service.UserService
	validate *validator.Validate
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}
	if err := h.validate.Struct(req); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		return mapServiceError(err)
	}

	return c.Status(http.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	user, err := h.service.GetUserByID(c.Context(), int32(id))
	if err != nil {
		return mapServiceError(err)
	}

	return c.JSON(user)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, limit, err := parsePagination(c)
	if err != nil {
		return err
	}

	users, err := h.service.ListUsers(c.Context(), page, limit)
	if err != nil {
		return mapServiceError(err)
	}

	return c.JSON(users)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}
	if err := h.validate.Struct(req); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.UpdateUser(c.Context(), int32(id), req)
	if err != nil {
		return mapServiceError(err)
	}

	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user id")
	}

	if err := h.service.DeleteUser(c.Context(), int32(id)); err != nil {
		return mapServiceError(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

func parsePagination(c *fiber.Ctx) (int, int, error) {
	page := 1
	limit := 10

	if raw := c.Query("page"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 1 {
			return 0, 0, fiber.NewError(http.StatusBadRequest, "invalid page")
		}
		page = value
	}

	if raw := c.Query("limit"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 1 {
			return 0, 0, fiber.NewError(http.StatusBadRequest, "invalid limit")
		}
		if value > 100 {
			value = 100
		}
		limit = value
	}

	return page, limit, nil
}

func mapServiceError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fiber.NewError(http.StatusNotFound, "user not found")
	}
	return fiber.NewError(http.StatusInternalServerError, err.Error())
}
