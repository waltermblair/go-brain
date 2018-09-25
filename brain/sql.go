package brain

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBClient interface {
	FetchConfig(int) (int, []int, string)
	HealthCheck() error
}

type DBClientImpl struct {
	dbURL		string
	db 			*sql.DB
}

func NewDBClient(url string) (client DBClient, err error) {

	db, err := sql.Open("mysql", url)

	client = &DBClientImpl{
		dbURL: url,
		db: db,
	}

	return client, err
}

// TODO - check in schema from db volume
// TODO - test
func (c *DBClientImpl) FetchConfig(id int) (int, []int, string) {

	var row ConfigRecord
	var nextKeys []int

	query := fmt.Sprintf(
		"SELECT c.this, c.function, k.next FROM configurations c " +
		"JOIN next_keys k ON c.this = k.this " +
		"WHERE c.this = %d", id)

	log.Println("executing query: ", query)

	rows, err := c.db.Query(query)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {

		if err := rows.Scan(&row.ID, &row.Function, &row.NextKey); err != nil {
			log.Fatal(err)
		}
		nextKeys = append(nextKeys, row.NextKey)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("retrieved from db config for id: ", row.ID)
	return row.ID, nextKeys, row.Function

}

func (c *DBClientImpl) HealthCheck() error {
	return c.db.Ping()
}