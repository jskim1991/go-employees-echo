package repository

import (
	"employees-echo/dto"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"regexp"
	"testing"
)

func initializeGormForTest(t *testing.T) (db *gorm.DB, mock sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatal(err)
	}
	return gormDB, mock
}

func TestRepository(t *testing.T) {
	t.Run("FindAll returns employee slice", func(t *testing.T) {
		db, mock := initializeGormForTest(t)

		statement := "SELECT * FROM `employees` WHERE `employees`.`deleted_at` IS NULL"
		mock.ExpectQuery(regexp.QuoteMeta(statement)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "salary", "age"}).
				AddRow(1, "Jay", "100", 30))

		repository := DefaultRepository{
			DB: db,
		}
		queryResult := repository.FindAll()

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Equal(t, uint(1), queryResult[0].Model.ID)
		assert.Equal(t, "Jay", queryResult[0].Name)
		assert.Equal(t, "100", queryResult[0].Salary)
		assert.Equal(t, 30, queryResult[0].Age)
	})

	t.Run("RegisterEmployee inserts a new employee and returns id", func(t *testing.T) {
		db, mock := initializeGormForTest(t)

		e := dto.EmployeeRequest{
			Name:   "Jay",
			Salary: "100",
			Age:    30,
		}

		statement := "INSERT INTO `employees` (`created_at`,`updated_at`,`deleted_at`,`name`,`salary`,`age`) VALUES (?,?,?,?,?,?)"
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(statement)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, e.Name, e.Salary, e.Age).
			WillReturnResult(sqlmock.NewResult(19, 1))
		mock.ExpectCommit()

		repository := DefaultRepository{
			DB: db,
		}
		newId := repository.InsertEmployee(e)

		assert.Equal(t, nil, mock.ExpectationsWereMet())
		assert.Equal(t, uint(19), newId)
	})

	t.Run("Update updates employee information", func(t *testing.T) {
		db, mock := initializeGormForTest(t)

		e := dto.EmployeeRequest{
			Name:   "Jay",
			Salary: "100",
			Age:    30,
		}

		statement := "INSERT INTO `employees` (`created_at`,`updated_at`,`deleted_at`,`name`,`salary`,`age`,`id`) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`salary`=VALUES(`salary`),`age`=VALUES(`age`)"
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(statement)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, e.Name, e.Salary, e.Age, 199).
			WillReturnResult(sqlmock.NewResult(199, 1))
		mock.ExpectCommit()

		repository := DefaultRepository{
			DB: db,
		}
		updatedResult := repository.Update(199, e)

		assert.Equal(t, nil, mock.ExpectationsWereMet())
		assert.Equal(t, uint(199), updatedResult.ID)
		assert.Equal(t, "Jay", updatedResult.Name)
		assert.Equal(t, "100", updatedResult.Salary)
		assert.Equal(t, 30, updatedResult.Age)
	})
}
