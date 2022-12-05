package controller

import (
	"employees-echo/dto"
	"employees-echo/model"
	"employees-echo/repository"
	"fmt"
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
	found := m.Repository.FindAll()
	employees := mapToResponse(found)
	return c.JSON(http.StatusOK, employees)
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
	return c.String(http.StatusCreated, fmt.Sprint(id))
}

func (m *Controller) UpdateEmployee(c echo.Context) error {
	e := &dto.EmployeeRequest{}
	err := c.Bind(e)
	if err != nil {
		log.Fatalln(err)
	}

	pathId := c.Param("id")
	employeeId, err := strconv.ParseUint(pathId, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	updateEmployee := dto.EmployeeResponse{
		Name:   e.Name,
		Salary: e.Salary,
		Age:    e.Age,
	}

	m.Repository.Update(uint(employeeId), updateEmployee)
	return c.NoContent(http.StatusOK)
}

func mapToResponse(models []model.Employee) []dto.EmployeeResponse {
	var responses []dto.EmployeeResponse
	for i := 0; i < len(models); i++ {
		e := models[i]
		r := dto.EmployeeResponse{
			Id:     e.ID,
			Name:   e.Name,
			Salary: e.Salary,
			Age:    e.Age,
		}
		responses = append(responses, r)
	}
	return responses
}
