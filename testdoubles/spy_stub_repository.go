package testdoubles

import "employees-echo/models"

type SpyStubRepository struct {
	FindAll_returnValue []models.Employee
	FindAll_invocation  int

	InsertEmployee_invocation  int
	InsertEmployee_argument    models.Employee
	InsertEmployee_returnValue int
}

func (m *SpyStubRepository) FindAll() []models.Employee {
	m.FindAll_invocation++
	return m.FindAll_returnValue
}

func (m *SpyStubRepository) InsertEmployee(e models.Employee) int {
	m.InsertEmployee_argument = e
	m.InsertEmployee_invocation++
	return m.InsertEmployee_returnValue
}
