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
		"SELECT c.this, c.function, k.next FROM `configurations` c " +
		"JOIN `next_keys` k ON c.this=k.this " +
		"WHERE c.this=$1")

	rows, err := c.db.Query(query, id)
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&row); err != nil {
			log.Fatal(err)
		}
		nextKeys = append(nextKeys, row.NextKey)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return row.ID, nextKeys, row.Function

}

func (c *DBClientImpl) HealthCheck() error {
	return c.db.Ping()
}