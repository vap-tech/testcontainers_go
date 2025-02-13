package rates

import (
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	"log"
	"time"
)

type Repository struct {
	conn *sql.DB
}

var data = [4]Rates{
	{
		day:   time.Date(2025, time.January, 23, 0, 0, 0, 0, time.UTC),
		value: 98.2804,
	},
	{
		day:   time.Date(2025, time.January, 24, 0, 0, 0, 0, time.UTC),
		value: 99.0978,
	},
	{
		day:   time.Date(2025, time.January, 25, 0, 0, 0, 0, time.UTC),
		value: 98.2636,
	},
	{
		day:   time.Date(2025, time.January, 28, 0, 0, 0, 0, time.UTC),
		value: 97.132,
	},
}

func NewRepository(connStr string) (*Repository, error) {
	conn, err := sql.Open("oracle", connStr)
	if err != nil {
		log.Printf("database connection error: %v", err)
		return nil, err
	}

	sqlText := "CREATE TABLE rates (day DATE NOT NULL, value FLOAT, PRIMARY KEY (day))"
	_, err = conn.Exec(sqlText)
	if err != nil {
		log.Print(err)
	}

	for _, rate := range data { // набиваем табличку данными
		stmt := fmt.Sprintf("INSERT INTO rates (day, value) VALUES (TO_DATE('%s','yyyy.mm.dd'), '%.4f')", rate.day.Format("2006-01-02"), rate.value)
		_, err = conn.Exec(stmt)
		if err != nil {
			log.Fatalf("table rates inserting error: %v", err)
			return nil, err
		}
	}

	return &Repository{
		conn: conn,
	}, nil
}

func (r Repository) GetRates(d time.Time) (rate Rates, err error) {

	d = d.Add(-24 * time.Hour) // Минус один день от целевой даты, т.к. курс актуален только на следующий день

	stmt := fmt.Sprintf("SELECT day, value FROM rates WHERE day <= TO_DATE('%s','yyyy.mm.dd') ORDER BY day DESC FETCH FIRST 1 ROW ONLY", d.Format("2006-01-02"))
	err = r.conn.QueryRow(stmt).Scan(&rate.day, &rate.value)
	if err != nil {
		return Rates{}, err
	}

	log.Printf("date params: %s", d)
	log.Printf("date for query: %v", d.Add(-24*time.Hour))
	log.Printf("result: %v, %v", rate.day, rate.value)

	return rate, nil
}
