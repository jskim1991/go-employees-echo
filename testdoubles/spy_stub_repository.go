package testdoubles

import (
	"employees-echo/dto"
	"employees-echo/model"
)

type SpyStubRepository struct {
	FindAllReturnValue []model.Employee
	FindAllInvocation  int

	InsertEmployeeInvocation  int
	InsertEmployeeArgument    dto.EmployeeRequest
	InsertEmployeeReturnValue uint

	UpdateInvocation       int
	UpdateArgumentId       uint
	UpdateArgumentEmployee dto.EmployeeRequest
	UpdateReturnValue      model.Employee
}

func (m *SpyStubRepository) FindAll() []model.Employee {
	m.FindAllInvocation++
	return m.FindAllReturnValue
}

func (m *SpyStubRepository) InsertEmployee(e dto.EmployeeRequest) uint {
	m.InsertEmployeeArgument = e
	m.InsertEmployeeInvocation++
	return m.InsertEmployeeReturnValue
}

func (m *SpyStubRepository) Update(id uint, e dto.EmployeeRequest) model.Employee {
	m.UpdateInvocation++
	m.UpdateArgumentId = id
	m.UpdateArgumentEmployee = e
	return m.UpdateReturnValue
}
