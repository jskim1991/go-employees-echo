package repository

import "employees-echo/models"

type Repository interface {
	FindAll() []models.Employee
}

type DefaultRepository struct {
}

func (m *DefaultRepository) FindAll() []models.Employee {
	var result []models.Employee

	employee := models.Employee{
		Id:     "1",
		Name:   "Jay",
		Salary: "100",
		Age:    "30",
	}
	result = append(result, employee)

	return result
}
