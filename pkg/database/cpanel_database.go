package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func GetPostgresDB() *sqlx.DB {
	var (
		_        = godotenv.Load("./database.env")
		host     = viper.GetString("postgres.host")
		port     = viper.GetString("postgres.port")
		user     = viper.GetString("postgres.user")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = viper.GetString("postgres.db")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Panic(err)
		}
	}()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Postgres succesfully connected !")

	return db
}
