package testdoubles

import (
	"employees-echo/dto"
)

type SpyStubRepository struct {
	FindAll_returnValue []dto.EmployeeResponse
	FindAll_invocation  int

	InsertEmployee_invocation  int
	InsertEmployee_argument    dto.EmployeeResponse
	InsertEmployee_returnValue int

	Update_invocation        int
	Update_argument_id       int
	Update_argument_employee dto.EmployeeResponse
}

func (m *SpyStubRepository) FindAll() []dto.EmployeeResponse {
	m.FindAll_invocation++
	return m.FindAll_returnValue
}

func (m *SpyStubRepository) InsertEmployee(e dto.EmployeeResponse) int {
	m.InsertEmployee_argument = e
	m.InsertEmployee_invocation++
	return m.InsertEmployee_returnValue
}

func (m *SpyStubRepository) Update(id int, e dto.EmployeeResponse) {
	m.Update_invocation++
	m.Update_argument_id = id
	m.Update_argument_employee = e
}
