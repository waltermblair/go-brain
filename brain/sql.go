package brain

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBClient interface {
	FetchNumInputs(int) (int, error)
	SetComponentStatus(int, string) error
	FetchConfig(Config) (int, []int, string, error)
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

// TODO factor out query and make this general FetchCount function
// FetchNumInputs returns number of inputs that this component should expect to receive in current configuration
func (c *DBClientImpl) FetchNumInputs(id int) (count int, err error) {

	var rows		*sql.Rows

	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM configurations c " +
		"JOIN next_keys k ON c.this = k.this " +
		"WHERE k.next = %d AND c.status != 'down'", id)

	log.Println("executing query: ", query)

	rows, err = c.DB.Query(query)

	if err != nil {
		log.Println("error fetching count from database")
		return count, err
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			log.Println("error unmarshalling query result row: ", err.Error())
			return count, err
		}
	}

	return count, err
}

// TODO factor out query and make this general Update function
func (c *DBClientImpl) SetComponentStatus(id int, status string) error {
	query := fmt.Sprintf(
		"UPDATE configurations c " +
		"SET status = '%s' " +
		"WHERE c.this = %d", status, id)
	_, err := c.DB.Exec(query)
	return err
}

// Returns config for the component associated with given routing-key
func (c *DBClientImpl) FetchConfig(config Config) (int, []int, string, error) {

	var rows		*sql.Rows
	var row 		ConfigRecord
	var numInputs	int
	var nextKeys 	[]int
	var err			error

	id := config.ID

	query := fmt.Sprintf(
		"SELECT c.function, k.next FROM configurations c " +
		"JOIN next_keys k ON c.this = k.this " +
		"WHERE c.this = %d", id)

	log.Println("executing query: ", query)

	rows, err = c.DB.Query(query)

	if err != nil {
		log.Println("error fetching config from database")
		return numInputs, nextKeys, row.Function, err
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		if err := rows.Scan(&row.Function, &row.NextKey); err != nil {
			log.Println("error unmarshalling query result row: ", err.Error())
			return numInputs, nextKeys, row.Function, err
		}
		nextKeys = append(nextKeys, row.NextKey)
	}

	log.Println("retrieved from db config for routing key: ", id)

	log.Println("updating status for this component in db...")
	err = c.SetComponentStatus(id, config.Status)
	if err != nil {
		log.Println("error updating status for this component: ", err.Error())
		return numInputs, nextKeys, row.Function, err
	}
	log.Printf("updated status for %d to %s", id, config.Status)

	log.Println("determining number of inputs...")
	numInputs, err = c.FetchNumInputs(id)
	if err != nil {
		log.Println("error determining number of inputs")
		return numInputs, nextKeys, row.Function, err
	}
	log.Printf("number of inputs for %d: %d", id, numInputs)

	return numInputs, nextKeys, row.Function, err

}

func (c *DBClientImpl) HealthCheck() error {
	return c.DB.Ping()
}