package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"github.com/rustingoff/admin_panel_rep/internal/service"
	"github.com/rustingoff/admin_panel_rep/pkg/hash"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UserController interface {
	CreateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	Login(c *gin.Context)
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
		return
	}

	c.JSON(http.StatusOK, "deleted")
}

func (cc *userController) Login(c *gin.Context) {
	var u model.UserLogin
	if err := c.ShouldBindJSON(&u); err != nil {
		log.Printf("invalid json provided, error: %v", err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	user, err := cc.cService.GetUserByUsername(u.Username)
	if err != nil {
		log.Printf("user not found, error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid credentials")
		return
	}

	//todo check user credentials...
	ok := hash.CheckPasswordHash(u.Password, user.Password)
	if !ok {
		log.Printf("invalid password\n")
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid credentials")
		return
	}

	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
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

func CreateToken(userid uint) (string, error) {
	var err error
	//Creating Access Token
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	_ = godotenv.Load("./database.env")

	acs := os.Getenv("ACCESS_SECRET")

	fmt.Println(acs)

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(acs))
	if err != nil {
		return "", err
	}
	return token, nil
}
