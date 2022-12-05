package main

import (
	"bytes"
	"database/sql"
	"employees-echo/controller"
	"employees-echo/dto"
	"employees-echo/model"
	"employees-echo/repository"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	DatabaseName = "journey"
	os.Exit(m.Run())
}

func run(t *testing.T) (e *echo.Echo, handler controller.Controller) {
	t.Helper()
	sqlDB, err := sql.Open("mysql", generateDatasource())
	if err != nil {
		panic(err)
	}

	db := repository.ConnectDB(sqlDB)

	db.AutoMigrate(&model.Employee{})

	defaultRepository := repository.DefaultRepository{
		DB: db,
	}
	handler = controller.Controller{
		Repository: &defaultRepository,
	}

	echo := echo.New()
	echo.Use(middleware.Logger())

	echo.GET("/employees", handler.GetAllEmployees)
	echo.POST("/employee", handler.RegisterEmployee)
	echo.PUT("/employee/:id", handler.UpdateEmployee)

	return echo, handler
}

func beforeEach(t *testing.T) {
	t.Helper()
	sqlDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/")
	_, err = sqlDB.Exec("CREATE DATABASE " + DatabaseName)
	if err != nil {
		panic(err)
	}
	_, err = sqlDB.Exec("USE " + DatabaseName)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
}

func afterEach(t *testing.T) {
	t.Helper()
	sqlDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/")
	_, err = sqlDB.Exec("DROP SCHEMA " + DatabaseName)
	if err != nil {
		panic(err)
	}
}

func TestServer(t *testing.T) {
	t.Run("Persona can register a new employee and can see a list of employees", func(t *testing.T) {
		beforeEach(t)
		e, handler := run(t)

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
		assert.Equal(t, uint(1), returnedEmployees[0].Id)
		assert.Equal(t, "Jay", returnedEmployees[0].Name)
		assert.Equal(t, "100", returnedEmployees[0].Salary)
		assert.Equal(t, 30, returnedEmployees[0].Age)

		afterEach(t)
	})

	t.Run("Persona can update existing employee", func(t *testing.T) {
		beforeEach(t)
		e, handler := run(t)

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
		assert.Equal(t, uint(1), returnedEmployees[0].Id)
		assert.Equal(t, "Jay Kim", returnedEmployees[0].Name)
		assert.Equal(t, "1000", returnedEmployees[0].Salary)
		assert.Equal(t, 31, returnedEmployees[0].Age)

		afterEach(t)
	})
}
