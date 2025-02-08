package customer

import (
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	"log"
)

type Repository struct {
	conn *sql.DB
}

func NewRepository(connStr string) (*Repository, error) {
	conn, err := sql.Open("oracle", connStr)
	if err != nil {
		log.Printf("database connection error: %v", err)
		return nil, err
	}

	sqlText := "CREATE TABLE customers (id number(10) NOT NULL, name varchar2(50) NOT NULL, email varchar2(50), PRIMARY KEY (id))"
	_, err = conn.Exec(sqlText)
	if err != nil {
		log.Print(err)
	}
	sqlText = "INSERT INTO customers (id, name, email) VALUES (0, 'Anna', 'a.petrenko@gmail.com')"
	_, err = conn.Exec(sqlText)
	if err != nil {
		log.Print(err)
	}

	return &Repository{
		conn: conn,
	}, nil
}

func (r Repository) CreateCustomer(customer Customer) (Customer, error) {

	sqlText := fmt.Sprintf("INSERT INTO customers (id, name, email) VALUES (%d, %s, %s)", customer.Id, customer.Name, customer.Email)
	_, _ = r.conn.Exec(sqlText)

	sqlText = fmt.Sprintf("SELECT id, name, email FROM customers WHERE id = %d", customer.Id)
	rows, _ := r.conn.Query(sqlText)

	var err error
	for rows.Next() {
		err = rows.Scan(&customer.Id, &customer.Name, &customer.Email)
	}

	return customer, err
}

func (r Repository) GetCustomerByEmail(email string) (Customer, error) {
	var customer Customer
	sqlText := fmt.Sprintf("SELECT id, name, email FROM customers WHERE email = '%s'", email)
	err := r.conn.QueryRow(sqlText).Scan(&customer.Id, &customer.Name, &customer.Email)
	if err != nil {
		log.Print(err)
	}

	return customer, err
}
