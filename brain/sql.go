package brain

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBClient interface {
	FetchConfig(int) ([]int, string, error)
	HealthCheck() error
}

type DBClientImpl struct {
	DBURL string
	DB    *sql.DB
}

func NewDBClient(url string) (client DBClient, err error) {

	db, err := sql.Open("mysql", url)

	client = &DBClientImpl{
		DBURL: url,
		DB:    db,
	}

	return client, err
}

// TODO - check in schema from db volume
// Returns config for the component associated with given routing-key
func (c *DBClientImpl) FetchConfig(id int) ([]int, string, error) {

	var rows		*sql.Rows
	var row 		ConfigRecord
	var nextKeys 	[]int
	var err			error

	query := fmt.Sprintf(
		"SELECT c.function, k.next FROM configurations c " +
		"JOIN next_keys k ON c.this = k.this " +
		"WHERE c.this = %d", id)

	log.Println("executing query: ", query)

	rows, err = c.DB.Query(query)

	if err != nil {
		log.Println("error fetching config from database")
		return nextKeys, row.Function, err
	} else {
		defer rows.Close()
	}

	for rows.Next() {

		if err := rows.Scan(&row.Function, &row.NextKey); err != nil {
			log.Println("error unmarshalling query result row")
			return nextKeys, row.Function, err
		}
		nextKeys = append(nextKeys, row.NextKey)
	}
	if err = rows.Err(); err != nil {
		log.Println("error parsing query result rows")
		return nextKeys, row.Function, err
	}

	log.Println("retrieved from db config for routing key: ", id)

	return nextKeys, row.Function, err

}

func (c *DBClientImpl) HealthCheck() error {
	return c.DB.Ping()
}