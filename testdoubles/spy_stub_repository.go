package testdoubles

import "employees-echo/models"

type SpyStubRepository struct {
	FindAll_returnValue []models.Employee
	FindAll_invocation  int
}

func (m *SpyStubRepository) FindAll() []models.Employee {
	m.FindAll_invocation++
	return m.FindAll_returnValue
}
