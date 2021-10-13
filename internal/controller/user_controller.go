package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"github.com/rustingoff/admin_panel_rep/internal/service"
	"github.com/rustingoff/admin_panel_rep/pkg/hash"
	"github.com/rustingoff/admin_panel_rep/pkg/jwt"
	_ "github.com/spf13/viper"
	_ "github.com/twinj/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserController interface {
	CreateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
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

	user.Password, err = hash.HashPassword(user.Password)
	if err != nil {
		log.Printf("FAILED to hash a password with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "invalid password")
		return
	}

	tokenAuth, err := jwt.ExtractTokenMetadata(c.Request)
	if err != nil {
		log.Printf("FAILED to extract token metadata with error: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userId, err := jwt.FetchAuth(tokenAuth)
	if err != nil {
		log.Printf("FAILED to fetch auth with error: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	err = cc.cService.CreateUser(user)
	if err != nil {
		log.Printf("FAILED create user with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "server can't create a user")
		return
	}

	c.JSON(http.StatusCreated, userId)
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
		return
	}

	c.JSON(http.StatusOK, "deleted")
}

func (cc *userController) Login(c *gin.Context) {
	var u = model.UserLogin{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}

	user, err := cc.cService.GetUserByUsername(u.Username)
	if err != nil {
		log.Printf("user not found, error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid credentials")
		return
	}

	ok := hash.CheckPasswordHash(u.Password, user.Password)
	if !ok {
		log.Printf("invalid password\n")
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid credentials")
		return
	}

	token, err := jwt.CreateToken(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := jwt.CreateAuth(user.ID, token)
	if saveErr != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, saveErr.Error())
		return
	}
	fmt.Println(time.Now(), token.AtExpires)
	c.SetCookie("Authorization", token.AccessToken, 8*60*60, "", "localhost", false, true)
	c.JSON(http.StatusOK, token.AccessToken)
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

func (cc *userController) Logout(c *gin.Context) {
	au, err := jwt.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := jwt.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}
