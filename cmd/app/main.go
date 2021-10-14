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
	"net/http"
)

var (
	postgresDB = database.GetPostgresDB()
	vld        = validator.New()

	clientRepository = repository.GetClientRepository(postgresDB)
	clientService    = service.GetClientService(clientRepository, vld)
	clientController = controller.GetClientController(clientService)

	userRepository = repository.GetUserRepository(postgresDB)
	userService    = service.GetUserService(userRepository, vld)
	userController = controller.GetUserController(userService)
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.LoadHTMLFiles("html/login.html", "html/index.html", "html/client.html")

	panelRouter := router.Group("/api")
	{
		panelRouter.Static("/static/", "html/")

		panelRouter.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusUnauthorized, "login.html", nil)
		})
		panelRouter.POST("/home", userController.Login)

		panelRouter.GET("/clients", clientController.GetAllClients)
	}

	srv := new(server.Server)

	if err := srv.Run(viper.GetString("MAIN_PORT"), router); err != nil {
		panic(err)
	}
}
