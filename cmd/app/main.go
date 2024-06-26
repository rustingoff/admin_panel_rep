package main

import (
	"github.com/gin-gonic/gin"
	server "github.com/rustingoff/admin_panel_rep"
	"github.com/rustingoff/admin_panel_rep/internal/controller"
	"github.com/rustingoff/admin_panel_rep/internal/middleware"
	"github.com/rustingoff/admin_panel_rep/internal/repository"
	"github.com/rustingoff/admin_panel_rep/internal/service"
	"github.com/rustingoff/admin_panel_rep/pkg/database"
	"github.com/rustingoff/admin_panel_rep/pkg/jwt"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
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

	// generate pdf from client contract

	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(gin.Logger())
	//router.Use(gin.Recovery())
	router.LoadHTMLFiles("html/login.html", "html/index.html", "html/client.html", "html/client_update.html")

	panelRouter := router.Group("/api")
	{
		panelRouter.Static("/static/", "html/")

		panelRouter.GET("/", func(c *gin.Context) {
			err := jwt.TokenValid(c.Request)
			if err != nil {
				c.HTML(http.StatusUnauthorized, "login.html", nil)
			} else {
				c.Redirect(http.StatusPermanentRedirect, "home/")
			}
		})

		panelRouter.POST("/", func(c *gin.Context) {
			c.HTML(http.StatusUnauthorized, "login.html", nil)
		})

		panelRouter.POST("/home", userController.Login)

		panelRouter.GET("/home", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
		panelRouter.POST("/clients", middleware.TokenAuthMiddleware(), clientController.CreateClient)
		panelRouter.GET("/clients", middleware.TokenAuthMiddleware(), clientController.GetAllClients)
		//panelRouter.GET("/clients/:client_id", middleware.TokenAuthMiddleware(), clientController.GetClient)

		panelRouter.POST("/clients/update/", middleware.TokenAuthMiddleware(), clientController.UpdateClient)
		panelRouter.GET("/clients/update/", middleware.TokenAuthMiddleware(), func(c *gin.Context) {
			idQ := c.Query("id")
			id, err := strconv.Atoi(idQ)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, "/clients")
			}
			c.HTML(http.StatusOK, "client_update.html", id)
		})

		panelRouter.GET("/clients/delete/", middleware.TokenAuthMiddleware(), clientController.DeleteClient)

		panelRouter.GET("/log_out", middleware.TokenAuthMiddleware(), userController.Logout)
	}

	srv := new(server.Server)

	if err := srv.Run(viper.GetString("MAIN_PORT"), router); err != nil {
		panic(err)
	}
}
