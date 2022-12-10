package controller

import (
	"bytes"
	"employees-echo/dto"
	"employees-echo/model"
	"employees-echo/testdoubles"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ControllerTestSuite struct {
	suite.Suite
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (suite *ControllerTestSuite) TestGetAllEmployees() {
	suite.Run("returns 200 ok", func() {
		context, rec := NewEchoContext(http.MethodGet, "/employees", nil)
		controller := Controller{Repository: &testdoubles.SpyStubRepository{
			FindAllReturnValue: []model.Employee{},
		}}

		controller.GetAllEmployees(context)

		assert.Equal(suite.T(), http.StatusOK, rec.Code)
	})

	suite.Run("invokes Repository::FindAll() [test doubles]", func() {
		context, _ := NewEchoContext(http.MethodGet, "/employees", nil)
		spyStubRepository := testdoubles.SpyStubRepository{
			FindAllReturnValue: []model.Employee{},
		}
		controller := Controller{Repository: &spyStubRepository}

		controller.GetAllEmployees(context)

		assert.Equal(suite.T(), 1, spyStubRepository.FindAllInvocation)
	})

	suite.Run("invokes Repository::FindAll() [testify]", func() {
		context, _ := NewEchoContext(http.MethodGet, "/employees", nil)
		mockRepository := testdoubles.MockRepository{}
		controller := Controller{Repository: &mockRepository}
		mockRepository.On("FindAll", mock.Anything).Return([]model.Employee{}, nil)

		controller.GetAllEmployees(context)

		mockRepository.AssertExpectations(suite.T())
	})

	suite.Run("returns employee slice [test doubles]", func() {
		context, rec := NewEchoContext(http.MethodGet, "/employees", nil)
		employee := model.Employee{
			Model:  gorm.Model{ID: 199},
			Name:   "Jay",
			Salary: "100",
			Age:    30,
		}
		stubEmployees := []model.Employee{employee}
		controller := Controller{Repository: &testdoubles.SpyStubRepository{
			FindAllReturnValue: stubEmployees,
		}}

		controller.GetAllEmployees(context)

		var returnedEmployees []dto.EmployeeResponse
		responseBody, _ := io.ReadAll(rec.Body)
		json.Unmarshal(responseBody, &returnedEmployees)
		assert.Equal(suite.T(), 1, len(returnedEmployees))
		returnedEmployee := returnedEmployees[0]
		assert.Equal(suite.T(), uint(199), returnedEmployee.Id)
		assert.Equal(suite.T(), "Jay", returnedEmployee.Name)
		assert.Equal(suite.T(), "100", returnedEmployee.Salary)
		assert.Equal(suite.T(), 30, returnedEmployee.Age)
	})

	suite.Run("returns employee slice [testify]", func() {
		context, rec := NewEchoContext(http.MethodGet, "/employees", nil)

		mockRepository := testdoubles.MockRepository{}
		controller := Controller{Repository: &mockRepository}

		employee := model.Employee{
			Model:  gorm.Model{ID: 199},
			Name:   "Jay",
			Salary: "100",
			Age:    30,
		}
		mockRepository.On("FindAll", mock.Anything).Return([]model.Employee{employee}, nil)
		controller.GetAllEmployees(context)

		var returnedEmployees []dto.EmployeeResponse
		responseBody, _ := io.ReadAll(rec.Body)
		json.Unmarshal(responseBody, &returnedEmployees)
		assert.Equal(suite.T(), 1, len(returnedEmployees))
		returnedEmployee := returnedEmployees[0]
		assert.Equal(suite.T(), uint(199), returnedEmployee.Id)
		assert.Equal(suite.T(), "Jay", returnedEmployee.Name)
		assert.Equal(suite.T(), "100", returnedEmployee.Salary)
		assert.Equal(suite.T(), 30, returnedEmployee.Age)
	})
}

func (suite *ControllerTestSuite) TestRegisterEmployee() {
	suite.Run("returns status created", func() {
		context, rec := NewEchoContext(http.MethodPost, "/employee", nil)
		controller := Controller{Repository: &testdoubles.SpyStubRepository{}}

		controller.RegisterEmployee(context)

		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
	})

	suite.Run("invokes Repository::InsertEmployee with given request body [test doubles]", func() {
		employeeRequest := dto.EmployeeRequest{Name: "Jay", Salary: "100", Age: 30}
		postJson, _ := json.Marshal(employeeRequest)
		context, _ := NewEchoContext(http.MethodPost, "/employee", bytes.NewReader(postJson))

		spyStubRepository := testdoubles.SpyStubRepository{
			InsertEmployeeReturnValue: uint(999),
		}
		controller := Controller{Repository: &spyStubRepository}

		controller.RegisterEmployee(context)

		assert.Equal(suite.T(), 1, spyStubRepository.InsertEmployeeInvocation)
		assert.Equal(suite.T(), "Jay", spyStubRepository.InsertEmployeeArgument.Name)
		assert.Equal(suite.T(), "100", spyStubRepository.InsertEmployeeArgument.Salary)
		assert.Equal(suite.T(), 30, spyStubRepository.InsertEmployeeArgument.Age)
	})

	suite.Run("invokes Repository::InsertEmployee with given request body [testify]", func() {
		employeeRequest := dto.EmployeeRequest{Name: "Jay", Salary: "100", Age: 30}
		postJson, _ := json.Marshal(employeeRequest)
		context, _ := NewEchoContext(http.MethodPost, "/employee", bytes.NewReader(postJson))

		mockRepository := testdoubles.MockRepository{}
		mockRepository.On("InsertEmployee", employeeRequest).Return(uint(999), nil)
		controller := Controller{Repository: &mockRepository}

		controller.RegisterEmployee(context)

		mockRepository.AssertExpectations(suite.T())
	})

	suite.Run("returns id of the registered employee [test doubles]", func() {
		employeeRequest := dto.EmployeeRequest{Name: "Jay", Salary: "100", Age: 30}
		postJson, _ := json.Marshal(employeeRequest)
		context, rec := NewEchoContext(http.MethodPost, "/employee", bytes.NewReader(postJson))

		spyStubRepository := testdoubles.SpyStubRepository{
			InsertEmployeeReturnValue: uint(999),
		}
		controller := Controller{Repository: &spyStubRepository}

		controller.RegisterEmployee(context)

		assert.Equal(suite.T(), "999", rec.Body.String())
	})

	suite.Run("returns id of the registered employee [testify]", func() {
		employeeRequest := dto.EmployeeRequest{Name: "Jay", Salary: "100", Age: 30}
		postJson, _ := json.Marshal(employeeRequest)
		context, rec := NewEchoContext(http.MethodPost, "/employee", bytes.NewReader(postJson))

		mockRepository := testdoubles.MockRepository{}
		mockRepository.On("InsertEmployee", employeeRequest).Return(uint(999), nil)
		controller := Controller{Repository: &mockRepository}

		controller.RegisterEmployee(context)

		assert.Equal(suite.T(), "999", rec.Body.String())
	})
}

func (suite *ControllerTestSuite) TestUpdateEmployee() {
	suite.Run("returns status 204 no content", func() {
		employee := dto.EmployeeRequest{Name: "Sam", Salary: "1", Age: 40}
		postJson, _ := json.Marshal(employee)
		context, rec := NewEchoContext(http.MethodPut, "/employee", bytes.NewReader(postJson))
		context.SetPath("/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		controller := Controller{Repository: &testdoubles.SpyStubRepository{}}

		controller.UpdateEmployee(context)

		assert.Equal(suite.T(), http.StatusNoContent, rec.Code)
	})

	suite.Run("invokes Repository::Update with given id and employee request body [test doubles]", func() {
		employee := dto.EmployeeRequest{Name: "Sam", Salary: "1", Age: 40}
		postJson, _ := json.Marshal(employee)
		context, _ := NewEchoContext(http.MethodPut, "/employee", bytes.NewReader(postJson))
		context.SetPath("/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		spyStubRepository := testdoubles.SpyStubRepository{}
		controller := Controller{Repository: &spyStubRepository}

		controller.UpdateEmployee(context)

		assert.Equal(suite.T(), 1, spyStubRepository.UpdateInvocation)
		assert.Equal(suite.T(), uint(1), spyStubRepository.UpdateArgumentId)
		assert.Equal(suite.T(), "Sam", spyStubRepository.UpdateArgumentEmployee.Name)
		assert.Equal(suite.T(), "1", spyStubRepository.UpdateArgumentEmployee.Salary)
		assert.Equal(suite.T(), 40, spyStubRepository.UpdateArgumentEmployee.Age)
	})

	suite.Run("invokes Repository::Update with given id and employee request body [testify]", func() {
		employee := dto.EmployeeRequest{Name: "Sam", Salary: "1", Age: 40}
		postJson, _ := json.Marshal(employee)
		context, _ := NewEchoContext(http.MethodPut, "/employee", bytes.NewReader(postJson))
		context.SetPath("/:id")
		context.SetParamNames("id")
		context.SetParamValues("199")

		mockRepository := testdoubles.MockRepository{}
		mockRepository.On("Update", uint(199), employee).Return(model.Employee{}, nil)
		controller := Controller{Repository: &mockRepository}

		controller.UpdateEmployee(context)

		mockRepository.AssertExpectations(suite.T())
	})
}

func NewEchoContext(method, url string, requestBody io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	request := httptest.NewRequest(method, url, requestBody)
	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	return c, recorder
}
