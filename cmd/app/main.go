package main

import (
	"github.com/gin-gonic/gin"
	server "github.com/rustingoff/admin_panel_rep"
	"github.com/rustingoff/admin_panel_rep/internal/controller"
	"github.com/rustingoff/admin_panel_rep/internal/repository"
	"github.com/rustingoff/admin_panel_rep/internal/service"
	"github.com/rustingoff/admin_panel_rep/pkg/database"
	"github.com/rustingoff/admin_panel_rep/pkg/redis"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

var (
	postgresDB = database.GetPostgresDB()
	vld        = validator.New()
	_          = redis.GetRedisConnection()

	clientRepository = repository.GetClientRepository(postgresDB)
	clientService    = service.GetClientService(clientRepository, vld)
	clientController = controller.GetClientController(clientService)

	userRepository = repository.GetUserRepository(postgresDB)
	userService    = service.GetUserService(userRepository, vld)
	userController = controller.GetUserController(userService)
)

func main() {
	router := gin.Default()

	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})

	panelRouter := router.Group("/api")
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

		adminRouter := panelRouter.Group("/cmd")
		{
			adminRouter.POST("/login", userController.Login)
			adminRouter.POST("/", userController.CreateUser)
		}
	}

	srv := new(server.Server)

	if err := srv.Run(viper.GetString("MAIN_PORT"), router); err != nil {
		panic(err)
	}
}
