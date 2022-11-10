package repository

import "employees-echo/models"

type SpyRepository struct {
	FindAll_returnValue []models.Employee
}

func (m *SpyRepository) FindAll() []models.Employee {
	return m.FindAll_returnValue
}
