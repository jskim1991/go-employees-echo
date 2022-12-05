package dto

type EmployeeRequest struct {
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Age    int    `json:"age"`
}

type EmployeeResponse struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Age    int    `json:"age"`
}
