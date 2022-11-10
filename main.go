package main

import (
	"employees-echo/handlers"
	"employees-echo/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	defaultRepository := repository.DefaultRepository{}
	handler := handlers.Handler{
		Repository: &defaultRepository,
	}

	e := echo.New()

	e.GET("/employees", handler.GetAllEmployees)

	e.Start(":8080")
}
