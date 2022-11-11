package repository

import (
	"context"
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
	Update(id int, e models.Employee)
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var employees []models.Employee

	query := `select id, name, salary, age from employee`
	rows, err := m.DB.QueryContext(ctx, query)
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newId int
	query := `insert into employee (name, salary, age, created_at, updated_at) values ($1, $2, $3, $4, $5) returning id`
	err := m.DB.QueryRowContext(ctx, query, e.Name, e.Salary, e.Age, time.Now(), time.Now()).Scan(&newId)
	if err != nil {
		log.Println(err)
	}

	return newId
}

func (m *DefaultRepository) Update(id int, e models.Employee) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update employee set name=$1, salary=$2, age=$3, updated_at=$4 where id = $5`
	_, err := m.DB.ExecContext(ctx, query, e.Name, e.Salary, e.Age, time.Now(), id)
	if err != nil {
		log.Println(err)
	}
}
