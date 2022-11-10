package models

import "time"

type Employee struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Age    int    `json:"age"`
}

type EmployeeEntity struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Salary    string    `db:"salary"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
