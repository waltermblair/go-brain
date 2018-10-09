package brain

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

type DBClient interface {
	FetchConfig(int) ([]int, string)
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
func (c *DBClientImpl) FetchConfig(id int) ([]int, string) {

	var row ConfigRecord
	var nextKeys []int

	query := fmt.Sprintf(
		"SELECT c.function, k.next FROM configurations c " +
		"JOIN next_keys k ON c.this = k.this " +
		"WHERE c.this = %d", id)

	log.Println("executing query: ", query)

	rows, err := c.DB.Query(query)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {

		if err := rows.Scan(&row.Function, &row.NextKey); err != nil {
			log.Fatal(err)
		}
		fmt.Println("APPENDING: " + strconv.Itoa(row.NextKey))
		nextKeys = append(nextKeys, row.NextKey)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("retrieved from db config for routing key: ", id)
	fmt.Println("YARRRR: " + row.Function)
	fmt.Println("YARRRR: " + strconv.Itoa(len(nextKeys)))
	return nextKeys, row.Function

}

func (c *DBClientImpl) HealthCheck() error {
	return c.DB.Ping()
}