package models

type EmployeeRequest struct {
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Age    int    `json:"age"`
}
