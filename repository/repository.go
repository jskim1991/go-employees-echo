package repository

import (
	"database/sql"
	"employees-echo/models"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type Repository interface {
	FindAll() []models.Employee
	InsertEmployee(e models.Employee) int
}

type DefaultRepository struct {
	DB *sql.DB
}

func ConnectDB(datasource string) *sql.DB {
	db, err := sql.Open("pgx", datasource)
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}
	return db
}

func (m *DefaultRepository) FindAll() []models.Employee {
	var employees []models.Employee

	query := `select id, name, salary, age from employee`
	rows, err := m.DB.Query(query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var e models.Employee
		err := rows.Scan(&e.Id, &e.Name, &e.Salary, &e.Age)
		if err != nil {
			log.Println(err)
		}
		employees = append(employees, e)
	}

	return employees
}

func (m *DefaultRepository) InsertEmployee(e models.Employee) int {
	var newId int
	query := `insert into employee (name, salary, age, created_at, updated_at) values ($1, $2, $3, $4, $5) returning id`
	err := m.DB.QueryRow(query, e.Name, e.Salary, e.Age, time.Now(), time.Now()).Scan(&newId)
	if err != nil {
		log.Println(err)
	}

	return newId
}
