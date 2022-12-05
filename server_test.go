package main

import (
	"bytes"
	"database/sql"
	"employees-echo/controller"
	"employees-echo/dto"
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
		db := connectInMemoryDB(t.Name())
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

		postJson, _ := json.Marshal(dto.EmployeeRequest{
			Name:   "Jay",
			Salary: "100",
			Age:    30,
		})
		request := httptest.NewRequest(http.MethodPost, "/employee", bytes.NewReader(postJson))
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

		var returnedEmployees []dto.EmployeeResponse
		json.Unmarshal(response.Body.Bytes(), &returnedEmployees)
		assert.Equal(t, 1, len(returnedEmployees))
		assert.Equal(t, 1, returnedEmployees[0].Id)
		assert.Equal(t, "Jay", returnedEmployees[0].Name)
		assert.Equal(t, "100", returnedEmployees[0].Salary)
		assert.Equal(t, 30, returnedEmployees[0].Age)
	})

	t.Run("Insert then update and return all", func(t *testing.T) {
		db := connectInMemoryDB(t.Name())
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

		postJson, _ := json.Marshal(dto.EmployeeRequest{
			Name:   "Jay",
			Salary: "100",
			Age:    30,
		})
		request := httptest.NewRequest(http.MethodPut, "/employee", bytes.NewReader(postJson))
		request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		handler.RegisterEmployee(c)

		newId := response.Body.String()

		putJson, _ := json.Marshal(dto.EmployeeRequest{
			Name:   "Jay Kim",
			Salary: "1000",
			Age:    31,
		})
		request = httptest.NewRequest(http.MethodPut, "/employee", bytes.NewReader(putJson))
		request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response = httptest.NewRecorder()
		c = e.NewContext(request, response)
		c.SetPath(":id")
		c.SetParamNames("id")
		c.SetParamValues(newId)

		err := handler.UpdateEmployee(c)
		if err != nil {
			t.Error(err)
		}

		request = httptest.NewRequest(http.MethodGet, "/employees", nil)
		response = httptest.NewRecorder()
		c = e.NewContext(request, response)

		err = handler.GetAllEmployees(c)
		if err != nil {
			t.Error(err)
		}

		var returnedEmployees []dto.EmployeeResponse
		json.Unmarshal(response.Body.Bytes(), &returnedEmployees)
		assert.Equal(t, 1, len(returnedEmployees))
		assert.Equal(t, 1, returnedEmployees[0].Id)
		assert.Equal(t, "Jay Kim", returnedEmployees[0].Name)
		assert.Equal(t, "1000", returnedEmployees[0].Salary)
		assert.Equal(t, 31, returnedEmployees[0].Age)
	})
}

func connectInMemoryDB(datasource string) *sql.DB {
	db, err := sql.Open("ramsql", datasource)
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
