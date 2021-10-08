package main

import (
	"github.com/gin-gonic/gin"
	server "github.com/rustingoff/admin_panel_rep"
	"github.com/rustingoff/admin_panel_rep/internal/controller"
	"github.com/rustingoff/admin_panel_rep/internal/repository"
	"github.com/rustingoff/admin_panel_rep/internal/service"
	"github.com/rustingoff/admin_panel_rep/pkg/database"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

var (
	postgresDB = database.GetPostgresDB()
	vld        = validator.New()

	clientRepository = repository.GetClientRepository(postgresDB)
	clientService    = service.GetClientService(clientRepository, vld)
	clientController = controller.GetClientController(clientService)
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

	panelRouter := router.Group("/cpanel")
	{
		//GET todo
		panelRouter.GET("/", clientController.GetAllClients)
		panelRouter.GET("/:clientID", clientController.GetClient)
		// panelRouter.GET()

		//CREATE todo
		panelRouter.POST("/", clientController.CreateClient)

		//UPDATE todo
		panelRouter.PATCH("/:clientID", clientController.UpdateClient)

		//DELETE todo
		panelRouter.DELETE("/:clientID", clientController.DeleteClient)
	}

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("MAIN_PORT"), router); err != nil {
		panic(err)
	}
}
