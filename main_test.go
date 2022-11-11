package main

import (
	"bytes"
	"database/sql"
	"employees-echo/controller"
	"employees-echo/models"
	"employees-echo/repository"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/proullon/ramsql/driver"
)

func TestIntegration(t *testing.T) {
	t.Run("Register an employee and return all", func(t *testing.T) {
		db := connectInMemoryDB()
		defer db.Close()

		defaultRepository := repository.DefaultRepository{
			DB: db,
		}
		handler := controller.Controller{
			Repository: &defaultRepository,
		}

		e := echo.New()
		e.Use(middleware.Logger())

		e.GET("/employees", handler.GetAllEmployees)
		e.POST("/employee", handler.RegisterEmployee)

		postJson, _ := json.Marshal(models.NewEmployeeRequest{
			Name:   "Jay",
			Salary: "100",
			Age:    30,
		})
		request := httptest.NewRequest(http.MethodGet, "/employees", bytes.NewReader(postJson))
		request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		err := handler.RegisterEmployee(c)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "1", response.Body.String())

		request = httptest.NewRequest(http.MethodGet, "/employees", nil)
		response = httptest.NewRecorder()
		c = e.NewContext(request, response)

		err = handler.GetAllEmployees(c)
		if err != nil {
			t.Error(err)
		}

		var returnedEmployees []models.Employee
		json.Unmarshal([]byte(response.Body.String()), &returnedEmployees)
		assert.Equal(t, 1, len(returnedEmployees))
		assert.Equal(t, "Jay", returnedEmployees[0].Name)
		assert.Equal(t, "100", returnedEmployees[0].Salary)
		assert.Equal(t, 30, returnedEmployees[0].Age)
	})
}

func connectInMemoryDB() *sql.DB {
	db, err := sql.Open("ramsql", "TestIntegration")
	if err != nil {
		log.Fatalln(err)
	}

	setupInitialData(db)

	return db
}

func setupInitialData(db *sql.DB) {
	batch := []string{
		`CREATE TABLE employee (
			id BIGSERIAL PRIMARY KEY,
			"name" TEXT,
			salary TEXT,
			age INT,
			created_at timestamp NOT NULL,
			updated_at timestamp NOT NULL);`,
	}

	for _, query := range batch {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
