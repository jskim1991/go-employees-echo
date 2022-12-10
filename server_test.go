package main

import (
	"bytes"
	"database/sql"
	"employees-echo/controller"
	"employees-echo/dto"
	"employees-echo/model"
	"employees-echo/repository"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ServerTestSuite struct {
	suite.Suite
	TestDatabaseName string
}

func (suite *ServerTestSuite) SetupTest() {
	suite.TestDatabaseName = "journey"
}

func (suite *ServerTestSuite) SetupSubTest() {
	sqlDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/")

	_, err = sqlDB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s", suite.TestDatabaseName))
	if err != nil {
		panic(err)
	}

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", suite.TestDatabaseName))
	if err != nil {
		panic(err)
	}
	_, err = sqlDB.Exec(fmt.Sprintf("USE %s", suite.TestDatabaseName))
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
}

func (suite *ServerTestSuite) TearDownSubTest() {
	sqlDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/")
	_, err = sqlDB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s", suite.TestDatabaseName))
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) TestJourneys() {
	suite.Run("Persona can register a new employee and can see a list of employees", func() {
		e, handler := suite.runMain()

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
			suite.T().Error(err)
		}
		assert.Equal(suite.T(), "1", response.Body.String())

		request = httptest.NewRequest(http.MethodGet, "/employees", nil)
		response = httptest.NewRecorder()
		c = e.NewContext(request, response)

		err = handler.GetAllEmployees(c)
		if err != nil {
			suite.T().Error(err)
		}

		var returnedEmployees []dto.EmployeeResponse
		b, err := io.ReadAll(response.Body)
		if err != nil {
			suite.T().Error(err)
		}
		json.Unmarshal(b, &returnedEmployees)
		assert.Equal(suite.T(), 1, len(returnedEmployees))
		assert.Equal(suite.T(), uint(1), returnedEmployees[0].Id)
		assert.Equal(suite.T(), "Jay", returnedEmployees[0].Name)
		assert.Equal(suite.T(), "100", returnedEmployees[0].Salary)
		assert.Equal(suite.T(), 30, returnedEmployees[0].Age)
	})

	suite.Run("Persona can update existing employee", func() {
		e, handler := suite.runMain()

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
			suite.T().Error(err)
		}

		request = httptest.NewRequest(http.MethodGet, "/employees", nil)
		response = httptest.NewRecorder()
		c = e.NewContext(request, response)

		err = handler.GetAllEmployees(c)
		if err != nil {
			suite.T().Error(err)
		}

		var returnedEmployees []dto.EmployeeResponse
		b, err := io.ReadAll(response.Body)
		if err != nil {
			suite.T().Error(err)
		}
		json.Unmarshal(b, &returnedEmployees)
		assert.Equal(suite.T(), 1, len(returnedEmployees))
		assert.Equal(suite.T(), uint(1), returnedEmployees[0].Id)
		assert.Equal(suite.T(), "Jay Kim", returnedEmployees[0].Name)
		assert.Equal(suite.T(), "1000", returnedEmployees[0].Salary)
		assert.Equal(suite.T(), 31, returnedEmployees[0].Age)
	})
}

func (suite *ServerTestSuite) runMain() (e *echo.Echo, handler controller.Controller) {
	sqlDB, err := sql.Open("mysql", generateDatasource(suite.TestDatabaseName))
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
