package mysql_easy

import (
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	backoff "github.com/cenkalti/backoff"
	mysql "github.com/go-sql-driver/mysql"
)

func ConfiguraDB() (db *sql.DB, err error) {

	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
	)

	cfg := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASS"),
		Net:    "tcp",
		Addr:   addr,
		DBName: os.Getenv("MYSQL_DB"),
	}

	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = 20 * time.Second

	err = backoff.Retry(
		func() (err error) {
			db, err = sql.Open("mysql", cfg.FormatDSN())

			if err != nil {
				return
			}

			err = db.Ping()

			if err != nil {
				log.Println("Error connecting to MySQL...")
			}

			return
		},

		bf,
	)

	if err != nil {
		return
	}

	log.Println("Connected!")

	return
}
