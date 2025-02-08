package customer

import (
	"fmt"
	"os"

	"database/sql"
	_ "github.com/sijms/go-ora/v2"
)

type Repository struct {
	conn *sql.DB
}

func NewRepository(connStr string) (*Repository, error) {
	conn, err := sql.Open("oracle", connStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "database connection error: %v\n", err)
		return nil, err
	}
	return &Repository{
		conn: conn,
	}, nil
}

func (r Repository) CreateCustomer(customer Customer) (Customer, error) {
	err := r.conn.QueryRow(
		"INSERT INTO customers (name, email) VALUES ($1, $2) RETURNING id",
		customer.Name, customer.Email).Scan(&customer.Id)
	return customer, err
}

func (r Repository) GetCustomerByEmail(email string) (Customer, error) {
	var customer Customer
	query := "SELECT id, name, email FROM customers WHERE email = $1"
	err := r.conn.QueryRow(query, email).
		Scan(&customer.Id, &customer.Name, &customer.Email)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}
