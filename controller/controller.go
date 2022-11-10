package controller

import (
	"employees-echo/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	Repository repository.Repository
}

func (m *Controller) NewHandler(r repository.Repository) {
	m.Repository = r
}

func (m *Controller) GetAllEmployees(c echo.Context) error {

	result := m.Repository.FindAll()

	return c.JSON(http.StatusOK, result)
}
