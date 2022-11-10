package repository

import (
	"database/sql"
	"employees-echo/models"
	"log"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository interface {
	FindAll() []models.Employee
}

type DefaultRepository struct {
	DB *sql.DB
}

func ConnectDB() *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5432 dbname=test user=postgres password=")
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
