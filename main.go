package main

import (
	"employees-echo/controller"
	"employees-echo/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	datasource := "host=localhost port=5432 dbname=test user=postgres password="
	db := repository.ConnectDB(datasource)
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
	e.PUT("/employee/:id", handler.UpdateEmployee)

	e.Logger.Fatal(e.Start(":8080"))
}
