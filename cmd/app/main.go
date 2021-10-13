package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	server "github.com/rustingoff/admin_panel_rep"
	"github.com/rustingoff/admin_panel_rep/internal/controller"
	"github.com/rustingoff/admin_panel_rep/internal/middleware"
	"github.com/rustingoff/admin_panel_rep/internal/repository"
	"github.com/rustingoff/admin_panel_rep/internal/service"
	"github.com/rustingoff/admin_panel_rep/pkg/database"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
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

	fmt.Println(time.Now().Format(time.RFC3339))

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
			c.HTML(http.StatusOK, "login.html", nil)
		})

		panelRouter.POST("/home", userController.Login, clientController.GetAllClients)

		panelRouter.GET("/client", middleware.TokenAuthMiddleware(), clientController.GetAllClients)
		panelRouter.GET("/client/:clientID", middleware.TokenAuthMiddleware(), clientController.GetClient)
		panelRouter.POST("/client", middleware.TokenAuthMiddleware(), clientController.CreateClient)
		panelRouter.PATCH("/client/:clientID", middleware.TokenAuthMiddleware(), clientController.UpdateClient)
		panelRouter.DELETE("/client/:clientID", middleware.TokenAuthMiddleware(), clientController.DeleteClient)

		adminRouter := panelRouter.Group("/cmd")
		{
			adminRouter.POST("/", middleware.TokenAuthMiddleware(), userController.CreateUser)
			adminRouter.GET("/logout", middleware.TokenAuthMiddleware(), userController.Logout)
		}
	}

	srv := new(server.Server)

	if err := srv.Run(viper.GetString("MAIN_PORT"), router); err != nil {
		panic(err)
	}
}
