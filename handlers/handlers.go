package handlers

import (
	"employees-echo/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	Repository repository.Repository
}

func (h *Handler) NewHandler(r repository.Repository) {
	h.Repository = r
}

func (h *Handler) GetAllEmployees(c echo.Context) error {

	result := h.Repository.FindAll()

	return c.JSON(http.StatusOK, result)
}
