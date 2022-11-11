package main

import (
	"employees-echo/controller"
	"employees-echo/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := repository.ConnectDB()
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

	e.Logger.Fatal(e.Start(":8080"))
}
