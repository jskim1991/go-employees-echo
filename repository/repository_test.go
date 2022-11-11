package repository

import (
	"database/sql/driver"
	"employees-echo/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestRepository(t *testing.T) {
	t.Run("FindAll returns employee slice", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		mock.ExpectQuery("select id, name, salary, age from employee").
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "salary", "age"}).
				AddRow(1, "Jay", "100", 30))

		repository := DefaultRepository{
			DB: db,
		}
		queryResult := repository.FindAll()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("error")
		}

		assert.Equal(t, 1, queryResult[0].Id)
		assert.Equal(t, "Jay", queryResult[0].Name)
		assert.Equal(t, "100", queryResult[0].Salary)
		assert.Equal(t, 30, queryResult[0].Age)
	})

	t.Run("RegisterEmployee inserts a new employee and returns id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		employeeToInsert := models.Employee{
			Name:   "Jay",
			Salary: "100",
			Age:    30,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`insert into employee (name, salary, age, created_at, updated_at) values ($1, $2, $3, $4, $5) returning id`)).
			WithArgs(employeeToInsert.Name, employeeToInsert.Salary, employeeToInsert.Age, AnyTime{}, AnyTime{}).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		repository := DefaultRepository{
			DB: db,
		}
		newId := repository.InsertEmployee(employeeToInsert)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, 1, newId)
	})

	t.Run("Update updates employee information", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		updateEmployee := models.Employee{
			Name:   "Jay",
			Salary: "100",
			Age:    30,
		}

		mock.ExpectExec(regexp.QuoteMeta(`update employee set name=$1, salary=$2, age=$3, updated_at=$4 where id = $5`)).
			WithArgs(updateEmployee.Name, updateEmployee.Salary, updateEmployee.Age, AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repository := DefaultRepository{
			DB: db,
		}
		repository.Update(1, updateEmployee)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf(err.Error())
		}
	})
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
