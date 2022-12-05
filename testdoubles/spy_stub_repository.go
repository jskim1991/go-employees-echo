package testdoubles

import (
	"employees-echo/dto"
	"employees-echo/model"
)

type SpyStubRepository struct {
	FindAllReturnValue []model.Employee
	FindAllInvocation  int

	InsertEmployeeInvocation  int
	InsertEmployeeArgument    dto.EmployeeResponse
	InsertEmployeeReturnValue uint

	UpdateInvocation       int
	UpdateArgumentId       uint
	UpdateArgumentEmployee dto.EmployeeResponse
	UpdateReturnValue      model.Employee
}

func (m *SpyStubRepository) FindAll() []model.Employee {
	m.FindAllInvocation++
	return m.FindAllReturnValue
}

func (m *SpyStubRepository) InsertEmployee(e dto.EmployeeResponse) uint {
	m.InsertEmployeeArgument = e
	m.InsertEmployeeInvocation++
	return m.InsertEmployeeReturnValue
}

func (m *SpyStubRepository) Update(id uint, e dto.EmployeeResponse) model.Employee {
	m.UpdateInvocation++
	m.UpdateArgumentId = id
	m.UpdateArgumentEmployee = e
	return m.UpdateReturnValue
}
