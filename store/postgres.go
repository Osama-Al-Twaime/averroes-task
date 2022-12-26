package store

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"

	_ "github.com/lib/pq"
)

var (
	dbClient     *sql.DB
	dbClientOnce sync.Once
)

func GetPostgres() *sql.DB {
	if dbClient == nil {
		dbClientOnce.Do(func() {
			port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
			connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				os.Getenv("DB_HOST"),
				port,
				os.Getenv("DB_USER"),
				os.Getenv("DB_PASSWORD"),
				os.Getenv("DB_NAME"),
			)

			db, err := sql.Open("postgres", connectionString)

			if err != nil {
				fmt.Println("Error connecting to db")
			}

			dbClient = db

			// check db
			err = db.Ping()
			if err != nil {
				fmt.Println("Error pinging the db")
			}
		})
	}

	return dbClient
}
