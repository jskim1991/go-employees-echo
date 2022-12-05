package controller

import (
	"employees-echo/dto"
	"employees-echo/repository"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
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

func (m *Controller) RegisterEmployee(c echo.Context) error {
	e := &dto.EmployeeRequest{}
	err := c.Bind(e)
	if err != nil {
		log.Fatalln(err)
	}

	newEmployee := dto.EmployeeResponse{
		Name:   e.Name,
		Salary: e.Salary,
		Age:    e.Age,
	}

	id := m.Repository.InsertEmployee(newEmployee)
	return c.String(http.StatusCreated, strconv.Itoa(id))
}

func (m *Controller) UpdateEmployee(c echo.Context) error {
	e := &dto.EmployeeRequest{}
	err := c.Bind(e)
	if err != nil {
		log.Fatalln(err)
	}

	path := c.Param("id")
	id, err := strconv.Atoi(path)
	if err != nil {
		log.Fatalln(err)
	}

	updateEmployee := dto.EmployeeResponse{
		Name:   e.Name,
		Salary: e.Salary,
		Age:    e.Age,
	}

	m.Repository.Update(id, updateEmployee)
	return c.NoContent(http.StatusOK)
}
