package repository

import (
	"database/sql"
	"employees-echo/dto"
	"employees-echo/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
)

type Repository interface {
	FindAll() []model.Employee
	InsertEmployee(e dto.EmployeeResponse) uint
	Update(id uint, e dto.EmployeeResponse) model.Employee
}

type DefaultRepository struct {
	DB *gorm.DB
}

func ConnectDB(conn *sql.DB) *gorm.DB {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: conn}), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalln(err)
	}
	return gormDB
}

func (m *DefaultRepository) FindAll() []model.Employee {
	var employees []model.Employee
	m.DB.Find(&employees)
	return employees
}

func (m *DefaultRepository) InsertEmployee(e dto.EmployeeResponse) uint {
	newEmployee := model.Employee{
		Name:   e.Name,
		Salary: e.Salary,
		Age:    e.Age,
	}
	m.DB.Create(&newEmployee)
	return newEmployee.ID
}

func (m *DefaultRepository) Update(employeeId uint, e dto.EmployeeResponse) model.Employee {
	employee := model.Employee{
		Model:  gorm.Model{ID: employeeId},
		Name:   e.Name,
		Salary: e.Salary,
		Age:    e.Age,
	}
	employee.Model.ID = employeeId
	m.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "salary", "age"}),
	}).Create(&employee)

	return employee
}
