package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init() {

	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load("./database.env")
}
