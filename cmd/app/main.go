package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rustingoff/admin_panel_rep/internal/database"
	"github.com/spf13/viper"
)

func main() {

	router := gin.New()
	router.Use(gin.Recovery())

	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.GetString("MAIN_PORT"))

	panelRouter := router.Group("/cpanel")
	{
		//GET todo
		panelRouter.GET("/")
		// panelRouter.GET()

		//CREATE todo
		panelRouter.POST("/")

		//UPDATE todo
		panelRouter.PATCH("/")
		panelRouter.PUT("/")

		//DELETE todo
		panelRouter.DELETE("/")
	}

	_ = database.GetPostgresDB()

	router.Run(viper.GetString("MAIN_PORT"))
}
