package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
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
}
