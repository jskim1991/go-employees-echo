package controller

import (
	"bytes"
	"employees-echo/dto"
	"employees-echo/testdoubles"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	t.Run("GetAllEmployees invokes Repository::FindAll()", func(t *testing.T) {
		e := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/employees", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		spyStubRepository := testdoubles.SpyStubRepository{
			FindAll_returnValue: []dto.EmployeeResponse{},
		}
		controller := Controller{Repository: &spyStubRepository}

		controller.GetAllEmployees(c)

		assert.Equal(t, 1, spyStubRepository.FindAll_invocation)
	})

	t.Run("GetAllEmployees returns status ok with employee slice", func(t *testing.T) {
		e := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/employees", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		employee := dto.EmployeeResponse{Id: 199, Name: "Jay", Salary: "100", Age: 30}
		stubEmployees := []dto.EmployeeResponse{employee}
		controller := Controller{Repository: &testdoubles.SpyStubRepository{
			FindAll_returnValue: stubEmployees,
		}}

		controller.GetAllEmployees(c)

		assert.Equal(t, http.StatusOK, response.Code)
		var returnedEmployees []dto.EmployeeResponse
		json.Unmarshal([]byte(response.Body.String()), &returnedEmployees)
		returnedEmployee := returnedEmployees[0]
		assert.Equal(t, 199, returnedEmployee.Id)
		assert.Equal(t, "Jay", returnedEmployee.Name)
		assert.Equal(t, "100", returnedEmployee.Salary)
		assert.Equal(t, 30, returnedEmployee.Age)
	})

	t.Run("RegisterEmployee returns status created", func(t *testing.T) {
		e := echo.New()
		employee := dto.EmployeeRequest{Name: "Jay", Salary: "100", Age: 30}
		postJson, _ := json.Marshal(employee)
		request := httptest.NewRequest(http.MethodPost, "/employee", bytes.NewReader(postJson))
		request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		controller := Controller{Repository: &testdoubles.SpyStubRepository{}}
		controller.RegisterEmployee(c)

		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("RegistrationEmployee returns id of the registered employee", func(t *testing.T) {
		e := echo.New()
		employee := dto.EmployeeRequest{Name: "Jay", Salary: "100", Age: 30}
		postJson, _ := json.Marshal(employee)
		request := httptest.NewRequest(http.MethodPost, "/employee", bytes.NewReader(postJson))
		request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		controller := Controller{Repository: &testdoubles.SpyStubRepository{
			InsertEmployee_returnValue: 199,
		}}
		controller.RegisterEmployee(c)

		assert.Equal(t, "199", response.Body.String())
	})

	t.Run("RegisterEmployee invokes Repository::InsertEmployee() with given employee", func(t *testing.T) {
		e := echo.New()
		employee := dto.EmployeeRequest{Name: "Jay", Salary: "100", Age: 30}
		postJson, _ := json.Marshal(employee)
		request := httptest.NewRequest(http.MethodPost, "/employee", bytes.NewReader(postJson))
		request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		spyStubRepository := &testdoubles.SpyStubRepository{}
		controller := Controller{Repository: spyStubRepository}
		controller.RegisterEmployee(c)

		assert.Equal(t, 1, spyStubRepository.InsertEmployee_invocation)
		assert.Equal(t, employee.Name, spyStubRepository.InsertEmployee_argument.Name)
		assert.Equal(t, employee.Salary, spyStubRepository.InsertEmployee_argument.Salary)
		assert.Equal(t, employee.Age, spyStubRepository.InsertEmployee_argument.Age)
	})

	t.Run("UpdateEmployee invokes Repository::Update() with given id and employee info", func(t *testing.T) {
		e := echo.New()
		employee := dto.EmployeeRequest{Name: "Sam", Salary: "1", Age: 40}
		postJson, _ := json.Marshal(employee)
		request := httptest.NewRequest(http.MethodPut, "/employee", bytes.NewReader(postJson))
		request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		spyStubRepository := &testdoubles.SpyStubRepository{}
		controller := Controller{Repository: spyStubRepository}
		controller.UpdateEmployee(c)

		assert.Equal(t, 1, spyStubRepository.Update_invocation)
		assert.Equal(t, 1, spyStubRepository.Update_argument_id)
		assert.Equal(t, "Sam", spyStubRepository.Update_argument_employee.Name)
		assert.Equal(t, "1", spyStubRepository.Update_argument_employee.Salary)
		assert.Equal(t, 40, spyStubRepository.Update_argument_employee.Age)
	})

	t.Run("UpdateEmployee returns status ok with no content", func(t *testing.T) {
		e := echo.New()
		employee := dto.EmployeeRequest{Name: "Sam", Salary: "1", Age: 40}
		postJson, _ := json.Marshal(employee)
		request := httptest.NewRequest(http.MethodPut, "/employee", bytes.NewReader(postJson))
		request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		spyStubRepository := &testdoubles.SpyStubRepository{}
		controller := Controller{Repository: spyStubRepository}
		controller.UpdateEmployee(c)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 0, response.Body.Len())
	})
}
