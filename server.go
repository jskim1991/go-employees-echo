package main

import (
	"database/sql"
	"employees-echo/controller"
	"employees-echo/model"
	"employees-echo/repository"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var DatabaseName string

func main() {
	DatabaseName = "employees-echo"
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.ParseTime = true
	cfg.User = "root"
	cfg.Addr = fmt.Sprintf("%s:%s", "127.0.0.1", "3307")
	cfg.DBName = DatabaseName
	datasource := cfg.FormatDSN()
	sqlDB, err := sql.Open("mysql", datasource)
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
