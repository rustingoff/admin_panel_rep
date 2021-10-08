package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"github.com/rustingoff/admin_panel_rep/internal/service"
	"log"
	"net/http"
	"strconv"
)

type UserController interface {
	CreateUser(c *gin.Context)
	DeleteUser(c *gin.Context)

	GetAllUsers(c *gin.Context)
	GetUser(c *gin.Context)
}

type userController struct {
	cService service.UserService
}

func GetUserController(s service.UserService) UserController {
	return &userController{cService: s}
}

func (cc *userController) CreateUser(c *gin.Context) {
	var user model.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Printf("FAILED bind json to user structure with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid request")
		return
	}

	err = cc.cService.CreateUser(user)
	if err != nil {
		log.Printf("FAILED create user with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "server can't create a user")
		return
	}

	c.JSON(http.StatusCreated, "created")
}

func (cc *userController) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil || userID < 1 {
		log.Printf("FAILED convert param to int with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid url")
		return
	}

	err = cc.cService.DeleteUser(uint(userID))
	if err != nil {
		log.Printf("FAILED delete useruser with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "server can't delete a user")
	}

	c.JSON(http.StatusOK, "deleted")
}

func (cc *userController) GetAllUsers(c *gin.Context) {

	response, err := cc.cService.GetAllUsers()
	if err != nil {
		log.Printf("FAILED to get users with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "invalid server response")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (cc *userController) GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil || userID < 1 {
		log.Printf("FAILED convert param to int with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid url")
		return
	}

	response, err := cc.cService.GetUser(uint(userID))
	if err != nil {
		log.Printf("FAILED to get user with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "invalid server response")
		return
	}

	c.JSON(http.StatusOK, response)
}
