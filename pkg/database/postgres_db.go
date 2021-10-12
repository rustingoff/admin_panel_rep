package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"

	"github.com/spf13/viper"
)

func GetPostgresDB() *gorm.DB {

	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

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

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	log.Printf("Postgres succesfully connected !")

	err = db.AutoMigrate(
		&model.Client{},
		&model.User{},
	)
	if err != nil {
		panic(err)
	}

	log.Println("Successful migration")

	return db
}
