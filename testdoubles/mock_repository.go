package testdoubles

import (
	"employees-echo/dto"
	"employees-echo/model"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindAll() []model.Employee {
	args := m.Called()
	return args.Get(0).([]model.Employee)
}

func (m *MockRepository) InsertEmployee(e dto.EmployeeRequest) uint {
	args := m.Called(e)
	return args.Get(0).(uint)
}

func (m *MockRepository) Update(id uint, e dto.EmployeeRequest) model.Employee {
	args := m.Called(id, e)
	return args.Get(0).(model.Employee)
}
