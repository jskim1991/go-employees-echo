package main

import (
	"database/sql"
	"employees-echo/controller"
	"employees-echo/model"
	"employees-echo/repository"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var DatabaseName string

func main() {
	DatabaseName = "employees-echo"
	sqlDB, err := sql.Open("mysql", generateDatasource())
	if err != nil {
		panic(err)
	}

	db := repository.ConnectDB(sqlDB)
	db.AutoMigrate(&model.Employee{})

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

func generateDatasource() string {
	options := "charset=utf8&parseTime=True"
	return fmt.Sprintf("root:@tcp(127.0.0.1:3307)/%s?%s", DatabaseName, options)
}
