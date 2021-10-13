package main

import (
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
	"text/template"
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

	panelRouter := router.Group("/api")
	{
		panelRouter.Static("/static/", "html/")
		tmpl := template.Must(template.ParseFiles("html/login.html"))
		panelRouter.GET("/", func(c *gin.Context) {
			err := tmpl.Execute(c.Writer, nil)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, "can't execute a template")
				return
			}
		})

		panelRouter.POST("/login", userController.Login)

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
