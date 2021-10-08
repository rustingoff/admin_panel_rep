package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

func GetRedisConnection() *redis.Client {

	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var (
		_        = godotenv.Load("./database.env")
		addr     = viper.GetString("REDIS_HOST")
		port     = viper.GetString("REDIS_PORT")
		password = os.Getenv("REDIS_PASSWORD")
	)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr + port,
		Password: password,
		DB:       0,
	})

	_, err = rdb.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}

	log.Println("Redis successfully connected !")
	return rdb
}
