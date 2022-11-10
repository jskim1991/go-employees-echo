package handlers

import (
	"employees-echo/models"
	"employees-echo/repository"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	t.Run("GetAllEmployees returns status ok", func(t *testing.T) {
		e := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/employees", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		employee := models.Employee{Id: "199", Name: "Jay", Salary: "100", Age: "30"}
		stubEmployees := []models.Employee{employee}
		h := Handler{Repository: &repository.SpyRepository{
			FindAll_returnValue: stubEmployees,
		}}

		h.GetAllEmployees(c)

		assert.Equal(t, http.StatusOK, response.Code)

		var returnedEmployees []models.Employee
		json.Unmarshal([]byte(response.Body.String()), &returnedEmployees)
		returnedEmployee := returnedEmployees[0]
		assert.Equal(t, "199", returnedEmployee.Id)
		assert.Equal(t, "Jay", returnedEmployee.Name)
		assert.Equal(t, "100", returnedEmployee.Salary)
		assert.Equal(t, "30", returnedEmployee.Age)
	})
}
